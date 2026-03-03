package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/workermill-examples/flagdeck/api/internal/models"
	"github.com/workermill-examples/flagdeck/api/internal/services"
)

type ExperimentsHandler struct {
	experimentsCollection *mongo.Collection
	auditService          *services.AuditService
}

type ExperimentsListResponse struct {
	Data  []models.Experiment `json:"data"`
	Total int64               `json:"total"`
}

type CreateExperimentRequest struct {
	Key         string                     `json:"key" validate:"required"`
	Name        string                     `json:"name" validate:"required"`
	Description string                     `json:"description"`
	FlagKey     string                     `json:"flag_key" validate:"required"`
	Environment string                     `json:"environment" validate:"required"`
	Status      string                     `json:"status" validate:"required,oneof=draft running paused completed"`
	StartDate   *time.Time                 `json:"start_date"`
	EndDate     *time.Time                 `json:"end_date"`
	Variants    []models.ExperimentVariant `json:"variants"`
	Results     map[string]interface{}     `json:"results"`
}

type UpdateExperimentRequest struct {
	Name        *string                    `json:"name,omitempty"`
	Description *string                    `json:"description,omitempty"`
	FlagKey     *string                    `json:"flag_key,omitempty"`
	Environment *string                    `json:"environment,omitempty"`
	Status      *string                    `json:"status,omitempty"`
	StartDate   *time.Time                 `json:"start_date,omitempty"`
	EndDate     *time.Time                 `json:"end_date,omitempty"`
	Variants    []models.ExperimentVariant `json:"variants,omitempty"`
	Results     map[string]interface{}     `json:"results,omitempty"`
}

type TrackExperimentRequest struct {
	UserID      string  `json:"user_id" validate:"required"`
	VariantName string  `json:"variant_name" validate:"required"`
	Conversion  bool    `json:"conversion"`
	Revenue     float64 `json:"revenue"`
}

func NewExperimentsHandler(experimentsCollection *mongo.Collection, auditService *services.AuditService) *ExperimentsHandler {
	return &ExperimentsHandler{
		experimentsCollection: experimentsCollection,
		auditService:          auditService,
	}
}

// GET /api/v1/experiments
func (h *ExperimentsHandler) ListExperiments(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if limit < 1 || limit > 100 {
		limit = 50
	}
	skip := (page - 1) * limit

	// Get total count
	totalCount, err := h.experimentsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to count experiments",
			},
		})
	}

	// Get experiments with pagination
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := h.experimentsCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch experiments",
			},
		})
	}
	defer cursor.Close(ctx)

	var experiments []models.Experiment
	if err = cursor.All(ctx, &experiments); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to decode experiments",
			},
		})
	}

	if experiments == nil {
		experiments = []models.Experiment{}
	}

	response := ExperimentsListResponse{
		Data:  experiments,
		Total: totalCount,
	}

	return c.JSON(response)
}

// GET /api/v1/experiments/:key
func (h *ExperimentsHandler) GetExperiment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Experiment key is required",
			},
		})
	}

	var experiment models.Experiment
	err := h.experimentsCollection.FindOne(ctx, bson.M{"key": key}).Decode(&experiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "EXPERIMENT_NOT_FOUND",
					"message": "Experiment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch experiment",
			},
		})
	}

	return c.JSON(experiment)
}

// POST /api/v1/experiments
func (h *ExperimentsHandler) CreateExperiment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req CreateExperimentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_JSON",
				"message": "Invalid JSON in request body",
			},
		})
	}

	// Validate required fields
	if req.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Key is required",
			},
		})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Name is required",
			},
		})
	}
	if req.FlagKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Flag key is required",
			},
		})
	}
	if req.Environment == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Environment is required",
			},
		})
	}

	// Validate status
	validStatuses := []string{"draft", "running", "paused", "completed"}
	validStatus := false
	for _, status := range validStatuses {
		if req.Status == status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Status must be one of: draft, running, paused, completed",
			},
		})
	}

	// Check if experiment key already exists
	count, err := h.experimentsCollection.CountDocuments(ctx, bson.M{"key": req.Key})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to check experiment existence",
			},
		})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "EXPERIMENT_EXISTS",
				"message": "Experiment with this key already exists",
			},
		})
	}

	// Initialize variants and results if not provided
	variants := req.Variants
	if variants == nil {
		variants = []models.ExperimentVariant{}
	}

	results := req.Results
	if results == nil {
		results = make(map[string]interface{})
	}

	experiment := models.Experiment{
		ID:          primitive.NewObjectID(),
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		FlagKey:     req.FlagKey,
		Environment: req.Environment,
		Status:      req.Status,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Variants:    variants,
		Results:     results,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = h.experimentsCollection.InsertOne(ctx, experiment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create experiment",
			},
		})
	}

	// Log audit entry
	user := c.Locals("user").(*models.User)
	h.auditService.LogCreate("experiment", experiment.Key, user.ID, user.Email, map[string]interface{}{
		"name":        experiment.Name,
		"flag_key":    experiment.FlagKey,
		"environment": experiment.Environment,
		"status":      experiment.Status,
	})

	return c.Status(fiber.StatusCreated).JSON(experiment)
}

