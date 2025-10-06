package texteditor

import (
	"fmt"
	"testing"

	mock_texteditor "github.com/frozenkro/go-agent/internal/tools/texteditor/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getMockResult(x any) (string, error) {
	m, ok := x.(map[string]any)
	if !ok {
		return "", fmt.Errorf("Parse Error in getMockResult")
	}
	command, ok := m["Command"].(string)
	if !ok {
		return "", fmt.Errorf("Parse Error in getMockResult")
	}
	return command, nil
}
func TestInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	workerMock := mock_texteditor.NewMockTextEditorWorker(ctrl)
	tool := TextEditorTool{
		w: workerMock,
	}

	tests := []string{
		"view",
		"str_replace",
		"create",
		"insert",
		"undo_edit",
	}

	workerMock.EXPECT().HandleView(gomock.Any()).DoAndReturn(getMockResult).AnyTimes()
	workerMock.EXPECT().HandleStrReplace(gomock.Any()).DoAndReturn(getMockResult).AnyTimes()
	workerMock.EXPECT().HandleCreate(gomock.Any()).DoAndReturn(getMockResult).AnyTimes()
	workerMock.EXPECT().HandleInsert(gomock.Any()).DoAndReturn(getMockResult).AnyTimes()
	workerMock.EXPECT().HandleUndoEdit(gomock.Any()).DoAndReturn(getMockResult).AnyTimes()

	for _, testCase := range tests {
		testName := fmt.Sprintf("TextEditorToolInvoke-%v", testCase)
		t.Run(testName, func(t *testing.T) {
			params := map[string]interface{}{}
			params["Command"] = testCase

			res, err := tool.Invoke(params)
			assert.Nil(t, err)
			assert.Equal(t, testCase, res)
		})
	}
}
