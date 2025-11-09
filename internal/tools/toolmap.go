package tools

import (
	"fmt"

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

func InitToolMap() *ToolMap {
	toolNameMap := make(map[string]ToolMeta)

	toolNameMap[toolschema.BASH] = ToolMeta{
		Name: toolschema.BASH,
		Tool: bash.BashTool{},
	}
	toolNameMap[toolschema.TEXT_EDITOR] = ToolMeta{
		Name: toolschema.TEXT_EDITOR,
		Tool: texteditor.NewTextEditorTool(),
	}

	return &ToolMap{
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
