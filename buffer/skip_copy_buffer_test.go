package buffer

import "testing"

func TestSkipCopyBuffer(t *testing.T) {
	t.Run("next-success", func(t *testing.T) {
		input := []byte("123456")
		bb := NewBytesBuffer(input)
		scb := NewSkipCopyBuffer(bb, 32)

		buf, err := scb.Next(3)
		assertNil(t, err, "err")
		assertEqual(t, []byte("123"), buf, "buf")
	})

	t.Run("next-fail", func(t *testing.T) {
		input := []byte("123456")
		bb := NewBytesBuffer(input)
		scb := NewSkipCopyBuffer(bb, 32)

		_, err := scb.Next(8)
		assertNonNil(t, err, "err")
	})

	t.Run("skip-success", func(t *testing.T) {
		input := []byte("123456")
		bb := NewBytesBuffer(input)
		scb := NewSkipCopyBuffer(bb, 32)

		err := scb.Skip(3)
		assertNil(t, err, "err")
	})

	t.Run("skip-fail", func(t *testing.T) {
		input := []byte("123456")
		bb := NewBytesBuffer(input)
		scb := NewSkipCopyBuffer(bb, 32)

		err := scb.Skip(8)
		assertNonNil(t, err, "err")
	})

	t.Run("buffer", func(t *testing.T) {
		input := []byte("123456")
		bb := NewBytesBuffer(input)
		scb := NewSkipCopyBuffer(bb, 32)

		_, err := scb.Next(3)
		assertNil(t, err, "err3")

		buf := scb.Buffer()
		assertEqual(t, []byte("123"), buf, "buf3")

		buf, err = scb.Next(3)
		assertNil(t, err, "err6")

		buf = scb.Buffer()
		assertEqual(t, input, buf, "buf6")

		_, err = scb.Next(1)
		assertNonNil(t, err, "err7")
	})
}
