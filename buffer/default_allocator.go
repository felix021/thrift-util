package buffer

import "github.com/felix021/thrift-util/buffer/iface"

var (
	defaultAllocator iface.Allocator = &DefaultAllocator{}
)

type DefaultAllocator struct{}

func (d DefaultAllocator) Allocate(size int) []byte {
	return make([]byte, 0, size)
}

func (d DefaultAllocator) Recycle(buf []byte) {
	// nothing to do
}
