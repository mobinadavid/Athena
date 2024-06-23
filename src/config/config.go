package config

import (
	"github.com/joho/godotenv"
	"sync"
)

// Config represents the application's configuration with thread-safe access.
type Config struct {
	data map[string]string
	mu   sync.RWMutex
}

var (
	instance *Config
	once     sync.Once
)

// GetInstance provides a thread-safe, lazy-loaded singleton instance of Config.
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{
			data: make(map[string]string),
		}
	})
	return instance
}

// Init initializes the configuration using .env file.
// It returns an error if loading the .env file fails.
func Init() error {
	return GetInstance().Load()
}

// Get retrieves the value for a given key, thread-safe.
func (c *Config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

// Set assigns a value to a key, thread-safe.
func (c *Config) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Load reads configurations from a .env file and stores them, thread-safe.
// Returns an error if reading the file fails.
func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	env, err := godotenv.Read(".env")
	if err != nil {
		return err
	}

	for key, value := range env {
		c.data[key] = value
	}

	return nil
}
