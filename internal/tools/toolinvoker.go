package tools

type ToolInvoker struct {
	ToolMap *ToolMap
}

func NewToolInvoker() ToolInvoker {
	toolMap := InitToolMap()
	return ToolInvoker{
		ToolMap: toolMap,
	}
}

// Invoke will invoke any tool in the ToolMap
func (t *ToolInvoker) Invoke(toolName string, input any) (string, error) {
	toolMeta, err := t.ToolMap.ToolMetaByName(toolName)
	if err != nil {
		return "", err
	}
	return toolMeta.Tool.Invoke(input)
}
