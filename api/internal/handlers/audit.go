package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

type AuditHandler struct {
	auditCollection *mongo.Collection
}

type AuditLogResponse struct {
	Data  []models.AuditLogEntry `json:"data"`
	Total int64                  `json:"total"`
}

func NewAuditHandler(auditCollection *mongo.Collection) *AuditHandler {
	return &AuditHandler{
		auditCollection: auditCollection,
	}
}

// GET /api/v1/audit-log
func (h *AuditHandler) GetAuditLog(c *fiber.Ctx) error {
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

	// Parse query parameters for filtering
	resource := c.Query("resource")
	action := c.Query("action")
	userID := c.Query("user_id")
	resourceID := c.Query("resource_id")

	// Parse custom limit and offset (takes precedence over page-based pagination)
	if limitStr := c.Query("limit"); limitStr != "" {
		if customLimit, err := strconv.Atoi(limitStr); err == nil && customLimit > 0 && customLimit <= 1000 {
			limit = customLimit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if customOffset, err := strconv.Atoi(offsetStr); err == nil && customOffset >= 0 {
			skip = customOffset
		}
	}

	// Build filter
	filter := bson.M{}
	if resource != "" {
		filter["resource"] = resource
	}
	if action != "" {
		filter["action"] = action
	}
	if userID != "" {
		filter["user_id"] = userID
	}
	if resourceID != "" {
		filter["resource_id"] = resourceID
	}

	// Parse date range filters
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			if filter["timestamp"] == nil {
				filter["timestamp"] = bson.M{}
			}
			filter["timestamp"].(bson.M)["$gte"] = startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			if filter["timestamp"] == nil {
				filter["timestamp"] = bson.M{}
			}
			filter["timestamp"].(bson.M)["$lte"] = endDate
		}
	}

	// Get total count
	totalCount, err := h.auditCollection.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to count audit log entries",
			},
		})
	}

	// Get audit log entries with pagination and sorting
	opts := options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := h.auditCollection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch audit log entries",
			},
		})
	}
	defer cursor.Close(ctx)

	var auditEntries []models.AuditLogEntry
	if err = cursor.All(ctx, &auditEntries); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to decode audit log entries",
			},
		})
	}

	if auditEntries == nil {
		auditEntries = []models.AuditLogEntry{}
	}

	response := AuditLogResponse{
		Data:  auditEntries,
		Total: totalCount,
	}

	return c.JSON(response)
}
