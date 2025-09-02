package tools

type Tool interface {
	Invoke(params any) (string, error)
}
