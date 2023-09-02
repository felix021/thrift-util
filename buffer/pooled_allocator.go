package buffer

import "sync"

type PolledAllocator struct {
	pools     []sync.Pool
	maxSize   int
	poolIndex []int
}

// NewPooledAllocator maxSize are expected to be a power of 2
func NewPooledAllocator(maxSize int) *PolledAllocator {
	nPool := GetExponent(maxSize)
	if maxSize < 4 || nPool <= 0 {
		return nil
	}
	allocator := &PolledAllocator{
		maxSize:   maxSize,
		pools:     make([]sync.Pool, nPool-1, nPool-1),
		poolIndex: make([]int, maxSize+1, maxSize+1),
	}
	i, j, n := 0, 0, 4
	for n <= maxSize {
		for ; j <= n; j++ {
			allocator.poolIndex[j] = i
		}
		i, n = i+1, n*2
	}
	return allocator
}

func (p *PolledAllocator) Allocate(size int) []byte {
	var capSize int
	if size <= p.maxSize {
		index := p.poolIndex[size]
		if buf := p.pools[index].Get(); buf != nil {
			return buf.([]byte)[:0]
		}
		capSize = 1 << (index + 2)
	} else {
		capSize = size
	}
	// pool[i] is empty || size > maxSize
	return make([]byte, 0, capSize)
}

func (p *PolledAllocator) Recycle(buf []byte) {
	if size := cap(buf); size <= p.maxSize {
		index := p.poolIndex[size]
		p.pools[index].Put(buf)
	}
}
