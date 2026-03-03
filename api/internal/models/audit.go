package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLogEntry struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Action     string                 `bson:"action" json:"action"`
	Resource   string                 `bson:"resource" json:"resource"`
	ResourceID string                 `bson:"resource_id" json:"resource_id"`
	UserID     primitive.ObjectID     `bson:"user_id" json:"user_id"`
	UserEmail  string                 `bson:"user_email" json:"user_email"`
	Changes    map[string]interface{} `bson:"changes" json:"changes"`
	Metadata   map[string]interface{} `bson:"metadata" json:"metadata"`
	Timestamp  time.Time              `bson:"timestamp" json:"timestamp"`
}
