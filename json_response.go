package scuter

import (
	"encoding/json"
	"net/http"
)

type JSONResponse[T any] struct {
	StatusCode int
	Content    T
}

type JSONResponder[T any] struct{}

func (this *JSONResponder[T]) RespondResponse(writer http.ResponseWriter, response *JSONResponse[T]) {
	this.Respond(writer, response.StatusCode, response.Content)
}
func (this *JSONResponder[T]) Respond(writer http.ResponseWriter, code int, content T) {
	// TODO: receive request and ensure that the Accept header doesn't conflict with a JSON response.
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	_ = json.NewEncoder(writer).Encode(content) // TODO: can we use json/v2 to achieve 100% reuse?
}
