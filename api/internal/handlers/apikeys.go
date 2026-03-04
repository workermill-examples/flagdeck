package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/models"
	"github.com/workermill-examples/flagdeck/api/internal/services"
)

type ApiKeysHandler struct {
	apiKeysCollection *mongo.Collection
	auditService      *services.AuditService
}

type ApiKeysListResponse struct {
	Data  []models.ApiKey `json:"data"`
	Total int64           `json:"total"`
}

type CreateApiKeyRequest struct {
	Name        string `json:"name" validate:"required"`
	Environment string `json:"environment" validate:"required"`
}

type CreateApiKeyResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	KeyPrefix   string             `json:"key_prefix"`
	RawKey      string             `json:"raw_key"`
	Environment string             `json:"environment"`
	LastUsedAt  *time.Time         `json:"last_used_at"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func NewApiKeysHandler(apiKeysCollection *mongo.Collection, auditService *services.AuditService) *ApiKeysHandler {
	return &ApiKeysHandler{
		apiKeysCollection: apiKeysCollection,
		auditService:      auditService,
	}
}

// generateAPIKey generates a cryptographically secure random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "fd_" + hex.EncodeToString(bytes), nil
}

// GET /api/v1/api-keys
func (h *ApiKeysHandler) ListApiKeys(c *fiber.Ctx) error {
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
	total, err := h.apiKeysCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count API keys")
	}

	// Get API keys with pagination
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := h.apiKeysCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return middleware.NewDatabaseError("Failed to retrieve API keys")
	}
	defer cursor.Close(context.Background())

	var apiKeys []models.ApiKey
	if err := cursor.All(context.Background(), &apiKeys); err != nil {
		return middleware.NewDatabaseError("Failed to decode API keys")
	}

	if apiKeys == nil {
		apiKeys = []models.ApiKey{}
	}

	return c.JSON(ApiKeysListResponse{
		Data:  apiKeys,
		Total: total,
	})
}

// POST /api/v1/api-keys
func (h *ApiKeysHandler) CreateApiKey(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)

	var req CreateApiKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewBadRequestError("Invalid request body")
	}

	// Validate required fields
	if req.Name == "" {
		return middleware.NewValidationError("Name is required")
	}
	if req.Environment == "" {
		return middleware.NewValidationError("Environment is required")
	}

	// Generate raw API key
	rawKey, err := generateAPIKey()
	if err != nil {
		return middleware.NewInternalError("Failed to generate API key")
	}

	// Hash the key for storage
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return middleware.NewInternalError("Failed to hash API key")
	}

	// Extract key prefix (first 8 characters after "fd_")
	keyPrefix := rawKey[:11] // "fd_" + first 8 hex chars

	now := time.Now()
	apiKey := models.ApiKey{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		KeyPrefix:   keyPrefix,
		KeyHash:     string(hashedKey),
		Environment: req.Environment,
		LastUsedAt:  nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = h.apiKeysCollection.InsertOne(context.Background(), apiKey)
	if err != nil {
		return middleware.NewDatabaseError("Failed to create API key")
	}

	// Log audit entry
	changes := map[string]interface{}{
		"name":        apiKey.Name,
		"key_prefix":  apiKey.KeyPrefix,
		"environment": apiKey.Environment,
	}
	h.auditService.LogCreate("api_key", apiKey.ID.Hex(), user.ID, user.Email, changes)

	// Return response with the raw key (only time it's exposed)
	response := CreateApiKeyResponse{
		ID:          apiKey.ID,
		Name:        apiKey.Name,
		KeyPrefix:   apiKey.KeyPrefix,
		RawKey:      rawKey,
		Environment: apiKey.Environment,
		LastUsedAt:  apiKey.LastUsedAt,
		CreatedAt:   apiKey.CreatedAt,
		UpdatedAt:   apiKey.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// DELETE /api/v1/api-keys/:id
func (h *ApiKeysHandler) DeleteApiKey(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.UserContext)
	idStr := c.Params("id")

	if idStr == "" {
		return middleware.NewBadRequestError("API key ID is required")
	}

	// Convert string to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return middleware.NewBadRequestError("Invalid API key ID format")
	}

	// Get existing API key for audit
	var existingApiKey models.ApiKey
	err = h.apiKeysCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingApiKey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("API key not found")
		}
		return middleware.NewDatabaseError("Failed to retrieve API key")
	}

	result, err := h.apiKeysCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return middleware.NewDatabaseError("Failed to delete API key")
	}
	if result.DeletedCount == 0 {
		return middleware.NewNotFoundError("API key not found")
	}

	// Log audit entry
	deletedData := map[string]interface{}{
		"name":        existingApiKey.Name,
		"key_prefix":  existingApiKey.KeyPrefix,
		"environment": existingApiKey.Environment,
	}
	h.auditService.LogDelete("api_key", idStr, user.ID, user.Email, deletedData)

	return c.SendStatus(fiber.StatusNoContent)
}
