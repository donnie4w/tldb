// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level1

import (
	"sync"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

var statAdmin = NewStatAdmin()

func NewStatAdmin() (_r *_statAdmin) {
	_r = &_statAdmin{mux: &sync.Mutex{}, nodeStatMap: NewMapL[int64, int8](), injectMap: NewMapL[int64, func()]()}
	return
}

type _statAdmin struct {
	mux         *sync.Mutex
	nodeStatMap *MapL[int64, int8]
	injectMap   *MapL[int64, func()]
}

func (this *_statAdmin) put(uuid int64) {
	this.nodeStatMap.Put(uuid, 0)
}

func (this *_statAdmin) addInject(id int64, f func()) {
	this.injectMap.Put(id, f)
}

func (this *_statAdmin) syncRunInject() {
	defer errRecover()
	defer this.mux.Unlock()
	this.mux.Lock()
	this.injectMap.Range(func(_ int64, f func()) bool {
		go func() {
			defer errRecover()
			f()
		}()
		return true
	})
}

func (this *_statAdmin) clear() {
	this.nodeStatMap.Range(func(k int64, _ int8) bool {
		this.nodeStatMap.Del(k)
		return true
	})
}

func (this *_statAdmin) isSyncOver() bool {
	t1, _ := tlog.getCacheStat()
	t2, _ := tlog.getStat()
	sy := t1 <= t2 && this.nodeStatMap.Len() >= int64(sys.CLUSTER_NUM)-1 && this.nodeStatMap.Len() == int64(len(nodeWare.GetRemoteUUIDS()))
	runids := nodeWare.GetRemoteRunUUID()
	getall := true
	for _, runid := range runids {
		if !this.nodeStatMap.Has(runid) {
			getall = false
			break
		}
	}
	return sy && getall
}

func (this *_statAdmin) pullData() (txid int64) {
	defer this.mux.Unlock()
	this.mux.Lock()
	defer statLock.Unlock()
	statLock.Lock()
	if this.nodeStatMap.Len() != int64(len(nodeWare.GetRemoteUUIDS())) && len(nodeWare.GetALLUUID()) >= sys.CLUSTER_NUM {
		m2 := nodeWare.GetRemoteUUIDS()
		for _, m := range m2 {
			if nodeWare.GetTlContext(m).stat != sys.PROXY {
				if _txid := tlog.pullData(m); txid < _txid {
					txid = _txid
				}
			}
		}
	}
	return
}

func (this *_statAdmin) amendPull() {
	if !this.isSyncOver() {
		t1, _ := tlog.getCacheStat()
		t2, _ := tlog.getStat()
		if t1 > t2 {
			if this.nodeStatMap.Len() >= int64(sys.CLUSTER_NUM)-1 && this.nodeStatMap.Len() == int64(len(nodeWare.GetRemoteUUIDS())) {
				runids := nodeWare.GetRemoteRunUUID()
				getall := true
				for _, runid := range runids {
					if !this.nodeStatMap.Has(runid) {
						getall = false
						break
					}
				}
				if getall {
					this.nodeStatMap.Range(func(k int64, _ int8) bool {
						this.nodeStatMap.Del(k)
						return true
					})
				}
			}
		}
	}
}

/***********************************************************************************/
var statLock = &sync.Mutex{}

func checkAndResetStat() {
	defer statLock.Unlock()
	statLock.Lock()
	if (len(nodeWare.GetALLUUID()) < sys.CLUSTER_NUM) && sys.IsRUN() {
		sys.SetStat(sys.READY, 0)
		statAdmin.clear()
		go pos_stat(sys.READY, 0)
	}
}

func setStat(stat sys.STATTYPE, timenano time.Duration) {
	defer statLock.Unlock()
	statLock.Lock()
	sys.SetStat(stat, timenano)
}

func fatalError(err error) {
	defer statLock.Unlock()
	statLock.Lock()
	sys.FmtLog("Fatal Error: set stat to PROXY,[", sys.UUID, "] Error:", err)
	sys.SetStat(sys.PROXY, 0)
	go pos_stat(sys.PROXY, 0)
}

func tryRunToProxy() (err error) {
	if statLock.TryLock() {
		defer statLock.Unlock()
	} else {
		err = util.Errors(sys.ERR_SETSTAT)
		return
	}
	if sys.IsRUN() {
		statAdmin.clear()
		sys.SetStat(sys.PROXY, 0)
		go pos_stat(sys.PROXY, 0)
	} else {
		err = util.Errors(sys.ERR_SETSTAT)
	}
	return
}

func tryProxyToReady() (err error) {
	if statLock.TryLock() {
		defer statLock.Unlock()
	} else {
		err = util.Errors(sys.ERR_SETSTAT)
		return
	}
	if sys.IsPROXY() {
		statAdmin.clear()
		sys.SetStat(sys.READY, 0)
		go pos_stat(sys.READY, 0)
	} else {
		err = util.Errors(sys.ERR_SETSTAT)
	}
	return
}

func checkAndSetRunToReady() (err error) {
	defer statLock.Unlock()
	statLock.Lock()
	if len(nodeWare.GetAllRunUUID()) < sys.CLUSTER_NUM {
		if sys.IsRUN() {
			statAdmin.clear()
			sys.SetStat(sys.READY, 0)
			go pos_stat(sys.READY, 0)
		} else {
			err = util.Errors(sys.ERR_SETSTAT)
		}
	}
	return
}

/***********************************************************************************/
