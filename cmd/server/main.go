// Package main is the entry point for the message service application.
//
// This service implements a modern microservice architecture with support for
// multiple protocols and data stores. It serves as a boilerplate for building
// scalable and maintainable microservices in Go.
//
// Architecture Overview:
// - HTTP Server: RESTful API using Echo framework
// - gRPC Server: High-performance RPC using Protocol Buffers
// - Database: PostgreSQL with sqlc for type-safe queries
// - Cache: Redis for performance optimization
// - Message Queue: Kafka for event-driven architecture
// - Documentation: Swagger/OpenAPI specification
//
// Key Components:
// - Request validation and error handling
// - Middleware for security, logging, and metrics
// - Graceful shutdown handling
// - Configuration management
// - Structured logging
//
// Environment Variables:
// - HTTP_PORT: Port for HTTP server (default: 3000)
// - GRPC_PORT: Port for gRPC server (default: 50051)
// - DB_URL: PostgreSQL connection string
// - REDIS_URL: Redis connection string
// - KAFKA_BROKERS: Comma-separated list of Kafka brokers
//
// @title Message Service API
// @version 1.0
// @description This is a RESTful API for managing messages
// @host localhost:3000
// @BasePath /api/v1
package main

import (
	"context"
	"fmt"
	"go-boilerplate/config"
	"go-boilerplate/internal/api/grpc"
	"go-boilerplate/internal/api/http"
	"go-boilerplate/internal/cache"
	"go-boilerplate/internal/kafka"
	"go-boilerplate/internal/middleware"
	"go-boilerplate/internal/service"
	pb "go-boilerplate/proto/message/v1"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	grpc_server "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "go-boilerplate/docs" // Import generated docs
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database connection pool
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	logger.Info("Database connection string", zap.String("connStr", connStr))
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Fatal("Failed to parse database config", zap.Error(err))
	}

	// Set connection pool settings with safe defaults
	if cfg.Database.MaxOpenConns <= 0 {
		logger.Warn("invalid max open connections, using default", zap.Int32("max_open_conns", cfg.Database.MaxOpenConns))
		dbConfig.MaxConns = 10 // safe default
	} else {
		dbConfig.MaxConns = cfg.Database.MaxOpenConns
	}

	if cfg.Database.MaxIdleConns <= 0 {
		logger.Warn("invalid max idle connections, using default", zap.Int32("max_idle_conns", cfg.Database.MaxIdleConns))
		dbConfig.MinConns = 2 // safe default
	} else {
		dbConfig.MinConns = cfg.Database.MaxIdleConns
	}

	// Create connection pool
	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}
	defer consumer.Close()

	// Initialize services
	messageService := service.NewMessageService(db, redisCache, producer)

	// Initialize HTTP handlers
	messageHandler := http.NewMessageHandler(messageService)

	// Initialize gRPC server
	grpcServer := grpc.NewMessageServer(messageService)

	// Start servers
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start HTTP server
	go func() {
		e := echo.New()

		// Set custom validator
		e.Validator = &middleware.CustomValidator{Validator: middleware.GetValidator()}

		// Middleware
		e.Use(echomiddleware.Logger())
		e.Use(echomiddleware.Recover())
		e.Use(echomiddleware.CORS())

		// Swagger docs
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		// Health check endpoints
		healthHandler := http.NewHealthHandler(db, redisCache)
		e.GET("/health", healthHandler.Health)
		e.GET("/health/live", healthHandler.LivenessProbe)
		e.GET("/health/ready", healthHandler.ReadinessProbe)

		// API routes
		v1 := e.Group("/api/v1")
		{
			messages := v1.Group("/messages")
			messages.POST("", messageHandler.CreateMessage)
			messages.GET("", messageHandler.ListMessages)
			messages.GET("/:id", messageHandler.GetMessage)
			messages.PUT("/:id", messageHandler.UpdateMessage)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
		}

		// Start server
		if err := e.Start(":" + cfg.Server.Port); err != nil {
			errChan <- fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}()

	// Start gRPC server
	go func() {
		listener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
		if err != nil {
			errChan <- fmt.Errorf("failed to listen: %w", err)
			return
		}

		server := grpc_server.NewServer()

		pb.RegisterMessageServiceServer(server, grpcServer)
		reflection.Register(server)

		logger.Info("Starting gRPC server", zap.String("port", cfg.GRPC.Port))
		if err := server.Serve(listener); err != nil {
			errChan <- fmt.Errorf("failed to start gRPC server: %w", err)
		}
	}()

	// Start Kafka consumer
	go func() {
		if err := consumer.Start(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start Kafka consumer: %w", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("Server error", zap.Error(err))
	case sig := <-sigChan:
		logger.Info("Received signal", zap.String("signal", sig.String()))
	}

	// Cleanup and shutdown
	cancel() // Stop Kafka consumer
	logger.Info("Shutting down servers")
}
