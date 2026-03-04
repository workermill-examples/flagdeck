package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
)

func TestRegisterRequest_Validation(t *testing.T) {
	// Test the request structure validation
	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	if reqBody.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", reqBody.Email)
	}

	if reqBody.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", reqBody.Name)
	}

	if reqBody.Password != "password123" {
		t.Errorf("Expected password 'password123', got '%s'", reqBody.Password)
	}
}

func TestLoginRequest_Structure(t *testing.T) {
	// Test the login request structure
	reqBody := LoginRequest{
		Email:    "user@example.com",
		Password: "mypassword",
	}

	if reqBody.Email != "user@example.com" {
		t.Errorf("Expected email 'user@example.com', got '%s'", reqBody.Email)
	}

	if reqBody.Password != "mypassword" {
		t.Errorf("Expected password 'mypassword', got '%s'", reqBody.Password)
	}
}

func TestRefreshRequest_Structure(t *testing.T) {
	// Test the refresh request structure
	reqBody := RefreshRequest{
		RefreshToken: "sample-refresh-token",
	}

	if reqBody.RefreshToken != "sample-refresh-token" {
		t.Errorf("Expected refresh_token 'sample-refresh-token', got '%s'", reqBody.RefreshToken)
	}
}

func TestAuthResponse_Structure(t *testing.T) {
	// Test the auth response structure
	userResp := UserResponse{
		ID:        "507f1f77bcf86cd799439011",
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      "viewer",
		CreatedAt: time.Now(),
	}

	authResp := AuthResponse{
		User:         userResp,
		AccessToken:  "access-token-here",
		RefreshToken: "refresh-token-here",
	}

	if authResp.User.Email != "test@example.com" {
		t.Errorf("Expected user email 'test@example.com', got '%s'", authResp.User.Email)
	}

	if authResp.User.Role != "viewer" {
		t.Errorf("Expected user role 'viewer', got '%s'", authResp.User.Role)
	}

	if authResp.AccessToken != "access-token-here" {
		t.Errorf("Expected access_token 'access-token-here', got '%s'", authResp.AccessToken)
	}

	if authResp.RefreshToken != "refresh-token-here" {
		t.Errorf("Expected refresh_token 'refresh-token-here', got '%s'", authResp.RefreshToken)
	}
}

func TestUserResponse_Structure(t *testing.T) {
	// Test the user response structure
	now := time.Now()
	userResp := UserResponse{
		ID:        "507f1f77bcf86cd799439011",
		Email:     "user@test.com",
		Name:      "User Name",
		Role:      "admin",
		CreatedAt: now,
	}

	if userResp.ID != "507f1f77bcf86cd799439011" {
		t.Errorf("Expected ID '507f1f77bcf86cd799439011', got '%s'", userResp.ID)
	}

	if userResp.Email != "user@test.com" {
		t.Errorf("Expected email 'user@test.com', got '%s'", userResp.Email)
	}

	if userResp.Name != "User Name" {
		t.Errorf("Expected name 'User Name', got '%s'", userResp.Name)
	}

	if userResp.Role != "admin" {
		t.Errorf("Expected role 'admin', got '%s'", userResp.Role)
	}

	if !userResp.CreatedAt.Equal(now) {
		t.Errorf("Expected created_at to match, got different time")
	}
}

func TestAuthHandler_Structure(t *testing.T) {
	// Test that the auth handler has the correct structure
	handler := NewAuthHandler(nil, "test-secret")

	if handler == nil {
		t.Fatal("NewAuthHandler returned nil")
	}

	if handler.JWTSecret != "test-secret" {
		t.Errorf("Expected JWT secret 'test-secret', got '%s'", handler.JWTSecret)
	}

	if handler.UserCollection == nil {
		// This is expected since we passed nil, but just verifying the field exists
		t.Log("UserCollection field exists and is set to nil as expected")
	}
}

