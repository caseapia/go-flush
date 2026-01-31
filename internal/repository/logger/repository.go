package loggerrepository

import (
	"github.com/uptrace/bun"
)

type LoggerRepository struct {
	db *bun.DB
}

func NewLoggerRepository(db *bun.DB) *LoggerRepository {
	return &LoggerRepository{db: db}
}
