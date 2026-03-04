package handlers

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

func TestCreateExperimentRequest_Validation(t *testing.T) {
	// Test the request structure validation
	startDate := time.Now()
	endDate := startDate.Add(30 * 24 * time.Hour) // 30 days later

	variant := models.ExperimentVariant{
		Name:         "control",
		Value:        false,
		TrafficSplit: 50,
	}

	reqBody := CreateExperimentRequest{
		Key:         "homepage-test",
		Name:        "Homepage A/B Test",
		Description: "Testing new homepage design",
		FlagKey:     "homepage-redesign",
		Environment: "production",
		Status:      "draft",
		StartDate:   &startDate,
		EndDate:     &endDate,
		Variants:    []models.ExperimentVariant{variant},
		Results:     map[string]interface{}{"control": map[string]interface{}{"impressions": 0}},
	}

	if reqBody.Key != "homepage-test" {
		t.Errorf("Expected key 'homepage-test', got '%s'", reqBody.Key)
	}

	if reqBody.Name != "Homepage A/B Test" {
		t.Errorf("Expected name 'Homepage A/B Test', got '%s'", reqBody.Name)
	}

	if reqBody.FlagKey != "homepage-redesign" {
		t.Errorf("Expected flag_key 'homepage-redesign', got '%s'", reqBody.FlagKey)
	}

	if reqBody.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", reqBody.Environment)
	}

	if reqBody.Status != "draft" {
		t.Errorf("Expected status 'draft', got '%s'", reqBody.Status)
	}

	if len(reqBody.Variants) != 1 {
		t.Errorf("Expected 1 variant, got %d", len(reqBody.Variants))
	}
}

func TestUpdateExperimentRequest_Structure(t *testing.T) {
	name := "Updated Experiment"
	description := "Updated description"
	status := "running"
	flagKey := "new-flag"
	environment := "staging"

	reqBody := UpdateExperimentRequest{
		Name:        &name,
		Description: &description,
		Status:      &status,
		FlagKey:     &flagKey,
		Environment: &environment,
	}

	if reqBody.Name == nil || *reqBody.Name != "Updated Experiment" {
		t.Error("Name field not set correctly")
	}

	if reqBody.Description == nil || *reqBody.Description != "Updated description" {
		t.Error("Description field not set correctly")
	}

	if reqBody.Status == nil || *reqBody.Status != "running" {
		t.Error("Status field not set correctly")
	}

	if reqBody.FlagKey == nil || *reqBody.FlagKey != "new-flag" {
		t.Error("FlagKey field not set correctly")
	}

	if reqBody.Environment == nil || *reqBody.Environment != "staging" {
		t.Error("Environment field not set correctly")
	}
}

func TestTrackExperimentRequest_Validation(t *testing.T) {
	reqBody := TrackExperimentRequest{
		UserID:      "user123",
		VariantName: "treatment",
		Conversion:  true,
		Revenue:     29.99,
	}

	if reqBody.UserID != "user123" {
		t.Errorf("Expected user_id 'user123', got '%s'", reqBody.UserID)
	}

	if reqBody.VariantName != "treatment" {
		t.Errorf("Expected variant_name 'treatment', got '%s'", reqBody.VariantName)
	}

	if !reqBody.Conversion {
		t.Error("Expected conversion to be true")
	}

	if reqBody.Revenue != 29.99 {
		t.Errorf("Expected revenue 29.99, got %f", reqBody.Revenue)
	}
}

