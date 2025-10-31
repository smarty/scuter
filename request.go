package scuter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUnsupportedRequestContentType = Error{
		Fields:  []string{"Content-Type"},
		Name:    "unsupported-content-type",
		Message: "The content-type was not supported.",
	}
	ErrInvalidRequestJSONBody = Error{
		Fields:  []string{"body"},
		Name:    "malformed-request-payload",
		Message: "The body did not contain well-formed data and could not be properly deserialized.",
	}
)

// ReadJSONRequestBody ensures the Content-Type header indicates JSON and if so, proceeds to unmarshal the body into
// the provided value. Failure at any point results in a JSON error which can be sent to the client with Flush.
func ReadJSONRequestBody(request *http.Request, v any) (ResponseOption, bool) {
	if !isJSONContent(request) {
		return Response.JSONErrors(http.StatusBadRequest, ErrUnsupportedRequestContentType), false
	}
	if err := json.NewDecoder(request.Body).Decode(&v); err != nil { // FUTURE: upgrade to json/v2's json.UnmarshalRead
		return Response.JSONErrors(http.StatusBadRequest, ErrInvalidRequestJSONBody), false
	}
	return nil, true
}

func isJSONContent(request *http.Request) bool {
	for _, contentType := range request.Header[headerContentType] {
		if strings.Contains(contentType, "json") {
			return true
		}
	}
	return false
}

// ReadUint64Header parses the first header value corresponding with key as a uint64.
func ReadUint64Header(headers http.Header, key string) uint64 {
	if values, contains := headers[key]; contains && len(values) > 0 {
		parsed, _ := strconv.ParseUint(values[0], 10, 64)
		return parsed
	}
	return 0
}

// ReadTimeHeader parses the first header value corresponding with key as a time.Time, else returns the Epoch time.
func ReadTimeHeader(headers http.Header, format, key string) (zero time.Time) {
	if values, contains := headers[key]; contains && len(values) > 0 {
		parsed, _ := time.Parse(format, values[0])
		return parsed
	}
	return zero
}

// ReadPathElement returns the path element immediately following the supplied label.
func ReadPathElement(rawPath, label string) string {
	found := false
	for element := range strings.SplitSeq(rawPath, "/") {
		if found == true {
			return element
		}
		if element == label {
			found = true
		}
	}
	return ""
}

// ReadNumericPathElement returns the path element immediately following the supplied label, parsed as a uint64
func ReadNumericPathElement(rawPath, label string) uint64 {
	raw := ReadPathElement(rawPath, label)
	value, _ := strconv.ParseUint(raw, 10, 64)
	return value
}
