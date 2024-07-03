package main

import (
	"UserServiceAuth/internal/config"
	auth "UserServiceAuth/internal/router/auth"
	router "UserServiceAuth/internal/router/publickeygrpc"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
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
	log.Info("старт приложения",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.GRPC.Port))

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	routerGrpc := router.NewGrpcApi(grpcServer)
	_ = routerGrpc

	// Используем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup

	// Запуск gRPC сервера в отдельной горутине
	wg.Add(1)
	go func(grpcServer *grpc.Server) {
		defer wg.Done()
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
		if err != nil {
			log.Error("ошибка при запуске gRPC сервера", "error", err)
			return
		}
		defer lis.Close()

		log.Info("gRPC сервер запущен", slog.String("addr", lis.Addr().String()))

		if err := grpcServer.Serve(lis); err != nil {
			log.Error("ошибка при запуске gRPC сервера", "error", err)
		}
	}(grpcServer)

	// Создание сервера Echo
	e := echo.New()

	// Создание роутера
	httpRouter := auth.NewHttpRouter(e)
	_ = httpRouter

	// Запуск сервера Echo в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := e.Start(cfg.HTTP.Address); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start HTTP server", "error", err)
		}
	}()

	// Ожидание сигнала для остановки сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	grpcServer.GracefulStop()
	log.Info("gRPC сервер успешно остановлен")

	log.Info("Shutting down server...")

	// Создание контекста с таймаутом для плавной остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Остановка сервера Echo
	if err := e.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown the server", "error", err)
	} else {
		log.Info("Server gracefully stopped")
	}

	// Ожидание завершения всех горутин
	wg.Wait()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
