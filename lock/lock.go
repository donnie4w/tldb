// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package lock

import (
	"errors"
	"fmt"
	"hash/maphash"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/sys"
)

type TlLock[T any] interface {
	Lock(key T)
	UnLock(key T)
}

// func NewTlLock[T any]() TlLock[T] {
// 	return &lockware[T]{mux: &sync.Mutex{}, lm: NewMap[T, *_qlock]()}
// }

type _qlock struct {
	idx int64
	mux *sync.Mutex
}

type lockware[T any] struct {
	mux *sync.Mutex
	lm  *Map[T, *_qlock]
}

func (this *lockware[T]) Lock(key T) {
	if l, ok := this.lm.Get(key); ok {
		atomic.AddInt64(&l.idx, 1)
		l.mux.Lock()
	} else {
		func() {
			this.mux.Lock()
			defer this.mux.Unlock()
			if !this.lm.Has(key) {
				this.lm.Put(key, &_qlock{0, &sync.Mutex{}})
			}
		}()
		this.Lock(key)
	}
}

func (this *lockware[T]) UnLock(key T) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if l, ok := this.lm.Get(key); ok {
		if atomic.AddInt64(&l.idx, -1) <= 0 {
			this.lm.Del(key)
		}
		l.mux.Unlock()
	}
}

// ////////////////////////////////////////

type Numlock struct {
	lockm  *Map[int64, *sync.Mutex]
	muxNum int
}

func NewNumLock(muxNum int) *Numlock {
	ml := &Numlock{NewMap[int64, *sync.Mutex](), muxNum}
	for i := 0; i < muxNum; i++ {
		ml.lockm.Put(int64(i), &sync.Mutex{})
	}
	return ml
}

func (this *Numlock) Lock(key int64) {
	l, _ := this.lockm.Get(int64(uint64(key) % uint64(this.muxNum)))
	l.Lock()
}

func (this *Numlock) Unlock(key int64) {
	l, _ := this.lockm.Get(int64(uint64(key) % uint64(this.muxNum)))
	l.Unlock()
}

/*****************************************************/
type Strlock struct {
	lockm  *Map[int64, *sync.RWMutex]
	muxNum int
}

func NewStrlock(muxNum int) *Strlock {
	ml := &Strlock{NewMap[int64, *sync.RWMutex](), muxNum}
	for i := 0; i < muxNum; i++ {
		ml.lockm.Put(int64(i), &sync.RWMutex{})
	}
	return ml
}

func (this *Strlock) Lock(key string) {
	u := hash([]byte(key))
	l, _ := this.lockm.Get(int64(u % uint64(this.muxNum)))
	l.Lock()
}

func (this *Strlock) Unlock(key string) {
	u := hash([]byte(key))
	l, _ := this.lockm.Get(int64(u % uint64(this.muxNum)))
	l.Unlock()
}

func (this *Strlock) RLock(key string) {
	u := hash([]byte(key))
	l, _ := this.lockm.Get(int64(u % uint64(this.muxNum)))
	l.RLock()
}

func (this *Strlock) RUnlock(key string) {
	u := hash([]byte(key))
	l, _ := this.lockm.Get(int64(u % uint64(this.muxNum)))
	l.RUnlock()
}

// func (this *lockware) getKeyLock(key string) (_r *sync.Mutex) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok {
// 		atomic.AddInt64(&l, 1)
// 	} else {
// 		this.keyQuote[key] = 1
// 	}
// 	if l, ok := this.keyLockMap.Get(key); ok {
// 		return l
// 	} else {
// 		_r = new(sync.Mutex)
// 		this.keyLockMap.Add(key, _r)
// 	}
// 	return
// }

// func (this *lockware) delKeyLock(key string) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok && l > 1 {
// 		atomic.AddInt64(&l, -1)
// 	} else {
// 		// delete(this.keyLockMap, key)
// 		this.keyLockMap.Del(key)
// 		delete(this.keyQuote, key)
// 	}
// 	return
// }

// /////////////////////////////////////////////////////////
// type rwlockware struct {
// 	mux        *sync.Mutex
// 	keyLockMap map[string]*sync.RWMutex
// 	keyQuote   map[string]int64
// }

// func NewRWLockWare() *rwlockware {
// 	return &rwlockware{mux: &sync.Mutex{}, keyLockMap: make(map[string]*sync.RWMutex, 0), keyQuote: make(map[string]int64, 0)}
// }
// func (this *rwlockware) Lock(key string) {
// 	this.getKeyLock(key).Lock()
// }

