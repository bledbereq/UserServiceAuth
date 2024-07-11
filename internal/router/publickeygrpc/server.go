package publickeygrpc

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	ssov1 "UserServiceAuth/gen/go"

	"google.golang.org/grpc"
)

type GrpcApi struct {
	ssov1.UnimplementedGetPublicKeyServer
	publicKey string
}

func NewGrpcApi(server *grpc.Server, publicKeyPath string) *GrpcApi {

	publicKeyBase64, _ := LoadPublicKeyFromFile(publicKeyPath)

	router := &GrpcApi{
		publicKey: publicKeyBase64,
	}
	ssov1.RegisterGetPublicKeyServer(server, router)
	return router
}

func (s *GrpcApi) PublicKey(ctx context.Context, req *ssov1.PublicKeyRequest) (*ssov1.PublicKeyResponse, error) {
	return &ssov1.PublicKeyResponse{PublicKey: s.publicKey}, nil
}

func LoadPublicKeyFromFile(filepath string) (string, error) {
	publicKeyBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key file: %v", err)
	}

	// Assuming the file contains PEM encoded public key
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block from public key file")
	}
	if block.Type != "PUBLIC KEY" {
		return "", fmt.Errorf("invalid PEM block type, expected PUBLIC KEY")
	}

	// Encode to Base64 to match the expected format
	publicKeyBase64 := base64.StdEncoding.EncodeToString(block.Bytes)

	return publicKeyBase64, nil
}
