package buffer

import (
	"testing"

	"github.com/felix021/thrift-util/buffer/iface"
)

func TestBytesBufferSkip(t *testing.T) {
	t.Run("skip-less-than-len", func(t *testing.T) {
		bb := NewBytesBuffer([]byte("12345"))
		assertNil(t, bb.Skip(3), "skip#1")
		assertEqual(t, 3, bb.pos, "n#1")
		assertNil(t, bb.Skip(2), "skip#2")
		assertEqual(t, 5, bb.pos, "n#2")
	})

	t.Run("skip-more-than-len", func(t *testing.T) {
		bb := NewBytesBuffer([]byte("12345"))
		assertNil(t, bb.Skip(3), "skip#1")
		assertEqual(t, 3, bb.pos, "n#1")
		assertEqual(t, iface.ErrOutOfBound, bb.Skip(3), "skip#2")
	})
}

func TestBytesBufferNext(t *testing.T) {
	t.Run("next-less-than-len", func(t *testing.T) {
		bb := NewBytesBuffer([]byte("12345"))

		buf, err := bb.Next(3)
		assertNil(t, err, "skip#1")
		assertEqual(t, []byte("123"), buf, "buf1")

		buf, err = bb.Next(2)
		assertNil(t, err, "skip#2")
		assertEqual(t, []byte("45"), buf, "buf2")
	})

	t.Run("next-more-than-len", func(t *testing.T) {
		bb := NewBytesBuffer([]byte("12345"))

		buf, err := bb.Next(3)
		assertNil(t, err, "skip#1")
		assertEqual(t, []byte("123"), buf, "buf1")

		buf, err = bb.Next(3)
		assertEqual(t, iface.ErrOutOfBound, bb.Skip(3), "skip#2")
	})
}
