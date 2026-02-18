package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
)

func (r *Repository) PopulateNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification

	err := r.db.NewSelect().
		Model(&notifications).
		Order("created_at DESC").
		Relation("Sender").
		Relation("User").
		Where("user_id = ?", userID).
		Limit(COLUMNS_LIMIT).
		Scan(ctx)

	return notifications, err
}

func (r *Repository) SendNotification(ctx context.Context, entry models.Notification) {
	_, err := r.db.NewInsert().
		Model(&entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to send notification!")
		return
	}

	slog.WithData(slog.M{"entryData": entry}).Debugf("notification sended successfully")
}

func (r *Repository) ReadNotifications(ctx context.Context, userID uint64) []models.Notification {
	var notifications []models.Notification

	err := r.db.NewUpdate().
		Model(&notifications).
		Set("is_readed = ?", 1).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("error occured when trying to mark notifications as readed")
	}

	return notifications
}

func (r *Repository) ClearNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification

	_, err := r.db.NewDelete().
		Model(&notifications).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err, "userID": userID}).Error("failed to clear user notifications")
		return nil, err
	}

	return notifications, nil
}

func (r *Repository) RemoveNotification(ctx context.Context, userID, notifyID uint64) (bool, error) {
	res, err := r.db.NewDelete().
		Model((*models.Notification)(nil)).
		Where("user_id = ? AND id = ?", userID, notifyID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
