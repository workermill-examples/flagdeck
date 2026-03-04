package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/models"
)

func createTestApiKeysApp(handler *ApiKeysHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Add test middleware that sets user context
	app.Use(func(c *fiber.Ctx) error {
		userID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
		c.Locals("user", middleware.UserContext{
			ID:    userID,
			Email: "test@example.com",
		})
		return c.Next()
	})

	app.Get("/api/v1/api-keys", handler.ListApiKeys)
	app.Post("/api/v1/api-keys", handler.CreateApiKey)
	app.Delete("/api/v1/api-keys/:id", handler.DeleteApiKey)

	return app
}

func TestCreateApiKeyRequest_Validation(t *testing.T) {
	// Test the request structure validation
	reqBody := CreateApiKeyRequest{
		Name:        "Test API Key",
		Environment: "production",
	}

	if reqBody.Name != "Test API Key" {
		t.Errorf("Expected name 'Test API Key', got '%s'", reqBody.Name)
	}

	if reqBody.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", reqBody.Environment)
	}
}

func TestCreateApiKeyResponse_Structure(t *testing.T) {
	// Test the response structure
	response := CreateApiKeyResponse{
		Name:        "Test API Key",
		KeyPrefix:   "fd_test1234",
		RawKey:      "fd_test1234567890abcdef",
		Environment: "production",
	}

	if response.Name != "Test API Key" {
		t.Error("Name field not set correctly")
	}

	if response.KeyPrefix != "fd_test1234" {
		t.Error("KeyPrefix field not set correctly")
	}

	if response.RawKey != "fd_test1234567890abcdef" {
		t.Error("RawKey field not set correctly")
	}

	if response.Environment != "production" {
		t.Error("Environment field not set correctly")
	}
}

func TestApiKeysHandler_CreateApiKey_InvalidJSON(t *testing.T) {
	// Test with a handler that has nil collections - this should fail gracefully
	handler := &ApiKeysHandler{
		apiKeysCollection: nil,
		auditService:      nil,
	}

	app := createTestApiKeysApp(handler)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/api-keys", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Check the error message structure
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Should contain error response
	if !strings.Contains(string(body), "error") {
		t.Error("Response should contain error field")
	}
}

func TestApiKeysHandler_CreateApiKey_MissingName(t *testing.T) {
	handler := &ApiKeysHandler{
		apiKeysCollection: nil,
		auditService:      nil,
	}

	app := createTestApiKeysApp(handler)

	reqBody := CreateApiKeyRequest{
		// Missing Name
		Environment: "production",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/api-keys", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestApiKeysHandler_CreateApiKey_MissingEnvironment(t *testing.T) {
	handler := &ApiKeysHandler{
		apiKeysCollection: nil,
		auditService:      nil,
	}

	app := createTestApiKeysApp(handler)

	reqBody := CreateApiKeyRequest{
		Name: "Test API Key",
		// Missing Environment
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/api-keys", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestApiKeysHandler_DeleteApiKey_InvalidID(t *testing.T) {
	handler := &ApiKeysHandler{
		apiKeysCollection: nil,
		auditService:      nil,
	}

	app := createTestApiKeysApp(handler)

	req := httptest.NewRequest("DELETE", "/api/v1/api-keys/invalid-id", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestApiKeysHandler_DeleteApiKey_MissingID(t *testing.T) {
	handler := &ApiKeysHandler{
		apiKeysCollection: nil,
		auditService:      nil,
	}

	app := createTestApiKeysApp(handler)

	req := httptest.NewRequest("DELETE", "/api/v1/api-keys/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Should get 405 Method Not Allowed since the route requires an ID parameter
	if resp.StatusCode != fiber.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", resp.StatusCode)
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key1, err := generateAPIKey()
	if err != nil {
		t.Fatalf("Failed to generate API key: %v", err)
	}

	key2, err := generateAPIKey()
	if err != nil {
		t.Fatalf("Failed to generate API key: %v", err)
	}

	// Keys should be different
	if key1 == key2 {
		t.Error("Generated API keys should be unique")
	}

	// Keys should start with "fd_"
	if len(key1) < 3 || key1[:3] != "fd_" {
		t.Errorf("Expected API key to start with 'fd_', got '%s'", key1)
	}

	if len(key2) < 3 || key2[:3] != "fd_" {
		t.Errorf("Expected API key to start with 'fd_', got '%s'", key2)
	}

	// Keys should be long enough (fd_ + 64 hex chars = 67 total)
	if len(key1) != 67 {
		t.Errorf("Expected API key length 67, got %d", len(key1))
	}

	if len(key2) != 67 {
		t.Errorf("Expected API key length 67, got %d", len(key2))
	}
}

func TestApiKeysHandler_Structure(t *testing.T) {
	// Test that the handler has the correct structure
	handler := NewApiKeysHandler(nil, nil)

	if handler == nil {
		t.Fatal("NewApiKeysHandler returned nil")
	}
}

func TestApiKeysListResponse_Structure(t *testing.T) {
	// Test the list response structure
	response := ApiKeysListResponse{
		Data:  []models.ApiKey{},
		Total: 0,
	}

	if response.Data == nil {
		t.Error("Data field should be initialized")
	}

	if response.Total != 0 {
		t.Error("Total field not set correctly")
	}
}
