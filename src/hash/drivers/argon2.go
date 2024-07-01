package drivers

import (
	"athena/src/pkg/utils"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// Argon2Hash is an implementation of the Hash interface using Argon2.
type Argon2Hash struct {
	// Time represents the number of passed over the specified memory.
	Time uint32
	// Memory to be used.
	Memory uint32
	// Threads for parallelism aspect of the algorithm.
	Threads uint8
	// HashLength the length of the generate hash key.
	HashLength uint32
	// SaltLength the length of the salt used.
	SaltLength uint32
}

// Generate hashes a []byte using Argon2.
func (argon *Argon2Hash) Generate(str []byte) ([]byte, error) {
	salt, err := utils.GenerateSalt(argon.SaltLength)
	if err != nil {
		return []byte(""), err
	}
	hashedStr := argon2.IDKey(str, salt, argon.Time, argon.Memory, argon.Threads, argon.HashLength)
	// Return the base64-encoded hash and salt
	return []byte(fmt.Sprintf(
		"%s$%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hashedStr),
	)), nil
}

// Verify compares an Argon2 hashed []byte with a possible plaintext equivalent.
func (argon *Argon2Hash) Verify(hashedStr, str string) (bool, error) {
	parts := strings.Split(hashedStr, "$")
	if len(parts) != 2 {
		return false, fmt.Errorf("incorrect hash format")
	}
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	hash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(str), salt, argon.Time, argon.Memory, argon.Threads, uint32(len(hash)))
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
