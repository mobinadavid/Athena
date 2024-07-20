package migrations

import (
	"athena/src/config"
	"athena/src/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

var blockchainMigration *migrate.Migrate
var walletMigration *migrate.Migrate
var explorerMigration *migrate.Migrate

func init() {
	config.Init()
	database.Init()

	db := database.GetInstance()
	driver, err := postgres.WithInstance(db.GetDB(), &postgres.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	blockchainMigration, err = migrate.NewWithDatabaseInstance(
		"file://src/services/blockchain/migration",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalln(err)
	}

	walletMigration, err = migrate.NewWithDatabaseInstance(
		"file://src/services/wallet-address/migration",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalln(err)
	}

	explorerMigration, err = migrate.NewWithDatabaseInstance(
		"file://src/services/blockchain-explorer/migration",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func Up() error {
	if err := blockchainMigration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err := walletMigration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err := explorerMigration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func Down() error {
	if err := blockchainMigration.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err := walletMigration.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err := explorerMigration.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
