package tools

import (
	"fmt"

	"github.com/frozenkro/goagent/internal/tools/bash"
	"github.com/frozenkro/goagent/internal/tools/texteditor"
	"github.com/frozenkro/goagent/models/anthropic"
)

type ToolMeta struct {
	Name anthropic.ToolName
	Spec anthropic.AnthropicToolSpec
	Tool Tool
}

type ToolMap struct {
	Map map[anthropic.ToolName]ToolMeta
}

func InitToolMap() *ToolMap {
	toolNameMap := make(map[anthropic.ToolName]ToolMeta)

	toolNameMap[anthropic.BASH] = ToolMeta{
		Name: anthropic.BASH,
		Spec: anthropic.NewBashTool(),
		Tool: bash.BashTool{},
	}
	toolNameMap[anthropic.TEXT_EDITOR] = ToolMeta{
		Name: anthropic.TEXT_EDITOR,
		Spec: anthropic.NewTextEditorTool(),
		Tool: texteditor.NewTextEditorTool(),
	}

	return &ToolMap{
		Map: toolNameMap,
	}
}

func (t *ToolMap) ToolMetaByName(name anthropic.ToolName) (*ToolMeta, error) {
	meta, ok := t.Map[name]
	if !ok {
		return nil, fmt.Errorf("No tool found with name %v", name)
	}

	return &meta, nil
}
