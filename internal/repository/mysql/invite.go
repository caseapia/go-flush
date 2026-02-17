package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"

)

func (r *Repository) SearchAllInvites(ctx context.Context) ([]models.Invite, error) {
	var invites []models.Invite

	err := r.db.NewSelect().
		Model(&invites).
		Relation("Creator").
		Relation("User").
		Order("created_at DESC").
		Limit(COLUMNS_LIMIT).
		Scan(ctx)

	return invites, err
}

func (r *Repository) CreateInvite(ctx context.Context, invite *models.Invite) error {
	_, err := r.db.NewInsert().
		Model(invite).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteInvite(ctx context.Context, inviteID uint64) error {
	_, err := r.db.NewDelete().
		Model((*models.Invite)(nil)).
		Where("id = ?", inviteID).
		Exec(ctx)

	return err
}

func (r *Repository) SearchInviteByCode(ctx context.Context, code string) (*models.Invite, error) {
	invite := new(models.Invite)

	err := r.db.NewSelect().
		Model(invite).
		Where("code = ?", code).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func (r *Repository) SearchInviteByID(ctx context.Context, id uint64) (*models.Invite, error) {
	invite := new(models.Invite)

	err := r.db.NewSelect().
		Model(invite).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return invite, err
}

func (r *Repository) MarkInviteAsUsed(ctx context.Context, inviteID, usedBy uint64) error {
	res, err := r.db.NewUpdate().
		Model((*models.Invite)(nil)).
		Set("used = ?", true).
		Set("used_by = ?", usedBy).
		Where("id = ?", inviteID).
		Where("used = ?", false).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &fiber.Error{Code: 404, Message: "invite already used or not found"}
	}

	return nil
}
