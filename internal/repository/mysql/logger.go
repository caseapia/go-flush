package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
)

func (r *Repository) fetchLogs(ctx context.Context, tableName string) ([]models.LogDTO, error) {
	var logs []models.LogDTO

	err := r.db.NewSelect().
		Model(&logs).
		ModelTableExpr("flushproject_logs." + tableName + " AS l").
		ColumnExpr("l.*").
		ColumnExpr("u_admin.name AS sender_name").
		ColumnExpr("u_target.name AS user_name").
		Join("LEFT JOIN flushproject.users AS u_admin ON u_admin.id = l.admin_id").
		Join("LEFT JOIN flushproject.users AS u_target ON u_target.id = l.user_id").
		Order("l.date DESC").
		Scan(ctx)

	return logs, err
}

func (r *Repository) GetCommonLogs(ctx context.Context) ([]models.LogDTO, error) {
	return r.fetchLogs(ctx, "admin_common")
}

func (r *Repository) GetPunishmentLogs(ctx context.Context) ([]models.LogDTO, error) {
	return r.fetchLogs(ctx, "admin_punishments")
}

func (l *Repository) SavePunishmentLog(ctx context.Context, entry interface{}) error {
	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to insert action log!")
		return err
	}

	slog.WithData(slog.M{
		"entryData": entry,
	}).Debugf("log inserted successfully")

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

	slog.WithData(slog.M{
		"entryData": entry,
	}).Debugf("log inserted successfully")

	return nil
}
