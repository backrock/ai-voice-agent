package providers

import (
	"context"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
}

func NewOpenAIProvider(apiKey, baseURL string) (LLMProvider, error) {
	if apiKey == "" {
		return nil, NewError(ErrInvalidConfig, "OpenAI API key is required")
	}

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	return &OpenAIProvider{
		client: openai.NewClientWithConfig(config),
	}, nil
}

func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

func (p *OpenAIProvider) Type() string {
	return "openai"
}

func (p *OpenAIProvider) GetAvailableModels() []string {
	return []string{
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	}
}

func (p *OpenAIProvider) ChatCompletion(ctx context.Context, model string, messages []Message) (string, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: openaiMessages,
	})
	if err != nil {
		return "", NewError(ErrAPIError, err.Error())
	}

	if len(resp.Choices) == 0 {
		return "", NewError(ErrNoResponse, "no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

func (p *OpenAIProvider) ChatCompletionStream(ctx context.Context, model string, messages []Message) (<-chan string, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: openaiMessages,
	})
	if err != nil {
		return nil, NewError(ErrAPIError, err.Error())
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
				ch <- resp.Choices[0].Delta.Content
			}
		}
	}()

	return ch, nil
}

func init() {
	Register("openai", NewOpenAIProvider)
}
