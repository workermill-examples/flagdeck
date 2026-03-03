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

type EnvironmentsHandler struct {
	environmentsCollection *mongo.Collection
	auditService           *services.AuditService
}

type EnvironmentsListResponse struct {
	Data  []models.Environment `json:"data"`
	Total int64                `json:"total"`
}

type CreateEnvironmentRequest struct {
	Key         string `json:"key" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateEnvironmentRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

func NewEnvironmentsHandler(environmentsCollection *mongo.Collection, auditService *services.AuditService) *EnvironmentsHandler {
	return &EnvironmentsHandler{
		environmentsCollection: environmentsCollection,
		auditService:           auditService,
	}
}

// ListEnvironments handles GET /api/v1/environments
func (h *EnvironmentsHandler) ListEnvironments(c *fiber.Ctx) error {
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
	total, err := h.environmentsCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count environments")
	}

	// Get environments with pagination
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := h.environmentsCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve environments")
	}
	defer cursor.Close(context.Background())

	var environments []models.Environment
	if err := cursor.All(context.Background(), &environments); err != nil {
		return middleware.NewDatabaseError("Failed to decode environments")
	}

	if environments == nil {
		environments = []models.Environment{}
	}

	return c.JSON(EnvironmentsListResponse{
		Data:  environments,
		Total: total,
	})
}

// GetEnvironment handles GET /api/v1/environments/:key
func (h *EnvironmentsHandler) GetEnvironment(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return middleware.NewBadRequestError("Environment key is required")
	}

	var environment models.Environment
	err := h.environmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&environment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Environment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve environment")
	}

	return c.JSON(environment)
}

// CreateEnvironment handles POST /api/v1/environments
func (h *EnvironmentsHandler) CreateEnvironment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)

	var req CreateEnvironmentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.Key == "" {
		return middleware.NewValidationError("Environment key is required")
	}
	if req.Name == "" {
		return middleware.NewValidationError("Environment name is required")
	}

	// Check if environment with this key already exists
	count, err := h.environmentsCollection.CountDocuments(context.Background(), bson.M{"key": req.Key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to check existing environment")
	}
	if count > 0 {
		return middleware.NewConflictError("Environment with this key already exists")
	}

	now := time.Now()
	environment := models.Environment{
		ID:          primitive.NewObjectID(),
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = h.environmentsCollection.InsertOne(context.Background(), environment)
	if err != nil {
		return middleware.NewDatabaseError("Failed to create environment")
	}

	// Log audit entry
	changes := map[string]interface{}{
		"key":         environment.Key,
		"name":        environment.Name,
		"description": environment.Description,
		"color":       environment.Color,
	}
	h.auditService.LogCreate("environment", environment.Key, user.ID, user.Email, changes)

	return c.Status(fiber.StatusCreated).JSON(environment)
}

// UpdateEnvironment handles PUT /api/v1/environments/:key
func (h *EnvironmentsHandler) UpdateEnvironment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Environment key is required")
	}

	var req UpdateEnvironmentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Get existing environment
	var existingEnvironment models.Environment
	err := h.environmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingEnvironment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Environment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve environment")
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Name != nil {
		update["$set"].(bson.M)["name"] = *req.Name
	}
	if req.Description != nil {
		update["$set"].(bson.M)["description"] = *req.Description
	}
	if req.Color != nil {
		update["$set"].(bson.M)["color"] = *req.Color
	}

	result, err := h.environmentsCollection.UpdateOne(
		context.Background(),
		bson.M{"key": key},
		update,
	)
	if err != nil {
		return middleware.NewDatabaseError("Failed to update environment")
	}
	if result.MatchedCount == 0 {
		return middleware.NewNotFoundError("Environment not found")
	}

	// Get updated environment
	var updatedEnvironment models.Environment
	err = h.environmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&updatedEnvironment)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve updated environment")
	}

	// Log audit entry
	beforeData := map[string]interface{}{
		"name":        existingEnvironment.Name,
		"description": existingEnvironment.Description,
		"color":       existingEnvironment.Color,
	}
	afterData := map[string]interface{}{
		"name":        updatedEnvironment.Name,
		"description": updatedEnvironment.Description,
		"color":       updatedEnvironment.Color,
	}
	h.auditService.LogUpdate("environment", key, user.ID, user.Email, beforeData, afterData)

	return c.JSON(updatedEnvironment)
}

// DeleteEnvironment handles DELETE /api/v1/environments/:key
func (h *EnvironmentsHandler) DeleteEnvironment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Environment key is required")
	}

	// Get existing environment for audit
	var existingEnvironment models.Environment
	err := h.environmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingEnvironment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Environment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve environment")
	}

	result, err := h.environmentsCollection.DeleteOne(context.Background(), bson.M{"key": key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to delete environment")
	}
	if result.DeletedCount == 0 {
		return middleware.NewNotFoundError("Environment not found")
	}

	// Log audit entry
	deletedData := map[string]interface{}{
		"key":         existingEnvironment.Key,
		"name":        existingEnvironment.Name,
		"description": existingEnvironment.Description,
		"color":       existingEnvironment.Color,
	}
	h.auditService.LogDelete("environment", key, user.ID, user.Email, deletedData)

	return c.SendStatus(fiber.StatusNoContent)
}
