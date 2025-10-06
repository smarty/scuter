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
	var statusCode int
	var result any

	defer func() {
		err := scuter.SerializeJSON(response, statusCode, result)
		if err != nil {
			this.logger.Printf("[WARN] JSON serialization error: %v", err)
		}
	}()

	id, err := strconv.ParseUint(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		statusCode, result = http.StatusBadRequest, scuter.NewErrors(errBadRequestInvalidID)
		return
	}

	command := app.DeleteTaskCommand{ID: id}
	this.handler.Handle(request.Context(), &command)

	switch {
	case command.Result.Error == nil:
		statusCode = http.StatusBadRequest
	case errors.Is(command.Result.Error, app.ErrTaskNotFound):
		statusCode = http.StatusNotFound
	default:
		statusCode, result = http.StatusInternalServerError, scuter.NewErrors(errInternalServerError)
	}
}

var ( // TODO: serialize these once and write bytes directly thereafter
	errBadRequestInvalidID = scuter.Error{
		Fields:  []string{"id"},
		Name:    "invalid-id",
		Message: "The id was invalid or not supplied.",
	}
)
