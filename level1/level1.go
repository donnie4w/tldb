// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level1

import (
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/level0"
	. "github.com/donnie4w/tldb/lock"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"
)

var netLock = NewStrlock(1 << 13)

var Level1 = &level1{}

type level1 struct{}

func (this *level1) GetLocalLike(key string) (datamap map[string][]byte, err error) {
	defer errRecover()
	datamap, err = Level0.GetByPrefix(key)
	return
}

func (this *level1) GetLocalKeysLike(key string) (keys []string, err error) {
	defer errRecover()
	keys, err = Level0.GetKeysPrefix(key)
	return
}

func (this *level1) GetLocalKeysLikeLimit(key string, limit int) (keys []string, err error) {
	defer errRecover()
	keys, err = Level0.GetKeysPrefixLimit(key, limit)
	return
}

func (this *level1) GetLocal(key string) (value []byte, err error) {
	defer errRecover()
	value, err = Level0.Get(key)
	return
}

func (this *level1) GetValue(key string) (value []byte, err error) {
	defer errRecover()
	var bs []byte
	if bs, err = Level0.Get(key); err == nil {
		data := decodeDataBeen(bs)
		if data.IsLocal {
			value = data.Value
		} else {
			value, err = pos_Get(key, data.Uuids)
		}
	}
	return
}

func (this *level1) GetValues(key []string) (_r map[string][]byte, err error) {
	defer errRecover()
	return
}

func (this *level1) Has(key string) (_r bool) {
	return Level0.Has(key)
}

func (this *level1) Batch(optype int16, kvmap map[string][]byte, dels []string) (err error) {
	if !sys.IsRUN() || !nodeWare.IsClusRun() {
		return Errors(sys.ERR_NO_RUNSTAT)
	}
	defer errRecover()
	pon := Pon()
	excuuuids := pon.ExcuBatchNodes(kvmap, dels)
	txid := newTxId()
	pb := newPonBean(txid, 2, optype, excuuuids, pon.uuids)
	pb.Batch, pb.IsFirst = &BatchPacket{kvmap, dels}, true
	if isLocal(excuuuids) {
		pb.IsExcu = true
	}
	task, _ := taskWare.pubTaskPB(pb)
	if err = pubTxSyncWithTimeOut(pb, sys.ReadTimeout/2).Error(); err != nil {
		logError.Error("Batch err:", err)
		pb2 := newPonBean(txid, 2, optype, nil, nil)
		pb2.IsVerify, pb2.DoCommit = true, false
		pubTxSync(pb2)
		return
	}
	if nodeWare.tlMap.Len() > 0 {
		pb2 := newPonBean(txid, 2, optype, nil, nil)
		pb2.IsVerify = true
		task.addPB(pb2)
		if !task.isDone() {
			select {
			case ch := <-task.ch():
				if ch == 2 {
					logError.Error("no match Batch:txid", txid)
					err = Errors(sys.ERR_CLUS_NOMATCH)
				}
			case <-time.After(sys.ReadTimeout):
				err = Errors(sys.ERR_TIMEOUT)
				logError.Error("batch over time:txid", txid)
			}
		}
		if task.onError() == nil && err == nil {
			if exce := pubTxSyncAlltime(pb2); exce.Code != 0 {
				pb.IsVerify, pb.DoCommit, pb.Pbtime = true, true, Int64ToBytes(Time().UnixNano())
				tlog.WriteBackLog(pb, MapToArray(exce.Map()))
			}
		} else {
			pb2.DoCommit = false
			pubTxSync(pb2)
		}
	}
	return
}

func (this *level1) GetAndIncrKey(key string) (_r int64, _err error) {
	if !nodeWare.IsClusRun() || !sys.IsRUN() {
		_err = Errors(sys.ERR_NO_CLUSTER)
		return
	}
	defer errRecover()
	if v, _ := pos_incrNetKey(key); v != nil {
		_r = *v
	}
	if _r <= 0 {
		_err = Errors(sys.ERR_INCR_SEQ)
	}
	return
}

func pos_incrNetKey(key string) (_r *int64, err error) {
	defer errRecover()
	if !nodeWare.IsClusRun() || !sys.IsRUN() {
		err = Errors(sys.ERR_NO_CLUSTER)
		return
	}
	/////////////////////////////////////////////////////////////
	excuUUid := Pon().ExcuNode(key)
	value := int64(0)
	if isLocal([]int64{excuUUid}) {
		if v, er := incrNetKeyExcu(key, nil); er == nil && v > 0 {
			_r = &v
		} else {
			return
		}
	} else {
		txid := newTxId()
		pb := newPonBean(txid, 1, 1, []int64{excuUUid}, nodeWare.GetAllRunUUID())
		pb.Key, pb.IsFirst = []byte(key), true
		// go pubSingleTx(pb, excuUUid)
		// task, _ := taskWare.pubTaskPB(pb)
		// if err = task.waitWithTimeout(sys.ReadTimeout); err == nil {
		// 	v := task.getPB(excuUUid).Value
		// 	value = BytesToInt64(v)
		// 	_r = &value
		// 	if *_r > 0 {
		// 		if tlog.writeKV(txid, key, v) == nil {
		// 			if err := Level0.Put(key, v); err != nil {
		// 				fatalError(err)
		// 			}
		// 		}
		// 	} else {
		// 		logger.Error("pos_incrNetKey value=0:", key)
		// 		err = util.Errors(sys.ERR_INCR_SEQ)
		// 	}
		// } else {
		// 	logger.Error("pos_incrNetKey over time:", key)
		// }
		if _r == nil || *_r == 0 {
			var rst *PonBean
			if rst, err = requestSingleTx(pb, excuUUid); err == nil && rst != nil && rst.Value != nil {
				if value = BytesToInt64(rst.Value); value > 0 {
					_r = &value
				} else {
					err = Errors(sys.ERR_INCR_SEQ)
				}
			}
		}
	}
	return
}

