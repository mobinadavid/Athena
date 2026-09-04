package jobs

import (
	"athena/src/config"
	"athena/src/pkg/logger"
	"athena/src/providers"
	"strings"
	"time"

	"go.uber.org/zap"
)

func StartPaymentTracker() error {
	if !trackerEnabled() {
		logger.GetInstance().Info("payment tracker is disabled")
		return nil
	}

	interval := trackerInterval()
	logger.GetInstance().Info("payment tracker started", zap.Duration("interval", interval))

	runOnce()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		runOnce()
	}
	return nil
}

func runOnce() {
	if err := providers.GetContainer().DepositService.HandleDeposits(); err != nil {
		logger.GetInstance().Warn("payment tracker run failed", zap.Error(err))
	}
}

func trackerEnabled() bool {
	raw := strings.ToLower(strings.TrimSpace(config.GetInstance().Get("TRACKER_ENABLED")))
	return raw == "" || raw == "true" || raw == "1"
}

func trackerInterval() time.Duration {
	if raw := config.GetInstance().Get("TRACKER_INTERVAL"); raw != "" {
		if duration, err := time.ParseDuration(raw); err == nil && duration >= 30*time.Second {
			return duration
		}
	}
	if raw := config.GetInstance().Get("TRACKER_INTERVAL_SECONDS"); raw != "" {
		if duration, err := time.ParseDuration(raw + "s"); err == nil && duration >= 30*time.Second {
			return duration
		}
	}
	return 2 * time.Minute
}
