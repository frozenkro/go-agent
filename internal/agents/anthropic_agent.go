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

	toolNameMap := make(map[anthropic.ToolName]anthropic.Tool)
	toolNameMap[anthropic.BASH] = anthropic.NewBashTool()
	toolNameMap[anthropic.TEXT_EDITOR] = anthropic.NewTextEditorTool()

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
		Tools: []anthropic.Tool{},
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

func (a *AnthropicAgent) HandleResponse(response *anthropic.AnthropicMessagesResponse) (*anthropic.AnthropicMessagesRequest, bool, error) {
	return &anthropic.AnthropicMessagesRequest{}, false, nil
}
