// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level1

import (
	"context"
	"strings"
	"sync"
	"time"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"
)

/********************************************/
type PON_TYPE int8

const (
	PON  PON_TYPE = 1
	PON2 PON_TYPE = 2
)

var delNodeMap = NewMapL[string, byte]()

type _nodeWare struct {
	tlMap        *MapL[*tlContext, int8]
	_remoteUUIDS []int64
}

var nodeWare = NewNodeWare()

func NewNodeWare() (_r *_nodeWare) {
	_r = new(_nodeWare)
	_r.tlMap = NewMapL[*tlContext, int8]()
	return
}

func (this *_nodeWare) add(tc *tlContext) (err error) {
	logger.Warn("add tc:", tc.RemoteUuid)
	if this.tlMap.Has(tc) {
		err = Errors(sys.ERR_UUID_REUSE)
		return
	}
	if sys.UUID == tc.RemoteUuid || this.hasUUID(tc.RemoteUuid) {
		logError.Error("AddNode and Close :", tc.RemoteUuid, " ", tc.RemoteAddr)
		tc._do_reconn = true
		tc.Close()
		err = Errors(sys.ERR_UUID_REUSE)
		return
	}
	delNodeMap.Del(tc.RemoteAddr)
	this.tlMap.Put(tc, 1)
	if ss := strings.Split(tc.RemoteAddr, ":"); len(ss) == 2 && ss[0] != "" {
		tc.RemoteHost = ss[0]
	}
	this._remoteUUIDS = this.remoteRunUUID()
	tokenroute.add(tc.RemoteUuid)
	return
}

func (this *_nodeWare) has(tc *tlContext) bool {
	return this.tlMap.Has(tc)
}

func (this *_nodeWare) hasUUID(uuid int64) (b bool) {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if uuid == k.RemoteUuid {
			b = true
			return false
		}
		return true
	})
	return
}

func (this *_nodeWare) setStat(uuid int64, _stat sys.STATTYPE, timenano int64) {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if uuid == k.RemoteUuid {
			k.stat = _stat
		}
		return true
	})
	this._remoteUUIDS = this.remoteRunUUID()
	return
}

func (this *_nodeWare) del(tc *tlContext) {
	if !tc.isClose {
		tc.Close()
	}
	this._remoteUUIDS = this.remoteRunUUID()
	this.tlMap.Del(tc)
	if !this.hasUUID(tc.RemoteUuid) {
		tokenroute.del(tc.RemoteUuid)
	}
}

func (this *_nodeWare) delAndNoReconn(tc *tlContext) {
	tc._do_reconn = true
	this.del(tc)
}

// remove node and no reconnect
func (this *_nodeWare) removeNode(uuid int64) {
	// this.delMap.Add(uuid, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if uuid == k.RemoteUuid {
			// this.Del(k)
			this.delAndNoReconn(k)
		}
		return true
	})
}

func (this *_nodeWare) LenByAddr(addr string) (i int32) {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if addr == k.RemoteAddr {
			i++
		}
		return true
	})
	return
}

// all the run stat node
func (this *_nodeWare) GetAllRunUUID() (_r []int64) {
	_r = this.GetRemoteRunUUID()
	if sys.IsRUN() {
		_r = append(_r, sys.UUID)
	}
	return
}

// the remote run stat node
func (this *_nodeWare) GetRemoteRunUUID() (_r []int64) {
	if this._remoteUUIDS == nil {
		this._remoteUUIDS = this.remoteRunUUID()
	}
	return this._remoteUUIDS
}

// the remote run stat node
func (this *_nodeWare) remoteRunUUID() (_r []int64) {
	_r = make([]int64, 0)
	_m := make(map[int64]int8, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if k.stat == sys.RUN {
			_m[k.RemoteUuid] = 0
		}
		return true
	})
	for k := range _m {
		_r = append(_r, k)
	}
	return
}

