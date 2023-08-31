package thrift_util

const (
	initialMaxLength = 64 * 1024
)

type SegmentedBuffer struct {
	segments    [][]byte
	segLength   int
	pos         int
	total       int
	noCopyLimit int
}

func NewByteSegments(segLength int, noCopyLimit int) *SegmentedBuffer {
	s := &SegmentedBuffer{
		segments:    make([][]byte, 0, initialMaxLength/segLength),
		segLength:   segLength,
		pos:         -1,
		noCopyLimit: noCopyLimit,
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

func (b *SegmentedBuffer) appendSegment(buf []byte) {
	b.segments = append(b.segments, buf)
	b.pos += 1
}

func (b *SegmentedBuffer) Append(buf []byte) {
	b.total += len(buf)
	if len(buf) > b.noCopyLimit {
		b.appendSegment(buf)
		return
	}
	seg := b.segments[b.pos]
	available := cap(seg) - len(seg)
	if available < len(buf) {
		seg = b.newSegment()
	}
	b.segments[b.pos] = append(seg, buf...)
}

func (b *SegmentedBuffer) Buffer() []byte {
	result := make([]byte, 0, b.total)
	for _, seg := range b.segments {
		result = append(result, seg...)
	}
	return result
}
