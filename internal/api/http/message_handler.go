package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-boilerplate/internal/middleware"
	"go-boilerplate/internal/models"
	"go-boilerplate/internal/service"
	"gorm.io/gorm"
	"net/http"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// CreateMessage godoc
// @Summary Create a new message
// @Description Create a new message with the provided content
// @Tags messages
// @Accept json
// @Produce json
// @Param message body CreateMessageRequest true "Message content"
// @Success 201 {object} models.Message
// @Router /api/v1/messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	req := c.MustGet("validated").(*CreateMessageRequest)

	message := &models.Message{
		Content: req.Content,
	}

	if err := h.messageService.CreateMessage(c.Request.Context(), message); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, message)
}

// GetMessage godoc
// @Summary Get a message by ID
// @Description Get a message by its unique identifier
// @Tags messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Router /api/v1/messages/{id} [get]
func (h *MessageHandler) GetMessage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(&middleware.ValidationError{
			Field:   "id",
			Message: "invalid UUID format",
		})
		return
	}

	message, err := h.messageService.GetMessage(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	if message == nil {
		c.Error(gorm.ErrRecordNotFound)
		return
	}

	c.JSON(http.StatusOK, message)
}

// ListMessages godoc
// @Summary List all messages
// @Description Get a list of all messages
// @Tags messages
// @Produce json
// @Success 200 {array} models.Message
// @Router /api/v1/messages [get]
func (h *MessageHandler) ListMessages(c *gin.Context) {
	req := c.MustGet("validated").(*ListMessagesRequest)

	messages, total, err := h.messageService.ListMessagesPaginated(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":  messages,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

type CreateMessageRequest struct {
	Content string `json:"content" binding:"required" validate:"required,min=1,max=1000"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required" validate:"required,min=1,max=1000"`
}

type ListMessagesRequest struct {
	Page     int `form:"page,default=1" validate:"min=1"`
	PageSize int `form:"page_size,default=10" validate:"min=1,max=100"`
}

// UpdateMessage godoc
// @Summary Update a message
// @Description Update a message's content by its ID
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param message body UpdateMessageRequest true "Updated message content"
// @Success 200 {object} models.Message
// @Router /api/v1/messages/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(&middleware.ValidationError{
			Field:   "id",
			Message: "invalid UUID format",
		})
		return
	}

	req := c.MustGet("validated").(*UpdateMessageRequest)

	message := &models.Message{
		ID:      id,
		Content: req.Content,
	}

	if err := h.messageService.UpdateMessage(c.Request.Context(), message); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, message)
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Delete a message by its ID
// @Tags messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 204 "No Content"
// @Router /api/v1/messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(&middleware.ValidationError{
			Field:   "id",
			Message: "invalid UUID format",
		})
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
