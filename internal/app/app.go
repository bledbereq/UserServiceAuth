package app

import (
	grpcapp "UserServiceAuth/internal/app/grpc"
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
	//authService authgrpc.Auth,
) *App {
	//Инициализация хранилища
	// Инициализация сервисного слоя

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}

}