// PUT /api/v1/experiments/:key
func (h *ExperimentsHandler) UpdateExperiment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Experiment key is required",
			},
		})
	}

	var req UpdateExperimentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_JSON",
				"message": "Invalid JSON in request body",
			},
		})
	}

	// Check if experiment exists
	var existingExperiment models.Experiment
	err := h.experimentsCollection.FindOne(ctx, bson.M{"key": key}).Decode(&existingExperiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "EXPERIMENT_NOT_FOUND",
					"message": "Experiment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch experiment",
			},
		})
	}

	// Build update document
	updateDoc := bson.M{"updated_at": time.Now()}

	if req.Name != nil {
		updateDoc["name"] = *req.Name
	}
	if req.Description != nil {
		updateDoc["description"] = *req.Description
	}
	if req.FlagKey != nil {
		updateDoc["flag_key"] = *req.FlagKey
	}
	if req.Environment != nil {
		updateDoc["environment"] = *req.Environment
	}
	if req.Status != nil {
		// Validate status
		validStatuses := []string{"draft", "running", "paused", "completed"}
		validStatus := false
		for _, status := range validStatuses {
			if *req.Status == status {
				validStatus = true
				break
			}
		}
		if !validStatus {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "VALIDATION_ERROR",
					"message": "Status must be one of: draft, running, paused, completed",
				},
			})
		}
		updateDoc["status"] = *req.Status
	}
	if req.StartDate != nil {
		updateDoc["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updateDoc["end_date"] = *req.EndDate
	}
	if req.Variants != nil {
		updateDoc["variants"] = req.Variants
	}
	if req.Results != nil {
		updateDoc["results"] = req.Results
	}

	_, err = h.experimentsCollection.UpdateOne(ctx, bson.M{"key": key}, bson.M{"$set": updateDoc})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update experiment",
			},
		})
	}

	// Get updated experiment
	var updatedExperiment models.Experiment
	err = h.experimentsCollection.FindOne(ctx, bson.M{"key": key}).Decode(&updatedExperiment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch updated experiment",
			},
		})
	}

	// Log audit entry
	user := c.Locals("user").(*models.User)
	changes := make(map[string]interface{})
	for field, value := range updateDoc {
		if field != "updated_at" {
			changes[field] = value
		}
	}
	h.auditService.LogUpdate("experiment", key, user.ID, user.Email, map[string]interface{}{}, changes)

	return c.JSON(updatedExperiment)
}

// DELETE /api/v1/experiments/:key
func (h *ExperimentsHandler) DeleteExperiment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Experiment key is required",
			},
		})
	}

	// Check if experiment exists
	var experiment models.Experiment
	err := h.experimentsCollection.FindOne(ctx, bson.M{"key": key}).Decode(&experiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "EXPERIMENT_NOT_FOUND",
					"message": "Experiment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch experiment",
			},
		})
	}

	_, err = h.experimentsCollection.DeleteOne(ctx, bson.M{"key": key})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete experiment",
			},
		})
	}

	// Log audit entry
	user := c.Locals("user").(*models.User)
	h.auditService.LogDelete("experiment", key, user.ID, user.Email, map[string]interface{}{
		"name":        experiment.Name,
		"flag_key":    experiment.FlagKey,
		"environment": experiment.Environment,
	})

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// POST /api/v1/experiments/:key/track
func (h *ExperimentsHandler) TrackExperiment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Experiment key is required",
			},
		})
	}

	var req TrackExperimentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_JSON",
				"message": "Invalid JSON in request body",
			},
		})
	}

	// Validate required fields
	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "User ID is required",
			},
		})
	}
	if req.VariantName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Variant name is required",
			},
		})
	}

	// Check if experiment exists
	var experiment models.Experiment
	err := h.experimentsCollection.FindOne(ctx, bson.M{"key": key}).Decode(&experiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "EXPERIMENT_NOT_FOUND",
					"message": "Experiment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch experiment",
			},
		})
	}

	// Verify variant exists
	variantExists := false
	for _, variant := range experiment.Variants {
		if variant.Name == req.VariantName {
			variantExists = true
			break
		}
	}
	if !variantExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VARIANT_NOT_FOUND",
				"message": "Variant not found in experiment",
			},
		})
	}

	// Initialize results if not present
	if experiment.Results == nil {
		experiment.Results = make(map[string]interface{})
	}

	// Get or initialize variant results
	variantResultsInterface, exists := experiment.Results[req.VariantName]
	var variantResults map[string]interface{}
	if exists {
		if vr, ok := variantResultsInterface.(map[string]interface{}); ok {
			variantResults = vr
		} else {
			variantResults = make(map[string]interface{})
		}
	} else {
		variantResults = make(map[string]interface{})
	}

	// Update impression count
	impressions, _ := variantResults["impressions"].(float64)
	variantResults["impressions"] = impressions + 1

	// Update conversion count if conversion occurred
	if req.Conversion {
		conversions, _ := variantResults["conversions"].(float64)
		variantResults["conversions"] = conversions + 1
	}

	// Update revenue
	revenue, _ := variantResults["revenue"].(float64)
	variantResults["revenue"] = revenue + req.Revenue

	// Save updated results
	experiment.Results[req.VariantName] = variantResults

	// Update experiment in database
	_, err = h.experimentsCollection.UpdateOne(
		ctx,
		bson.M{"key": key},
		bson.M{
			"$set": bson.M{
				"results":    experiment.Results,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update experiment results",
			},
		})
	}

	// Return tracking confirmation
	return c.JSON(fiber.Map{
		"tracked":    true,
		"experiment": key,
		"variant":    req.VariantName,
		"user_id":    req.UserID,
		"conversion": req.Conversion,
		"revenue":    req.Revenue,
		"tracked_at": time.Now(),
	})
}
