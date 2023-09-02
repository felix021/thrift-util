package decoder

import (
	"encoding/binary"
	"errors"

	"github.com/felix021/thrift-util/buffer/iface"
)

var (
	ErrInvalidType = errors.New("invalid type")
)

func FieldSize(tp byte, buf iface.NextBuffer) (size int, err error) {
	switch tp {
	case TYPE_STOP:
		return
	case TYPE_BOOL, TYPE_I8:
		size = 1
	case TYPE_I16:
		size = 2
	case TYPE_I32:
		size = 4
	case TYPE_I64, TYPE_DOUBLE:
		size = 8
	case TYPE_UUID:
		size = 16
	case TYPE_BINARY:
		var n int
		if n, err = readInt(buf); err != nil {
			return
		}
		if err = buf.Skip(n); err != nil {
			return
		}
		return n + 4, nil
	case TYPE_STRUCT:
		return StructSize(buf)
	case TYPE_MAP:
		return MapSize(buf)
	case TYPE_SET, TYPE_LIST:
		return ListSize(buf)
	default:
		return -1, ErrInvalidType
	}
	if size > 0 {
		if err = buf.Skip(size); err != nil {
			return 0, err
		}
	}
	return
}

func StructSize(buf iface.NextBuffer) (size int, err error) {
	var tp []byte
	var fieldSize int
	for {
		if tp, err = buf.Next(1); err != nil {
			return
		}
		if tp[0] == TYPE_STOP {
			return size + 1, nil
		}
		if err = buf.Skip(2); err != nil { // field id takes up 2 bytes
			return
		}
		if fieldSize, err = FieldSize(tp[0], buf); err != nil || fieldSize == 0 {
			return
		}
		size += fieldSize + 3 // 1(tp) + 2(id)
	}
}

func ListSize(buf iface.NextBuffer) (size int, err error) {
	var header []byte
	var tp byte
	var n, valueSize int
	if header, err = buf.Next(5); err != nil {
		return
	}
	tp = header[0]
	n = int(binary.BigEndian.Uint32(header[1:]))
	size += 5 // 1 + 4
	for i := 0; i < n; i++ {
		if valueSize, err = FieldSize(tp, buf); err != nil {
			return
		}
		size += valueSize
	}
	return
}

func MapSize(buf iface.NextBuffer) (size int, err error) {
	var header []byte
	var keyType, valueType byte
	var n, keySize, valueSize int
	if header, err = buf.Next(6); err != nil {
		return
	}
	keyType = header[0]
	valueType = header[1]
	n = int(binary.BigEndian.Uint32(header[2:]))
	size += 6 // 1 + 1 + 4
	for i := 0; i < n; i++ {
		if keySize, err = FieldSize(keyType, buf); err != nil {
			return
		}
		if valueSize, err = FieldSize(valueType, buf); err != nil {
			return
		}
		size += keySize + valueSize
	}
	return
}

func readInt(buf iface.NextBuffer) (n int, err error) {
	var data []byte
	if data, err = buf.Next(4); err != nil {
		return
	}
	return int(binary.BigEndian.Uint32(data)), nil
}
