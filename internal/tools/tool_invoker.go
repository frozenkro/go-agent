package tools

import "github.com/frozenkro/go-agent/internal/models/anthropic"

type ToolInvoker struct {
	ToolMap *ToolMap
}

func NewToolInvoker() ToolInvoker {
	toolMap := InitToolMap()
	return ToolInvoker{
		ToolMap: toolMap,
	}
}

func (t *ToolInvoker) Invoke(toolName anthropic.ToolName, params any) (string, error) {
	toolMeta, err := t.ToolMap.ToolMetaByName(toolName)
	if err != nil {
		return "", err
	}
	return toolMeta.Tool.Invoke(params)
}
