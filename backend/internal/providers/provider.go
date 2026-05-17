package providers

import (
	"context"
)

// LLMProvider 定义大模型提供商接口
type LLMProvider interface {
	// Name 返回提供商名称
	Name() string
	// Type 返回提供商类型
	Type() string
	// GetAvailableModels 获取可用模型列表
	GetAvailableModels() []string
	// ChatCompletion 发送聊天消息并获取完整响应
	ChatCompletion(ctx context.Context, model string, messages []Message) (string, error)
	// ChatCompletionStream 发送聊天消息并流式返回响应
	ChatCompletionStream(ctx context.Context, model string, messages []Message) (<-chan string, error)
}

// Message 定义消息结构
type Message struct {
	Role    string `json:"role"` // "user" 或 "assistant"
	Content string `json:"content"`
}

// ProviderFactory 工厂函数类型
type ProviderFactory func(apiKey, baseURL string) (LLMProvider, error)

// registry 提供商注册表
var registry = make(map[string]ProviderFactory)

// Register 注册提供商
func Register(providerType string, factory ProviderFactory) {
	registry[providerType] = factory
}

// Create 创建提供商实例
func Create(providerType, apiKey, baseURL string) (LLMProvider, error) {
	factory, exists := registry[providerType]
	if !exists {
		return nil, NewError(ErrProviderNotFound, "provider not found: "+providerType)
	}
	return factory(apiKey, baseURL)
}

// GetRegisteredProviders 获取所有已注册的提供商类型
func GetRegisteredProviders() []string {
	providers := make([]string, 0, len(registry))
	for k := range registry {
		providers = append(providers, k)
	}
	return providers
}
