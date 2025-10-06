package texteditor

import (
	"fmt"
	"strings"
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
	command, ok := m["command"].(string)
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
			params["command"] = testCase

			res, err := tool.Invoke(params)
			assert.Nil(t, err)
			assert.Equal(t, testCase, res)
		})
	}
}

func TestHandleView_NoRange(t *testing.T) {
	params := map[string]any{}
	params["path"] = "test_data/test_file.txt"

	sut := TextEditorWorkerImpl{}
	res, err := sut.HandleView(params)
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(res, "1. Test Line 1"))
	assert.True(t, strings.HasSuffix(strings.Trim(res, "\n"), "3. Test Line 3"))
}

// func TestHandleStrReplace(t *testing.T) {
// 	src, err := os.Open("test_data/test_file.txt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = os.Mkdir("test_results", 0644)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	dst, err := os.Create("test_results/test_file.txt")
// 	_, err = io.Copy(dst, src)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
