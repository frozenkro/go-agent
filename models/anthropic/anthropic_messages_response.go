package anthropic

import (
	"encoding/json"
)

type MessagesBaseResponse struct {
	Type string `json:"type"`
}

func (r MessagesBaseResponse) GetType() string {
	return r.Type
}

type MessagesResponse struct {
	MessagesBaseResponse
	ID         string     `json:"id"`
	Role       string     `json:"role"`
	Content    []Content  `json:"-"`
	Model      string     `json:"model"`
	StopReason StopReason `json:"stop_reason"`
	Usage      any        `json:"usage"`
	Container  Container  `json:"container,omitempty"`
}

type MessagesErrorResponse struct {
	MessagesBaseResponse
	Error     MessagesError `json:"error"`
	RequestId string        `json:"request_id"`
}

type MessagesError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Custom unmarshaling for the response
func (r *MessagesResponse) UnmarshalJSON(data []byte) error {
	type Alias MessagesResponse
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

type MessagesUsage struct {
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

type StopReason string

const (
	SR_END_TURN      StopReason = "end_turn"
	SR_MAX_TOKENS    StopReason = "max_tokens"
	SR_STOP_SEQUENCE StopReason = "stop_sequence"
	SR_TOOL_USE      StopReason = "tool_use"
	SR_PAUSE_TURN    StopReason = "pause_turn"
	SR_REFUSAL       StopReason = "refusal"
)
