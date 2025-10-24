package http

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
	"github.com/smarty/scuter"
)

func TestHTTPFixture(t *testing.T) {
	gunit.Run(new(HTTPFixture), t)
}

type HTTPFixture struct {
	*gunit.Fixture
	ctx context.Context
}

func (this *HTTPFixture) Handle(ctx context.Context, messages ...any) {
	this.So(ctx.Value("testing"), should.Equal, this.Name())
	for _, msg := range messages {
		_ = msg // TODO
	}
}
func (this *HTTPFixture) Setup() {
	this.ctx = context.WithValue(this.T().Context(), "testing", this.Name())
}

func (this *HTTPFixture) serve(request *http.Request) *httptest.ResponseRecorder {
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
func (this *HTTPFixture) assertFullHTTP(method, target string, req scuter.RequestOption, res scuter.ResponseOption) {
	actual := this.serve(scuter.NewTestRequest(this.ctx, method, target, req))
	expected := httptest.NewRecorder()
	scuter.Flush(expected, res)

	this.So(actual.Code, should.Equal, expected.Code)
	this.So(actual.Header(), should.Equal, expected.Header())
	if strings.Contains(actual.Header().Get("Content-Type"), "json") {
		var actualBody, expectedBody any
		_ = json.Unmarshal(actual.Body.Bytes(), &actualBody)
		_ = json.Unmarshal(expected.Body.Bytes(), &expectedBody)
		this.So(actualBody, should.Equal, expectedBody)
	} else {
		this.So(actual.Body.String(), should.Equal, expected.Body.String())
	}
}

func (this *HTTPFixture) TestCreateTask_InvalidJSONRequestBody() {
	this.assertFullHTTP(http.MethodPut, "/tasks",
		scuter.Request.With(
			scuter.Request.Body(strings.NewReader("not valid json")),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusBadRequest),
			scuter.Response.JSONError(scuter.Error{
				Fields:  []string{"body"},
				Name:    "malformed-request-payload",
				Message: "The body did not contain well-formed data and could not be properly deserialized.",
			}),
		),
	)
}
func (this *HTTPFixture) TestCreateTask_InvalidFields() {
	var empty CreateTaskModel
	this.assertFullHTTP(http.MethodPut, "/tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(empty.Request),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusUnprocessableEntity),
			scuter.Response.JSONError(scuter.Error{
				Fields:  []string{"due_date"},
				Name:    "missing-due-date",
				Message: "The due date is required.",
			}),
			scuter.Response.JSONError(scuter.Error{
				Fields:  []string{"details"},
				Name:    "missing-details",
				Message: "The details of the task are required.",
			}),
		),
	)
}
