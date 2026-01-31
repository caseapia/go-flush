package repository

import (
	"context"
	"database/sql"
	"errors"

	models "github.com/caseapia/goproject-flush/internal/models/user"
)

// GetByID
func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	u := new(models.User)
	err := r.db.NewSelect().
		Model(u).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// GetAll
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.NewSelect().
		Model(&users).
		Scan(ctx)
	return users, err
}

// GetByName
func (r *UserRepository) GetByName(ctx context.Context, name string) (*models.User, error) {
	u := new(models.User)
	err := r.db.NewSelect().
		Model(u).
		Where("name = ?", name).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}
