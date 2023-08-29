// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	. "github.com/donnie4w/tldb/level0"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
	. "github.com/donnie4w/tldb/util"
)

func pas_ponIncrkey(pb *PonBean) {
	txid := pb.Txid
	key, value := string(pb.GetKey()), pb.GetValue()
	if pb.IsExcu {
		if sys.IsREADY() {
			pb.Batch = &BatchPacket{Kvmap: map[string][]byte{string(pb.GetKey()): pb.GetValue()}}
			tlog.WriteCacheLog(pb)
		} else if sys.IsRUN() {
			if v, err := Level0.Get(key); err == nil && v != nil {
				if BytesToInt64(value) <= BytesToInt64(v) {
					return
				}
			}
			if err := tlog.writeKV(key, txid, value); err != nil {
				fatalError(err)
			}
			incrExcuPool.del(key)
		}
	} else if pb.ExecNode[0] == sys.UUID {
		ackPb := newPonBean(txid, 1, pb.OpType, []int64{sys.UUID}, nodeWare.GetAllRunUUID())
		ackPb.Key = pb.GetKey()
		if !sys.IsRUN() {
			ackPb.Value = Int64ToBytes(0)
			GoPool.Go(func() { pubSingleTx(ackPb, pb.Fromuuid) })
		} else {
			incrNetKeyExcu(key, ackPb)
		}
	}
}

func pas_ponGet(pb *PonBean) {
	ackpb := newPonBean(pb.Txid, 3, 0, []int64{pb.Fromuuid}, nodeWare.GetAllRunUUID())
	key := string(pb.Key)
	if value, err := Level1.GetValue(key); err == nil {
		ackpb.Value = value
	}
	pubSingleTx(ackpb, pb.Fromuuid)
}

func pas_ponGetRemote(pb *PonBean) {
	ackpb := newPonBean(pb.Txid, 31, 0, nil, nil)
	key := string(pb.Key)
	if value, err := Level1.GetLocal(key); err == nil {
		if data := decodeDataBeen(value); data != nil && data.IsLocal {
			ackpb.Value = data.Value
		}
	}
	go pubSingleTx(ackpb, pb.Fromuuid)
}

func pas_ponBatch(pb *PonBean) (err error) {
	if !pb.IsFirst {
		return
	}
	if pb.IsFirst && pb.IsVerify && pb.DoCommit {
		t := int64(0)
		if pb.Pbtime != nil {
			t = BytesToInt64(pb.Pbtime)
		} else {
			t = Time().UnixNano()
		}
		if err = tlog.writePonBean(t, pb); err == nil {
			trashTx.Put(pb.Txid, 0)
		}
		if err != nil {
			fatalError(err)
		}
		return
	}

	if !sys.IsRUN() && util.ArrayContain(pb.ClusterNode, sys.UUID) {
		err = Errors(sys.ERR_NO_RUNSTAT)
		return
	}
	switch sys.BatchMode {
	case 1:
		pb.Ptype = 22
		taskWare.pubTaskPB(pb)
		if sys.IsRUN() {
			excuuuids := Pon().ExcuBatchNodes(pb.Batch.Kvmap, pb.Batch.Dels)
			ackPb := newPonBean(pb.Txid, 22, pb.OpType, excuuuids, nodeWare.GetAllRunUUID())
			GoPool.Go(func() { pubSingleTx(ackPb, pb.Fromuuid) })
		}
	case 2:
		taskWare.pubTaskPB(pb)
		if sys.IsRUN() {
			excuuuids := Pon().ExcuBatchNodes(pb.Batch.Kvmap, pb.Batch.Dels)
			ackPb := newPonBean(pb.Txid, pb.Ptype, pb.OpType, excuuuids, nodeWare.GetAllRunUUID())
			GoPool.Go(func() { pubTxAsync(ackPb) })
			taskWare.pubTaskPB(ackPb)
		}
	}
	return
}

func pas_ponSyncKey(pb *PonBean) {
	ackpb := newPonBean(pb.Txid, 4, 0, nil, nil)
	key := string(pb.Key)
	// netLock.Lock(key)
	// defer netLock.Unlock(key)
	ackpb.Key = pb.Key
	ackpb.Value = Int64ToBytes(0)
	if v, ok := incrExcuPool.getWithoutLock(key); ok {
		ackpb.Value = util.Int64ToBytes(*v)
		if Pon().ExcuNode(key) != sys.UUID {
			incrExcuPool.del(key)
		}
	} else if value, err := Level0.Get(key); err == nil {
		ackpb.Value = value
	}
	pubSingleTx(ackpb, pb.Fromuuid)
}

/*********************************************************************/

func pas_ponStat(pb *PonBean) {
	taskWare.pubTaskPB(pb)
	ackpb := newPonBean(pb.Txid, 5, 0, nil, nil)
	// taskWare.pubTaskPB(ackpb)
	pub2SingleTx(ackpb, pb.Fromuuid)
}

func pas_ponLoad(pb *PonBean) {
	lb := pb.LoadBean
	ackpb := newPonBean(pb.Txid, 6, 0, pb.ExecNode, nodeWare.GetAllRunUUID())
	ackpb.LoadBean = lb
	pub2SingleTx(ackpb, pb.Fromuuid)
}

func pas_ponTime(pb *PonBean) {
	ackpb := newPonBean(pb.Txid, 7, 0, []int64{pb.Fromuuid}, nodeWare.GetAllRunUUID())
	ackpb.TimeBean = NewTimeBean()
	pon := PonUUIDS(nodeWare.GetALLUUID())
	excuuuid := pon.ExcuNode("time")
	if sys.IsREADY() {
		if excuuuid == pb.ExecNode[0] {
			ackpb.TimeBean.IsPass, ackpb.TimeBean.RunUUID = true, 0
		} else {
			ackpb.TimeBean.IsPass, ackpb.TimeBean.RunUUID = false, 0
		}
	} else {
		if sys.IsRUN() || sys.TIME_DEVIATION != 0 {
			ackpb.TimeBean.IsPass, ackpb.TimeBean.RunUUID = false, sys.UUID
		} else if uuids := nodeWare.GetAllRunUUID(); len(uuids) > 0 {
			ackpb.TimeBean.IsPass, ackpb.TimeBean.RunUUID = false, uuids[0]
		} else {
			ackpb.TimeBean.IsPass, ackpb.TimeBean.RunUUID = true, 0
		}
	}
	pub2SingleTx(ackpb, pb.Fromuuid)
}

func pas_ponSeq(pb *PonBean) {
	ackpb := newPonBean(pb.Txid, 8, 0, []int64{pb.Fromuuid}, nil)
	pub2SingleTx(ackpb, pb.Fromuuid)
}
