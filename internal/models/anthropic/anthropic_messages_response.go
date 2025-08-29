package anthropic

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type ToolConfiguration struct {
	AllowedTools []string `json:"allowed_tools"`
	Enabled      bool     `json:"enabled"`
}

type MCPServer struct {
	Name               string            `json:"name"`
	Type               string            `json:"type"`
	Url                string            `json:"url"`
	AuthorizationToken string            `json:"authorization_token"`
	ToolConfiguration  ToolConfiguration `json:"tool_configuration"`
}

type Metadata struct {
	UserId string `json:"user_id"`
}

type ThinkingData struct {
	BudgetTokens int    `json:"budget_tokens"`
	Type         string `json:"type"`
}

type Tool interface {
	GetType() string
	GetName() string
}

type BaseTool struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (t *BaseTool) GetType() string {
	return t.Type
}
func (t *BaseTool) GetName() string {
	return t.Name
}

type BashTool struct {
	BaseTool
	CacheControl CacheControl `json:"cache_control"`
}

type TextEditorTool struct {
	BaseTool
	CacheControl  CacheControl `json:"cache_control"`
	MaxCharacters int          `json:"max_characters"`
}

type CacheTTL string

const (
	TTL_5m CacheTTL = "5m"
	TTL_1h CacheTTL = "1h"
)

type CacheControl struct {
	Type string   `json:"type"`
	TTL  CacheTTL `json:"ttl"`
}

type AnthropicMessagesRequest struct {
	Model         string       `json:"model"`
	Messages      []Message    `json:"messages"`
	MaxTokens     int          `json:"max_tokens"`
	Container     string       `json:"container"`
	MCPServers    []MCPServer  `json:"mcp_servers"`
	Metadata      Metadata     `json:"metadata"`
	ServiceTier   string       `json:"service_tier"`
	StopSequences []string     `json:"stop_sequences"`
	Stream        bool         `json:"stream"`
	System        string       `json:"system"` //System prompt
	Temperature   float32      `json:"temperature"`
	Thinking      ThinkingData `json:"thinking"`
	ToolChoice    any          `json:"tool_choice"`
	Tools         []Tool       `json:"tools"`
	TopK          int          `json:"top_k"`
	TopP          int          `json:"top_p"`
}
