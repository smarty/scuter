package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
	"github.com/smarty/scuter"
	"github.com/smarty/scuter/example/internal/app"
)

func TestDeleteTaskFixture(t *testing.T) {
	gunit.Run(new(DeleteTaskFixture), t)
}

type DeleteTaskFixture struct {
	*gunit.Fixture
	*HTTPFixture
}

func (this *DeleteTaskFixture) Setup() {
	this.HTTPFixture = &HTTPFixture{Fixture: this.Fixture}
	this.HTTPFixture.Setup()
}

func (this *DeleteTaskFixture) TestInvalidID() {
	this.assertFullHTTP("DELETE /tasks",
		scuter.Request.With(
			scuter.Request.Query("id", "INVALID"),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusBadRequest),
			scuter.Response.JSONError(testErrBadRequestInvalidID),
		),
	)
}
func (this *DeleteTaskFixture) TestUnrecognizedApplicationError() {
	this.app = func(v any) { v.(*app.DeleteTaskCommand).Result.Error = errors.New("boink") }
	this.assertFullHTTP("DELETE /tasks",
		scuter.Request.Query("id", "42"),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusInternalServerError),
			scuter.Response.JSONError(testErrInternalServerError),
		),
	)
}
func (this *DeleteTaskFixture) TestTaskNotFound() {
	this.app = func(v any) { v.(*app.DeleteTaskCommand).Result.Error = app.ErrTaskNotFound }
	this.assertFullHTTP("DELETE /tasks",
		scuter.Request.Query("id", "42"),
		scuter.Response.StatusCode(http.StatusOK),
	)
}
func (this *DeleteTaskFixture) TestSuccess() {
	this.app = func(v any) {
		command := v.(*app.DeleteTaskCommand)
		this.So(command.ID, should.Equal, 42)
		command.Result.Error = nil
	}
	this.assertFullHTTP("DELETE /tasks",
		scuter.Request.Query("id", "42"),
		scuter.Response.StatusCode(http.StatusOK),
	)
}
