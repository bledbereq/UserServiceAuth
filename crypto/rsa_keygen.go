package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func generateKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	privateFile, err := os.Create("gen/key/private.pem")
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, privatePEM); err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey
	publicFile, err := os.Create("gen/key/public.pem")
	if err != nil {
		return err
	}
	defer publicFile.Close()

	publicPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicPEM,
	}
	if err := pem.Encode(publicFile, publicPEMBlock); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := generateKeyPair(); err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	fmt.Println("RSA key pair generated and saved to private.pem and public.pem")
}
