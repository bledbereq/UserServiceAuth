package publickeygrpc

import (
	"context"

	ssov1 "UserServiceAuth/gen/go"

	"google.golang.org/grpc"
)

func RegisterGrpcRouter(server *grpc.Server) {
	ssov1.RegisterGetPublicKeyServer(server, &grpcApi{})
}

type grpcApi struct {
	ssov1.UnimplementedGetPublicKeyServer
}

func (s *grpcApi) PublicKey(ctx context.Context, req *ssov1.PublicKeyRequest) (*ssov1.PublicKeyResponse, error) {
	return &ssov1.PublicKeyResponse{PublicKey: "publickey12731723929381"}, nil
}
