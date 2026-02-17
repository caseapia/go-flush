package models

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:flushproject.users"`

	ID            uint64       `bun:"id,pk,autoincrement,unique" json:"id"`
	Name          string       `bun:"name,unique,notnull" json:"name"`
	Email         string       `bun:"email" json:"-"`
	Password      string       `bun:"password" json:"-"`
	IsVerified    bool         `bun:"is_verified" json:"isVerified"`
	IsDeleted     bool         `bun:"is_deleted" json:"isDeleted,omitempty"`
	StaffRank     int          `bun:"staff_rank,default:1" json:"staffRank"`
	DeveloperRank int          `bun:"developer_rank,default:1" json:"developerRank"`
	Flags         *[]string    `bun:"staff_flags" json:"staffFlags"`
	CreatedAt     time.Time    `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time    `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt     *time.Time   `bun:"deleted_at,nullzero" json:"-"`
	TokenVersion  int          `bun:"token_version" json:"-"`
	ActiveBanID   *uint64      `bun:"active_ban" json:"-"`
	ActiveBan     *BanModelDTO `bun:"rel:has-one,join:active_ban=id" json:"activeBan"`
	LastLogin     *time.Time   `bun:"last_login" json:"lastLogin"`
	RegisterIP    string       `bun:"register_ip" json:"-"`
	LastIP        string       `bun:"last_ip" json:"-"`
}

type BanRequest struct {
	UnbanDate time.Time `json:"unbanDate"`
	Reason    string    `json:"reason"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RankSetterRequest struct {
	Status int `json:"status"`
}

type ChangeUserDataRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type EditUserFlagsRequest struct {
	NewFlags []string `json:"flags"`
}

func (u *User) SetStaffRank(rank int) (*User, error) {
	if u.IsDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	u.StaffRank = rank
	u.UpdatedAt = time.Now()
	return u, nil
}

func (u *User) SetDeveloperRank(rank int) (*User, error) {
	if u.IsDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	u.DeveloperRank = rank
	u.UpdatedAt = time.Now()
	return u, nil
}

func (u *User) UserHasFlag(flag string) bool {
	if u.Flags == nil {
		return false
	}

	for _, f := range *u.Flags {
		if f == flag {
			return true
		}
	}

	return false
}

func (u *User) GetPrivateData() map[string]interface{} {
	return map[string]interface{}{
		"email":      u.Email,
		"registerIP": u.RegisterIP,
		"lastIP":     u.LastIP,
	}
}
