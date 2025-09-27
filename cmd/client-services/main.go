package main

import (
	"client-services/internal/config"
	"client-services/internal/graph"
	"client-services/internal/server/middlewares/logger"
	"log/slog"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/dikkadev/prettyslog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

// адрес query-запроса
const (
	query = "/query"
)

func main() {
	// загружаем .env файл для параметров окружения
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
		slog.String("storage-type", cfg.GetStorageLink()),
	)
	slog.Debug("debug messages are enabled")
	slog.Error("error messaages are enabled")

	//TODO бд

	router := initRouter(log)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{
			//TODO: Подача в структуру в resolver.go
		}},
	))

	router.Handle(query, srv)
}

// TODO: вынести в отдельный пакет
func initRouter(log *slog.Logger) *chi.Mux {
	slog.Info("starting router")
	router := chi.NewRouter()

	// подключаем middlewares
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer) // защита от паник

	return router
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
