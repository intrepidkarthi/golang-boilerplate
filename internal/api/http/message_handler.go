package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-boilerplate/internal/models"
	"go-boilerplate/internal/service"
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
func (h *MessageHandler) CreateMessage(c echo.Context) error {
	req := new(CreateMessageRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	message := &models.Message{
		Content: req.Content,
	}

	if err := h.messageService.CreateMessage(c.Request().Context(), message); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, message)
}

// GetMessage godoc
// @Summary Get a message by ID
// @Description Get a message by its unique identifier
// @Tags messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Router /api/v1/messages/{id} [get]
func (h *MessageHandler) GetMessage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid UUID format")
	}

	message, err := h.messageService.GetMessage(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if message == nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	return c.JSON(http.StatusOK, message)
}

// ListMessages godoc
// @Summary List all messages
// @Description Get a list of all messages
// @Tags messages
// @Produce json
// @Success 200 {array} models.Message
// @Router /api/v1/messages [get]
func (h *MessageHandler) ListMessages(c echo.Context) error {
	req := &ListMessagesRequest{}

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Set defaults if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	messages, total, err := h.messageService.ListMessagesPaginated(c.Request().Context(), int32(req.Page), int32(req.PageSize))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages":  messages,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type ListMessagesRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
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
func (h *MessageHandler) UpdateMessage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid UUID format")
	}

	req := new(UpdateMessageRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	message := &models.Message{
		ID:      id,
		Content: req.Content,
	}

	if err := h.messageService.UpdateMessage(c.Request().Context(), message); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, message)
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Delete a message by its ID
// @Tags messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 204 "No Content"
// @Router /api/v1/messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid UUID format")
	}

	if err := h.messageService.DeleteMessage(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
