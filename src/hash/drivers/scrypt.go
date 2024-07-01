package drivers

import (
	"athena/src/pkg/utils"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

// ScryptHash struct implements the IHash interface for scrypt hashing.
// It includes parameters for the scrypt algorithm that can be adjusted based on security requirements.
type ScryptHash struct {
	Memory     int // CPU/memory cost parameter
	BlockSize  int // Block size parameter
	Threads    int // Parallelization parameter
	HashLength int // Length of the generated key
	SaltLength int // Size of the salt in bytes
}

// Generate creates a hash of a given byte slice using the scrypt algorithm.
// It returns the base64-encoded hash along with the salt, following a similar format to your Argon2 implementation.
func (scryptHash *ScryptHash) Generate(str []byte) ([]byte, error) {
	salt, err := utils.GenerateSalt(uint32(scryptHash.SaltLength)) // Assuming GenerateSalt is implemented in your utils package
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	hashedStr, err := scrypt.Key(str, salt, scryptHash.Memory, scryptHash.BlockSize, scryptHash.Threads, scryptHash.HashLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate scrypt hash: %w", err)
	}
	return []byte(fmt.Sprintf("%s$%s", base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hashedStr))), nil
}

// Verify checks whether a given plaintext string matches the scrypt hash.
func (scryptHash *ScryptHash) Verify(hashedStr, str string) (bool, error) {
	// Extract salt and hash from the hashedStr, assuming a format similar to the one used in Generate
	parts := strings.Split(hashedStr, "$")
	if len(parts) != 2 {
		return false, errors.New("incorrect hash format")
	}
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}
	hash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}
	computedHash, err := scrypt.Key([]byte(str), salt, scryptHash.Memory, scryptHash.BlockSize, scryptHash.Threads, len(hash))
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
