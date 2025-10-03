package should

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// RunTests is a minimal xUnit-style test runner for Go: https://www.smarty.com/blog/lets-build-xunit-in-go
func RunTests(fixturePointer any, t *testing.T) {
	fixtureType := reflect.TypeOf(fixturePointer)

	for x := 0; x < fixtureType.NumMethod(); x++ {
		testMethodName := fixtureType.Method(x).Name
		if strings.HasPrefix(testMethodName, "Test") {
			instance := reflect.New(fixtureType.Elem())

			setupMethod := instance.MethodByName("Setup")
			if setupMethod.IsValid() {
				setupMethod.Call([]reflect.Value{})
			}

			testMethod := instance.MethodByName(testMethodName)
			callableTest := testMethod.Interface().(func(t *testing.T))
			t.Run(testMethodName, callableTest)
		}
	}
}

type assertion func(actual any, expected ...any) error

func So(t *testing.T, actual any, assertion assertion, expected ...any) {
	if err := assertion(actual, expected...); err != nil {
		t.Helper()
		t.Error(err)
	}
}

var NOT negated

type negated struct{}

func Equal(actual any, expected ...any) error {
	if equalTimes(actual, expected[0]) {
		return nil
	}
	if reflect.DeepEqual(actual, expected[0]) {
		return nil
	}
	return fmt.Errorf("\nExpected: %s\nActual:   %s", format(expected[0]), format(actual))
}
func (negated) Equal(actual any, expected ...any) error {
	if Equal(actual, expected...) != nil {
		return nil
	}
	return fmt.Errorf("\nExpected:     %s\nto not equal: %s\n(but it did)", format(expected[0]), format(actual))
}
func BeTrue(actual any, _ ...any) error          { return Equal(actual, true) }
func BeFalse(actual any, _ ...any) error         { return Equal(actual, false) }
func BeNil(actual any, _ ...any) error           { return Equal(actual, nil) }
func (negated) BeNil(actual any, _ ...any) error { return NOT.Equal(actual, nil) }

func format(v any) string {
	if isTime(v) {
		return fmt.Sprintf("%v", v)
	}
	if v == nil {
		return fmt.Sprintf("%v", v)
	}
	if t := reflect.TypeOf(v); t.Kind() <= reflect.Float64 {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("%#v", v)
}

func equalTimes(a, b any) bool {
	return isTime(a) && isTime(b) && a.(time.Time).Equal(b.(time.Time))
}
func isTime(v any) bool {
	_, ok := v.(time.Time)
	return ok
}
