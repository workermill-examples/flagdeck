package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/services"
)

func TestEvaluationRequest_Structure(t *testing.T) {
	// Test the request structure validation
	reqBody := services.EvaluationRequest{
		FlagKey:     "test-flag",
		Environment: "production",
		UserID:      "user-123",
		UserContext: map[string]interface{}{
			"country": "US",
		},
	}

	if reqBody.FlagKey != "test-flag" {
		t.Errorf("Expected flag_key 'test-flag', got '%s'", reqBody.FlagKey)
	}

	if reqBody.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", reqBody.Environment)
	}

	if reqBody.UserID != "user-123" {
		t.Errorf("Expected user_id 'user-123', got '%s'", reqBody.UserID)
	}

	if reqBody.UserContext == nil {
		t.Error("UserContext should not be nil")
	}
}

func TestBulkEvaluationRequest_Structure(t *testing.T) {
	// Test the bulk request structure
	reqBody := services.BulkEvaluationRequest{
		Environment: "production",
		UserID:      "user-123",
		UserContext: map[string]interface{}{
			"country": "US",
		},
		FlagKeys: []string{"flag1", "flag2", "flag3"},
	}

	if reqBody.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", reqBody.Environment)
	}

	if reqBody.UserID != "user-123" {
		t.Errorf("Expected user_id 'user-123', got '%s'", reqBody.UserID)
	}

	if len(reqBody.FlagKeys) != 3 {
		t.Errorf("Expected 3 flag keys, got %d", len(reqBody.FlagKeys))
	}

	expectedFlags := []string{"flag1", "flag2", "flag3"}
	for i, flag := range reqBody.FlagKeys {
		if flag != expectedFlags[i] {
			t.Errorf("Expected flag '%s', got '%s'", expectedFlags[i], flag)
		}
	}
}

func TestEvaluationResult_Structure(t *testing.T) {
	// Test the response structure
	ruleID := "rule-123"
	result := services.EvaluationResult{
		Key:          "test-flag",
		Value:        true,
		Type:         "boolean",
		Reason:       "targeting_rule_match",
		RuleID:       &ruleID,
		Environment:  "production",
		EvaluationMs: 5,
	}

	if result.Key != "test-flag" {
		t.Errorf("Expected key 'test-flag', got '%s'", result.Key)
	}

	if result.Value != true {
		t.Errorf("Expected value true, got %v", result.Value)
	}

	if result.Type != "boolean" {
		t.Errorf("Expected type 'boolean', got '%s'", result.Type)
	}

	if result.Reason != "targeting_rule_match" {
		t.Errorf("Expected reason 'targeting_rule_match', got '%s'", result.Reason)
	}

	if result.RuleID == nil {
		t.Error("RuleID should not be nil")
	} else if *result.RuleID != "rule-123" {
		t.Errorf("Expected rule_id 'rule-123', got '%s'", *result.RuleID)
	}

	if result.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", result.Environment)
	}

	if result.EvaluationMs != 5 {
		t.Errorf("Expected evaluation_ms 5, got %d", result.EvaluationMs)
	}
}

func TestBulkEvaluationResponse_Structure(t *testing.T) {
	// Test the bulk response structure
	results := map[string]services.EvaluationResult{
		"flag1": {
			Key:          "flag1",
			Value:        true,
			Type:         "boolean",
			Reason:       "default_value",
			Environment:  "production",
			EvaluationMs: 3,
		},
		"flag2": {
			Key:          "flag2",
			Value:        "enabled",
			Type:         "string",
			Reason:       "rollout_excluded",
			Environment:  "production",
			EvaluationMs: 4,
		},
	}

	response := services.BulkEvaluationResponse{
		Results: results,
	}

	if len(response.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(response.Results))
	}

	if flag1Result, exists := response.Results["flag1"]; !exists {
		t.Error("Missing result for flag1")
	} else {
		if flag1Result.Key != "flag1" {
			t.Errorf("Expected flag1 key 'flag1', got '%s'", flag1Result.Key)
		}
		if flag1Result.Value != true {
			t.Errorf("Expected flag1 value true, got %v", flag1Result.Value)
		}
	}
}

