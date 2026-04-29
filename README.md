# Go Service Testing Demo

This project is a high-performance demonstration of how to build and test Go services using the **In-Memory Repository Pattern**.

## The Architecture: Why it's Fantastic

By decoupling business logic from the database layer using interfaces, we gain an incredible advantage:
- **Lightning-Fast Tests**: Run hundreds of business logic tests in milliseconds.
- **Cheap Infrastructure**: No need to spin up a Dockerized database or manage migrations just to verify your use cases.
- **Total Isolation**: Test edge cases, repository failures, and complex logic without side effects.
- **Developer Productivity**: Immediate feedback loop during development.

---

## Deep Dive: How In-Memory Repos Work

You might notice the In-Memory repository files contain a significant amount of code. This is intentional and necessary for a robust, production-grade testing environment.

### 1. Full Interface Implementation
In Go, an interface is a contract. To use the `InMemoryRepository` in place of the real GORM repository, it must implement **every single method** of the interface (List, Get, Create, Update, Delete). This ensures our tests are "real" and catch signature mismatches.

### 2. Thread Safety (Concurrency)
Unlike a simple script, a web API handles many requests at once. Our In-Memory repos use `sync.RWMutex` to prevent data races. This allows you to run integration tests where multiple Goroutines might access the storage simultaneously.

### 3. Programmable Behavior (The Hooks)
The "secret sauce" of these repos is the **Hook System** (e.g., `OnCreate`, `OnList`).
- **Standard Mode**: By default, the repo uses internal maps to store data (Stateful).
- **Mock Mode**: You can "program" a method to fail for a specific test:
  ```go
  repo.OnCreate = func(...) error { return errors.New("DB connection lost") }
  ```
This is why we have 100% test coverage—we can force the "unhappy paths" that are impossible to trigger with a healthy database.

---

## Getting Started

### Prerequisites
- Go 1.25.6+

### Setup
1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## How to Run & Test

### Run the API
```bash
make run
```

### Run the Tests
```bash
make test
```
This command runs all service-level tests with **100% coverage**, proving that all business logic is correct without ever touching a real database.

### Manual Testing
A pre-configured `api_tests.http` file is included for use with any REST client.
