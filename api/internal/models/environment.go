package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Environment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key         string             `bson:"key" json:"key"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Color       string             `bson:"color" json:"color"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
