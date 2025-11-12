package sessionmgr

type Statement struct {
	Role       Role
	Type       Type
	Text       string
	ToolCall   ToolCall
	ToolCallId string
}

type Role int

const (
	SYSTEM Role = iota + 1
	USER
	ASSISTANT
	TOOL
)

type Type int

const (
	TEXT Type = iota + 1
	TOOL_CALL
	TOOL_RESPONSE
	THINKING
)

type ToolCall struct {
	Name   string
	Params map[string]any
}
