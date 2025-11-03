# Go Agent Library

A Go library for creating AI agents that can execute tasks using Anthropic's Claude AI with various tools.

## Usage

### Basic Usage

```go
package main

import (
    "log"
    "github.com/frozenkro/go-agent/goagent"
)

func main() {
    // Create a new agent with default settings
    agent, err := goagent.NewAgent()
    if err != nil {
        log.Fatal("Failed to create agent:", err)
    }

    // Run a task
    err = agent.Run("List all files in the current directory")
    if err != nil {
        log.Fatal("Failed to run task:", err)
    }
}
```

### Advanced Usage

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/frozenkro/go-agent/goagent"
    "github.com/frozenkro/go-agent/models/anthropic"
)

func main() {
    // Create an agent with custom options
    agent, err := goagent.NewAgent(
        goagent.WithModel(anthropic.SONNET_4),
        goagent.WithTools(anthropic.BASH, anthropic.TEXT_EDITOR),
        goagent.WithMaxTokens(2048),
        goagent.WithAPIKey("your-api-key"), // Optional, uses GA_ANTHROPIC_API_KEY env var by default
    )
    if err != nil {
        log.Fatal("Failed to create agent:", err)
    }

    // Run a task with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    err = agent.RunWithContext(ctx, "Create a simple Go function that calculates fibonacci numbers")
    if err != nil {
        log.Fatal("Failed to run task:", err)
    }
}
```

## Configuration

### Environment Variables

- `GA_ANTHROPIC_API_KEY`: Your Anthropic API key (required unless provided via `WithAPIKey` option)

### Agent Options

- `WithModel(model)`: Set the AI model to use (default: `anthropic.SONNET_4`)
- `WithTools(tools...)`: Set available tools (default: `anthropic.BASH`, `anthropic.TEXT_EDITOR`)
- `WithAPIKey(key)`: Set a custom API key
- `WithMaxTokens(tokens)`: Set maximum tokens for responses (default: 1024)

### Available Models

- `anthropic.SONNET_4`: Claude 3.5 Sonnet (default)

### Available Tools

- `anthropic.BASH`: Execute bash commands
- `anthropic.TEXT_EDITOR`: Edit text files

## API Reference

### `NewAgent(opts ...AgentOption) (*Agent, error)`

Creates a new agent with the given options.

### `(a *Agent) Run(task string) error`

Executes the given task using the agent.

### `(a *Agent) RunWithContext(ctx context.Context, task string) error`

Executes the given task using the agent with a custom context for cancellation/timeout.

## Error Handling

All methods return detailed error messages that wrap the underlying errors. Common error scenarios:

- Missing API key
- Network/HTTP errors
- Anthropic API errors
- JSON marshaling/unmarshaling errors

## Examples

See the `/example` directory for complete working examples.