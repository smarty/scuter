package http

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
	"time"

	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
	"github.com/smarty/scuter"
)

type HTTPFixture struct {
	*gunit.Fixture
	now time.Time
	ctx context.Context
	app func(any)
}

func NewHTTPFixture(inner *gunit.Fixture) *HTTPFixture {
	return &HTTPFixture{
		Fixture: inner,
		now:     time.Now().Truncate(time.Second),
		ctx:     context.WithValue(inner.T().Context(), "testing", inner.Name()),
		app:     func(any) {},
	}
}
func (this *HTTPFixture) Handle(ctx context.Context, messages ...any) {
	this.So(ctx.Value("testing"), should.Equal, this.Name())
	for _, msg := range messages {
		this.app(msg)
	}
}

func (this *HTTPFixture) Serve(request *http.Request) *httptest.ResponseRecorder {
	if testing.Verbose() {
		requestDump, err := httputil.DumpRequest(request, true)
		this.So(err, should.BeNil)
		for line := range strings.SplitSeq(string(requestDump), "\n") {
			this.Println("> ", line)
		}
	}

	actual := httptest.NewRecorder()
	router := New(this, this)
	router.ServeHTTP(actual, request)

	if testing.Verbose() {
		this.Println()
		responseDump, err := httputil.DumpResponse(actual.Result(), true)
		this.So(err, should.BeNil)
		for line := range strings.SplitSeq(string(responseDump), "\n") {
			this.Println("< ", line)
		}
	}
	return actual
}
func (this *HTTPFixture) AssertFullHTTP(route string, req scuter.RequestOption, res scuter.ResponseOption) {
	fields := strings.Fields(route)
	method, target := fields[0], fields[1]
	actual := this.Serve(scuter.NewTestRequest(this.ctx, method, target, req))
	expected := httptest.NewRecorder()
	scuter.Flush(expected, res)

	this.So(actual.Code, should.Equal, expected.Code)
	this.So(actual.Header(), should.Equal, expected.Header())
	if strings.TrimSpace(actual.Body.String()) == strings.TrimSpace(expected.Body.String()) {
		return
	}
	var actualBody, expectedBody any
	_ = json.Unmarshal(actual.Body.Bytes(), &actualBody)
	_ = json.Unmarshal(expected.Body.Bytes(), &expectedBody)
	this.So(actualBody, should.Equal, expectedBody)
}
