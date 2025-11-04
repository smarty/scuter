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
	this.HTTPFixture = NewHTTPFixture(this.Fixture)
}

func (this *DeleteTaskFixture) TestInvalidID() {
	this.assertFullHTTP("DELETE /tasks",
		scuter.Request.Query("id", "INVALID"),
		scuter.Response.JSONErrors(http.StatusBadRequest, testErrBadRequestInvalidID),
	)
}
func (this *DeleteTaskFixture) TestUnrecognizedApplicationError() {
	this.app = func(v any) { v.(*app.DeleteTaskCommand).Result.Error = errors.New("boink") }
	this.assertFullHTTP("DELETE /tasks",
		scuter.Request.Query("id", "42"),
		scuter.Response.JSONErrors(http.StatusInternalServerError, testErrInternalServerError),
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
