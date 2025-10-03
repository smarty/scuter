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

func (this *JSONResponder[T]) SerializeJSON(writer http.ResponseWriter, response *JSONResponse[T]) {
	writer.WriteHeader(response.StatusCode)
	_ = json.NewEncoder(writer).Encode(response.Content) // TODO: can we use json to achieve 100% reuse?
}
