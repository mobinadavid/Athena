package drivers

import (
	"bytes"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
)

// OpenPgpEncrypt represents the OpenPGP encryption algorithm driver.
type OpenPgpEncrypt struct {
	PublicKey  []byte // PublicKey is the public key used for encryption.
	PrivateKey []byte // PrivateKey is the private key used for decryption. It should be encrypted with a passphrase.
	Passphrase []byte // Passphrase is the passphrase used to decrypt the private key.
}

// Encrypt encrypts the given byte slice using OpenPGP encryption with a public key.
func (o *OpenPgpEncrypt) Encrypt(data []byte) ([]byte, error) {
	pubKeyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(o.PublicKey))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	w, err := armor.Encode(buf, "PGP MESSAGE", nil)
	if err != nil {
		return nil, err
	}

	plaintext, err := openpgp.Encrypt(w, pubKeyRing, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	_, err = plaintext.Write(data)
	if err != nil {
		return nil, err
	}
	err = plaintext.Close()
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decrypt decrypts the given byte slice using OpenPGP encryption with a private key.
func (o *OpenPgpEncrypt) Decrypt(encryptedData []byte) ([]byte, error) {
	privateKeyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(o.PrivateKey))
	if err != nil {
		return nil, err
	}

	decBuffer := bytes.NewBuffer(encryptedData)
	block, err := armor.Decode(decBuffer)
	if err != nil {
		return nil, err
	}

	md, err := openpgp.ReadMessage(block.Body, privateKeyRing, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return o.Passphrase, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (o *OpenPgpEncrypt) Sign(data []byte) ([]byte, error) {
	return nil, nil
}
