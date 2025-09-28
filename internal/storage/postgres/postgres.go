package postgres

import (
	"client-services/internal/config"

	"github.com/go-pg/pg/v10"
)

type Storage struct {
	DB pg.DB
}

func NewStorage(cfg config.Config) {

}
