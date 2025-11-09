// derived from docs: https://docs.ollama.com/api/chat
package ollama

import (
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
	KeepAlive string        `json:"keep_alive"`
}

type Role string

const (
	SYSTEM    Role = "system"
	USER      Role = "user"
	ASSISTANT Role = "assistant"
	TOOL      Role = "tool"
)

type Message struct {
	Role      Role   `json:"role"`
	Content   string `json:"content"`
	Thinking  string `json:"thinking,omitempty"`
	ToolCalls []any  `json:"tool_calls,omitempty"`
}

type ToolSpec struct {
	Type     ToolType     `json:"type"` // always "function"
	Function ToolFunction `json:"function"`
}

type ToolType string

const Function ToolType = "function"

type ToolFunction struct {
	Name        string `json:"name"`
	Parameters  any    `json:"parameters"`
	Description string `json:"description"`
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
	Seed        int      `json:"seed"`
	Temperature int      `json:"temperature"`
	TopK        int      `json:"top_k"`
	TopP        int      `json:"top_p"`
	MinP        int      `json:"min_p"`
	Stop        []string `json:"stop"`
	NumCtx      int      `json:"num_ctx"`
	NumPredict  int      `json:"num_predict"`
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
	// TODO
	return nil
}

func (r *OllamaRequest) SetMaxTokens(val int) error {
	// TODO
	return nil
}

func (r *OllamaRequest) AddMessageGroup([]sessionmgr.Message) error {
	// TODO
	return nil
}
