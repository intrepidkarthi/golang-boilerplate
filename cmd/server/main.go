package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-boilerplate/config"
	"go-boilerplate/internal/api/grpc"
	"go-boilerplate/internal/api/http"
	"go-boilerplate/internal/cache"
	"go-boilerplate/internal/kafka"
	"go-boilerplate/internal/service"
	pb "go-boilerplate/proto/message/v1"
	"go.uber.org/zap"
	grpc_middleware "google.golang.org/grpc/middleware"
	"google.golang.org/grpc/reflection"
	grpc_server "google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	))
	if err != nil {
		logger.Fatal("Failed to parse database config", zap.Error(err))
	}

	// Set connection pool settings
	dbConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	dbConfig.MinConns = int32(cfg.Database.MaxIdleConns)

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
		router := gin.Default()
		
		// API routes
		v1 := router.Group("/api/v1")
		{
			messages := v1.Group("/messages")
			{
				messages.POST("", messageHandler.CreateMessage)
				messages.GET("", messageHandler.ListMessages)
				messages.GET("/:id", messageHandler.GetMessage)
			}
		}

		if err := router.Run(":" + cfg.Server.Port); err != nil {
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

		server := grpc_server.NewServer(
			grpc_server.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(),
			),
		)

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
