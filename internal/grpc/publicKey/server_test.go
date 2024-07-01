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
	publickeygrpc "UserServiceAuth/internal/grpc/publicKey"
)

// startTestGRPCServer запускает тестовый gRPC сервер
func startTestGRPCServer(t *testing.T, port int) (*grpc.Server, net.Listener, chan struct{}) {
	server := grpc.NewServer()
	publickeygrpc.Register(server)

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

	// Подождите немного, чтобы сервер точно успел запуститься
	time.Sleep(500 * time.Millisecond)

	return server, lis, done
}

func TestGRPCServer(t *testing.T) {
	port := 44044

	// Запуск тестового сервера
	server, lis, done := startTestGRPCServer(t, port)
	defer func() {
		server.GracefulStop()
		<-done
		lis.Close()
	}()

	// Создание gRPC клиента
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := ssov1.NewGetPublicKeyClient(conn)

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Отправка запроса и получение ответа
	req := &ssov1.PublicKeyRequest{}
	resp, err := client.PublicKey(ctx, req)
	if err != nil {
		t.Fatalf("PublicKey request failed: %v", err)
	}

	// Проверка ответа
	expectedPublicKey := "publickey12731723929381"
	assert.Equal(t, expectedPublicKey, resp.PublicKey, "PublicKey should match")
}
