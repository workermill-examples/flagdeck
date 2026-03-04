package handlers

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

func TestCreateSegmentRequest_Validation(t *testing.T) {
	// Test the request structure validation
	rule := models.SegmentRule{
		Attribute: "country",
		Operator:  "equals",
		Value:     "US",
	}

	reqBody := CreateSegmentRequest{
		Key:         "us-users",
		Name:        "US Users",
		Description: "Users from the United States",
		Rules:       []models.SegmentRule{rule},
	}

	if reqBody.Key != "us-users" {
		t.Errorf("Expected key 'us-users', got '%s'", reqBody.Key)
	}

	if reqBody.Name != "US Users" {
		t.Errorf("Expected name 'US Users', got '%s'", reqBody.Name)
	}

	if reqBody.Description != "Users from the United States" {
		t.Errorf("Expected description 'Users from the United States', got '%s'", reqBody.Description)
	}

	if len(reqBody.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(reqBody.Rules))
	}

	if reqBody.Rules[0].Attribute != "country" {
		t.Errorf("Expected rule attribute 'country', got '%s'", reqBody.Rules[0].Attribute)
	}
}

func TestUpdateSegmentRequest_Structure(t *testing.T) {
	name := "Updated Segment"
	description := "Updated description"
	rules := []models.SegmentRule{
		{
			Attribute: "age",
			Operator:  "greater_than",
			Value:     18,
		},
	}

	reqBody := UpdateSegmentRequest{
		Name:        &name,
		Description: &description,
		Rules:       rules,
	}

	if reqBody.Name == nil || *reqBody.Name != "Updated Segment" {
		t.Error("Name field not set correctly")
	}

	if reqBody.Description == nil || *reqBody.Description != "Updated description" {
		t.Error("Description field not set correctly")
	}

	if len(reqBody.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(reqBody.Rules))
	}

	if reqBody.Rules[0].Attribute != "age" {
		t.Errorf("Expected rule attribute 'age', got '%s'", reqBody.Rules[0].Attribute)
	}
}

func TestSegmentRule_Structure(t *testing.T) {
	// Test various rule structures
	tests := []struct {
		name string
		rule models.SegmentRule
	}{
		{
			name: "String equals rule",
			rule: models.SegmentRule{
				Attribute: "country",
				Operator:  "equals",
				Value:     "US",
			},
		},
		{
			name: "Number greater than rule",
			rule: models.SegmentRule{
				Attribute: "age",
				Operator:  "greater_than",
				Value:     21,
			},
		},
		{
			name: "String contains rule",
			rule: models.SegmentRule{
				Attribute: "email",
				Operator:  "contains",
				Value:     "@company.com",
			},
		},
		{
			name: "Array in rule",
			rule: models.SegmentRule{
				Attribute: "plan",
				Operator:  "in",
				Value:     []string{"premium", "enterprise"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rule.Attribute == "" {
				t.Error("Attribute should not be empty")
			}
			if tt.rule.Operator == "" {
				t.Error("Operator should not be empty")
			}
			if tt.rule.Value == nil {
				t.Error("Value should not be nil")
			}
		})
	}
}

func TestSegmentsListResponse_Structure(t *testing.T) {
	// Test the response structure
	segment := models.Segment{
		ID:          primitive.NewObjectID(),
		Key:         "test-segment",
		Name:        "Test Segment",
		Description: "A test segment",
		Rules: []models.SegmentRule{
			{
				Attribute: "country",
				Operator:  "equals",
				Value:     "US",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := SegmentsListResponse{
		Data:  []models.Segment{segment},
		Total: 1,
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(response.Data))
	}

	if response.Total != 1 {
		t.Errorf("Expected total 1, got %d", response.Total)
	}

	if response.Data[0].Key != "test-segment" {
		t.Errorf("Expected key 'test-segment', got '%s'", response.Data[0].Key)
	}

	if len(response.Data[0].Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(response.Data[0].Rules))
	}
}

func TestValidateSegmentRules_ValidOperators(t *testing.T) {
	handler := &SegmentsHandler{}

	validOperators := []string{
		"equals",
		"not_equals",
		"in",
		"not_in",
		"contains",
		"starts_with",
		"ends_with",
		"greater_than",
		"greater_than_equal",
		"less_than",
		"less_than_equal",
	}

	for _, operator := range validOperators {
		t.Run("operator_"+operator, func(t *testing.T) {
			rules := []models.SegmentRule{
				{
					Attribute: "test",
					Operator:  operator,
					Value:     "test_value",
				},
			}

			err := handler.validateSegmentRules(rules)
			if err != nil {
				t.Errorf("Expected operator '%s' to be valid, got error: %v", operator, err)
			}
		})
	}
}

func TestValidateSegmentRules_InvalidOperators(t *testing.T) {
	handler := &SegmentsHandler{}

	invalidOperators := []string{
		"invalid",
		"unknown",
		"equal", // should be "equals"
		"gt",    // should be "greater_than"
	}

	for _, operator := range invalidOperators {
		t.Run("invalid_operator_"+operator, func(t *testing.T) {
			rules := []models.SegmentRule{
				{
					Attribute: "test",
					Operator:  operator,
					Value:     "test_value",
				},
			}

			err := handler.validateSegmentRules(rules)
			if err == nil {
				t.Errorf("Expected operator '%s' to be invalid, but validation passed", operator)
			}
		})
	}
}

func TestValidateSegmentRules_RequiredFields(t *testing.T) {
	handler := &SegmentsHandler{}

	tests := []struct {
		name        string
		rule        models.SegmentRule
		shouldError bool
	}{
		{
			name: "Valid rule",
			rule: models.SegmentRule{
				Attribute: "country",
				Operator:  "equals",
				Value:     "US",
			},
			shouldError: false,
		},
		{
			name: "Missing attribute",
			rule: models.SegmentRule{
				Operator: "equals",
				Value:    "US",
			},
			shouldError: true,
		},
		{
			name: "Missing operator",
			rule: models.SegmentRule{
				Attribute: "country",
				Value:     "US",
			},
			shouldError: true,
		},
		{
			name: "Missing value",
			rule: models.SegmentRule{
				Attribute: "country",
				Operator:  "equals",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateSegmentRules([]models.SegmentRule{tt.rule})

			if tt.shouldError && err == nil {
				t.Error("Expected validation error, but got none")
			} else if !tt.shouldError && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}

func TestSegmentsHandler_Structure(t *testing.T) {
	// Test that the segments handler has the correct structure
	handler := &SegmentsHandler{
		segmentsCollection: nil, // Would be a real collection in practice
		auditService:       nil, // Would be a real service in practice
	}

	if handler == nil {
		t.Fatal("SegmentsHandler should not be nil")
	}

	// Test NewSegmentsHandler constructor
	newHandler := NewSegmentsHandler(nil, nil)
	if newHandler == nil {
		t.Fatal("NewSegmentsHandler returned nil")
	}
}
