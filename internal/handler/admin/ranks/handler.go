package adminRanks

import (
	adminRanks "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *adminRanks.RanksService
}

func NewHandler(s *adminRanks.RanksService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(app fiber.Router) {
	ranks := app.Group("/admin/rank")

	ranks.Get("/list", h.GetRanksList)  // Get ranks list
	ranks.Post("/create", h.CreateRank) // Create rank

	// TODO: Ranks deletion & Rank flags edit
}
