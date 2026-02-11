package adminUser

import (
	loggermodel "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
	"github.com/gofiber/fiber/v2"
)

func (s *AdminUserService) CreateUser(ctx *fiber.Ctx, adminID int, name string) (*models.User, error) {
	existing, err := s.repo.GetByName(ctx.UserContext(), name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, UserError.UserAlreadyExists()
	}

	if name == "" || len(name) < 3 || len(name) > 30 {
		return nil, UserError.UserInvalidUsername()
	}

	user := &models.User{
		Name: name,
	}

	if err := s.adminUser.Create(ctx.UserContext(), user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx.UserContext(), loggermodel.CommonLogger, uint64(adminID), nil, loggermodel.Create, "as "+name)

	return user, nil
}