func incrNetKeyExcu(key string, pb *PonBean) (value int64, err error) {
	if pb == nil {
		pb = newPonBean(newTxId(), 1, 1, []int64{sys.UUID}, nodeWare.GetAllRunUUID())
	}
	pb.Key, pb.IsExcu = []byte(key), true

	//the execution uuid of key, add the key to the lru
	if v, ok := incrExcuPool.get(key); ok {
		value = atomic.AddInt64(v, 1)
	} else {
		if value, err = func() (_cv int64, err error) {
			defer netLock.Unlock(key)
			netLock.Lock(key)
			if v, ok = incrExcuPool.get(key); ok {
				_cv = atomic.AddInt64(v, 1)
			} else {
				if _cv, err = pos_syncKey(key); err == nil {
					if _lv, er := Level0.Get(key); er == nil {
						if v := BytesToInt64(_lv); v > _cv {
							_cv = v
						}
					}
					_cv++
					incrExcuPool.add(key, &_cv)
				} else {
					checkAndResetStat()
					logger.Error("incrNetKeyExcu pos_syncKey:", err)
					return
				}
			}
			return
		}(); err != nil {
			return
		}
	}
	if value > 0 {
		pb.Value = Int64ToBytes(value)
		if err = tlog.writeKV(key, pb.Txid, pb.Value); err != nil {
			fatalError(err)
		} else {
			GoPool.Go(func() { pubTxAsync(pb) })
		}
	} else {
		err = Errors(sys.ERR_INCR_SEQ)
	}
	return
}

func pos_Get(key string, uuids []int64) (value []byte, err error) {
	defer errRecover()
	if key == "" || uuids == nil || len(uuids) == 0 {
		return nil, Errors(sys.ERR_NO_MATCH_PARAM)
	}
	// uuid := getRandUUID(uuids, nodeWare.GetALLUUID())
	// if uuid == 0 {
	// 	return nil, Errors(sys.ERR_NODE_NOFOUND)
	// }
	// txid := newTxId()
	// pb := newPonBean(txid, 3, 0, []int64{uuid}, nil)
	// pb.Key, pb.IsFirst = []byte(key), true
	// task, _ := taskWare.pubTaskPB(pb)
	// if err = PubSingleTx(pb, uuid); err == nil {
	// 	if err = task.waitWithTimeout(sys.ReadTimeout); err == nil {
	// 		p := task.getPB(uuid)
	// 		value = p.Value
	// 	}
	// }
	if uuids != nil {
		txid := newTxId()
		pb := newPonBean(txid, 3, 0, nil, nil)
		pb.Key = []byte(key)
		var ackPb *PonBean
		for _, uuid := range uuids {
			if uuid != sys.UUID {
				if ackPb, err = requestSingleTx(pb, uuid); err == nil {
					value = ackPb.Value
					break
				}
			}
		}
	} else {
		err = Errors(sys.ERR_DATA_NOEXIST)
	}
	return
}

func pos_GetRemote(key string) (value []byte, err error) {
	defer errRecover()
	if key == "" {
		return nil, Errors(sys.ERR_NO_MATCH_PARAM)
	}
	txid := newTxId()
	pb := newPonBean(txid, 31, 0, nodeWare.GetAllRunUUID(), nil)
	pb.Key, pb.IsFirst = []byte(key), true
	task, _ := taskWare.pubTaskPB(pb)
	go pubTxAsync(pb)
	if err = task.waitWithTimeout(sys.ReadTimeout); err == nil {
		if value = task.getFirstPB().GetValue(); value == nil {
			err = Errors(sys.ERR_GETDATA)
		}
	}
	return
}

func pos_syncKey(key string) (value int64, err error) {
	defer errRecover()
	if nodeWare.tlMap.Len() == 0 {
		return
	}
	txid := newTxId()
	pb := newPonBean(txid, 4, 0, nil, nodeWare.GetAllRunUUID())
	pb.Key, pb.IsFirst = []byte(key), true
	task, _ := taskWare.pubTaskPB(pb)
	GoPool.Go(func() { pubTxByRunAsync(pb, sys.ReadTimeout) })
	if err = task.waitWithTimeout(sys.ReadTimeout); err == nil {
		m := task.getAllPB()
		m.Range(func(uuid int64, v *PonBean) bool {
			if uuid != sys.UUID {
				if i := BytesToInt64(v.Value); value < i {
					value = i
				}
			}
			return true
		})
	} else {
		logger.Error("pos_syncKey over time:", key)
	}
	return
}

// //////////////////////////////////////////////////////////////////
func isLocal(uuids []int64) bool {
	for _, id := range uuids {
		if sys.UUID == id {
			return true
		}
	}
	return false
}

// //////////////////////////////////////////////////////////////
func txCallbackFunc(pbMap *Map[int64, *PonBean], excuUUids []int64) (b bool) {
	b = true
	for _, _uuid := range excuUUids {
		if !pbMap.Has(_uuid) {
			b = false
			break
		}
	}
	return
}
