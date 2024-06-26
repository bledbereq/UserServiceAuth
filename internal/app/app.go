package app

import (
	grpcapp "UserServiceAuth/internal/app/grpc"
	authgrpc "UserServiceAuth/internal/grpc/auth"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	authService authgrpc.Auth,
) *App {
	//Инициализация хранилища
	// Инициализация сервисного слоя

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}

}
