# Go Microservices Boilerplate

A production-grade Go microservices boilerplate with support for gRPC, REST, Kafka, Redis, and PostgreSQL. This boilerplate provides a solid foundation for building scalable, maintainable, and well-structured microservices.

## 🚀 Features

### Core Features
- **Dual Protocol Support**: REST API and gRPC services
- **Database Integration**: PostgreSQL with sqlc for type-safe SQL
- **Caching**: Redis for improved performance
- **Message Streaming**: Kafka for event-driven architecture
- **API Documentation**: Swagger/OpenAPI documentation

### Technical Features
- **Validation**: Request validation using go-playground/validator
- **Error Handling**: Centralized error handling with middleware
- **Pagination**: Built-in support for paginated responses
- **Testing**: Comprehensive unit and integration tests
- **Docker**: Containerization with Docker and Docker Compose
- **Graceful Shutdown**: Proper shutdown handling
- **Structured Logging**: Using Zap logger

### Development Features
- **Hot Reload**: Live reload during development
- **Make Commands**: Easy-to-use Make commands
- **Migration Tools**: Database migration support
- **Environment Config**: Environment-based configuration

## 📁 Project Structure

```bash
.
├── api/                 # API layer
│   ├── grpc/           # gRPC service implementations
│   └── http/           # HTTP handlers and routes
├── cmd/                # Application entrypoints
│   └── server/         # Main server application
├── config/             # Configuration management
├── internal/           # Internal packages
│   ├── cache/          # Redis cache implementation
│   ├── database/       # Database operations
│   ├── kafka/          # Kafka producer/consumer
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   └── service/        # Business logic
├── migrations/         # Database migrations
├── proto/              # Protocol buffer definitions
└── scripts/           # Utility scripts
```

## 🛠 Prerequisites

Before you begin, ensure you have the following installed:

1. **Go 1.21 or later**
   ```bash
   # Check Go version
   go version
   ```

2. **Docker and Docker Compose**
   ```bash
   # Check Docker version
   docker --version
   docker-compose --version
   ```

3. **Protocol Buffers Compiler**
   ```bash
   # macOS
   brew install protobuf
   
   # Check protoc version
   protoc --version
   ```

4. **Go Protocol Buffers plugins**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

5. **sqlc**
   ```bash
   # Install sqlc
   make sqlc-install
   ```

## ⚡️ Quick Start

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd go-boilerplate
   ```

2. **Set Up Environment Variables**
   ```bash
   # Copy example environment file
   cp .env.example .env
   
   # Edit .env with your configuration
   vim .env
   ```

3. **Start Infrastructure Services**
   ```bash
   # Start PostgreSQL, Redis, and Kafka
   make docker-up
   ```

4. **Run Database Migrations**
   ```bash
   make migrate-up
   ```

5. **Generate Code**
   ```bash
   # Generate Protocol Buffers
   make proto
   
   # Generate SQL code
   make sqlc-generate
   ```

6. **Run the Application**
   ```bash
   # Development mode with hot reload
   make run
   
   # Or build and run
   make build
   ./bin/server
   ```

## 🔧 Development

### Available Make Commands

```bash
# Build the application
make build

# Run with hot reload
make run

# Generate Protocol Buffers
make proto

# Generate SQL code
make sqlc-generate

# Run tests
make test
make test-coverage

# Database operations
make migrate-up
make migrate-down

# Docker operations
make docker-build
make docker-up
make docker-down

# Clean build artifacts
make clean
```

### API Documentation

#### REST Endpoints

```bash
# Create a message
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello, World!"}'}

# Get a message
curl http://localhost:8080/api/v1/messages/{id}

# Update a message
curl -X PUT http://localhost:8080/api/v1/messages/{id} \
  -H "Content-Type: application/json" \
  -d '{"content":"Updated content"}'}

# Delete a message
curl -X DELETE http://localhost:8080/api/v1/messages/{id}

# List messages (with pagination)
curl http://localhost:8080/api/v1/messages?page=1&page_size=10
```

#### gRPC Service

The gRPC service is available at `localhost:50051` with the following methods:
- `CreateMessage`
- `GetMessage`
- `UpdateMessage`
- `DeleteMessage`
- `ListMessages`
- `StreamMessages`

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/service -run TestCreateMessage
```

## 🔒 Security

- All inputs are validated using go-playground/validator
- Proper error handling and sanitization
- Rate limiting middleware available
- Secure headers middleware included
- Environment-based configuration

## 📦 Infrastructure

### PostgreSQL
- **Host**: localhost (default)
- **Port**: 5432
- **Database**: messagedb
- **Migrations**: Located in `/migrations`

### Redis
- **Host**: localhost (default)
- **Port**: 6379
- **Usage**: Message caching

### Kafka
- **Broker**: localhost:9092
- **Topic**: message-events
- **Usage**: Event streaming

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Configuration

The application can be configured using environment variables or a `.env` file:

```env
# Server Configuration
PORT=8080
GRPC_PORT=50051

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=messagedb
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=message-events
```

## 📚 Additional Documentation

- [API Documentation](docs/api.md)
- [Database Schema](docs/schema.md)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Kafka](https://kafka.apache.org)
- [Redis](https://redis.io)


```
.
├── api/                # API Definitions
│   ├── grpc/          # gRPC Services
│   └── http/          # HTTP Handlers
├── cmd/
│   └── server/        # Application Entry Point
├── config/            # Configuration
├── internal/
│   ├── cache/         # Redis Operations
│   ├── database/      # PostgreSQL Operations
│   ├── kafka/         # Kafka Producer/Consumer
│   ├── models/        # Domain Models
│   └── service/       # Business Logic
├── migrations/        # Database Migrations
├── pkg/               # Shared Packages
├── proto/             # Protocol Buffers
└── scripts/          # Utility Scripts
```

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Protocol Buffers Compiler (protoc)
- Make

## Quick Start

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd go-boilerplate
   ```

2. **Set Up Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Infrastructure**
   ```bash
   make dev-setup
   ```

## Development

### Available Make Commands

```bash
# Build and Run
make build        # Build the application
make run          # Run the application

# Protocol Buffers
make proto        # Generate gRPC code

# Testing
make test         # Run tests
make test-coverage # Run tests with coverage

# Docker
make docker-build # Build Docker images
make docker-up    # Start Docker containers
make docker-down  # Stop Docker containers

# Database
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations

# Development
make dev-setup    # Set up development environment
```

## API Examples

### REST Endpoints

```bash
# Create Message
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello, World!"}'}

# Get Message
curl http://localhost:8080/api/v1/messages/{id}

# List Messages
curl http://localhost:8080/api/v1/messages
```

### gRPC Service

The gRPC service is available at `localhost:50051` with the following methods:
- `CreateMessage`
- `GetMessage`
- `StreamMessages`

### Kafka Events

Messages are automatically published to Kafka topic `message-events` when created.

## Infrastructure

### PostgreSQL
- Host: localhost
- Port: 5432
- Database: messagedb

### Redis
- Host: localhost
- Port: 6379

### Kafka
- Broker: localhost:9092
- Topic: message-events

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details
