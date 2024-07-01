// Package hash provides a flexible hashing system that supports multiple hashing algorithms.
// It defines interfaces and implementations for generating and verifying hashes, utilizing a driver-based approach.
package hash

import (
	"athena/src/config"
	"athena/src/hash/drivers"
	"fmt"
	"strings"
)

// IHash defines an interface for hash operations, including generating and verifying hashes.
// This allows for the implementation of various hashing algorithms with a consistent interface.
type IHash interface {
	Generate(str []byte) ([]byte, error)        // Generate creates a hash of the given byte slice.
	Verify(hashedStr, str string) (bool, error) // Verify checks if the provided string matches the hashed string.
}

// Hash encapsulates hashing operations using a specific hashing algorithm driver.
// It serves as the main entry point for hashing functionality in the application.
type Hash struct {
	driver     IHash  // The driver implementing the hashing algorithm.
	driverName string // The name of the driver.
}

// Generate generates a hash of the input byte slice using the configured hash driver.
// Generated result: <driver-name>:<hash>
func (hash *Hash) Generate(str []byte) ([]byte, error) {
	hashedStr, err := hash.driver.Generate(str)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s:%s", hash.driverName, hashedStr)), nil
}

// Verify checks if the provided hash corresponds to the given input string,
// using the configured hash driver.
func (hash *Hash) Verify(hashedStr, str string) (bool, error) {
	return hash.driver.Verify(hashedStr, str)
}

var (
	instance *Hash // Singleton instance of Hash.
)

// GetInstance returns an instance of Hash.
// It optionally takes a driver name and configures the Hash instance to use that driver.
// If no driver is specified, it uses the default driver or one specified by environment variables.
func GetInstance(driverNameArg ...string) *Hash {
	driverName := getDriverName(driverNameArg...)
	driver, _ := hashFactory(driverName)

	instance = &Hash{
		driver:     driver,
		driverName: driverName,
	}

	return instance
}

// getDriverName determines the hash driver name to use based on the input and environment configuration.
// It prioritizes the user input, falls back to an environment variable, or defaults to "bcrypt" if neither is specified.
func getDriverName(args ...string) string {
	if len(args) > 0 && args[0] != "" {
		return args[0]
	}

	configs := config.GetInstance()
	envDriver := configs.Get("HASH_DRIVER")
	if envDriver != "" {
		return envDriver
	}

	return "argon2" // Default to Argon2 if no driver is specified.
}

// hashFactory returns an IHash instance based on the driver name.
// It allows for the dynamic selection of hashing algorithms based on configuration or application needs.
func hashFactory(driverName string) (IHash, error) {
	switch driverName {
	case "argon2":
		return &drivers.Argon2Hash{
			Time:       3,
			Memory:     64 * 1024,
			Threads:    4,
			HashLength: 32,
			SaltLength: 16,
		}, nil
	case "pbkdf2":
		return &drivers.Pbkdf2Hash{
			Iteration:  600000,
			HashLength: 32,
			SaltLength: 16,
		}, nil
	case "bcrypt":
		return &drivers.BcryptHash{
			Cost: 12,
		}, nil
	case "scrypt":
		return &drivers.ScryptHash{
			Memory:     32 * 1024,
			BlockSize:  8,
			Threads:    3,
			HashLength: 32,
			SaltLength: 16,
		}, nil
	case "sha256":
		return &drivers.SHA256Hash{
			SaltLength: 16,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported hash driver: %s", driverName)
	}
}

// VerifyStoredHash checks if the provided hash corresponds to the given input string,
// parsing the driver name from the stored hash and using the appropriate hash driver.
func VerifyStoredHash(storedHash []byte, inputStr string) (bool, error) {
	parts := strings.Split(string(storedHash), ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid stored hash format")
	}
	driverName, hashedStr := parts[0], parts[1]

	driver := GetInstance(driverName)
	if driver == nil {
		return false, fmt.Errorf("failed to get hash instance for driver: %s", driverName)
	}

	return driver.Verify(hashedStr, inputStr)
}
