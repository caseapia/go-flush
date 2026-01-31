package repository

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}
