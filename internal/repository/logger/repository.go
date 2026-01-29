package loggerrepository

import "gorm.io/gorm"

type LoggerRepository struct {
	db *gorm.DB
}

func NewLoggerRepository(db *gorm.DB) *LoggerRepository {
	return &LoggerRepository{db: db}
}
