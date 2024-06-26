package main

import (
	"UserServiceAuth/internal/app"
	"UserServiceAuth/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := setupLogger(cfg.Env)
	log.Info("starting app",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.GRPS.Port))
	log.Debug("debage message")
	log.Error("error message")
	log.Warn("warn message")

	application := app.New(log, cfg.GRPS.Port, cfg.StoragePath, cfg.TokenTTL)
	application.GRPCSrv.MustRun()

	// grpc сервер и хендлеры
	// инициализировать точку входа в приложение

	//Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	log.Info("Gracefully stopped")

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
