package main

import (
	"UserServiceAuth/internal/app"
	"UserServiceAuth/internal/config"
	"UserServiceAuth/lib/logger/handler/slogpretty"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
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

	// Запуск RestApi сервера
	go func() {
		s := echo.New()
		s.POST("/register", RegisterHandler)
		s.POST("/login", LognHandler)
		if err := s.Start(cfg.HttpServer.Adress); err != nil {
			log.Error("failed to start HTTP server")
		}
	}()

	// Запуск gRPC сервера
	application := app.New(log, cfg.GRPS.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.Run()
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
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
