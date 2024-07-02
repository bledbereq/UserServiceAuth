package main

import (
	"UserServiceAuth/internal/config"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"UserServiceAuth/internal/httphandler/login"
	"UserServiceAuth/internal/httphandler/register"

	"github.com/labstack/echo/v4"
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

	// Создание сервера Echo
	e := echo.New()
	e.POST("/register", register.RegisterHandler)
	e.POST("/login", login.LognHandler)

	// Запуск сервера в отдельной горутине
	go func() {
		if err := e.Start(cfg.HttpServer.Adress); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start HTTP server", "error", err)
		}
	}()

	// Ожидание сигнала для остановки сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Shutting down server...")

	// Создание контекста с таймаутом для плавной остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown the server", "error", err)
	} else {
		log.Info("Server gracefully stopped")
	}
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
