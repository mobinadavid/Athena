package migrations

import (
	"athena/src/config"
	"athena/src/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
)

var migration *migrate.Migrate

func init() {
	config.Init()
	database.Init()

	db := database.GetInstance()
	driver, err := postgres.WithInstance(db.GetDB(), &postgres.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://src/database/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalln(err)
	}

	migration = m
}

func Up() error {
	return migration.Up()
}

func Down() error {
	return migration.Down()
}
