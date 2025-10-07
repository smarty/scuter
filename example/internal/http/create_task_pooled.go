package http

import (
	"errors"
	"net/http"

	"github.com/mdw-go/scuter"
	"github.com/mdw-go/scuter/example/internal/app"
)

type (
	// CreateTaskModel is intended as a pooled resource that encapsulates all data belonging to this use case.
	CreateTaskModel struct {
		Request struct {
			Details string `json:"details"`
		}
		Command  *app.CreateTaskCommand
		Response struct {
			ID      uint64 `json:"id,omitempty"`
			Details string `json:"details,omitempty"`
		}
	}
)

// CreateTaskShell is intended to be a long-lived, concurrent-safe structure for serving all HTTP requests routed here.
type CreateTaskShell struct {
	*scuter.PooledModelFramework[*CreateTaskModel]
	logger  app.Logger
	handler app.Handler
}

func NewCreateTaskShell(logger app.Logger, handler app.Handler) *CreateTaskShell {
	return &CreateTaskShell{
		logger:               logger,
		handler:              handler,
		PooledModelFramework: scuter.NewPooledModelFramework(logger, newCreateTaskModel, resetCreateTaskModel),
	}
}
func resetCreateTaskModel(result *CreateTaskModel) {
	result.Request.Details = ""
	result.Command.Details = ""
	result.Command.Result.Error = nil
	result.Command.Result.ID = 0
	result.Response.ID = 0
	result.Response.Details = ""
}
func newCreateTaskModel() *CreateTaskModel {
	result := new(CreateTaskModel)
	result.Request.Details = "."
	result.Command = new(app.CreateTaskCommand)
	result.Command.Details = "."
	result.Command.Result.ID = 42
	result.Command.Result.Error = errors.New(".")
	result.Response.ID = 42
	result.Response.Details = "."
	return result
}
func (this *CreateTaskShell) serveHTTP(request *http.Request, model *CreateTaskModel) scuter.ResponseOption {
	if err := scuter.DeserializeJSON(request, &model.Request); err != nil {
		return this.badRequest(model)
	}

	model.Command.Details = model.Request.Details
	this.handler.Handle(request.Context(), model.Command)

	switch {
	case model.Command.Result.Error == nil && model.Command.Result.ID > 0:
		return this.ok(model)
	case errors.Is(model.Command.Result.Error, app.ErrTaskTooHard):
		return this.taskTooHard(model)
	default:
		return this.internalServerError(model)
	}
}

func (this *CreateTaskShell) badRequest(model *CreateTaskModel) scuter.ResponseOption {
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusBadRequest),
		scuter.Response.JSONError(errBadRequestInvalidJSON),
	)
}
func (this *CreateTaskShell) ok(model *CreateTaskModel) scuter.ResponseOption {
	model.Response.Details = model.Request.Details
	model.Response.ID = model.Command.Result.ID
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusCreated),
		scuter.Response.JSONBody(model.Response),
	)
}
func (this *CreateTaskShell) taskTooHard(model *CreateTaskModel) scuter.ResponseOption {
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusTeapot),
		scuter.Response.JSONError(errTaskTooHard),
	)
}
func (this *CreateTaskShell) internalServerError(model *CreateTaskModel) scuter.ResponseOption {
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusInternalServerError),
		scuter.Response.JSONError(errInternalServerError),
	)
}

var (
	errBadRequestInvalidJSON = scuter.Error{
		Fields:  []string{"body"},
		Name:    "malformed-request-payload",
		Message: "The body did not contain well-formed data and could not be properly deserialized.",
	}
	errTaskTooHard = scuter.Error{
		Fields:  []string{"details"},
		ID:      12345,
		Name:    "task-too-hard",
		Message: "the specified task was deemed overly difficult",
	}
	errInternalServerError = scuter.Error{
		ID:      54321,
		Name:    "internal-server-error",
		Message: "Internal Server Error",
	}
)
