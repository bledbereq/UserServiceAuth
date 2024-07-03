package main

import (
	"UserServiceAuth/internal/config"
	router "UserServiceAuth/internal/router"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "путь к файлу конфигурации")
	flag.Parse()

	// Проверка аргумента "config"
	if configPath == "" {
		fmt.Println("Использование: ./вашприложение -config <путь_к_файлу_конфигурации>")
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

	// Запуск gRPC сервера в отдельной горутине
	go func() {
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
	}()

	// Ожидание сигнала для остановки сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	grpcServer.GracefulStop()
	log.Info("приложение успешно остановлено")
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
