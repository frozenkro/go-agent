package tools

type ToolInvoker struct{}

// Invoke will invoke any tool in the ToolMap
func (t *ToolInvoker) Invoke(toolName string, input any) (string, error) {
	toolMeta, err := ToolMapInst.ToolMetaByName(toolName)
	if err != nil {
		return "", err
	}
	return toolMeta.Tool.Invoke(input)
}
