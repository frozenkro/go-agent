package toolschema

type BashToolInput struct {
	Command string `json:"command"`
	Restart bool   `json:"restart"`
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

type TextEditorToolInputUndoEdit struct {
	BaseTextEditorToolInput
}
