// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package container

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

type LruCache[T any] struct {
	cache *lru.Cache
	lock  *sync.Mutex
}

func NewLruCache[T any](maxEntries int) *LruCache[T] {
	return &LruCache[T]{cache: lru.New(maxEntries), lock: new(sync.Mutex)}
}

func (this *LruCache[T]) Get(key lru.Key) (value T, b bool) {
	defer recover()
	defer this.lock.Unlock()
	this.lock.Lock()
	if v, ok := this.cache.Get(key); ok {
		if v != nil {
			value = v.(T)
		}
		b = true
	}
	return
}

func (this *LruCache[T]) Remove(key lru.Key) {
	defer recover()
	defer this.lock.Unlock()
	this.lock.Lock()
	this.cache.Remove(key)
}

func (this *LruCache[T]) RemoveMulti(keys []string) {
	defer recover()
	defer this.lock.Unlock()
	this.lock.Lock()
	for _, k := range keys {
		this.cache.Remove(k)
	}
}

func (this *LruCache[T]) Add(key lru.Key, value T) {
	defer recover()
	defer this.lock.Unlock()
	this.lock.Lock()
	this.cache.Add(key, value)
}

// //////////////////////////////////////////////////
// type FiFo[K, T any] struct {
// 	arr   *list.List
// 	m     *Map[any, T]
// 	mux   *sync.Mutex
// 	limit int
// }

// func NewFiFo[K, T any](_limit int) *FiFo[K, T] {
// 	return &FiFo[K, T]{list.New(), NewMap[any, T](), &sync.Mutex{}, _limit}
// }

// func (this *FiFo[K, T]) Add(key K, value T) {
// 	this.arr.PushFront(key)
// 	for this.arr.Len() > this.limit {
// 		if e := this.arr.Back(); e != nil {
// 			this.m.Del(e.Value)
// 			this.remove(e)
// 		}
// 	}
// 	this.m.Put(key, value)
// }

// func (this *FiFo[K, T]) remove(e *list.Element) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	this.arr.Remove(e)
// }

// func (this *FiFo[K, T]) Get(key K) (_r T, ok bool) {
// 	return this.m.Get(key)
// }

// func (this *FiFo[K, T]) Has(key K) (ok bool) {
// 	return this.m.Has(key)
// }

// /////////////////////////////////
// type List[T any] struct {
// 	l   *list.List
// 	mux *sync.RWMutex
// }

// func NewList[T any]() *List[T] {
// 	return &List[T]{list.New(), &sync.RWMutex{}}
// }

// func (this *List[T]) Len() int {
// 	return this.l.Len()
// }

// func (this *List[T]) PushFront(t T) {
// 	this.l.PushFront(t)
// }

// func (this *List[T]) PushBack(t T) {
// 	this.l.PushBack(t)
// }

// func (this *List[T]) Front() (t T) {
// 	if v := this.l.Front(); v != nil && v.Value != nil {
// 		t = (v.Value).(T)
// 	}
// 	return
// }

// // if you call AddFront() in your func(t T), it may cause an infinite loop
// func (this *List[T]) BackPrev(f func(t T) bool) {
// 	for e := this.l.Back(); e != nil; e = e.Prev() {
// 		if e.Value != nil {
// 			t := (e.Value).(T)
// 			if !f(t) {
// 				break
// 			}
// 		}
// 	}
// }

// // if you call AddBack() in your func(t T), it may cause an infinite loop
// func (this *List[T]) FrontNext(f func(t T) bool) {
// 	for e := this.l.Front(); e != nil; e = e.Next() {
// 		if e.Value != nil {
// 			t := (e.Value).(T)
// 			if !f(t) {
// 				break
// 			}
// 		}
// 	}
// }

// func (this *List[T]) GetBack() (t T) {
// 	if v := this.l.Back(); v != nil && v.Value != nil {
// 		t = (v.Value).(T)
// 	}
// 	return
// }

// func (this *List[T]) DelAndLoadBack() (t T) {
// 	if v := this.l.Back(); v != nil && v.Value != nil {
// 		t = (v.Value).(T)
// 		this.remove(v)
// 	}
// 	return
// }

// func (this *List[T]) DelBack() {
// 	if v := this.l.Back(); v != nil && v.Value != nil {
// 		this.remove(v)
// 	}
// }

// func (this *List[T]) remove(e *list.Element) {
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	this.l.Remove(e)
// }

func myRecovr() {
	if err := recover(); err != nil {
	}
}
