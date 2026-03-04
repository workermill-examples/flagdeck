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

type FlagsHandler struct {
	flagsCollection *mongo.Collection
	auditService    *services.AuditService
}

type FlagsListResponse struct {
	Data  []models.Flag `json:"data"`
	Total int64         `json:"total"`
}

type CreateFlagRequest struct {
	Key          string                            `json:"key" validate:"required"`
	Name         string                            `json:"name" validate:"required"`
	Description  string                            `json:"description"`
	Type         string                            `json:"type" validate:"required,oneof=boolean string number json"`
	Environments map[string]models.FlagEnvironment `json:"environments"`
}

type UpdateFlagRequest struct {
	Name         *string                           `json:"name,omitempty"`
	Description  *string                           `json:"description,omitempty"`
	Type         *string                           `json:"type,omitempty"`
	Environments map[string]models.FlagEnvironment `json:"environments,omitempty"`
}

type ToggleFlagRequest struct {
	Environment *string `json:"environment,omitempty"`
}

func NewFlagsHandler(flagsCollection *mongo.Collection, auditService *services.AuditService) *FlagsHandler {
	return &FlagsHandler{
		flagsCollection: flagsCollection,
		auditService:    auditService,
	}
}

// ListFlags handles GET /api/v1/flags
func (h *FlagsHandler) ListFlags(c *fiber.Ctx) error {
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
	total, err := h.flagsCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count flags")
	}

	// Get flags with pagination
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := h.flagsCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve flags")
	}
	defer cursor.Close(context.Background())

	var flags []models.Flag
	if err := cursor.All(context.Background(), &flags); err != nil {
		return middleware.NewDatabaseError("Failed to decode flags")
	}

	if flags == nil {
		flags = []models.Flag{}
	}

	return c.JSON(FlagsListResponse{
		Data:  flags,
		Total: total,
	})
}

// GetFlag handles GET /api/v1/flags/:key
func (h *FlagsHandler) GetFlag(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return middleware.NewBadRequestError("Flag key is required")
	}

	var flag models.Flag
	err := h.flagsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&flag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Flag not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve flag")
	}

	return c.JSON(flag)
}

// CreateFlag handles POST /api/v1/flags
func (h *FlagsHandler) CreateFlag(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)

	var req CreateFlagRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.Key == "" {
		return middleware.NewValidationError("Flag key is required")
	}
	if req.Name == "" {
		return middleware.NewValidationError("Flag name is required")
	}
	if req.Type == "" {
		return middleware.NewValidationError("Flag type is required")
	}

	// Validate type
	if req.Type != "boolean" && req.Type != "string" && req.Type != "number" && req.Type != "json" {
		return middleware.NewValidationError("Flag type must be one of: boolean, string, number, json")
	}

	// Check if flag with this key already exists
	count, err := h.flagsCollection.CountDocuments(context.Background(), bson.M{"key": req.Key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to check existing flag")
	}
	if count > 0 {
		return middleware.NewConflictError("Flag with this key already exists")
	}

	// Initialize environments if not provided
	if req.Environments == nil {
		req.Environments = make(map[string]models.FlagEnvironment)
	}

	now := time.Now()
	flag := models.Flag{
		ID:           primitive.NewObjectID(),
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		IsActive:     true,
		Type:         req.Type,
		Environments: req.Environments,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = h.flagsCollection.InsertOne(context.Background(), flag)
	if err != nil {
		return middleware.NewDatabaseError("Failed to create flag")
	}

	// Log audit entry
	changes := map[string]interface{}{
		"key":          flag.Key,
		"name":         flag.Name,
		"description":  flag.Description,
		"type":         flag.Type,
		"is_active":    flag.IsActive,
		"environments": flag.Environments,
	}
	h.auditService.LogCreate("flag", flag.Key, user.ID, user.Email, changes)

	return c.Status(fiber.StatusCreated).JSON(flag)
}

// UpdateFlag handles PUT /api/v1/flags/:key
func (h *FlagsHandler) UpdateFlag(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Flag key is required")
	}

	var req UpdateFlagRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Get existing flag
	var existingFlag models.Flag
	err := h.flagsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingFlag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Flag not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve flag")
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Name != nil {
		update["$set"].(bson.M)["name"] = *req.Name
	}
	if req.Description != nil {
		update["$set"].(bson.M)["description"] = *req.Description
	}
	if req.Type != nil {
		// Validate type
		if *req.Type != "boolean" && *req.Type != "string" && *req.Type != "number" && *req.Type != "json" {
			return middleware.NewValidationError("Flag type must be one of: boolean, string, number, json")
		}
		update["$set"].(bson.M)["type"] = *req.Type
	}
	if req.Environments != nil {
		update["$set"].(bson.M)["environments"] = req.Environments
	}

	result, err := h.flagsCollection.UpdateOne(
		context.Background(),
		bson.M{"key": key},
		update,
	)
	if err != nil {
		return middleware.NewDatabaseError("Failed to update flag")
	}
	if result.MatchedCount == 0 {
		return middleware.NewNotFoundError("Flag not found")
	}

	// Get updated flag
	var updatedFlag models.Flag
	err = h.flagsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&updatedFlag)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve updated flag")
	}

	// Log audit entry
	beforeData := map[string]interface{}{
		"name":         existingFlag.Name,
		"description":  existingFlag.Description,
		"type":         existingFlag.Type,
		"environments": existingFlag.Environments,
	}
	afterData := map[string]interface{}{
		"name":         updatedFlag.Name,
		"description":  updatedFlag.Description,
		"type":         updatedFlag.Type,
		"environments": updatedFlag.Environments,
	}
	h.auditService.LogUpdate("flag", key, user.ID, user.Email, beforeData, afterData)

	return c.JSON(updatedFlag)
}

