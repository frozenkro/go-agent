package sessionmgr

import (
	"errors"
	"fmt"
	"log"

	"github.com/frozenkro/goagent/internal/tools"
	"github.com/frozenkro/goagent/models"
)

type SessionMgr struct {
	req         RequestHandler
	toolInvoker tools.ToolInvoker
}

type SessionMgrOption func(*SessionMgr) error

func WithTools(toolNames ...string) SessionMgrOption {

	return func(s *SessionMgr) error {
		for _, toolName := range toolNames {
			toolMeta, err := tools.ToolMapInst.ToolMetaByName(toolName)

			if err == nil {
				if err = s.req.AddTool(toolMeta); err != nil {
					return err
				}
			} else {
				log.Printf("%v", err.Error())
				return err
			}
		}

		return nil
	}
}

func WithMaxTokens(maxTokens int) SessionMgrOption {
	return func(s *SessionMgr) error {
		return s.req.SetMaxTokens(maxTokens)
	}
}

func NewSessionMgr(req RequestHandler, model, prompt string, opts ...SessionMgrOption) (SessionMgr, error) {

	s := SessionMgr{
		req:         req,
		toolInvoker: tools.ToolInvoker{},
	}

	s.req.Init(model, prompt)

	for _, opt := range opts {
		if err := opt(&s); err != nil {

			return s, err
		}
	}

	return s, nil
}

// HandleResponse accepts the llm response, which should implement ResponseHandler
// This handles stop sequences and invocation of tool calls
// Returns true if task is complete
func (s *SessionMgr) HandleResponse(response ResponseHandler, out chan models.AgentEvent) (bool, error) {
	complete := response.IsComplete()

	statements, err := response.GetStatementGroup()
	if err != nil {
		return complete, fmt.Errorf("Failed to retrieve statements from response: %w", err)
	}
	if err := s.req.AddStatementGroup(statements); err != nil {
		return complete, fmt.Errorf("Failed to append llm response statements to request context: %w", err)
	}

	newStatements := []Statement{}
	errs := []error{}

	// Iterate over all response content to build subsequent request content and task context
	for _, m := range statements {
		switch m.Type {
		case TOOL_CALL:
			out <- models.AgentEvent{Message: fmt.Sprintf("[Tool Call]: %v", m.ToolCall.Name)}

			toolResult, err := s.toolInvoker.Invoke(m.ToolCall.Name, m.ToolCall.Params)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			out <- models.AgentEvent{Message: fmt.Sprintf("[Tool Call Result]: %v", toolResult)}
			newStatements = append(newStatements, Statement{
				Role:     TOOL,
				Type:     TOOL_RESPONSE,
				Text:     toolResult,
				ToolCall: ToolCall{Name: m.ToolCall.Name, Id: m.ToolCall.Id},
			})

		case TEXT:
			out <- models.AgentEvent{Message: m.Text}

		case THINKING:
			out <- models.AgentEvent{Message: fmt.Sprintf("[Thinking]: %v", m.Text)}

		}
	}

	if err := s.req.AddStatementGroup(newStatements); err != nil {
		return complete, fmt.Errorf("Failed to append client request messages to request context: %w", err)
	}

	if len(errs) > 0 {
		err := errors.Join(errs...)
		return complete, err
	}
	return complete, nil
}