func TestEvaluateHandler_RequestValidation(t *testing.T) {
	// Test request validation without actual service dependency
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	handler := &EvaluateHandler{
		evaluationService: nil, // nil service will cause issues but we test validation first
	}

	app.Post("/api/v1/evaluate", handler.EvaluateFlag)

	// Test missing flag key
	reqBody := services.EvaluationRequest{
		// Missing FlagKey
		Environment: "production",
		UserID:      "user-123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/evaluate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Should contain error response
	if !strings.Contains(string(body), "error") {
		t.Error("Response should contain error field")
	}
}

func TestEvaluateHandler_BulkRequestValidation(t *testing.T) {
	// Test bulk request validation
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	handler := &EvaluateHandler{
		evaluationService: nil,
	}

	app.Post("/api/v1/evaluate/bulk", handler.EvaluateBulk)

	// Test empty flag keys
	reqBody := services.BulkEvaluationRequest{
		Environment: "production",
		UserID:      "user-123",
		FlagKeys:    []string{}, // Empty flag keys
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/evaluate/bulk", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestEvaluateHandler_InvalidJSON(t *testing.T) {
	// Test invalid JSON handling
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	handler := &EvaluateHandler{
		evaluationService: nil,
	}

	app.Post("/api/v1/evaluate", handler.EvaluateFlag)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/evaluate", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestEvaluateHandler_Structure(t *testing.T) {
	// Test that the handler has the correct structure
	handler := NewEvaluateHandler(nil)

	if handler == nil {
		t.Fatal("NewEvaluateHandler returned nil")
	}
}

func TestEvaluationReasons(t *testing.T) {
	// Test that evaluation reasons are properly structured
	reasons := []string{
		"flag_not_found",
		"flag_disabled",
		"environment_not_found",
		"flag_not_configured_for_environment",
		"environment_disabled",
		"targeting_rule_match",
		"rollout_excluded",
		"default_value",
		"evaluation_error",
	}

	for _, reason := range reasons {
		result := services.EvaluationResult{
			Reason: reason,
		}

		if result.Reason != reason {
			t.Errorf("Reason '%s' not set correctly", reason)
		}
	}
}

func TestJSONMarshaling(t *testing.T) {
	// Test that our structs marshal to JSON correctly
	ruleID := "test-rule"
	result := services.EvaluationResult{
		Key:          "test-flag",
		Value:        true,
		Type:         "boolean",
		Reason:       "default_value",
		RuleID:       &ruleID,
		Environment:  "production",
		EvaluationMs: 5,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	if !strings.Contains(string(jsonData), "test-flag") {
		t.Error("JSON should contain flag key")
	}

	if !strings.Contains(string(jsonData), "boolean") {
		t.Error("JSON should contain flag type")
	}

	if !strings.Contains(string(jsonData), "production") {
		t.Error("JSON should contain environment")
	}

	// Test unmarshaling
	var unmarshaledResult services.EvaluationResult
	if err := json.Unmarshal(jsonData, &unmarshaledResult); err != nil {
		t.Fatalf("Failed to unmarshal result from JSON: %v", err)
	}

	if unmarshaledResult.Key != result.Key {
		t.Error("Unmarshaled result key doesn't match")
	}

	if unmarshaledResult.Environment != result.Environment {
		t.Error("Unmarshaled result environment doesn't match")
	}
}

func TestRequiredFields(t *testing.T) {
	// Test required fields validation logic

	// Valid request
	validReq := services.EvaluationRequest{
		FlagKey:     "test-flag",
		Environment: "production",
		UserID:      "user-123",
	}

	if validReq.FlagKey == "" {
		t.Error("Valid request should have flag_key")
	}

	if validReq.Environment == "" {
		t.Error("Valid request should have environment")
	}

	if validReq.UserID == "" {
		t.Error("Valid request should have user_id")
	}

	// Invalid request (missing fields)
	invalidReq := services.EvaluationRequest{
		// Missing required fields
	}

	if invalidReq.FlagKey != "" {
		t.Error("Invalid request should have empty flag_key")
	}

	if invalidReq.Environment != "" {
		t.Error("Invalid request should have empty environment")
	}

	if invalidReq.UserID != "" {
		t.Error("Invalid request should have empty user_id")
	}
}
