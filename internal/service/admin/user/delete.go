package adminUser

import (
	"context"
	"time"

	loggermodel "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	AdminErrorConstructor "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"

)

func (s *AdminUserService) DeleteUser(ctx context.Context, id uint64) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	r, err := s.rankService.GetByID(ctx, u.StaffRank)

	if err != nil && u == nil {
		return nil, err
	}

	if r.HasFlag("MANAGER") {
		_ = s.logger.Log(ctx, loggermodel.CommonLogger, 0, &id, loggermodel.TriedToDeleteManager)

		return nil, AdminErrorConstructor.CantDeleteManager()
	}

	if u.IsDeleted {
		_ = s.logger.Log(ctx, loggermodel.CommonLogger, 0, &id, loggermodel.HardDelete)

		if err := s.adminUser.HardDelete(ctx, id); err != nil {
			return nil, err
		}

		return nil, nil
	}

	_ = s.logger.Log(ctx, loggermodel.CommonLogger, 0, &id, loggermodel.SoftDelete)

	u.IsDeleted = true
	u.UpdatedAt = time.Now()

	if err := s.adminUser.SoftDelete(ctx, u); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AdminUserService) RestoreUser(ctx context.Context, id uint64) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil && u == nil {
		return nil, err
	}

	if !u.IsDeleted {
		return u, UserError.UserAlreadyExists()
	}

	_ = s.logger.Log(ctx, loggermodel.CommonLogger, 0, &id, loggermodel.RestoreUser)

	u.IsDeleted = false
	u.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
