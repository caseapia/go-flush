package user

import (
	"context"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s *UserService) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, UserError.UserNotFound()
	}

	_ = s.logger.Log(ctx, loggermodule.CommonLogger, 0, &id, loggermodule.SearchByUserID)

	return u, nil
}

func (s *UserService) GetUsersList(ctx context.Context) ([]models.User, error) {
	u, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}
