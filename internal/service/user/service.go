package user

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"golang.org/x/crypto/bcrypt"
)

type Logger interface {
	Log(ctx context.Context, loggerType models.LoggerType, adminID uint64, userID *uint64, action interface{}, additional ...string)
}

type Repository interface {
	SearchUserByID(ctx context.Context, id uint64) (*models.User, error)
	SearchAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	SearchUserByName(ctx context.Context, name string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	SoftDelete(ctx context.Context, u *models.User) error
	HardDelete(ctx context.Context, id uint64) error
	Restore(ctx context.Context, user *models.User) error
	CreateBan(ctx context.Context, ban *models.BanModel) error
	GetActiveBan(ctx context.Context, userID uint64) (*models.BanModelDTO, error)
	DeleteBan(ctx context.Context, userID uint64) error
	ChangeUserData(ctx context.Context, u *models.User, updateName, updateEmail, updatePassword bool) error

	SearchRankByID(ctx context.Context, id int) (*models.RankStructure, error)
	SetStaffRank(ctx context.Context, userID uint64, rankID int) (*models.User, error)
	SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*models.User, error)
}

type Service struct {
	repo   Repository
	logger Logger
}

func NewService(r Repository, l Logger) *Service {
	return &Service{
		repo:   r,
		logger: l,
	}
}

func (s *Service) SearchUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fiber.ErrNotFound
	}

	ban, _ := s.repo.GetActiveBan(ctx, id)
	user.ActiveBan = ban

	if id != adminID {
		s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.SearchByUserID)
	}

	return user, nil
}

func (s *Service) GetUsersList(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.SearchAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) GetOwnAccount(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &fiber.Error{Code: 401, Message: "not authorized to get their own info"}
	}

	ban, _ := s.repo.GetActiveBan(ctx, id)
	user.ActiveBan = ban

	return user, nil
}

// ! Admin actions
func (s *Service) BanUser(ctx context.Context, adminID, userID uint64, unbanDate time.Time, reason string) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user == nil || user.IsDeleted {
		slog.WithData(slog.M{
			"error": err,
			"user":  user,
		}).Error("error occured")
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if user.UserHasFlag("NONBANNABLE") {
		return nil, fiber.NewError(fiber.StatusForbidden, "ban of this user is not allowed")
	}

	ban := &models.BanModel{
		IssuedBy:       adminID,
		IssuedTo:       userID,
		Date:           time.Now(),
		ExpirationDate: unbanDate,
		Reason:         reason,
	}

	if err := s.repo.CreateBan(ctx, ban); err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("reason: %s\nuntil: %s", reason, unbanDate.String())
	s.logger.Log(ctx, models.PunishmentLogger, adminID, &userID, models.Ban, addInfo)

	user.ActiveBan = &models.BanModelDTO{
		BanModel: *ban,
	}
	return user, nil
}

func (s *Service) UnbanUser(ctx context.Context, adminID, userID uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user == nil || user.IsDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	activeBan, _ := s.repo.GetActiveBan(ctx, userID)
	if activeBan != nil {
		if err := s.repo.DeleteBan(ctx, userID); err != nil {
			return nil, err
		}
	}

	user.ActiveBanID = nil
	user.ActiveBan = nil

	s.logger.Log(ctx, models.PunishmentLogger, adminID, &userID, models.Unban)

	return user, nil
}

