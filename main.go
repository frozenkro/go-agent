package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const ANTHROPIC_MESSAGES_URL = "https://api.anthropic.com/v1/messages"
const TEST_REQUEST = `
{
    "model": "claude-sonnet-4-20250514",
    "max_tokens": 1024,
    "messages": [
        {"role": "user", "content": "Hello, world"}
    ]
}
`

const EXEC_TOOL_SCHEMA = `
{
	"name": "execute_command",
	"description": "Run any common bash command, such as 'ls', 'cd', 'cat', 'sed', etc.",
	"input_schema": {
		"type": "object",
		"properties": {
			"command": {
					"type": "string",
					"description": "Arbitrary bash command"
			}
		},
		"required": ["command"]
	}
}
`

func main() {
	ctx := context.Background()
	godotenv.Load()

	postMessage(ctx, TEST_REQUEST)
}

func postMessage(ctx context.Context, body string) (any, error) {
	apiKey := os.Getenv("GA_ANTHROPIC_API_KEY")
	bodyReader := bytes.NewReader([]byte(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ANTHROPIC_MESSAGES_URL, bodyReader)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Response:\n%v", string(content))
	return string(content), nil
}
