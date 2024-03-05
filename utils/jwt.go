package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
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
	} else {
		key, err = x509.ParsePKIXPublicKey(block.Bytes)
	}
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetTokenAndKeyFromAlgo returns a token and a key from an algorithm and a secret
func GetTokenAndKeyFromAlgo(algo, secret, keyPath string) (*jwt.Token, interface{}, error) {
	// Create token
	var token *jwt.Token
	var key interface{}
	var err error

	if algo == "HS512" {
		if len(secret) < 8 {
			return nil, nil, errors.New("secret must have at least 8 characters")
		}

		token = jwt.New(jwt.SigningMethodHS512)

		key = []byte(secret)
	} else if algo == "ES384" {
		token = jwt.New(jwt.SigningMethodES384)

		key, err = LoadECDSAKeyFromFile(keyPath, true)
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, errors.New("unsupported JWT algo: must be HS512 or ES384")
	}

	return token, key, nil
}

// GetKeyFromAlgo returns a key from an algorithm and a secret
func GetKeyFromAlgo(algo, secret, keyPath string) (interface{}, error) {
	var key interface{}
	var err error

	if algo == jwtware.HS512 {
		key = []byte(secret)
	} else if algo == jwtware.ES384 {
		key, err = LoadECDSAKeyFromFile(keyPath, false)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported JWT algo: must be HS512 or ES384")
	}

	return key, nil
}
