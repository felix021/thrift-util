package thrift_util

type NextBuffer interface {
	Next(n int) (p []byte, err error)
	NextByte() (p byte, err error)
	Skip(n int) (err error)
}

type ByteBuffer struct {
	buf []byte
	len int
	pos int
}

func NewByteBuffer(buf []byte) *ByteBuffer {
	return &ByteBuffer{
		buf: buf,
		len: len(buf),
		pos: 0,
	}
}

func (b *ByteBuffer) Skip(n int) error {
	b.pos += n
	if b.pos > b.len {
		return ErrOutOfBound
	}
	return nil
}

func (b *ByteBuffer) Next(n int) (p []byte, err error) {
	end := b.pos + n
	if end > b.len {
		return nil, ErrOutOfBound
	}
	p = b.buf[b.pos:end]
	b.pos = end
	return p, nil
}

func (b *ByteBuffer) NextByte() (p byte, err error) {
	if b.pos >= b.len {
		return 0, ErrOutOfBound
	}
	p = b.buf[b.pos]
	b.pos += 1
	return p, nil
}
