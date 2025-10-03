package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mdw-go/scuter"
	"github.com/mdw-go/scuter/example/internal/app"
)

type DeleteTaskShell struct {
	*scuter.JSONResponder[*scuter.Errors]
	logger  app.Logger
	handler app.Handler
}

func NewDeleteTaskShell(logger app.Logger, handler app.Handler) *DeleteTaskShell {
	return &DeleteTaskShell{
		logger:        logger,
		handler:       handler,
		JSONResponder: scuter.NewJSONResponder[*scuter.Errors](logger, json.DefaultOptionsV1()),
	}
}

func (this *DeleteTaskShell) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseUint(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		this.Respond(response, http.StatusBadRequest, scuter.NewErrors(errBadRequestInvalidID))
		return
	}
	command := &app.DeleteTaskCommand{ID: id}
	this.handler.Handle(request.Context(), command)
	switch {
	case command.Result.Error == nil:
		this.Respond(response, http.StatusNoContent, nil)
	case errors.Is(command.Result.Error, app.ErrTaskNotFound):
		this.Respond(response, http.StatusNotFound, nil)
	default:
		this.Respond(response, http.StatusInternalServerError, scuter.NewErrors(errInternalServerError))
	}
}

var ( // TODO: serialize these once and write bytes directly thereafter
	errBadRequestInvalidID = scuter.Error{
		Fields:  []string{"id"},
		Name:    "invalid-id",
		Message: "The id was invalid or not supplied.",
	}
)
