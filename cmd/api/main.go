//go:generate swag init -g cmd/api/main.go -o internal/docs --outputTypes json --parseDependency --parseInternal

// @title           Organizations & Contracts API
// @version         1.0.0
// @description     REST API for managing organizations and their contracts.
// @description
// @description     ## Idempotency
// @description     POST and PUT endpoints accept an optional Idempotency-Key header (UUID). Repeated requests with the same key return the cached response for 1 hour, with X-Idempotency-Replayed: true.
// @description
// @description     ## Errors
// @description     All errors are RFC 7807 application/problem+json.
// @host            localhost:3000
// @BasePath        /api

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tests/internal/config"
	"tests/internal/docs"
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

	// Docs
	app.Get("/docs/openapi.json", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/json")
		return c.Send(docs.Spec)
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/html")
		return c.SendString(scalarPage(c.BaseURL()))
	})

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

// scalarPage returns the Scalar API reference HTML, pointing the spec URL at
// baseURL/docs/openapi.json so it works on any host or port.
func scalarPage(baseURL string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width,initial-scale=1"/>
  <title>API Reference</title>
</head>
<body>
  <script
    id="api-reference"
    data-url="%s/docs/openapi.json"
    data-configuration='{"theme":"purple"}'
  ></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`, baseURL)
}
