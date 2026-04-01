package scuter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/smarty/scuter/internal/should"
)

func TestError(t *testing.T) {
	err := Error{
		Fields:  []string{"hi"},
		ID:      1,
		Name:    "name",
		Message: "message",
	}
	err2 := err
	should.So(t, err.Error(), should.Equal, "message")
	should.So(t, errors.Is(err2, err), should.BeTrue)
	should.So(t, errors.Is(fmt.Errorf("%w", err), err), should.BeTrue)
}
