package buffer

import "testing"

func TestPooledAllocator(t *testing.T) {
	t.Run("invalid max size", func(t *testing.T) {
		allocator := NewPooledAllocator(15)
		assertNil(t, allocator, "max size")
	})

	t.Run("len(pool) and len(poolIndex)", func(t *testing.T) {
		allocator := NewPooledAllocator(16)
		assertEqual(t, 3, len(allocator.pools), "len(pool)")
		assertEqual(t, 17, len(allocator.poolIndex), "len(poolIndex)")

		assertEqual(t, 0, allocator.poolIndex[1], "index#1")
		assertEqual(t, 0, allocator.poolIndex[2], "index#2")
		assertEqual(t, 0, allocator.poolIndex[3], "index#3")
		assertEqual(t, 0, allocator.poolIndex[4], "index#4")
		assertEqual(t, 1, allocator.poolIndex[5], "index#5")
		assertEqual(t, 2, allocator.poolIndex[16], "index#16")
	})

	t.Run("first-allocate", func(t *testing.T) {
		allocator := NewPooledAllocator(16)

		var buf []byte

		buf = allocator.Allocate(1)
		assertEqual(t, 0, len(buf), "len(buf1) #1")
		assertEqual(t, 4, cap(buf), "cap(buf1) #1")

		buf = allocator.Allocate(4)
		assertEqual(t, 0, len(buf), "len(buf) #4")
		assertEqual(t, 4, cap(buf), "cap(buf) #4")

		buf = allocator.Allocate(5)
		assertEqual(t, 0, len(buf), "len(buf) #5")
		assertEqual(t, 8, cap(buf), "cap(buf) #5")

		buf = allocator.Allocate(15)
		assertEqual(t, 0, len(buf), "len(buf) #15")
		assertEqual(t, 16, cap(buf), "cap(buf) #15")

		buf = allocator.Allocate(16)
		assertEqual(t, 0, len(buf), "len(buf) #16")
		assertEqual(t, 16, cap(buf), "cap(buf) #16")

		buf = allocator.Allocate(17)
		assertEqual(t, 0, len(buf), "len(buf) #17")
		assertEqual(t, 17, cap(buf), "cap(buf) #17")
	})

	t.Run("allocate-after-recycle", func(t *testing.T) {
		allocator := NewPooledAllocator(16)
		var buf []byte
		expected := byte('x')

		buf = allocator.Allocate(16)
		buf = append(buf, expected)
		allocator.Recycle(buf)

		buf = allocator.Allocate(15)
		buf = buf[:1]
		assertEqual(t, expected, buf[0], "buf[0]")
	})
}
