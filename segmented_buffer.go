package thrift_util

const (
	initialMaxLength     = 64 * 1024
	defaultSegmentLength = 512
)

type SegmentedBuffer struct {
	segments  [][]byte
	segLength int
	pos       int
	total     int
}

func NewByteSegments(segLength int) *SegmentedBuffer {
	s := &SegmentedBuffer{
		segments:  make([][]byte, 0, initialMaxLength/segLength),
		segLength: segLength,
		pos:       -1,
	}
	s.newSegment()
	return s
}

func (b *SegmentedBuffer) newSegment() []byte {
	s := make([]byte, 0, b.segLength)
	b.segments = append(b.segments, s)
	b.pos += 1
	return s
}

func (b *SegmentedBuffer) Append(buf []byte) {
	b.total += len(buf)
	for {
		seg := b.segments[b.pos]
		available := cap(seg) - len(seg)
		left := len(buf)
		if available > 0 {
			size := available
			if left < available {
				size = left
			}
			b.segments[b.pos] = append(seg, buf[:size]...)
		}
		if available >= left {
			return
		}
		b.newSegment()
		buf = buf[available:]
	}
}

func (b *SegmentedBuffer) Buffer() []byte {
	result := make([]byte, 0, b.total)
	for _, seg := range b.segments {
		result = append(result, seg...)
	}
	return result
}
