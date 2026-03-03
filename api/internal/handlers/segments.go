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

type SegmentsHandler struct {
	segmentsCollection *mongo.Collection
	auditService       *services.AuditService
}

type SegmentsListResponse struct {
	Data  []models.Segment `json:"data"`
	Total int64            `json:"total"`
}

type CreateSegmentRequest struct {
	Key         string               `json:"key" validate:"required"`
	Name        string               `json:"name" validate:"required"`
	Description string               `json:"description"`
	Rules       []models.SegmentRule `json:"rules"`
}

type UpdateSegmentRequest struct {
	Name        *string              `json:"name,omitempty"`
	Description *string              `json:"description,omitempty"`
	Rules       []models.SegmentRule `json:"rules,omitempty"`
}

func NewSegmentsHandler(segmentsCollection *mongo.Collection, auditService *services.AuditService) *SegmentsHandler {
	return &SegmentsHandler{
		segmentsCollection: segmentsCollection,
		auditService:       auditService,
	}
}

// ListSegments handles GET /api/v1/segments
func (h *SegmentsHandler) ListSegments(c *fiber.Ctx) error {
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
	total, err := h.segmentsCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count segments")
	}

	// Get segments with pagination
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := h.segmentsCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve segments")
	}
	defer cursor.Close(context.Background())

	var segments []models.Segment
	if err := cursor.All(context.Background(), &segments); err != nil {
		return middleware.NewDatabaseError("Failed to decode segments")
	}

	if segments == nil {
		segments = []models.Segment{}
	}

	return c.JSON(SegmentsListResponse{
		Data:  segments,
		Total: total,
	})
}

// GetSegment handles GET /api/v1/segments/:key
func (h *SegmentsHandler) GetSegment(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return middleware.NewBadRequestError("Segment key is required")
	}

	var segment models.Segment
	err := h.segmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&segment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Segment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve segment")
	}

	return c.JSON(segment)
}

// CreateSegment handles POST /api/v1/segments
func (h *SegmentsHandler) CreateSegment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)

	var req CreateSegmentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.Key == "" {
		return middleware.NewValidationError("Segment key is required")
	}
	if req.Name == "" {
		return middleware.NewValidationError("Segment name is required")
	}

	// Validate segment rules
	if err := h.validateSegmentRules(req.Rules); err != nil {
		return err
	}

	// Check if segment with this key already exists
	count, err := h.segmentsCollection.CountDocuments(context.Background(), bson.M{"key": req.Key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to check existing segment")
	}
	if count > 0 {
		return middleware.NewConflictError("Segment with this key already exists")
	}

	// Initialize rules if not provided
	if req.Rules == nil {
		req.Rules = []models.SegmentRule{}
	}

	now := time.Now()
	segment := models.Segment{
		ID:          primitive.NewObjectID(),
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = h.segmentsCollection.InsertOne(context.Background(), segment)
	if err != nil {
		return middleware.NewDatabaseError("Failed to create segment")
	}

	// Log audit entry
	changes := map[string]interface{}{
		"key":         segment.Key,
		"name":        segment.Name,
		"description": segment.Description,
		"rules":       segment.Rules,
	}
	h.auditService.LogCreate("segment", segment.Key, user.ID, user.Email, changes)

	return c.Status(fiber.StatusCreated).JSON(segment)
}

// UpdateSegment handles PUT /api/v1/segments/:key
func (h *SegmentsHandler) UpdateSegment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Segment key is required")
	}

	var req UpdateSegmentRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate segment rules if provided
	if req.Rules != nil {
		if err := h.validateSegmentRules(req.Rules); err != nil {
			return err
		}
	}

	// Get existing segment
	var existingSegment models.Segment
	err := h.segmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingSegment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Segment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve segment")
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Name != nil {
		update["$set"].(bson.M)["name"] = *req.Name
	}
	if req.Description != nil {
		update["$set"].(bson.M)["description"] = *req.Description
	}
	if req.Rules != nil {
		update["$set"].(bson.M)["rules"] = req.Rules
	}

	result, err := h.segmentsCollection.UpdateOne(
		context.Background(),
		bson.M{"key": key},
		update,
	)
	if err != nil {
		return middleware.NewDatabaseError("Failed to update segment")
	}
	if result.MatchedCount == 0 {
		return middleware.NewNotFoundError("Segment not found")
	}

	// Get updated segment
	var updatedSegment models.Segment
	err = h.segmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&updatedSegment)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve updated segment")
	}

	// Log audit entry
	beforeData := map[string]interface{}{
		"name":        existingSegment.Name,
		"description": existingSegment.Description,
		"rules":       existingSegment.Rules,
	}
	afterData := map[string]interface{}{
		"name":        updatedSegment.Name,
		"description": updatedSegment.Description,
		"rules":       updatedSegment.Rules,
	}
	h.auditService.LogUpdate("segment", key, user.ID, user.Email, beforeData, afterData)

	return c.JSON(updatedSegment)
}

// DeleteSegment handles DELETE /api/v1/segments/:key
func (h *SegmentsHandler) DeleteSegment(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	key := c.Params("key")

	if key == "" {
		return middleware.NewBadRequestError("Segment key is required")
	}

	// Get existing segment for audit
	var existingSegment models.Segment
	err := h.segmentsCollection.FindOne(context.Background(), bson.M{"key": key}).Decode(&existingSegment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("Segment not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve segment")
	}

	result, err := h.segmentsCollection.DeleteOne(context.Background(), bson.M{"key": key})
	if err != nil {
		return middleware.NewDatabaseError("Failed to delete segment")
	}
	if result.DeletedCount == 0 {
		return middleware.NewNotFoundError("Segment not found")
	}

	// Log audit entry
	deletedData := map[string]interface{}{
		"key":         existingSegment.Key,
		"name":        existingSegment.Name,
		"description": existingSegment.Description,
		"rules":       existingSegment.Rules,
	}
	h.auditService.LogDelete("segment", key, user.ID, user.Email, deletedData)

	return c.SendStatus(fiber.StatusNoContent)
}

// validateSegmentRules validates the structure and operators of segment rules
func (h *SegmentsHandler) validateSegmentRules(rules []models.SegmentRule) error {
	validOperators := map[string]bool{
		"equals":             true,
		"not_equals":         true,
		"in":                 true,
		"not_in":             true,
		"contains":           true,
		"starts_with":        true,
		"ends_with":          true,
		"greater_than":       true,
		"greater_than_equal": true,
		"less_than":          true,
		"less_than_equal":    true,
	}

	for _, rule := range rules {
		if rule.Attribute == "" {
			return middleware.NewValidationError("Rule attribute is required")
		}
		if rule.Operator == "" {
			return middleware.NewValidationError("Rule operator is required")
		}
		if !validOperators[rule.Operator] {
			return middleware.NewValidationError("Invalid rule operator: " + rule.Operator)
		}
		if rule.Value == nil {
			return middleware.NewValidationError("Rule value is required")
		}
	}

	return nil
}
