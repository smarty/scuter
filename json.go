//go:build goexperiment.jsonv2

package scuter

import (
	"encoding/json/v2"
	"io"
	"net/http"
)

// DeserializeJSON is a convenience function that masks io.EOF from json.UnmarshalRead.
func DeserializeJSON(request *http.Request, v any, options ...json.Options) error {
	err := json.UnmarshalRead(request.Body, &v, options...)
	if err == io.EOF {
		err = nil
	}
	return err
}
