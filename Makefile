.PHONY: build run test proto clean docker-up docker-down migrate-up migrate-down sqlc-install sqlc-generate sqlc-verify check security lint install-tools

# Code quality and security commands
check: security lint

security:
	@echo "Running security checks..."
	@~/go/bin/gosec -quiet ./...

lint:
	@echo "Running linter..."
	@~/go/bin/golangci-lint run --timeout=5m

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Build commands
build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

# Proto commands
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/message/v1/*.proto

# Test commands
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...

# Migration commands
migrate-up:
	go run cmd/migrate/main.go -direction up

migrate-down:
	go run cmd/migrate/main.go -direction down

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development setup
dev-setup: docker-down docker-build docker-up migrate-up

# Cleanup
clean:
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
deps:
	go mod download

# SQLC commands
sqlc-install:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc-generate:
	sqlc generate

sqlc-verify:
	sqlc verify

.DEFAULT_GOAL := build
