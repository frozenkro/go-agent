package goagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	apimgrs "github.com/frozenkro/go-agent/internal/api_mgrs"
	"github.com/frozenkro/go-agent/models"
	"github.com/frozenkro/go-agent/models/anthropic"
	"github.com/joho/godotenv"
)

const ANTHROPIC_MESSAGES_URL = "https://api.anthropic.com/v1/messages"

// Agent represents a Go agent that can execute tasks using AI
type Agent struct {
	model        anthropic.Model
	tools        []anthropic.ToolName
	apiKey       string
	maxTokens    int
	outputWriter *io.Writer
}

// AgentOption is a function type for configuring the Agent
type AgentOption func(*Agent)

// WithModel sets the AI model to use
func WithModel(model anthropic.Model) AgentOption {
	return func(a *Agent) {
		a.model = model
	}
}

// WithTools sets the available tools for the agent
func WithTools(tools ...anthropic.ToolName) AgentOption {
	return func(a *Agent) {
		a.tools = tools
	}
}

// WithAPIKey sets a custom API key (otherwise uses GA_ANTHROPIC_API_KEY env var)
func WithAPIKey(apiKey string) AgentOption {
	return func(a *Agent) {
		a.apiKey = apiKey
	}
}

// WithMaxTokens sets the maximum number of tokens for responses
func WithMaxTokens(maxTokens int) AgentOption {
	return func(a *Agent) {
		a.maxTokens = maxTokens
	}
}

// NewAgent creates a new Agent
func NewAgent(opts ...AgentOption) (*Agent, error) {
	godotenv.Load()

	agent := &Agent{
		model:     anthropic.SONNET_4,
		apiKey:    os.Getenv("GA_ANTHROPIC_API_KEY"),
		maxTokens: 1024,
	}

	for _, opt := range opts {
		opt(agent)
	}

	if agent.apiKey == "" {
		return nil, fmt.Errorf("API key not provided. Set GA_ANTHROPIC_API_KEY environment variable or use WithAPIKey option")
	}

	return agent, nil
}

// Run executes the given task using the agent
func (a *Agent) Run(task string) <-chan models.AgentEvent {
	return a.RunWithContext(context.Background(), task)
}

// RunWithContext executes the given task using the agent within a provided context
// This contains the elusive "loop" of the agent
func (a *Agent) RunWithContext(ctx context.Context, task string) <-chan models.AgentEvent {
	out := make(chan models.AgentEvent)

	go func() {
		defer close(out)

		anthropicClient, err := apimgrs.NewAnthropicClient(a.model,
			task,
			apimgrs.WithTools(a.tools...),
			apimgrs.WithMaxTokens(a.maxTokens))
		if err != nil {
			out <- models.AgentEvent{Error: fmt.Errorf("failed to create anthropic agent: %w", err)}
			return
		}

		for {
			request := anthropicClient.GetRequest()
			reqJson, err := json.Marshal(request)
			if err != nil {
				out <- models.AgentEvent{Error: fmt.Errorf("failed to marshal request: %w", err)}
				return
			}

			resBytes, err := a.postMessage(ctx, string(reqJson))
			if err != nil {
				out <- models.AgentEvent{Error: fmt.Errorf("failed to send message: %w", err)}
				return
			}

			err = a.checkMessagesResponseErr(resBytes)
			if err != nil {
				out <- models.AgentEvent{Error: fmt.Errorf("API error: %w", err)}
				return
			}

			response := &anthropic.MessagesResponse{}
			err = json.Unmarshal(resBytes, response)
			if err != nil {
				out <- models.AgentEvent{Error: fmt.Errorf("failed to unmarshal response: %w", err)}
				return
			}

			done, err := anthropicClient.HandleResponse(response, out)
			if done {
				break
			}
			if err != nil {
				out <- models.AgentEvent{Error: fmt.Errorf("failed to handle response: %w", err)}
				return
			}
		}
	}()
	return out
}

// postMessage posts `body` to the anthropic messages api
func (a *Agent) postMessage(ctx context.Context, body string) ([]byte, error) {
	bodyReader := bytes.NewReader([]byte(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ANTHROPIC_MESSAGES_URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-api-key", a.apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return content, nil
}

// checkMessagesResponseErr checks if anthropic response returned an error according to their schema
func (a *Agent) checkMessagesResponseErr(data []byte) error {
	baseRes := &anthropic.MessagesBaseResponse{}
	if err := json.Unmarshal(data, baseRes); err != nil {
		return fmt.Errorf("failed to unmarshal base response: %w", err)
	}

	if baseRes.Type == "error" {
		errRes := &anthropic.MessagesErrorResponse{}
		if err := json.Unmarshal(data, errRes); err != nil {
			return fmt.Errorf("failed to unmarshal error response: %w", err)
		}

		return fmt.Errorf("Anthropic API error - type: %s, message: %s", errRes.Error.Type, errRes.Error.Message)
	}
	return nil
}
