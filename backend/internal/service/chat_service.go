package service

import (
	"context"
	"fmt"
	"time"

	"github.com/backrock/ai-voice-agent/backend/internal/models"
	"github.com/backrock/ai-voice-agent/backend/internal/providers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService struct {
	db               *gorm.DB
	providersConfig map[string]interface{}
	providerCache   map[string]providers.LLMProvider
}

func NewChatService(db *gorm.DB, providersConfig interface{}) *ChatService {
	return &ChatService{
		db:               db,
		providersConfig: providersConfig.(map[string]interface{}),
		providerCache:   make(map[string]providers.LLMProvider),
	}
}

func (s *ChatService) CreateSession(title, provider, model string) (*models.Session, error) {
	session := &models.Session{
		ID:        uuid.New().String(),
		Title:     title,
		Provider:  provider,
		Model:     model,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (s *ChatService) GetSessions() ([]models.Session, error) {
	var sessions []models.Session
	if err := s.db.Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *ChatService) GetSessionByID(id string) (*models.Session, error) {
	var session models.Session
	if err := s.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *ChatService) DeleteSession(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除消息
		if err := tx.Delete(&models.Message{}, "session_id = ?", id).Error; err != nil {
			return err
		}
		// 删除会话
		return tx.Delete(&models.Session{}, "id = ?", id).Error
	})
}

func (s *ChatService) SendMessage(ctx context.Context, sessionID, content string) (*models.Message, error) {
	// 获取会话
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %v", err)
	}

	// 保存用户消息
	userMsg := &models.Message{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(userMsg).Error; err != nil {
		return nil, err
	}

	// 获取提供商
	provider, err := s.getProvider(session.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	// 准备消息列表
	messages := make([]providers.Message, 0)
	for _, msg := range session.Messages {
		messages = append(messages, providers.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	messages = append(messages, providers.Message{
		Role:    "user",
		Content: content,
	})

	// 获取AI响应
	response, err := provider.ChatCompletion(ctx, session.Model, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}

	// 保存AI消息
	assistantMsg := &models.Message{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "assistant",
		Content:   response,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(assistantMsg).Error; err != nil {
		return nil, err
	}

	return assistantMsg, nil
}

func (s *ChatService) SendMessageStream(ctx context.Context, sessionID, content string) (<-chan string, error) {
	// 获取会话
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %v", err)
	}

	// 保存用户消息
	userMsg := &models.Message{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(userMsg).Error; err != nil {
		return nil, err
	}

	// 获取提供商
	provider, err := s.getProvider(session.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	// 准备消息列表
	messages := make([]providers.Message, 0)
	for _, msg := range session.Messages {
		messages = append(messages, providers.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	messages = append(messages, providers.Message{
		Role:    "user",
		Content: content,
	})

	// 获取流式响应
	stream, err := provider.ChatCompletionStream(ctx, session.Model, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream response: %v", err)
	}

	// 创建输出通道并在后台保存完整响应
	ch := make(chan string)
	go func() {
		var fullResponse string
		defer func() {
			close(ch)
			// 保存完整的AI响应
			if fullResponse != "" {
				assistantMsg := &models.Message{
					ID:        uuid.New().String(),
					SessionID: sessionID,
					Role:      "assistant",
					Content:   fullResponse,
					CreatedAt: time.Now(),
				}
				s.db.Create(assistantMsg)
			}
		}()

		for chunk := range stream {
			fullResponse += chunk
			ch <- chunk
		}
	}()

	return ch, nil
}

func (s *ChatService) getProvider(providerType string) (providers.LLMProvider, error) {
	// 检查缓存
	if provider, exists := s.providerCache[providerType]; exists {
		return provider, nil
	}

	// 查询数据库获取提供商配置
	var dbProvider models.Provider
	if err := s.db.First(&dbProvider, "type = ? AND enabled = ?", providerType, true).Error; err != nil {
		return nil, fmt.Errorf("provider not found or not enabled")
	}

	// 创建提供商
	provider, err := providers.Create(providerType, dbProvider.APIKey, dbProvider.BaseURL)
	if err != nil {
		return nil, err
	}

	// 缓存
	s.providerCache[providerType] = provider
	return provider, nil
}
