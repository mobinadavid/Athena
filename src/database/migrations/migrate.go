package migrations

import (
	"athena/src/config"
	"athena/src/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

var migrations []*migrate.Migrate

func init() {

	services := []string{
		"blockchain",
		"wallet-address",
		"blockchain-explorer",
	}

	for _, service := range services {
		config.Init()
		database.Init()

		db := database.GetInstance()
		driver, err := postgres.WithInstance(db.GetDB(), &postgres.Config{})
		if err != nil {
			log.Fatalf("Failed to create database driver: %v", err)
		}

		migration, err := migrate.NewWithDatabaseInstance(
			"file://src/services/"+service+"/migration",
			"postgres",
			driver,
		)
		if err != nil {
			log.Fatalf("Failed to create migration instance for %s: %v", service, err)
		}
		migrations = append(migrations, migration)
	}
}

func Up() error {
	for _, m := range migrations {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Failed to apply migration: %v", err)
			return err
		}
	}
	return nil
}

func Down() error {
	for _, m := range migrations {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Failed to revert migration: %v", err)
			return err
		}
	}
	return nil
}
