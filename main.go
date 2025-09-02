package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/frozenkro/go-agent/internal/agents"
	"github.com/frozenkro/go-agent/internal/models/anthropic"
	"github.com/frozenkro/go-agent/internal/options"
	"github.com/joho/godotenv"
)

const ANTHROPIC_MESSAGES_URL = "https://api.anthropic.com/v1/messages"
const TEST_PROMPT = "List all files in the current directory"

type AnthropicHandler interface {
	HandleResponse(anthropic.AnthropicMessagesResponse) (anthropic.AnthropicMessagesRequest, bool, error)
	InitRequest(anthropic.Model, string, ...anthropic.AnthropicMessagesRequestOption)
}

func main() {
	var (
		request  *anthropic.AnthropicMessagesRequest
		response *anthropic.AnthropicMessagesResponse
		done     bool
	)

	ctx := context.Background()
	godotenv.Load()

	anthropicAgent, err := agents.NewAnthropicAgent()
	if err != nil {
		log.Fatal(err.Error())
	}

	request = anthropicAgent.InitRequest(anthropic.SONNET_4, TEST_PROMPT, options.WithTools(anthropic.BASH))

	for {
		reqJson, err := json.Marshal(request)
		if err != nil {
			log.Fatal(err.Error())
		}

		resBytes, err := postMessage(ctx, string(reqJson))
		if err != nil {
			log.Fatal(err.Error())
		}

		response = &anthropic.AnthropicMessagesResponse{}
		err = json.Unmarshal(resBytes, response)
		if err != nil {
			log.Fatal(err.Error())
		}

		request, done, err = anthropicAgent.HandleResponse(response)
		if done {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func postMessage(ctx context.Context, body string) ([]byte, error) {
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
	return content, nil
}
