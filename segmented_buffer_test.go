package thrift_util

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, success bool, fmt string, args ...interface{}) bool {
	if success {
		return false
	}
	t.Errorf("assert failed: "+fmt, args...)
	return false
}

func TestSegmentedBuffer(t *testing.T) {
	type tCase struct {
		name        string
		input       string
		expectedLen int
	}

	candidates := []tCase{
		{"0", "0123", 1},
		{"1", "4", 1},
		{"2", "56", 1},
		{"3", "7", 1},
		{"4", "89abcdef", 1},
		{"5", "01234567", 2}, // newSeg
		{"6", "89abcdef", 2},
		{"7", "012345678", 3}, // > 8, appendSeg
		{"8", "0123456", 4},
		{"9", "789abcd", 4},
		{"a", "ef", 4},
		{"b", "0", 5},                // newSeg
		{"c", "0123456789abcdef", 6}, // > 8, appendSeg
	}

	sb := NewByteSegments(16, 8)
	assert(t, cap(sb.segments) == initialMaxLength/sb.segLength, "cap(segments)")
	assert(t, len(sb.segments) == 1, "len(segments)")

	var expected []byte
	for _, c := range candidates {
		buf := make([]byte, len(c.input), len(c.input))
		copy(buf, c.input)
		sb.Append(buf)
		assert(t, len(sb.segments) == c.expectedLen, "case[%s]: len(segments) != %d", c.name, c.expectedLen)
		expected = append(expected, buf...)
	}

	got := sb.Buffer()
	assert(t, reflect.DeepEqual(got, expected), "got = %s, expected= %s", got, expected)
}
