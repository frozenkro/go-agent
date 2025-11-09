package goagent

import (
	"github.com/frozenkro/goagent/models/anthropic"
	"github.com/frozenkro/goagent/models/toolschema"
)

const (
	BASH_TOOL        = toolschema.BASH
	TEXT_EDITOR_TOOL = toolschema.TEXT_EDITOR
	ANTH_SONNET_4    = anthropic.SONNET_4
)

type Provider int

const (
	ANTHROPIC Provider = iota + 1
	OLLAMA
)
