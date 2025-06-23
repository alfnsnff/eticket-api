package runner

import (
	"context"
	"eticket-api/internal/job"
	"log"

	"github.com/robfig/cron/v3"
)

type CleanupRunner struct {
	Job *job.CleanupJob
}

func NewCleanupRunner(job *job.CleanupJob) *CleanupRunner {
	return &CleanupRunner{
		Job: job,
	}
}

func (r *CleanupRunner) Start() {
	c := cron.New()

	_, err := c.AddFunc("@every 1h", func() {
		ctx := context.Background()
		err := r.Job.Run(ctx)
		if err != nil {
			log.Printf("[Job] Cleanup failed: %v", err)
		} else {
			log.Println("[Job] Cleanup completed successfully")
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule CleanupJob: %v", err)
	}

	log.Println("CleanupJob scheduled every hour")

	go func() {
		log.Println("Running initial CleanupJob now")
		ctx := context.Background()
		if err := r.Job.Run(ctx); err != nil {
			log.Printf("Initial Cleanup failed: %v", err)
		} else {
			log.Println("Initial Cleanup completed successfully")
		}
	}()

	c.Start()
}
