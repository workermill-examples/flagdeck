package handlers

import (
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestStatsResponse_JSONMarshaling(t *testing.T) {
	response := StatsResponse{
		Flags:           10,
		Environments:    5,
		Segments:        3,
		Experiments:     7,
		APIKeys:         2,
		AuditLogEntries: 100,
	}

	// Test that the struct marshals to JSON correctly
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal StatsResponse to JSON: %v", err)
	}

	// Unmarshal back to verify field names
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check that all expected fields are present with correct names
	expectedFields := map[string]int64{
		"flags":             10,
		"environments":      5,
		"segments":          3,
		"experiments":       7,
		"api_keys":          2,
		"audit_log_entries": 100,
	}

	for fieldName, expectedValue := range expectedFields {
		value, exists := unmarshaled[fieldName]
		if !exists {
			t.Errorf("Expected field '%s' not found in JSON", fieldName)
			continue
		}

		// JSON unmarshals numbers as float64
		if floatValue, ok := value.(float64); ok {
			if int64(floatValue) != expectedValue {
				t.Errorf("Expected field '%s' to have value %d, got %v", fieldName, expectedValue, floatValue)
			}
		} else {
			t.Errorf("Expected field '%s' to be a number, got %T", fieldName, value)
		}
	}
}

func TestNewStatsHandler(t *testing.T) {
	// Test that the constructor returns a non-nil handler with all fields set
	var (
		flagsCollection        *mongo.Collection
		environmentsCollection *mongo.Collection
		segmentsCollection     *mongo.Collection
		experimentsCollection  *mongo.Collection
		apiKeysCollection      *mongo.Collection
		auditLogCollection     *mongo.Collection
	)

	handler := NewStatsHandler(
		flagsCollection,
		environmentsCollection,
		segmentsCollection,
		experimentsCollection,
		apiKeysCollection,
		auditLogCollection,
	)

	if handler == nil {
		t.Fatal("NewStatsHandler returned nil")
	}

	if handler.FlagsCollection != flagsCollection {
		t.Error("FlagsCollection field not set correctly")
	}

	if handler.EnvironmentsCollection != environmentsCollection {
		t.Error("EnvironmentsCollection field not set correctly")
	}

	if handler.SegmentsCollection != segmentsCollection {
		t.Error("SegmentsCollection field not set correctly")
	}

	if handler.ExperimentsCollection != experimentsCollection {
		t.Error("ExperimentsCollection field not set correctly")
	}

	if handler.APIKeysCollection != apiKeysCollection {
		t.Error("APIKeysCollection field not set correctly")
	}

	if handler.AuditLogCollection != auditLogCollection {
		t.Error("AuditLogCollection field not set correctly")
	}
}

func TestStatsResponse_Structure(t *testing.T) {
	// Test the response structure fields and types
	response := StatsResponse{
		Flags:           42,
		Environments:    7,
		Segments:        15,
		Experiments:     3,
		APIKeys:         8,
		AuditLogEntries: 250,
	}

	if response.Flags != 42 {
		t.Errorf("Expected Flags field to be 42, got %d", response.Flags)
	}

	if response.Environments != 7 {
		t.Errorf("Expected Environments field to be 7, got %d", response.Environments)
	}

	if response.Segments != 15 {
		t.Errorf("Expected Segments field to be 15, got %d", response.Segments)
	}

	if response.Experiments != 3 {
		t.Errorf("Expected Experiments field to be 3, got %d", response.Experiments)
	}

	if response.APIKeys != 8 {
		t.Errorf("Expected APIKeys field to be 8, got %d", response.APIKeys)
	}

	if response.AuditLogEntries != 250 {
		t.Errorf("Expected AuditLogEntries field to be 250, got %d", response.AuditLogEntries)
	}
}
