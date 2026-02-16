package models

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID            uint64       `bun:"id,pk,autoincrement,unique" json:"id"`
	Name          string       `bun:"name,unique,notnull" json:"name"`
	Email         string       `bun:"email" json:"-"`
	Password      string       `bun:"password" json:"-"`
	IsVerified    bool         `bun:"is_verified" json:"isVerified"`
	IsDeleted     bool         `bun:"is_deleted" json:"isDeleted,omitempty"`
	StaffRank     int          `bun:"staff_rank,default:1" json:"staffRank"`
	DeveloperRank int          `bun:"developer_rank,default:1" json:"developerRank"`
	Flags         []string     `bun:"staff_flags" json:"staffFlags"`
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

func (u *User) EditFlags(flags []string) (*User, error) {
	if u.IsDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	u.Flags = flags
	u.UpdatedAt = time.Now()

	return u, nil
}

func (u *User) UserHasFlag(flag string) bool {
	for _, f := range u.Flags {
		if f == "MANAGER" {
			return true
		}

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
