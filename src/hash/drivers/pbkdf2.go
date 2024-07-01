package drivers

import (
	"athena/src/pkg/utils"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

type Pbkdf2Hash struct {
	Iteration  uint32
	SaltLength uint32
	HashLength uint32
}

// Generate hashes a password using PBKDF2.
func (pbkdf2hash *Pbkdf2Hash) Generate(str []byte) ([]byte, error) {
	// Generate a random salt
	salt, err := utils.GenerateSalt(pbkdf2hash.SaltLength)
	if err != nil {
		return []byte(""), err
	}
	// Generate the hash
	hash := pbkdf2.Key(str, salt, int(pbkdf2hash.Iteration), int(pbkdf2hash.HashLength), sha256.New) // 10000 iterations, 32-byte key length

	return []byte(fmt.Sprintf("%s$%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	)), nil
}

// Verify compares a PBKDF2 hashed password with a possible plaintext equivalent.
func (pbkdf2hash *Pbkdf2Hash) Verify(hashedStr, str string) (bool, error) {
	parts := strings.Split(hashedStr, "$")
	if len(parts) != 2 {
		return false, errors.New("hashed string format is incorrect")
	}

	// Decode the salt and hash from the stored value
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	originalHash, err := hex.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	// Generate a hash from the provided string using the same salt
	hash := pbkdf2.Key([]byte(str), salt, int(pbkdf2hash.Iteration), int(pbkdf2hash.HashLength), sha256.New) // Using the same parameters as when hashing

	// Compare the newly generated hash with the original hash
	return subtle.ConstantTimeCompare(originalHash, hash) == 1, nil
}
