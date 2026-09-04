package api

import (
	"athena/src/jobs"
	"athena/src/pkg/logger"

	"go.uber.org/zap"
)

func startPaymentTracker() {
	go func() {
		if err := jobs.StartPaymentTracker(); err != nil {
			logger.GetInstance().Error("payment tracker stopped", zap.Error(err))
		}
	}()
}