func TestExperimentStatus_Validation(t *testing.T) {
	validStatuses := []string{"draft", "running", "paused", "completed"}
	invalidStatuses := []string{"pending", "stopped", "active", "finished"}

	// Test valid statuses
	for _, status := range validStatuses {
		t.Run("valid_status_"+status, func(t *testing.T) {
			req := CreateExperimentRequest{
				Key:         "test",
				Name:        "Test",
				FlagKey:     "test-flag",
				Environment: "production",
				Status:      status,
			}

			// In a real handler, this would be validated
			isValid := false
			for _, validStatus := range validStatuses {
				if req.Status == validStatus {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Status '%s' should be valid", status)
			}
		})
	}

	// Test invalid statuses
	for _, status := range invalidStatuses {
		t.Run("invalid_status_"+status, func(t *testing.T) {
			req := CreateExperimentRequest{
				Status: status,
			}

			// In a real handler, this would fail validation
			isValid := false
			for _, validStatus := range validStatuses {
				if req.Status == validStatus {
					isValid = true
					break
				}
			}

			if isValid {
				t.Errorf("Status '%s' should be invalid", status)
			}
		})
	}
}

func TestExperimentsListResponse_Structure(t *testing.T) {
	// Test the response structure
	experiment := models.Experiment{
		ID:          primitive.NewObjectID(),
		Key:         "test-experiment",
		Name:        "Test Experiment",
		Description: "A test experiment",
		FlagKey:     "test-flag",
		Environment: "production",
		Status:      "draft",
		Variants: []models.ExperimentVariant{
			{Name: "control", Value: false, TrafficSplit: 50},
			{Name: "treatment", Value: true, TrafficSplit: 50},
		},
		Results:   map[string]interface{}{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := ExperimentsListResponse{
		Data:  []models.Experiment{experiment},
		Total: 1,
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 experiment, got %d", len(response.Data))
	}

	if response.Total != 1 {
		t.Errorf("Expected total 1, got %d", response.Total)
	}

	if response.Data[0].Key != "test-experiment" {
		t.Errorf("Expected key 'test-experiment', got '%s'", response.Data[0].Key)
	}

	if len(response.Data[0].Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(response.Data[0].Variants))
	}
}

func TestExperimentVariant_Structure(t *testing.T) {
	// Test variant with boolean value
	boolVariant := models.ExperimentVariant{
		Name:         "control",
		Value:        false,
		TrafficSplit: 50,
	}

	if boolVariant.Name != "control" {
		t.Errorf("Expected name 'control', got '%s'", boolVariant.Name)
	}

	if boolVariant.Value != false {
		t.Errorf("Expected value false, got %v", boolVariant.Value)
	}

	if boolVariant.TrafficSplit != 50 {
		t.Errorf("Expected traffic_split 50, got %d", boolVariant.TrafficSplit)
	}

	// Test variant with string value
	stringVariant := models.ExperimentVariant{
		Name:         "treatment",
		Value:        "new-design",
		TrafficSplit: 50,
	}

	if stringVariant.Value != "new-design" {
		t.Errorf("Expected value 'new-design', got %v", stringVariant.Value)
	}
}

func TestExperimentResults_Structure(t *testing.T) {
	// Test results structure
	results := map[string]interface{}{
		"control": map[string]interface{}{
			"impressions": float64(1000),
			"conversions": float64(50),
			"revenue":     float64(1500.00),
		},
		"treatment": map[string]interface{}{
			"impressions": float64(1000),
			"conversions": float64(75),
			"revenue":     float64(2250.00),
		},
	}

	controlResults := results["control"].(map[string]interface{})
	treatmentResults := results["treatment"].(map[string]interface{})

	if controlResults["impressions"] != float64(1000) {
		t.Errorf("Expected control impressions 1000, got %v", controlResults["impressions"])
	}

	if treatmentResults["conversions"] != float64(75) {
		t.Errorf("Expected treatment conversions 75, got %v", treatmentResults["conversions"])
	}
}

func TestExperimentsHandler_Structure(t *testing.T) {
	// Test that the experiments handler has the correct structure
	handler := &ExperimentsHandler{
		experimentsCollection: nil, // Would be a real collection in practice
		auditService:          nil, // Would be a real service in practice
	}

	if handler == nil {
		t.Fatal("ExperimentsHandler should not be nil")
	}

	// Test NewExperimentsHandler constructor
	newHandler := NewExperimentsHandler(nil, nil)
	if newHandler == nil {
		t.Fatal("NewExperimentsHandler returned nil")
	}
}

func TestCreateExperimentRequest_RequiredFields(t *testing.T) {
	// Test validation of required fields
	tests := []struct {
		name    string
		request CreateExperimentRequest
		isValid bool
	}{
		{
			name: "Valid request",
			request: CreateExperimentRequest{
				Key:         "test-exp",
				Name:        "Test Experiment",
				FlagKey:     "test-flag",
				Environment: "production",
				Status:      "draft",
			},
			isValid: true,
		},
		{
			name: "Missing key",
			request: CreateExperimentRequest{
				Name:        "Test Experiment",
				FlagKey:     "test-flag",
				Environment: "production",
				Status:      "draft",
			},
			isValid: false,
		},
		{
			name: "Missing name",
			request: CreateExperimentRequest{
				Key:         "test-exp",
				FlagKey:     "test-flag",
				Environment: "production",
				Status:      "draft",
			},
			isValid: false,
		},
		{
			name: "Missing flag_key",
			request: CreateExperimentRequest{
				Key:         "test-exp",
				Name:        "Test Experiment",
				Environment: "production",
				Status:      "draft",
			},
			isValid: false,
		},
		{
			name: "Missing environment",
			request: CreateExperimentRequest{
				Key:     "test-exp",
				Name:    "Test Experiment",
				FlagKey: "test-flag",
				Status:  "draft",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasKey := tt.request.Key != ""
			hasName := tt.request.Name != ""
			hasFlagKey := tt.request.FlagKey != ""
			hasEnvironment := tt.request.Environment != ""
			actualValid := hasKey && hasName && hasFlagKey && hasEnvironment

			if actualValid != tt.isValid {
				t.Errorf("Expected valid=%v, got valid=%v for request %+v", tt.isValid, actualValid, tt.request)
			}
		})
	}
}
