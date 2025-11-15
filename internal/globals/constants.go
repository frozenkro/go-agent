package globals

type Provider int

const (
	ANTHROPIC Provider = iota + 1
	OLLAMA
)

// Anthropic-specific tool names
const (
	ANTH_BASH        string = "bash"
	ANTH_TEXT_EDITOR string = "str_replace_based_edit_tool"
)
