package main

import (
	"log"

	"github.com/frozenkro/go-agent/goagent"
	"github.com/frozenkro/go-agent/models/anthropic"
)

func main() {
	// Create a new agent with default settings
	agent, err := goagent.NewAgent()
	if err != nil {
		log.Fatal("Failed to create agent:", err)
	}

	// Run a simple task
	err = agent.Run("List all files in the current directory")
	if err != nil {
		log.Fatal("Failed to run task:", err)
	}

	// Create an agent with custom options
	customAgent, err := goagent.NewAgent(
		goagent.WithModel(anthropic.SONNET_4),
		goagent.WithTools(anthropic.BASH, anthropic.TEXT_EDITOR),
		goagent.WithMaxTokens(2048),
	)
	if err != nil {
		log.Fatal("Failed to create custom agent:", err)
	}

	// Run a more complex task
	err = customAgent.Run("Create a simple Go function that calculates fibonacci numbers and save it to fib.go")
	if err != nil {
		log.Fatal("Failed to run custom task:", err)
	}
}
