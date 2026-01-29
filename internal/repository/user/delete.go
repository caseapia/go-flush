package repository

import models "github.com/caseapia/goproject-flush/internal/models/user"

func (r *UserRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}
