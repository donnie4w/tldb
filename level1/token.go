// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package level1

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/lock"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

func init() {
	sys.Lock = getToken
	sys.UnLock = tokenroute.backToken
	sys.ForceUnLock = tokenroute.forceRecvToken
	sys.TryLock = tokenroute.tryGetToken
	go processBroadChan()
	go clearToken()
}

var broadChan = make(chan *resp, 1<<8)
var strLock = NewStrlock(1 << 13)

type req struct {
	uuid  int64
	strId int64
}

type resp struct {
	uuid  int64
	strId int64
	token string
	str   string
}

type strToken struct {
	str         string
	waitChan    chan *req
	waitCount   int32
	tokenFlag   int64
	token       string
	status      int8
	mux         *sync.Mutex
	lastUseTime int64
}

func NewStrToken(sl string, status int8, tokenFlag int64) (_r *strToken) {
	_r = &strToken{str: sl, status: status, waitChan: make(chan *req, 1<<15), mux: &sync.Mutex{}, tokenFlag: tokenFlag}
	if _r.tokenFlag == 0 {
		_r.tokenFlag = util.NewTxId()
	}
	_r.lastUseTime = time.Now().Unix()
	return
}

func (this *strToken) newtoken() string {
	this.token = string(append(util.Int64ToBytes(this.tokenFlag), util.Int64ToBytes(time.Now().UnixNano())...))
	return this.token
}

func (this *strToken) getToken() (token string, ok bool) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.status == 0 {
		this.status = 1
		ok = true
		token = this.newtoken()
	}
	this.updateTime()
	return
}

func (this *strToken) register(r *req) (_r bool) {
	logger.Info("register>>>", *r)
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.status == 0 {
		this.status = 1
		broadChan <- &resp{str: this.str, uuid: r.uuid, strId: r.strId, token: this.newtoken()}
		_r = true
	} else {
		atomic.AddInt32(&this.waitCount, 1)
		go func() { this.waitChan <- r }()
		_r = false
	}
	this.updateTime()
	return
}

func (this *strToken) updateTime() {
	this.lastUseTime = time.Now().Unix()
}

func (this *strToken) recvToken(token string) (isEmpty bool, ok bool) {
	defer errRecover()
	logger.Info("recvToken>>>", token)
	this.mux.Lock()
	defer this.mux.Unlock()
	if bs := []byte(token); len(bs) == 16 {
		reflag := util.BytesToInt64(bs[0:8])
		retime := util.BytesToInt64(bs[8:16])
		if this.tokenFlag == reflag {
			if this.token == "" {
				this.status = 0
			} else {
				if retime >= util.BytesToInt64([]byte(this.token)[8:16]) {
					this.status = 0
				}
			}
		}
	} else if token == "" {
		this.status = 0
	}
	logger.Info("this.waitCount>>", this.waitCount)
	logger.Info("this.status>>", this.status)
	if this.status == 0 {
		ok = true
	}
	if this.status == 0 && this.waitCount > 0 {
		this.newtoken()
		this.status = 1
		r := <-this.waitChan
		broadChan <- &resp{str: this.str, uuid: r.uuid, strId: r.strId, token: this.token}
		atomic.AddInt32(&this.waitCount, -1)
	}
	logger.Info("this.waitCount>>", this.waitCount)
	logger.Info("this.status>>", this.status)
	this.updateTime()
	return this.status == 0 && this.waitCount == 0, ok
}

func processBroadChan() {
	for {
		select {
		case resp := <-broadChan:
			func() {
				defer errRecover()
				if resp.uuid > 0 {
					tt := &TokenTrans{ReqType: TOKEN_RESPUUID, Str: resp.str, Status: 0, TokenFlag: 0, Token: resp.token}
					if tt := nodeWare.broadToken(nil, 0, resp.uuid, tt, false, nil); tt == nil || tt.ReqType != TOKEN_ACK {
						logger.Error("TOKEN_RESPUUID err:", resp.str)
						tokenroute.forceRecvToken(resp.str)
						tokenroute.processReqToken(resp.str, &req{uuid: resp.uuid}) //redo this request
					}
				} else {
					logger.Info("resp.token>>", resp.token)
					if err := sys.ReqToken(resp.strId, resp.token); err != nil {
						tokenroute.forceRecvToken(resp.str)
					}
				}
			}()
		}
	}
}

/*****************************************************************/
var tokenroute = &tokenRoute{strMap: NewMap[string, *strToken](), cmap: NewConsistenthash(1 << 13), mux: &sync.Mutex{}}

type tokenRoute struct {
	strMap *Map[string, *strToken]
	cmap   *Consistenthash
	mux    *sync.Mutex
	isInit bool
}

func (this *tokenRoute) _initCmap() {
	this.cmap.Add(nodeWare.GetALLUUID()...)
}

func (this *tokenRoute) del(uuid int64) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.cmap.Del(uuid)
}

func (this *tokenRoute) add(uuid int64) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.isInit {
		this.cmap.Add(uuid)
	}
}

