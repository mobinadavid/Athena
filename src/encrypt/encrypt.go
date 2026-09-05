// Package encrypt provides a flexible encryption system that supports multiple encryption algorithms.
// It defines interfaces and implementations for encrypting and decrypting data, utilizing a driver-based approach.
package encrypt

import (
	"athena/src/config"
	"athena/src/encrypt/drivers"
	"errors"
	"fmt"
	"sync"
)

// IEncrypt defines an interface for encryption operations, including encrypting and decrypting data.
// This allows for the implementation of various encryption algorithms with a consistent interface.
type IEncrypt interface {
	Encrypt(data []byte) ([]byte, error)          // Encrypt encrypts the given byte slice.
	Decrypt(encryptedData []byte) ([]byte, error) // Decrypt decrypts the given byte slice.
	Sign(data []byte) ([]byte, error)             // Sign signs the given data.
}

// Encryptor encapsulates encryption operations using a specific encryption algorithm driver.
// It serves as the main entry point for encryption functionality in the application.
type Encryptor struct {
	driver     IEncrypt // The driver implementing the encryption algorithm.
	driverName string   // The name of the driver.
}

// Encrypt encrypts the input byte slice using the configured encryption driver.
func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	encryptedData, err := e.driver.Encrypt(data)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

// Decrypt decrypts the provided byte slice using the configured encryption driver.
func (e *Encryptor) Decrypt(encryptedData []byte) ([]byte, error) {
	return e.driver.Decrypt(encryptedData)
}

// Sign signs the given data using the configured encryption driver.
func (e *Encryptor) Sign(data []byte) ([]byte, error) {
	return e.driver.Sign(data)
}

var (
	once     sync.Once  // Ensures that the Encryptor instance is created only once.
	instance *Encryptor // Singleton instance of Encryptor.
)

// GetInstance returns a thread-safe singleton instance of Encryptor.
// It optionally takes a driver name and configures the Encryptor instance to use that driver.
// If no driver is specified, it uses the default driver or one specified by environment variables.
func GetInstance(driverNameArg ...string) *Encryptor {
	var err error
	once.Do(func() {
		driverName := getDriverName(driverNameArg...)
		var driver IEncrypt
		driver, err = encryptFactory(driverName)
		if err != nil {
			return
		}

		instance = &Encryptor{
			driver:     driver,
			driverName: driverName,
		}
	})

	if instance == nil {
		return nil
	}

	return instance
}

func getDriverName(args ...string) string {
	if len(args) > 0 && args[0] != "" {
		return args[0]
	}

	configs := config.GetInstance()
	envDriver := configs.Get("ENCRYPT_DRIVER")
	if envDriver != "" {
		return envDriver
	}

	return "aes" // Default to AES if no driver is specified.
}

// encryptFactory returns an IEncrypt instance based on the driver name.
func encryptFactory(driverName string) (IEncrypt, error) {
	key := config.GetInstance().Get("ENCRYPTION_SECRET")
	if key == "" {
		return nil, errors.New("ENCRYPTION_SECRET is required")
	}
	switch driverName {
	case "aes":
		return &drivers.AesEncrypt{
			Key: []byte(key),
		}, nil
	case "openpgp":
		return &drivers.OpenPgpEncrypt{}, nil
	case "rsa":
		rsaDriver, err := drivers.NewRsaEncrypt([]byte(""))
		if err != nil {
			return nil, err
		}
		return rsaDriver, nil
	default:
		return nil, fmt.Errorf("unsupported encrypt driver: %s", driverName)
	}
}
