package anthropic

import (
	"encoding/json"
	"fmt"
)

type ContentTypes string

const (
	TEXT                       ContentTypes = "text"
	THINKING                   ContentTypes = "thinking"
	REDACTED_THINKING          ContentTypes = "redacted_thinking"
	TOOL_USE                   ContentTypes = "tool_use"
	SERVER_TOOL_USE            ContentTypes = "server_tool_use"
	WEB_SEARCH_TOOL_RESULT     ContentTypes = "web_search_tool_result"
	CODE_EXECUTION_TOOL_RESULT ContentTypes = "code_execution_tool_result"
	MCP_TOOL_USE               ContentTypes = "mcp_tool_use"
	MCP_TOOL_RESULT            ContentTypes = "mcp_tool_result"
	CONTAINER_UPLOAD           ContentTypes = "container_upload"
)

// Content interface that all content types implement
type Content interface {
	GetType() ContentTypes
}

// Base struct for common fields
type BaseContent struct {
	Type ContentTypes `json:"type"`
}

func (b BaseContent) GetType() ContentTypes {
	return b.Type
}

// Specific content type structs
type TextContent struct {
	BaseContent
	Text string `json:"text"`
}

type ThinkingContent struct {
	BaseContent
	Signature string `json:"signature"`
	Thinking  string `json:"thinking"`
}

type RedactedThinkingContent struct {
	BaseContent
	Data string `json:"data"`
}

type ToolUseContent struct {
	BaseContent
	Id    string `json:"id"`
	Name  string `json:"name"`
	Input any    `json:"input"`
}

type WebSearchToolResultContent struct {
	BaseContent
	ToolUseId string                            `json:"tool_use_id"`
	Content   WebSearchToolResultContentContent `json:"content"`
}

type WebSearchToolResultContentContent struct {
	BaseContent
	ErrorCode string `json:"error_code"`
}

type CodeExecutionToolResultContent struct {
	BaseContent
	Content   Content `json:"-"`
	ToolUseId string  `json:"tool_use_id"`
}

type CodeExecutionToolResultContentContent struct {
	BaseContent
	Content    CodeExecutionToolResultContentContentContent `json:"content"`
	ReturnCode int                                          `json:"return_code"`
	StdErr     string                                       `json:"stderr"`
	StdOut     string                                       `json:"stdout"`
}

type CodeExecutionToolResultContentErr struct {
	BaseContent
	ErrorCode string `json:"error_code"`
}

type CodeExecutionToolResultContentContentContent struct {
	BaseContent
	FileId string `json:"file_id"`
}

type MCPToolUseContent struct {
	BaseContent
	Id         string `json:"id"`
	Name       string `json:"name"`
	Input      any    `json:"input"`
	ServerName string `json:"server_name"`
}

type MCPToolResultContent struct {
	BaseContent
	Content   string `json:"content"`
	IsError   bool   `json:"is_error"`
	ToolUseId string `json:"tool_use_id"`
}

type ContainerUploadContent struct {
	BaseContent
	FileId string `json:"file_id"`
}

func UnmarshalContents(data []byte) ([]Content, error) {
	var rawContents []json.RawMessage
	if err := json.Unmarshal(data, &rawContents); err != nil {
		return nil, err
	}

	contents := make([]Content, len(rawContents))
	for i, raw := range rawContents {
		var base BaseContent
		if err := json.Unmarshal(raw, &base); err != nil {
			return nil, err
		}

		var content Content
		switch base.Type {
		case TEXT:
			content = &TextContent{}
		case THINKING:
			content = &ThinkingContent{}
		case REDACTED_THINKING:
			content = &RedactedThinkingContent{}
		case TOOL_USE:
			content = &ToolUseContent{}
		case SERVER_TOOL_USE:
			content = &ToolUseContent{}
		case WEB_SEARCH_TOOL_RESULT:
			content = &WebSearchToolResultContent{}
		case CODE_EXECUTION_TOOL_RESULT:
			content = &CodeExecutionToolResultContent{}
		case MCP_TOOL_USE:
			content = &MCPToolUseContent{}
		case MCP_TOOL_RESULT:
			content = &MCPToolResultContent{}
		case CONTAINER_UPLOAD:
			content = &ContainerUploadContent{}
		default:
			return nil, fmt.Errorf("unknown content type: %s", base.Type)
		}

		if err := json.Unmarshal(raw, content); err != nil {
			return nil, err
		}
		contents[i] = content
	}

	return contents, nil
}

// Updated response struct
type AnthropicMessagesResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Role       string    `json:"role"`
	Content    []Content `json:"-"`
	Model      string    `json:"model"`
	StopReason string    `json:"stop_reason"`
	Usage      any       `json:"usage"`
	Container  Container `json:"container"`
}

// Custom unmarshaling for the response
func (r *AnthropicMessagesResponse) UnmarshalJSON(data []byte) error {
	type Alias AnthropicMessagesResponse
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	contents, err := UnmarshalContents(aux.Content)
	if err != nil {
		return err
	}
	r.Content = contents

	return nil
}

type CacheCreation struct {
	Ephemeral1hInputTokens int `json:"ephemeral_1h_input_tokens"`
	Ephemeral5mInputTokens int `json:"ephemeral_5m_input_tokens"`
}

type ServerToolUse struct {
	WebSearchRequests int `json:"web_search_requests"`
}

type AnthropicMessagesUsage struct {
	CacheCreation            CacheCreation `json:"cache_creation"`
	CacheCreationInputTokens int           `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int           `json:"cache_read_input_tokens"`
	InputTokens              int           `json:"input_tokens"`
	OutputTokens             int           `json:"output_tokens"`
	ServerToolUse            ServerToolUse `json:"server_tool_use"`
	ServiceTier              string        `json:"service_tier"`
}

type Container struct {
	ExpiresAt string `json:"expires_at"`
	Id        string `json:"id"`
}
