package run

import (
	"client-services/internal/config"
	"client-services/internal/graph"
	uqmutex "client-services/internal/graph/unique-mutex"
	"client-services/internal/server/middlewares/logger"
	"client-services/internal/services"
	in_memory "client-services/internal/storage/in-memory"
	"client-services/internal/storage/postgres"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vektah/gqlparser/v2/ast"
)

func Run(cfg *config.Config, log *slog.Logger) {
	router := initRouter(log)

	resolver, err := initResolver(cfg)
	if err != nil {
		slog.Error("failed to init resolver",
			slog.String("storage", cfg.Storage),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	srv := initGraphQL(cfg.QueryCache, resolver)

	router.Handle("/pground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	err = startServer(cfg.HTTPServer, router, log)
	if err != nil {
		log.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// TODO: вынести в отдельный пакет
func startServer(cfg *config.HTTPServer, router *chi.Mux, log *slog.Logger) error {
	address := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	errChan := make(chan error, 1)
	go func() {
		slog.Info("starting http server", slog.String("address", address))
		var err error
		if err = srv.ListenAndServe(); err != nil {
			slog.Error("failed to start http server", slog.String("error", err.Error()))
		}
		errChan <- err
	}()

	if err := <-errChan; err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	slog.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("expected shutdown failed", slog.String("error", err.Error()))
		return err
	} else {
		slog.Info("server stopped")
	}
	return nil
}

func initGraphQL(queryCache int, resolver *graph.Resolver) *handler.Server {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: time.Second * 10,
	})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](queryCache))
	srv.Use(extension.Introspection{})

	slog.Info("graphql initialized successfully")
	return srv
}

// TODO: вынести в отдельный пакет
func initResolver(cfg *config.Config) (*graph.Resolver, error) {
	var resolver *graph.Resolver

	switch cfg.Storage {
	case "in-memory":
		storage := in_memory.NewStorage()
		resolver = &graph.Resolver{
			Log:      slog.Default(),
			Storage:  storage,
			Post_:    storage.NewPostStorage(),
			Comment_: storage.NewCommentStorage(),
			UqMutex:  uqmutex.NewUqMutex(),
		}
	case "postgres":
		storage, err := postgres.NewStorage(*cfg.StorageConnect)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize postgres database: %w", err)
		}

		resolver = &graph.Resolver{
			Log:      slog.Default(),
			Storage:  storage,
			Post_:    services.NewPostService(&storage.DB),
			Comment_: services.NewCommentService(&storage.DB),
			UqMutex:  uqmutex.NewUqMutex(),
		}
	default:
		return nil, fmt.Errorf("unknown storage type")
	}

	slog.Info("resolver initialized successfully", slog.String("storage type", cfg.Storage))
	return resolver, nil
}

// TODO: вынести в отдельный пакет
func initRouter(log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	slog.Info("router started")
	return router
}
