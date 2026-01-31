package repository

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	_, err := r.db.NewDelete().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

func (r *UserRepository) Restore(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}
