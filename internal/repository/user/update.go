package repository

import (
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
