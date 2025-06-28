package job

import (
	"context"
	"time"

	"eticket-api/internal/common/logger"
	"eticket-api/internal/usecase"

	"github.com/robfig/cron/v3"
)

type ClaimSessionJob struct {
	Log     logger.Logger
	Usecase *usecase.ClaimSessionUsecase
}

func NewClaimSessionJob(log logger.Logger, usecase *usecase.ClaimSessionUsecase) *ClaimSessionJob {
	return &ClaimSessionJob{Log: log, Usecase: usecase}
}

func (j *ClaimSessionJob) CleanExpiredClaimSession() {
	j.Log.Info("[ClaimSessionJob] Scheduler starting...")

	c := cron.New()
	c.AddFunc("@every 1h", func() {
		j.Log.Info("[ClaimSessionJob] Scheduled cleanup triggered")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := j.Usecase.DeleteExpiredClaimSession(ctx); err != nil {
			j.Log.WithError(err).Error("[ClaimSessionJob] Cleanup failed")
		} else {
			j.Log.Info("[ClaimSessionJob] Cleanup completed successfully")
		}
	})
	c.Start()

	// Optionally run once at startup
	go func() {
		j.Log.Info("[ClaimSessionJob] Initial cleanup triggered")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := j.Usecase.DeleteExpiredClaimSession(ctx); err != nil {
			j.Log.WithError(err).Error("[ClaimSessionJob] Initial cleanup failed")
		} else {
			j.Log.Info("[ClaimSessionJob] Initial cleanup completed successfully")
		}
	}()
}
