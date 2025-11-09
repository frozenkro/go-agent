package goagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/frozenkro/goagent/internal/sessionmgr"
	"github.com/frozenkro/goagent/models"
	"github.com/frozenkro/goagent/models/anthropic"
	"github.com/frozenkro/goagent/models/ollama"
	"github.com/joho/godotenv"
)

const ANTHROPIC_MESSAGES_URL = "https://api.anthropic.com/v1/messages"

// Agent represents a Go agent that can execute tasks using AI
type Agent struct {
	model        string
	provider     Provider
	tools        []string
	apiKey       string
	maxTokens    int
	outputWriter *io.Writer
	url          string
}

// AgentOption is a function type for configuring the Agent
type AgentOption func(*Agent)

// WithModel sets the AI model to use
func WithModel(model string) AgentOption {
	return func(a *Agent) {
		a.model = model
	}
}

// WithTools sets the available tools for the agent
func WithTools(tools ...string) AgentOption {
	return func(a *Agent) {
		a.tools = tools
	}
}

// WithAPIKey sets a custom API key (otherwise uses GA_API_KEY env var)
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

func WithProviderUrl(url string) AgentOption {
	return func(a *Agent) {
		a.url = url
	}
}

// NewAgent creates a new Agent
func NewAgent(provider Provider, opts ...AgentOption) (*Agent, error) {
	godotenv.Load()

	agent := &Agent{
		provider:  provider,
		model:     anthropic.SONNET_4,
		apiKey:    os.Getenv("GA_API_KEY"),
		maxTokens: 1024,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent.validate()
}

func (a *Agent) validate() (*Agent, error) {
	if a.provider == ANTHROPIC && a.apiKey == "" {
		return nil, fmt.Errorf("No anthropic api key provided. Either pass WithAPIKey AgentOption or set 'GA_API_KEY' environment variable")
	}

	if a.provider == OLLAMA && a.url == "" {
		return nil, fmt.Errorf("No ollama url provided. Use WithProviderUrl AgentOption.")
	}

	return a, nil
}

// Run executes the given task using the agent
// Returns a channel which can be looped over to retrieve output in real time
func (a *Agent) Run(task string) <-chan models.AgentEvent {
	return a.RunWithContext(context.Background(), task)
}

// RunWithContext executes the given task using the agent within a provided context
// This contains the elusive "loop" of the agent
// Returns a channel which can be looped over to retrieve output in real time
func (a *Agent) RunWithContext(ctx context.Context, task string) <-chan models.AgentEvent {
	out := make(chan models.AgentEvent)

	go func() {
		defer close(out)

		request, err := a.newRequestHandler()
		if err != nil {
			out <- models.AgentEvent{Error: fmt.Errorf("failed to initialize provider: %w", err)}
		}

		opts := []sessionmgr.SessionMgrOption{}
		if a.maxTokens > 0 {
			opts = append(opts, sessionmgr.WithMaxTokens(a.maxTokens))
		}
		if len(a.tools) > 0 {
			opts = append(opts, sessionmgr.WithTools(a.tools...))
		}
		mgr, err := sessionmgr.NewSessionMgr(request,
			a.model,
			task,
			opts...,
		)
		if err != nil {
			out <- models.AgentEvent{Error: fmt.Errorf("failed to initialize session: %w", err)}
			return
		}

		for {
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

			response, err := a.newResponseHandler()
			if err != nil {
				out <- models.AgentEvent{Error: err}
				return
			}

			err = response.Init(resBytes)
			if err != nil {
				out <- models.AgentEvent{Error: fmt.Errorf("Failed to parse response from provider: %w", err)}
				return
			}

			done, err := mgr.HandleResponse(response, out)
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

// postMessage posts `body` to the llm provider's chat endpoint
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

func (a *Agent) newRequestHandler() (sessionmgr.RequestHandler, error) {
	switch a.provider {
	case ANTHROPIC:
		return &anthropic.MessagesRequest{}, nil
	case OLLAMA:
		return &ollama.OllamaRequest{}, nil
	}

	return nil, fmt.Errorf("No request type available for provider: '%v'", a.provider)
}

func (a *Agent) newResponseHandler() (sessionmgr.ResponseHandler, error) {
	switch a.provider {
	case ANTHROPIC:
		return &anthropic.MessagesResponse{}, nil
	case OLLAMA:
		return &ollama.OllamaResponse{}, nil
	}

	return nil, fmt.Errorf("No response type available for provider: '%v'", a.provider)
}
