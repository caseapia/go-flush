package service

func (s *UserService) DeleteUser(id uint) (bool, error) {
	user, err := s.repo.GetByID(id)

	if err != nil {
		return false, ErrUserNotFound
	}

	if user.IsDeleted {
		return user.IsDeleted, ErrUserNotFound
	}

	user.IsDeleted = true

	if err := s.repo.Update(user); err != nil {
		return false, err
	}

	return true, nil
}
