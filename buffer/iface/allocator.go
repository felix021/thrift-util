package iface

type Allocator interface {
	Allocate(size int) []byte
	Recycle(buf []byte)
}
