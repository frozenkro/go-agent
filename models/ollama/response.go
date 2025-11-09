package ollama

import "github.com/frozenkro/goagent/internal/sessionmgr"

type OllamaResponse struct {
	Model              string  `json:"model"`
	CreatedAt          string  `json:"created_at"`
	Message            Message `json:"message"`
	Done               bool    `json:"done"`
	DoneReason         bool    `json:"done_reason"`
	TotalDuration      int     `json:"total_duration"`
	LoadDuration       int     `json:"load_duration"`
	PromptEvalCount    int     `json:"prompt_eval_count"`
	PromptEvalDuration int     `json:"prompt_eval_duration"`
	EvalCount          int     `json:"eval_count"`
	EvalDuration       int     `json:"eval_duration"`
}

func (r *OllamaResponse) IsComplete() bool {
	return r.Done
}

func (r *OllamaResponse) GetMessageGroup() []sessionmgr.Message {
	// TODO
	return nil
}

func (r *OllamaResponse) Init(data []byte) error {
	// TODO

	return nil
}
