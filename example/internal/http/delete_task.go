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
	statusCode, body := this.serveHTTP(request)
	err := scuter.SerializeJSON(response, statusCode, body)
	if err != nil {
		this.logger.Printf("[WARN] JSON serialization error: %v", err)
	}
}
func (this *DeleteTaskShell) serveHTTP(request *http.Request) (code int, body any) {
	id, err := strconv.ParseUint(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return http.StatusBadRequest, scuter.NewErrors(errBadRequestInvalidID)
	}

	command := app.DeleteTaskCommand{ID: id}
	this.handler.Handle(request.Context(), &command)

	switch {
	case command.Result.Error == nil:
		return http.StatusBadRequest, nil
	case errors.Is(command.Result.Error, app.ErrTaskNotFound):
		return http.StatusNotFound, nil
	default:
		return http.StatusInternalServerError, scuter.NewErrors(errInternalServerError)
	}
}

var (
	errBadRequestInvalidID = scuter.Error{
		Fields:  []string{"id"},
		Name:    "invalid-id",
		Message: "The id was invalid or not supplied.",
	}
)
