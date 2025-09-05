package tools

import (
	"github.com/frozenkro/go-agent/models/anthropic/content"
	toolModels "github.com/frozenkro/go-agent/models/anthropic/tools"
)

type ToolInvoker struct {
	ToolMap *ToolMap
}

func NewToolInvoker() ToolInvoker {
	toolMap := InitToolMap()
	return ToolInvoker{
		ToolMap: toolMap,
	}
}

func (t *ToolInvoker) Invoke(toolUseContent content.ToolUseContent) (content.ToolResultContent, error) {
	toolMeta, err := t.ToolMap.ToolMetaByName(toolModels.ToolName(toolUseContent.Name))
	if err != nil {
		return content.ToolResultContent{}, err
	}
	result, err := toolMeta.Tool.Invoke(toolUseContent.Input)
	if err != nil {
		return content.ToolResultContent{}, err
	}

	toolResultContent := content.ToolResultContent{
		BaseContent: content.BaseContent{Type: content.TOOL_RESULT},
		ToolUseId:   toolUseContent.Id,
		Content:     result,
	}
	return toolResultContent, nil
}
