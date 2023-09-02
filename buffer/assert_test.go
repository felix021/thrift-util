package buffer

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, ok bool, format string, args ...interface{}) {
	if !ok {
		t.Errorf(format, args...)
	}
}

func isNil(got interface{}) bool {
	if got == nil {
		return true
	}
	v := reflect.ValueOf(got)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func assertEqual(t *testing.T, expected interface{}, got interface{}, format string, args ...interface{}) {
	info := fmt.Sprintf(format, args...)
	assert(t, reflect.DeepEqual(expected, got), "assertEqual: expected=%v, got=%v, info=%v", expected, got, info)
}

func assertNil(t *testing.T, got interface{}, format string, args ...interface{}) {
	assert(t, isNil(got), "assertNil: got=%v, info=%v", got, fmt.Sprintf(format, args...))
}

func assertNonNil(t *testing.T, got interface{}, format string, args ...interface{}) {
	assert(t, !isNil(got), "assertNonNil: got=%v, info=%v", got, fmt.Sprintf(format, args...))
}
