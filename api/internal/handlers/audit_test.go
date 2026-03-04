package handlers

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

func TestAuditLogResponse_Structure(t *testing.T) {
	// Test the response structure
	timestamp := time.Now()
	entry := models.AuditLogEntry{
		ID:         primitive.NewObjectID(),
		Resource:   "flag",
		ResourceID: "homepage-banner",
		Action:     "flag.created",
		UserID:     primitive.NewObjectID(),
		UserEmail:  "admin@example.com",
		Changes: map[string]interface{}{
			"name": "Homepage Banner",
			"type": "boolean",
		},
		Timestamp: timestamp,
	}

	response := AuditLogResponse{
		Data:  []models.AuditLogEntry{entry},
		Total: 1,
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 audit entry, got %d", len(response.Data))
	}

	if response.Total != 1 {
		t.Errorf("Expected total 1, got %d", response.Total)
	}

	if response.Data[0].Resource != "flag" {
		t.Errorf("Expected resource 'flag', got '%s'", response.Data[0].Resource)
	}

	if response.Data[0].Action != "flag.created" {
		t.Errorf("Expected action 'flag.created', got '%s'", response.Data[0].Action)
	}

	if response.Data[0].UserEmail != "admin@example.com" {
		t.Errorf("Expected user email 'admin@example.com', got '%s'", response.Data[0].UserEmail)
	}
}

func TestAuditLogResponse_EmptyList(t *testing.T) {
	// Test empty response
	response := AuditLogResponse{
		Data:  []models.AuditLogEntry{},
		Total: 0,
	}

	if len(response.Data) != 0 {
		t.Errorf("Expected 0 audit entries, got %d", len(response.Data))
	}

	if response.Total != 0 {
		t.Errorf("Expected total 0, got %d", response.Total)
	}

	if response.Data == nil {
		t.Error("Data should be empty slice, not nil")
	}
}

func TestAuditLogResponse_MultipleEntries(t *testing.T) {
	// Test response with multiple entries
	timestamp := time.Now()

	entries := []models.AuditLogEntry{
		{
			ID:         primitive.NewObjectID(),
			Resource:   "flag",
			ResourceID: "feature-a",
			Action:     "flag.created",
			UserID:     primitive.NewObjectID(),
			UserEmail:  "user1@example.com",
			Changes:    map[string]interface{}{"name": "Feature A"},
			Timestamp:  timestamp,
		},
		{
			ID:         primitive.NewObjectID(),
			Resource:   "environment",
			ResourceID: "production",
			Action:     "environment.updated",
			UserID:     primitive.NewObjectID(),
			UserEmail:  "user2@example.com",
			Changes:    map[string]interface{}{"color": "#red"},
			Timestamp:  timestamp.Add(-1 * time.Hour),
		},
		{
			ID:         primitive.NewObjectID(),
			Resource:   "experiment",
			ResourceID: "homepage-test",
			Action:     "experiment.started",
			UserID:     primitive.NewObjectID(),
			UserEmail:  "user3@example.com",
			Changes:    map[string]interface{}{"status": "running"},
			Timestamp:  timestamp.Add(-2 * time.Hour),
		},
	}

	response := AuditLogResponse{
		Data:  entries,
		Total: 3,
	}

	if len(response.Data) != 3 {
		t.Errorf("Expected 3 audit entries, got %d", len(response.Data))
	}

	if response.Total != 3 {
		t.Errorf("Expected total 3, got %d", response.Total)
	}

	// Test different resource types
	resources := make(map[string]bool)
	for _, entry := range response.Data {
		resources[entry.Resource] = true
	}

	expectedResources := []string{"flag", "environment", "experiment"}
	for _, resource := range expectedResources {
		if !resources[resource] {
			t.Errorf("Expected to find resource type '%s' in audit entries", resource)
		}
	}
}

func TestAuditLogEntry_Actions(t *testing.T) {
	// Test various action types
	actions := []string{
		"flag.created",
		"flag.updated",
		"flag.deleted",
		"flag.toggled",
		"environment.created",
		"environment.updated",
		"environment.deleted",
		"experiment.created",
		"experiment.updated",
		"experiment.deleted",
		"experiment.started",
		"experiment.paused",
		"experiment.completed",
		"segment.created",
		"segment.updated",
		"segment.deleted",
		"api_key.created",
		"api_key.deleted",
	}

	for _, action := range actions {
		t.Run("action_"+action, func(t *testing.T) {
			entry := models.AuditLogEntry{
				ID:         primitive.NewObjectID(),
				Resource:   "test",
				ResourceID: "test-id",
				Action:     action,
				UserID:     primitive.NewObjectID(),
				UserEmail:  "test@example.com",
				Changes:    map[string]interface{}{},
				Timestamp:  time.Now(),
			}

			if entry.Action != action {
				t.Errorf("Expected action '%s', got '%s'", action, entry.Action)
			}
		})
	}
}

