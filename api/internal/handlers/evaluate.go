package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/services"
)

type EvaluateHandler struct {
	evaluationService *services.EvaluationService
}

func NewEvaluateHandler(evaluationService *services.EvaluationService) *EvaluateHandler {
	return &EvaluateHandler{
		evaluationService: evaluationService,
	}
}

// EvaluateFlag handles POST /api/v1/evaluate
func (h *EvaluateHandler) EvaluateFlag(c *fiber.Ctx) error {
	var req services.EvaluationRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.FlagKey == "" {
		return middleware.NewValidationError("flag_key is required")
	}
	if req.Environment == "" {
		return middleware.NewValidationError("environment is required")
	}
	if req.UserID == "" {
		return middleware.NewValidationError("user_id is required")
	}

	result, err := h.evaluationService.EvaluateFlag(req)
	if err != nil {
		return middleware.NewInternalError("Failed to evaluate flag")
	}

	return c.JSON(result)
}

// EvaluateBulk handles POST /api/v1/evaluate/bulk
func (h *EvaluateHandler) EvaluateBulk(c *fiber.Ctx) error {
	var req services.BulkEvaluationRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.Environment == "" {
		return middleware.NewValidationError("environment is required")
	}
	if req.UserID == "" {
		return middleware.NewValidationError("user_id is required")
	}
	if len(req.FlagKeys) == 0 {
		return middleware.NewValidationError("flag_keys cannot be empty")
	}

	result, err := h.evaluationService.EvaluateFlags(req)
	if err != nil {
		return middleware.NewInternalError("Failed to evaluate flags")
	}

	return c.JSON(result)
}
