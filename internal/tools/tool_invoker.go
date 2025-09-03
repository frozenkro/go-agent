package tools

import "github.com/frozenkro/go-agent/models/anthropic"

type ToolInvoker struct {
	ToolMap *ToolMap
}

func NewToolInvoker() ToolInvoker {
	toolMap := InitToolMap()
	return ToolInvoker{
		ToolMap: toolMap,
	}
}

func (t *ToolInvoker) Invoke(toolUseContent anthropic.ToolUseContent) (anthropic.ToolResultContent, error) {
	toolMeta, err := t.ToolMap.ToolMetaByName(anthropic.ToolName(toolUseContent.Name))
	if err != nil {
		return anthropic.ToolResultContent{}, err
	}
	result, err := toolMeta.Tool.Invoke(toolUseContent.Input)

	toolResultContent := anthropic.ToolResultContent{
		BaseContent: anthropic.BaseContent{Type: anthropic.TOOL_RESULT},
		ToolUseId:   toolUseContent.Id,
		Content:     result,
	}
	return toolResultContent, nil
}
