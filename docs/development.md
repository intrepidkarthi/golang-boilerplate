# Development Guide

## Development Environment Setup

### Required Tools
1. **Go 1.21+**
   ```bash
   # Install Go
   brew install go
   
   # Verify installation
   go version
   ```

2. **Docker & Docker Compose**
   ```bash
   # Install Docker Desktop (includes Docker Compose)
   brew install --cask docker
   ```

3. **Protocol Buffers**
   ```bash
   # Install protoc compiler
   brew install protobuf
   
   # Install Go plugins
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

4. **Development Tools**
   ```bash
   # Install air for hot reload
   go install github.com/cosmtrek/air@latest
   
   # Install golang-migrate
   brew install golang-migrate

   # Install sqlc
   make sqlc-install
   ```

## Project Setup

1. **Clone Repository**
   ```bash
   git clone <repository-url>
   cd go-boilerplate
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod tidy

   # Generate SQL code
   make sqlc-generate
   ```

3. **Environment Setup**
   ```bash
   # Copy example environment file
   cp .env.example .env
   
   # Edit environment variables
   vim .env
   ```

4. **Start Infrastructure**
   ```bash
   # Start required services
   make docker-up
   
   # Run migrations
   make migrate-up
   
   # Generate protobuf
   make proto
   ```

## Development Workflow

### Running the Application

1. **Development Mode (with hot reload)**
   ```bash
   make run
   ```

2. **Build and Run**
   ```bash
   make build
   ./bin/server
   ```

3. **Docker Development**
   ```bash
   # Build image
   make docker-build
   
   # Run container
   make docker-run
   ```

### Making Changes

1. **Code Style**
   - Follow Go standard formatting
   - Use gofmt/goimports
   - Follow project structure

2. **Adding Dependencies**
   ```bash
   # Add new dependency
   go get github.com/example/package
   
   # Update dependencies
   go mod tidy
   ```

3. **Database Changes**
   ```bash
   # Create new migration
   make migrate-create name=add_new_table
   
   # Apply migration
   make migrate-up
   
   # Rollback migration
   make migrate-down
   ```

4. **Protocol Buffer Changes**
   ```bash
   # Edit proto files in /proto directory
   # Generate new code
   make proto
   ```

### Testing

1. **Running Tests**
   ```bash
   # Run all tests
   make test
   
   # Run with coverage
   make test-coverage
   
   # Run specific test
   go test ./internal/service -run TestCreateMessage
   ```

2. **Writing Tests**
   - Place tests in same package as code
   - Use table-driven tests
   - Mock external dependencies
   - Use testify assertions

### Debugging

1. **Logs**
   ```bash
   # View application logs
   make logs
   
   # View specific service logs
   make logs-db
   make logs-redis
   make logs-kafka
   ```

2. **Database**
   ```bash
   # Connect to database
   make db-connect
   
   # View migrations
   make migrate-status
   ```

3. **Profiling**
   ```bash
   # Enable profiling
   export ENABLE_PPROF=true
   
   # Access profiles at
   http://localhost:8080/debug/pprof/
   ```

## Code Organization

### Project Structure
```
.
├── api/                 # API layer
│   ├── grpc/           # gRPC implementations
│   └── http/           # HTTP handlers
├── cmd/                # Entry points
├── config/             # Configuration
├── internal/           # Internal packages
├── migrations/         # Database migrations
├── proto/              # Protocol buffers
└── scripts/           # Utility scripts
```

### Package Guidelines
1. **Keep packages focused**
   - Single responsibility
   - Clear dependencies
   - Well-defined interfaces

2. **Use internal packages**
   - Hide implementation details
   - Control API surface
   - Maintain modularity

3. **Organize by feature**
   - Group related functionality
   - Minimize cross-package dependencies
   - Clear ownership

## Best Practices

### Code Style
1. **Formatting**
   - Use gofmt
   - Follow Go conventions
   - Consistent naming

2. **Documentation**
   - Document public APIs
   - Add examples
   - Keep docs updated

3. **Error Handling**
   - Use custom errors
   - Proper error wrapping
   - Meaningful messages

### Testing
1. **Unit Tests**
   - Table-driven tests
   - Mock interfaces
   - Clear assertions

2. **Integration Tests**
   - Test real dependencies
   - Clean test data
   - Proper setup/teardown

3. **Performance Tests**
   - Benchmark critical paths
   - Profile bottlenecks
   - Load testing

### Security
1. **Input Validation**
   - Validate all input
   - Sanitize data
   - Use prepared statements

2. **Authentication**
   - Secure tokens
   - Proper encryption
   - Rate limiting

3. **Secrets Management**
   - Use environment variables
   - Encrypt sensitive data
   - Rotate credentials

## Troubleshooting

### Common Issues

1. **Database Connection**
   ```bash
   # Check database status
   make db-status
   
   # Reset database
   make db-reset
   ```

2. **Protocol Buffer Generation**
   ```bash
   # Clean generated files
   make proto-clean
   
   # Regenerate
   make proto
   ```

3. **Docker Issues**
   ```bash
   # Reset containers
   make docker-down
   make docker-up
   
   # Clean volumes
   make docker-clean
   ```

### Getting Help
1. Check documentation
2. Review issue tracker
3. Ask in team chat
4. Create new issue

## Deployment

### Building
```bash
# Build binary
make build

# Build Docker image
make docker-build
```

### Configuration
1. Set environment variables
2. Update configs
3. Check dependencies

### Verification
1. Run tests
2. Check logs
3. Monitor metrics

## Additional Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Project Wiki](docs/wiki)
- [API Documentation](docs/api.md)

### Tools
- [Go Tools](https://golang.org/cmd/go/)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
