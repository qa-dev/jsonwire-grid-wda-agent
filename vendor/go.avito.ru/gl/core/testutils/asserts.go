package testutils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// AssertResponse проверяет что ответ содержит имеет определенный статус ответа и Content-type.
func AssertResponse(t *testing.T,
	w *httptest.ResponseRecorder,
	code int,
	contentType string) {
	if w.Code != code {
		t.Errorf(
			`Status code, expected: %v, actual: %v for body "%v"`,
			code, w.Code, w.Body)
	}

	if w.HeaderMap.Get("Content-type") != contentType {
		t.Errorf(
			"Content type, expected: %v, actual: %v.",
			contentType,
			w.HeaderMap.Get("Content-type"))
	}
}

// AssertJSON проверяет что данные содержат правильные json-данные.
func AssertJSON(t *testing.T, body []byte, data interface{}) {
	err := json.Unmarshal(body, data)
	if err != nil {
		t.Fatalf("Invalid json %v.", body)
	}
}
