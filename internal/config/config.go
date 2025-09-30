package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string          `yaml:"env" env:"ENV" env-default:"local" env-requered:"true"`
	Storage        string          `yaml:"storage" env-default:"in-memory"`
	QueryCache     int             `yaml:"query-cache" env-default:"100"`
	StorageConnect *StorageConnect `yaml:"storage_connect"`
	HTTPServer     *HTTPServer     `yaml:"http_server"`
}

type StorageConnect struct {
	SQLDriver   string `yaml:"sql_driver" env-default:"postgres"`
	SQLUser     string `yaml:"sql_user" env-default:"postgres"`
	SQLPassword string `yaml:"sql_password" env-default:"postgres"`
	SQLAddress  string `yaml:"sql_address" env-default:"localhost"`
	SQLPort     string `yaml:"sql_port" env-default:"5432"`
	SQLDBName   string `yaml:"sql_dbname" env-default:"client-service"`
	SQLSSLMode  string `yaml:"sql_sslmode" env-default:"disable"`
}

type HTTPServer struct {
	URL          string        `yaml:"url" env-default:"localhost"`
	Port         string        `yaml:"port" env-default:":8080"`
	Timeout      time.Duration `yaml:"timeout" env-default:"10s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		slog.Error("CONFIG_PATH is not set")
		os.Exit(1)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("cannot find config-file by path",
			slog.String("path", configPath),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("failed to read config file",
			slog.String("path", configPath),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	return &cfg
}
