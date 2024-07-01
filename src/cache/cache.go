package cache

import (
	cacheDrivers "athena/src/cache/drivers"
	"athena/src/config"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

var (
	connectOnce     sync.Once // Ensures cache connection is established only once.
	getInstanceOnce sync.Once // Ensures a single instance of Cache is created.
	instance        *Cache    // Holds the singleton instance of Cache.
)

// ICacheDriver defines the interface for cache drivers.
// It specifies the methods required for a cache driver to be compatible with the Cache struct.
type ICacheDriver interface {
	Connect() error           // Connect establishes a connection to the cache.
	Close() error             // Close terminates the connection to the cache.
	GetClient() *redis.Client // GetClient returns the underlying client for direct cache operations.
}

// Cache encapsulates the cache operations and driver.
// It serves as a central point for cache interactions, leveraging a driver that implements the ICacheDriver interface.
type Cache struct {
	driver ICacheDriver // The cache driver, implementing ICacheDriver for cache operations.
}

// Init initializes the cache by establishing a connection.
// It retrieves the singleton instance of the Cache and calls Connect on it.
func Init() (err error) {
	return GetInstance().Connect()
}

// Connect establishes a connection to the cache if not already connected.
// It uses connectOnce to ensure that the cache connection is established only once,
// preventing multiple connections in a concurrent environment.
func (cache *Cache) Connect() (err error) {
	connectOnce.Do(func() {
		configs := config.GetInstance() // Retrieve configurations
		// Initialize the driver with configuration values
		cache.driver = &cacheDrivers.Redis{
			Host:     configs.Get("REDIS_HOST"),
			Port:     configs.Get("REDIS_PORT"),
			Password: configs.Get("REDIS_PASSWORD"),
		}

		if err = cache.driver.Connect(); err != nil {
			log.Fatalln(err) // Log and halt on error during connection setup
		}
	})

	return
}

// Close terminates the cache connection.
// It delegates the close operation to the cache driver and logs the closure.
func (cache *Cache) Close() (err error) {
	if err = cache.driver.Close(); err != nil {
		log.Println("Error closing cache connection:", err)
	} else {
		log.Println("Cache Service: Disconnected Successfully.")
	}
	return
}

// GetClient retrieves the redis.Client from the cache driver.
// It allows for direct cache operations using the client.
func (cache *Cache) GetClient() *redis.Client {
	return cache.driver.GetClient()
}

// GetInstance returns the singleton instance of the Cache.
// It ensures that only one instance of Cache is created and used throughout the application,
// leveraging getInstanceOnce to enforce this constraint.
func GetInstance() *Cache {
	getInstanceOnce.Do(func() {
		instance = &Cache{} // Initialize the singleton instance
	})
	return instance
}
