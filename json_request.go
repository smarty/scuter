//go:build goexperiment.jsonv2

package scuter

import (
	"encoding/json/v2"
	"io"
	"net/http"
)

type JSONRequest struct {
	options []json.Options
	logger  Logger
}

func NewJSONRequest(logger Logger, options ...json.Options) *JSONRequest {
	return &JSONRequest{
		logger:  logger,
		options: options,
	}
}
func (this *JSONRequest) DeserializeJSON(request *http.Request, v any) bool {
	// TODO: ensure the Content-Type doesn't conflict with JSON decoding.
	err := json.UnmarshalRead(request.Body, &v, this.options...)
	if err == nil {
		return true
	}
	if err == io.EOF {
		return true
	}
	return false
}
