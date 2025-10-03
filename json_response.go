//go:build goexperiment.jsonv2

package scuter

import (
	"encoding/json/v2"
	"net/http"
)

type JSONResponse[T any] struct {
	StatusCode int
	Content    T
}

type JSONResponder[T any] struct {
	logger  Logger
	options []json.Options
}

func NewJSONResponder[T any](logger Logger, options ...json.Options) *JSONResponder[T] {
	return &JSONResponder[T]{logger: logger, options: options}
}
func (this *JSONResponder[T]) RespondResponse(writer http.ResponseWriter, response *JSONResponse[T]) {
	this.Respond(writer, response.StatusCode, response.Content)
}
func (this *JSONResponder[T]) Respond(writer http.ResponseWriter, code int, content T) {
	// TODO: receive request and ensure that the Accept header doesn't conflict with a JSON response.
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	_ = json.MarshalWrite(writer, content, this.options...)
}
