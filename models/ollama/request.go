// derived from docs: https://docs.ollama.com/api/chat
package ollama

type OllamaRequest struct {
	Model     string              `json:"model"`
	Messages  []Message           `json:"messages"`
	Tools     []Tool              `json:"tools"`
	Format    OllamaRequestFormat `json:"format"` // this can also be a json schema, which is not elaborated on in the docs
	Options   *Options            `json:"options,omitempty"`
	Stream    bool                `json:"stream"`
	Think     bool                `json:"think"`
	KeepAlive string              `json:"keep_alive"`
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
	Thinking  string `json:"thinking"`
	ToolCalls []any  `json:"tool_calls,omitempty"`
}

type Tool struct {
	Type OllamaToolType // always "function"
}

type OllamaToolType string

const Function OllamaToolType = "function"

type OllamaToolFunction struct {
	Name        string `json:"name"`
	Parameters  any    `json:"parameters"`
	Description string `json:"description"`
}

type OllamaRequestFormat string

const JSON_FORMAT OllamaRequestFormat = "json"

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
