package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"UserServiceAuth/internal/config"
	auth "UserServiceAuth/internal/router/auth"
	router "UserServiceAuth/internal/router/publickeygrpc"
	"UserServiceAuth/internal/router/repositories"
	services "UserServiceAuth/internal/uscase"
	"UserServiceAuth/storage"
	"context"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
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

	if _, err := os.Stat(cfg.KeyPrivatePath); os.IsNotExist(err) {
		fmt.Println("Private key not found, generating a new one...")
		if err := generateKeyPair(cfg.KeyPrivatePath, cfg.KeyPublicPath); err != nil {
			fmt.Printf("Error: failed to generate key pair: %v\n", err)
			return
		}
		fmt.Println("Key pair generated successfully.")
	} else {
		fmt.Println("Private key already exists.")
	}

	privateKey, err := loadPrivateKeyFromFile(cfg.KeyPrivatePath)
	if err != nil {
		log.Error("ошибка при загрузке приватного ключа", "error", err)
		return
	}

	tokenService := services.NewTokenService(privateKey, &privateKey.PublicKey)

	// Инициализация UserService с использованием tokenService
	userService := services.NewUserService(userRepo, tokenService)

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

func loadPrivateKeyFromFile(privateKeyFile string) (*rsa.PrivateKey, error) {
	// Читаем содержимое файла
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла приватного ключа: %w", err)
	}

	// Декодируем PEM-блок
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("ошибка декодирования PEM блока")
	}

	// Парсим закрытый ключ RSA
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга приватного ключа: %w", err)
	}

	return privateKey, nil
}

func generateKeyPair(privateKeyPath, publicKeyPath string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(privateKeyPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create and write private key file
	privateFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateFile.Close()

	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, privatePEM); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	// Generate public key
	publicKey := &privateKey.PublicKey

	// Create and write public key file
	publicFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer publicFile.Close()

	publicPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	publicPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicPEM,
	}
	if err := pem.Encode(publicFile, publicPEMBlock); err != nil {
		return fmt.Errorf("failed to write public key: %v", err)
	}

	return nil
}
