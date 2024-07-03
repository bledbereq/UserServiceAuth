package main

import (
	app "UserServiceAuth/internal/app"
	"UserServiceAuth/internal/config"
	"flag"
	"fmt"
	"log/slog"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	// Валидация аргумента "config"
	if configPath == "" {
		fmt.Println("Usage: ./yourapp -config <path_to_config_file>")
		os.Exit(1)
	}

	cfg := config.MustLoadByPath(configPath)
	fmt.Println(cfg)
	log := setupLogger(cfg.Env)
	log.Info("starting app",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.GRPC.Port))

	// Запуск gRPC сервера
	application := app.New(log, cfg.GRPC.Port)
	go application.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	log.Info("Gracefully stopped")
}

const (
	envDev  = "dev"
	envProd = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
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
