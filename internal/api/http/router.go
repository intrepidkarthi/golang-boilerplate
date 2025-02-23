package http

import (
	"github.com/gin-gonic/gin"
	"go-boilerplate/internal/middleware"
	"go-boilerplate/internal/service"
)

// SetupRouter initializes the HTTP router and registers all routes
func SetupRouter(messageService *service.MessageService) *gin.Engine {
	router := gin.Default()

	// Add global middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.ValidationMiddleware())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		messages := v1.Group("/messages")
		handler := NewMessageHandler(messageService)

		messages.POST("", middleware.RequestValidator(&CreateMessageRequest{}), handler.CreateMessage)
		messages.GET("", middleware.RequestValidator(&ListMessagesRequest{}), handler.ListMessages)
		messages.GET("/:id", handler.GetMessage)
		messages.PUT("/:id", middleware.RequestValidator(&UpdateMessageRequest{}), handler.UpdateMessage)
		messages.DELETE("/:id", handler.DeleteMessage)
	}

	return router
}
