// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

package level1

import (
	"github.com/donnie4w/tldb/sys"
)

func pos_stat(_stat sys.STATTYPE, timenano int64) (err error) {
	defer errRecover()
	txid := newTxId()
	pb := newPonBean(txid, 5, 0, nil, nil)
	pb.Stat, pb.IsFirst = &Stat{Stat: int8(_stat), Timenano: timenano}, true
	task, _ := taskWare.pubTaskPB(pb)
	if err = pub2TxSync(pb).Error(); err == nil {
		pb2 := newPonBean(txid, 5, 0, nil, nil)
		pb2.IsVerify = true
		if err = pub2TxSync(pb2).Error(); err == nil {
			task.addPB(pb2)
			err = task.waitWithTimeout(sys.ReadTimeout)
			// go setStat(_stat, time.Duration(timenano-Time().UnixNano()))
		}
	}
	return
}

// func pos_load(lb *LoadBean, uuid int64) (m map[int64]*LogKV, err error) {
// 	defer myRecovr()
// 	txid := newTxId()
// 	logger.Info("pos_load==================>", txid)
// 	pb := newPonBean(txid, 6, 0, []int64{uuid}, nil)
// 	pb.LoadBean, pb.IsFirst = lb, true
// 	task := taskWare.pubTaskPB(pb)
// 	if err = Pub2SingleTx(pb, uuid); err == nil {
// 		if err = task.waitWithTimeout(sys.WaitTimeout); err == nil {
// 			m = task.getPB(uuid).LoadBean.Logm
// 		}
// 	}
// 	return
// }

func pos_time() (err error) {
	txid := newTxId()
	pon := PonUUIDS(nodeWare.GetRemoteUUIDS())
	excuuuid := pon.ExcuNode("time")
	pb := newPonBean(txid, 7, 0, []int64{excuuuid}, nil)
	pb.IsFirst = true
	task, _ := taskWare.pubTaskPB(pb)
	if err = pub2TxSync(pb).Error(); err == nil {
		if err = task.waitWithTimeout(sys.ReadTimeout); err == nil {
			m := task.getAllPB()
			var timeBeans = make([]*TimeBean, 0)
			m.Range(func(k int64, v *PonBean) bool {
				if k != sys.UUID {
					timeBeans = append(timeBeans, v.TimeBean)
				}
				return true
			})
			ispass, _excuuuid := 0, int64(0)
			for _, v := range timeBeans {
				if v.IsPass {
					ispass++
				} else if v.RunUUID > 0 {
					_excuuuid = v.RunUUID
				}
			}
			length := len(pon.uuids)
			if ispass >= length/2 {
				if !isLocal([]int64{excuuuid}) {
					syncTime(nodeWare.GetTlContext(excuuuid))
				}
			} else if _excuuuid > 0 {
				syncTime(nodeWare.GetTlContext(_excuuuid))
			} else {
				logger.Warn("pos_time warn,excuuuid:", excuuuid)
			}
		}
	}
	return
}

// func pos_seq() (err error) {
// 	defer myRecovr()
// 	txid := newTxId()
// 	logger.Info("pos_seq=================>", txid)
// 	pb := newPonBean(txid, 8, 0, nil, nil)
// 	pb.IsFirst = true
// 	task := taskWare.pubTaskPB(pb)
// 	if err = Pub2TxSync(pb).Error(); err == nil {
// 		if err = task.waitWithTimeout(sys.ReadTimeout); err == nil {
// 			task.getAllPB().Range(func(k int64, v *PonBean) bool {
// 				if k != sys.UUID {
// 					syncDataObj.setNodeKM(k, v.SeqBean.Seq)
// 				}
// 				return true
// 			})
// 		}
// 	}
// 	return
// }

// /////////////////////////////////////////////
