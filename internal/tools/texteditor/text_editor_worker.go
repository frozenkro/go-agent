package texteditor

type TextEditorWorker interface {
	HandleView(any) (string, error)
	HandleStrReplace(any) (string, error)
	HandleCreate(any) (string, error)
	HandleInsert(any) (string, error)
	HandleUndoEdit(any) (string, error)
}
