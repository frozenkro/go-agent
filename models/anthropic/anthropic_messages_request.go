package anthropic

import (
	"github.com/frozenkro/go-agent/models/anthropic/content"
	"github.com/frozenkro/go-agent/models/anthropic/tools"
)

type Message struct {
	Role    Role              `json:"role"`
	Content []content.Content `json:"content"`
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

type AnthropicMessagesRequest struct {
	Model         Model                     `json:"model"`
	Messages      []Message                 `json:"messages"`
	MaxTokens     int                       `json:"max_tokens"`
	Container     string                    `json:"container,omitempty"`
	MCPServers    []MCPServer               `json:"mcp_servers,omitempty"`
	Metadata      *Metadata                 `json:"metadata,omitempty"`
	ServiceTier   string                    `json:"service_tier,omitempty"`
	StopSequences []string                  `json:"stop_sequences,omitempty"`
	Stream        bool                      `json:"stream,omitempty"`
	System        string                    `json:"system,omitempty"` //System prompt
	Temperature   float32                   `json:"temperature,omitempty"`
	Thinking      *ThinkingData             `json:"thinking,omitempty"`
	ToolChoice    any                       `json:"tool_choice,omitempty"`
	Tools         []tools.AnthropicToolSpec `json:"tools,omitempty"`
	TopK          int                       `json:"top_k,omitempty"`
	TopP          int                       `json:"top_p,omitempty"`
}

type Model string

const (
	SONNET_4 Model = "claude-sonnet-4-20250514"
)