func (s *Service) CreateUser(ctx *fiber.Ctx, adminID uint64, name, email, password string) (*models.User, error) {
	existing, err := s.repo.SearchUserByName(ctx.UserContext(), name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, fiber.ErrBadRequest
	}

	if name == "" || len(name) < 3 || len(name) > 30 || len(password) < 6 {
		return nil, fiber.ErrBadRequest
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	if err := s.repo.CreateUser(ctx.UserContext(), user); err != nil {
		return nil, err
	}

	s.logger.Log(ctx.UserContext(), models.CommonLogger, uint64(adminID), nil, models.Create, "with nickname "+name)

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, id)
	r, err := s.repo.SearchRankByID(ctx, u.StaffRank)

	if err != nil && u == nil {
		return nil, err
	}

	if r.HasFlag("MANAGER") {
		s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.TriedToDeleteManager)

		return nil, fiber.ErrForbidden
	}

	if u.IsDeleted {
		s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.HardDelete)

		if err := s.repo.HardDelete(ctx, id); err != nil {
			return nil, err
		}

		return nil, nil
	}

	s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.SoftDelete)

	u.IsDeleted = true
	u.UpdatedAt = time.Now()

	if err := s.repo.SoftDelete(ctx, u); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) RestoreUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, id)
	if err != nil && u == nil {
		return nil, err
	}

	if !u.IsDeleted {
		return u, fiber.ErrBadRequest
	}

	s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.RestoreUser)

	u.IsDeleted = false
	u.UpdatedAt = time.Now()

	if err := s.repo.Restore(ctx, u); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) SetStaffRank(ctx context.Context, adminID uint64, userID uint64, rankID int) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	currentRank, _ := s.repo.SearchRankByID(ctx, u.StaffRank)
	oldRankName := "None"
	if currentRank != nil {
		oldRankName = fmt.Sprintf("%s (%d)", currentRank.Name, currentRank.ID)
	}

	newRank, err := s.repo.SearchRankByID(ctx, rankID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "target rank not found")
	}

	if newRank.HasFlag("DEV") {
		return nil, fiber.NewError(fiber.StatusForbidden, "developer rank cannot be issued here")
	}

	updatedUser, err := s.repo.SetStaffRank(ctx, userID, rankID)
	if err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("Before: %s\nAfter: %s (%d)", oldRankName, newRank.Name, newRank.ID)

	s.logger.Log(ctx, models.CommonLogger, adminID, &userID, models.SetStaffRank, addInfo)

	return updatedUser, nil
}

func (s *Service) SetDeveloperRank(ctx context.Context, adminID uint64, userId uint64, rankID int) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, userId)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	currentRank, _ := s.repo.SearchRankByID(ctx, u.DeveloperRank)
	oldRankInfo := "None"
	if currentRank != nil {
		oldRankInfo = fmt.Sprintf("%s (%d)", currentRank.Name, currentRank.ID)
	}

	r, err := s.repo.SearchRankByID(ctx, rankID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "target rank not found")
	}

	if !r.HasFlag("DEV") && r.Name != "None" && r.Name != "Player" {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userId,
		}).Error("Attempt to set non-DEV rank via SetDeveloperRank")

		return nil, fiber.NewError(fiber.StatusForbidden, "this function is only for developer ranks")
	}

	setRank, err := s.repo.SetDeveloperRank(ctx, userId, rankID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUser(ctx, setRank); err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("Before: %s\nAfter: %s (%d)", oldRankInfo, r.Name, r.ID)

	s.logger.Log(ctx, models.CommonLogger, adminID, &userId, models.SetDeveloperRank, addInfo)

	return setRank, nil
}

func (s *Service) ChangeUser(ctx context.Context, adminID uint64, userID uint64, name *string, email *string, password *string) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	var hash []byte
	if password != nil {
		var genErr error
		hash, genErr = bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if genErr != nil {
			return nil, genErr
		}
		u.Password = string(hash)
	}
	oldInfo := fmt.Sprintf("Name: %s, Email: %s", u.Name, u.Email)

	if name != nil {
		u.Name = *name
	}
	if email != nil {
		u.Email = *email
	}
	if password != nil {
		u.Password = string(hash)
	}

	err = s.repo.ChangeUserData(ctx, u, name != nil, email != nil, hash != nil)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newInfo := fmt.Sprintf("Name: %s, Email: %s", u.Name, u.Email)

	if password == nil {
		addInfo := fmt.Sprintf("Before: %s\nAfter: %s", oldInfo, newInfo)
		s.logger.Log(ctx, models.CommonLogger, adminID, &userID, models.ChangeUserData, addInfo)
	} else {
		s.logger.Log(ctx, models.CommonLogger, adminID, &userID, models.ChangeUserPassword)
	}

	return u, nil
}
