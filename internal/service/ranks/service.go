package ranks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Logger interface {
	Log(ctx context.Context, loggerType models.LoggerType, adminID uint64, userID *uint64, action interface{}, additional ...string) error
}

type Repository interface {
	SearchUserByID(ctx context.Context, id uint64) (*models.User, error)
	SearchAllRanks(ctx context.Context) ([]models.RankStructure, error)
	SearchRankByID(ctx context.Context, id int) (*models.RankStructure, error)
	SearchRankByName(ctx context.Context, rankName string) (*models.RankStructure, error)
	CreateRank(ctx context.Context, rank *models.RankStructure) error
	DeleteRank(ctx context.Context, rank *models.RankStructure) error
	EditRank(ctx context.Context, rank *models.RankStructure) (*models.RankStructure, error)
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

func (s *Service) CreateRank(ctx *fiber.Ctx, adminID uint64, rankName string, rankColor string, rankFlags []string) (*models.RankStructure, error) {
	u, err := s.repo.SearchUserByID(ctx.UserContext(), adminID)
	if err != nil {
		return nil, err
	}

	r, err := s.repo.SearchRankByID(ctx.UserContext(), int(u.StaffRank))
	if err != nil {
		return nil, err
	}

	if !r.HasFlag("STAFFMANAGEMENT") {
		return nil, &fiber.Error{Code: 403, Message: "you're not allowed to use this function"}
	}

	existing, err := s.repo.SearchRankByName(ctx.UserContext(), rankName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, &fiber.Error{Code: 409, Message: "rank with that name already exists"}
	}

	if rankName == "" || len(rankName) < 3 || len(rankName) > 30 {
		return nil, &fiber.Error{Code: 400, Message: "invalid length of rank name"}
	}

	rank := &models.RankStructure{Name: rankName, Color: rankColor, Flags: rankFlags}

	if err := s.repo.CreateRank(ctx.UserContext(), rank); err != nil {
		return nil, err
	}

	addInfo := "with name: " + rankName + ", with color: " + rankColor + "with flags: " + strings.Join(rankFlags, ", ")

	_ = s.logger.Log(ctx.UserContext(), models.CommonLogger, adminID, nil, models.CreateRank, addInfo)

	return rank, nil
}

func (s *Service) DeleteRank(ctx *fiber.Ctx, adminID uint64, id int) (bool, error) {
	u, err := s.repo.SearchUserByID(ctx.UserContext(), adminID)
	if err != nil {
		return false, err
	}
	uRank, err := s.repo.SearchRankByID(ctx.UserContext(), u.StaffRank)
	if err != nil {
		return false, err
	}

	if !uRank.HasFlag("STAFFMANAGEMENT") {
		return false, &fiber.Error{Code: 403, Message: "you're not allowed to use this function"}
	}

	r, err := s.repo.SearchRankByID(ctx.UserContext(), id)
	if err != nil {
		return false, err
	}
	if r == nil {
		return false, &fiber.Error{Code: 404, Message: "rank with that name not found"}
	}

	if err := s.repo.DeleteRank(ctx.UserContext(), r); err != nil {
		return false, err
	}

	addInfo := "with ID: " + strconv.FormatInt(r.ID, 10) + ", with name: " + r.Name

	_ = s.logger.Log(ctx.UserContext(), models.CommonLogger, 0, nil, models.DeleteRank, addInfo)

	return true, nil
}

func (s *Service) SearchAllRanks(ctx *fiber.Ctx) ([]models.RankStructure, error) {
	ranks, err := s.repo.SearchAllRanks(ctx.UserContext())
	if err != nil {
		return nil, err
	}

	return ranks, nil
}

func (s *Service) SearchRankByID(ctx *fiber.Ctx, id int) (*models.RankStructure, error) {
	rank, err := s.repo.SearchRankByID(ctx.UserContext(), id)
	if err != nil {
		return nil, err
	}

	return rank, nil
}

func (s *Service) EditRank(ctx context.Context, sender uint64, rank *models.RankStructure) (*models.RankStructure, error) {
	oldRank, err := s.repo.SearchRankByID(ctx, int(rank.ID))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "rank not found")
	}

	oldFlags := strings.Join(oldRank.Flags, ", ")
	oldInfo := fmt.Sprintf("Name: %s, Color: %s, Flags: %v", oldRank.Name, oldRank.Color, oldFlags)

	updatedRank, err := s.repo.EditRank(ctx, rank)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: err.Error()}
	}

	newFlags := strings.Join(updatedRank.Flags, ", ")
	newInfo := fmt.Sprintf("Name: %s, Color: %s, Flags: %v", updatedRank.Name, updatedRank.Color, newFlags)

	addInfo := "Before: " + oldInfo + "\nAfter: " + newInfo
	_ = s.logger.Log(ctx, models.CommonLogger, sender, nil, models.EditRank, addInfo)

	return updatedRank, nil
}
