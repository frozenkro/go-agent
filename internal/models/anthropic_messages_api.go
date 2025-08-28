package models

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
)

type ContentContentContent struct {
	FileId string `json:"file_id"`
	Type   string `json:"type"`
}

type ContentContent struct {
	ErrorCode  string                `json:"error_code"`
	Type       string                `json:"type"`
	Content    ContentContentContent `json:"content"`
	ReturnCode int                   `json:"return_code"`
	StdErr     string                `json:"stderr"`
	StdOut     string                `json:"stdout"`
}
type AnthropicMessagesContent struct {
	Text       string         `json:"text"`
	Type       string         `json:"type"`
	Id         string         `json:"id"`
	Input      string         `json:"input"`
	Name       string         `json:"name"`
	Data       string         `json:"data"`
	Thinking   string         `json:"thinking"`
	Signature  string         `json:"signature"`
	Content    ContentContent `json:"content"` // TODO - this can be a string on mcp_tool_result, we need to map out each object individually
	ToolUseId  string         `json:"tool_use_id"`
	ServerName string         `json:"server_name"`
	IsError    bool           `json:"is_error"`
}

type AnthropicMessagesResponse struct {
	ID      string                     `json:"id"`
	Type    string                     `json:"type"`
	Role    string                     `json:"role"`
	Content []AnthropicMessagesContent `json:"content"`
}
