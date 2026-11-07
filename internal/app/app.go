package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/akkul143/ai-analysis/internal/browser"
	"github.com/akkul143/ai-analysis/internal/config"
	"github.com/akkul143/ai-analysis/internal/observability"
	"github.com/akkul143/ai-analysis/internal/profile"
	runservice "github.com/akkul143/ai-analysis/internal/run"
	"github.com/akkul143/ai-analysis/internal/store/postgres"
)

type Options struct {
	ConfigPath string
	ProfilePath string
	LogOutput  io.Writer
}

type Application struct {
	config        config.Config
	profile       profile.Profile
	store         *postgres.Store
	runs          *runservice.Service
	browserWorker *browser.Client
	logger        *slog.Logger
}

func New(ctx context.Context, options Options) (*Application, error) {
	cfg, err := config.Load(options.ConfigPath)
	if err != nil {
		return nil, err
	}
	logger, err := observability.New(cfg.App.LogLevel, options.LogOutput)
	if err != nil {
		return nil, err
	}
	candidate, err := profile.Load(options.ProfilePath)
	if err != nil {
		return nil, err
	}
	databaseURL := os.Getenv(cfg.Database.URLEnv)
	if databaseURL == "" {
		return nil, fmt.Errorf("environment variable %s is required", cfg.Database.URLEnv)
	}
	store, err := postgres.Open(ctx, databaseURL, time.Duration(cfg.Database.ConnectTimeoutSeconds)*time.Second)
	if err != nil {
		return nil, err
	}
	if err := store.Migrate(ctx); err != nil {
		store.Close()
		return nil, err
	}

	return &Application{
		config:        cfg,
		profile:       candidate,
		store:         store,
		runs:          runservice.NewService(store, logger),
		browserWorker: browser.NewClient(cfg.BrowserWorker),
		logger:        logger,
	}, nil
}

func (application *Application) Close() {
	application.store.Close()
}

func (application *Application) DryRun(ctx context.Context) (runservice.Summary, error) {
	metadata := map[string]any{
		"candidate_name":       application.profile.Identity.Name,
		"preferred_titles":     len(application.profile.Titles.Preferred),
		"direct_sources":       len(application.config.DirectCareerSources),
		"third_party_sources":  len(application.config.ThirdPartySources),
		"browser_worker_enabled": application.browserWorker.Enabled(),
	}
	return application.runs.Execute(ctx, "dry_run", metadata, func(ctx context.Context, runID string) error {
		application.logger.InfoContext(ctx, "configuration loaded", "run_id", runID, "app", application.config.App.Name)
		application.logger.InfoContext(ctx, "candidate profile loaded", "run_id", runID, "candidate", application.profile.Identity.Name)
		if application.browserWorker.Enabled() {
			response, err := application.browserWorker.Health(ctx)
			if err != nil {
				return err
			}
			application.logger.InfoContext(ctx, "browser worker available", "run_id", runID, "worker", response.Result["worker"])
		}
		return nil
	})
}
