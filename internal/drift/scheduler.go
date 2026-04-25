package drift

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/user/driftwatch/internal/config"
)

// ScheduledJob holds the configuration and state for a recurring drift check.
type ScheduledJob struct {
	ServiceName string
	Config      *config.ServiceConfig
	Interval    time.Duration
}

// Scheduler runs drift checks on a configurable interval for one or more services.
// It uses a Watcher internally to perform each check and records results via HistoryStore.
type Scheduler struct {
	watcher *Watcher
	history *HistoryStore
	alerts  *AlertManager
	jobs    []ScheduledJob
	mu      sync.Mutex
}

// NewScheduler creates a Scheduler with the provided Watcher, HistoryStore, and AlertManager.
func NewScheduler(w *Watcher, h *HistoryStore, a *AlertManager) *Scheduler {
	return &Scheduler{
		watcher: w,
		history: h,
		alerts:  a,
	}
}

// Register adds a service configuration to the scheduler with the given check interval.
// Duplicate service names are allowed but will result in multiple independent jobs.
func (s *Scheduler) Register(name string, cfg *config.ServiceConfig, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, ScheduledJob{
		ServiceName: name,
		Config:      cfg,
		Interval:    interval,
	})
}

// Start launches a goroutine for each registered job and blocks until the context is cancelled.
// Each goroutine ticks at the job's interval, runs a drift check, records history, and fires alerts.
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	jobs := make([]ScheduledJob, len(s.jobs))
	copy(jobs, s.jobs)
	s.mu.Unlock()

	if len(jobs) == 0 {
		log.Println("scheduler: no jobs registered, exiting")
		return
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j ScheduledJob) {
			defer wg.Done()
			s.runJob(ctx, j)
		}(job)
	}
	wg.Wait()
}

// runJob executes drift checks for a single job on its configured interval until ctx is cancelled.
func (s *Scheduler) runJob(ctx context.Context, job ScheduledJob) {
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	log.Printf("scheduler: starting job for service %q every %s", job.ServiceName, job.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("scheduler: stopping job for service %q", job.ServiceName)
			return
		case t := <-ticker.C:
			s.executeCheck(job, t)
		}
	}
}

// executeCheck performs a single drift check, stores the result in history, and emits alerts.
func (s *Scheduler) executeCheck(job ScheduledJob, at time.Time) {
	if job.Config == nil {
		log.Printf("scheduler: skipping check for %q — nil config", job.ServiceName)
		return
	}

	result := s.watcher.Check(job.Config)

	entry := HistoryEntry{
		Service:   job.ServiceName,
		Timestamp: at,
		Drifted:   len(result.Diffs) > 0,
		Diffs:     result.Diffs,
	}

	if err := s.history.Append(job.ServiceName, entry); err != nil {
		log.Printf("scheduler: failed to record history for %q: %v", job.ServiceName, err)
	}

	if s.alerts != nil {
		if summary := s.alerts.Evaluate(result.Diffs); summary.Level != AlertLevelNone {
			log.Printf("scheduler: alert [%s] for service %q — %s",
				summary.Level, job.ServiceName, SummaryLine(summary))
		}
	}
}
