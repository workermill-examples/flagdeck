package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsHandler struct {
	FlagsCollection        *mongo.Collection
	EnvironmentsCollection *mongo.Collection
	SegmentsCollection     *mongo.Collection
	ExperimentsCollection  *mongo.Collection
	APIKeysCollection      *mongo.Collection
	AuditLogCollection     *mongo.Collection
}

type StatsResponse struct {
	Flags           int64 `json:"flags"`
	Environments    int64 `json:"environments"`
	Segments        int64 `json:"segments"`
	Experiments     int64 `json:"experiments"`
	APIKeys         int64 `json:"api_keys"`
	AuditLogEntries int64 `json:"audit_log_entries"`
}

func NewStatsHandler(
	flagsCollection *mongo.Collection,
	environmentsCollection *mongo.Collection,
	segmentsCollection *mongo.Collection,
	experimentsCollection *mongo.Collection,
	apiKeysCollection *mongo.Collection,
	auditLogCollection *mongo.Collection,
) *StatsHandler {
	return &StatsHandler{
		FlagsCollection:        flagsCollection,
		EnvironmentsCollection: environmentsCollection,
		SegmentsCollection:     segmentsCollection,
		ExperimentsCollection:  experimentsCollection,
		APIKeysCollection:      apiKeysCollection,
		AuditLogCollection:     auditLogCollection,
	}
}

func (h *StatsHandler) GetStats(c *fiber.Ctx) error {
	ctx := context.Background()

	// Count documents in each collection
	flagsCount, err := h.FlagsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count flags")
	}

	environmentsCount, err := h.EnvironmentsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count environments")
	}

	segmentsCount, err := h.SegmentsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count segments")
	}

	experimentsCount, err := h.ExperimentsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count experiments")
	}

	apiKeysCount, err := h.APIKeysCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count API keys")
	}

	auditLogCount, err := h.AuditLogCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return middleware.NewDatabaseError("Failed to count audit log entries")
	}

	response := StatsResponse{
		Flags:           flagsCount,
		Environments:    environmentsCount,
		Segments:        segmentsCount,
		Experiments:     experimentsCount,
		APIKeys:         apiKeysCount,
		AuditLogEntries: auditLogCount,
	}

	return c.JSON(response)
}
