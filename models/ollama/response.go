package ollama

import (
	"encoding/json"
	"strconv"

	"github.com/frozenkro/goagent/internal/sessionmgr"
)

type OllamaResponse struct {
	Model              string  `json:"model"`
	CreatedAt          string  `json:"created_at"`
	Message            Message `json:"message"`
	Done               bool    `json:"done"`
	DoneReason         string  `json:"done_reason,omitempty"`
	TotalDuration      int     `json:"total_duration,omitempty"`
	LoadDuration       int     `json:"load_duration,omitempty"`
	PromptEvalCount    int     `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int     `json:"prompt_eval_duration,omitempty"`
	EvalCount          int     `json:"eval_count,omitempty"`
	EvalDuration       int     `json:"eval_duration,omitempty"`
}

func (r *OllamaResponse) IsComplete() bool {
	return r.Done
}

func (r *OllamaResponse) GetStatementGroup() ([]sessionmgr.Statement, error) {
	sts := []sessionmgr.Statement{}

	sts = append(sts, sessionmgr.Statement{
		Role: sessionmgr.ASSISTANT,
		Type: sessionmgr.TEXT,
		Text: r.Message.Content,
	})

	for i, v := range r.Message.ToolCalls {
		sts = append(sts, sessionmgr.Statement{
			Role: sessionmgr.ASSISTANT,
			Type: sessionmgr.TOOL_CALL,
			ToolCall: sessionmgr.ToolCall{
				Name:   v.Function.Name,
				Params: v.Function.Arguments,
				Id:     strconv.Itoa(i),
			},
		})
	}

	return sts, nil
}

func (r *OllamaResponse) Init(data []byte) error {
	// TODO -- is there any ollama error response schema documented? Can't find it
	return json.Unmarshal(data, r)
}
