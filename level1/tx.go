// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package level1

// import (
// 	. "github.com/donnie4w/tldb/util"
// )

// var TxWare = NewTxWare(sys.TransTimeout)

// func NewTxWare(t time.Duration) (_r *_txWare) {
// 	_r = &_txWare{NewMap[int64, *txBeen](), NewTlLock[int64]()}
// 	go _r.timer(t)
// 	return
// }

// type _txWare struct {
// 	waremap *Map[int64, *txBeen]
// 	mux     TlLock[int64]
// }

// func (this *_txWare) NewTx(txid int64) (tb *txBeen) {
// 	this.mux.Lock(txid)
// 	defer this.mux.UnLock(txid)
// 	logger.Info("NewTx txid:", txid)
// 	if TrashTx.Has(txid) {
// 		return
// 	}
// 	var ok bool
// 	if tb, ok = this.waremap.Get(txid); !ok {
// 		tb = newTxBeen(txid)
// 		this.waremap.Put(txid, tb)
// 	}
// 	return
// }

// func (this *_txWare) AddPonPB(pb *PonBean) {
// 	this.mux.Lock(pb.Txid)
// 	defer this.mux.UnLock(pb.Txid)
// 	logger.Info("AddPonPB txid:", pb.Txid)
// 	if !pb.IsFirst && !TrashTx.Has(pb.Txid) {
// 		if tb, ok := this.waremap.Get(pb.Txid); ok {
// 			tb.AddTxPB(pb)
// 		} else {

// 		}
// 	}
// }

// func (this *_txWare) GetTx(txid int64) (tb *txBeen, ok bool) {
// 	tb, ok = this.waremap.Get(txid)
// 	return
// }

// func (this *_txWare) HasTx(txid int64) (ok bool) {
// 	return this.waremap.Has(txid)
// }

// func (this *_txWare) DelTx(txid int64) {
// 	this.mux.Lock(txid)
// 	defer this.mux.UnLock(txid)
// 	TrashTx.Add(txid, 0)
// 	this.waremap.Del(txid)
// }

// func (this *_txWare) AddPosPB(pb *PonBean) *txBeen {
// 	// defer this.mux.UnLock(pb.Txid)
// 	// this.mux.Lock(pb.Txid)
// 	if tb, ok := this.waremap.Get(pb.Txid); ok {
// 		logger.Info("AddPosPB: txid:", pb.Txid, ", Fromuuid:", pb.Fromuuid)
// 		tb.AddTxPB(pb)
// 		tb.tryDone()
// 		return tb
// 	} else if TrashTx.Has(pb.Txid) {
// 		logger.Warn("tx has finish,txid:", pb.Txid)
// 	} else {
// 		logger.Error("tx not exist,txid:", pb.Txid)
// 	}
// 	return nil
// }

// func (this *_txWare) timer(t time.Duration) {
// 	ticker := time.NewTicker(t)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			this.__timer(t)
// 		}
// 	}
// }

// func (this *_txWare) __timer(t time.Duration) {
// 	defer myRecovr()
// 	<-time.After(t)
// 	warearr := make([]*txBeen, 0)
// 	this.waremap.Range(func(k int64, v *txBeen) bool {
// 		warearr = append(warearr, v)
// 		return true
// 	})
// 	<-time.After(t)
// 	for _, tb := range warearr {
// 		if ok := this.waremap.Has(tb.Txid); ok {
// 			logger.Error("tx cancel tb:", tb.Txid)
// 			tb.failedAndcancel()
// 			this.DelTx(tb.Txid)
// 			// this.waremap.Del(tb.Txid)
// 		}
// 	}
// }

// // //////////////////////////////////////////////////
// type txBeen struct {
// 	mux     *sync.Mutex
// 	Txid    int64
// 	txPbMap *Map[int64, *PonBean] //uuid-> pb
// 	txSync  *txBeenSync
// }

// func newTxBeen(_txid int64) *txBeen {
// 	return &txBeen{mux: &sync.Mutex{}, Txid: _txid, txPbMap: NewMap[int64, *PonBean]()}
// }

// func (this *txBeen) AddTxPB(pb *PonBean) {
// 	this.txPbMap.Put(pb.Fromuuid, pb)
// }

// func (this *txBeen) GetPonBean(uuid int64) (pb *PonBean, ok bool) {
// 	pb, ok = this.txPbMap.Get(uuid)
// 	return
// }

// func (this *txBeen) doneFunc(doneFunc func(pbMap *Map[int64, *PonBean]) bool) {
// 	this.txSync = newtxBeenSync(doneFunc)
// }

// func (this *txBeen) wait() bool {
// 	return this.txSync.wait()
// }

// func (this *txBeen) tryDone() {
// 	if this.txSync != nil {
// 		this.txSync.tryDone(this.txPbMap)
// 	}
// }

// func (this *txBeen) len() int {
// 	return int(this.txPbMap.Len())
// }

// func (this *txBeen) failedAndcancel() {
// 	if this.txSync != nil {
// 		this.txSync.failedAndcancel()
// 	}
// }

// // ///////////////////////////////////////////////////////
// type txBeenSync struct {
// 	ch        chan bool
// 	isSetChan bool
// 	doneFunc  func(pbMap *Map[int64, *PonBean]) bool //uuid-> pb
// 	mux       *sync.Mutex
// }

// func newtxBeenSync(doneFunc func(pbMap *Map[int64, *PonBean]) bool) *txBeenSync {
// 	return &txBeenSync{make(chan bool, 1), false, doneFunc, &sync.Mutex{}}
// }

// func (this *txBeenSync) wait() bool {
// 	return <-this.ch
// }

// func (this *txBeenSync) tryDone(pbMap *Map[int64, *PonBean]) {
// 	defer myRecovr()
// 	this.mux.Lock()
// 	defer this.mux.Unlock()
// 	if this.doneFunc != nil && pbMap != nil && this.doneFunc(pbMap) && !this.isSetChan {
// 		this.isSetChan = true
// 		this.ch <- true
// 	}
// }

// func (this *txBeenSync) failedAndcancel() {
// 	defer myRecovr()
// 	if !this.isSetChan {
// 		this.isSetChan = true
// 		this.ch <- false
// 	}
// }