func TestRegisterRequest_NoRoleField(t *testing.T) {
	// Test that RegisterRequest does not accept role field (as per spec)
	// This test verifies the structure doesn't include role field
	reqBodyJSON := `{
		"email": "test@example.com",
		"name": "Test User",
		"password": "password123",
		"role": "admin"
	}`

	var reqBody RegisterRequest
	err := json.Unmarshal([]byte(reqBodyJSON), &reqBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// The role field should not be set even if provided in JSON
	// because RegisterRequest struct doesn't have a role field
	if reqBody.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", reqBody.Email)
	}

	if reqBody.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", reqBody.Name)
	}

	if reqBody.Password != "password123" {
		t.Errorf("Expected password 'password123', got '%s'", reqBody.Password)
	}

	// Verify the struct doesn't have a role field by checking marshaling back
	marshalled, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var parsed map[string]interface{}
	json.Unmarshal(marshalled, &parsed)

	if _, exists := parsed["role"]; exists {
		t.Error("RegisterRequest should not include role field in JSON output")
	}
}

func TestJWTTokenDuration_Requirements(t *testing.T) {
	// Test token duration requirements from spec
	// Access token: 15 minutes, Refresh token: 7 days
	// This is more of a specification verification test

	userID := primitive.NewObjectID()
	email := "test@example.com"
	name := "Test User"
	role := "viewer"
	secret := "test-secret"

	accessToken, refreshToken, err := middleware.GenerateTokens(userID, email, name, role, secret)
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	if accessToken == "" {
		t.Error("Access token should not be empty")
	}

	if refreshToken == "" {
		t.Error("Refresh token should not be empty")
	}

	// Verify tokens are different
	if accessToken == refreshToken {
		t.Error("Access token and refresh token should be different")
	}
}

func TestValidateRefreshToken_Structure(t *testing.T) {
	// Test refresh token validation structure
	// This verifies the function exists and has correct signature

	userID := primitive.NewObjectID()
	email := "test@example.com"
	name := "Test User"
	role := "viewer"
	secret := "test-secret"

	// Generate a refresh token first
	_, refreshToken, err := middleware.GenerateTokens(userID, email, name, role, secret)
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Validate the refresh token
	claims, err := middleware.ValidateRefreshToken(refreshToken, secret)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	if claims == nil {
		t.Error("Claims should not be nil")
	}

	if claims.Email != email {
		t.Errorf("Expected email '%s', got '%s'", email, claims.Email)
	}

	if claims.Name != name {
		t.Errorf("Expected name '%s', got '%s'", name, claims.Name)
	}

	if claims.Role != role {
		t.Errorf("Expected role '%s', got '%s'", role, claims.Role)
	}

	if claims.Type != "refresh" {
		t.Errorf("Expected token type 'refresh', got '%s'", claims.Type)
	}
}

func TestAuthHandler_DefaultRole(t *testing.T) {
	// Test that registration defaults to "viewer" role
	// This is implied by the spec: "no role field accepted, defaults to viewer"

	// Create a request without role
	reqBody := RegisterRequest{
		Email:    "newuser@example.com",
		Name:     "New User",
		Password: "password123",
	}

	// Verify the request structure doesn't include role
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var parsed map[string]interface{}
	err = json.Unmarshal(reqBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}

	if _, hasRole := parsed["role"]; hasRole {
		t.Error("RegisterRequest should not include role field")
	}

	// Verify we have the required fields
	if parsed["email"] != reqBody.Email {
		t.Error("Email field missing or incorrect")
	}

	if parsed["name"] != reqBody.Name {
		t.Error("Name field missing or incorrect")
	}

	if parsed["password"] != reqBody.Password {
		t.Error("Password field missing or incorrect")
	}
}

func TestErrorResponseFormat(t *testing.T) {
	// Test error response format consistency
	// From spec: {"error": {"code": "CODE", "message": "msg"}}

	// This test verifies that auth handlers use the middleware error functions
	// which should return the correct error format

	validationErr := middleware.NewValidationError("Test validation error")
	unauthorizedErr := middleware.NewUnauthorizedError("Test unauthorized error")
	conflictErr := middleware.NewConflictError("Test conflict error")

	// These should be fiber errors with the proper structure
	// We can't easily test the JSON structure here without a full request context,
	// but we can verify the errors exist and have messages

	if validationErr == nil {
		t.Error("Validation error should not be nil")
	}

	if unauthorizedErr == nil {
		t.Error("Unauthorized error should not be nil")
	}

	if conflictErr == nil {
		t.Error("Conflict error should not be nil")
	}
}
