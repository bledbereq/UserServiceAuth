package main

import (
	"UserServiceAuth/internal/config"
	auth "UserServiceAuth/internal/router/auth"
	router "UserServiceAuth/internal/router/publickeygrpc"
	"UserServiceAuth/internal/router/repositories"
	services "UserServiceAuth/internal/uscase"
	"UserServiceAuth/storage"
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

	"github.com/go-playground/validator/v10"
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
	log := setupLogger(cfg.Env)
	log.Info("старт приложения",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg))

	var wg sync.WaitGroup

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	routerGrpc := router.NewGrpcApi(grpcServer)
	_ = routerGrpc

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Error("ошибка при запуске gRPC сервера", "error", err)
		return
	}
	defer lis.Close()

	log.Info("gRPC сервер запущен", slog.String("addr", lis.Addr().String()))

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("ошибка при запуске gRPC сервера", "error", err)
		}
	}()

	// Создание сервера Echo
	e := echo.New()

	db := storage.InitDB(cfg)
	userRepo := repositories.NewUserRepository(db)

	// Создание сервиса
	userService := services.NewUserService(userRepo)

	// Создание валидатора
	validator := validator.New()

	// Создание и настройка HTTP роутера
	authRouter := auth.NewHttpRouter(e, userService, validator)
	_ = authRouter

	// Запуск сервера Echo
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.Start(cfg.HTTP.Address); err != nil && err != http.ErrServerClosed {
			log.Error("ошибка при запуске HTTP сервера", "error", err)
		}
	}()

	// Ожидание сигнала для остановки серверов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Получен сигнал на остановку серверов")

	// Остановка gRPC сервера
	grpcServer.GracefulStop()
	log.Info("gRPC сервер успешно остановлен")

	// Создание контекста с таймаутом для плавной остановки HTTP сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Остановка HTTP сервера
	if err := e.Shutdown(ctx); err != nil {
		log.Error("ошибка при остановке HTTP сервера", "error", err)
	} else {
		log.Info("HTTP сервер успешно остановлен")
	}

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Info("Сервера успешно остановлены")
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
