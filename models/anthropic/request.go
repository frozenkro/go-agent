package anthropic

import (
	"fmt"

	"github.com/frozenkro/goagent/internal/globals"
	"github.com/frozenkro/goagent/internal/sessionmgr"
	"github.com/frozenkro/goagent/internal/tools"
	"github.com/frozenkro/goagent/models/toolschema"
)

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

type AnthropicToolSpec interface {
	GetType() string
	GetName() string
}

type BaseTool struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (t BaseTool) GetType() string {
	return t.Type
}
func (t BaseTool) GetName() string {
	return t.Name
}

type BashTool struct {
	BaseTool
	CacheControl *CacheControl `json:"cache_control"`
}

func NewBashTool() BashTool {
	return BashTool{
		BaseTool: BaseTool{Type: "bash_20250124", Name: globals.ANTH_BASH},
	}
}

type TextEditorTool struct {
	BaseTool
	MaxCharacters int           `json:"max_characters"`
	CacheControl  *CacheControl `json:"cache_control"`
}

func NewTextEditorTool() TextEditorTool {
	return TextEditorTool{
		BaseTool:      BaseTool{Type: "text_editor_20250728", Name: globals.ANTH_TEXT_EDITOR},
		MaxCharacters: 10000,
	}
}

type CustomTool struct {
	BaseTool
	Description  string        `json:"description"`
	InputSchema  InputSchema   `json:"input_schema"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

type InputSchema struct {
	Type       string             `json:"type"` // Always "object"
	Required   []string           `json:"required"`
	Properties *map[string]string `json:"properties,omitempty"`
}

func NewCustomTool(meta *tools.ToolMeta) CustomTool {
	props := map[string]string{}
	for k, v := range meta.Tool.Spec().Parameters.Properties {
		props[k] = v.Type
	}

	return CustomTool{
		BaseTool: BaseTool{
			Type: "custom",
			Name: meta.Name,
		},
		Description: meta.Tool.Spec().Description,
		InputSchema: InputSchema{
			Type:       "object",
			Required:   meta.Tool.Spec().Parameters.Required,
			Properties: &props,
		},
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

type MessagesRequest struct {
	Model         string              `json:"model"`
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

const (
	SONNET_4 string = "claude-sonnet-4-20250514"
)

func (r *MessagesRequest) Init(model, prompt string) error {
	r.Model = model
	// TODO validate?

	r.Messages = append(r.Messages, Message{
		Role: "user",
		Content: []Content{
			TextContent{
				BaseContent: BaseContent{
					Type: TEXT,
				},
				Text: prompt,
			},
		},
	})

	return nil
}

func (r *MessagesRequest) AddTool(meta *tools.ToolMeta) error {
	// Our bash and text editor tools adhere to anthropic-specific specs
	// so we can ignore their defined schema and just tell anthropic they are
	// their custom specs
	if meta.Name == toolschema.BASH {
		r.Tools = append(r.Tools, NewBashTool())
	} else if meta.Name == toolschema.TEXT_EDITOR {
		r.Tools = append(r.Tools, NewTextEditorTool())
	} else {
		r.Tools = append(r.Tools, NewCustomTool(meta))
	}
	return nil
}

func (r *MessagesRequest) SetMaxTokens(val int) error {
	r.MaxTokens = val
	return nil
}

func (r *MessagesRequest) AddStatementGroup(messages []sessionmgr.Statement) error {
	if len(messages) == 0 {
		return nil
	}
	var role Role
	if messages[0].Role == sessionmgr.ASSISTANT {
		role = ASSISTANT
	} else {
		role = USER
	}

	conts := []Content{}
	for _, m := range messages {
		switch m.Type {

		case sessionmgr.TEXT:
			conts = append(conts, TextContent{
				BaseContent: BaseContent{Type: TEXT},
				Text:        m.Text,
			})
		case sessionmgr.TOOL_CALL:
			conts = append(conts, ToolUseContent{
				BaseContent: BaseContent{Type: TOOL_USE},
				Id:          m.ToolCall.Id,
				Name:        m.ToolCall.Name,
				Input:       m.ToolCall.Params,
			})
		case sessionmgr.TOOL_RESPONSE:
			conts = append(conts, ToolResultContent{
				BaseContent: BaseContent{Type: TOOL_RESULT},
				ToolUseId:   m.ToolCall.Id,
				Content:     m.Text,
			})
		case sessionmgr.THINKING:
			conts = append(conts, ThinkingContent{
				BaseContent: BaseContent{Type: THINKING},
				Thinking:    m.Text,
			})
		default:
			return fmt.Errorf("Unsupported content type '%v'\nContent: %v", m.Type, m.Text)
		}
	}

	r.Messages = append(r.Messages, Message{
		Role:    role,
		Content: conts,
	})

	return nil
}