func (this *_nodeWare) GetUUIDNode() (_r *Node) {
	_r = &Node{UUID: sys.UUID, Addr: sys.ADDR, Stat: int8(sys.SYS_STAT), MqAddr: sys.MQADDR, AdminAddr: sys.WEBADMINADDR, CliAddr: sys.CLIADDR}
	_m := make(map[int64]string, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		_m[k.RemoteUuid] = k.RemoteAddr
		return true
	})
	_r.Nodekv = _m
	return
}

func (this *_nodeWare) GetRemoteUUIDS() (_r []int64) {
	_r = make([]int64, 0)
	_m := make(map[int64]int8, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if _, ok := _m[k.RemoteUuid]; !ok {
			_m[k.RemoteUuid] = 0
			_r = append(_r, k.RemoteUuid)
		}
		return true
	})
	return
}

func (this *_nodeWare) GetALLUUID() []int64 {
	return append(this.GetRemoteUUIDS(), sys.UUID)
}

func (this *_nodeWare) GetAllTlContext() (tcs []*tlContext) {
	tcs = make([]*tlContext, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		tcs = append(tcs, k)
		return true
	})
	return
}

func (this *_nodeWare) GetTlContext(uuid int64) (tc *tlContext) {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if uuid == k.RemoteUuid {
			tc = k
			return false
		}
		return true
	})
	return
}

func (this *_nodeWare) close() {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		k._do_reconn = true
		this.del(k)
		return true
	})
}

func (this *_nodeWare) IsClusRun() (b bool) {
	if n := len(this.GetAllRunUUID()); n > 0 && n >= sys.CLUSTER_NUM {
		b = true
	}
	return
}

func (this *_nodeWare) getRemoteNodes() (_r []*RemoteNode) {
	_r = make([]*RemoteNode, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		r := &RemoteNode{Addr: k.RemoteAddr, UUID: k.RemoteUuid, Host: k.RemoteHost, AdminAddr: k.RemoteAdminAddr, MQAddr: k.RemoteMqAddr, Stat: int8(k.stat), CliAddr: k.RemoteCliAddr}
		_r = append(_r, r)
		return true
	})
	return
}

func (this *_nodeWare) broadcastSync(pb *PonBean, pontype PON_TYPE, to time.Duration, alltime bool) (exce *sys.Exception) {
	exce = sys.NewException()
	if this.tlMap.Len() > 0 {
		var wg sync.WaitGroup
		uuidMap := make(map[int64]int8, 0)
		this.tlMap.Range(func(k *tlContext, _ int8) bool {
			if _, ok := uuidMap[k.RemoteUuid]; !ok {
				wg.Add(1)
				GoPool.Go(func() { this.sendPb(exce, pb, pontype, k.RemoteUuid, &wg, to, alltime) })
				uuidMap[k.RemoteUuid] = 0
			}
			return true
		})
		wg.Wait()
	}
	return
}

func (this *_nodeWare) broadcast(pb *PonBean, pontype PON_TYPE, t time.Duration, mustRun bool) {
	exce := sys.NewException()
	if this.tlMap.Len() > 0 {
		uuidMap := make(map[int64]int8, 0)
		this.tlMap.Range(func(k *tlContext, _ int8) bool {
			if mustRun && k.stat != sys.RUN {
				return true
			}
			if _, ok := uuidMap[k.RemoteUuid]; !ok {
				GoPool.Go(func() { this.sendPb(exce, pb, pontype, k.RemoteUuid, nil, t, false) })
				uuidMap[k.RemoteUuid] = 0
			}
			return true
		})
	}
}

func (this *_nodeWare) singleBroad(pb *PonBean, uuid int64, pontype PON_TYPE, t time.Duration, alltime bool) (err error) {
	exce := sys.NewException()
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if k.RemoteUuid == uuid {
			err = this.sendPb(exce, pb, pontype, k.RemoteUuid, nil, t, alltime)
			return false
		}
		return true
	})
	err = exce.Error()
	return
}

