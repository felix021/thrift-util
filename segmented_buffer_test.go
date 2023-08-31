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
	sb := NewByteSegments(4)
	assert(t, cap(sb.segments) == initialMaxLength/4, "cap(segments)")
	assert(t, len(sb.segments) == 1, "len(segments)")
	sb.Append([]byte("1234"))
	assert(t, len(sb.segments) == 1, "len(segments)")
	sb.Append([]byte("1"))
	assert(t, len(sb.segments) == 2, "len(segments)")
	sb.Append([]byte("2341"))
	assert(t, len(sb.segments) == 3, "len(segments)")
	sb.Append([]byte("2341234"))
	assert(t, len(sb.segments) == 4, "len(segments)")
	sb.Append([]byte("12341"))
	assert(t, len(sb.segments) == 6, "len(segments)")
	sb.Append([]byte("23412345"))
	assert(t, len(sb.segments) == 8, "len(segments)")
	assert(t, reflect.DeepEqual(sb.Buffer(), []byte("12341234123412341234123412345")), "buf")
}
