package agents

import (
	"fmt"
	"log"

	"github.com/frozenkro/go-agent/internal/tools"
	"github.com/frozenkro/go-agent/models/anthropic"
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

func (a *AnthropicRequestMgr) GetRequest() *anthropic.MessagesRequest {
	return a.requestContext
}

func (a *AnthropicRequestMgr) HandleResponse(response *anthropic.MessagesResponse) (*anthropic.MessagesRequest, bool, error) {
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
		usrMsg, err := a.getToolCallResponses(response.Content)
		if err != nil {
			return a.requestContext, complete, err
		}
		a.requestContext.Messages = append(a.requestContext.Messages, usrMsg)
	}

	return a.requestContext, complete, nil
}

func (a *AnthropicRequestMgr) getToolCallResponses(content []anthropic.Content) (anthropic.Message, error) {
	usrMsg := anthropic.Message{
		Role:    anthropic.USER,
		Content: []anthropic.Content{},
	}

	for _, c := range content {

		if c.GetType() == anthropic.TOOL_USE {

			toolUseContent, ok := c.(*anthropic.ToolUseContent)
			if !ok {
				return usrMsg, fmt.Errorf("Response content did not properly parse")
			}

			toolResultContent, err := a.toolInvoker.Invoke(*toolUseContent)
			if err != nil {
				return usrMsg, fmt.Errorf("Error occurred during tool invocation for tool '%v':\n%w", toolUseContent.Name, err)
			}

			usrMsg.Content = append(usrMsg.Content, toolResultContent)
		}
	}

	return usrMsg, nil

}
