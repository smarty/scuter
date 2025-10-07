package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/mdw-go/scuter"
	"github.com/mdw-go/scuter/example/internal/app"
)

type DeleteTaskShell struct {
	logger  app.Logger
	handler app.Handler
}

func NewDeleteTaskShell(logger app.Logger, handler app.Handler) *DeleteTaskShell {
	return &DeleteTaskShell{
		logger:  logger,
		handler: handler,
	}
}

func (this *DeleteTaskShell) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := scuter.Flush(response, this.serveHTTP(request))
	if err != nil {
		this.logger.Printf("[WARN] JSON serialization error: %v", err)
	}
}
func (this *DeleteTaskShell) serveHTTP(request *http.Request) scuter.ResponseOption {
	id, err := strconv.ParseUint(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return this.badRequest()
	}

	command := app.DeleteTaskCommand{ID: id}
	this.handler.Handle(request.Context(), &command)

	switch {
	case command.Result.Error == nil:
		return nil
	case errors.Is(command.Result.Error, app.ErrTaskNotFound):
		return scuter.Response.StatusCode(http.StatusNotFound)
	default:
		return this.internalServerError()
	}
}

func (this *DeleteTaskShell) badRequest() scuter.ResponseOption {
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusBadRequest),
		scuter.Response.JSONError(errBadRequestInvalidID),
	)
}
func (this *DeleteTaskShell) internalServerError() scuter.ResponseOption {
	return scuter.Response.With(
		scuter.Response.StatusCode(http.StatusInternalServerError),
		scuter.Response.JSONError(errInternalServerError),
	)
}

var errBadRequestInvalidID = scuter.Error{
	Fields:  []string{"id"},
	Name:    "invalid-id",
	Message: "The id was invalid or not supplied.",
}
