package http

import (
	"net/http"
	"strings"
	"testing"

	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
	"github.com/smarty/scuter"
	"github.com/smarty/scuter/example/internal/app"
)

func TestCreateTaskFixture(t *testing.T) {
	gunit.Run(new(CreateTaskFixture), t)
}

type CreateTaskFixture struct {
	*gunit.Fixture
	*HTTPFixture
}

func (this *CreateTaskFixture) Setup() {
	this.HTTPFixture = &HTTPFixture{Fixture: this.Fixture}
	this.HTTPFixture.Setup()
}

func (this *CreateTaskFixture) TestInvalidJSONRequestBody() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.Body(strings.NewReader("invalid json")),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusBadRequest),
			scuter.Response.JSONError(testErrBadRequestInvalidJSON),
		),
	)
}
func (this *CreateTaskFixture) TestInvalidFields() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(nil),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusUnprocessableEntity),
			scuter.Response.JSONError(testErrMissingDueDate),
			scuter.Response.JSONError(testErrMissingDetails),
		),
	)
}
func (this *CreateTaskFixture) TestNoID() {
	this.app = func(v any) { v.(*app.CreateTaskCommand).Result.ID = 0 }

	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(map[string]any{
				"details":  "Details",
				"due_date": this.now,
			}),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusInternalServerError),
			scuter.Response.JSONError(testErrInternalServerError),
		),
	)
}
func (this *CreateTaskFixture) TestTaskTooHard() {
	this.app = func(v any) { v.(*app.CreateTaskCommand).Result.Error = app.ErrTaskTooHard }

	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(map[string]any{
				"details":  "Details",
				"due_date": this.now,
			}),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusTeapot),
			scuter.Response.JSONError(testErrTaskTooHard),
		),
	)
}
func (this *CreateTaskFixture) TestHappyPath() {
	this.app = func(v any) {
		command := v.(*app.CreateTaskCommand)
		this.So(command.Details, should.Equal, "Details")
		command.Result.ID = 42
	}

	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(map[string]any{
				"details":  "Details",
				"due_date": this.now,
			}),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusCreated),
			scuter.Response.JSONBody(map[string]any{
				"id":      42,
				"details": "Details",
			}),
		),
	)
}
