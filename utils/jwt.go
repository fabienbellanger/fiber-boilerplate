package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// LoadECDSAKeyFromFile loads an ECDSA private or public key from a file
func LoadECDSAKeyFromFile(filename string, isPrivate bool) (any, error) {
	// Read file
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decode PEM
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("error when decoding .pem file")
	}

	// Parse key
	var key any
	if isPrivate {
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		key, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	return key, nil
}
