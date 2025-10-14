package goagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/frozenkro/go-agent/agents"
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

func WithOutput(writer *io.Writer) AgentOption {
	return func(a *Agent) {
		a.outputWriter = writer
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
func (a *Agent) Run(task string) error {
	return a.RunWithContext(context.Background(), task)
}

// RunWithContext executes the given task using the agent within a provided context
func (a *Agent) RunWithContext(ctx context.Context, task string) error {
	anthropicClient, err := agents.NewAnthropicClient(a.model, task, agents.WithTools(a.tools...))
	if err != nil {
		return fmt.Errorf("failed to create anthropic agent: %w", err)
	}

	request := anthropicClient.GetRequest()
	request.MaxTokens = a.maxTokens

	for {
		reqJson, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		resBytes, err := a.postMessage(ctx, string(reqJson))
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		err = a.checkMessagesResponseErr(resBytes)
		if err != nil {
			return fmt.Errorf("API error: %w", err)
		}

		response := &anthropic.MessagesResponse{}
		err = json.Unmarshal(resBytes, response)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		nextRequest, done, err := anthropicClient.HandleResponse(response)
		if done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to handle response: %w", err)
		}
		request = nextRequest
	}

	return nil
}

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
