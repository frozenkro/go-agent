// Package texteditor provides a text editor tool to be invoked by llm agents
// Current implementation is specific to anthropic spec:
// https://anthropic.mintlify.app/en/docs/agents-and-tools/tool-use/text-editor-tool#view

package texteditor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	toolschema "github.com/frozenkro/go-agent/models/anthropic/tool_schema"
	"github.com/mitchellh/mapstructure"
)

type TextEditorTool struct {
	w TextEditorWorker
}

func NewTextEditorTool() *TextEditorTool {
	return &TextEditorTool{
		w: TextEditorWorker{},
	}
}

type TextEditorWorker struct{}

func (t TextEditorTool) Invoke(params any) (string, error) {
	var base toolschema.BaseTextEditorToolInput
	err := mapstructure.Decode(params, &base)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	}

	switch base.Command {
	case "view":
		return t.w.HandleView(params)
	case "str_replace":
		return t.w.HandleStrReplace(params)
	case "create":
		return t.w.HandleCreate(params)
	case "insert":
		return t.w.HandleInsert(params)
	case "undo_edit":
		return t.w.HandleUndoEdit(params)
	default:
		return "", fmt.Errorf("Unrecognized 'command' in tool invocation parameters: '%v'", base.Command)
	}
}

// Handle request to view a file or directory
func (w TextEditorWorker) HandleView(params any) (string, error) {
	var input toolschema.TextEditorToolInputView
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", err
	}

	// Handle directory view
	dirChar := input.Path[len(input.Path)-1]
	if dirChar == '/' || dirChar == '\\' {
		entries, err := os.ReadDir(input.Path)
		if err != nil {
			return "", fmt.Errorf("Failed to list directory contents for dir '%v': %w", input.Path, err)
		}

		result := strings.Builder{}
		for _, v := range entries {
			var entry string
			if v.IsDir() {
				// Indicate to the requesting agent that the entry is a dir
				entry = fmt.Sprintf("%v%c", v.Name(), dirChar)
			} else {
				entry = v.Name()
			}
			fmt.Fprintf(&result, "%v\n", entry)
		}

		return result.String(), nil
	}

	// Handle file view
	file, err := os.Open(input.Path)
	if err != nil {
		return "", fmt.Errorf("Failed to open file at path '%v': %w", input.Path, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	start, end := getViewRange(input.ViewRange)

	for i := 0; i < start; i++ {
		// advance reader to first line we want to read
		scanner.Scan()
	}

	scanning := true
	result := strings.Builder{}
	for i := start; i == end+1; i++ {
		scanning = scanner.Scan()
		if !scanning {
			break
		}

		line := scanner.Text()
		fmt.Fprintf(&result, "%c. %v\n", i+1, line)
	}

	return result.String(), nil
}

// Handle request to replace a string within a file
func (w TextEditorWorker) HandleStrReplace(params any) (string, error) {
	var input toolschema.TextEditorToolInputStrReplace
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	}

	// TODO

	return "", nil
}

// Handle request to create a file
func (w TextEditorWorker) HandleCreate(params any) (string, error) {
	var input toolschema.TextEditorToolInputCreate
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	}

	// TODO

	return "", nil
}

// Handle request to insert a string into a file at a specified line number
func (w TextEditorWorker) HandleInsert(params any) (string, error) {
	var input toolschema.TextEditorToolInputInsert
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	}

	// TODO

	return "", nil
}

// Unsupported on claude 4, skipping implementation for the time being
// Note for implementation - we will need to keep a stack of edits per-file per-session
func (w TextEditorWorker) HandleUndoEdit(params any) (string, error) {
	var input toolschema.TextEditorToolInputUndoEdit
	err := mapstructure.Decode(params, &input)
	if err != nil {
		return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	}

	// TODO

	return "", nil
}

func getViewRange(inputViewRange string) (int, int) {
	if inputViewRange == "" {
		return 0, -1
	}

	rangeArr := make([]int, 2)
	err := json.Unmarshal([]byte(inputViewRange), rangeArr)
	if err != nil {
		return 0, -1
	}

	// Requested lines are 1-indexed
	return rangeArr[0] - 1, rangeArr[1] - 1
}
