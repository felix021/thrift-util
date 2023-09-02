package buffer

import "github.com/felix021/thrift-util/buffer/iface"

const (
	defaultInitialSegmentsCap = 16
)

type Segment struct {
	buf      []byte
	borrowed bool
}

func newSegment(buf []byte, borrowed bool) Segment {
	return Segment{
		buf:      buf,
		borrowed: borrowed,
	}
}

type SegmentedBuffer struct {
	segments           []Segment
	segmentLength      int
	pos                int
	total              int
	borrowLimit        int
	allowBorrow        bool
	allocator          iface.Allocator
	initialSegmentsCap int
}

type SegmentedBufferOption func(s *SegmentedBuffer)

func NewSegmentedBuffer(segmentLength int, opts ...SegmentedBufferOption) *SegmentedBuffer {
	b := &SegmentedBuffer{
		segmentLength:      segmentLength,
		allowBorrow:        false,
		pos:                -1,
		allocator:          defaultAllocator,
		initialSegmentsCap: defaultInitialSegmentsCap,
	}
	for _, opt := range opts {
		opt(b)
	}
	b.Initialize()
	return b
}

func WithInitialSegmentsCap(initialSegmentsCap int) SegmentedBufferOption {
	return func(s *SegmentedBuffer) {
		s.initialSegmentsCap = initialSegmentsCap
	}
}

func WithAllocator(allocator iface.Allocator) SegmentedBufferOption {
	return func(s *SegmentedBuffer) {
		s.allocator = allocator
	}
}

func WithBorrowLimit(borrowLimit int) SegmentedBufferOption {
	return func(s *SegmentedBuffer) {
		s.borrowLimit = borrowLimit
		s.allowBorrow = true
	}
}

func WithAlwaysCopy() SegmentedBufferOption {
	return func(s *SegmentedBuffer) {
		s.borrowLimit = -1
		s.allowBorrow = false
	}
}

// Initialize can be used after sync.Pool.Get()
func (b *SegmentedBuffer) Initialize() {
	b.segments = make([]Segment, 0, b.initialSegmentsCap)
	b.newSegment()
}

func (b *SegmentedBuffer) newSegment() []byte {
	seg := b.allocator.Allocate(b.segmentLength)
	b.segments = append(b.segments, newSegment(seg, false))
	b.pos += 1
	return seg
}

func (b *SegmentedBuffer) appendSegment(buf []byte) {
	b.segments = append(b.segments, newSegment(buf, true))
	b.pos += 1
}

// Append appends the []byte given.
// NOTE:
//   1. Do not modify buf util you call Buffer() which will do a final copy.
//   2. If there's really the need, set borrowLimit to -1
func (b *SegmentedBuffer) Append(buf []byte) {
	bufLen := len(buf)
	b.total += bufLen
	if b.allowBorrow && (bufLen) > b.borrowLimit {
		b.appendSegment(buf)
		b.newSegment() // do not append to borrowed buf
		return
	}
	seg := b.segments[b.pos].buf
	available := cap(seg) - len(seg)
	if available < bufLen {
		if bufLen > b.segmentLength { // copy to a new segment
			newBuf := b.allocator.Allocate(bufLen)
			newBuf = append(newBuf, buf...)
			b.appendSegment(newBuf)
			return
		}
		seg = b.newSegment()
	}
	b.segments[b.pos].buf = append(seg, buf...)
}

func (b *SegmentedBuffer) Buffer() []byte {
	result := b.allocator.Allocate(b.total)
	for _, seg := range b.segments {
		result = append(result, seg.buf...)
	}
	return result
}

func (b *SegmentedBuffer) Clean() {
	for _, seg := range b.segments {
		if !seg.borrowed {
			b.allocator.Recycle(seg.buf)
		}
		seg.buf = nil
	}
	b.segments = nil
}