// func (this *rwlockware) UnLock(key string) {
// 	this.getKeyLock(key).Unlock()
// 	this.delKeyLock(key)
// }

// func (this *rwlockware) RLock(key string) {
// 	this.getKeyLock(key).RLock()
// }

// func (this *rwlockware) RUnlock(key string) {
// 	if l, ok := this.keyLockMap[key]; ok {
// 		l.RUnlock()
// 	}
// 	this.delKeyLock(key)
// }

// func (this *rwlockware) getKeyLock(key string) (_r *sync.RWMutex) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok {
// 		atomic.AddInt64(&l, 1)
// 	} else {
// 		this.keyQuote[key] = 1
// 	}
// 	if l, ok := this.keyLockMap[key]; ok {
// 		return l
// 	} else {
// 		_r = new(sync.RWMutex)
// 		this.keyLockMap[key] = _r
// 	}
// 	return
// }

// func (this *rwlockware) delKeyLock(key string) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok && l > 1 {
// 		atomic.AddInt64(&l, -1)
// 	} else {
// 		delete(this.keyLockMap, key)
// 		delete(this.keyQuote, key)
// 	}
// 	return
// }

// ////////////////////////////////////////////////////////////////////
// func NewTransLockWare() *transLockware {
// 	return &transLockware{mux: &sync.Mutex{}, keyLockMap: make(map[string]*sync.Mutex, 0), keyQuote: make(map[string]int64, 0)}
// }

// type transLockware struct {
// 	mux        *sync.Mutex
// 	keyLockMap map[string]*sync.Mutex
// 	keyQuote   map[string]int64
// }

// func (this *transLockware) LockBatch(key string, keyMap map[string][]byte, delarr []string) {
// 	if key != "" {
// 		this.getKeyLock(key).Lock()
// 	}
// 	if keyMap != nil {
// 		for key := range keyMap {
// 			this.getKeyLock(key).Lock()
// 		}
// 	}
// 	if delarr != nil {
// 		for _, key := range delarr {
// 			this.getKeyLock(key).Lock()
// 		}
// 	}
// }

// func (this *transLockware) UnLockBatch(key string, keyMap map[string][]byte, delarr []string) {
// 	if key == "" {
// 		if l, ok := this.keyLockMap[key]; ok {
// 			l.Unlock()
// 			this.keyUnLock(key)
// 		}
// 	}
// 	if keyMap != nil {
// 		for key := range keyMap {
// 			// this.getKeyLock(key).Unlock()
// 			if l, ok := this.keyLockMap[key]; ok {
// 				l.Unlock()
// 			}
// 			this.keyUnLock(key)
// 		}
// 	}
// 	if delarr != nil {
// 		for _, key := range delarr {
// 			// this.getKeyLock(key).Unlock()
// 			if l, ok := this.keyLockMap[key]; ok {
// 				l.Unlock()
// 			}
// 			this.keyUnLock(key)
// 		}
// 	}
// }

// func (this *transLockware) UnLockAndCommitBatch(keyMap map[string][]byte, delarr []string) {
// 	if keyMap != nil {
// 		for key := range keyMap {
// 			// this.getKeyLock(key).Unlock()
// 			if l, ok := this.keyLockMap[key]; ok {
// 				l.Unlock()
// 			}
// 			this.keyUnLock(key)
// 			this.delKeyLock(key)
// 		}
// 	}
// 	if delarr != nil {
// 		for _, key := range delarr {
// 			// this.getKeyLock(key).Unlock()
// 			if l, ok := this.keyLockMap[key]; ok {
// 				l.Unlock()
// 			}
// 			this.keyUnLock(key)
// 			this.delKeyLock(key)
// 		}
// 	}
// }

// func (this *transLockware) Lock(key string) {
// 	this.getKeyLock(key).Lock()
// }

// func (this *transLockware) UnLock(key string) {
// 	// this.getKeyLock(key).Unlock()
// 	if l, ok := this.keyLockMap[key]; ok {
// 		l.Unlock()
// 	}
// 	this.keyUnLock(key)
// }

// func (this *transLockware) UnLockAndCommit(key string) {
// 	// this.getKeyLock(key).Unlock()
// 	if l, ok := this.keyLockMap[key]; ok {
// 		l.Unlock()
// 	}
// 	this.keyUnLock(key)
// 	this.delKeyLock(key)
// }

// func (this *transLockware) getKeyLock(key string) (_r *sync.Mutex) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok {
// 		atomic.AddInt64(&l, 1)
// 	} else {
// 		this.keyQuote[key] = 1
// 	}
// 	if l, ok := this.keyLockMap[key]; ok {
// 		return l
// 	} else {
// 		_r = new(sync.Mutex)
// 		this.keyLockMap[key] = _r
// 	}
// 	return
// }

