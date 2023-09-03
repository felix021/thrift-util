package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/felix021/thrift-util/buffer"
	"github.com/felix021/thrift-util/buffer/iface"
	"github.com/felix021/thrift-util/decoder"
	"tests/kitex_gen/test"
)

func assert(t *testing.T, success bool, fmt string, args ...interface{}) bool {
	if success {
		return false
	}
	t.Errorf("assert failed: "+fmt, args...)
	return false
}

func mustEncode(v *test.Demo) []byte {
	transport := thrift.NewTMemoryBuffer()
	protocol := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(transport)

	err := v.Write(protocol)
	if err != nil {
		panic(err)
	}
	return transport.Bytes()
}

func testCase(t *testing.T, name string, v *test.Demo) {
	t.Run(name, func(t *testing.T) {
		buf := mustEncode(v)
		size, err := decoder.StructSize(buffer.NewBytesBuffer(buf))
		if !assert(t, err == nil, "err = %v, buf = %v", err, buf) {
			return
		}
		assert(t, size == len(buf), "size not match, expected=%d, actual=%d, buf=%v", len(buf), size, buf)
	})
}

var (
	demoBool1          = true
	demoDouble1        = float64(1)
	demoByte1          = int8(2)
	demoInt16          = int16(3)
	demoInt32          = int32(4)
	demoInt64          = int64(5)
	demoString1        = "string"
	demoStringEmpty    = ""
	demoMapStringInt32 = map[string]int32{"a": 1, "b": 2}
	demoSetByte        = []int8{1, 2, 3}
	demoSetI32         = []int32{1, 2, 3}
	demoListByte       = []int8{1, 2, 3}
	demoListString     = []string{"a", "b", "c"}
	demoMapStringDemo  = map[string]*test.Demo{
		"a": {Int16: &demoInt16},
		"b": {Int32: &demoInt32},
	}
	demoListMapDemo = []map[int64]*test.Demo{
		{
			1: {Int64: &demoInt64},
			2: {String1: &demoString1},
		},
	}
	demoEnum1 = test.EnumType_A

	demoDemo = &test.Demo{
		Bool1:          &demoBool1,
		Double1:        &demoDouble1,
		Byte1:          &demoByte1,
		Int16:          &demoInt16,
		Int32:          &demoInt32,
		Int64:          &demoInt64,
		String1:        &demoString1,
		StructDemo:     &test.Demo{Bool1: &demoBool1},
		MapStringInt32: demoMapStringInt32,
		MapStringDemo:  demoMapStringDemo,
		SetByte:        demoSetByte,
		SetI32:         demoSetI32,
		ListByte:       demoListByte,
		ListString:     demoListString,
		ListMapDemo:    demoListMapDemo,
		Enum1:          &demoEnum1,
	}
)

func TestAll(t *testing.T) {
	testCase(t, "empty", &test.Demo{})
	testCase(t, "bool", &test.Demo{Bool1: &demoBool1})
	testCase(t, "int8", &test.Demo{Byte1: &demoByte1})
	testCase(t, "int16", &test.Demo{Int16: &demoInt16})
	testCase(t, "int32", &test.Demo{Int32: &demoInt32})
	testCase(t, "int64", &test.Demo{Int64: &demoInt64})
	testCase(t, "double", &test.Demo{Double1: &demoDouble1})
	testCase(t, "empty-string", &test.Demo{String1: &demoStringEmpty})
	testCase(t, "string", &test.Demo{String1: &demoString1})
	testCase(t, "struct", &test.Demo{StructDemo: &test.Demo{Bool1: &demoBool1, Int16: &demoInt16}})
	testCase(t, "map-string-int32", &test.Demo{MapStringInt32: demoMapStringInt32})
	testCase(t, "map-string-demo", &test.Demo{MapStringDemo: demoMapStringDemo})
	testCase(t, "set-byte", &test.Demo{SetByte: demoSetByte})
	testCase(t, "set-int32", &test.Demo{SetI32: demoSetI32})
	testCase(t, "list-byte", &test.Demo{ListByte: demoListByte})
	testCase(t, "list-string", &test.Demo{ListString: demoListString})
	testCase(t, "list-map-demo", &test.Demo{ListMapDemo: demoListMapDemo})
	testCase(t, "enum", &test.Demo{Enum1: &demoEnum1})
	testCase(t, "full", demoDemo)
}

var _ iface.NextBuffer = (*skipAppendBuffer)(nil)

type skipAppendBuffer struct {
	in  *buffer.BytesBuffer
	out []byte
}

func newSkipAppendBuffer(buf []byte) *skipAppendBuffer {
	return &skipAppendBuffer{
		in: buffer.NewBytesBuffer(buf),
	}
}

func (s *skipAppendBuffer) Next(n int) (buf []byte, err error) {
	if buf, err = s.in.Next(n); err != nil {
		s.out = append(s.out, buf...)
	}
	return
}

