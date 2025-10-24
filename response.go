package scuter

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

// Flush applies the options, which may be supplied in any order, to the provide ResponseWriter.
// IMPORTANT: errors that occur from IO operations involving the response body are silently ignored.
func Flush(response http.ResponseWriter, options ...ResponseOption) {
	config := responseConfigs.Get()
	defer responseConfigs.Put(config)
	config.reset(response.Header())
	Response.With(options...)(config)

	response.WriteHeader(config.status)

	if len(config.jsonErrors.Errors) > 0 {
		_ = json.MarshalWrite(response, config.jsonErrors)
	} else if config.dataJSON != nil {
		_ = json.MarshalWrite(response, config.dataJSON)
	} else if config.dataReader != nil {
		config.writeFromReader(response, config.dataReader)
	} else if config.data.Len() > 0 {
		config.writeFromReader(response, &config.data)
	}
}

// ResponseOption is a callback func with an opportunity to modify the *responseConfig.
type ResponseOption func(*responseConfig)

// Response is the 'namespace' for all methods that return a ResponseOption.
var Response responseSingleton

type responseSingleton struct{}

// With returns a 'composite' option which will be the result of calling all options in the provided order.
func (responseSingleton) With(options ...ResponseOption) ResponseOption {
	return func(config *responseConfig) {
		for _, option := range options {
			if option != nil {
				option(config)
			}
		}
	}
}

// Header adds the value associated with key.
func (responseSingleton) Header(key, value string) ResponseOption {
	return func(config *responseConfig) { config.header.Add(key, value) }
}

// ContentType sets the 'Content-Type' header.
func (responseSingleton) ContentType(mime string) ResponseOption {
	return func(config *responseConfig) { config.header.Add(headerContentType, mime) }
}

// JSONContentType sets the 'Content-Type' header to a sensible value representing JSON.
func (responseSingleton) JSONContentType() ResponseOption {
	return func(config *responseConfig) { config.header.Add(headerContentType, jsonContentType) }
}

// StatusCode sets the status code (and writes all headers).
func (responseSingleton) StatusCode(code int) ResponseOption {
	return func(config *responseConfig) { config.status = code }
}

// BytesBody writes the bytes to the ResponseWriter and returns any error.
func (responseSingleton) BytesBody(b []byte) ResponseOption {
	return func(config *responseConfig) {
		_, _ = config.data.Write(b)
	}
}

// JSONBody uses json.MarshalWrite to serialize v to the ResponseWriter using the provided options and returning any error.
func (responseSingleton) JSONBody(v any, options ...json.Options) ResponseOption {
	return func(config *responseConfig) {
		config.header.Set(headerContentType, jsonContentType)
		config.jsonOptions = append(config.jsonOptions, options...)
		config.dataJSON = v
	}
}

func (responseSingleton) JSONError(err Error, options ...json.Options) ResponseOption {
	return func(config *responseConfig) {
		config.header.Set(headerContentType, jsonContentType)
		config.jsonErrors.Append(err)
		config.jsonOptions = append(config.jsonOptions, options...)
	}
}

// BodyFromReader copies from the provided io.Reader into the http.ResponseWriter and
// calls Close() on the reader (if implemented), returning any and all errors.
func (responseSingleton) BodyFromReader(r io.Reader) ResponseOption {
	return func(config *responseConfig) { config.dataReader = r }
}

// BodyWithAttachment sets headers to deliver the provided content as a downloaded attachment
// with Content-Type set dynamically according to the file extension.
func (responseSingleton) BodyWithAttachment(filename string, content io.Reader) ResponseOption {
	return func(config *responseConfig) {
		config.header.Set(headerContentDisposition, fmt.Sprintf(attachmentDisposition, filename))
		config.header.Set(headerContentType, mime.TypeByExtension(filepath.Ext(filename)))
		config.status = http.StatusOK
		config.dataReader = content
	}
}

var (
	headerContentType        = "Content-Type"
	headerContentDisposition = "Content-Disposition"

	attachmentDisposition = `attachment; filename="%s"`
	jsonContentType       = "application/json; charset=utf-8"
)

type responseConfig struct {
	header      http.Header
	status      int
	dataReader  io.Reader
	data        bytes.Buffer
	dataJSON    any
	jsonErrors  *Errors
	jsonOptions []json.Options
}

func (this *responseConfig) writeFromReader(response http.ResponseWriter, reader io.Reader) {
	defer func() {
		closer, ok := reader.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}()
	_, _ = io.Copy(response, reader)
}

func (this *responseConfig) reset(header http.Header) {
	this.header = header
	this.status = http.StatusOK
	this.data.Reset()
	this.dataJSON = nil
	this.dataReader = nil
	this.jsonErrors.Errors = this.jsonErrors.Errors[:0]
	this.jsonOptions = this.jsonOptions[:0]
}

var responseConfigs = NewPool[*responseConfig](func() *responseConfig {
	config := &responseConfig{jsonErrors: NewErrors()}
	config.reset(nil)
	return config
})
