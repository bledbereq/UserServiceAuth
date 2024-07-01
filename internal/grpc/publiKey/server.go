package publickeygrpc

import (
	"context"

	ssov1 "UserServiceAuth/gen/go"

	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedGetPublicKeyServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterGetPublicKeyServer(gRPC, &serverAPI{})

}

func (s *serverAPI) PublicKey(
	ctx context.Context,
	req *ssov1.PublicKeyRequest,
) (*ssov1.PublicKeyResponse, error) {
	return &ssov1.PublicKeyResponse{
		PublicKey: "publickey12731723929381",
	}, nil

}
