package migrations

import (
	"athena/src/config"
	"athena/src/database"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"
)

//go:embed *.sql
var migrationFS embed.FS

var (
	migration *migrate.Migrate
	db        = database.GetInstance()
)

func init() {
	config.Init()
	database.Init()

	driver, err := postgres.WithInstance(db.GetDB(), &postgres.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	source, err := iofs.New(migrationFS, ".")
	if err != nil {
		log.Fatalf("Migration service error:%v", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalf("Migration service error:%v", err)
	}

	migration = m
}

func Up() error {
	return migration.Up()
}

func Down() error {
	return migration.Down()
}