func (this *_nodeWare) singleBroadAlltime(pb *PonBean, uuid int64, pontype PON_TYPE, t time.Duration, isBack bool) {
	this.sendPbAlltime(pb, pontype, uuid, t, isBack)
}

// func (this *_nodeWare) sendPb(exce *sys.Exception, pb *PonBean, pontype PON_TYPE, uuid, sendId int64, wg *sync.WaitGroup, toTime time.Duration) (err error) {
// 	defer myRecovr()
// 	if sendId == 0 {
// 		sendId = newTxId()
// 	}
// 	ch := await.Get(sendId)
// 	defer await.DelAndClose(sendId)
// 	err = this._sendPb(pb, pontype, uuid, sendId)
// 	if err != nil {
// 		if Time().UnixNano() > int64(toTime) {
// 			exce.Code = sys.ERR_TIMEOUT
// 			exce.Put(uuid)
// 			return
// 		}
// 		<-time.After(3 * time.Second)
// 		this.sendPb(exce, pb, pontype, uuid, sendId, wg, toTime)
// 	} else {
// 		var C <-chan time.Time
// 		if toTime > 0 {
// 			C = time.After(toTime - time.Duration(Time().UnixNano()))
// 		} else {
// 			C = make(<-chan time.Time, 1)
// 		}
// 		for {
// 			select {
// 			case r := <-ch:
// 				if r == 1 {
// 					exce.Code = sys.ERR_BATCHFAIL
// 					exce.Put(uuid)
// 				}
// 				goto END
// 			case <-time.After(sys.ReadTimeout):
// 				logger.Warn("re sendPb,pb.Txid:", pb.Txid, ",pb.Ptype:", pb.Ptype, ",sendId:", sendId)
// 				this._sendPb(pb, pontype, uuid, sendId)
// 			case <-C:
// 				exce.Code = sys.ERR_TIMEOUT_2
// 				exce.Put(uuid)
// 				goto END
// 			}
// 		}
// 	END:
// 		if wg != nil {
// 			wg.Done()
// 		}
// 	}
// 	return
// }

func (this *_nodeWare) sendPb(exce *sys.Exception, pb *PonBean, pontype PON_TYPE, uuid int64, wg *sync.WaitGroup, timeout time.Duration, alltime bool) (err error) {
	defer errRecover()
	spTime := 5 * time.Second
	for timeout > 0 {
		sendId := newTxId()
		ch := await.Get(sendId)
		defer await.DelAndClose(sendId)
		if tc := this.GetTlContext(uuid); tc != nil {
			this._sendPb(pb, pontype, uuid, sendId)
			select {
			case r := <-ch:
				if r == 1 {
					exce.Code = sys.ERR_BATCHFAIL
					exce.Put(uuid)
				}
				goto END
			case <-time.After(spTime):
				timeout = timeout - spTime
			}
		} else {
			<-time.After(spTime)
			timeout = timeout - spTime
		}
	}
END:
	if timeout <= 0 {
		exce.Code = sys.ERR_TIMEOUT
		exce.Put(uuid)
		err = Errors(sys.ERR_TIMEOUT)
		if alltime {
			go this.sendPbAlltime(pb, pontype, uuid, 240*time.Hour, true)
		}
	}
	if wg != nil {
		wg.Done()
	}
	return
}

func (this *_nodeWare) sendPbAlltime(pb *PonBean, pontype PON_TYPE, uuid int64, timeout time.Duration, isBack bool) {
	defer errRecover()
	spTime := 5 * time.Second
	backcount := 0
	for timeout > 0 {
		if tc := this.GetTlContext(uuid); tc != nil {
			sendId := newTxId()
			defer await.DelAndClose(sendId)
			ch := await.Get(sendId)
			this._sendPb(pb, pontype, uuid, sendId)
			select {
			case f := <-ch:
				backcount++
				if isBack {
					if f == 0 || backcount > 3 {
						tlog.RemoveBackLogUUid(pb.Txid, uuid)
						goto END
					}
				} else {
					goto END
				}
			case <-time.After(spTime):
				timeout = timeout - spTime
			}
		} else {
			<-time.After(spTime)
			timeout = timeout - spTime
		}
	}
END:
}

