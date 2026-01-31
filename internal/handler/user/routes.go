package handler

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router, h *UserHandler) {
	users := app.Group("/users")
	user := app.Group("/user")

	users.Get("/", h.GetUsersList)
	user.Get("/:id", h.GetUser)
	user.Put("/ban/:id", h.BanUser)
	user.Delete("/unban/:id", h.UnbanUser)
	user.Put("/create/", h.CreateUser)
	user.Delete("/delete/:id", h.DeleteUser)
	user.Post("/restore/:id", h.RestoreUser)
}
