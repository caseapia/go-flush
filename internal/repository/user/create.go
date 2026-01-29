package repository

import (
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
