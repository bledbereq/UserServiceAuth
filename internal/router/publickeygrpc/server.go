package publickeygrpc

import (
	"context"

	ssov1 "UserServiceAuth/gen/go"

	"google.golang.org/grpc"
)

func NewGrpcApi(server *grpc.Server) *GrpcApi {
	router := &GrpcApi{}
	ssov1.RegisterGetPublicKeyServer(server, router)
	return router
}

type GrpcApi struct {
	ssov1.UnimplementedGetPublicKeyServer
}

func (s *GrpcApi) PublicKey(ctx context.Context, req *ssov1.PublicKeyRequest) (*ssov1.PublicKeyResponse, error) {
	return &ssov1.PublicKeyResponse{PublicKey: "publickey12731723929381"}, nil
}
