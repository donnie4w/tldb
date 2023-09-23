// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/lock"
	. "github.com/donnie4w/tldb/stub"
)

/*****************************************/
var mergeSyncTxWare = NewMap[int64, Merge[*syncBean, map[int64]int8]]()
var numlock = NewNumLock(1 << 6)

type UUID int64

func (this UUID) syncTxMerge(syncList map[int64]int8) (err error) {
	defer myRecovr()
	i := 10
	for i > 0 {
		i--
		if tc := nodeWare.GetTlContext(int64(this)); tc != nil {
			if err = tc.Conn.SyncTxMerge(context.Background(), syncList); err == nil {
				break
			}
		} else {
			<-time.After(1 * time.Second)
		}
	}
	return
}

func (this UUID) syncTxAdd(sb *syncBean, cancelfunc func() bool, f func(bs map[int64]int8) error) {
	var mst Merge[*syncBean, map[int64]int8]
	var ok bool
	if mst, ok = mergeSyncTxWare.Get(int64(this)); !ok {
		numlock.Lock(int64(this))
		defer numlock.Unlock(int64(this))
		if mst, ok = mergeSyncTxWare.Get(int64(this)); !ok {
			mst = NewMergeSyncTx()
			mergeSyncTxWare.Put(int64(this), mst)
		}
	}
	mst.Add(sb)
	go mst.CallBack(cancelfunc, f)
}

func (this UUID) syncTxDel() {
	mergeSyncTxWare.Del(int64(this))
}

/*********************************************************************/
func NewMergeSyncTx() Merge[*syncBean, map[int64]int8] {
	return &mergeSyncTx[*syncBean, map[int64]int8]{mmap: NewMapL[int64, *syncBean](), mux: &sync.Mutex{}}
}

type syncBean struct {
	SyncId int64
	Result int8
}

type mergeSyncTx[T *syncBean, V map[int64]int8] struct {
	mmap      *MapL[int64, *syncBean]
	seq       int64
	cursor    int64
	mux       *sync.Mutex
	mergeSize int8
}

func (this *mergeSyncTx[T, V]) Add(sb *syncBean) {
	i := atomic.AddInt64(&this.seq, 1)
	this.mmap.Put(i, sb)
}

func (this *mergeSyncTx[T, V]) Del(endId int64) {
	this.mmap.Del(endId)
}

func (this *mergeSyncTx[T, V]) Get(id int64) (_r *syncBean, ok bool) {
	_r, ok = this.mmap.Get(id)
	return
}

func (this *mergeSyncTx[T, V]) Len() int64 {
	return this.mmap.Len()
}

func (this *mergeSyncTx[T, V]) SetZlib(on bool) {
}

func (this *mergeSyncTx[T, V]) MergeSize(size int8) {
	this.mergeSize = size
}

func (this *mergeSyncTx[T, V]) CallBack(cancelfunc func() bool, f func(bs map[int64]int8) error) (err error) {
	if this.mux.TryLock() {
		defer this.mux.Unlock()
	} else {
		return
	}
	for this.Len() > 0 && this.seq > this.cursor && !cancelfunc() {
		if mb, err := this.getMergeBean(); err == nil {
			f(mb)
		}
	}
	return
}

func (this *mergeSyncTx[T, V]) getMergeBean() (_r V, err error) {
	defer myRecovr()
	size := 1 << 20
	_r = make(map[int64]int8, 0)
	start, end := this.cursor+1, this.cursor+1
	for i := start; size > 0 && i <= this.seq; i++ {
		var sb *syncBean
		var ok bool
		for {
			if sb, ok = this.Get(i); ok {
				break
			} else {
				<-time.After(1 * time.Nanosecond)
			}
		}
		size = size - 9
		_r[sb.SyncId] = sb.Result
		end = i
	}
	if len(_r) > 0 {
		for i := start; i <= end; i++ {
			this.Del(i)
		}
		this.cursor = end
	} else {
		err = errors.New("")
	}
	return
}
