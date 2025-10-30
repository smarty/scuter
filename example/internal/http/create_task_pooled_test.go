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

func (this *CreateTaskFixture) TestUnsupportedContentType() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.Header("Content-Type", "wrong"),
			scuter.Request.Body(strings.NewReader(`{"valid":"json"}`)),
		),
		scuter.Response.JSONErrors(http.StatusBadRequest, scuter.ErrUnsupportedRequestContentType),
	)
}
func (this *CreateTaskFixture) TestInvalidJSONRequestBody() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.Header("Content-Type", "application/json; charset=utf-8"),
			scuter.Request.Body(strings.NewReader("invalid json")),
		),
		scuter.Response.JSONErrors(http.StatusBadRequest, scuter.ErrInvalidRequestJSONBody),
	)
}
func (this *CreateTaskFixture) TestInvalidFields() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.JSONBody(nil),
		scuter.Response.JSONErrors(http.StatusUnprocessableEntity, testErrMissingDueDate, testErrMissingDetails),
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
		scuter.Response.JSONErrors(http.StatusInternalServerError, testErrInternalServerError),
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
		scuter.Response.JSONErrors(http.StatusTeapot, testErrTaskTooHard),
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
