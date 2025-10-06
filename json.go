//go:build goexperiment.jsonv2

package scuter

import (
	v1 "encoding/json"
	"encoding/json/v2"
	"io"
	"net/http"
)

// JSONOptionsV1 will likely be a more convenient way to reference this value from the v1 encoding/json package.
var JSONOptionsV1 = v1.DefaultOptionsV1()

// DeserializeJSON is a convenience function that masks io.EOF from json.UnmarshalRead.
func DeserializeJSON(request *http.Request, v any, options ...json.Options) error {
	err := json.UnmarshalRead(request.Body, &v, options...)
	if err == io.EOF {
		err = nil
	}
	return err
}

// SerializeJSON is a convenience function that utilizes the Flush API defined in this package to write a response.
func SerializeJSON(writer http.ResponseWriter, code int, content any, options ...json.Options) error {
	return Flush(writer,
		Response.StatusCode(code),
		Response.JSONBody(content, options...),
	)
}
