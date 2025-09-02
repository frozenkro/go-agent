package options

import (
	"log"

	"github.com/frozenkro/go-agent/internal/models/anthropic"
	"github.com/frozenkro/go-agent/internal/tools"
)

func WithTools(toolNames ...anthropic.ToolName) anthropic.AnthropicMessagesRequestOption {

	return func(a *anthropic.AnthropicMessagesRequest) {

		toolMap := tools.InitToolMap()

		for _, toolName := range toolNames {
			toolMeta, err := toolMap.ToolMetaByName(toolName)

			if err == nil {
				a.Tools = append(a.Tools, toolMeta.Spec)
			} else {
				log.Printf(err.Error())
			}
		}
	}
}
