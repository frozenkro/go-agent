// derived from docs: https://docs.ollama.com/api/chat
package ollama

import (
	"fmt"
	"strconv"

	"github.com/frozenkro/goagent/internal/sessionmgr"
	"github.com/frozenkro/goagent/internal/tools"
)

type OllamaRequest struct {
	Model     string        `json:"model"`
	Messages  []Message     `json:"messages"`
	Tools     []ToolSpec    `json:"tools"`
	Format    RequestFormat `json:"format"` // this can also be a json schema, which is not elaborated on in the docs
	Options   *Options      `json:"options,omitempty"`
	Stream    bool          `json:"stream"`
	Think     bool          `json:"think"`
	KeepAlive string        `json:"keep_alive,omitempty"`
}

type Role string

const (
	SYSTEM    Role = "system"
	USER      Role = "user"
	ASSISTANT Role = "assistant"
	TOOL      Role = "tool"
)

func (r Role) ToDomainRole() (sessionmgr.Role, error) {
	switch r {
	case SYSTEM:
		return sessionmgr.SYSTEM, nil
	case USER:
		return sessionmgr.USER, nil
	case TOOL:
		return sessionmgr.TOOL, nil
	case ASSISTANT:
		return sessionmgr.ASSISTANT, nil
	default:
		return 0, fmt.Errorf("Unsupported role '%v'", r)
	}
}

type Message struct {
	Role      Role             `json:"role"`
	Content   string           `json:"content"`
	Thinking  string           `json:"thinking,omitempty"`
	ToolName  string           `json:"tool_name,omitempty"`
	ToolCalls []OllamaToolCall `json:"tool_calls,omitempty"`
}

type OllamaToolCall struct {
	Function ToolCallFunction `json:"function"`
}
type ToolCallFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Arguments   map[string]any `json:"arguments"`
}

type ToolSpec struct {
	Type     ToolType     `json:"type"` // always "function"
	Function ToolFunction `json:"function"`
}

type ToolType string

const Function ToolType = "function"

type ToolFunction struct {
	Name        string     `json:"name"`
	Parameters  ToolParams `json:"parameters"`
	Description string     `json:"description"`
}

type ToolParams struct {
	Type       string                  `json:"type"` // always "object"
	Required   []string                `json:"required"`
	Properties map[string]ToolProperty `json:"properties"`
}

type ToolProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RequestFormat string

const JSON_FORMAT RequestFormat = "json"

type Options struct {
	Seed        int      `json:"seed,omitempty"`
	Temperature int      `json:"temperature,omitempty"`
	TopK        int      `json:"top_k,omitempty"`
	TopP        int      `json:"top_p,omitempty"`
	MinP        int      `json:"min_p,omitempty"`
	Stop        []string `json:"stop,omitempty"`
	NumCtx      int      `json:"num_ctx,omitempty"`
	NumPredict  int      `json:"num_predict,omitempty"`
}

func (r *OllamaRequest) Init(model, prompt string) error {
	r.Model = model

	r.Messages = make([]Message, 1)
	r.Messages[0] = Message{
		Role:    USER,
		Content: prompt,
	}

	return nil
}

func (r *OllamaRequest) AddTool(meta *tools.ToolMeta) error {
	spec := meta.Tool.Spec()
	props := map[string]ToolProperty{}
	for key, p := range spec.Parameters.Properties {
		props[key] = ToolProperty{
			Type:        p.Type,
			Description: p.Description,
		}
	}

	r.Tools = append(r.Tools, ToolSpec{
		Type: "function",
		Function: ToolFunction{
			Name:        meta.Name,
			Description: spec.Description,
			Parameters: ToolParams{
				Type:       "object",
				Required:   spec.Parameters.Required,
				Properties: props,
			},
		},
	})
	return nil
}

func (r *OllamaRequest) SetMaxTokens(val int) error {
	if r.Options == nil {
		r.Options = &Options{}
	}
	r.Options.NumPredict = val
	return nil
}

func (r *OllamaRequest) AddStatementGroup(sts []sessionmgr.Statement) error {
	if len(sts) == 0 {
		return nil
	}

	switch sts[0].Role {
	case sessionmgr.ASSISTANT:
		r.Messages = append(r.Messages, toAssistantMessage(sts))
	case sessionmgr.TOOL:
		res, err := toToolCallResults(sts)
		if err != nil {
			return err
		}
		r.Messages = append(r.Messages, res...)
	case sessionmgr.USER:
		r.Messages = append(r.Messages, toTextMessages(sts, USER)...)
	case sessionmgr.SYSTEM:
		r.Messages = append(r.Messages, toTextMessages(sts, SYSTEM)...)
	}
	return nil
}

func toAssistantMessage(sts []sessionmgr.Statement) Message {
	msg := Message{}

	for _, v := range sts {
		switch v.Type {
		case sessionmgr.TOOL_CALL:
			msg.ToolCalls = append(msg.ToolCalls, OllamaToolCall{
				Function: ToolCallFunction{
					Name:      v.ToolCall.Name,
					Arguments: v.ToolCall.Params,
				},
			})
		}
	}

	return msg
}

func toToolCallResults(sts []sessionmgr.Statement) ([]Message, error) {
	msgs := make([]Message, len(sts))

	for _, v := range sts {
		index, err := strconv.Atoi(v.ToolCall.Id)
		if err != nil {
			return nil, fmt.Errorf("Failed to cast tool call index '%v' to int", v.ToolCall.Id)
		}

		msgs[index] = Message{
			Role:     TOOL,
			ToolName: v.ToolCall.Name,
			Content:  v.Text,
		}
	}

	return msgs, nil
}

func toTextMessages(sts []sessionmgr.Statement, role Role) []Message {
	msgs := []Message{}

	for _, v := range sts {
		msgs = append(msgs, Message{
			Role:    role,
			Content: v.Text,
		})
	}

	return msgs
}
