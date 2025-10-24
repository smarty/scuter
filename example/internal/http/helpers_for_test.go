package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var boink = errors.New("boink")

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewRequest(ctx context.Context, method, target string, options ...requestOption) *http.Request {
	var config requestConfig
	config.query = make(url.Values)
	config.headers = make(http.Header)

	for _, opt := range options {
		opt(&config)
	}

	request := httptest.NewRequestWithContext(ctx, method, target, config.body)
	request.URL.RawQuery = config.query.Encode()
	for key, values := range config.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return request
}

type requestConfig struct {
	query   url.Values
	headers http.Header
	body    io.Reader
}

type (
	requestOption  func(*requestConfig)
	requestOptions struct{}
)

var RequestOptions requestOptions

func (requestOptions) Body(r io.Reader) requestOption {
	return func(c *requestConfig) { c.body = r }
}

func (requestOptions) Header(key, value string) requestOption {
	return func(c *requestConfig) { c.headers.Add(key, value) }
}

func (requestOptions) Query(key, value string) requestOption {
	return func(c *requestConfig) { c.query.Add(key, value) }
}

func (requestOptions) JSONBody(v any) requestOption {
	return func(c *requestConfig) {
		body, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		RequestOptions.Header("Content-Type", "application/json; charset=utf-8")(c)
		RequestOptions.Body(bytes.NewReader(body))(c)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
