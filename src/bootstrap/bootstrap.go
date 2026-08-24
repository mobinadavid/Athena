package bootstrap

import (
	"athena/src/api"
	"athena/src/cache"
	"athena/src/database"
	"athena/src/pkg/i18n"
	"athena/src/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Init() (err error) {
	defer func() {
		log.Println("Goodbye!")
		os.Exit(0)
	}()

	err = logger.Init()
	if err != nil {
		log.Fatalf("Logger Service: Failed to Initialize. %v", err)
	}
	logger.GetInstance().Info("Initialized Successfully.", zap.String("Service", "Logger"), zap.Time("timestamp", time.Now()))

	// Initialize i18n
	err = i18n.Init()
	if err != nil {
		log.Fatalf("I18n Service: Failed to Initialize. %v", err)
	}
	log.Println("I18n Service: Initialized Successfully.")

	// Initialize database
	err = database.Init()
	if err != nil {
		log.Fatalf("Database Service: Failed to Initialize. %v", err)
	}
	log.Println("Database Service: Initialized Successfully.")

	defer func() {
		if err = database.GetInstance().Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	// Initialize Cache
	err = cache.Init()
	if err != nil {
		logger.GetInstance().Fatal("Failed to Initialize", zap.String("Service", "Cache"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}
	logger.GetInstance().Info("Initialized Successfully.", zap.String("Service", "Cache"), zap.Time("timestamp", time.Now()))

	defer func() {
		if err = cache.GetInstance().Close(); err != nil {
			logger.GetInstance().Fatal("Failed to close cache connection:", zap.String("Service", "Cache"), zap.Error(err), zap.Time("timestamp", time.Now()))
		}
	}()

	// Initialize API
	go func() {
		err = api.Init()
		if err != nil {
			log.Fatalf("API Service: Failed to Initialize. %v", err)
		}
		log.Println("API Service: Initialized Successfully.")
	}()

	log.Println("Application is now running.\nPress CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Application is shutting down...")

	return
}