func (this *tokenRoute) init() {
	if this.isInit {
		return
	}
	this.mux.Lock()
	defer this.mux.Unlock()
	this._initCmap()
	defer func() { this.isInit = true }()
	for _, uuid := range nodeWare.GetRemoteUUIDS() {
		tt := &TokenTrans{ReqType: TOKEN_REQ_INIT, Str: "", Status: 0, TokenFlag: 0, Token: ""}
		_r := nodeWare.broadToken(nil, 0, uuid, tt, false, nil)
		if _r.List != nil {
			for _, v := range _r.List {
				this.newStrTokenSaft(v.Str, v.Status, v.TokenFlag)
			}
		}
	}
}

func (this *tokenRoute) _newStrToken(s string, status int8, tokenflag int64) {
	logger.Info("_newStrToken>>>", s)
	if !this.strMap.Has(s) {
		this.strMap.Put(s, NewStrToken(s, status, tokenflag))
	}
}

func (this *tokenRoute) newStrTokenSaft(s string, status int8, tokenflag int64) (_r *strToken) {
	logger.Info("newStrTokenSaft>>>", s)
	strLock.Lock(s)
	defer strLock.Unlock(s)
	this._newStrToken(s, status, tokenflag)
	_r, _ = this.strMap.Get(s)
	return
}

// The client requests a token
func (this *tokenRoute) processReqToken(s string, r *req) {
	logger.Info("processReqToken>>>", s, " ", *r)
	for !this.isInit {
		<-time.After(time.Millisecond)
	}
	strLock.Lock(s)
	defer strLock.Unlock(s)
	if !this.strMap.Has(s) {
		if this.getNode(s) == sys.UUID {
			this._newStrToken(s, 0, 0)
		} else {
			this.strMap.Put(s, this.reqNewStrToken(this.getNode(s), s))
		}
	}
	if sm, ok := this.strMap.Get(s); ok {
		if !sm.register(r) {
			if this.getNode(s) != sys.UUID {
				this.reqNetToken(s)
			}
		}
	}
}

// The network requests  reqNetToken
func (this *tokenRoute) reqNetToken(s string) (err error) {
	logger.Info("reqNetToken>>>", s)
	if node := this.getNode(s); node != sys.UUID {
		tt := &TokenTrans{ReqType: TOKEN_REQ_TOKEN, Str: s, Status: 0, TokenFlag: 0, Token: ""}
		if tt := nodeWare.broadToken(nil, 0, node, tt, false, nil); tt == nil || tt.ReqType != TOKEN_ACK {
			err = errors.New("req failed")
		}
	}
	return
}

// The network requests newStrToken
func (this *tokenRoute) reqNewStrToken(toUUID int64, s string) (_r *strToken) {
	logger.Info("reqNewStrToken>>>", toUUID, " , ", s)
	tt := &TokenTrans{ReqType: TOKEN_REQ_NEW, Str: s, Status: 0, TokenFlag: 0, Token: ""}
	_tt := nodeWare.broadToken(nil, 0, toUUID, tt, false, nil)
	_r = NewStrToken(s, _tt.Status, _tt.TokenFlag)
	return
}

func (this *tokenRoute) forceRecvToken(s string) {
	logger.Info("forceRecvToken>>>", s)
	if st, ok := this.strMap.Get(s); ok {
		this.recvToken(st.str, st.tokenFlag, "", 0, true)
	} else {
		this.recvToken(s, 0, "", 0, true)
	}
}

func (this *tokenRoute) backToken(s string, token string) (ok bool) {
	logger.Info("backToken>>>>", token)
	if s != "" && len([]byte(token)) == 16 {
		_, ok = this.recvToken(s, 0, token, 0, false)
	}
	return
}

// Receive tokens (including node notifications and terminal returns)
func (this *tokenRoute) recvToken(s string, tokenflag int64, token string, uuid int64, force bool) (isEmpty bool, ok bool) {
	logger.Info("recvToken>>>", s, ",token>>", token, " ,uuid>>", uuid)
	this.mux.Lock()
	defer this.mux.Unlock()
	node := this.getNode(s)
	// if uuid > 0 && !this.strMap.Has(s) && node == sys.UUID {
	// 	this.newStrTokenSaft(s, 1, util.BytesToInt64([]byte(token)[:8]))
	// }
	tokenBack := false
	if uuid > 0 || node == sys.UUID {
		if !this.strMap.Has(s) {
			tf := tokenflag
			if tf == 0 && token != "" {
				tf = util.BytesToInt64([]byte(token)[:8])
			}
			this.newStrTokenSaft(s, 1, tf)
		}
		if st, b := this.strMap.Get(s); b {
			isEmpty, ok = st.recvToken(token)
			if isEmpty && node != sys.UUID {
				tokenBack = true
			}
		}
	}

	// if st, _ok := this.strMap.Get(s); _ok {
	// 	if uuid > 0 || node == sys.UUID {
	// 		isEmpty, ok = st.recvToken(token)
	// 		if isEmpty && node != sys.UUID {
	// 			tokenBack = true
	// 		}
	// 	}
	// }
	if (uuid == 0 && node != sys.UUID) || tokenBack {
		tt := &TokenTrans{ReqType: TOKEN_BACK, Str: s, Status: 0, TokenFlag: 0, Token: token, ForceBack: &force}
		_r := nodeWare.broadToken(nil, 0, node, tt, false, nil)
		ok = *_r.BackOk
		logger.Info("TOKEN_BACK>>>>", node)
	}
	return
}

