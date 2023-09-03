package tests

import (
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/felix021/thrift-util/buffer"
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
