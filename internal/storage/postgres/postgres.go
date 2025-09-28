package postgres

import (
	"client-services/internal/config"
	"client-services/internal/graph/model"
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Storage struct {
	DB pg.DB
}

func NewStorage(cfg config.StorageConnect) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	BDAddr := fmt.Sprintf("%s:%s", cfg.SQLAddress, cfg.SQLPort)
	conn := pg.Connect(&pg.Options{
		Addr:     BDAddr,
		User:     cfg.SQLUser,
		Password: cfg.SQLPassword,
		Database: cfg.SQLDBName,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to connect database: %w", op, err)
	}

	s := &Storage{DB: *conn}

	if err := migrate(s); err != nil {
		return nil, fmt.Errorf("%s: failed to migrate: %w", op, err)
	}

	return s, nil
}

func migrate(s *Storage) error {
	const op = "storage.postgres.migrate"
	_ = op

	schemas := []interface{}{
		(*model.Post)(nil),
		(*model.Comment)(nil),
	}

	for _, schem := range schemas {
		err := s.DB.Model(schem).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

func (s *Storage) CloseDB() error {
	return s.DB.Close()
}
