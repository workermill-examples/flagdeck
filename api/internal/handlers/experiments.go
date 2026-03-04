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

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
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

// ListExperiments handles GET /api/v1/experiments
func (h *ExperimentsHandler) ListExperiments(c *fiber.Ctx) error {
	// Get pagination parameters
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get total count
	total, err := h.experimentsCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count experiments")
	}

	// Get experiments with pagination
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := h.experimentsCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve experiments")
	}
	defer cursor.Close(context.Background())

	var experiments []models.Experiment
	if err := cursor.All(context.Background(), &experiments); err != nil {
		return middleware.NewDatabaseError("Failed to decode experiments")
	}

	if experiments == nil {
		experiments = []models.Experiment{}
	}

	return c.JSON(ExperimentsListResponse{
		Data:  experiments,
		Total: total,
	})
}

// GetExperiment handles GET /api/v1/experiments/:key
func (h *ExperimentsHandler) GetExperiment(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return middleware.NewBadRequestError("Experiment key is required")
	}

	var experiment models.Experiment
	err := h.experimentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&experiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Experiment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve experiment")
	}

	return c.JSON(experiment)
}

// CreateExperiment handles POST /api/v1/experiments
func (h *ExperimentsHandler) CreateExperiment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)

	var req CreateExperimentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.Key == "" {
		return middleware.NewValidationError("Experiment key is required")
	}
	if req.Name == "" {
		return middleware.NewValidationError("Experiment name is required")
	}
	if req.FlagKey == "" {
		return middleware.NewValidationError("Flag key is required")
	}
	if req.Environment == "" {
		return middleware.NewValidationError("Environment is required")
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
		return middleware.NewValidationError("Status must be one of: draft, running, paused, completed")
	}

	// Check if experiment with this key already exists
	count, err := h.experimentsCollection.CountDocuments(context.Background(), bson.M{"key": req.Key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to check existing experiment")
	}
	if count > 0 {
		return middleware.NewConflictError("Experiment with this key already exists")
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

	now := time.Now()
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
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = h.experimentsCollection.InsertOne(context.Background(), experiment)
	if err != nil {
		return middleware.NewDatabaseError("Failed to create experiment")
	}

	// Log audit entry
	changes := map[string]interface{}{
		"key":         experiment.Key,
		"name":        experiment.Name,
		"flag_key":    experiment.FlagKey,
		"environment": experiment.Environment,
		"status":      experiment.Status,
	}
	h.auditService.LogCreate("experiment", experiment.Key, user.ID, user.Email, changes)

	return c.Status(fiber.StatusCreated).JSON(experiment)
}

// UpdateExperiment handles PUT /api/v1/experiments/:key
func (h *ExperimentsHandler) UpdateExperiment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Experiment key is required")
	}

	var req UpdateExperimentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Get existing experiment
	var existingExperiment models.Experiment
	err := h.experimentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingExperiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Experiment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve experiment")
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Name != nil {
		update["$set"].(bson.M)["name"] = *req.Name
	}
	if req.Description != nil {
		update["$set"].(bson.M)["description"] = *req.Description
	}
	if req.FlagKey != nil {
		update["$set"].(bson.M)["flag_key"] = *req.FlagKey
	}
	if req.Environment != nil {
		update["$set"].(bson.M)["environment"] = *req.Environment
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
			return middleware.NewValidationError("Status must be one of: draft, running, paused, completed")
		}
		update["$set"].(bson.M)["status"] = *req.Status
	}
	if req.StartDate != nil {
		update["$set"].(bson.M)["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		update["$set"].(bson.M)["end_date"] = *req.EndDate
	}
	if req.Variants != nil {
		update["$set"].(bson.M)["variants"] = req.Variants
	}
	if req.Results != nil {
		update["$set"].(bson.M)["results"] = req.Results
	}

	result, err := h.experimentsCollection.UpdateOne(
		context.Background(),
		bson.M{"key": key},
		update,
	)
	if err != nil {
		return middleware.NewDatabaseError("Failed to update experiment")
	}
	if result.MatchedCount == 0 {
		return middleware.NewNotFoundError("Experiment not found")
	}

	// Get updated experiment
	var updatedExperiment models.Experiment
	err = h.experimentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&updatedExperiment)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve updated experiment")
	}

	// Log audit entry
	beforeData := map[string]interface{}{
		"name":        existingExperiment.Name,
		"description": existingExperiment.Description,
		"flag_key":    existingExperiment.FlagKey,
		"environment": existingExperiment.Environment,
		"status":      existingExperiment.Status,
	}
	afterData := map[string]interface{}{
		"name":        updatedExperiment.Name,
		"description": updatedExperiment.Description,
		"flag_key":    updatedExperiment.FlagKey,
		"environment": updatedExperiment.Environment,
		"status":      updatedExperiment.Status,
	}
	h.auditService.LogUpdate("experiment", key, user.ID, user.Email, beforeData, afterData)

	return c.JSON(updatedExperiment)
}

// DeleteExperiment handles DELETE /api/v1/experiments/:key
func (h *ExperimentsHandler) DeleteExperiment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Experiment key is required")
	}

	// Get existing experiment for audit
	var existingExperiment models.Experiment
	err := h.experimentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingExperiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Experiment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve experiment")
	}

	result, err := h.experimentsCollection.DeleteOne(context.Background(), bson.M{"key": key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to delete experiment")
	}
	if result.DeletedCount == 0 {
		return middleware.NewNotFoundError("Experiment not found")
	}

	// Log audit entry
	deletedData := map[string]interface{}{
		"key":         existingExperiment.Key,
		"name":        existingExperiment.Name,
		"flag_key":    existingExperiment.FlagKey,
		"environment": existingExperiment.Environment,
		"status":      existingExperiment.Status,
	}
	h.auditService.LogDelete("experiment", key, user.ID, user.Email, deletedData)

	return c.SendStatus(fiber.StatusNoContent)
}

// TrackExperiment handles POST /api/v1/experiments/:key/track
func (h *ExperimentsHandler) TrackExperiment(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return middleware.NewBadRequestError("Experiment key is required")
	}

	var req TrackExperimentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.UserID == "" {
		return middleware.NewValidationError("User ID is required")
	}
	if req.VariantName == "" {
		return middleware.NewValidationError("Variant name is required")
	}

	// Check if experiment exists
	var experiment models.Experiment
	err := h.experimentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&experiment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Experiment not found")
		}
		return middleware.NewDatabaseError("Failed to fetch experiment")
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
		return middleware.NewBadRequestError("Variant not found in experiment")
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
		context.Background(),
		bson.M{"key": key},
		bson.M{
			"$set": bson.M{
				"results":    experiment.Results,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return middleware.NewDatabaseError("Failed to update experiment results")
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
