package texteditor

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	mock_texteditor "github.com/frozenkro/goagent/internal/tools/texteditor/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	testResultsDir string = "test_results"
	testDataDir    string = "test_data"
)

func TestMain(m *testing.M) {
	_ = os.RemoveAll(testResultsDir)
	m.Run()
}

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

func TestHandleStrReplace(t *testing.T) {
	fileName := copyTestFile(t)

	line1 := "Butt Line 1"
	line3 := "Butt Line 3"

	params := map[string]any{}
	params["path"] = fileName
	params["old_str"] = "Test"
	params["new_str"] = "Butt"

	sut := TextEditorWorkerImpl{}
	res, err := sut.HandleStrReplace(params)
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(res, line1))
	assert.True(t, strings.HasSuffix(strings.Trim(res, "\n"), line3))

	written, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	res = string(written)
	assert.True(t, strings.HasPrefix(res, line1))
	assert.True(t, strings.HasSuffix(strings.Trim(res, "\n"), line3))
}

func TestHandleCreate(t *testing.T) {
	fileName := fmt.Sprintf("%v/%v.txt", testResultsDir, t.Name())
	line1 := "Line 1"
	line2 := "Line 2"
	params := map[string]any{}
	params["path"] = fileName
	params["file_text"] = fmt.Sprintf("%v\n%v", line1, line2)

	sut := TextEditorWorkerImpl{}
	_, err := sut.HandleCreate(params)
	assert.Nil(t, err)

	written, err := os.ReadFile(fileName)
	res := string(written)
	assert.True(t, strings.HasPrefix(res, line1))
	assert.True(t, strings.HasSuffix(res, line2))
}

func TestHandleCreate_Nested(t *testing.T) {
	fileName := fmt.Sprintf("%v/nested/%v.txt", testResultsDir, t.Name())
	line1 := "Line 1"
	line2 := "Line 2"
	params := map[string]any{}
	params["path"] = fileName
	params["file_text"] = fmt.Sprintf("%v\n%v", line1, line2)

	sut := TextEditorWorkerImpl{}
	_, err := sut.HandleCreate(params)
	assert.Nil(t, err)

	written, err := os.ReadFile(fileName)
	res := string(written)
	assert.True(t, strings.HasPrefix(res, line1))
	assert.True(t, strings.HasSuffix(res, line2))
}

func TestHandleInsert(t *testing.T) {
	fileName := copyTestFile(t)
	insert := "Line one point five"
	params := map[string]any{}
	params["path"] = fileName
	params["insert_line"] = 2
	params["new_str"] = insert

	sut := TextEditorWorkerImpl{}
	_, err := sut.HandleInsert(params)
	assert.Nil(t, err)

	b, err := os.ReadFile(fileName)
	res := strings.Split(string(b), "\n")

	exp := []string{
		"Test Line 1",
		insert,
		"Test Line 2",
		"Test Line 3",
	}

	for i, v := range exp {
		assert.Equal(t, v, res[i])
	}

}

func copyTestFile(t *testing.T) string {
	fileName := fmt.Sprintf("%v/%v.txt", testResultsDir, t.Name())

	src, err := os.Open(fmt.Sprintf("%v/test_file.txt", testDataDir))
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(testResultsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	dst, err := os.Create(fileName)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		t.Fatal(err)
	}

	return fileName
}
