package buffer

import (
	"github.com/felix021/thrift-util/buffer/iface"
)

type BytesBuffer struct {
	buf []byte
	len int
	pos int
}

func NewBytesBuffer(buf []byte) *BytesBuffer {
	return &BytesBuffer{
		buf: buf,
		len: len(buf),
		pos: 0,
	}
}

func (b *BytesBuffer) Skip(n int) error {
	b.pos += n
	if b.pos > b.len {
		return iface.ErrOutOfBound
	}
	return nil
}

func (b *BytesBuffer) Next(n int) (p []byte, err error) {
	end := b.pos + n
	if end > b.len {
		return nil, iface.ErrOutOfBound
	}
	p = b.buf[b.pos:end]
	b.pos = end
	return p, nil
}
