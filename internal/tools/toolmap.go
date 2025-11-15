package tools

import (
	"fmt"

	"github.com/frozenkro/goagent/internal/globals"
	"github.com/frozenkro/goagent/internal/tools/bash"
	"github.com/frozenkro/goagent/internal/tools/texteditor"
	"github.com/frozenkro/goagent/models/toolschema"
)

type ToolMeta struct {
	Name string
	Tool Tool
}

type ToolMap struct {
	Map map[string]ToolMeta
}

var ToolMapInst ToolMap

func InitToolMap(provider globals.Provider) {
	toolNameMap := make(map[string]ToolMeta)

	var bashToolName string
	var textEditorToolName string

	if provider == globals.ANTHROPIC {
		bashToolName = globals.ANTH_BASH
		textEditorToolName = globals.ANTH_TEXT_EDITOR
	} else {
		bashToolName = toolschema.BASH
		textEditorToolName = toolschema.TEXT_EDITOR
	}

	toolNameMap[toolschema.BASH] = ToolMeta{
		Name: bashToolName,
		Tool: bash.BashTool{},
	}
	toolNameMap[toolschema.TEXT_EDITOR] = ToolMeta{
		Name: textEditorToolName,
		Tool: texteditor.NewTextEditorTool(),
	}

	ToolMapInst = ToolMap{
		Map: toolNameMap,
	}
}

func (t *ToolMap) ToolMetaByName(name string) (*ToolMeta, error) {
	meta, ok := t.Map[name]
	if !ok {
		return nil, fmt.Errorf("No tool found with name %v", name)
	}

	return &meta, nil
}