func (s *skipAppendBuffer) Skip(n int) (err error) {
	_, err = s.Next(n)
	return
}

var _ iface.NextBuffer = (*skipAppendBuffer)(nil)

type skipSegmentedAppendBuffer struct {
	in *buffer.BytesBuffer
	sb *buffer.SegmentedBuffer
}

func newSkipSegmentedAppendBuffer(buf []byte, segmentLength int, opts ...buffer.SegmentedBufferOption) *skipSegmentedAppendBuffer {
	return &skipSegmentedAppendBuffer{
		in: buffer.NewBytesBuffer(buf),
		sb: buffer.NewSegmentedBuffer(segmentLength, opts...),
	}
}

func (s *skipSegmentedAppendBuffer) Next(n int) (buf []byte, err error) {
	if buf, err = s.in.Next(n); err != nil {
		s.sb.Append(buf)
	}
	return
}

func (s *skipSegmentedAppendBuffer) Skip(n int) (err error) {
	_, err = s.Next(n)
	return
}

func (s *skipSegmentedAppendBuffer) Buffer() []byte {
	buf := s.sb.Buffer()
	s.sb.Clean()
	return buf
}

func BenchmarkSkipDecoderSmall(b *testing.B) {
	bufSmall := mustEncode(demoDemo)
	b.Logf("len(buffSmall) = %v", len(bufSmall))

	b.Run("small-append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := decoder.StructSize(newSkipAppendBuffer(bufSmall))
			if err != nil {
				b.Errorf("decode failed")
			}
		}
	})

	segLen := 128
	options := []buffer.SegmentedBufferOption{
		buffer.WithBorrowLimit(32),
	}

	b.Run("small-segmented", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := decoder.StructSize(newSkipSegmentedAppendBuffer(bufSmall, segLen, options...))
			if err != nil {
				b.Errorf("decode failed")
			}
		}
	})

	b.Run("small-segmented-pooled-allocator", func(b *testing.B) {
		allocator := buffer.NewPooledAllocator(1024)
		opts := append(options, buffer.WithAllocator(allocator))
		for i := 0; i < b.N; i++ {
			_, err := decoder.StructSize(newSkipSegmentedAppendBuffer(bufSmall, segLen, opts...))
			if err != nil {
				b.Errorf("decode failed")
			}
		}
	})
}

func BenchmarkSkipDecoderAllSizesSmallObjects(b *testing.B) {
	multiplies := []int{1, 2, 4, 8, 16, 32, 64, 128}

	segLen := 128
	options := []buffer.SegmentedBufferOption{
		buffer.WithBorrowLimit(32),
	}

	for _, multiplier := range multiplies {
		var buf []byte
		demo := &test.Demo{
			MapStringDemo: map[string]*test.Demo{},
		}
		for i := 0; i < multiplier; i++ {
			demo.MapStringDemo[strconv.Itoa(i)] = demoDemo
		}
		buf = mustEncode(demo)

		prefix := fmt.Sprintf("size(%d)-", len(buf))

		benchmark(b, prefix, buf, segLen, options)
	}
}

func BenchmarkSkipDecoderAllSizesLongString(b *testing.B) {
	multiplies := []int{1, 2, 4, 8, 16, 32, 64, 128}
	single := string(mustEncode(demoDemo))

	segLen := 128
	options := []buffer.SegmentedBufferOption{
		buffer.WithBorrowLimit(32),
	}

	for _, multiplier := range multiplies {
		var buf []byte
		s := ""
		demo := &test.Demo{String1: &s}
		for i := 0; i < multiplier; i++ {
			*demo.String1 += single
		}
		buf = mustEncode(demo)
		prefix := fmt.Sprintf("size(%d)-", len(buf))
		benchmark(b, prefix, buf, segLen, options)
	}
}

func benchmark(b *testing.B, prefix string, buf []byte, segLen int, options []buffer.SegmentedBufferOption) {
	b.Run(prefix+"append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := decoder.StructSize(newSkipAppendBuffer(buf))
			if err != nil {
				b.Errorf("decode failed")
			}
		}
	})

	b.Run(prefix+"segmented", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := decoder.StructSize(newSkipSegmentedAppendBuffer(
				buf, 128, buffer.WithBorrowLimit(32)))
			if err != nil {
				b.Errorf("decode failed")
			}
		}
	})

	b.Run(prefix+"segmented+pooled-allocator", func(b *testing.B) {
		allocator := buffer.NewPooledAllocator(1024)
		opts := append(options, buffer.WithAllocator(allocator))
		for i := 0; i < b.N; i++ {
			_, err := decoder.StructSize(newSkipSegmentedAppendBuffer(buf, segLen, opts...))
			if err != nil {
				b.Errorf("decode failed")
			}
		}
	})
}
