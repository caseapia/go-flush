package auth

import (
	"context"
	"errors"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/hash"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository mysql.Repository
}

func NewService(userRepo mysql.Repository) *Service {
	return &Service{repository: userRepo}
}

var ErrInvalidToken = &fiber.Error{Code: 400, Message: "invalid token"}

func (s *Service) Register(ctx context.Context, name, invite, email, password, ip string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		Password:     string(hash),
		TokenVersion: 1,
		IsVerified:   true,
		RegisterIP:   ip,
	}

	return user, s.repository.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, login, password, userAgent, ip string) (string, string, error) {
	user, err := s.repository.SearchByLogin(ctx, login)
	if err != nil {
		return "", "", fiber.NewError(401, "invalid credentials")
	}

	if !hash.CheckPassword(user.Password, password) {
		return "", "", fiber.NewError(401, "invalid credentials")
	}

	sessionID := uuid.NewString()
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	refreshHash := hash.HashToken(refreshToken)
	session := &models.Session{
		ID:          sessionID,
		UserID:      user.ID,
		RefreshHash: refreshHash,
		UserAgent:   userAgent,
		IPLast:      ip,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:   time.Now(),
	}

	if err := s.repository.CreateSession(ctx, session); err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, sessionID, user.TokenVersion)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	refreshHash := hash.HashToken(refreshToken)

	session, err := s.repository.GetSessionByHash(ctx, refreshHash)
	if err != nil || session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("invalid or expired session")
	}

	user, err := s.repository.SearchUserByID(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	session.RefreshHash = hash.HashToken(newRefreshToken)
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	if err := s.repository.UpdateSession(ctx, session); err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, session.ID, user.TokenVersion)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	return s.repository.RevokeSession(ctx, sessionID)
}

func (s *Service) ParseJWT(tokenString string) (*models.User, *utils.Claims, error) {
	claims, err := utils.ParseAccessToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.repository.SearchUserByID(context.Background(), claims.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		slog.WithData(slog.M{
			"error":  err,
			"user":   user,
			"claims": claims,
		}).Error("user seems to be nil on JWT Parsing")
		return nil, nil, errors.New("user not found")
	}

	if user.TokenVersion != claims.TokenVer {
		return nil, nil, errors.New("invalid token version")
	}

	return user, claims, nil
}
