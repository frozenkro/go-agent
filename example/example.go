package main

import (
	"log"
	"os"

	goagent "github.com/frozenkro/goagent/agent"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-ba":
			basicAnthropic()
		case "-aa":
			advancedAnthropic()
		case "-bo":
			basicOllama()
		case "-ao":
			advancedOllama()
		default:
			defaultExamples()
		}
	} else {
		defaultExamples()
	}

}

func defaultExamples() {
	basicAnthropic()
	advancedAnthropic()
}

func basicAnthropic() {
	// Create a minimal agent with the bash tool and default settings
	agent, err := goagent.NewAgent(goagent.ANTHROPIC, goagent.WithTools(goagent.BASH_TOOL))
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
}

func advancedAnthropic() {
	// Create an agent with custom options
	customAgent, err := goagent.NewAgent(
		goagent.ANTHROPIC,
		goagent.WithModel(goagent.ANTH_SONNET_4),
		goagent.WithTools(goagent.BASH_TOOL, goagent.TEXT_EDITOR_TOOL),
		goagent.WithMaxTokens(2048),
	)
	if err != nil {
		log.Fatal("Failed to create custom agent:", err)
	}

	// clean up between example runs
	os.Remove("fib/fib.go")
	// Run a more complex task
	ch := customAgent.Run("Create a simple Go function that calculates fibonacci numbers and save it to fib/fib.go. Then, test the function.")
	for event := range ch {
		if event.Error != nil {
			log.Fatalf("Error: %v", event.Error)
		}

		log.Println(event.Message)
	}
}

func basicOllama() {
	// Create a minimal agent with the bash tool and default settings
	agent, err := goagent.NewAgent(
		goagent.OLLAMA,
		goagent.WithModel("gpt-oss:20b"),
		// goagent.WithModel("llama3.1:8b"),
		goagent.WithTools(goagent.BASH_TOOL),
	)
	if err != nil {
		log.Fatal("Failed to create agent:", err)
	}

	// Run a simple task
	ch := agent.Run("List all the files in the current directory, then explain what each one is likely for.")
	for event := range ch {
		if event.Error != nil {
			log.Fatalf("Error: %v", event.Error)
		}

		log.Println(event.Message)
	}
}

func advancedOllama() {
	// Create an agent with custom options
	customAgent, err := goagent.NewAgent(
		goagent.OLLAMA,
		goagent.WithModel("gpt-oss:20b"),
		// goagent.WithModel("llama3.1:8b"),
		goagent.WithProviderUrl("http://localhost:11434/api/chat"),
		goagent.WithTools(goagent.BASH_TOOL, goagent.TEXT_EDITOR_TOOL),
		goagent.WithMaxTokens(2048),
	)
	if err != nil {
		log.Fatal("Failed to create custom agent:", err)
	}

	// clean up between example runs
	os.Remove("fib/fib.go")
	// Run a more complex task
	ch := customAgent.Run("Create a simple Go function that calculates fibonacci numbers and save it to fib/fib.go. Then, test the function.")
	for event := range ch {
		if event.Error != nil {
			log.Fatalf("Error: %v", event.Error)
		}

		log.Println(event.Message)
	}
}
