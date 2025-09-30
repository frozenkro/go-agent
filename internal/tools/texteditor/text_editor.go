package texteditor

import (
	"fmt"

	toolschema "github.com/frozenkro/go-agent/models/anthropic/tool_schema"
	"github.com/mitchellh/mapstructure"
)

type TextEditorTool struct {
	te TextEditor
}

// maybe have this be a separate thing from TextEditorTool? easier to test
type TextEditor struct{}

func (t TextEditorTool) Invoke(params any) (string, error) {
	var base toolschema.BaseTextEditorToolInput
	err := mapstructure.Decode(params, &base)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	}

	switch base.Command {
	case "view":
		return t.HandleView(params)
	case "str_replace":
		return t.HandleStrReplace(params)
	case "create":
		return t.HandleCreate(params)
	case "insert":
		return t.HandleInsert(params)
	case "undo_edit":
		return t.HandleUndoEdit(params)
	default:
		return "", fmt.Errorf("Unrecognized 'command' in tool invocation parameters: '%v'", base.Command)
	}
}

func (t TextEditorTool) HandleView(params any) (string, error) {
	var input toolschema.TextEditorToolInputView
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", err
	}

	// TODO

	return "", nil
}
func (t TextEditorTool) HandleStrReplace(params any) (string, error) {
	var input toolschema.TextEditorToolInputStrReplace
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", err
	}

	// TODO

	return "", nil
}
func (t TextEditorTool) HandleCreate(params any) (string, error) {
	var input toolschema.TextEditorToolInputCreate
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", err
	}

	// TODO

	return "", nil
}
func (t TextEditorTool) HandleInsert(params any) (string, error) {
	var input toolschema.TextEditorToolInputInsert
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", err
	}

	// TODO

	return "", nil
}
func (t TextEditorTool) HandleUndoEdit(params any) (string, error) {
	var input toolschema.TextEditorToolInputUndoEdit
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", err
	}

	// TODO

	return "", nil
}
