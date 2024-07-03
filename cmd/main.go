package main

import (
	"UserServiceAuth/internal/config"
	auth "UserServiceAuth/internal/router"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	// Создание роутера
	httpRouter := auth.NewHttpRouter(e)
	_ = httpRouter

	// Используем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup

	// Добавление одной горутины в WaitGroup
	wg.Add(1)

	// Запуск сервера Echo в отдельной горутине
	go func() {
		defer wg.Done()

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

	// Ожидание завершения всех горутин
	wg.Wait()

	// Остановка сервера Echo
	if err := e.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown the server", "error", err)
	} else {
		log.Info("Server gracefully stopped")
	}
	wg.Done()

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
