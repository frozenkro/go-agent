package toolschema

type BashToolInput struct {
	Command string `json:"command"`
	Restart bool   `json:"restart"`
}

var BashToolParams = ToolParams{
	Required: []string{"command"},
	Properties: map[string]ToolProperty{
		"command": {Type: "string", Description: "Bash command to be executed"},
		"restart": {Type: "boolean", Description: "Restart persistent bash session"},
	},
}

type BaseTextEditorToolInput struct {
	Command string `json:"command"`
	Path    string `json:"path"`
}

type TextEditorToolInputView struct {
	BaseTextEditorToolInput
	ViewRange string `json:"view_range,omitempty"`
}

type TextEditorToolInputStrReplace struct {
	BaseTextEditorToolInput
	OldStr string `json:"old_str"`
	NewStr string `json:"new_str"`
}

type TextEditorToolInputCreate struct {
	BaseTextEditorToolInput
	FileText string `json:"file_text"`
}

type TextEditorToolInputInsert struct {
	BaseTextEditorToolInput
	InsertLine int    `json:"insert_line"`
	NewStr     string `json:"new_str"`
}

// Not currently supported on claude 4, or by this project
type TextEditorToolInputUndoEdit struct {
	BaseTextEditorToolInput
}

var TextEditorToolParams = ToolParams{
	Required: []string{"command", "path"},
	Properties: map[string]ToolProperty{
		"command":     {Type: "string", Description: "Allowed values: ['view', 'str_replace', 'create', 'insert']."},
		"path":        {Type: "string", Description: "Relative path to the file to be viewed, created, or edited"},
		"view_range":  {Type: "string", Description: "For commands: ['view']. An array of two integers specifying the start and end line numbers to view. Line numbers are 1-indexed, and -1 for the end line means read to the end of the file. This parameter only applies when viewing files, not directories. Ex: [1, 50]"},
		"old_str":     {Type: "string", Description: "For commands: ['str_replace']. The text to replace (must match exactly, including whitespace and indentation)."},
		"new_str":     {Type: "string", Description: "For commands: ['str_replace', 'insert']. New text to be inserted."},
		"file_text":   {Type: "string", Description: "For commands: ['create']. The content of the file to be created."},
		"insert_line": {Type: "string", Description: "For commands: ['insert']. Line number to begin insert."},
	},
}
