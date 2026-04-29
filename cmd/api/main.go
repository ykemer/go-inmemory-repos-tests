package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tests/internal/config"
	"tests/internal/handlers"
	"tests/internal/idempotency"
	"tests/internal/middleware"
	"tests/internal/repositories"
	"tests/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()
	db := cfg.InitDatabase()

	// Repositories
	orgRepo := repositories.NewOrgsRepository(db)
	contractRepo := repositories.NewContractsRepository(db)

	// Services
	orgService := services.NewOrganizationService(orgRepo)
	contractService := services.NewContractService(contractRepo, orgRepo)

	// Handlers
	orgHandler := handlers.NewOrganizationHandler(orgService)
	contractHandler := handlers.NewContractHandler(contractService)

	idempotencyStore := idempotency.NewInMemoryStore(1 * time.Hour)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ProblemJSONErrorHandler,
	})
	app.Use(middleware.Recover())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded, try again later")
		},
	}))
	app.Use(middleware.Idempotency(idempotencyStore))

	api := app.Group("/api")

	// Organizations
	orgs := api.Group("/organizations")
	orgs.Post("/", orgHandler.Create)
	orgs.Get("/", orgHandler.List)
	orgs.Get("/:id", orgHandler.Get)
	orgs.Put("/:id", orgHandler.Update)
	orgs.Delete("/:id", orgHandler.Delete)

	// Contracts
	orgs.Post("/:orgId/contracts", contractHandler.Create)
	orgs.Get("/:orgId/contracts", contractHandler.List)
	orgs.Get("/:orgId/contracts/:contractId", contractHandler.Get)
	orgs.Put("/:orgId/contracts/:contractId", contractHandler.Update)
	orgs.Delete("/:orgId/contracts/:contractId", contractHandler.Delete)

	// Start the server in a goroutine so the main goroutine can wait for a
	// shutdown signal.
	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received, draining connections (timeout: 10s)...")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server stopped cleanly")
}
