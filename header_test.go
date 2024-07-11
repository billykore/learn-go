package learning

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RequestHeader(writer http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("content-type")
	fmt.Fprint(writer, contentType)
}

func TestHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "localhost:8000", nil)
	request.Header.Add("Content-Type", "application/json")

	RequestHeader(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	fmt.Printf(string(body))
}

func ResponseHeader(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("X-Powered-By", "Billy Kore")
	fmt.Fprint(writer, "OK")
}

func TestResponseHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "localhost:8000", nil)

	ResponseHeader(recorder, request)

	poweredBy := recorder.Header().Get("x-powered-by")
	fmt.Println(poweredBy)
}
