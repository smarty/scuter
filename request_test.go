package scuter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/smarty/scuter/internal/should"
)

func TestReadJSONRequestBody_MissingContentType(t *testing.T) {
	request := NewTestRequest(t.Context(), "PUT", "/")
	v := make(map[string]any)

	actual, ok := ReadJSONRequestBody(request, v)

	should.So(t, ok, should.BeFalse)
	assertResponseEqual(t, Response.JSONErrors(http.StatusUnsupportedMediaType, ErrUnsupportedRequestContentType), actual)
}
func TestReadJSONRequestBody_MalformedJSON(t *testing.T) {
	request := NewTestRequest(t.Context(), "PUT", "/", Request.With(
		Request.Header("Content-Type", "application/json"),
		Request.Body(strings.NewReader(`{invalid`)),
	))
	v := make(map[string]any)

	actual, ok := ReadJSONRequestBody(request, v)

	should.So(t, ok, should.BeFalse)
	assertResponseEqual(t, Response.JSONErrors(http.StatusBadRequest, ErrInvalidRequestJSONBody), actual)
}
func TestReadJSONRequestBody(t *testing.T) {
	request := NewTestRequest(t.Context(), "PUT", "/", Request.JSONBody(map[string]any{"a": 1, "b": 2}))
	v := make(map[string]any)

	actual, ok := ReadJSONRequestBody(request, v)

	should.So(t, ok, should.BeTrue)
	should.So(t, actual, should.BeNil)
}
func assertResponseEqual(t *testing.T, expected, actual ResponseOption) {
	t.Helper()
	EXPECTED := httptest.NewRecorder()
	ACTUAL := httptest.NewRecorder()

	Flush(EXPECTED, expected)
	Flush(ACTUAL, actual)

	should.So(t, ACTUAL.Code, should.Equal, EXPECTED.Code)
	should.So(t, ACTUAL.Header(), should.Equal, EXPECTED.Header())
	should.So(t, ACTUAL.Body.String(), should.Equal, EXPECTED.Body.String())
}
func TestReadUint64Header(t *testing.T) {
	should.So(t, ReadUint64Header(http.Header{"a": []string{"1"}}, "a"), should.Equal, uint64(1))
	should.So(t, ReadUint64Header(http.Header{"a": []string{"NaN"}}, "a"), should.Equal, uint64(0))
	should.So(t, ReadUint64Header(http.Header{"a": []string{"1"}}, "nope"), should.Equal, uint64(0))
}
func TestReadTimeHeader(t *testing.T) {
	date := time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC)
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"2025-01-02"}}, "2006-01-02", "a"), should.Equal, date)
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"nope"}}, "2006-01-02", "a"), should.Equal, time.Time{})
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"2025-01-02"}}, "nope", "a"), should.Equal, time.Time{})
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"2025-01-02"}}, "2006-01-02", "nope"), should.Equal, time.Time{})
}

func TestReadPathElement(t *testing.T) {
	assertPathElement(t, "/path/users/value", "users", "value")
	assertPathElement(t, "/path/users/value/", "users", "value")
	assertPathElement(t, "/path/users/value/other/stuff", "users", "value")
	assertPathElement(t, "/not/found/", "users", "")
}
func assertPathElement(t *testing.T, raw, element, expected string) {
	should.So(t, ReadPathElement(raw, element), should.Equal, expected)
}

func TestReadNumericPathElement(t *testing.T) {
	assertNumericPathElement(t, "/path/users/123", "users", 123)
	assertNumericPathElement(t, "/path/users/123/", "users", 123)
	assertNumericPathElement(t, "/path/users/123/other/stuff", "users", 123)
	assertNumericPathElement(t, "/path/users/-123/other/stuff", "users", 0)
	assertNumericPathElement(t, "/path/users/abc/other/stuff", "users", 0)
	assertNumericPathElement(t, "/path/users/abc123/other/stuff", "users", 0)
	assertNumericPathElement(t, "/path/users/_123/other/stuff", "users", 0)
}
func assertNumericPathElement(t *testing.T, raw, element string, expected uint64) {
	should.So(t, ReadNumericPathElement(raw, element), should.Equal, expected)
}
