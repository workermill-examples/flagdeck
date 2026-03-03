package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/workermill-examples/flagdeck/api/internal/config"
	"github.com/workermill-examples/flagdeck/api/internal/database"
	"github.com/workermill-examples/flagdeck/api/internal/handlers"
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

	// Initialize middleware
	// TODO: Implement middleware constructors
	// jwtAuth := middleware.NewJWTAuthMiddleware(mongodb.UsersCollection())
	// apiKeyAuth := middleware.NewApiKeyAuthMiddleware(mongodb.APIKeysCollection())
	// authRateLimit := middleware.NewRateLimitMiddleware(redisdb, "auth", 5, 60)    // 5 req/min for auth
	// apiRateLimit := middleware.NewRateLimitMiddleware(redisdb, "api", 100, 900)   // 100 req/15min for API
	// evalRateLimit := middleware.NewRateLimitMiddleware(redisdb, "eval", 1000, 60) // 1000 req/min for evaluation

	// Public routes
	app.Get("/health", healthHandler.GetHealth)

	// Authentication routes (public but rate-limited)
	authGroup := app.Group("/auth") // TODO: Add authRateLimit
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	// authGroup.Post("/refresh", authHandler.RefreshToken) // TODO: Implement RefreshToken method
	authGroup.Post("/logout", authHandler.Logout)

	// Protected authentication routes (require JWT)
	authGroup.Get("/me", authHandler.GetMe) // TODO: Add jwtAuth middleware

	// Protected API routes (require JWT authentication)
	apiGroup := app.Group("/api/v1") // TODO: Add jwtAuth, apiRateLimit middleware

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

	// API key authenticated routes (for flag evaluation and experiment tracking)
	evalGroup := app.Group("/api/v1") // TODO: Add apiKeyAuth, evalRateLimit middleware
	evalGroup.Post("/evaluate", evaluateHandler.EvaluateFlag)
	evalGroup.Post("/evaluate/bulk", evaluateHandler.EvaluateBulk)
	evalGroup.Post("/experiments/:key/track", experimentsHandler.TrackExperiment)
}
