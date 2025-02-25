package http

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go-boilerplate/internal/service"
)

// SetupRouter initializes the HTTP router and registers all routes
func SetupRouter(messageService *service.MessageService) *echo.Echo {
	e := echo.New()

	// Add global middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// API v1 routes
	handler := NewMessageHandler(messageService)
	v1 := e.Group("/api/v1")
	messages := v1.Group("/messages")

	messages.POST("", handler.CreateMessage)
	messages.GET("", handler.ListMessages)
	messages.GET("/:id", handler.GetMessage)
	messages.PUT("/:id", handler.UpdateMessage)
	messages.DELETE("/:id", handler.DeleteMessage)

	return e
}
