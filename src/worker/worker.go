package worker

import (
	"athena/src/providers"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

func Init(ctx context.Context) error {
	serviceContainer := providers.GetContainer()
	cronScheduler := cron.New()

	// Schedule the service method to run every 30 seconds
	_, err := cronScheduler.AddFunc("@every 30s", func() {
		err := serviceContainer.DepositService.HandleDeposits()
		if err != nil {
			log.Printf("Error handling new deposits: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule HandleDeposit: %w", err)
	}

	// Start the cron scheduler in a separate goroutine
	cronScheduler.Start()
	// Listen for context cancellation to stop the scheduler gracefully
	go func() {
		<-ctx.Done()
		log.Println("Shutting down cron scheduler...")
		cronScheduler.Stop()
	}()

	return nil
}
