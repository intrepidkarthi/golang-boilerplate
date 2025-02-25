# Go Microservices Boilerplate

A production-grade Go microservices boilerplate with support for gRPC, REST, Kafka, Redis, and PostgreSQL. This boilerplate provides a solid foundation for building scalable, maintainable, and well-structured microservices.

## 🚀 Features

### Core Features
- **Dual Protocol Support**: REST API and gRPC services
- **Database Integration**: PostgreSQL with sqlc for type-safe SQL
- **Caching**: Redis for improved performance
- **Message Streaming**: Kafka for event-driven architecture
- **API Documentation**: Swagger/OpenAPI documentation
- **Health Monitoring**: Comprehensive health check endpoints

### Technical Features
- **Validation**: Request validation using go-playground/validator
- **Error Handling**: Centralized error handling with middleware
- **Pagination**: Built-in support for paginated responses
- **Testing**: Comprehensive unit and integration tests
- **Docker**: Containerization with Docker and Docker Compose
- **Graceful Shutdown**: Proper shutdown handling
- **Structured Logging**: Using Zap logger
- **Type-safe SQL**: Using sqlc for compile-time SQL validation
- **Security Scanning**: Automated security checks with gosec and golangci-lint

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
│   ├── db/             # Database operations and sqlc generated code
│   ├── kafka/          # Kafka producer/consumer
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   └── service/        # Business logic
├── migrations/         # Database migrations
├── proto/              # Protocol buffer definitions
└── scripts/           # Utility scripts
```

## 🔒 Security and Code Quality

The project includes comprehensive security scanning and code quality tools:

### Security Scanning with Gosec

Gosec is specifically designed to scan Go code for security vulnerabilities:

```bash
# Install Gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run Gosec
~/go/bin/gosec -quiet ./...
```

Key security checks include:
- Buffer overflow vulnerabilities
- Integer overflow risks
- SQL injection vulnerabilities
- Command injection risks
- Cryptographic weakness
- Hardcoded credentials
- Insecure file operations

### Code Quality with GolangCI-Lint

GolangCI-Lint provides comprehensive static code analysis:

```bash
# Install GolangCI-Lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Run GolangCI-Lint
~/go/bin/golangci-lint run --timeout=5m
```

Lint checks include:
- Code style and formatting
- Potential bugs and errors
- Performance issues
- Code complexity
- Security anti-patterns
- Best practices violations

### Continuous Integration

Both security scanning and linting are integrated into the CI pipeline:

```bash
# Run all checks
make check

# Run security checks only
make security

# Run lint checks only
make lint
```

### Configuration

- Security rules: `.gosec.config.json`
- Linting rules: `.golangci.yml`
- CI configuration: `.github/workflows/security.yml`

Customize these files to adjust severity levels, exclude patterns, or add custom rules.

## 🚀 Running the Application

### Prerequisites

1. **Go 1.21 or later**
2. **PostgreSQL**
3. **Redis**
4. **Kafka**

### Local Development

1. **Build the application**:
```bash
go build -o bin/server ./cmd/server
```

2. **Set up environment variables**:
Create a `.env` file in the root directory with the following content:
```env
# Server Configuration
PORT=8080
READ_TIMEOUT=5s
WRITE_TIMEOUT=10s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=1h

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=messages

# gRPC Configuration
GRPC_PORT=50051
```

3. **Start required services** (if using Docker):
```bash
docker-compose up -d postgres redis kafka
```

4. **Run the application**:
```bash
./bin/server
```

### Available Endpoints

- **HTTP API**: `http://localhost:8080`
  - Health Check: `GET /health`
  - Messages API: `GET /api/v1/messages`

- **gRPC**: `localhost:50051`

### Monitoring

- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`

## 🏥 Health Monitoring

The service includes comprehensive health check endpoints for monitoring and operational readiness:

### Health Check Endpoints

1. **Main Health Check**
   ```bash
   curl http://localhost:3000/health
   ```
   Returns detailed health status of the service and its dependencies:
   ```json
   {
       "status": "ok",
       "timestamp": "2025-02-25T16:55:21+05:30",
       "version": "1.0.0",
       "services": {
           "database": {
               "status": "up"
           },
           "cache": {
               "status": "up"
           }
       }
   }
   ```

2. **Kubernetes Liveness Probe**
   ```bash
   curl http://localhost:3000/health/live
   ```
   Quick check to verify if the service is running:
   ```json
   {
       "status": "alive"
   }
   ```

3. **Kubernetes Readiness Probe**
   ```bash
   curl http://localhost:3000/health/ready
   ```
   Verifies if the service is ready to handle requests:
   ```json
   {
       "status": "ready"
   }
   ```

### Health Check Features
- Real-time dependency status monitoring
- Kubernetes-compatible health probes
- Detailed service status reporting
- Timestamp and version information
- Individual component health status

### Integration with Monitoring
The health endpoints can be integrated with:
- Kubernetes health checks
- Load balancer health monitoring
- Prometheus metrics collection
- Custom monitoring solutions
- CI/CD pipeline checks

## 📚 API Documentation with Swagger

### Installing Swagger Tools

1. **Install Swag CLI**
   ```bash
   # Install swag CLI tool
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Verify installation
   ~/go/bin/swag --version
   ```

2. **Install Echo Swagger**
   ```bash
   # Install echo-swagger package
   go get -u github.com/swaggo/echo-swagger
   ```

### Generating API Documentation

1. **Add Swagger Annotations**
   Add Swagger annotations to your handlers. Example:
   ```go
   // @Summary Create a new message
   // @Description Create a new message with the provided content
   // @Tags messages
   // @Accept json
   // @Produce json
   // @Param message body CreateMessageRequest true "Message content"
   // @Success 200 {object} models.Message
   // @Router /messages [post]
   ```

2. **Generate Swagger Files**
   ```bash
   # Generate swagger documentation
   ~/go/bin/swag init -g cmd/server/main.go
   ```

3. **View Documentation**
   - Start the server: `go run cmd/server/main.go`
   - Open your browser and navigate to: `http://localhost:3000/swagger/index.html`

### Updating Documentation

Whenever you make changes to your API annotations:
1. Regenerate the documentation: `~/go/bin/swag init -g cmd/server/main.go`
2. Restart the server to see the changes

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

3. **Infrastructure Services**
   The following services are required and will be automatically started via Docker Compose:
   - PostgreSQL 15+ (Database)
   - Redis 7+ (Caching)
   - Apache Kafka 3+ (Message Streaming)
   - Zookeeper (Required for Kafka)

   **Manual Installation (macOS)**:
   ```bash
   # Install PostgreSQL
   brew install postgresql@15
   brew services start postgresql@15

   # Install Redis
   brew install redis
   brew services start redis

   # Install Kafka (includes Zookeeper)
   brew install kafka
   brew services start zookeeper
   brew services start kafka
   ```

   **Default Configurations**:
   - PostgreSQL:
     ```
     Host: localhost
     Port: 5432
     Default Database: postgres
     Default User: postgres
     Default Password: postgres
     ```

   - Redis:
     ```
     Host: localhost
     Port: 6379
     No default password
     ```

   - Kafka:
     ```
     Bootstrap Servers: localhost:9092
     Zookeeper: localhost:2181
     Default Topics Created: messages
     Replication Factor: 1
     Number of Partitions: 3
     ```

   **Configuration Files**:
   - PostgreSQL: ~/Library/Application Support/postgres/postgresql.conf
   - Redis: /usr/local/etc/redis.conf
   - Kafka: /usr/local/etc/kafka/server.properties
   - Zookeeper: /usr/local/etc/kafka/zookeeper.properties

   **Common Commands**:
   ```bash
   # PostgreSQL
   psql -U postgres                    # Connect to PostgreSQL
   createdb mydb                       # Create a database
   dropdb mydb                         # Drop a database

   # Redis
   redis-cli                          # Connect to Redis
   redis-cli ping                      # Test Redis connection
   redis-cli monitor                   # Monitor Redis commands

   # Kafka
   kafka-topics --create \
     --bootstrap-server localhost:9092 \
     --topic my-topic \
     --partitions 3 \
     --replication-factor 1           # Create a Kafka topic

   kafka-topics --list \
     --bootstrap-server localhost:9092 # List topics

   kafka-console-producer \
     --bootstrap-server localhost:9092 \
     --topic my-topic                 # Produce messages

   kafka-console-consumer \
     --bootstrap-server localhost:9092 \
     --topic my-topic \
     --from-beginning                 # Consume messages
   ```

   **Health Checks**:
   ```bash
   # PostgreSQL
   pg_isready -h localhost -p 5432

   # Redis
   redis-cli ping

   # Kafka
   kafka-topics --bootstrap-server localhost:9092 --list
   ```