// DeleteFlag handles DELETE /api/v1/flags/:key
func (h *FlagsHandler) DeleteFlag(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Flag key is required")
	}

	// Get existing flag for audit
	var existingFlag models.Flag
	err := h.flagsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingFlag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Flag not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve flag")
	}

	result, err := h.flagsCollection.DeleteOne(context.Background(), bson.M{"key": key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to delete flag")
	}
	if result.DeletedCount == 0 {
		return middleware.NewNotFoundError("Flag not found")
	}

	// Log audit entry
	deletedData := map[string]interface{}{
		"key":          existingFlag.Key,
		"name":         existingFlag.Name,
		"description":  existingFlag.Description,
		"type":         existingFlag.Type,
		"is_active":    existingFlag.IsActive,
		"environments": existingFlag.Environments,
	}
	h.auditService.LogDelete("flag", key, user.ID, user.Email, deletedData)

	return c.SendStatus(fiber.StatusNoContent)
}

// ToggleFlag handles POST /api/v1/flags/:key/toggle
func (h *FlagsHandler) ToggleFlag(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Flag key is required")
	}

	var req ToggleFlagRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Get existing flag
	var existingFlag models.Flag
	err := h.flagsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingFlag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Flag not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve flag")
	}

	var update bson.M
	var fromState, toState bool
	var environment string

	if req.Environment != nil && *req.Environment != "" {
		// Toggle environment-specific enabled flag
		environment = *req.Environment
		envConfig, exists := existingFlag.Environments[environment]
		if !exists {
			return middleware.NewNotFoundError("Environment not found in flag")
		}

		fromState = envConfig.Enabled
		toState = !fromState

		// Update the environment enabled state
		update = bson.M{
			"$set": bson.M{
				"environments." + environment + ".enabled": toState,
				"updated_at": time.Now(),
			},
		}
	} else {
		// Toggle global is_active flag
		fromState = existingFlag.IsActive
		toState = !fromState

		update = bson.M{
			"$set": bson.M{
				"is_active":  toState,
				"updated_at": time.Now(),
			},
		}
	}

	result, err := h.flagsCollection.UpdateOne(
		context.Background(),
		bson.M{"key": key},
		update,
	)
	if err != nil {
		return middleware.NewDatabaseError("Failed to toggle flag")
	}
	if result.MatchedCount == 0 {
		return middleware.NewNotFoundError("Flag not found")
	}

	// Get updated flag
	var updatedFlag models.Flag
	err = h.flagsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&updatedFlag)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve updated flag")
	}

	// Log audit entry
	h.auditService.LogToggle("flag", key, user.ID, user.Email, environment, fromState, toState)

	return c.JSON(updatedFlag)
}
