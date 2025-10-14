package main

import (
	"log"

	"github.com/frozenkro/go-agent/goagent"
	"github.com/frozenkro/go-agent/models/anthropic"
)

// const TEST_PROMPT = "List all files in the current directory"
const TEST_PROMPT = "Use the text editor tool to insert lines into my fizzbuzz.sh"

func main() {
	agent, err := goagent.NewAgent(
		goagent.WithModel(anthropic.SONNET_4),
		goagent.WithTools(anthropic.BASH, anthropic.TEXT_EDITOR),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = agent.Run(TEST_PROMPT)
	if err != nil {
		log.Fatal(err.Error())
	}
}
