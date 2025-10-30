package http

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/smarty/scuter"
	"github.com/smarty/scuter/example/internal/app"
)

type (
	// CreateTaskModel is intended as a pooled resource that encapsulates all data belonging to this use case.
	CreateTaskModel struct {
		Request struct {
			DueDate time.Time `json:"due_date"`
			Details string    `json:"details"`
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
	pool    *scuter.Pool[*CreateTaskModel]
	logger  app.Logger
	handler app.Handler
}

func NewCreateTaskShell(logger app.Logger, handler app.Handler) *CreateTaskShell {
	return &CreateTaskShell{
		pool:    scuter.NewPool(newCreateTaskModel),
		logger:  logger,
		handler: handler,
	}
}
func newCreateTaskModel() *CreateTaskModel {
	return &CreateTaskModel{Command: &app.CreateTaskCommand{}}
}
func resetCreateTaskModel(result *CreateTaskModel) {
	result.Request.Details = ""
	result.Command.Details = ""
	result.Command.Result.Error = nil
	result.Command.Result.ID = 0
	result.Response.ID = 0
	result.Response.Details = ""
}
func (this *CreateTaskShell) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	model := this.pool.Get()
	defer this.pool.Put(model)
	resetCreateTaskModel(model)
	scuter.Flush(response, this.serveHTTP(request, model))
}
func (this *CreateTaskShell) serveHTTP(request *http.Request, model *CreateTaskModel) (result scuter.ResponseOption) {
	result, ok := scuter.ReadJSONRequestBody(request, &model.Request)
	if !ok {
		return result
	}
	if model.Request.DueDate.IsZero() {
		result = scuter.Response.With(result, scuter.Response.JSONError(errMissingDueDate))
	}
	model.Request.Details = strings.TrimSpace(model.Request.Details)
	if model.Request.Details == "" {
		result = scuter.Response.With(result, scuter.Response.JSONError(errMissingDetails))
	}
	if result != nil {
		return scuter.Response.With(result, scuter.Response.StatusCode(http.StatusUnprocessableEntity))
	}

	model.Command.Details = model.Request.Details
	this.handler.Handle(request.Context(), model.Command)

	switch {
	case model.Command.Result.Error == nil && model.Command.Result.ID > 0:
		return this.ok(model)
	case errors.Is(model.Command.Result.Error, app.ErrTaskTooHard):
		return scuter.Response.JSONErrors(http.StatusTeapot, errTaskTooHard)
	default:
		return scuter.Response.JSONErrors(http.StatusInternalServerError, errInternalServerError)
	}
}

func (this *CreateTaskShell) ok(model *CreateTaskModel) scuter.ResponseOption {
	model.Response.Details = model.Request.Details
	model.Response.ID = model.Command.Result.ID
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusCreated),
		scuter.Response.JSONBody(model.Response),
	)
}
