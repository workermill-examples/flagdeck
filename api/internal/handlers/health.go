package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/workermill-examples/flagdeck/api/internal/database"
)

type HealthHandler struct {
	MongoDB *database.MongoDB
	Redis   *database.RedisDB
}

type HealthResponse struct {
	Status  string `json:"status"`
	MongoDB string `json:"mongodb"`
	Redis   string `json:"redis"`
}

func NewHealthHandler(mongoDB *database.MongoDB, redis *database.RedisDB) *HealthHandler {
	return &HealthHandler{
		MongoDB: mongoDB,
		Redis:   redis,
	}
}

func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:  "ok",
		MongoDB: "connected",
		Redis:   "connected",
	}

	statusCode := fiber.StatusOK

	// Check MongoDB connection
	if err := h.MongoDB.Ping(); err != nil {
		response.MongoDB = "disconnected"
		response.Status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	// Check Redis connection
	if err := h.Redis.Ping(); err != nil {
		response.Redis = "disconnected"
		response.Status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	// If both are down, status should be "down"
	if response.MongoDB == "disconnected" && response.Redis == "disconnected" {
		response.Status = "down"
	}

	return c.Status(statusCode).JSON(response)
}
