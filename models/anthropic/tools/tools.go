package tools

type BashToolInput struct {
	Command string `json:"command"`
	Restart bool   `json:"restart"`
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
