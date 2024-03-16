package helper

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func GetRSAPrivateKey() (*rsa.PrivateKey, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %s", err)
	}

	privateKeyPath := filepath.Join(rootDir, "private_pkcs1.pem")
	keyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	// Read the private key data
	keyBytes, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, err
	}

	// Decode PEM-encoded data
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Parse the RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GetRSAPublicKey() (*rsa.PublicKey, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %s", err)
	}

	publicKeyPath := filepath.Join(rootDir, "public.pem")
	keyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	// Read the public key data
	keyBytes, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, err
	}

	// Decode PEM-encoded data
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the RSA public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPublicKey, nil
}
