package buffer

import (
	"reflect"
	"testing"
)

func TestSegmentedBuffer(t *testing.T) {
	t.Run("initial-segment-cap", func(t *testing.T) {
		sb := NewSegmentedBuffer(8, WithInitialSegmentsCap(4))
		assertEqual(t, 4, cap(sb.segments), "segments cap")
	})

	t.Run("allowBorrow && bufLen > borrowLimit", func(t *testing.T) {
		sb := NewSegmentedBuffer(8, WithBorrowLimit(4))
		buf := []byte("123456")
		sb.Append(buf)
		assert(t, len(sb.segments) == 3, "len") // first empty, appendSeg, newSeg
		assert(t, reflect.DeepEqual(sb.Buffer(), buf), "buffer")
	})

	t.Run("alwaysCopy, segLen=8, append 4", func(t *testing.T) {
		sb := NewSegmentedBuffer(8, WithAlwaysCopy())
		buf := []byte("123456")
		sb.Append(buf)
		assert(t, len(sb.segments) == 1, "len")
		assert(t, reflect.DeepEqual(sb.Buffer(), buf), "buffer")
	})

	t.Run("alwaysCopy, segLen=8, append 10", func(t *testing.T) {
		sb := NewSegmentedBuffer(8, WithAlwaysCopy())
		buf := []byte("1234567890")
		sb.Append(buf)
		assert(t, len(sb.segments) == 2, "len")
		assert(t, reflect.DeepEqual(sb.Buffer(), buf), "buffer")
	})

	t.Run("alwaysCopy, segLen=8, append 4, append 4", func(t *testing.T) {
		sb := NewSegmentedBuffer(8, WithAlwaysCopy())
		buf1 := []byte("1234")
		sb.Append(buf1)
		assert(t, len(sb.segments) == 1, "len")
		assert(t, reflect.DeepEqual(sb.Buffer(), buf1), "buffer")

		buf2 := []byte("1234")
		sb.Append(buf2)
		assert(t, len(sb.segments) == 1, "len")
		assert(t, reflect.DeepEqual(sb.Buffer(), append(buf1, buf2...)), "buffer")
	})

	t.Run("alwaysCopy, segLen=8, append 4, append 6", func(t *testing.T) {
		sb := NewSegmentedBuffer(8, WithAlwaysCopy())
		buf1 := []byte("1234")
		sb.Append(buf1)

		buf2 := []byte("123456")
		sb.Append(buf2)
		assert(t, len(sb.segments) == 2, "len")
		assert(t, reflect.DeepEqual(sb.Buffer(), append(buf1, buf2...)), "buffer")
	})

	t.Run("allocator", func(t *testing.T) {
		allocator := &countedAllocator{}
		sb := NewSegmentedBuffer(8, WithAllocator(allocator), WithBorrowLimit(4)) // new
		assertEqual(t, 1, allocator.count, "count")
		sb.Append([]byte("12"))     // copy
		sb.Append([]byte("123456")) // borrow && new
		sb.Append([]byte("12"))     // copy
		sb.Append([]byte("3456"))   // copy
		assertEqual(t, 2, allocator.count, "count")
		sb.Append([]byte("7890")) // new & copy
		assertEqual(t, 3, allocator.count, "count")
		sb.Clean()
		assertEqual(t, 0, allocator.count, "count")
	})
}

type countedAllocator struct {
	count int
}

func (d *countedAllocator) Allocate(size int) []byte {
	d.count += 1
	return defaultAllocator.Allocate(size)
}

func (d *countedAllocator) Recycle(buf []byte) {
	d.count -= 1
}
