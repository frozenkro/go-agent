package anthropic

import (
	"fmt"
	"log"
)

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

type ToolName string

const (
	BASH        ToolName = "bash"
	TEXT_EDITOR ToolName = "str_replace_based_edit_tool"
)

type Tool interface {
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
	CacheControl CacheControl `json:"cache_control"`
}

func NewBashTool() BashTool {
	return BashTool{
		BaseTool:     BaseTool{Type: "bash_20250124", Name: BASH},
		CacheControl: CacheControl{},
	}
}

type TextEditorTool struct {
	BaseTool
	CacheControl  CacheControl `json:"cache_control"`
	MaxCharacters int          `json:"max_characters"`
}

func NewTextEditorTool() TextEditorTool {
	return TextEditorTool{
		BaseTool:      BaseTool{Type: "text_editor_20250728", Name: TEXT_EDITOR},
		MaxCharacters: 10000,
		CacheControl:  CacheControl{},
	}
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

type Model string

const (
	SONNET_4 Model = "claude-sonnet-4-20250514"
)

type AnthropicMessagesRequestOption func(*AnthropicMessagesRequest)

func WithTools(tools ...ToolName) AnthropicMessagesRequestOption {
	return func(a *AnthropicMessagesRequest) {

		for _, toolName := range tools {
			tool, err := ToolByName(toolName)

			if err == nil {
				a.Tools = append(a.Tools, tool)
			} else {
				log.Printf(err.Error())
			}
		}
	}
}

func ToolByName(name ToolName) (Tool, error) {
	switch name {
	case BASH:
		return NewBashTool(), nil
	case TEXT_EDITOR:
		return NewTextEditorTool(), nil
	default:
		return nil, fmt.Errorf("No tool found with name %v", name)
	}
}
