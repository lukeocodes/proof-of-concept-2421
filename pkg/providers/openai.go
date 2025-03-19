package providers

import (
	"context"

	"github.com/openai/openai-go"
)

// OpenAIClient implements the Provider interface for OpenAI
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient() (*OpenAIClient, error) {
	client := openai.NewClient()
	// TODO: Implement proper configuration loading
	return &OpenAIClient{
		client: client,
	}, nil
}

// ChatCompletion implements the Provider interface for OpenAI
func (c *OpenAIClient) ChatCompletion(ctx context.Context, messages []any) (any, error) {
	// Type assert the messages to OpenAI's expected type
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		if m, ok := msg.(openai.ChatCompletionMessageParamUnion); ok {
			openaiMessages[i] = m
		}
	}

	result, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.F(openai.ChatModelO1),
		Messages: openai.F(openaiMessages),
	})
	if err != nil {
		return nil, err
	}

	return result.Choices[0].Message, nil
}

// SummariseMessages implements the Provider interface for OpenAI
func (c *OpenAIClient) SummariseMessages(messages []any) (any, error) {
	// Type assert the messages to OpenAI's expected type
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		if m, ok := msg.(openai.ChatCompletionMessageParamUnion); ok {
			openaiMessages[i] = m
		}
	}

	// TODO: Implement OpenAI message summarisation
	return openai.ChatCompletionMessage{}, nil
}

func MapOpenAIProviderMessage(message ProviderMessage) openai.ChatCompletionMessageParamUnion {
	return openai.ChatCompletionMessage{
		Content:   message.Content,
		Role:      openai.ChatCompletionMessageRole(message.Role),
		ToolCalls: message.ToolCalls,
	}
}

func UnmapOpenAIProviderMessage(message openai.ChatCompletionMessage) ProviderMessage {
	return ProviderMessage{
		Content:   message.Content,
		Role:      ProviderMessageRole(message.Role),
		ToolCalls: message.ToolCalls,
	}
}
