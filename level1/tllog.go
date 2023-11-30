// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level1

// import (
// 	"fmt"
// 	"sync"
// 	"sync/atomic"

// 	. "github.com/donnie4w/tldb/contanier"
// 	. "github.com/donnie4w/tldb/key"
// 	. "github.com/donnie4w/tldb/level0"
// 	. "github.com/donnie4w/tldb/util"
// )

// type _tlLog struct {
// 	mux      *sync.RWMutex
// 	_cache   *LimitLinkedMap[string, int8]
// 	seq      int64
// 	cacheSeq int64
// }

// var __tLlog *_tlLog

// func tLlog() *_tlLog {
// 	if __tLlog == nil {
// 		__tLlog = newTlLog()
// 	}
// 	return __tLlog
// }

// func newTlLog() (t *_tlLog) {
// 	t = &_tlLog{mux: &sync.RWMutex{}, _cache: NewLimitLinkedMap[string, int8](1 << 15)}
// 	if i, err := Level0_Log.Get(KeyLevel0.LogSeq()); err == nil {
// 		t.seq = BytesToInt64(i)
// 	}
// 	return
// }

// func (this *_tlLog) addkvWithSeq(seq int64, key string) error {
// 	flushmap := make(map[string][]byte, 0)
// 	if this.seq < seq {
// 		this.seq = seq
// 		flushmap[KeyLevel0.LogSeq()] = Int64ToBytes(seq)
// 	}
// 	flushmap[fmt.Sprint(seq)] = []byte(key)
// 	flushmap[key] = Int64ToBytes(seq)
// 	return Level0_Log.Batch(flushmap, nil)
// }

// func (this *_tlLog) addkv(key string) (err error) {
// 	if this._cache.Has(key) {
// 		return nil
// 	}
// 	this._cache.Put(key, 0)
// 	if !Level0_Log.Has(key) {
// 		flushmap := make(map[string][]byte, 0)
// 		seq := atomic.AddInt64(&this.seq, 1)
// 		flushmap[KeyLevel0.LogSeq()] = Int64ToBytes(seq)
// 		flushmap[fmt.Sprint(seq)] = []byte(key)
// 		flushmap[key] = Int64ToBytes(seq)
// 		err = Level0_Log.Batch(flushmap, nil)
// 	}
// 	return
// }

// func (this *_tlLog) addBatch(put map[string][]byte, del []string) error {
// 	var flushmap map[string][]byte
// 	if put != nil {
// 		flushmap = make(map[string][]byte, 0)
// 		for k := range put {
// 			if this._cache.Has(k) {
// 				continue
// 			}
// 			this._cache.Put(k, 0)
// 			if !Level0_Log.Has(k) {
// 				seq := atomic.AddInt64(&this.seq, 1)
// 				flushmap[fmt.Sprint(seq)] = []byte(k)
// 				flushmap[k] = Int64ToBytes(seq)
// 			}
// 		}
// 		flushmap[KeyLevel0.LogSeq()] = Int64ToBytes(this.seq)
// 	}
// 	var dels []string
// 	if del != nil {
// 		dels = make([]string, 0)
// 		for _, v := range del {
// 			this._cache.Del(v)
// 			if bs, err := Level0_Log.Get(v); err == nil {
// 				dels = append(dels, fmt.Sprint(BytesToInt64(bs)), v)
// 			}
// 		}
// 	}
// 	return Level0_Log.Batch(flushmap, dels)
// }

// func (this *_tlLog) addCachekv(key string, value []byte) error {
// 	seq := atomic.AddInt64(&this.cacheSeq, 1)
// 	lcb := &LogCacheBean{Put: map[string][]byte{key: value}}
// 	return Level0_Log.Put(KeyLevel0.CacheKey(seq), lcb.encode())
// }

// func (this *_tlLog) addCacheBatch(put map[string][]byte, del []string) error {
// 	logger.Info("addCacheBatch===>", put)
// 	flushmap := make(map[string][]byte, 0)
// 	if put != nil {
// 		for k := range put {
// 			seq := atomic.AddInt64(&this.cacheSeq, 1)
// 			lcb := &LogCacheBean{Put: map[string][]byte{KeyLevel0.CacheKey(seq): []byte(k)}}
// 			flushmap[KeyLevel0.CacheKey(seq)] = lcb.encode()
// 		}
// 	}
// 	if del != nil {
// 		for _, v := range del {
// 			seq := atomic.AddInt64(&this.cacheSeq, 1)
// 			lcb := &LogCacheBean{Del: []string{v}}
// 			flushmap[KeyLevel0.CacheKey(seq)] = lcb.encode()
// 		}
// 	}
// 	return Level0_Log.Batch(flushmap, nil)
// }

// func (this *_tlLog) getKey(s, t int64) (_r map[int64]*LogKV) {
// 	_r = make(map[int64]*LogKV, 0)
// 	for i := s; i <= t; i++ {
// 		if bs, err := Level0_Log.Get(fmt.Sprint(i)); err == nil {
// 			k := string(bs)
// 			if d, err := Level0.Get(k); err == nil {
// 				_r[i] = &LogKV{i, k, d}
// 			}
// 		}
// 	}
// 	return
// }

// func (this *_tlLog) saveAndClearCache() (err error) {
// 	logger.Info("saveAndClearCache")
// 	for i := int64(1); i < this.cacheSeq; i++ {
// 		if bs, err := Level0_Log.Get(KeyLevel0.CacheKey(i)); err == nil {
// 			lcb := &LogCacheBean{}
// 			lcb.Decode(bs)
// 			this.addBatch(lcb.Put, lcb.Del)
// 			Level0_Log.Del(KeyLevel0.CacheKey(i))
// 		}
// 	}
// 	this.cacheSeq = 0
// 	return
// }

// // //////////////////////////////////////////////
// type LogCacheBean struct {
// 	Put map[string][]byte
// 	Del []string
// }

// func (this *LogCacheBean) encode() (bs []byte) {
// 	bs, _ = Encoder(this)
// 	return
// }

// func (this *LogCacheBean) Decode(bs []byte) {
// 	Decoder(bs, this)
// 	return
// }
