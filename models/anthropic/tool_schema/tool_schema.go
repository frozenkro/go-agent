package toolschema

type BashToolInput struct {
	Command string `json:"command"`
	Restart bool   `json:"restart"`
}