// func (this *transLockware) keyUnLock(key string) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok && l > 1 {
// 		atomic.AddInt64(&l, -1)
// 	}
// 	return
// }

// func (this *transLockware) delKeyLock(key string) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if l, ok := this.keyQuote[key]; ok && l == 1 {
// 		delete(this.keyLockMap, key)
// 		delete(this.keyQuote, key)
// 	}
// 	return
// }

/*********************************************************/
type _tsLockPool struct {
	pool     map[string]*sync.Mutex
	keyQuote map[string]int64
	mux      *sync.Mutex
}

func NewTsLockPool() *_tsLockPool {
	return &_tsLockPool{make(map[string]*sync.Mutex, 0), make(map[string]int64, 0), &sync.Mutex{}}
}

func (this *_tsLockPool) GetLock(keys []string) (ts *TsLock) {
	this.mux.Lock()
	defer this.mux.Unlock()
	ts = newTsLock(keys)
	for _, k := range keys {
		if m, ok := this.pool[k]; ok {
			ts.add(m)
		} else {
			m = &sync.Mutex{}
			this.pool[k] = m
			ts.add(m)
		}
		if i, ok := this.keyQuote[k]; ok {
			this.keyQuote[k] = i + 1
		} else {
			this.keyQuote[k] = 1
		}
	}
	return
}

func (this *_tsLockPool) RmLock(keys []string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, k := range keys {
		if i, ok := this.keyQuote[k]; ok {
			if i > 1 {
				this.keyQuote[k] = i - 1
			} else {
				delete(this.keyQuote, k)
				delete(this.pool, k)
			}
		}
	}
	return
}

/*********************************************************/
type TsLock struct {
	locks []*sync.Mutex
	keys  []string
}

func newTsLock(_keys []string) *TsLock {
	return &TsLock{locks: make([]*sync.Mutex, 0), keys: _keys}
}

func (this *TsLock) add(mux *sync.Mutex) {
	this.locks = append(this.locks, mux)
}

func (this *TsLock) Lock() {
	for _, l := range this.locks {
		l.Lock()
	}
}

func (this *TsLock) UnLock() {
	for _, l := range this.locks {
		l.Unlock()
	}
}

func (this *TsLock) GetKey() []string {
	return this.keys
}

/*********************************************************/
type Await[T any] struct {
	m   *Map[int64, chan T]
	mux *Numlock
}

func NewAwait[T any](muxlimit int) *Await[T] {
	return &Await[T]{NewMap[int64, chan T](), NewNumLock(muxlimit)}
}

func (this *Await[T]) Get(idx int64) (ch chan T) {
	defer this.mux.Unlock(idx)
	this.mux.Lock(idx)
	var ok bool
	if ch, ok = this.m.Get(idx); !ok {
		ch = make(chan T, 1)
		this.m.Put(idx, ch)
	}
	return
}

func (this *Await[T]) DelAndClose(idx int64) {
	defer recover()
	defer this.mux.Unlock(idx)
	this.mux.Lock(idx)
	if o, ok := this.m.Get(idx); ok {
		close(o)
		this.m.Del(idx)
	}
}

func (this *Await[T]) DelAndPut(idx int64, v T) {
	defer recover()
	defer this.mux.Unlock(idx)
	this.mux.Lock(idx)
	if o, ok := this.m.Get(idx); ok {
		o <- v
		this.m.Del(idx)
	}
}

/*********************************************************/
type LimitLock struct {
	ch      chan int
	count   int64
	_count  int64
	timeout time.Duration
}

func NewLimitLock(limit int, timeout time.Duration) (_r *LimitLock) {
	ch := make(chan int, limit)
	_r = &LimitLock{ch: ch, timeout: timeout}
	return
}

func (this *LimitLock) Lock() (err error) {
	select {
	case <-time.After(this.timeout):
		err = errors.New(fmt.Sprint(sys.ERR_TIMEOUT))
	case this.ch <- 1:
	}
	atomic.AddInt64(&this.count, 1)
	return
}

func (this *LimitLock) Unlock() {
	<-this.ch
	atomic.AddInt64(&this._count, 1)
}

// concurrency num
func (this *LimitLock) Cc() int64 {
	return this.count - this._count
}

func (this *LimitLock) LockCount() int64 {
	return this.count
}

var seed = maphash.MakeSeed()

func hash(key []byte) uint64 {
	return maphash.Bytes(seed, key)
}
