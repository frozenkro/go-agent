package tools

import (
	"fmt"

	"github.com/frozenkro/go-agent/internal/tools/bash"
	"github.com/frozenkro/go-agent/models/anthropic"
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
		Tool: nil, //TODO
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
