package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"

)

func (r *Repository) GetCommonLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.CommonLog, int, error) {
	var logs []models.CommonLog

	query := r.db.NewSelect().
		Model(&logs).
		Relation("Admin").
		Relation("User").
		Order("date DESC").
		Limit(LOGS_COLUMNS_LIMIT)

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	if keywords != "" {
		query = query.WhereGroup("AND", func(q *bun.SelectQuery) *bun.SelectQuery {
			keyword := "%" + keywords + "%"
			return q.Where("LOWER(action) LIKE LOWER(?)", keyword).
				WhereOr("LOWER(additional_information) LIKE LOWER(?)", keyword)
		})
	}

	err := query.Scan(ctx)
	return logs, COLUMNS_LIMIT, err
}

func (r *Repository) GetPunishmentLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.PunishmentLog, int, error) {
	var logs []models.PunishmentLog

	query := r.db.NewSelect().
		Model(&logs).
		Relation("Admin").
		Relation("User").
		Order("date DESC").
		Limit(COLUMNS_LIMIT)

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	if keywords != "" {
		query = query.WhereGroup("AND", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("action ILIKE ?", "%"+keywords+"%").
				WhereOr("additional_information ILIKE ?", "%"+keywords+"%")
		})
	}

	err := query.Scan(ctx)
	return logs, COLUMNS_LIMIT, err
}

func (l *Repository) SavePunishmentLog(ctx context.Context, entry interface{}) error {
	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to insert action log!")
		return err
	}

	slog.WithData(slog.M{"entryData": entry}).Debugf("log inserted successfully")
	return nil
}

func (l *Repository) SaveCommonLog(ctx context.Context, entry interface{}) error {
	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to insert action log!")
		return err
	}

	slog.WithData(slog.M{"entryData": entry}).Debugf("log inserted successfully")
	return nil
}
