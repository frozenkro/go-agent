package agents

import (
	"fmt"
	"log"

	"github.com/frozenkro/go-agent/internal/tools"
	"github.com/frozenkro/go-agent/models/anthropic"
)

type AnthropicAgent struct {
	requestContext *anthropic.AnthropicMessagesRequest
	toolInvoker    tools.ToolInvoker
}

type AnthropicAgentOption func(*anthropic.AnthropicMessagesRequest)

func WithTools(toolNames ...anthropic.ToolName) AnthropicAgentOption {

	return func(a *anthropic.AnthropicMessagesRequest) {

		toolMap := tools.InitToolMap()

		for _, toolName := range toolNames {
			toolMeta, err := toolMap.ToolMetaByName(toolName)

			if err == nil {
				a.Tools = append(a.Tools, toolMeta.Spec)
			} else {
				log.Printf(err.Error())
			}
		}
	}
}

func NewAnthropicAgent(model anthropic.Model, prompt string, opts ...AnthropicAgentOption) (AnthropicAgent, error) {
	ti := tools.NewToolInvoker()

	messages := []anthropic.Message{
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
	}

	req := &anthropic.AnthropicMessagesRequest{
		Model:     model,
		MaxTokens: 1024,
		Messages:  messages,
		Tools:     []anthropic.AnthropicToolSpec{},
	}

	for _, opt := range opts {
		opt(req)
	}

	return AnthropicAgent{
		requestContext: &anthropic.AnthropicMessagesRequest{},
		toolInvoker:    ti,
	}, nil
}

func (a *AnthropicAgent) GetRequest() *anthropic.AnthropicMessagesRequest {
	return a.requestContext
}

func (a *AnthropicAgent) HandleResponse(response *anthropic.AnthropicMessagesResponse) (*anthropic.AnthropicMessagesRequest, bool, error) {

	currentContent := response.Content[len(response.Content)-1]

	newMessage := anthropic.Message{
		Role: anthropic.USER,
	}

	switch currentContent.GetType() {
	case anthropic.TOOL_USE:
		toolUseContent, ok := currentContent.(*anthropic.ToolUseContent)
		if !ok {
			return nil, false, fmt.Errorf("Response content did not properly parse")
		}

		toolResultContent, err := a.toolInvoker.Invoke(*toolUseContent)
		if err != nil {
			return nil, false, fmt.Errorf("Error occurred during tool invocation for tool '%v':\n%w", toolUseContent.Name, err)
		}

		newMessage.Content = []anthropic.Content{toolResultContent}

	}

	return &anthropic.AnthropicMessagesRequest{}, false, nil
}
