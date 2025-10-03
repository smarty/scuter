package app

import (
	"context"
	"errors"
)

var (
	ErrTaskTooHard  = errors.New("task too hard")
	ErrTaskNotFound = errors.New("task not found")
)

type Logger interface {
	Printf(format string, args ...any)
}
type Handler interface {
	Handle(context.Context, ...any)
}

type CreateTaskCommand struct {
	Details string
	Result  struct {
		ID    uint64
		Error error
	}
}
type DeleteTaskCommand struct {
	ID     uint64
	Result struct {
		Error error
	}
}

type Application struct{}

func (this *Application) Handle(ctx context.Context, messages ...any) {
	for _, message := range messages {
		switch message := message.(type) {
		case *CreateTaskCommand:
			message.Result.Error = ErrTaskTooHard
		case *DeleteTaskCommand:
			message.Result.Error = nil
		}
	}
}
