package buffer

import (
	"testing"
)

func TestDefaultAllocator(t *testing.T) {
	for i := 1; i < 100; i++ {
		buf := defaultAllocator.Allocate(i)
		assertEqual(t, 0, len(buf), "len")
		assertEqual(t, i, cap(buf), "cap")
	}
}
