package api

import (
	"github.com/backrock/ai-voice-agent/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, chatService *service.ChatService, configService *service.ConfigService) {
	// Chat API
	chat := router.Group("/api/v1/chat")
	{
		// 会话管理
		chat.POST("/sessions", CreateSession(chatService))
		chat.GET("/sessions", GetSessions(chatService))
		chat.GET("/sessions/:id", GetSession(chatService))
		chat.DELETE("/sessions/:id", DeleteSession(chatService))

		// 消息
		chat.POST("/messages", SendMessage(chatService))
		chat.POST("/messages/stream", SendMessageStream(chatService))
	}

	// Admin/Config API
	admin := router.Group("/api/v1/admin")
	{
		// 提供商管理
		admin.GET("/providers", GetProviders(configService))
		admin.GET("/models", GetAvailableModels(configService))
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
