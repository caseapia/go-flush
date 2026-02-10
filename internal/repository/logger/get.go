package logger

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (l *LoggerRepository) GetAllLogs(ctx context.Context) ([]models.BaseLog, error) {
	var logs []models.BaseLog

	err := l.db.NewSelect().
		Model(&logs).
		Scan(ctx)

	return logs, err
}

func (l *LoggerRepository) GetCommonLogs(ctx context.Context) ([]models.CommonLog, error) {
	var logs []models.CommonLog

	err := l.db.NewSelect().
		Model(&logs).
		Scan(ctx)

	return logs, err
}

func (l *LoggerRepository) GetPunishmentLogs(ctx context.Context) ([]models.PunishmentLog, error) {
	var logs []models.PunishmentLog

	err := l.db.NewSelect().
		Model(&logs).
		Scan(ctx)

	return logs, err
}
