package thrift_util

import (
	"encoding/binary"
	"errors"
)

const (
	TYPE_STOP   byte = 0
	TYPE_BOOL   byte = 2
	TYPE_I8     byte = 3
	TYPE_DOUBLE byte = 4
	TYPE_I16    byte = 6
	TYPE_I32    byte = 8
	TYPE_I64    byte = 10
	TYPE_BINARY byte = 11
	TYPE_STRUCT byte = 12
	TYPE_MAP    byte = 13
	TYPE_SET    byte = 14
	TYPE_LIST   byte = 15
	TYPE_UUID   byte = 16
)

var (
	ErrOutOfBound  = errors.New("out of bound")
	ErrInvalidType = errors.New("invalid type")
)

func StructSize(buf NextBuffer) (size int, err error) {
	return skipStruct(buf)
}

func skipField(tp byte, buf NextBuffer) (size int, err error) {
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
	case TYPE_BINARY:
		var n int
		if n, err = readInt(buf); err != nil {
			return
		}
		if err = buf.Skip(n); err != nil {
			return
		}
		return n + 4, nil
	case TYPE_UUID:
		size = 16
	case TYPE_STRUCT:
		size, err = skipStruct(buf)
		return
	case TYPE_MAP:
		size, err = skipMap(buf)
		return
	case TYPE_SET, TYPE_LIST:
		size, err = skipList(buf)
		return
	default:
		return 0, ErrInvalidType
	}
	if size > 0 {
		if err = buf.Skip(size); err != nil {
			return 0, err
		}
	}
	return
}

func skipList(buf NextBuffer) (size int, err error) {
	var tp byte
	var n, valueSize int
	if tp, err = buf.NextByte(); err != nil {
		return
	}
	if n, err = readInt(buf); err != nil {
		return
	}
	size += 5 // 1 + 4
	for i := 0; i < n; i++ {
		if valueSize, err = skipField(tp, buf); err != nil {
			return
		}
		size += valueSize
	}
	return
}

func skipMap(buf NextBuffer) (size int, err error) {
	var keyType, valueType byte
	var n, keySize, valueSize int
	if keyType, err = buf.NextByte(); err != nil {
		return
	}
	if valueType, err = buf.NextByte(); err != nil {
		return
	}
	if n, err = readInt(buf); err != nil {
		return
	}
	size += 6 // 1 + 1 + 4
	for i := 0; i < n; i++ {
		if keySize, err = skipField(keyType, buf); err != nil {
			return
		}
		if valueSize, err = skipField(valueType, buf); err != nil {
			return
		}
		size += keySize + valueSize
	}
	return
}

func readInt(buf NextBuffer) (n int, err error) {
	var data []byte
	if data, err = buf.Next(4); err != nil {
		return
	}
	return int(binary.BigEndian.Uint32(data)), nil
}

func skipStruct(buf NextBuffer) (size int, err error) {
	var tp byte
	var fieldSize int
	for {
		if tp, err = buf.NextByte(); err != nil {
			return
		}
		if tp == TYPE_STOP {
			return size + 1, nil
		}
		if err = buf.Skip(2); err != nil {
			return
		}
		if fieldSize, err = skipField(tp, buf); err != nil || fieldSize == 0 {
			return
		}
		size += 3 + fieldSize
	}
}
