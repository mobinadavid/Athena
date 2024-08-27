package database

import (
	"athena/src/config"
	databaseDrivers "athena/src/database/drivers"
	"database/sql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
)

// Package-level variables to enforce singleton pattern for database connection and instance retrieval.
var (
	connectOnce     sync.Once // Ensures database connection is established only once.
	getInstanceOnce sync.Once // Ensures a single instance of Database is created.
	instance        *Database // Holds the singleton instance of Database.
)

// IDatabaseDriver defines the interface for database drivers.
// It specifies the methods required for a database driver to be compatible with the Database struct.
type IDatabaseDriver interface {
	Connect() error      // Connect establishes a connection to the database.
	Close() error        // Close terminates the connection to the database.
	GetClient() *gorm.DB // GetClient returns the underlying client for direct operations.
	GetDB() *sql.DB      // GetDB returns the underlying client for direct operations.
}

// Database encapsulates the database operations and driver.
// It serves as a central point for database interactions, leveraging a driver that implements the IDatabaseDriver interface.
type Database struct {
	driver IDatabaseDriver // The database driver, implementing IDatabaseDriver for database operations.
}

// Init initializes the database by establishing a connection.
// It retrieves the singleton instance of the Database and calls Connect on it.
func Init() (err error) {
	return GetInstance().Connect()
}

// Connect establishes a connection to the database if not already connected.
// It uses connectOnce to ensure that the database connection is established only once,
// preventing multiple connections in a concurrent environment.
func (database *Database) Connect() (err error) {
	connectOnce.Do(func() {
		configs := config.GetInstance()                            // Retrieve configurations
		dbPort, _ := strconv.Atoi(configs.Get("DB_INTERNAL_PORT")) // Convert port to int
		// Initialize the driver with configuration values
		database.driver = &databaseDrivers.Postgres{
			Username: configs.Get("DB_USERNAME"),
			Password: configs.Get("DB_PASSWORD"),
			Host:     configs.Get("DB_HOST"),
			Port:     dbPort,
			Database: configs.Get("DB_DATABASE"),
			SSLMode:  configs.Get("DB_SSL_MODE"),
			Timezone: configs.Get("APP_TZ"),
		}
		// Establish connection
		err = database.driver.Connect()
	})
	return
}

// Close terminates the database connection.
// It delegates the close operation to the database driver and logs the closure.
func (database *Database) Close() (err error) {
	err = database.driver.Close()
	log.Println("Database Service: Disconnected Successfully.")
	return
}

// GetClient retrieves the gorm.DB client from the database driver.
// It allows for direct database operations using the ORM.
func (database *Database) GetClient() *gorm.DB {
	return database.driver.GetClient()
}

// GetDB retrieves the sql.DB client from the database driver.
// It allows for direct database operations using the ORM.
func (database *Database) GetDB() *sql.DB {
	return database.driver.GetDB()
}

// GetInstance returns the singleton instance of the Database.
// It ensures that only one instance of Database is created and used throughout the application,
// leveraging getInstanceOnce to enforce this constraint.
func GetInstance() *Database {
	getInstanceOnce.Do(func() {
		instance = &Database{} // Initialize the singleton instance
	})
	return instance
}
