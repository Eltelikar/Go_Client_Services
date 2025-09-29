package main

import (
	"client-services/internal/config"
	"client-services/internal/graph"
	"client-services/internal/server/middlewares/logger"
	"client-services/internal/services"
	in_memory "client-services/internal/storage/in-memory"
	"client-services/internal/storage/postgres"
	"fmt"
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
		slog.String("storage-type", cfg.Storage),
	)
	slog.Debug("debug messages are enabled")
	slog.Error("error messaages are enabled")

	resolver, err := initResolver(cfg)
	if err != nil {
		slog.Error("failed to init resolver",
			slog.String("storage", cfg.Storage),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	router := initRouter(log)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	router.Handle(query, srv)
}

// TODO: вынести в отдельный пакет
func initResolver(cfg *config.Config) (*graph.Resolver, error) {
	var resolver *graph.Resolver

	switch cfg.Storage {
	case "in-memory":
		storage := in_memory.NewStorage()
		resolver = &graph.Resolver{
			Storage:  storage,
			Post_:    storage.NewPostStorage(),
			Comment_: storage.NewCommentStorage(),
		}
	case "postgres":
		storage, err := postgres.NewStorage(*cfg.StorageConnect)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize postgres database")
		}

		resolver = &graph.Resolver{
			Storage:  storage,
			Post_:    services.NewPostService(&storage.DB),
			Comment_: services.NewCommentService(&storage.DB),
		}
	default:
		return nil, fmt.Errorf("unknown storage type")
	}

	return resolver, nil
}

// TODO: вынести в отдельный пакет
func initRouter(log *slog.Logger) *chi.Mux {
	slog.Info("starting router")
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

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
