package adminUser

import (
	AdminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *AdminUserService.AdminUserService
}

func NewAdminUserHandler(service *AdminUserService.AdminUserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app fiber.Router) {
	user := app.Group("/admin/user")

	user.Delete("/delete/:id", h.DeleteUser)   // Delete user
	user.Put("/restore/:id", h.RestoreUser)    // Restore deleted user
	user.Put("/create", h.CreateUser)          // Create user
	user.Patch("/ban/:id", h.BanUser)          // Ban user
	user.Delete("/unban/:id", h.UnbanUser)     // Unban banned user
	user.Patch("/flags/edit/:id", h.EditFlags) // Edit user staff flags
}
