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
		ColumnExpr("user.*").
		ColumnExpr("ab.id AS active_ban__id").
		ColumnExpr("ab.issued_by AS active_ban__issued_by").
		ColumnExpr("ab.issued_to AS active_ban__issued_to").
		ColumnExpr("ab.reason AS active_ban__reason").
		ColumnExpr("ab.date AS active_ban__date").
		ColumnExpr("ab.expiration_date AS active_ban__expiration_date").
		ColumnExpr("adm.name AS active_ban__admin_name").
		Join("LEFT JOIN bans AS ab ON ab.id = user.active_ban AND ab.expiration_date > NOW()").
		Join("LEFT JOIN users AS adm ON adm.id = ab.issued_by").
		Where("user.id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if u.ActiveBan != nil && u.ActiveBan.ID == 0 {
		u.ActiveBan = nil
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

	if u.ActiveBanID != nil {
		banInfo, _ := r.GetActiveBan(ctx, u.ID)
		if banInfo != nil {
			u.ActiveBan = banInfo
		}
	}

	return u, nil
}

func (r *Repository) SearchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.NewSelect().
		Model(&users).
		ColumnExpr("user.*").
		ColumnExpr("ab.id AS active_ban__id").
		ColumnExpr("ab.reason AS active_ban__reason").
		ColumnExpr("ab.expiration_date AS active_ban__expiration_date").
		ColumnExpr("ab.issued_by AS active_ban__issued_by").
		ColumnExpr("adm.name AS active_ban__admin_name").
		Join("LEFT JOIN bans AS ab ON ab.id = user.active_ban").
		Join("LEFT JOIN users AS adm ON adm.id = ab.issued_by").
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

func (r *Repository) ChangeUserData(ctx context.Context, u *models.User, updateName, updateEmail, updatePassword bool) error {
	var conditions []string
	if updateName {
		conditions = append(conditions, "name")
	}
	if updateEmail {
		conditions = append(conditions, "email")
	}
	if updatePassword {
		conditions = append(conditions, "password")
	}

	for _, col := range conditions {
		val := ""
		if col == "name" {
			val = u.Name
		}
		if col == "email" {
			val = u.Email
		}
		if col == "password" {
			val = u.Password
		}

		exists, err := r.db.NewSelect().
			Model((*models.User)(nil)).
			Where("? = ?", bun.Ident(col), val).
			Where("id != ?", u.ID).
			Exists(ctx)

		if err != nil {
			return err
		}
		if exists && col != "password" {
			return fiber.NewError(fiber.StatusConflict, "User with this "+col+" already exists")
		}
	}

	query := r.db.NewUpdate().Model(u).WherePK()

	if updateName {
		query.Column("name")
	}
	if updateEmail {
		query.Column("email")
	}
	if updatePassword {
		query.Column("password")
	}

	if !updateName && !updateEmail && !updatePassword {
		return nil
	}

	_, err := query.Exec(ctx)
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

func (r *Repository) GetActiveBan(ctx context.Context, userID uint64) (*models.BanModelDTO, error) {
	u := new(models.BanModelDTO)

	err := r.db.NewSelect().
		Model(u).
		TableExpr("bans AS b").
		ColumnExpr("b.id AS id").
		ColumnExpr("b.issued_by AS issued_by").
		ColumnExpr("b.issued_to AS issued_to").
		ColumnExpr("b.reason AS reason").
		ColumnExpr("b.date AS date").
		ColumnExpr("b.expiration_date AS expiration_date").
		ColumnExpr("admin.name AS admin_name").
		Join("LEFT JOIN users AS admin ON admin.id = b.issued_by").
		Where("b.issued_to = ? AND b.expiration_date > NOW()", userID).
		Order("b.date DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
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
