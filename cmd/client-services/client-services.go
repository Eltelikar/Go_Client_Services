package main

import (
	"client-servicec/internal/config"
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

// адреса http-запросов
const ()

func main() {
	// загружаем .env файл для параметров окружения
	if err := godotenv.Load(pathDocker); err != nil {
		slog.Error("failed to load .env file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	cfg := config.NewConfig()
	slog.Info("config file loaded successfully")

	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	slog.Info("Test",
		slog.String("some text", "test text"),
		slog.String("some text", "test text"),
		slog.String("some text", "test text"),
	)

	//TODO логгер
	//TODO бд
	//TODO сервер
	//TODO хендлеры
}

// TODO:вынести в отдельный пакет?
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
	}

	return log
}
