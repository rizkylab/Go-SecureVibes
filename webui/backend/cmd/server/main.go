package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/rizkylab/Go-SecureVibes/webui/backend/config"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/handlers"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "SecureVibes Dashboard v1.0",
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Server.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/login", handlers.Login(db, cfg))

	// Protected routes
	protected := api.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))

	// Dashboard
	protected.Get("/dashboard/summary", handlers.GetDashboardSummary(db))

	// Scans
	protected.Get("/scans", handlers.ListScans(db))
	protected.Get("/scans/:id", handlers.GetScan(db))
	protected.Post("/scan-result", handlers.UploadScanResult(db, cfg))
	protected.Delete("/scans/:id", handlers.DeleteScan(db, cfg))

	// Architecture
	protected.Get("/scans/:id/architecture", handlers.GetArchitecture(db, cfg))

	// Threat Model
	protected.Get("/scans/:id/threat-model", handlers.GetThreatModel(db, cfg))

	// Findings
	protected.Get("/scans/:id/findings", handlers.ListFindings(db))
	protected.Get("/scans/:id/findings/:finding_id", handlers.GetFinding(db))
	protected.Patch("/scans/:id/findings/:finding_id/status", handlers.UpdateFindingStatus(db))

	// DAST
	protected.Get("/scans/:id/dast-findings", handlers.GetDASTFindings(db, cfg))

	// Comparison
	protected.Get("/scans/compare", handlers.CompareScans(db, cfg))

	// Export
	protected.Get("/scans/:id/export", handlers.ExportScan(db, cfg))

	// API Tokens
	protected.Get("/tokens", handlers.ListAPITokens(db))
	protected.Post("/tokens", handlers.CreateAPIToken(db))
	protected.Delete("/tokens/:id", handlers.RevokeAPIToken(db))

	// Static files (for frontend)
	app.Static("/", "./public")

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
		})
	})

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server stopped")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
