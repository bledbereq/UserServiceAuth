package publickeygrpc

import (
	"context"

	ssov1 "UserServiceAuth/gen/go"

	"google.golang.org/grpc"
)

func RegisterRouter(gRPC *grpc.Server) {
	ssov1.RegisterGetPublicKeyServer(gRPC, &serverAPI{})
}

type serverAPI struct {
	ssov1.UnimplementedGetPublicKeyServer
}

func (s *serverAPI) PublicKey(
	ctx context.Context,
	req *ssov1.PublicKeyRequest,
) (*ssov1.PublicKeyResponse, error) {

	return &ssov1.PublicKeyResponse{PublicKey: "publickey12731723929381"}, nil
}
