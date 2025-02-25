package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-boilerplate/internal/cache"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/kafka"
	"go-boilerplate/internal/models"
	"time"
)

type MessageService struct {
	queries  *db.Queries
	pool     *pgxpool.Pool
	cache    *cache.RedisCache
	producer *kafka.Producer
}

func NewMessageService(pool *pgxpool.Pool, cache *cache.RedisCache, producer *kafka.Producer) *MessageService {
	return &MessageService{
		queries:  db.New(pool),
		pool:     pool,
		cache:    cache,
		producer: producer,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, message *models.Message) error {
	// Create message in database
	result, err := s.queries.CreateMessage(ctx, message.Content)
	if err != nil {
		return err
	}

	// Update message with database values
	*message = models.Message{
		ID:        result.ID,
		Content:   result.Content,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}

	// Cache the message
	if err := s.cache.Set(ctx, message.ID.String(), message, 24*time.Hour); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	// Publish message created event
	if err := s.producer.PublishMessage(message); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

func (s *MessageService) GetMessage(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	// Try to get from cache first
	var message models.Message
	if err := s.cache.Get(ctx, id.String(), &message); err == nil {
		return &message, nil
	}

	// If not in cache, get from database
	result, err := s.queries.GetMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	message = models.Message{
		ID:        result.ID,
		Content:   result.Content,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}

	// Cache the message for future requests
	if err := s.cache.Set(ctx, id.String(), &message, 24*time.Hour); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return &message, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	// Update message in database
	params := db.UpdateMessageParams{
		ID:      message.ID,
		Content: message.Content,
	}
	result, err := s.queries.UpdateMessage(ctx, params)
	if err != nil {
		return err
	}

	// Update message with latest values
	message.UpdatedAt = result.UpdatedAt.Time

	// Update cache
	if err := s.cache.Set(ctx, message.ID.String(), message, 24*time.Hour); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	// Publish message updated event
	if err := s.producer.PublishMessage(message); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	// Delete from database
	if err := s.queries.DeleteMessage(ctx, id); err != nil {
		return err
	}

	// Delete from cache
	if err := s.cache.Del(ctx, id.String()); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	// Publish message deleted event
	deleteEvent := &models.Message{ID: id}
	if err := s.producer.PublishMessage(deleteEvent); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

func (s *MessageService) ListMessages(ctx context.Context) ([]*models.Message, error) {
	params := db.ListMessagesParams{
		Limit:  100, // Default limit
		Offset: 0,
	}
	results, err := s.queries.ListMessages(ctx, params)
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, len(results))
	for i, result := range results {
		messages[i] = &models.Message{
			ID:        result.ID,
			Content:   result.Content,
			CreatedAt: result.CreatedAt.Time,
			UpdatedAt: result.UpdatedAt.Time,
		}
	}

	return messages, nil
}

func (s *MessageService) ListMessagesPaginated(ctx context.Context, page, pageSize uint32) ([]*models.Message, int64, error) {
	// Get total count
	total, err := s.queries.GetTotalMessages(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated messages
	offset := (page - 1) * pageSize
	
	// Safely convert uint32 to int32
	if pageSize > uint32(1<<31-1) || offset > uint32(1<<31-1) {
		return nil, 0, fmt.Errorf("pagination values too large")
	}
	
	results, err := s.queries.ListMessages(ctx, db.ListMessagesParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	// Convert to models
	messages := make([]*models.Message, len(results))
	for i, result := range results {
		messages[i] = &models.Message{
			ID:        result.ID,
			Content:   result.Content,
			CreatedAt: result.CreatedAt.Time,
			UpdatedAt: result.UpdatedAt.Time,
		}
	}

	return messages, total, nil
}
