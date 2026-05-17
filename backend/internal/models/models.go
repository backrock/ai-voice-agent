package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Messages  []Message `json:"messages,omitempty"`
}

type Message struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	SessionID string    `gorm:"index" json:"session_id"`
	Role      string    `json:"role"` // user, assistant
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

type Provider struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" gorm:"unique"`
	Type      string    `json:"type"` // openai, claude, ollama, etc
	Enabled   bool      `json:"enabled"`
	APIKey    string    `json:"-" gorm:"type:text"` // Not exposed in JSON
	BaseURL   string    `json:"base_url"`
	Models    string    `json:"models" gorm:"type:text"` // JSON array of available models
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSessionRequest struct {
	Title    string `json:"title" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	Model    string `json:"model" binding:"required"`
}

type SendMessageRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type UpdateProviderRequest struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}
