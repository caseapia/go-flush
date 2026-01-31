package repository

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}
