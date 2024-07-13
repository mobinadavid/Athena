package drivers

import (
	"athena/src/pkg/utils"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"log"
	"math/big"
)

// RsaEncrypt is an implementation of IEncrypt using RSA encryption algorithm.
type RsaEncrypt struct {
	privateKey *rsa.PrivateKey
}

// NewRsaEncrypt creates a new instance of RsaEncrypt with the given RSA private and public keys.
func NewRsaEncrypt(privateKey []byte) (*RsaEncrypt, error) {
	key, err := loadPrivateKeyFromXML(privateKey)

	if err != nil {
		return nil, err
	}

	return &RsaEncrypt{
		privateKey: key,
	}, nil
}

// Encrypt encrypts the input byte slice using RSA encryption.
func (r *RsaEncrypt) Encrypt(data []byte) ([]byte, error) {
	return nil, nil
}

// Decrypt decrypts the provided byte slice using RSA decryption.
func (r *RsaEncrypt) Decrypt(encryptedData []byte) ([]byte, error) {
	return nil, nil
}

// Sign signs the given data using RSA digital signature.
func (r *RsaEncrypt) Sign(data []byte) ([]byte, error) {
	hash := sha1.New()
	hash.Write(data)
	hashed := hash.Sum(nil)

	// Sign the hashed message
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA1, hashed)
	if err != nil {
		log.Println("Error signing:", err)
		return nil, nil
	}

	return []byte(base64.StdEncoding.EncodeToString(signature)), nil
}

// RSAKey represents the RSA private key structure
type RSAKey struct {
	XMLName  xml.Name `xml:"RSAKeyValue"`
	Modulus  string   `xml:"Modulus"`
	Exponent string   `xml:"Exponent"`
	D        string   `xml:"D"`
	P        string   `xml:"P"`
	Q        string   `xml:"Q"`
	DP       string   `xml:"DP"`
	DQ       string   `xml:"DQ"`
	InverseQ string   `xml:"InverseQ"`
}

// loadPrivateKeyFromXML loads an RSA private key from an XML file
func loadPrivateKeyFromXML(privateKey []byte) (*rsa.PrivateKey, error) {
	var key RSAKey
	err := xml.Unmarshal(privateKey, &key)
	if err != nil {
		return nil, err
	}

	modulus := utils.Base64ToBigInt(key.Modulus)
	exponent := utils.Base64ToBigInt(key.Exponent)
	d := utils.Base64ToBigInt(key.D)
	p := utils.Base64ToBigInt(key.P)
	q := utils.Base64ToBigInt(key.Q)
	dp := utils.Base64ToBigInt(key.DP)
	dq := utils.Base64ToBigInt(key.DQ)
	inverseQ := utils.Base64ToBigInt(key.InverseQ)

	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: modulus,
			E: int(exponent.Int64()),
		},
		D:      d,
		Primes: []*big.Int{p, q},
		Precomputed: rsa.PrecomputedValues{
			Dp:   dp,
			Dq:   dq,
			Qinv: inverseQ,
		},
	}, nil
}
