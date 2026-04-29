package main

import (
	"log"
	"tests/internal/config"
	"tests/internal/handlers"
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

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ProblemJSONErrorHandler,
	})
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded, try again later")
		},
	}))

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

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
