package bootstrap

import (
	"athena/src/api"
	"athena/src/cache"
	"athena/src/database"
	"athena/src/pkg/i18n"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Init() (err error) {
	defer func() {
		log.Println("Goodbye!")
		os.Exit(0)
	}()

	// Initialize i18n
	err = i18n.Init()
	if err != nil {
		log.Fatalf("I18n Service: Failed to Initialize. %v", err)
	}
	log.Println("I18n Service: Initialized Successfully.")

	// Initialize cache
	err = cache.Init()
	if err != nil {
		log.Fatalf("Cache Service: Failed to Initialize. %v", err)
	}
	log.Println("Cache Service: Initialized Successfully.")

	defer func() {
		if err = cache.GetInstance().Close(); err != nil {
			log.Fatalf("Failed to close cache connection: %v", err)
		}
	}()

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
