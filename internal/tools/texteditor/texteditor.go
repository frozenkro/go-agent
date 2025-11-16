// Package texteditor provides a text editor tool to be invoked by llm agents
// Current implementation is specific to anthropic spec:
// https://anthropic.mintlify.app/en/docs/agents-and-tools/tool-use/text-editor-tool#view

package texteditor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frozenkro/goagent/models/toolschema"
	"github.com/mitchellh/mapstructure"
)

type TextEditorTool struct {
	w TextEditorWorker
}

func NewTextEditorTool() *TextEditorTool {
	return &TextEditorTool{
		w: TextEditorWorkerImpl{},
	}
}

type TextEditorWorkerImpl struct{}

func (t TextEditorTool) Spec() toolschema.ToolSpec {
	return toolschema.ToolSpec{
		Name:        toolschema.TEXT_EDITOR,
		Description: "View or modify text files",
		Parameters:  toolschema.TextEditorToolParams,
	}
}

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
func (w TextEditorWorkerImpl) HandleView(params any) (string, error) {
	var input toolschema.TextEditorToolInputView
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal params during type conversion: %w", err)
	}
	err = json.Unmarshal(jsonParams, &input)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal params during type conversion: %w", err)
	}

	input.Path = sanitizeRelativePath(input.Path)

	// Handle directory view - basically just the results of `ls`
	fileInfo, err := os.Stat(input.Path)
	if err != nil {
		return "", fmt.Errorf("Failed to get stats for path '%v': %w", input.Path, err)
	}
	if fileInfo.IsDir() {
		entries, err := os.ReadDir(input.Path)
		if err != nil {
			return "", fmt.Errorf("Failed to list directory contents for dir '%v': %w", input.Path, err)
		}

		result := strings.Builder{}
		for _, v := range entries {
			var entry string
			if v.IsDir() {
				// Indicate to the requesting agent that this entry is a dir
				entry = fmt.Sprintf("%v%c", v.Name(), '/')
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
	for i := start; continueReadFile(i, end); i++ {
		scanning = scanner.Scan()
		if !scanning {
			break
		}

		line := scanner.Text()
		fmt.Fprintf(&result, "%v. %v\n", i+1, line)

	}

	return result.String(), nil
}

func continueReadFile(i int, end int) bool {
	// Keep reading until end of file
	if end == -1 {
		return true
	}

	return i <= end
}

// Handle request to replace a string within a file
func (w TextEditorWorkerImpl) HandleStrReplace(params any) (string, error) {
	var input toolschema.TextEditorToolInputStrReplace
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal params during type conversion: %w", err)
	}
	err = json.Unmarshal(jsonParams, &input)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal params during type conversion: %w", err)
	}

	input.Path = sanitizeRelativePath(input.Path)

	fileContent, err := getFileContent(input.Path)
	if err != nil {
		return "", fmt.Errorf("Failed to get file content at %v: %w", input.Path, err)
	}

	newContent := strings.ReplaceAll(fileContent, input.OldStr, input.NewStr)

	err = os.WriteFile(input.Path, []byte(newContent), 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write file %v: %w", input.Path, err)
	}

	return newContent, nil
}

// Handle request to create a file
func (w TextEditorWorkerImpl) HandleCreate(params any) (string, error) {
	var input toolschema.TextEditorToolInputCreate
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal params during type conversion: %w", err)
	}
	err = json.Unmarshal(jsonParams, &input)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal params during type conversion: %w", err)
	}

	input.Path = sanitizeRelativePath(input.Path)

	path := filepath.Dir(input.Path)
	if err = os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("Failed to create directory %v: %w", path, err)
	}

	err = os.WriteFile(input.Path, []byte(input.FileText), 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write file %v: %w", input.Path, err)
	}

	return input.FileText, nil
}

// Handle request to insert a string into a file at a specified line number
func (w TextEditorWorkerImpl) HandleInsert(params any) (string, error) {
	var input toolschema.TextEditorToolInputInsert
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal params during type conversion: %w", err)
	}
	err = json.Unmarshal(jsonParams, &input)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal params during type conversion: %w", err)
	}

	input.Path = sanitizeRelativePath(input.Path)

	file, err := os.Open(input.Path)
	if err != nil {
		return "", fmt.Errorf("Failed to open file at path '%v': %w", input.Path, err)
	}
	// defer file.Close() // Manually closed after loop
	scanner := bufio.NewScanner(file)

	b := strings.Builder{}
	scanning := true
	lineNumber := 0
	for {
		lineNumber += 1

		if lineNumber == input.InsertLine {
			_, err := b.Write([]byte(input.NewStr))
			b.Write([]byte("\n"))
			if err != nil {
				return "", fmt.Errorf("Failed to write to internal buffer: %w", err)
			}

			continue
		}

		scanning = scanner.Scan()
		if !scanning {
			break
		}

		line := scanner.Bytes()

		b.Write(line)
		b.Write([]byte("\n"))

	}
	file.Close()

	err = os.WriteFile(input.Path, []byte(b.String()), 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write to %v: %w", input.Path, err)
	}

	return b.String(), nil
}

// Unsupported on claude 4, skipping implementation for the time being
// Note for implementation - we will need to keep a stack of edits per-file per-session
func (w TextEditorWorkerImpl) HandleUndoEdit(params any) (string, error) {
	// var input toolschema.TextEditorToolInputUndoEdit
	// err := mapstructure.Decode(params, &input)
	// if err != nil {
	// 	return "", fmt.Errorf("Unable to parse invoke params for TextEditorTool: '%v'", params)
	// }

	return "", fmt.Errorf("UndoEdit operation unsupported")
}

func getViewRange(inputViewRange string) (int, int) {
	if inputViewRange == "" {
		return 0, -1
	}

	rangeArr := make([]int, 2)
	err := json.Unmarshal([]byte(inputViewRange), &rangeArr)
	if err != nil {
		return 0, -1
	}

	if len(rangeArr) != 2 {
		return 0, -1
	}

	start := rangeArr[0]
	end := rangeArr[1]

	// Requested lines are 1-indexed
	if start > 0 {
		start--
	}
	if end != -1 {
		end--
	}
	return start, end
}

func getFileContent(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Failed to open file at path '%v': %w", path, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	b := strings.Builder{}
	scanning := true
	for {
		scanning = scanner.Scan()
		line := scanner.Bytes()

		b.Write(line)
		b.Write([]byte("\n"))

		if !scanning {
			break
		}
	}

	return b.String(), nil
}

// sanitizeRelativePath ensures "myFile", "./myFile", and "/myFile" all result in a write to "./myFile"
// This is because agents should never operate outside of cwd
func sanitizeRelativePath(path string) string {
	cleanInputPath := strings.TrimPrefix(path, "/")
	return fmt.Sprintf("./%v", cleanInputPath)
}