4. **Protocol Buffers Compiler**
   ```bash
   # macOS
   brew install protobuf
   
   # Check protoc version
   protoc --version
   ```

5. **Go Protocol Buffers plugins**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

6. **sqlc**
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
curl -X POST http://localhost:3000/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello, World!"}'

# Get a message
curl http://localhost:3000/api/v1/messages/{id}

# Update a message
curl -X PUT http://localhost:3000/api/v1/messages/{id} \
  -H "Content-Type: application/json" \
  -d '{"content":"Updated content"}'

# Delete a message
curl -X DELETE http://localhost:3000/api/v1/messages/{id}

# List messages (with pagination)
curl http://localhost:3000/api/v1/messages?page=1&page_size=10
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
PORT=3000
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

## 🛠 Customizing the Services

### REST API Development

The REST API is built using the Echo framework. Here's how to add new endpoints:

1. **Create a New Handler**
   ```go
   // internal/api/http/your_handler.go
   type YourHandler struct {
       service *service.YourService
   }

   func NewYourHandler(service *service.YourService) *YourHandler {
       return &YourHandler{service: service}
   }

   func (h *YourHandler) CreateItem(c echo.Context) error {
       // Bind and validate request
       req := new(CreateItemRequest)
       if err := c.Bind(req); err != nil {
           return echo.NewHTTPError(http.StatusBadRequest, err.Error())
       }

       // Your handler logic
       return c.JSON(http.StatusCreated, response)
   }
   ```

2. **Add Routes**
   ```go
   // cmd/server/main.go in the router setup
   e := echo.New()

   // Middleware
   e.Use(middleware.Logger())
   e.Use(middleware.Recover())
   e.Use(middleware.CORS())

   // Routes
   v1 := e.Group("/api/v1")
   items := v1.Group("/items")
   items.POST("", yourHandler.CreateItem)
   items.GET("", yourHandler.ListItems)
   items.GET("/:id", yourHandler.GetItem)
   ```

3. **Key Files to Modify**:
   - `internal/api/http/` - Add new handlers
   - `internal/service/` - Implement business logic
   - `internal/db/query.sql` - Add new SQL queries
   - `internal/models/` - Define request/response structs

### gRPC Service Development

1. **Define New Service**
   ```protobuf
   // proto/your_service/v1/your_service.proto
   service YourService {
       rpc CreateItem(CreateItemRequest) returns (CreateItemResponse);
       rpc GetItem(GetItemRequest) returns (GetItemResponse);
       rpc ListItems(ListItemsRequest) returns (stream ListItemsResponse);
   }
   ```

2. **Generate gRPC Code**
   ```bash
   make proto
   ```

3. **Implement Service**
   ```go
   // internal/api/grpc/your_service.go
   type YourServiceServer struct {
       pb.UnimplementedYourServiceServer
       service *service.YourService
   }

   func (s *YourServiceServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
       // Your implementation
   }
   ```

4. **Register Service**
   ```go
   // cmd/server/main.go in the gRPC server setup
   pb.RegisterYourServiceServer(grpcServer, grpc.NewYourServiceServer(yourService))
   ```

5. **Key Files to Modify**:
   - `proto/` - Define new services and messages
   - `internal/api/grpc/` - Implement gRPC services
   - `internal/service/` - Add business logic
   - `internal/db/query.sql` - Add required SQL queries

### Best Practices

1. **Validation**
   - Use `validator` tags for REST API requests
   - Implement validation in gRPC services using interceptors

2. **Error Handling**
   - Use the provided error types in `internal/errors`
   - Map domain errors to appropriate HTTP/gRPC status codes

3. **Database Operations**
   - Add new queries to `internal/db/query.sql`
   - Run `make sqlc` to generate type-safe Go code

4. **Testing**
   - Add unit tests for new handlers and services
   - Use the provided test helpers in `internal/testutil`

## 📚 Additional Documentation

- [API Documentation](docs/api.md)
- [Database Schema](docs/schema.md)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Echo Framework](https://github.com/labstack/echo)
- [sqlc](https://sqlc.dev)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Kafka](https://kafka.apache.org)
- [Redis](https://redis.io)


n created.

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
