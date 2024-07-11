package main

import (
	"UserServiceAuth/keygen"
	"fmt"
)

func main() {
	key_private_path := "../gen/key/private.pem"
	key_public_path := "../gen/key/public.pem"
	if err := keygen.GenerateKeyPair(key_private_path, key_public_path); err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	fmt.Println("RSA key pair generated and saved to private.pem and public.pem")
}
