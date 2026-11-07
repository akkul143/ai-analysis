package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/akkul143/ai-analysis/internal/app"
)

const version = "0.1.0-phase1"

func Execute(arguments []string, stdout io.Writer, stderr io.Writer) int {
	if len(arguments) == 0 {
		printUsage(stderr)
		return 2
	}

	switch arguments[0] {
	case "run":
		if err := executeRun(arguments[1:], stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "jobfinder: %v\n", err)
			return 1
		}
		return 0
	case "version", "--version", "-version":
		fmt.Fprintln(stdout, version)
		return 0
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", arguments[0])
		printUsage(stderr)
		return 2
	}
}

func executeRun(arguments []string, stdout io.Writer, stderr io.Writer) error {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	flags.SetOutput(stderr)
	dryRun := flags.Bool("dry-run", false, "validate the Phase 1 foundation and persist a run record")
	configPath := flags.String("config", "config/config.local.yaml", "path to application configuration")
	profilePath := flags.String("profile", "config/profile.yaml", "path to candidate profile")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if !*dryRun {
		return errors.New("real job collection is introduced in Phase 2; use `jobfinder run --dry-run`")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	application, err := app.New(ctx, app.Options{
		ConfigPath:  *configPath,
		ProfilePath: *profilePath,
		LogOutput:   stderr,
	})
	if err != nil {
		return err
	}
	defer application.Close()

	summary, err := application.DryRun(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, "Phase 1 dry run completed")
	fmt.Fprintf(stdout, "Run ID:      %s\n", summary.ID)
	fmt.Fprintf(stdout, "Status:      %s\n", summary.Status)
	fmt.Fprintf(stdout, "Started:     %s\n", summary.StartedAt.Format(time.RFC3339))
	fmt.Fprintf(stdout, "Completed:   %s\n", summary.CompletedAt.Format(time.RFC3339))
	fmt.Fprintln(stdout, "PostgreSQL:  connected")
	fmt.Fprintln(stdout, "Run lock:    acquired and released")
	fmt.Fprintln(stdout, "Run record:  persisted")
	return nil
}

func printUsage(output io.Writer) {
	fmt.Fprintln(output, `ai-analysis Phase 1

Usage:
  jobfinder run --dry-run [--config PATH] [--profile PATH]
  jobfinder version
  jobfinder help

Real collection, matching, reports, scheduling, and notifications are added in later phases.`)
}
