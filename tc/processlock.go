// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package tc

import (
	"errors"
	"sync"
	"time"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/lock"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
	"github.com/donnie4w/tlnet"
)

func init() {
	sys.ReqToken = reqToken
	go exTime()
}

var strLock = NewStrlock(1 << 13)

var tokenKeeper = &keeper{NewMap[int64, *tokenBean](), &sync.Mutex{}}

type keeper struct {
	tokenMap *Map[int64, *tokenBean]
	mux      *sync.Mutex
}

func (this *keeper) new(strId int64, hc *tlnet.HttpContext, str string, extime int32, istry bool) (ch chan int8) {
	defer util.Recovr()
	strLock.Lock(str)
	defer strLock.Unlock(str)
	if extime <= 0 {
		extime = 60
	}
	if strId == 0 && hc != nil {
		hc.ResponseBytes(0, []byte{0})
		return
	}
	if st, ok := this.tokenMap.Get(strId); !ok {
		this.mux.Lock()
		defer this.mux.Unlock()
		if !this.tokenMap.Has(strId) {
			if !istry {
				sys.Lock(strId, str)
			}
			ch = make(chan int8)
			this.tokenMap.Put(strId, &tokenBean{0, str, hc, ch, extime, false, istry})
		}
	} else {
		st.hc, st.extime = hc, extime
		ch = st.ch
	}
	return
}

func (this *keeper) get(strId int64) (tb *tokenBean, ok bool) {
	return this.tokenMap.Get(strId)
}

func (this *keeper) del(strId int64) {
	if tb, ok := this.tokenMap.Get(strId); ok {
		tb.close()
		this.tokenMap.Del(strId)
	}
}

func (this *keeper) delExpired(strId int64) {
	if tb, ok := this.tokenMap.Get(strId); ok {
		tb.close()
		if !tb.istry {
			tb.hc.ResponseBytes(0, []byte{0})
		}
		this.tokenMap.Del(strId)
	}
}

type tokenBean struct {
	reqtime    int64
	str        string
	hc         *tlnet.HttpContext
	ch         chan int8
	extime     int32
	_close     bool
	istry      bool
}

func (this *tokenBean) close() {
	defer recover()
	if !this._close {
		this._close = true
		close(this.ch)
	}
}

func reqToken(strId int64, token string) (err error) {
	limit := 3
	bs := []byte(token)
	for limit > 0 {
		limit--
		if tk, ok := tokenKeeper.get(strId); ok {
			var n int
			if n, err = tk.hc.ResponseBytes(0, bs); err == nil && n > 0 {
				tk.reqtime = time.Now().Unix()
				tk.close()
			}
			break
		} else {
			<-time.After(time.Second)
		}
	}
	if limit == 0 {
		err = errors.New("over time")
	}
	return
}

func exTime() {
	tk := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recovr()
				tokenKeeper.tokenMap.Range(func(k int64, v *tokenBean) bool {
					if v.reqtime > 0 && v.reqtime+int64(v.extime) < time.Now().Unix() {
						logger.Error("forceunlock>>>>>>", *v)
						go sys.ForceUnLock(v.str)
						tokenKeeper.delExpired(k)
					} 
					return true
				})
			}()
		}
	}
}
