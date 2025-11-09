package sessionmgr

type ResponseHandler interface {
	Init([]byte) error
	IsComplete() bool
	GetMessageGroup() []Message
}
