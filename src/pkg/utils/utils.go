package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	mRand "math/rand"
	"reflect"
)

// GenerateSalt generates a new salt of the given length.
func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateRandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mRand.Intn(len(letterRunes))]
	}
	return string(b)
}

// IsLetter function to check if a character is a letter
func IsLetter(c byte) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

// IsDigit function to check if a character is a digit
func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsUpdateRequestEmpty checks if the update request has at least one non-nil value
func IsUpdateRequestEmpty(req interface{}) bool {
	// Get the reflect.Value of the update request
	reqValue := reflect.ValueOf(req)

	// Iterate through the fields of the update request
	for i := 0; i < reqValue.NumField(); i++ {
		fieldValue := reqValue.Field(i)

		// Check if the field is nil (i.e., not updated)
		if fieldValue.IsValid() && !fieldValue.IsNil() {
			return false
		}
	}

	return true
}

func Base64ToBigInt(b64 string) *big.Int {
	data, _ := base64.StdEncoding.DecodeString(b64)
	return new(big.Int).SetBytes(data)
}
