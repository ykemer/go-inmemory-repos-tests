.PHONY: run test build clean docs mcp-run mcp-build

# Regenerate OpenAPI spec from handler annotations (JSON only — no Go blob).
# Requires: go install github.com/swaggo/swag/cmd/swag@latest
docs:
	swag init -g cmd/api/main.go -o internal/docs --outputTypes json --parseDependency --parseInternal

# Run the API application (regenerates spec first)
run: docs
	go run cmd/api/main.go

# Build binary (regenerates spec first so the embed is always up to date)
build: docs
	go build -o bin/api cmd/api/main.go

# Run all tests with verbose output
test:
	go test ./... -v

test-no-cache:
	go test ./... -v -count=1

# Clean build artifacts
clean:
	rm -rf bin/

# Run the MCP server (stdio transport — connect via an MCP client or Claude Desktop)
mcp-run:
	go run cmd/mcp/main.go

# Build the MCP server binary
mcp-build:
	go build -o bin/mcp cmd/mcp/main.go
