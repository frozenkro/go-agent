package tools

import (
	"fmt"

	"github.com/frozenkro/go-agent/internal/tools/bash"
	toolModels "github.com/frozenkro/go-agent/models/anthropic/tools"
)

type ToolMeta struct {
	Name toolModels.ToolName
	Spec toolModels.AnthropicToolSpec
	Tool Tool
}

type ToolMap struct {
	Map map[toolModels.ToolName]ToolMeta
}

func InitToolMap() *ToolMap {
	toolNameMap := make(map[toolModels.ToolName]ToolMeta)

	toolNameMap[toolModels.BASH] = ToolMeta{
		Name: toolModels.BASH,
		Spec: toolModels.NewBashTool(),
		Tool: bash.BashTool{},
	}
	toolNameMap[toolModels.TEXT_EDITOR] = ToolMeta{
		Name: toolModels.TEXT_EDITOR,
		Spec: toolModels.NewTextEditorTool(),
		Tool: nil, //TODO
	}

	return &ToolMap{
		Map: toolNameMap,
	}
}

func (t *ToolMap) ToolMetaByName(name toolModels.ToolName) (*ToolMeta, error) {
	meta, ok := t.Map[name]
	if !ok {
		return nil, fmt.Errorf("No tool found with name %v", name)
	}

	return &meta, nil
}
