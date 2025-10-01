package main

import (
	"client-services/internal/config"
	"client-services/internal/run"
	"log/slog"
	"os"

	"github.com/dikkadev/prettyslog"
	"github.com/joho/godotenv"
)

// путь к .env-файлу
const (
	pathDocker = ".env"
)

// уровни логирования
const (
	envLocal = "local"
	envDebug = "debug"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(pathDocker); err != nil {
		slog.Error("failed to load .env file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	cfg := config.MustLoad()
	slog.Info("config file loaded successfully")

	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	slog.Info("starting service",
		slog.String("env", cfg.Env),
		slog.String("storage-type", cfg.Storage),
	)
	slog.Debug("debug messages are enabled")
	slog.Error("error messaages are enabled")

	run.Run(cfg, log)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(prettyslog.NewPrettyslogHandler("ClientServices",
			prettyslog.WithLevel(slog.LevelDebug),
		))
	case envDebug:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
