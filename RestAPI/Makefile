.PHONY: run swagger test build clean tidy

# Run the application
run:
	go run ./cmd/main.go

# Generate Swagger documentation
swagger:
	swag init -g cmd/main.go -o docs

# Run tests with coverage
test:
	go test ./... -v -cover

# Build binary
build:
	go build -o bin/restapi ./cmd/main.go

# Clean build artifacts and database
clean:
	rm -rf bin/ data/

# Tidy dependencies
tidy:
	go mod tidy

# Install swag CLI if not present
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Full setup: install tools, tidy, run
setup: install-swag tidy swagger run
