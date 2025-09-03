package anthropic

import (
	"encoding/json"
)

// Updated response struct
type AnthropicMessagesResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Role       string    `json:"role"`
	Content    []Content `json:"-"`
	Model      string    `json:"model"`
	StopReason string    `json:"stop_reason"`
	Usage      any       `json:"usage"`
	Container  Container `json:"container"`
}

// Custom unmarshaling for the response
func (r *AnthropicMessagesResponse) UnmarshalJSON(data []byte) error {
	type Alias AnthropicMessagesResponse
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	contents, err := UnmarshalContents(aux.Content)
	if err != nil {
		return err
	}
	r.Content = contents

	return nil
}

type CacheCreation struct {
	Ephemeral1hInputTokens int `json:"ephemeral_1h_input_tokens"`
	Ephemeral5mInputTokens int `json:"ephemeral_5m_input_tokens"`
}

type ServerToolUse struct {
	WebSearchRequests int `json:"web_search_requests"`
}

type AnthropicMessagesUsage struct {
	CacheCreation            CacheCreation `json:"cache_creation"`
	CacheCreationInputTokens int           `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int           `json:"cache_read_input_tokens"`
	InputTokens              int           `json:"input_tokens"`
	OutputTokens             int           `json:"output_tokens"`
	ServerToolUse            ServerToolUse `json:"server_tool_use"`
	ServiceTier              string        `json:"service_tier"`
}

type Container struct {
	ExpiresAt string `json:"expires_at"`
	Id        string `json:"id"`
}
