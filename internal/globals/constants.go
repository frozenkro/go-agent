package globals

type Provider int

const (
	ANTHROPIC Provider = iota + 1
	OLLAMA
)

const (
	ANTHROPIC_MESSAGES_URL = "https://api.anthropic.com/v1/messages"
	OLLAMA_DEFAULT_URL     = "http://localhost:11434/api/chat"
)
