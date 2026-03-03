package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

type AuditService struct {
	auditCollection *mongo.Collection
}

type AuditEntry struct {
	Action     string
	Resource   string
	ResourceID string
	UserID     primitive.ObjectID
	UserEmail  string
	Changes    map[string]interface{}
	Metadata   map[string]interface{}
}

func NewAuditService(auditCollection *mongo.Collection) *AuditService {
	return &AuditService{
		auditCollection: auditCollection,
	}
}

// LogEntry records an audit log entry for a mutating operation
func (a *AuditService) LogEntry(entry AuditEntry) error {
	auditLogEntry := models.AuditLogEntry{
		ID:         primitive.NewObjectID(),
		Action:     entry.Action,
		Resource:   entry.Resource,
		ResourceID: entry.ResourceID,
		UserID:     entry.UserID,
		UserEmail:  entry.UserEmail,
		Changes:    entry.Changes,
		Metadata:   entry.Metadata,
		Timestamp:  time.Now(),
	}

	_, err := a.auditCollection.InsertOne(context.Background(), auditLogEntry)
	if err != nil {
		log.Printf("Failed to log audit entry: %v", err)
		// Don't fail the main operation if audit logging fails
		return nil
	}

	return nil
}

// LogCreate logs a resource creation
func (a *AuditService) LogCreate(resource, resourceID string, userID primitive.ObjectID, userEmail string, data map[string]interface{}) {
	entry := AuditEntry{
		Action:     "create",
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		UserEmail:  userEmail,
		Changes:    data,
		Metadata:   map[string]interface{}{},
	}
	a.LogEntry(entry)
}

// LogUpdate logs a resource update with before/after changes
func (a *AuditService) LogUpdate(resource, resourceID string, userID primitive.ObjectID, userEmail string, before, after map[string]interface{}) {
	changes := make(map[string]interface{})

	// Track what changed by comparing before and after
	for key, afterValue := range after {
		if beforeValue, exists := before[key]; !exists || beforeValue != afterValue {
			changes[key] = map[string]interface{}{
				"before": beforeValue,
				"after":  afterValue,
			}
		}
	}

	// Check for removed fields
	for key, beforeValue := range before {
		if _, exists := after[key]; !exists {
			changes[key] = map[string]interface{}{
				"before": beforeValue,
				"after":  nil,
			}
		}
	}

	entry := AuditEntry{
		Action:     "update",
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		UserEmail:  userEmail,
		Changes:    changes,
		Metadata:   map[string]interface{}{},
	}
	a.LogEntry(entry)
}

// LogDelete logs a resource deletion
func (a *AuditService) LogDelete(resource, resourceID string, userID primitive.ObjectID, userEmail string, deletedData map[string]interface{}) {
	entry := AuditEntry{
		Action:     "delete",
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		UserEmail:  userEmail,
		Changes:    deletedData,
		Metadata:   map[string]interface{}{},
	}
	a.LogEntry(entry)
}

// LogToggle logs a flag toggle operation (special case)
func (a *AuditService) LogToggle(resource, resourceID string, userID primitive.ObjectID, userEmail string, environment string, fromState, toState bool) {
	changes := map[string]interface{}{
		"enabled": map[string]interface{}{
			"before": fromState,
			"after":  toState,
		},
	}

	metadata := map[string]interface{}{}
	if environment != "" {
		metadata["environment"] = environment
	} else {
		metadata["scope"] = "global"
	}

	entry := AuditEntry{
		Action:     "toggle",
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		UserEmail:  userEmail,
		Changes:    changes,
		Metadata:   metadata,
	}
	a.LogEntry(entry)
}

// LogCustom logs a custom operation with arbitrary action and metadata
func (a *AuditService) LogCustom(action, resource, resourceID string, userID primitive.ObjectID, userEmail string, changes, metadata map[string]interface{}) {
	entry := AuditEntry{
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		UserEmail:  userEmail,
		Changes:    changes,
		Metadata:   metadata,
	}
	a.LogEntry(entry)
}

// Helper functions to create changes maps from structs/maps

// CreateChangesMap creates a changes map suitable for audit logging from any data
func (a *AuditService) CreateChangesMap(data interface{}) map[string]interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		return v
	case models.User:
		return map[string]interface{}{
			"email": v.Email,
			"name":  v.Name,
			"role":  v.Role,
		}
	case models.AuditLogEntry:
		return map[string]interface{}{
			"action":      v.Action,
			"resource":    v.Resource,
			"resource_id": v.ResourceID,
		}
	default:
		// For unknown types, return empty map
		return map[string]interface{}{}
	}
}

// GetResourceName returns standardized resource names for audit logging
func (a *AuditService) GetResourceName(resourceType string) string {
	switch resourceType {
	case "flag", "flags":
		return "flag"
	case "environment", "environments":
		return "environment"
	case "segment", "segments":
		return "segment"
	case "experiment", "experiments":
		return "experiment"
	case "user", "users":
		return "user"
	case "api_key", "api_keys", "apikey":
		return "api_key"
	default:
		return resourceType
	}
}
