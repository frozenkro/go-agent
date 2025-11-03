package main

import (
	"log"
	"os"

	"github.com/frozenkro/go-agent/goagent"
	"github.com/frozenkro/go-agent/models/anthropic"
)

func main() {
	// Create a minimal agent with the bash tool and default settings
	agent, err := goagent.NewAgent(goagent.WithTools(anthropic.BASH))
	if err != nil {
		log.Fatal("Failed to create agent:", err)
	}

	// Run a simple task
	ch := agent.Run("List all files in the current directory")
	for event := range ch {
		if event.Error != nil {
			log.Fatalf("Error: %v", event.Error)
		}

		log.Println(event.Message)
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

	// clean up between example runs
	os.Remove("fib/fib.go")
	// Run a more complex task
	ch = customAgent.Run("Create a simple Go function that calculates fibonacci numbers and save it to fib/fib.go")
	for event := range ch {
		if event.Error != nil {
			log.Fatalf("Error: %v", event.Error)
		}

		log.Println(event.Message)
	}
}
