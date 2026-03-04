package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Experiment struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Key         string                 `bson:"key" json:"key"`
	Name        string                 `bson:"name" json:"name"`
	Description string                 `bson:"description" json:"description"`
	FlagKey     string                 `bson:"flag_key" json:"flag_key"`
	Environment string                 `bson:"environment" json:"environment"`
	Status      string                 `bson:"status" json:"status"`
	StartDate   *time.Time             `bson:"start_date" json:"start_date"`
	EndDate     *time.Time             `bson:"end_date" json:"end_date"`
	Variants    []ExperimentVariant    `bson:"variants" json:"variants"`
	Results     map[string]interface{} `bson:"results" json:"results"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

type ExperimentVariant struct {
	Name         string      `bson:"name" json:"name"`
	Value        interface{} `bson:"value" json:"value"`
	TrafficSplit int         `bson:"traffic_split" json:"traffic_split"`
}

type VariantResults struct {
	VariantName string  `bson:"variant_name" json:"variant_name"`
	Impressions int     `bson:"impressions" json:"impressions"`
	Conversions int     `bson:"conversions" json:"conversions"`
	Revenue     float64 `bson:"revenue" json:"revenue"`
}
