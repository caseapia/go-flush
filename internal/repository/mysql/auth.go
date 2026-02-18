package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
)

func (r *Repository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return user, err
}

func (r *Repository) SearchByLogin(ctx context.Context, login string) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("email = ? OR name = ?", login, login).
		Limit(1).
		Scan(ctx)

	return u, err
}

func (r *Repository) SearchByID(ctx context.Context, id uint64) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("id = ?", id).
		Scan(ctx)

	return u, err
}

func (r *Repository) UpdateTokenVersion(ctx context.Context, userID uint64, version int) error {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("token_version = ?", version).
		Where("id = ?", userID).
		Exec(ctx)

	return err
}

func (r *Repository) CheckMultiAccountByFingerprint(ctx context.Context, userID int64, fingerprint string) (bool, int, error) {
	var count int

	count, err := r.db.NewSelect().
		Model((*models.Fingerprint)(nil)).
		Where("hash = ? AND user_id != ?", fingerprint, userID).
		Count(ctx)
	if err != nil {
		return false, 0, err
	}

	if count > 0 {
		slog.WithData(slog.M{
			"count": count,
		}).Warn("This fingerprint is already used by another user!")
		return true, count, nil
	}

	return false, 0, nil
}

func (r *Repository) CheckMultiAccountByIP(ctx context.Context, userID int64, ip string) (bool, int, error) {
	var count int

	count, err := r.db.NewSelect().
		Model((*models.Fingerprint)(nil)).
		Where("ip = ? AND user_id = ?", ip, userID).
		Count(ctx)
	if err != nil {
		return false, 0, err
	}

	if count > 3 {
		slog.WithData(slog.M{
			"count": count,
		}).Warn("This IP is already have more than 3 accounts")
		return true, count, nil
	}

	return false, 0, nil
}

func (r *Repository) CheckMultiAccountByUA(ctx context.Context, userID int64, userAgent string) (bool, int, error) {
	var count int

	count, err := r.db.NewSelect().
		Model((*models.Fingerprint)(nil)).
		Where("user_agent = ? AND user_id != ?", userAgent, userID).
		Count(ctx)
	if err != nil {
		return false, 0, err
	}
	if count > 3 {
		fmt.Println("This UserAgent is already have more than 3 accounts")
		return true, count, nil
	}

	return false, 0, nil
}

func (r *Repository) RegisterFingerprint(ctx context.Context, userID int64, hash, ip, userAgent string) error {
	fp := &models.Fingerprint{
		UserID:    userID,
		Hash:      hash,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	_, err := r.db.NewInsert().Model(fp).Exec(ctx)
	return err
}
