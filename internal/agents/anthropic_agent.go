package agents

import (
	"github.com/frozenkro/go-agent/internal/models/anthropic"
	"github.com/frozenkro/go-agent/internal/tools/bash"
)

type AnthropicAgent struct {
	bashSession bash.BashSession
}

func NewAnthropicAgent() (AnthropicAgent, error) {
	bs, err := bash.NewBashSession()
	if err != nil {
		return AnthropicAgent{}, err
	}

	return AnthropicAgent{
		bashSession: *bs,
	}, nil
}

func (a *AnthropicAgent) InitRequest(model anthropic.Model, prompt string, opts ...anthropic.AnthropicMessagesRequestOption) *anthropic.AnthropicMessagesRequest {

	req := &anthropic.AnthropicMessagesRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 1024,
		Messages: []anthropic.Message{
			anthropic.Message{
				Role: "user",
				Content: []anthropic.Content{
					anthropic.TextContent{
						BaseContent: anthropic.BaseContent{
							Type: anthropic.TEXT,
						},
						Text: prompt,
					},
				},
			},
		},
		Tools: []anthropic.AnthropicToolSpec{},
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

func (a *AnthropicAgent) HandleResponse(response *anthropic.AnthropicMessagesResponse) (*anthropic.AnthropicMessagesRequest, bool, error) {

	currentContent := response.Content[len(response.Content)-1]

	switch currentContent.GetType() {
	case anthropic.TOOL_USE:

	}

	return &anthropic.AnthropicMessagesRequest{}, false, nil
}
