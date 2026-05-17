package api

import (
	"net/http"

	"github.com/backrock/ai-voice-agent/backend/internal/models"
	"github.com/backrock/ai-voice-agent/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func CreateSession(chatService *service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session, err := chatService.CreateSession(req.Title, req.Provider, req.Model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func GetSessions(chatService *service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions, err := chatService.GetSessions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}

func GetSession(chatService *service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		session, err := chatService.GetSessionByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func DeleteSession(chatService *service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := chatService.DeleteSession(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
	}
}

func SendMessage(chatService *service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SendMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		msg, err := chatService.SendMessage(c.Request.Context(), req.SessionID, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, msg)
	}
}

func SendMessageStream(chatService *service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SendMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stream, err := chatService.SendMessageStream(c.Request.Context(), req.SessionID, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		for chunk := range stream {
			c.Writer.WriteString("data: " + chunk + "\n\n")
			c.Writer.Flush()
		}
	}
}

func GetProviders(configService *service.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		providers, err := configService.GetProviders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, providers)
	}
}

func GetAvailableModels(configService *service.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		providerType := c.Query("type")
		if providerType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "provider type is required"})
			return
		}

		models, err := configService.GetAvailableModels(providerType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"models": models})
	}
}
