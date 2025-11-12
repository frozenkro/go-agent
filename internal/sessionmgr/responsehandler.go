package sessionmgr

type ResponseHandler interface {
	Init([]byte) error
	IsComplete() bool
	GetStatementGroup() ([]Statement, error)
}
