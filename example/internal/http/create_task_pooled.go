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
		Request  *CreateTaskRequest
		Command  *app.CreateTaskCommand
		Response *scuter.JSONResponse[*CreateTaskResponse]
	}
	CreateTaskRequest struct {
		Details string `json:"details"`
	}
	CreateTaskResponse struct {
		ID       uint64         `json:"id,omitempty"`
		Details  string         `json:"details,omitempty"`
		Failures *scuter.Errors `json:"failures,omitempty"`
	}
)

// CreateTaskShell is intended to be a long-lived, concurrent-safe structure for serving all HTTP requests routed here.
type CreateTaskShell struct {
	*scuter.Pool[*CreateTaskModel]
	*scuter.JSON[*CreateTaskResponse]
	handler app.Handler
}

func NewCreateTaskShell(handler app.Handler) *CreateTaskShell {
	return &CreateTaskShell{
		handler: handler,
		JSON:    &scuter.JSON[*CreateTaskResponse]{},
		Pool: scuter.NewPool(func() *CreateTaskModel {
			result := new(CreateTaskModel)
			result.Request.Details = "."
			result.Command = new(app.CreateTaskCommand)
			result.Command.Details = "."
			result.Command.Result.ID = 42
			result.Command.Result.Error = errors.New(".")
			result.Response = new(scuter.JSONResponse[*CreateTaskResponse])
			result.Response.Content = new(CreateTaskResponse)
			result.Response.Content.ID = 42
			result.Response.Content.Details = "."
			result.Response.Content.Failures = new(scuter.Errors)
			return result
		}),
	}
}
func (this *CreateTaskShell) initModel() *CreateTaskModel {
	result := this.Get()
	result.Request.Details = ""
	result.Command.Details = ""
	result.Command.Result.Error = nil
	result.Command.Result.ID = 0
	result.Response.StatusCode = http.StatusOK
	result.Response.Content.ID = 0
	result.Response.Content.Details = ""
	result.Response.Content.Failures.Errors = result.Response.Content.Failures.Errors[:0]
	return result
}

func (this *CreateTaskShell) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	model := this.initModel()
	defer this.Put(model)
	defer func() { this.RespondResponse(response, model.Response) }()

	if !this.DeserializeJSON(request, &model.Request) {
		this.badRequest(model.Response)
		return
	}
	model.Command.Details = model.Request.Details

	this.handler.Handle(request.Context(), model.Command)

	switch {
	case model.Command.Result.Error == nil && model.Command.Result.ID > 0:
		this.ok(model)
	case errors.Is(model.Command.Result.Error, app.ErrTaskTooHard):
		this.taskTooHard(model.Response)
	default:
		this.internalServerError(model.Response)
	}
}

func (this *CreateTaskShell) badRequest(response *scuter.JSONResponse[*CreateTaskResponse]) {
	response.StatusCode = http.StatusBadRequest
	response.Content.Failures.Errors = append(response.Content.Failures.Errors, errBadRequestInvalidJSON)
}
func (this *CreateTaskShell) ok(model *CreateTaskModel) {
	model.Response.StatusCode = http.StatusCreated
	model.Response.Content.Details = model.Request.Details
	model.Response.Content.ID = model.Command.Result.ID
}
func (this *CreateTaskShell) taskTooHard(response *scuter.JSONResponse[*CreateTaskResponse]) {
	response.StatusCode = http.StatusTeapot
	response.Content.Failures.Errors = append(response.Content.Failures.Errors, errTaskTooHard)
}
func (this *CreateTaskShell) internalServerError(response *scuter.JSONResponse[*CreateTaskResponse]) {
	response.StatusCode = http.StatusInternalServerError
	response.Content.Failures.Errors = append(response.Content.Failures.Errors, errInternalServerError)
}

var ( // TODO: serialize these once and write bytes directly thereafter
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
