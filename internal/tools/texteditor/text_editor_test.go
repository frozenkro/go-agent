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

func TestHandleView_Ranges(t *testing.T) {
	const line1 = "1. Test Line 1"
	const line2 = "2. Test Line 2"
	const line3 = "3. Test Line 3"

	testCases := []struct {
		Name      string
		Range     string
		FirstLine string
		LastLine  string
	}{
		{
			Name:      "No Range",
			Range:     "",
			FirstLine: line1,
			LastLine:  line3,
		},
		{
			Name:      "Line 2 to End",
			Range:     "[2,-1]",
			FirstLine: line2,
			LastLine:  line3,
		},
		{
			Name:      "Line 1 to 2",
			Range:     "[1,2]",
			FirstLine: line1,
			LastLine:  line2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			params := map[string]any{}
			params["path"] = "test_data/test_file.txt"
			params["view_range"] = tc.Range

			sut := TextEditorWorkerImpl{}
			res, err := sut.HandleView(params)
			assert.Nil(t, err)
			if !assert.True(t, strings.HasPrefix(res, tc.FirstLine)) {
				fmt.Printf("Test Results: %v", res)
			}
			if !assert.True(t, strings.HasSuffix(strings.Trim(res, "\n"), tc.LastLine)) {
				fmt.Printf("Test Results: %v", res)
			}
		})
	}

}

func TestHandleView_Directory(t *testing.T) {
	file1 := "empty_file.txt"
	file2 := "test_file.txt"
	params := map[string]any{}
	params["path"] = "test_data"

	sut := TextEditorWorkerImpl{}
	res, err := sut.HandleView(params)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(res, file1))
	assert.True(t, strings.Contains(res, file2))
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
