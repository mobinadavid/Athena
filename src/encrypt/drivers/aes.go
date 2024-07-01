package drivers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AesEncrypt represents the AES encryption algorithm driver.
type AesEncrypt struct {
	Key []byte // Key is the encryption key. Its length must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
}

// Encrypt encrypts the given byte slice using AES encryption.
func (a *AesEncrypt) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)
	return encryptedData, nil
}

// Decrypt decrypts the given byte slice using AES encryption.
func (a *AesEncrypt) Decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:gcm.NonceSize()], encryptedData[gcm.NonceSize():]
	data, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a *AesEncrypt) Sign(data []byte) ([]byte, error) {
	return nil, nil
}
