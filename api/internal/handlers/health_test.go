package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/workermill-examples/flagdeck/api/internal/database"
)

// Pingable interface for testable database connections
type Pingable interface {
	Ping() error
}

// TestableHealthHandler is a version of HealthHandler that accepts interfaces
type TestableHealthHandler struct {
	MongoDB Pingable
	Redis   Pingable
}

// GetHealth method for testable health handler
func (h *TestableHealthHandler) GetHealth(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:  "ok",
		MongoDB: "connected",
		Redis:   "connected",
	}

	statusCode := fiber.StatusOK

	// Check MongoDB connection
	if err := h.MongoDB.Ping(); err != nil {
		response.MongoDB = "disconnected"
		response.Status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	// Check Redis connection
	if err := h.Redis.Ping(); err != nil {
		response.Redis = "disconnected"
		response.Status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	// If both are down, status should be "down"
	if response.MongoDB == "disconnected" && response.Redis == "disconnected" {
		response.Status = "down"
	}

	return c.Status(statusCode).JSON(response)
}

// MockDB implements Pingable interface for testing
type MockDB struct {
	shouldFail bool
}

func (m *MockDB) Ping() error {
	if m.shouldFail {
		return fmt.Errorf("connection failed")
	}
	return nil
}

func TestHealthHandler_GetHealth_AllServicesConnected(t *testing.T) {
	// Test when both MongoDB and Redis are connected
	mockMongo := &MockDB{shouldFail: false}
	mockRedis := &MockDB{shouldFail: false}

	handler := &TestableHealthHandler{
		MongoDB: mockMongo,
		Redis:   mockRedis,
	}

	app := fiber.New()
	app.Get("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response HealthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}
	if response.MongoDB != "connected" {
		t.Errorf("Expected MongoDB 'connected', got '%s'", response.MongoDB)
	}
	if response.Redis != "connected" {
		t.Errorf("Expected Redis 'connected', got '%s'", response.Redis)
	}
}

func TestHealthHandler_GetHealth_MongoDBDown(t *testing.T) {
	// Test when MongoDB is down
	mockMongo := &MockDB{shouldFail: true}
	mockRedis := &MockDB{shouldFail: false}

	handler := &TestableHealthHandler{
		MongoDB: mockMongo,
		Redis:   mockRedis,
	}

	app := fiber.New()
	app.Get("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response HealthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", response.Status)
	}
	if response.MongoDB != "disconnected" {
		t.Errorf("Expected MongoDB 'disconnected', got '%s'", response.MongoDB)
	}
	if response.Redis != "connected" {
		t.Errorf("Expected Redis 'connected', got '%s'", response.Redis)
	}
}

func TestHealthHandler_GetHealth_RedisDown(t *testing.T) {
	// Test when Redis is down
	mockMongo := &MockDB{shouldFail: false}
	mockRedis := &MockDB{shouldFail: true}

	handler := &TestableHealthHandler{
		MongoDB: mockMongo,
		Redis:   mockRedis,
	}

	app := fiber.New()
	app.Get("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response HealthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", response.Status)
	}
	if response.MongoDB != "connected" {
		t.Errorf("Expected MongoDB 'connected', got '%s'", response.MongoDB)
	}
	if response.Redis != "disconnected" {
		t.Errorf("Expected Redis 'disconnected', got '%s'", response.Redis)
	}
}

func TestHealthHandler_GetHealth_BothServicesDown(t *testing.T) {
	// Test when both MongoDB and Redis are down
	mockMongo := &MockDB{shouldFail: true}
	mockRedis := &MockDB{shouldFail: true}

	handler := &TestableHealthHandler{
		MongoDB: mockMongo,
		Redis:   mockRedis,
	}

	app := fiber.New()
	app.Get("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response HealthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "down" {
		t.Errorf("Expected status 'down', got '%s'", response.Status)
	}
	if response.MongoDB != "disconnected" {
		t.Errorf("Expected MongoDB 'disconnected', got '%s'", response.MongoDB)
	}
	if response.Redis != "disconnected" {
		t.Errorf("Expected Redis 'disconnected', got '%s'", response.Redis)
	}
}

func TestHealthHandler_Structure(t *testing.T) {
	// Test that the health handler has the correct structure
	var mongoDB *database.MongoDB
	var redisDB *database.RedisDB

	handler := NewHealthHandler(mongoDB, redisDB)

	if handler == nil {
		t.Fatal("NewHealthHandler returned nil")
	}

	if handler.MongoDB != mongoDB {
		t.Error("MongoDB field not set correctly")
	}

	if handler.Redis != redisDB {
		t.Error("Redis field not set correctly")
	}
}

func TestHealthResponse_Structure(t *testing.T) {
	// Test the response structure
	response := HealthResponse{
		Status:  "ok",
		MongoDB: "connected",
		Redis:   "connected",
	}

	if response.Status != "ok" {
		t.Error("Status field not set correctly")
	}

	if response.MongoDB != "connected" {
		t.Error("MongoDB field not set correctly")
	}

	if response.Redis != "connected" {
		t.Error("Redis field not set correctly")
	}
}
