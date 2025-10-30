package http

import "github.com/smarty/scuter"

var (
	testErrBadRequestInvalidID = scuter.Error{
		Fields:  []string{"id"},
		Name:    "invalid-id",
		Message: "The id was invalid or not supplied.",
	}
	testErrMissingDueDate = scuter.Error{
		Fields:  []string{"due_date"},
		Name:    "missing-due-date",
		Message: "The due date is required.",
	}
	testErrMissingDetails = scuter.Error{
		Fields:  []string{"details"},
		Name:    "missing-details",
		Message: "The details of the task are required.",
	}
	testErrTaskTooHard = scuter.Error{
		Fields:  []string{"details"},
		ID:      12345,
		Name:    "task-too-hard",
		Message: "the specified task was deemed overly difficult",
	}
	testErrInternalServerError = scuter.Error{
		ID:      54321,
		Name:    "internal-server-error",
		Message: "Internal Server Error",
	}
)
