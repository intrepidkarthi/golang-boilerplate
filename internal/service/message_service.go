package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-boilerplate/internal/cache"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/kafka"
	"go-boilerplate/internal/models"
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
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	// Cache the message
	if err := s.cache.SetMessage(ctx, message); err != nil {
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
	message, err := s.cache.GetMessage(ctx, id.String())
	if err != nil {
		// Log error but continue to database
		// TODO: Add proper logging
	}

	if message != nil {
		return message, nil
	}

	// If not in cache, get from database
	result, err := s.queries.GetMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	message = &models.Message{
		ID:        result.ID,
		Content:   result.Content,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	// Cache the message for future requests
	if err := s.cache.SetMessage(ctx, message); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return message, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	// Update message in database
	result, err := s.queries.UpdateMessage(ctx, db.UpdateMessageParams{
		ID:      message.ID,
		Content: message.Content,
	})
	if err != nil {
		return err
	}

	// Update message with latest values
	message.UpdatedAt = result.UpdatedAt

	// Update cache
	if err := s.cache.SetMessage(ctx, message); err != nil {
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
	if err := s.cache.DeleteMessage(ctx, id.String()); err != nil {
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
	results, err := s.queries.ListMessages(ctx, db.ListMessagesParams{
		Limit:  100, // Default limit
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, len(results))
	for i, result := range results {
		messages[i] = &models.Message{
			ID:        result.ID,
			Content:   result.Content,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}
	}

	return messages, nil
}

func (s *MessageService) ListMessagesPaginated(ctx context.Context, page, pageSize int32) ([]*models.Message, int64, error) {
	// Get total count
	total, err := s.queries.GetTotalMessages(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated messages
	offset := (page - 1) * pageSize
	results, err := s.queries.ListMessages(ctx, db.ListMessagesParams{
		Limit:  pageSize,
		Offset: offset,
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
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}
	}

	return messages, total, nil
}
