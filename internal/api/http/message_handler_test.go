package http

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-boilerplate/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(ctx interface{}, message *models.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageService) GetMessage(ctx interface{}, id uuid.UUID) (*models.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageService) UpdateMessage(ctx interface{}, message *models.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageService) DeleteMessage(ctx interface{}, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMessageService) ListMessages(ctx interface{}) ([]models.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Message), args.Error(1)
}

func setupTestRouter(mockService *MockMessageService) *echo.Echo {
	e := echo.New()
	
	handler := NewMessageHandler(mockService)
	v1 := e.Group("/api/v1")
	messages := v1.Group("/messages")
	messages.POST("", handler.CreateMessage)
	messages.GET("", handler.ListMessages)
	messages.GET("/:id", handler.GetMessage)
	messages.PUT("/:id", handler.UpdateMessage)
	messages.DELETE("/:id", handler.DeleteMessage)
	
	return e
}

func TestCreateMessage(t *testing.T) {
	mockService := new(MockMessageService)
	router := setupTestRouter(mockService)

	message := &models.Message{
		Content: "Test message",
	}

	mockService.On("CreateMessage", mock.Anything, mock.MatchedBy(func(m *models.Message) bool {
		return m.Content == message.Content
	})).Return(nil)

	body, _ := json.Marshal(CreateMessageRequest{Content: message.Content})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)

	assert.NoError(t, handler.CreateMessage(c))

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetMessage(t *testing.T) {
	mockService := new(MockMessageService)
	router := setupTestRouter(mockService)

	id := uuid.New()
	message := &models.Message{
		ID:        id,
		Content:   "Test message",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetMessage", mock.Anything, id).Return(message, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages/"+id.String(), nil)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	assert.NoError(t, handler.GetMessage(c))

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, message.ID, response.ID)
	assert.Equal(t, message.Content, response.Content)

	mockService.AssertExpectations(t)
}

func TestUpdateMessage(t *testing.T) {
	mockService := new(MockMessageService)
	router := setupTestRouter(mockService)

	id := uuid.New()
	updatedContent := "Updated message"

	mockService.On("UpdateMessage", mock.Anything, mock.MatchedBy(func(m *models.Message) bool {
		return m.ID == id && m.Content == updatedContent
	})).Return(nil)

	body, _ := json.Marshal(UpdateMessageRequest{Content: updatedContent})
	req := httptest.NewRequest("PUT", "/api/v1/messages/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteMessage(t *testing.T) {
	mockService := new(MockMessageService)
	router := setupTestRouter(mockService)

	id := uuid.New()

	mockService.On("DeleteMessage", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/messages/"+id.String(), nil)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	assert.NoError(t, handler.DeleteMessage(c))

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestListMessages(t *testing.T) {
	mockService := new(MockMessageService)
	router := setupTestRouter(mockService)

	messages := []models.Message{
		{
			ID:        uuid.New(),
			Content:   "Message 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Content:   "Message 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockService.On("ListMessages", mock.Anything).Return(messages, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)

	assert.NoError(t, handler.ListMessages(c))

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, messages[0].Content, response[0].Content)
	assert.Equal(t, messages[1].Content, response[1].Content)

	mockService.AssertExpectations(t)
}
