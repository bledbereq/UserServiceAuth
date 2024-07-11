package publickeygrpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ssov1 "UserServiceAuth/gen/go"
	"UserServiceAuth/internal/router/publickeygrpc"
)

func startTestGRPCServer(t *testing.T, port int) (*grpc.Server, net.Listener, chan struct{}) {
	server := grpc.NewServer()

	// Path to your public key file
	publicKeyPath := "./public.pem"
	grpcService := publickeygrpc.NewGrpcApi(server, publicKeyPath)
	_ = grpcService

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	done := make(chan struct{})

	go func() {
		if err := server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("Failed to serve: %v", err)
		}
		close(done)
	}()

	return server, lis, done
}

func TestGRPCServer(t *testing.T) {
	port := 44050

	server, lis, done := startTestGRPCServer(t, port)
	defer func() {
		server.GracefulStop()
		<-done
		lis.Close()
	}()

	expectedPublicKey, err := publickeygrpc.LoadPublicKeyFromFile("../../../gen/key/public.pem")
	if err != nil {
		t.Fatalf("Failed to load expected public key: %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := ssov1.NewGetPublicKeyClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &ssov1.PublicKeyRequest{}
	resp, err := client.PublicKey(ctx, req)
	if err != nil {
		t.Fatalf("PublicKey request failed: %v", err)
	}

	assert.Equal(t, expectedPublicKey, resp.PublicKey, "Received public key does not match expected")
}
