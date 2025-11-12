package sessionmgr

import "github.com/frozenkro/goagent/internal/tools"

type RequestHandler interface {
	Init(model, prompt string) error
	AddTool(*tools.ToolMeta) error
	SetMaxTokens(int) error
	AddStatementGroup([]Statement) error
}
