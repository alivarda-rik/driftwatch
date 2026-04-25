// Package main is the entry point for the driftwatch CLI tool.
// It wires together configuration loading, drift detection, alerting,
// and reporting into a cohesive command-line interface.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/drift"
)

const (
	defaultSnapshotDir = ".driftwatch/snapshots"
	defaultBaselineDir = ".driftwatch/baselines"
	defaultHistoryDir  = ".driftwatch/history"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("driftwatch", flag.ContinueOnError)

	configFile := fs.String("config", "", "Path to the declared configuration file (YAML or TOML)")
	snapshotDir := fs.String("snapshot-dir", defaultSnapshotDir, "Directory for storing snapshots")
	baselineDir := fs.String("baseline-dir", defaultBaselineDir, "Directory for storing baselines")
	historyDir := fs.String("history-dir", defaultHistoryDir, "Directory for storing drift history")
	watchMode := fs.Bool("watch", false, "Enable continuous watch mode")
	watchInterval := fs.Duration("interval", 30*time.Second, "Poll interval for watch mode")
	captureBaseline := fs.Bool("capture-baseline", false, "Capture current state as baseline and exit")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	if *configFile == "" {
		fs.Usage()
		return fmt.Errorf("--config is required")
	}

	// Load declared configuration from file.
	cfg, err := config.Load(*configFile)
	if err != nil {
		return fmt.Errorf("loading config %q: %w", *configFile, err)
	}

	// Initialise stores.
	snapshotStore, err := drift.NewSnapshotStore(*snapshotDir)
	if err != nil {
		return fmt.Errorf("initialising snapshot store: %w", err)
	}

	baselineStore, err := drift.NewBaselineStore(*baselineDir)
	if err != nil {
		return fmt.Errorf("initialising baseline store: %w", err)
	}

	historyStore, err := drift.NewHistoryStore(*historyDir)
	if err != nil {
		return fmt.Errorf("initialising history store: %w", err)
	}

	// Build the core components.
	detector := drift.NewDetector(cfg)
	reporter := drift.NewReporter(os.Stdout)
	alerter := drift.NewAlerter(drift.DefaultAlertRules())
	alertManager := drift.NewAlertManager(detector, alerter)
	baselineManager := drift.NewBaselineManager(baselineStore, detector)

	_ = snapshotStore
	_ = historyStore

	// Capture baseline mode: record the current live state and exit.
	if *captureBaseline {
		if err := baselineManager.Capture(cfg.Name); err != nil {
			return fmt.Errorf("capturing baseline for %q: %w", cfg.Name, err)
		}
		fmt.Printf("Baseline captured for service %q\n", cfg.Name)
		return nil
	}

	if *watchMode {
		// Continuous watch: emit alerts on every tick.
		watcher := drift.NewWatcher(detector, *watchInterval)
		results := watcher.Watch()
		fmt.Printf("Watching %q every %s — press Ctrl+C to stop\n", cfg.Name, *watchInterval)
		for report := range results {
			alert := alertManager.Evaluate(report)
			if alert != nil {
				fmt.Println(drift.SummaryLine(alert))
			}
			reporter.Report(report)
		}
		return nil
	}

	// Single-shot check.
	report, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("detecting drift: %w", err)
	}

	alert := alertManager.Evaluate(report)
	if alert != nil {
		fmt.Println(drift.SummaryLine(alert))
	}
	reporter.Report(report)

	if report != nil && len(report.Diffs) > 0 {
		// Exit with a non-zero code so CI pipelines can react to drift.
		os.Exit(2)
	}
	return nil
}
