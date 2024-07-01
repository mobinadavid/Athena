package drivers

import (
	"athena/src/pkg/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

type SHA256Hash struct {
	SaltLength uint32
}

// Generate hashes the input using SHA-256 with a generated salt.
// The output format is "salt:hash" where both parts are hex-encoded.
func (sha256hash SHA256Hash) Generate(str []byte) ([]byte, error) {
	// Generate a random salt
	salt, err := utils.GenerateSalt(sha256hash.SaltLength)
	if err != nil {
		return []byte(""), err
	}

	// Prepend the salt to the input and hash
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write(str)
	hash := hasher.Sum(nil)

	// Return the formatted "salt:hash" output
	return []byte(fmt.Sprintf("%s$%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	)), nil
}

// Verify checks if the provided input matches the hashed value.
// It expects the hashedStr format to be "salt:hash" with both parts hex-encoded.
func (sha256hash SHA256Hash) Verify(hashedStr, input string) (bool, error) {
	parts := strings.Split(hashedStr, "$")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the input with the extracted salt
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	// Compare the computed hash with the expected hash
	return hmac.Equal(hash, expectedHash), nil
}
