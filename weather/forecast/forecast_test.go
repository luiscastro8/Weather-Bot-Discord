package forecast

import (
	//"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestProcessResponse(t *testing.T) {
	t.Run("It should return an error if the status code was not 2XX", func(t *testing.T) {
		s, err := processResponse(createMockResponse("", 500), "")
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if the status code was not 2XX and the response has error details", func(t *testing.T) {
		panic("this is not")
	})
}

func createMockResponse(body string, code int) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func createMockResponseBody() string {
	//bytes, err := json.Marshal()
	//if err != nil {
	//	return ""
	//}
	return ""
}

func createMockResponseBodyError(r errorResponse) string {
	return ""
}
