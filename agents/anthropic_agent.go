package agents

import (
	"fmt"
	"log"

	"github.com/frozenkro/go-agent/internal/tools"
	"github.com/frozenkro/go-agent/models/anthropic"
	"github.com/frozenkro/go-agent/models/anthropic/content"
	toolModels "github.com/frozenkro/go-agent/models/anthropic/tools"
)

type AnthropicAgent struct {
	requestContext *anthropic.AnthropicMessagesRequest
	toolInvoker    tools.ToolInvoker
}

type AnthropicAgentOption func(*anthropic.AnthropicMessagesRequest)

func WithTools(toolNames ...toolModels.ToolName) AnthropicAgentOption {

	return func(a *anthropic.AnthropicMessagesRequest) {

		toolMap := tools.InitToolMap()

		a.Tools = make([]toolModels.AnthropicToolSpec, len(toolNames))
		for i, toolName := range toolNames {
			toolMeta, err := toolMap.ToolMetaByName(toolName)

			if err == nil {
				a.Tools[i] = toolMeta.Spec
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
			Content: []content.Content{
				content.TextContent{
					BaseContent: content.BaseContent{
						Type: content.TEXT,
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
	}

	for _, opt := range opts {
		opt(req)
	}

	return AnthropicAgent{
		requestContext: req,
		toolInvoker:    ti,
	}, nil
}

func (a *AnthropicAgent) GetRequest() *anthropic.AnthropicMessagesRequest {
	return a.requestContext
}

func (a *AnthropicAgent) HandleResponse(response *anthropic.MessagesResponse) (*anthropic.AnthropicMessagesRequest, bool, error) {

	existingMessage := anthropic.Message{
		Role:    anthropic.ASSISTANT,
		Content: response.Content,
	}

	newMessage := anthropic.Message{
		Role: anthropic.USER,
	}
	newMessage.Content = make([]content.Content, 0)

	switch response.StopReason {
	case anthropic.TOOL_USE:

		for _, c := range response.Content {

			if c.GetType() == content.TOOL_USE {
				toolUseContent, ok := c.(*content.ToolUseContent)
				if !ok {
					return nil, false, fmt.Errorf("Response content did not properly parse")
				}

				toolResultContent, err := a.toolInvoker.Invoke(*toolUseContent)
				if err != nil {
					return nil, false, fmt.Errorf("Error occurred during tool invocation for tool '%v':\n%w", toolUseContent.Name, err)
				}

				newMessage.Content = append(newMessage.Content, toolResultContent)

			}
		}
	}

	a.requestContext.Messages = append(a.requestContext.Messages, existingMessage, newMessage)

	return a.requestContext, false, nil
}
