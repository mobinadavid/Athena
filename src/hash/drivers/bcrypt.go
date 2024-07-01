package drivers

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHash struct {
	Cost int
}

// Generate hashes a password using Bcrypt.
func (bcryptHash BcryptHash) Generate(str []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(str, bcryptHash.Cost)
}

// Verify compares a Bcrypt hashed password with a possible plaintext equivalent.
func (bcryptHash BcryptHash) Verify(hashedStr, str string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return err == nil, err
}
