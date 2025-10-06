package scuter

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
)

// Flush applies the options on the provide ResponseWriter.
// The options must be supplied in valid order, meaning that
// options which set headers should come before any that write
// a status code, which must come before any that write the body.
func Flush(response http.ResponseWriter, options ...ResponseOption) (err error) {
	for _, option := range options {
		err = option(response)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResponseOption is a callback func with an opportunity to call methods on http.ResponseWriter.
type ResponseOption func(http.ResponseWriter) error

// Response is the 'namespace' for all methods that return a ResponseOption.
var Response responseSingleton

type responseSingleton struct{}

// Header adds the value associated with key.
func (responseSingleton) Header(key, value string) ResponseOption {
	return func(response http.ResponseWriter) error {
		response.Header().Add(key, value)
		return nil
	}
}

// ContentType sets the 'Content-Type' header.
func (responseSingleton) ContentType(mime string) ResponseOption {
	return func(response http.ResponseWriter) error {
		return Response.Header(headerContentType, mime)(response)
	}
}

// JSONContentType sets the 'Content-Type' header to a sensible value representing JSON.
func (responseSingleton) JSONContentType() ResponseOption {
	return func(response http.ResponseWriter) error {
		return Response.ContentType(jsonContentType)(response)
	}
}

// StatusCode sets the status code (and writes all headers).
func (responseSingleton) StatusCode(code int) ResponseOption {
	return func(response http.ResponseWriter) error {
		response.WriteHeader(code)
		return nil
	}
}

// BytesBody writes the bytes to the ResponseWriter and returns any error.
func (responseSingleton) BytesBody(b []byte) ResponseOption {
	return func(response http.ResponseWriter) error {
		_, err := response.Write(b)
		return err
	}
}

// JSONBody uses json.MarshalWrite to serialize v to the ResponseWriter using the provided options and returning any error.
func (responseSingleton) JSONBody(v any, options ...json.Options) ResponseOption {
	return func(response http.ResponseWriter) error {
		return errors.Join(
			Response.JSONContentType()(response),
			json.MarshalWrite(response, v, options...),
		)
	}
}

// RawJSONBody writes the provided bytes after setting a JSON Content-Type
func (responseSingleton) RawJSONBody(b []byte) ResponseOption {
	return func(response http.ResponseWriter) error {
		return errors.Join(
			Response.JSONContentType()(response),
			Response.BytesBody(b)(response),
		)
	}
}

// BodyFromReader copies from the provided io.Reader into the http.ResponseWriter and
// calls Close() on the reader (if implemented), returning any and all errors.
func (responseSingleton) BodyFromReader(r io.Reader) ResponseOption {
	return func(w http.ResponseWriter) (err error) {
		defer func() {
			closer, ok := r.(io.Closer)
			if ok {
				err = errors.Join(err, closer.Close())
			}
		}()
		_, err = io.Copy(w, r)
		return err
	}
}

// BodyWithAttachment sets headers to deliver the provided content as a downloaded attachment
// with Content-Type set dynamically according to the file extension.
func (responseSingleton) BodyWithAttachment(filename string, content io.Reader) ResponseOption {
	return func(response http.ResponseWriter) error {
		return errors.Join(
			Response.ContentType(mime.TypeByExtension(path.Ext(filename)))(response),
			Response.Header(headerContentDisposition, fmt.Sprintf(attachmentDisposition, filename))(response),
			Response.StatusCode(http.StatusOK)(response),
			Response.BodyFromReader(content)(response),
		)
	}
}

var (
	headerContentType        = "Content-Type"
	headerContentDisposition = "Content-Disposition"

	attachmentDisposition = `attachment; filename="%s"`
	jsonContentType       = "application/json; charset=utf-8"
)
