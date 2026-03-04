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
	totalCount, err := h.apiKeysCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to count API keys",
			},
		})
	}

	// Get API keys with pagination
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := h.apiKeysCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch API keys",
			},
		})
	}
	defer cursor.Close(ctx)

	var apiKeys []models.ApiKey
	if err = cursor.All(ctx, &apiKeys); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to decode API keys",
			},
		})
	}

	if apiKeys == nil {
		apiKeys = []models.ApiKey{}
	}

	response := ApiKeysListResponse{
		Data:  apiKeys,
		Total: totalCount,
	}

	return c.JSON(response)
}

// POST /api/v1/api-keys
func (h *ApiKeysHandler) CreateApiKey(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req CreateApiKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_JSON",
				"message": "Invalid JSON in request body",
			},
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "Name is required",
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

	// Generate raw API key
	rawKey, err := generateAPIKey()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "KEY_GENERATION_ERROR",
				"message": "Failed to generate API key",
			},
		})
	}

	// Hash the key for storage
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "KEY_HASHING_ERROR",
				"message": "Failed to hash API key",
			},
		})
	}

	apiKey := models.ApiKey{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		KeyHash:     string(hashedKey),
		Environment: req.Environment,
		LastUsedAt:  nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = h.apiKeysCollection.InsertOne(ctx, apiKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create API key",
			},
		})
	}

	// Log audit entry
	user := c.Locals("user").(*models.User)
	h.auditService.LogCreate("api_key", apiKey.ID.Hex(), user.ID, user.Email, map[string]interface{}{
		"name":        apiKey.Name,
		"environment": apiKey.Environment,
	})

	// Return response with the raw key (only time it's exposed)
	response := CreateApiKeyResponse{
		ID:          apiKey.ID,
		Name:        apiKey.Name,
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "API key ID is required",
			},
		})
	}

	// Convert string to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_ID",
				"message": "Invalid API key ID format",
			},
		})
	}

	// Check if API key exists
	var apiKey models.ApiKey
	err = h.apiKeysCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&apiKey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "API_KEY_NOT_FOUND",
					"message": "API key not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to fetch API key",
			},
		})
	}

	_, err = h.apiKeysCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete API key",
			},
		})
	}

	// Log audit entry
	user := c.Locals("user").(*models.User)
	h.auditService.LogDelete("api_key", idStr, user.ID, user.Email, map[string]interface{}{
		"name":        apiKey.Name,
		"environment": apiKey.Environment,
	})

	return c.Status(fiber.StatusNoContent).Send(nil)
}
