package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// GenerateKeyPair generates an RSA key pair and saves them to specified paths.
func GenerateKeyPair(privateKeyPath, publicKeyPath string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(privateKeyPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create and write private key file
	privateFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateFile.Close()

	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, privatePEM); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	// Generate public key
	publicKey := &privateKey.PublicKey

	// Create and write public key file
	publicFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer publicFile.Close()

	publicPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	publicPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicPEM,
	}
	if err := pem.Encode(publicFile, publicPEMBlock); err != nil {
		return fmt.Errorf("failed to write public key: %v", err)
	}

	return nil
}

// LoadPrivateKeyFromFile loads an RSA private key from the specified file.
func LoadPrivateKeyFromFile(privateKeyFile string) (*rsa.PrivateKey, error) {
	// Read file content
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("error decoding PEM block")
	}

	// Parse RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	return privateKey, nil
}
