package apimgrs

import (
	"fmt"
	"log"

	"github.com/frozenkro/goagent/internal/tools"
	"github.com/frozenkro/goagent/models"
	"github.com/frozenkro/goagent/models/anthropic"
)

type AnthropicRequestMgr struct {
	requestContext *anthropic.MessagesRequest
	toolInvoker    tools.ToolInvoker
}

type AnthropicRequestMgrOption func(*anthropic.MessagesRequest)

func WithTools(toolNames ...anthropic.ToolName) AnthropicRequestMgrOption {

	return func(a *anthropic.MessagesRequest) {

		toolMap := tools.InitToolMap()

		a.Tools = make([]anthropic.AnthropicToolSpec, len(toolNames))
		for i, toolName := range toolNames {
			toolMeta, err := toolMap.ToolMetaByName(toolName)

			if err == nil {
				a.Tools[i] = toolMeta.Spec
			} else {
				log.Printf("%v", err.Error())
			}
		}
	}
}

func WithMaxTokens(maxTokens int) AnthropicRequestMgrOption {
	return func(a *anthropic.MessagesRequest) {
		a.MaxTokens = maxTokens
	}
}

func NewAnthropicClient(model anthropic.Model, prompt string, opts ...AnthropicRequestMgrOption) (AnthropicRequestMgr, error) {
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

	req := &anthropic.MessagesRequest{
		Model:    model,
		Messages: messages,
	}

	for _, opt := range opts {
		opt(req)
	}

	return AnthropicRequestMgr{
		requestContext: req,
		toolInvoker:    ti,
	}, nil
}

// GetRequest gets the current running `MessagesRequest` object
// This request effectively serves as a running context of the agentic task
func (a *AnthropicRequestMgr) GetRequest() *anthropic.MessagesRequest {
	return a.requestContext
}

// HandleResponse accepts the anthropic api "/messages" response
// This handles stop sequences and invocation of tool calls
// Returns true if task is complete
func (a *AnthropicRequestMgr) HandleResponse(response *anthropic.MessagesResponse, out chan models.AgentEvent) (bool, error) {
	complete := false

	sysMsg := anthropic.Message{
		Role:    anthropic.ASSISTANT,
		Content: response.Content,
	}
	a.requestContext.Messages = append(a.requestContext.Messages, sysMsg)

	// TODO Handle these reasons appropriately
	switch response.StopReason {
	case anthropic.SR_END_TURN:
		complete = true
	case anthropic.SR_MAX_TOKENS:
		complete = true
	case anthropic.SR_STOP_SEQUENCE:
		complete = true
	case anthropic.SR_PAUSE_TURN:
		complete = true
	case anthropic.SR_REFUSAL:
		complete = true
	case anthropic.SR_TOOL_USE:
		complete = false
	}

	// Iterate over all response content to build subsequent request content and task context
	usrMsg := anthropic.Message{
		Role:    anthropic.USER,
		Content: []anthropic.Content{},
	}
	for _, c := range response.Content {

		switch c.GetType() {
		case anthropic.TOOL_USE:

			toolResultContent, err := a.handleToolUse(c, out)
			if err != nil {
				return complete, err
			}

			usrMsg.Content = append(usrMsg.Content, toolResultContent)

		case anthropic.TEXT:
			textContent, ok := c.(*anthropic.TextContent)
			if !ok {
				return complete, fmt.Errorf("Response content did not properly parse")
			}

			out <- models.AgentEvent{Message: textContent.Text}
		}
	}
	a.requestContext.Messages = append(a.requestContext.Messages, usrMsg)

	return complete, nil
}

func (a *AnthropicRequestMgr) handleToolUse(content anthropic.Content, out chan models.AgentEvent) (anthropic.ToolResultContent, error) {
	toolUseContent, ok := content.(*anthropic.ToolUseContent)
	if !ok {
		return anthropic.ToolResultContent{}, fmt.Errorf("Response content did not properly parse")
	}

	out <- models.AgentEvent{Message: fmt.Sprintf("[Tool Call]: %v", toolUseContent.Name)}

	toolResultContent, err := a.toolInvoker.Invoke(*toolUseContent)
	if err != nil {
		return anthropic.ToolResultContent{}, fmt.Errorf("Error occurred during tool invocation for tool '%v':\n%w", toolUseContent.Name, err)
	}

	return toolResultContent, nil

}
