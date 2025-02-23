package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-boilerplate/internal/models"
	"testing"
	"time"
)

// MockDB is a mock implementation of the database interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) WithContext(ctx context.Context) interface{} {
	return m
}

func (m *MockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) error {
	args := m.Called(dest, conds)
	return args.Error(0)
}

func (m *MockDB) Save(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) error {
	args := m.Called(value, conds)
	return args.Error(0)
}

// MockCache is a mock implementation of the cache interface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) SetMessage(ctx context.Context, message *models.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockCache) GetMessage(ctx context.Context, id string) (*models.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Message), args.Error(1)
}

// MockProducer is a mock implementation of the Kafka producer
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) PublishMessage(message *models.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestMessageService_CreateMessage(t *testing.T) {
	// Setup
	mockDB := new(MockDB)
	mockCache := new(MockCache)
	mockProducer := new(MockProducer)
	
	service := NewMessageService(mockDB, mockCache, mockProducer)
	
	ctx := context.Background()
	message := &models.Message{
		Content: "Test message",
	}
	
	// Expectations
	mockDB.On("Create", message).Return(nil)
	mockCache.On("SetMessage", ctx, message).Return(nil)
	mockProducer.On("PublishMessage", message).Return(nil)
	
	// Test
	err := service.CreateMessage(ctx, message)
	
	// Assertions
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestMessageService_GetMessage(t *testing.T) {
	// Setup
	mockDB := new(MockDB)
	mockCache := new(MockCache)
	mockProducer := new(MockProducer)
	
	service := NewMessageService(mockDB, mockCache, mockProducer)
	
	ctx := context.Background()
	id := uuid.New()
	expectedMessage := &models.Message{
		ID:        id,
		Content:   "Test message",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Expectations
	mockCache.On("GetMessage", ctx, id.String()).Return(expectedMessage, nil)
	
	// Test
	message, err := service.GetMessage(ctx, id)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, message)
	mockCache.AssertExpectations(t)
}

func TestMessageService_UpdateMessage(t *testing.T) {
	// Setup
	mockDB := new(MockDB)
	mockCache := new(MockCache)
	mockProducer := new(MockProducer)
	
	service := NewMessageService(mockDB, mockCache, mockProducer)
	
	ctx := context.Background()
	id := uuid.New()
	message := &models.Message{
		ID:      id,
		Content: "Updated content",
	}
	
	// Expectations
	mockDB.On("Save", message).Return(nil)
	mockCache.On("SetMessage", ctx, message).Return(nil)
	mockProducer.On("PublishMessage", message).Return(nil)
	
	// Test
	err := service.UpdateMessage(ctx, message)
	
	// Assertions
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestMessageService_DeleteMessage(t *testing.T) {
	// Setup
	mockDB := new(MockDB)
	mockCache := new(MockCache)
	mockProducer := new(MockProducer)
	
	service := NewMessageService(mockDB, mockCache, mockProducer)
	
	ctx := context.Background()
	id := uuid.New()
	
	// Expectations
	mockDB.On("Delete", &models.Message{}, "id = ?", id).Return(nil)
	
	// Test
	err := service.DeleteMessage(ctx, id)
	
	// Assertions
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}
