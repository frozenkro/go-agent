package tools

import "github.com/frozenkro/goagent/models/toolschema"

type Tool interface {
	Invoke(params any) (string, error)
	Spec() toolschema.ToolSpec
}
