.PHONY: run test build clean docs

# Run the API application
run:
	go run cmd/api/main.go

# Run all tests with verbose output
test:
	go test ./... -v

test-no-cache:
	go test ./... -v -count=1

# Build the project
build:
	go build -o bin/api cmd/api/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Regenerate OpenAPI spec from handler annotations.
# Requires: go install github.com/swaggo/swag/cmd/swag@latest
docs:
	swag init -g cmd/api/main.go -o internal/docs --parseDependency --parseInternal