func (this *_nodeWare) _sendPb(pb *PonBean, pontype PON_TYPE, uuid, sendId int64) (err error) {
	defer errRecover()
	if tc := this.GetTlContext(uuid); tc != nil {
		switch pontype {
		case PON:
			err = tc.Conn.Pon(context.Background(), EncodePonBean(pb).Bytes(), sendId)
		case PON2:
			err = tc.Conn.Pon2(context.Background(), pb, sendId)
		}
	} else {
		logError.Error("send, uuid:", uuid, " not found")
		err = Errors(sys.ERR_NODE_NOFOUND)
	}
	return
}

func (this *_nodeWare) commitTx(uuid, sendId, txid int64, timeout time.Duration) (_r int8, err error) {
	defer errRecover()
	if tc := this.GetTlContext(uuid); tc != nil {
		if sendId == 0 {
			sendId = newTxId()
		}
		ch := await.Get(sendId)
		defer await.DelAndClose(sendId)
		tc.Conn.CommitTx(context.Background(), sendId, txid)
		select {
		case _r = <-ch:
		case <-time.After(timeout):
			err = Errors(sys.ERR_TIMEOUT)
		}
	} else {
		err = Errors(sys.ERR_NODE_NOFOUND)
	}
	return
}

func (this *_nodeWare) commitTx2(uuid, sendId, txid int64, commit bool, timeout time.Duration) (err error) {
	defer errRecover()
	if tc := this.GetTlContext(uuid); tc != nil {
		if sendId == 0 {
			sendId = newTxId()
		}
		ch := await.Get(sendId)
		defer await.DelAndClose(sendId)
		tc.Conn.CommitTx2(context.Background(), sendId, txid, commit)
		select {
		case <-ch:
		case <-time.After(timeout):
			err = Errors(sys.ERR_TIMEOUT)
		}
	} else {
		err = Errors(sys.ERR_NODE_NOFOUND)
	}
	return
}

func (this *_nodeWare) broadcastMQ(mqType int8, bs []byte) {
	if this.tlMap.Len() > 0 {
		uuidMap := make(map[int64]int8, 0)
		this.tlMap.Range(func(k *tlContext, _ int8) bool {
			if _, ok := uuidMap[k.RemoteUuid]; !ok {
				GoPool.Go(func() { this.pubMq(k.RemoteUuid, 0, mqType, bs, sys.ReadTimeout) })
				uuidMap[k.RemoteUuid] = 0
			}
			return true
		})
	}
}

func (this *_nodeWare) pubMq(uuid, sendId int64, mqType int8, bs []byte, timeout time.Duration) (err error) {
	defer errRecover()
	c := 10
	if sendId == 0 {
		sendId = newTxId()
	}
	defer await.DelAndClose(sendId)
	for c > 0 {
		c--
		if tc := this.GetTlContext(uuid); tc != nil {
			tc.Conn.PubMq(context.Background(), sendId, mqType, bs)
			ch := await.Get(sendId)
			select {
			case <-ch:
				goto END
			case <-time.After(timeout):
				err = Errors(sys.ERR_TIMEOUT)
			}
		} else {
			err = Errors(sys.ERR_NODE_NOFOUND)
			<-time.After(timeout)
		}
	}
END:
	return
}

