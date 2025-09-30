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

// resolver, err := initResolver(cfg)
// if err != nil {
// 	slog.Error("failed to init resolver",
// 		slog.String("storage", cfg.Storage),
// 		slog.String("error", err.Error()),
// 	)
// 	os.Exit(1)
// }

// router := initRouter(log)
// srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
// 	Resolvers: resolver,
// }))

// srv.AddTransport(transport.GET{})
// srv.AddTransport(transport.POST{})
// srv.AddTransport(transport.Options{})
// srv.AddTransport(transport.Websocket{
// 	KeepAlivePingInterval: time.Second * 10,
// })
// srv.SetQueryCache(lru.New[*ast.QueryDocument](100)) //TODO: задавать через конфиг
// srv.Use(extension.Introspection{})

// router.Handle("/pground", playground.Handler("GraphQL playground", query))
// router.Handle(query, srv)

// err = startServer(cfg.HTTPServer, router, log)
// if err != nil {
// 	log.Error("failed ")
// }

// // TODO: вынести в отдельный пакет
// func startServer(cfg *config.HTTPServer, router *chi.Mux, log *slog.Logger) error {
// 	address := fmt.Sprintf("%s:%s", cfg.URL, cfg.Port)
// 	srv := &http.Server{
// 		Addr:    address,
// 		Handler: router,
// 	}

// 	errChan := make(chan error, 1)
// 	go func() {
// 		slog.Info("starting http server", slog.String("address", address))
// 		var err error
// 		if err = srv.ListenAndServe(); err != nil {
// 			slog.Error("failed to start http server", slog.String("error", err.Error()))
// 		}
// 		errChan <- err
// 	}()

// 	if err := <-errChan; err != nil {
// 		return err
// 	}

// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 	defer stop()

// 	<-ctx.Done()
// 	log.Info("shutting down server...")
// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()

// 	if err := srv.Shutdown(shutdownCtx); err != nil {
// 		log.Error("expected shutdown failed", slog.String("error", err.Error()))
// 		return err
// 	} else {
// 		log.Info("server stopped")
// 	}
// 	return nil
// }

// // TODO: вынести в отдельный пакет
// func initResolver(cfg *config.Config) (*graph.Resolver, error) {
// 	var resolver *graph.Resolver

// 	switch cfg.Storage {
// 	case "in-memory":
// 		storage := in_memory.NewStorage()
// 		resolver = &graph.Resolver{
// 			Storage:  storage,
// 			Post_:    storage.NewPostStorage(),
// 			Comment_: storage.NewCommentStorage(),
// 		}
// 	case "postgres":
// 		storage, err := postgres.NewStorage(*cfg.StorageConnect)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to initialize postgres database")
// 		}

// 		resolver = &graph.Resolver{
// 			Storage:  storage,
// 			Post_:    services.NewPostService(&storage.DB),
// 			Comment_: services.NewCommentService(&storage.DB),
// 		}
// 	default:
// 		return nil, fmt.Errorf("unknown storage type")
// 	}

// 	return resolver, nil
// }

// // TODO: вынести в отдельный пакет
// func initRouter(log *slog.Logger) *chi.Mux {
// 	slog.Info("starting router")
// 	router := chi.NewRouter()

// 	router.Use(middleware.RequestID)
// 	router.Use(logger.New(log))
// 	router.Use(middleware.Recoverer)

// 	return router
// }
