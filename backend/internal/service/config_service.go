package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/backrock/ai-voice-agent/backend/internal/models"
	"github.com/backrock/ai-voice-agent/backend/internal/providers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

func (s *ConfigService) GetProviders() ([]models.Provider, error) {
	var providers []models.Provider
	if err := s.db.Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func (s *ConfigService) GetProvider(id string) (*models.Provider, error) {
	var provider models.Provider
	if err := s.db.First(&provider, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func (s *ConfigService) CreateProvider(name, providerType, apiKey, baseURL string) (*models.Provider, error) {
	provider := &models.Provider{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      providerType,
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *ConfigService) UpdateProvider(id string, updates map[string]interface{}) (*models.Provider, error) {
	var provider models.Provider
	if err := s.db.First(&provider, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&provider).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &provider, nil
}

func (s *ConfigService) DeleteProvider(id string) error {
	return s.db.Delete(&models.Provider{}, "id = ?", id).Error
}

func (s *ConfigService) GetAvailableModels(providerType string) ([]string, error) {
	// 这里可以调用实际的提供商来获取模型列表
	// 现在返回硬编码的列表
	modelMap := map[string][]string{
		"openai": {
			"gpt-4-turbo",
			"gpt-4",
			"gpt-3.5-turbo",
		},
		"claude": {
			"claude-3-opus",
			"claude-3-sonnet",
			"claude-3-haiku",
		},
		"ollama": {
			"llama2",
			"mistral",
			"neural-chat",
		},
	}

	if models, ok := modelMap[providerType]; ok {
		return models, nil
	}

	return nil, fmt.Errorf("provider type not found: %s", providerType)
}

func (s *ConfigService) TestProvider(providerType, apiKey, baseURL string) error {
	provider, err := providers.Create(providerType, apiKey, baseURL)
	if err != nil {
		return err
	}

	if provider == nil {
		return fmt.Errorf("provider not available")
	}

	return nil
}

func (s *ConfigService) GetProviderModels(id string) ([]string, error) {
	var provider models.Provider
	if err := s.db.First(&provider, "id = ?", id).Error; err != nil {
		return nil, err
	}

	var models []string
	if err := json.Unmarshal([]byte(provider.Models), &models); err != nil {
		return nil, err
	}

	return models, nil
}

func (s *ConfigService) UpdateProviderModels(id string, models []string) error {
	modelJSON, err := json.Marshal(models)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Provider{}).Where("id = ?", id).Update("models", string(modelJSON)).Error
}