func (this *_nodeWare) pullData(uuid, sendId int64, ldb *LogDataBean, timeout time.Duration, async bool) (_r *LogDataBean, err error) {
	defer errRecover()
	if tc := this.GetTlContext(uuid); tc != nil {
		if sendId == 0 && !async {
			sendId = newTxId()
		}
		tc.Conn.PullData(context.Background(), sendId, ldb)
		if !async {
			ch := awaitLog.Get(sendId)
			defer awaitLog.DelAndClose(sendId)
			select {
			case _r = <-ch:
			case <-time.After(timeout):
				err = Errors(sys.ERR_TIMEOUT)
			}
		}
	} else {
		err = Errors(sys.ERR_NODE_NOFOUND)
	}
	return
}

// func (this *_nodeWare) broadcast3Sync(sendId int64, ack bool, pb *PonBean, t time.Duration) (exce *sys.Exception) {
// 	var wg sync.WaitGroup
// 	exce = sys.NewException()
// 	uuidMap := make(map[int64]int8, 0)
// 	this.tlMap.Range(func(k *TlContext, v int8) bool {
// 		if _, ok := uuidMap[k.RemoteUuid]; !ok {
// 			wg.Add(1)
// 			go this.pon3(exce, k.RemoteUuid, sendId, ack, pb, t, &wg)
// 			uuidMap[k.RemoteUuid] = 0
// 		}
// 		return true
// 	})
// 	wg.Wait()
// 	return
// }

// func (this *_nodeWare) broadcast3(sendId int64, ack bool, pb *PonBean, t time.Duration, async bool) (exce *sys.Exception) {
// 	exce = sys.NewException()
// 	uuidMap := make(map[int64]int8, 0)
// 	this.tlMap.Range(func(k *TlContext, v int8) bool {
// 		if _, ok := uuidMap[k.RemoteUuid]; !ok {
// 			go this.pon3(exce, k.RemoteUuid, sendId, ack, pb, t, async, nil)
// 			uuidMap[k.RemoteUuid] = 0
// 		}
// 		return true
// 	})
// 	return
// }

func (this *_nodeWare) pon3(uuid, sendId int64, ack bool, pb *PonBean, timeout time.Duration) (_r *PonBean, err error) {
	spTime := 5 * time.Second
	for timeout > 0 {
		if tc := this.GetTlContext(uuid); tc != nil {
			if sendId == 0 {
				sendId = newTxId()
			} else if !ack {
				sendId = newTxId()
			}
			pb.SyncId = &sendId
			tc.Conn.Pon3(context.Background(), EncodePonBean(pb).Bytes(), sendId, ack)
			if !ack {
				ch := awaitPB.Get(sendId)
				defer awaitPB.DelAndClose(sendId)
				select {
				case _r = <-ch:
					goto END
				case <-time.After(spTime):
					timeout = timeout - spTime
				}
			} else {
				ch := await.Get(sendId)
				defer await.DelAndClose(sendId)
				select {
				case <-ch:
					goto END
				case <-time.After(spTime):
					timeout = timeout - spTime
				}
			}
		} else {
			timeout = timeout - time.Second
			<-time.After(timeout)
		}
	}
END:
	if timeout <= 0 {
		err = Errors(sys.ERR_TIMEOUT)
	}
	return
}

func (this *_nodeWare) broadcastReinit(sbean *SysBean) (exce *sys.Exception) {
	exce = sys.NewException()
	if this.tlMap.Len() > 0 {
		var wg sync.WaitGroup
		uuidMap := make(map[int64]int8, 0)
		this.tlMap.Range(func(k *tlContext, _ int8) bool {
			if _, ok := uuidMap[k.RemoteUuid]; !ok {
				wg.Add(1)
				GoPool.Go(func() { this.reinit(exce, k.RemoteUuid, sbean, &wg) })
				uuidMap[k.RemoteUuid] = 0
			}
			return true
		})
		wg.Wait()
	}
	return
}

