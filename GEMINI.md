# Project: tests

A Go-based REST API for managing Organizations and Contracts.

## Project Overview

This project is a Go application structured with a clean separation of concerns, utilizing the Repository pattern. It serves as a demonstration for why in-memory applications are beneficial for service testing.

### Key Components

- **DTOs & Validation**: Located in `internal/dtos`, using `github.com/go-playground/validator/v10`.
- **Services**: Business logic in `internal/services`.
- **Handlers**: Fiber-based handlers in `internal/handlers` (`organizations_handler.go`, `contracts_handler.go`).
- **Repositories**: Data access layer in `internal/repositories`.

## Building and Running

### Prerequisites

- Go 1.25.6 or later.
- A running PostgreSQL instance (for the main application).
- A `.env` file in the root directory.

### Commands

- **Install dependencies:**
  ```bash
  go mod tidy
  ```
- **Run the application:**
  ```bash
  go run cmd/api/main.go
  ```

## In-Memory Testing Demo

The project is designed to be easily testable using in-memory implementations of repositories and even an in-memory application instance for integration testing without a real database.

## TODO / Known Issues

- Implement unit tests for services using mocked repositories.
- Implement integration tests using an in-memory database or application setup.
