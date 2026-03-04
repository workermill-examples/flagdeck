package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

type APIKeyContext struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Environment string             `json:"environment"`
}

type APIKeyConfig struct {
	APIKeyCollection *mongo.Collection
}

func AuthenticateAPIKey(config APIKeyConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return NewUnauthorizedError("X-API-Key header required")
		}

		if len(apiKey) < 10 {
			return NewUnauthorizedError("Invalid API key format")
		}

		var apiKeys []models.ApiKey
		cursor, err := config.APIKeyCollection.Find(context.Background(), bson.M{})
		if err != nil {
			log.Printf("Database error while fetching API keys: %v", err)
			return NewDatabaseError("Failed to verify API key")
		}
		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), &apiKeys)
		if err != nil {
			log.Printf("Database error while decoding API keys: %v", err)
			return NewDatabaseError("Failed to verify API key")
		}

		var matchedKey *models.ApiKey
		for i := range apiKeys {
			err := bcrypt.CompareHashAndPassword([]byte(apiKeys[i].KeyHash), []byte(apiKey))
			if err == nil {
				matchedKey = &apiKeys[i]
				break
			}
		}

		if matchedKey == nil {
			return NewUnauthorizedError("Invalid API key")
		}

		now := time.Now()
		_, err = config.APIKeyCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": matchedKey.ID},
			bson.M{
				"$set": bson.M{
					"last_used_at": now,
					"updated_at":   now,
				},
			},
		)
		if err != nil {
			log.Printf("Failed to update API key last_used_at: %v", err)
		}

		apiKeyCtx := APIKeyContext{
			ID:          matchedKey.ID,
			Name:        matchedKey.Name,
			Environment: matchedKey.Environment,
		}

		c.Locals("apikey", apiKeyCtx)
		return c.Next()
	}
}

func RequireEnvironment(requiredEnvironment string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey, ok := c.Locals("apikey").(APIKeyContext)
		if !ok {
			return NewUnauthorizedError("API key context not found")
		}

		if apiKey.Environment != requiredEnvironment && requiredEnvironment != "" {
			return NewForbiddenError("API key not authorized for this environment")
		}

		return c.Next()
	}
}

func RequireEnvironments(requiredEnvironments []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey, ok := c.Locals("apikey").(APIKeyContext)
		if !ok {
			return NewUnauthorizedError("API key context not found")
		}

		for _, env := range requiredEnvironments {
			if apiKey.Environment == env {
				return c.Next()
			}
		}

		return NewForbiddenError("API key not authorized for required environments")
	}
}
