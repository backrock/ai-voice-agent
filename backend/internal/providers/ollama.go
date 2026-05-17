package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaProvider struct {
	baseURL    string
	httpClient *http.Client
}

func NewOllamaProvider(apiKey, baseURL string) (LLMProvider, error) {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &OllamaProvider{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}, nil
}

func (p *OllamaProvider) Name() string {
	return "Ollama"
}

func (p *OllamaProvider) Type() string {
	return "ollama"
}

func (p *OllamaProvider) GetAvailableModels() []string {
	return []string{
		"llama2",
		"mistral",
		"neural-chat",
		"orca-mini",
	}
}

type ollamaRequest struct {
	Model    string        `json:"model"`
	Messages []ollamaMsg   `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ollamaMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message ollamaMsg `json:"message"`
	Done    bool      `json:"done"`
}

func (p *OllamaProvider) ChatCompletion(ctx context.Context, model string, messages []Message) (string, error) {
	ollamaMessages := make([]ollamaMsg, len(messages))
	for i, msg := range messages {
		ollamaMessages[i] = ollamaMsg{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	reqBody := ollamaRequest{
		Model:    model,
		Messages: ollamaMessages,
		Stream:   false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", NewError(ErrInvalidRequest, err.Error())
	}

	resp, err := p.httpClient.Post(
		p.baseURL+"/api/chat",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", NewError(ErrAPIError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewError(ErrAPIError, fmt.Sprintf("Ollama API error: %d", resp.StatusCode))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", NewError(ErrParseError, err.Error())
	}

	return ollamaResp.Message.Content, nil
}

func (p *OllamaProvider) ChatCompletionStream(ctx context.Context, model string, messages []Message) (<-chan string, error) {
	ollamaMessages := make([]ollamaMsg, len(messages))
	for i, msg := range messages {
		ollamaMessages[i] = ollamaMsg{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	reqBody := ollamaRequest{
		Model:    model,
		Messages: ollamaMessages,
		Stream:   true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewError(ErrInvalidRequest, err.Error())
	}

	resp, err := p.httpClient.Post(
		p.baseURL+"/api/chat",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, NewError(ErrAPIError, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, NewError(ErrAPIError, fmt.Sprintf("Ollama API error: %d", resp.StatusCode))
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var ollamaResp ollamaResponse
			if err := decoder.Decode(&ollamaResp); err != nil {
				if err == io.EOF {
					return
				}
				return
			}

			if ollamaResp.Message.Content != "" {
				ch <- ollamaResp.Message.Content
			}

			if ollamaResp.Done {
				return
			}
		}
	}()

	return ch, nil
}

func init() {
	Register("ollama", NewOllamaProvider)
}
