package buffer

import "github.com/felix021/thrift-util/buffer/iface"

var _ iface.NextBuffer = (*SkipCopyBuffer)(nil)

func NewSkipCopyBuffer(in iface.NextBuffer, minLength int) *SkipCopyBuffer {
	return &SkipCopyBuffer{
		in:  in,
		out: make([]byte, 0, minLength),
	}
}

type SkipCopyBuffer struct {
	in  iface.NextBuffer
	out []byte
}

func (b *SkipCopyBuffer) Next(n int) (buf []byte, err error) {
	if buf, err = b.in.Next(n); err == nil {
		b.out = append(b.out, buf...)
	}
	return
}

func (b *SkipCopyBuffer) Skip(n int) (err error) {
	_, err = b.Next(n)
	return
}

func (b *SkipCopyBuffer) Buffer() []byte {
	return b.out
}
