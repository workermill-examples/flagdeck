package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flag struct {
	ID           primitive.ObjectID         `bson:"_id,omitempty" json:"id"`
	Key          string                     `bson:"key" json:"key"`
	Name         string                     `bson:"name" json:"name"`
	Description  string                     `bson:"description" json:"description"`
	IsActive     bool                       `bson:"is_active" json:"is_active"`
	Type         string                     `bson:"type" json:"type"`
	Environments map[string]FlagEnvironment `bson:"environments" json:"environments"`
	CreatedAt    time.Time                  `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time                  `bson:"updated_at" json:"updated_at"`
}

type FlagEnvironment struct {
	Enabled        bool            `bson:"enabled" json:"enabled"`
	DefaultValue   interface{}     `bson:"default_value" json:"default_value"`
	RolloutPercent int             `bson:"rollout_percent" json:"rollout_percent"`
	TargetingRules []TargetingRule `bson:"targeting_rules" json:"targeting_rules"`
}

type TargetingRule struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Priority   int                `bson:"priority" json:"priority"`
	Conditions []Condition        `bson:"conditions" json:"conditions"`
	Value      interface{}        `bson:"value" json:"value"`
}

type Condition struct {
	Attribute string      `bson:"attribute" json:"attribute"`
	Operator  string      `bson:"operator" json:"operator"`
	Value     interface{} `bson:"value" json:"value"`
}
