package anthropic

type Message struct {
	Role    Role      `json:"role"`
	Content []Content `json:"content"`
}

type Role string

const (
	USER      Role = "user"
	ASSISTANT Role = "assistant"
)

type ToolConfiguration struct {
	AllowedTools []string `json:"allowed_tools,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
}

type MCPServer struct {
	Name               string             `json:"name"`
	Type               string             `json:"type"`
	Url                string             `json:"url"`
	AuthorizationToken string             `json:"authorization_token,omitempty"`
	ToolConfiguration  *ToolConfiguration `json:"tool_configuration,omitempty"`
}

type Metadata struct {
	UserId string `json:"user_id"`
}

type ThinkingData struct {
	BudgetTokens int    `json:"budget_tokens"`
	Type         string `json:"type"`
}

type ToolName string

const (
	BASH        ToolName = "bash"
	TEXT_EDITOR ToolName = "str_replace_based_edit_tool"
)

type AnthropicToolSpec interface {
	GetType() string
	GetName() ToolName
}

type BaseTool struct {
	Name ToolName `json:"name"`
	Type string   `json:"type"`
}

func (t BaseTool) GetType() string {
	return t.Type
}
func (t BaseTool) GetName() ToolName {
	return t.Name
}

type BashTool struct {
	BaseTool
	CacheControl *CacheControl `json:"cache_control"`
}

func NewBashTool() BashTool {
	return BashTool{
		BaseTool: BaseTool{Type: "bash_20250124", Name: BASH},
	}
}

type TextEditorTool struct {
	BaseTool
	MaxCharacters int           `json:"max_characters"`
	CacheControl  *CacheControl `json:"cache_control"`
}

func NewTextEditorTool() TextEditorTool {
	return TextEditorTool{
		BaseTool:      BaseTool{Type: "text_editor_20250728", Name: TEXT_EDITOR},
		MaxCharacters: 10000,
	}
}

type CacheTTL string

const (
	TTL_5m CacheTTL = "5m"
	TTL_1h CacheTTL = "1h"
)

type CacheControl struct {
	Type string   `json:"type,omitempty"`
	TTL  CacheTTL `json:"ttl,omitempty"`
}

type AnthropicMessagesRequest struct {
	Model         Model               `json:"model"`
	Messages      []Message           `json:"messages"`
	MaxTokens     int                 `json:"max_tokens"`
	Container     string              `json:"container,omitempty"`
	MCPServers    []MCPServer         `json:"mcp_servers,omitempty"`
	Metadata      *Metadata           `json:"metadata,omitempty"`
	ServiceTier   string              `json:"service_tier,omitempty"`
	StopSequences []string            `json:"stop_sequences,omitempty"`
	Stream        bool                `json:"stream,omitempty"`
	System        string              `json:"system,omitempty"` //System prompt
	Temperature   float32             `json:"temperature,omitempty"`
	Thinking      *ThinkingData       `json:"thinking,omitempty"`
	ToolChoice    any                 `json:"tool_choice,omitempty"`
	Tools         []AnthropicToolSpec `json:"tools,omitempty"`
	TopK          int                 `json:"top_k,omitempty"`
	TopP          int                 `json:"top_p,omitempty"`
}

type Model string

const (
	SONNET_4 Model = "claude-sonnet-4-20250514"
)
