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

---

## API Reference (Scalar)

Once the server is running, open **[http://localhost:3000/docs](http://localhost:3000/docs)** in your browser.

[Scalar](https://scalar.com) is a modern API reference UI — a cleaner, interactive alternative to Swagger UI. Browse every endpoint, inspect request/response schemas, and send live requests directly from the browser.

| URL | What it serves |
|---|---|
| `GET /docs` | Scalar interactive UI |
| `GET /docs/openapi.json` | Raw Swagger 2.0 spec (auto-generated) |

### How the spec is generated

The spec is **auto-generated from source code** using [swag](https://github.com/swaggo/swag). Each handler carries `// @` annotations that describe its route, parameters, and response shapes. Running `make docs` regenerates `internal/docs/` from those annotations — so the spec always stays in sync with the code.

```bash
make docs   # regenerate after changing handler annotations
```

One-time setup (install the CLI):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Adding or updating an endpoint

1. Edit or add `// @` annotations on the handler function (see any existing handler for examples).
2. Run `make docs`.
3. Restart the server — the new spec is live at `/docs`.

---

## MCP Server

`cmd/mcp/main.go` exposes all 10 service operations as MCP tools over stdio, so any MCP-compatible client (Claude Desktop, Cursor, etc.) can call them directly.

```bash
make mcp-run    # start the MCP server (stdio)
make mcp-build  # compile to bin/mcp
```

### Inspecting with MCP Inspector

[MCP Inspector](https://github.com/modelcontextprotocol/inspector) is a browser-based UI for exploring and calling tools on any stdio MCP server. No install needed — `npx` pulls it on first run.

**Prerequisites:** Node.js 18+

**Run against the source directly:**
```bash
npx @modelcontextprotocol/inspector go run cmd/mcp/main.go
```

**Or against the compiled binary (faster start):**
```bash
make mcp-build
npx @modelcontextprotocol/inspector ./bin/mcp
```

Inspector opens at **[http://localhost:5173](http://localhost:5173)**. From there you can:

1. Click **Tools** in the left sidebar to see all 10 registered tools.
2. Select a tool (e.g. `create_organization`) — Inspector renders its input schema as a form.
3. Fill in the fields and hit **Run Tool** to send a live call to the server.
4. The raw JSON request and response are shown side by side.

> The MCP server connects to the real PostgreSQL database (same `.env` config as the HTTP server). Make sure the DB is running before launching it.