func (this *tokenRoute) getNode(s string) (_r int64) {
	return this.cmap.Get(int64(util.CRC64([]byte(s))))
}

func (this *tokenRoute) tryGetToken(s string) (token string, ok bool) {
	if node := this.getNode(s); node == sys.UUID {
		if !this.strMap.Has(s) {
			this.newStrTokenSaft(s, 0, 0)
		}
		if st, _ := this.strMap.Get(s); st != nil {
			token, ok = st.getToken()
		}
	} else {
		r := &TokenTrans{ReqType: TOKEN_TRYGET_TOKEN, Str: s, Status: 0, TokenFlag: 0, Token: ""}
		tt := nodeWare.broadToken(nil, 0, node, r, false, nil)
		if tt.Token != "" {
			token, ok = tt.Token, true
		}
	}
	return
}

func (this *tokenRoute) processTokenTrans(tt *TokenTrans, fuuid int64) {
	sendId := *tt.SyncId
	switch tt.ReqType {
	case TOKEN_BACK:
		logger.Info("TOKEN_BACK>>>", fuuid)
		force := false
		if tt.ForceBack != nil {
			force = *tt.ForceBack
		}
		_, ok := this.recvToken(tt.Str, tt.TokenFlag, tt.Token, fuuid, force)
		r := &TokenTrans{ReqType: TOKEN_ACK, Str: tt.Str, Status: 0, TokenFlag: 0, Token: "", BackOk: &ok}
		nodeWare.broadToken(nil, sendId, fuuid, r, true, nil)
		return
	case TOKEN_REQ_NEW:
		var st *strToken
		if node := this.getNode(tt.Str); node == sys.UUID {
			st = this.newStrTokenSaft(tt.Str, 0, 0)
		} else {
			st = this.reqNewStrToken(node, tt.Str)
		}
		r := &TokenTrans{ReqType: TOKEN_REQ_NEW, Str: st.str, Status: 1, TokenFlag: st.tokenFlag, Token: ""}
		nodeWare.broadToken(nil, sendId, fuuid, r, true, nil)
		return
	case TOKEN_RESPUUID:
		logger.Info("TOKEN_RESPUUID>>>", fuuid)
		this.recvToken(tt.Str, tt.TokenFlag, tt.Token, fuuid, false)
	case TOKEN_REQ_TOKEN:
		logger.Info("TOKEN_REQ_TOKEN>>>", fuuid)
		r := &req{uuid: fuuid}
		this.processReqToken(tt.Str, r)
	case TOKEN_REQ_INIT:
		r := &TokenTrans{ReqType: TOKEN_REQ_INIT, Str: "", Status: 0, TokenFlag: 0, Token: ""}
		rlist := make([]*TokenTrans, 0)
		this.strMap.Range(func(k string, v *strToken) bool {
			if this.getNode(k) == fuuid {
				t := &TokenTrans{ReqType: TOKEN_REQ_INIT, Str: v.str, Status: v.status, TokenFlag: v.tokenFlag, Token: v.token}
				rlist = append(rlist, t)
			}
			return true
		})
		if len(rlist) > 0 {
			r.List = rlist
		}
		nodeWare.broadToken(nil, sendId, fuuid, r, true, nil)
		return
	case TOKEN_TRYGET_TOKEN:
		r := &TokenTrans{ReqType: TOKEN_TRYGET_TOKEN, Str: tt.Str, Status: 0, TokenFlag: 0, Token: ""}
		r.Token, _ = tokenroute.tryGetToken(tt.Str)
		nodeWare.broadToken(nil, sendId, fuuid, r, true, nil)
		return
	}
	r := &TokenTrans{ReqType: TOKEN_ACK, Str: tt.Str, Status: 0, TokenFlag: 0, Token: ""}
	nodeWare.broadToken(nil, sendId, fuuid, r, true, nil)
}

/************************************************************/
const (
	TOKEN_ACK          int8 = 0
	TOKEN_BACK         int8 = 1
	TOKEN_REQ_NEW      int8 = 2
	TOKEN_RESPUUID     int8 = 3
	TOKEN_REQ_TOKEN    int8 = 4
	TOKEN_REQ_INIT     int8 = 5
	TOKEN_TRYGET_TOKEN int8 = 6
)

func getToken(strid int64, str string) {
	logger.Info("getToken>>>>", strid)
	r := &req{strId: strid}
	tokenroute.processReqToken(str, r)
}

func clearToken() {
	tk := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-tk.C:
			func() {
				defer recover()
				tokenroute.strMap.Range(func(k string, v *strToken) bool {
					freelimit := sys.FREELOCKTIME
					if freelimit <= 60 {
						freelimit = 60
					}
					if v.status == 0 && v.lastUseTime > 0 && v.lastUseTime+freelimit < time.Now().Unix() {
						tokenroute.strMap.Del(k)
					}
					return true
				})
			}()
		}
	}
}
