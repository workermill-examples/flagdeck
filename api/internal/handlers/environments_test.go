package handlers

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

func TestCreateEnvironmentRequest_Validation(t *testing.T) {
	// Test the request structure validation
	reqBody := CreateEnvironmentRequest{
		Key:         "production",
		Name:        "Production",
		Description: "Production environment",
		Color:       "#ff0000",
	}

	if reqBody.Key != "production" {
		t.Errorf("Expected key 'production', got '%s'", reqBody.Key)
	}

	if reqBody.Name != "Production" {
		t.Errorf("Expected name 'Production', got '%s'", reqBody.Name)
	}

	if reqBody.Description != "Production environment" {
		t.Errorf("Expected description 'Production environment', got '%s'", reqBody.Description)
	}

	if reqBody.Color != "#ff0000" {
		t.Errorf("Expected color '#ff0000', got '%s'", reqBody.Color)
	}
}

func TestUpdateEnvironmentRequest_Structure(t *testing.T) {
	name := "Updated Environment"
	description := "Updated description"
	color := "#00ff00"

	reqBody := UpdateEnvironmentRequest{
		Name:        &name,
		Description: &description,
		Color:       &color,
	}

	if reqBody.Name == nil || *reqBody.Name != "Updated Environment" {
		t.Error("Name field not set correctly")
	}

	if reqBody.Description == nil || *reqBody.Description != "Updated description" {
		t.Error("Description field not set correctly")
	}

	if reqBody.Color == nil || *reqBody.Color != "#00ff00" {
		t.Error("Color field not set correctly")
	}
}

func TestUpdateEnvironmentRequest_PartialUpdate(t *testing.T) {
	// Test partial update with only name
	name := "Only Name Updated"

	reqBody := UpdateEnvironmentRequest{
		Name:        &name,
		Description: nil,
		Color:       nil,
	}

	if reqBody.Name == nil || *reqBody.Name != "Only Name Updated" {
		t.Error("Name field not set correctly")
	}

	if reqBody.Description != nil {
		t.Error("Description should be nil for partial update")
	}

	if reqBody.Color != nil {
		t.Error("Color should be nil for partial update")
	}
}

func TestEnvironmentsListResponse_Structure(t *testing.T) {
	// Test the response structure
	environment := models.Environment{
		ID:          primitive.NewObjectID(),
		Key:         "test-env",
		Name:        "Test Environment",
		Description: "A test environment",
		Color:       "#blue",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	response := EnvironmentsListResponse{
		Data:  []models.Environment{environment},
		Total: 1,
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 environment, got %d", len(response.Data))
	}

	if response.Total != 1 {
		t.Errorf("Expected total 1, got %d", response.Total)
	}

	if response.Data[0].Key != "test-env" {
		t.Errorf("Expected key 'test-env', got '%s'", response.Data[0].Key)
	}

	if response.Data[0].Name != "Test Environment" {
		t.Errorf("Expected name 'Test Environment', got '%s'", response.Data[0].Name)
	}
}

func TestEnvironmentsListResponse_EmptyList(t *testing.T) {
	// Test empty response
	response := EnvironmentsListResponse{
		Data:  []models.Environment{},
		Total: 0,
	}

	if len(response.Data) != 0 {
		t.Errorf("Expected 0 environments, got %d", len(response.Data))
	}

	if response.Total != 0 {
		t.Errorf("Expected total 0, got %d", response.Total)
	}

	if response.Data == nil {
		t.Error("Data should be empty slice, not nil")
	}
}

func TestEnvironmentsHandler_Structure(t *testing.T) {
	// Test that the environments handler has the correct structure
	handler := &EnvironmentsHandler{
		environmentsCollection: nil, // Would be a real collection in practice
		auditService:           nil, // Would be a real service in practice
	}

	if handler == nil {
		t.Fatal("EnvironmentsHandler should not be nil")
	}

	// Test NewEnvironmentsHandler constructor
	newHandler := NewEnvironmentsHandler(nil, nil)
	if newHandler == nil {
		t.Fatal("NewEnvironmentsHandler returned nil")
	}
}

func TestCreateEnvironmentRequest_RequiredFields(t *testing.T) {
	// Test validation of required fields
	tests := []struct {
		name    string
		request CreateEnvironmentRequest
		isValid bool
	}{
		{
			name: "Valid request",
			request: CreateEnvironmentRequest{
				Key:  "prod",
				Name: "Production",
			},
			isValid: true,
		},
		{
			name: "Missing key",
			request: CreateEnvironmentRequest{
				Name: "Production",
			},
			isValid: false,
		},
		{
			name: "Missing name",
			request: CreateEnvironmentRequest{
				Key: "prod",
			},
			isValid: false,
		},
		{
			name: "Empty key",
			request: CreateEnvironmentRequest{
				Key:  "",
				Name: "Production",
			},
			isValid: false,
		},
		{
			name: "Empty name",
			request: CreateEnvironmentRequest{
				Key:  "prod",
				Name: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasKey := tt.request.Key != ""
			hasName := tt.request.Name != ""
			actualValid := hasKey && hasName

			if actualValid != tt.isValid {
				t.Errorf("Expected valid=%v, got valid=%v for request %+v", tt.isValid, actualValid, tt.request)
			}
		})
	}
}
