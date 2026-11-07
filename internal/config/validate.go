package config

import (
	"fmt"
	"strings"
	"time"
)

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.App.Name) == "" {
		return fmt.Errorf("app.name is required")
	}
	if cfg.App.LogLevel != "debug" && cfg.App.LogLevel != "info" && cfg.App.LogLevel != "warn" && cfg.App.LogLevel != "error" {
		return fmt.Errorf("app.logLevel must be debug, info, warn, or error")
	}
	if strings.TrimSpace(cfg.Database.URLEnv) == "" {
		return fmt.Errorf("database.urlEnv is required")
	}
	if cfg.Database.ConnectTimeoutSeconds < 1 || cfg.Database.ConnectTimeoutSeconds > 60 {
		return fmt.Errorf("database.connectTimeoutSeconds must be between 1 and 60")
	}
	if cfg.Schedule.EveryHours < 1 {
		return fmt.Errorf("schedule.everyHours must be at least 1")
	}
	if cfg.Schedule.MaxRunMinutes < 1 {
		return fmt.Errorf("schedule.maxRunMinutes must be at least 1")
	}
	if _, err := time.LoadLocation(cfg.Schedule.Timezone); err != nil {
		return fmt.Errorf("schedule.timezone %q is invalid: %w", cfg.Schedule.Timezone, err)
	}
	if cfg.BrowserWorker.Enabled {
		if strings.TrimSpace(cfg.BrowserWorker.Command) == "" {
			return fmt.Errorf("browserWorker.command is required when the worker is enabled")
		}
		if cfg.BrowserWorker.HealthTimeoutSeconds < 1 || cfg.BrowserWorker.HealthTimeoutSeconds > 60 {
			return fmt.Errorf("browserWorker.healthTimeoutSeconds must be between 1 and 60")
		}
	}
	if strings.TrimSpace(cfg.Reports.Directory) == "" {
		return fmt.Errorf("reports.directory is required")
	}

	seen := make(map[string]struct{})
	allSources := append(append([]SourceConfig{}, cfg.DirectCareerSources...), cfg.ThirdPartySources...)
	for index, source := range allSources {
		id := strings.TrimSpace(source.ID)
		if id == "" {
			return fmt.Errorf("source at index %d has no id", index)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("source id %q is duplicated", id)
		}
		seen[id] = struct{}{}
	}
	return nil
}