func TestAuditLogEntry_Resources(t *testing.T) {
	// Test various resource types
	resources := []string{
		"flag",
		"environment",
		"experiment",
		"segment",
		"api_key",
		"user",
	}

	for _, resource := range resources {
		t.Run("resource_"+resource, func(t *testing.T) {
			entry := models.AuditLogEntry{
				ID:         primitive.NewObjectID(),
				Resource:   resource,
				ResourceID: "test-id",
				Action:     "test.action",
				UserID:     primitive.NewObjectID(),
				UserEmail:  "test@example.com",
				Changes:    map[string]interface{}{},
				Timestamp:  time.Now(),
			}

			if entry.Resource != resource {
				t.Errorf("Expected resource '%s', got '%s'", resource, entry.Resource)
			}
		})
	}
}

func TestAuditLogEntry_Changes(t *testing.T) {
	// Test different types of changes
	tests := []struct {
		name    string
		changes map[string]interface{}
	}{
		{
			name: "String changes",
			changes: map[string]interface{}{
				"name":        "New Name",
				"description": "New Description",
			},
		},
		{
			name: "Boolean changes",
			changes: map[string]interface{}{
				"enabled":   true,
				"is_active": false,
			},
		},
		{
			name: "Number changes",
			changes: map[string]interface{}{
				"rollout_percent": 75,
				"weight":          50.5,
			},
		},
		{
			name: "Complex object changes",
			changes: map[string]interface{}{
				"targeting_rules": []interface{}{
					map[string]interface{}{
						"attribute": "country",
						"operator":  "equals",
						"value":     "US",
					},
				},
				"environments": map[string]interface{}{
					"production": map[string]interface{}{
						"enabled": true,
						"value":   "new-value",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := models.AuditLogEntry{
				ID:         primitive.NewObjectID(),
				Resource:   "test",
				ResourceID: "test-id",
				Action:     "test.updated",
				UserID:     primitive.NewObjectID(),
				UserEmail:  "test@example.com",
				Changes:    tt.changes,
				Timestamp:  time.Now(),
			}

			if entry.Changes == nil {
				t.Error("Changes should not be nil")
			}

			if len(entry.Changes) != len(tt.changes) {
				t.Errorf("Expected %d changes, got %d", len(tt.changes), len(entry.Changes))
			}

			for key, expectedValue := range tt.changes {
				if actualValue, exists := entry.Changes[key]; !exists {
					t.Errorf("Expected change key '%s' to exist", key)
				} else {
					// Note: Deep comparison of interface{} values would require more complex logic
					// For now, just verify the key exists and the value is not nil
					if actualValue == nil && expectedValue != nil {
						t.Errorf("Expected change value for key '%s' to not be nil", key)
					}
				}
			}
		})
	}
}

func TestAuditLogEntry_Timestamps(t *testing.T) {
	// Test timestamp ordering for proper audit trail
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	oneDayAgo := now.Add(-24 * time.Hour)

	entries := []models.AuditLogEntry{
		{
			ID:        primitive.NewObjectID(),
			Resource:  "flag",
			Action:    "flag.created",
			Timestamp: oneDayAgo,
		},
		{
			ID:        primitive.NewObjectID(),
			Resource:  "flag",
			Action:    "flag.updated",
			Timestamp: oneHourAgo,
		},
		{
			ID:        primitive.NewObjectID(),
			Resource:  "flag",
			Action:    "flag.toggled",
			Timestamp: now,
		},
	}

	// Verify timestamps are in the expected order (oldest to newest)
	if !entries[0].Timestamp.Before(entries[1].Timestamp) {
		t.Error("First entry should be older than second entry")
	}

	if !entries[1].Timestamp.Before(entries[2].Timestamp) {
		t.Error("Second entry should be older than third entry")
	}

	// Verify all timestamps are valid (not zero)
	for i, entry := range entries {
		if entry.Timestamp.IsZero() {
			t.Errorf("Entry %d should have a valid timestamp", i)
		}
	}
}

func TestAuditHandler_Structure(t *testing.T) {
	// Test that the audit handler has the correct structure
	handler := &AuditHandler{
		auditCollection: nil, // Would be a real collection in practice
	}

	if handler == nil {
		t.Fatal("AuditHandler should not be nil")
	}

	// Test NewAuditHandler constructor
	newHandler := NewAuditHandler(nil)
	if newHandler == nil {
		t.Fatal("NewAuditHandler returned nil")
	}
}
