package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/workermill-examples/flagdeck/api/internal/config"
	"github.com/workermill-examples/flagdeck/api/internal/database"
	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/routes"
)

func main() {
	cfg := config.Load()

	mongodb, err := database.NewMongoDB(cfg.MongodbURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close()

	redisdb, err := database.NewRedisDB(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisdb.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    10 * 1024 * 1024, // 10MB
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${pid} ${locals:requestid} ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(recover.New())
	app.Use(middleware.DefaultCORS())

	// Setup all application routes
	routes.SetupRoutes(app, mongodb, redisdb, cfg)

	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
