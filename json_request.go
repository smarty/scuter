package scuter

import (
	"encoding/json"
	"io"
	"net/http"
)

type JSONRequest struct{}

func (this *JSONRequest) DeserializeJSON(request *http.Request, v any) bool {
	// TODO: ensure the Content-Type doesn't conflict with JSON decoding.
	err := json.NewDecoder(request.Body).Decode(v) // TODO: can we use json/v2 to achieve 100% reuse?
	if err == nil {
		return true
	}
	if err == io.EOF {
		return true
	}
	return false
}