// over time = timeout*30
func (this *_nodeWare) reinit(exce *sys.Exception, uuid int64, sbean *SysBean, wg *sync.WaitGroup) {
	defer errRecover()
	limit := 30
	var err error
	sendId := newTxId()
	ch := await.Get(sendId)
	defer await.DelAndClose(sendId)
	for limit > 0 {
		err = nil
		limit--
		if tc := this.GetTlContext(uuid); tc != nil {
			err = tc.Conn.ReInit(context.Background(), sendId, sbean)
		}
		select {
		case <-ch:
			goto END
		case <-time.After(sys.ReadTimeout):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
END:
	if wg != nil {
		wg.Done()
	}
	if err != nil {
		exce.Code = sys.ERR_TIMEOUT
		exce.Put(uuid)
	}
	return
}

////////////////////////////////////////////////////////////////////////

func (this *_nodeWare) broadcastRunNodeProxy(paramData []byte, uuid int64, pType int8, ctype int8) (_r *TableParam, exce *sys.Exception) {
	exce = sys.NewException()
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		match := true
		if uuid > 0 {
			match = k.RemoteUuid == uuid
		}
		if k.stat == sys.RUN && match {
			_r = this.proxyCall(exce, k.RemoteUuid, paramData, pType, ctype, sys.WaitTimeout)
			if exce.Error() == nil {
				return false
			}
		}
		return true
	})
	return
}

func (this *_nodeWare) proxyCall(exce *sys.Exception, uuid int64, paramData []byte, pType int8, ctype int8, timeout time.Duration) (_r *TableParam) {
	defer errRecover()
	if tc := this.GetTlContext(uuid); tc != nil {
		syncId := newTxId()
		if err := tc.Conn.ProxyCall(context.Background(), syncId, paramData, pType, ctype); err == nil {
			ch := awaitProxy.Get(syncId)
			defer awaitProxy.DelAndClose(syncId)
			select {
			case _r = <-ch:
			case <-time.After(timeout):
				exce.Code = sys.ERR_TIMEOUT
			}
		} else {
			exce.Code = sys.ERR_BROADCAST
		}
	} else {
		exce.Code = sys.ERR_NODE_NOFOUND
	}
	return
}

func (this *_nodeWare) broadToken(exce *sys.Exception, sendId, uuid int64, tt *TokenTrans, ack bool, wg *sync.WaitGroup) (_r *TokenTrans) {
	defer errRecover()
	if uuid == sys.UUID {
		return
	}

	limit := 30
	var err error
	for limit > 0 {
		err = nil
		limit--
		if sendId == 0 {
			sendId = newTxId()
		} else if !ack {
			sendId = newTxId()
		}
		tt.SyncId = &sendId
		if tc := this.GetTlContext(uuid); tc != nil {
			err = tc.Conn.BroadToken(context.Background(), sendId, tt, ack)
		}
		if !ack {
			ch := awaitToken.Get(sendId)
			defer awaitToken.DelAndClose(sendId)
			select {
			case _r = <-ch:
				goto END
			case <-time.After(sys.ReadTimeout):
				err = Errors(sys.ERR_TIMEOUT)
			}
		} else {
			ch := await.Get(sendId)
			defer await.DelAndClose(sendId)
			select {
			case <-ch:
				goto END
			case <-time.After(sys.ReadTimeout):
				err = Errors(sys.ERR_TIMEOUT)
			}
		}

	}
END:
	if wg != nil {
		wg.Done()
	}
	if err != nil && exce != nil {
		exce.Code = sys.ERR_TIMEOUT
		exce.Put(uuid)
	}
	return
}

/**************************************************************************/

func syncTime(tc *tlContext) {
	defer errRecover()
	if tc != nil {
		tc.Conn.Time(context.Background(), time.Now().UnixNano(), 0, 200, false)
	}
}

func syncNode(tc *tlContext) {
	defer errRecover()
	if tc != nil {
		tc.Conn.SyncNode(context.Background(), nodeWare.GetUUIDNode(), true)
	}
}

func (this *_nodeWare) print() {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		logger.Debug(k.RemoteAddr, "====>", k.RemoteUuid)
		return true
	})
}
