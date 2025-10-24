package http

import (
	"context"
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

func (this *HTTPFixture) assertHTTP(request *http.Request, responseOptions ...scuter.ResponseOption) {
	if testing.Verbose() {
		requestDump, _ := httputil.DumpRequest(request, true)
		for _, line := range strings.Split(string(requestDump), "\n") {
			this.Println("> ", line)
		}
	}

	actual := httptest.NewRecorder()
	router := New(this, this)
	router.ServeHTTP(actual, request)

	if testing.Verbose() {
		responseDump, _ := httputil.DumpResponse(actual.Result(), true)
		for _, line := range strings.Split(string(responseDump), "\n") {
			this.Println("< ", line)
		}
	}

	expected := httptest.NewRecorder()
	_ = scuter.Flush(expected, responseOptions...)
	this.So(actual.Code, should.Equal, expected.Code)
	this.So(actual.Header(), should.Equal, expected.Header())
	this.So(actual.Body.String(), should.Equal, expected.Body.String())
}

func (this *HTTPFixture) TestCreateTask_InvalidJSONRequestBody() {
	this.assertHTTP(
		NewRequest(this.ctx, http.MethodPut, "/tasks",
			RequestOptions.Body(strings.NewReader("not valid json")),
		),
		scuter.Response.StatusCode(http.StatusBadRequest),
		scuter.Response.JSONError(scuter.Error{
			Fields:  []string{"body"},
			Name:    "malformed-request-payload",
			Message: "The body did not contain well-formed data and could not be properly deserialized.",
		}),
	)
}
func (this *HTTPFixture) TestCreateTask_InvalidFields() {
	var empty CreateTaskModel
	this.assertHTTP(
		NewRequest(this.ctx, http.MethodPut, "/tasks",
			RequestOptions.JSONBody(empty.Request),
		),
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
	)
}
