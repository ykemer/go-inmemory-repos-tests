# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make run          # Run the API server
make test         # Run all tests with verbose output
make test-no-cache  # Run all tests bypassing cache
make build        # Build binary to bin/api
make clean        # Remove build artifacts
```

Run a single test or test function:
```bash
go test ./internal/services/... -run TestContractService_Create -v
go test ./internal/services/... -run TestContractService_Create/Success -v
```

## Architecture

This is a Fiber-based REST API for managing **Organizations** and **Contracts** (contracts are nested under organizations). The app is intentionally structured to demonstrate service-layer testing without a real database.

### Dependency flow

```
main.go → handlers → services → repositories (interface)
                                      ↓
                            GORM impl (prod) or InMemory impl (tests)
```

### Key design: In-Memory Repository Pattern

The core of this project is how tests are structured. Repository interfaces (`IOrgsRepository`, `IContractsRepository`) in `internal/repositories/interfaces.go` allow services to be tested without a real PostgreSQL instance.

Two implementations exist for each interface:
- **GORM repos** (`orgs_repository.go`, `contracts_repository.go`) — used in production via `main.go`
- **InMemory repos** (`inmemory_orgs_repository.go`, `inmemory_contracts_repository.go`) — used in service tests

### Hook system (testing unhappy paths)

The InMemory repos expose programmable hooks (`OnCreate`, `OnList`, `OnGetByID`, etc.) that are `nil` by default (stateful mode). Assign a function to force a specific return value in a single test:

```go
contractRepo.OnCreate = func(ctx context.Context, data *models.Contract) error {
    return fmt.Errorf("db error")
}
```

Call `repo.Clear()` between tests to reset both stored data and all hooks.

### Route structure

```
POST   /api/organizations
GET    /api/organizations
GET    /api/organizations/:id
PUT    /api/organizations/:id
DELETE /api/organizations/:id

POST   /api/organizations/:orgId/contracts
GET    /api/organizations/:orgId/contracts
GET    /api/organizations/:orgId/contracts/:contractId
PUT    /api/organizations/:orgId/contracts/:contractId
DELETE /api/organizations/:orgId/contracts/:contractId
```

### Configuration

Config is loaded from `.env` via `internal/config/config.go` using `github.com/caarlos0/env`. Copy `.env.example` to `.env` before running the server (PostgreSQL required for the real app; not needed for tests).

### DTOs & validation

Request/response shapes live in `internal/dtos`. Input validation uses `github.com/go-playground/validator/v10` struct tags. Services return DTO response types directly, not raw models.

### Adding a new service test

1. Instantiate fresh InMemory repos and the service under test at the start of each subtest (no shared state between subtests).
2. Use hook overrides only for error-path subtests; let the stateful InMemory impl handle happy-path data.
3. Assert with `github.com/stretchr/testify/assert`.
