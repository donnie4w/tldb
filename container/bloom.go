// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
package container

import (
	"hash/maphash"
	"sync"
	"sync/atomic"
)

var seed = maphash.MakeSeed()

func hashcode(key []byte) uint32 {
	return uint32(maphash.Bytes(seed, key))
}

type Bloomfilter struct {
	dest       []byte
	dest2      []byte
	nBits      uint32
	k          uint8
	mux        *sync.RWMutex
	limit      int
	bitsPerKey int
	count      int64
}

func NewBloomFilter(limit int) *Bloomfilter {
	bitsPerKey := 10
	limit = limit * 2
	nBits := uint32(limit * bitsPerKey)
	if nBits < 64 {
		nBits = 64
	}
	k := uint8(bitsPerKey * 69 / 100) // 0.69 =~ ln(2)
	if k < 1 {
		k = 1
	} else if k > 30 {
		k = 30
	}
	return &Bloomfilter{
		dest:       _newdest(nBits, k),
		dest2:      _newdest(nBits, k),
		nBits:      nBits,
		k:          k,
		limit:      limit,
		bitsPerKey: bitsPerKey,
		mux:        &sync.RWMutex{},
	}
}

func _newdest(nBits uint32, k uint8) (dest []byte) {
	dest = make([]byte, int((nBits+7)/8)+1)
	dest[(nBits+7)/8] = k
	return dest
}

func (this *Bloomfilter) Add(key []byte) {
	if key != nil {
		this._add(key)
		c := atomic.AddInt64(&this.count, 1)
		if c%int64(this.limit) == int64(this.limit)/2 {
			this.dest = _newdest(this.nBits, this.k)
		}
		if c%int64(this.limit) == int64(this.limit)-1 {
			this.dest2 = _newdest(this.nBits, this.k)
		}
	}
}

func (this *Bloomfilter) _add(key []byte) {
	kh := hashcode(key)
	delta := (kh >> 17) | (kh << 15)
	for j := uint8(0); j < this.k; j++ {
		bitpos := kh % this.nBits
		this.dest[bitpos/8] |= (1 << (bitpos % 8))
		this.dest2[bitpos/8] |= (1 << (bitpos % 8))
		kh += delta
	}
}

func (this *Bloomfilter) Contains(key []byte) bool {
	c := this.count
	if c%int64(this.limit) >= int64(this.limit)/2 && c%int64(this.limit) < int64(this.limit)-1 {
		return contains(this.dest2, key)
	} else {
		return contains(this.dest, key)
	}
}

func contains(filter, key []byte) bool {
	nBytes := len(filter) - 1
	if nBytes < 1 {
		return false
	}
	nBits := uint32(nBytes * 8)
	k := filter[nBytes]
	if k > 30 {
		return true
	}
	kh := hashcode(key)
	delta := (kh >> 17) | (kh << 15)
	for j := uint8(0); j < k; j++ {
		bitpos := kh % nBits
		if (uint32(filter[bitpos/8]) & (1 << (bitpos % 8))) == 0 {
			return false
		}
		kh += delta
	}
	return true
}
