package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func (r *Repository) SearchUserByID(ctx context.Context, id uint64) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) SearchUserByName(ctx context.Context, name string) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("name = ?", name).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) SearchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.NewSelect().
		Model(&users).
		Scan(ctx)

	return users, err
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

// ! Admin actions
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

func (r *Repository) SoftDelete(ctx context.Context, u *models.User) error {
	u.Name = u.Name + "_old"

	_, err := r.db.NewUpdate().
		Model(u).
		Column("name").
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) ChangeUserData(ctx context.Context, u *models.User, updateName bool, updateEmail bool) error {
	if updateName && updateEmail {
		return nil
	}

	var conditions []string
	if updateName {
		conditions = append(conditions, "name")
	}
	if updateEmail {
		conditions = append(conditions, "email")
	}

	for _, col := range conditions {
		val := ""
		if col == "name" {
			val = u.Name
		}
		if col == "email" {
			val = u.Email
		}

		exists, err := r.db.NewSelect().
			Model((*models.User)(nil)).
			Where("? = ?", bun.Ident(col), val).
			Where("id != ?", u.ID).
			Exists(ctx)

		if err != nil {
			return err
		}
		if exists {
			fiber.NewError(fiber.StatusConflict, "User with this "+col+" already exists")
		}
	}

	query := r.db.NewUpdate().Model(u).WherePK()

	if updateName {
		query.Column("name")
	}
	if updateEmail {
		query.Column("email")
	}

	if !updateName && !updateEmail {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *Repository) ChangePassword(ctx context.Context, u *models.User, newPassword string) error {
	u.Password = newPassword

	_, err := r.db.NewUpdate().
		Model(u).
		Column("password").
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) HardDelete(ctx context.Context, id uint64) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *Repository) Restore(ctx context.Context, user *models.User) error {
	user.Name = strings.ReplaceAll(user.Name, "_old", "")

	_, err := r.db.NewUpdate().
		Model(user).
		Column("name").
		WherePK().
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			return &fiber.Error{Code: 403, Message: "cannot restore: nickname is already taken"}
		}
	}
	return err
}

func (r *Repository) SetStaffRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("staff_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("developer_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) CreateBan(ctx context.Context, ban *models.BanModel) error {
	_, err := r.db.NewInsert().
		Model(ban).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewUpdate().
		Model(&models.User{}).
		Set("active_ban = ?", ban.ID).
		Where("id = ?", ban.IssuedTo).
		Exec(ctx)
	return err
}

func (r *Repository) GetActiveBan(ctx context.Context, userID uint64) (*models.BanModel, error) {
	var ban models.BanModel

	err := r.db.NewSelect().
		Model(&ban).
		Where("issued_to = ? AND expiration_date > NOW()", userID).
		Order("date DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &fiber.Error{Code: 404, Message: "no rows"}
		}
		return nil, err
	}

	return &ban, nil
}

func (r *Repository) DeleteBan(ctx context.Context, userID uint64) error {
	_, err := r.db.NewDelete().
		Model(&models.BanModel{}).
		Where("issued_to = ?", userID).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewUpdate().
		Table("users").
		Set("active_ban = NULL").
		Where("id = ?", userID).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateLastLogin(ctx *fiber.Ctx, userID uint64) error {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("last_login = ?", time.Now()).
		Set("last_ip = ?", ctx.IP()).
		Where("id = ?", userID).
		Exec(ctx.UserContext())
	return err
}
