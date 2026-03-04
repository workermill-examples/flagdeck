package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/workermill-examples/flagdeck/api/internal/config"
	"github.com/workermill-examples/flagdeck/api/internal/database"
	"github.com/workermill-examples/flagdeck/api/internal/handlers"
	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/services"
)

// SetupRoutes configures all application routes with appropriate middleware
func SetupRoutes(app *fiber.App, mongodb *database.MongoDB, redisdb *database.RedisDB, cfg *config.Config) {
	// Initialize services
	auditService := services.NewAuditService(mongodb.AuditLogCollection())
	evaluationService := services.NewEvaluationService(
		mongodb.FlagsCollection(),
		mongodb.EnvironmentsCollection(),
		mongodb.SegmentsCollection(),
	)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(mongodb, redisdb)
	authHandler := handlers.NewAuthHandler(
		mongodb.UsersCollection(),
		cfg.JWTSecret,
	)
	flagsHandler := handlers.NewFlagsHandler(
		mongodb.FlagsCollection(),
		auditService,
	)
	evaluateHandler := handlers.NewEvaluateHandler(evaluationService)
	environmentsHandler := handlers.NewEnvironmentsHandler(
		mongodb.EnvironmentsCollection(),
		auditService,
	)
	segmentsHandler := handlers.NewSegmentsHandler(
		mongodb.SegmentsCollection(),
		auditService,
	)
	experimentsHandler := handlers.NewExperimentsHandler(
		mongodb.ExperimentsCollection(),
		auditService,
	)
	apiKeysHandler := handlers.NewApiKeysHandler(
		mongodb.APIKeysCollection(),
		auditService,
	)
	auditHandler := handlers.NewAuditHandler(mongodb.AuditLogCollection())
	statsHandler := handlers.NewStatsHandler(
		mongodb.FlagsCollection(),
		mongodb.EnvironmentsCollection(),
		mongodb.SegmentsCollection(),
		mongodb.ExperimentsCollection(),
		mongodb.APIKeysCollection(),
		mongodb.AuditLogCollection(),
	)

	// Initialize middleware
	jwtAuth := middleware.AuthenticateJWT(middleware.AuthConfig{
		JWTSecret:      cfg.JWTSecret,
		UserCollection: mongodb.UsersCollection(),
	})
	apiKeyAuth := middleware.AuthenticateAPIKey(middleware.APIKeyConfig{
		APIKeyCollection: mongodb.APIKeysCollection(),
	})
	authRateLimit := middleware.AuthRateLimit(redisdb.Client)     // 300 req/min for auth
	apiRateLimit := middleware.APIRateLimit(redisdb.Client)       // 1000 req/min for API
	evalRateLimit := middleware.EvaluateRateLimit(redisdb.Client) // 1000 req/min for evaluation

	// Public routes
	app.Get("/health", healthHandler.GetHealth)

	// Authentication routes (public but rate-limited)
	authGroup := app.Group("/auth", authRateLimit)
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.Refresh)
	authGroup.Post("/logout", authHandler.Logout)

	// Protected authentication routes (require JWT)
	authGroup.Get("/me", jwtAuth, authHandler.GetMe)

	// Protected API routes (require JWT authentication)
	apiGroup := app.Group("/api/v1", jwtAuth, apiRateLimit)

	// Flag management
	apiGroup.Get("/flags", flagsHandler.ListFlags)
	apiGroup.Get("/flags/:key", flagsHandler.GetFlag)
	apiGroup.Post("/flags", flagsHandler.CreateFlag)
	apiGroup.Put("/flags/:key", flagsHandler.UpdateFlag)
	apiGroup.Delete("/flags/:key", flagsHandler.DeleteFlag)
	apiGroup.Post("/flags/:key/toggle", flagsHandler.ToggleFlag)

	// Environment management
	apiGroup.Get("/environments", environmentsHandler.ListEnvironments)
	apiGroup.Get("/environments/:key", environmentsHandler.GetEnvironment)
	apiGroup.Post("/environments", environmentsHandler.CreateEnvironment)
	apiGroup.Put("/environments/:key", environmentsHandler.UpdateEnvironment)
	apiGroup.Delete("/environments/:key", environmentsHandler.DeleteEnvironment)

	// Segment management
	apiGroup.Get("/segments", segmentsHandler.ListSegments)
	apiGroup.Get("/segments/:key", segmentsHandler.GetSegment)
	apiGroup.Post("/segments", segmentsHandler.CreateSegment)
	apiGroup.Put("/segments/:key", segmentsHandler.UpdateSegment)
	apiGroup.Delete("/segments/:key", segmentsHandler.DeleteSegment)

	// Experiment management
	apiGroup.Get("/experiments", experimentsHandler.ListExperiments)
	apiGroup.Get("/experiments/:key", experimentsHandler.GetExperiment)
	apiGroup.Post("/experiments", experimentsHandler.CreateExperiment)
	apiGroup.Put("/experiments/:key", experimentsHandler.UpdateExperiment)
	apiGroup.Delete("/experiments/:key", experimentsHandler.DeleteExperiment)

	// API key management
	apiGroup.Get("/api-keys", apiKeysHandler.ListApiKeys)
	apiGroup.Post("/api-keys", apiKeysHandler.CreateApiKey)
	apiGroup.Delete("/api-keys/:id", apiKeysHandler.DeleteApiKey)

	// Audit log
	apiGroup.Get("/audit-log", auditHandler.GetAuditLog)

	// Stats endpoint
	apiGroup.Get("/stats", statsHandler.GetStats)

	// API key authenticated routes (for flag evaluation and experiment tracking)
	// These routes use API key authentication instead of JWT
	app.Post("/api/v1/evaluate", apiKeyAuth, evalRateLimit, evaluateHandler.EvaluateFlag)
	app.Post("/api/v1/evaluate/bulk", apiKeyAuth, evalRateLimit, evaluateHandler.EvaluateBulk)
	app.Post("/api/v1/experiments/:key/track", apiKeyAuth, evalRateLimit, experimentsHandler.TrackExperiment)
}
