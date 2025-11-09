package toolschema

type ToolSpec struct {
	Name        string
	Parameters  ToolParams
	Description string
}

type ToolParams struct {
	Required   []string
	Properties map[string]ToolProperty
}

type ToolProperty struct {
	Type        string
	Description string
}
