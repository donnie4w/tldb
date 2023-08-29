// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"runtime/debug"
	"sync/atomic"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	. "github.com/donnie4w/tldb/util"
)

var reTxId = NewBloomFilter(1 << 17)
var clientLinkCache = NewMapL[string, int8]()
var trashTx = NewLimitMap[int64, int8](1 << 20)
var rePon = NewLimitMap[int64, int8](1 << 20)
var reMq = NewLimitMap[int64, int8](1 << 20)
var await = NewAwait[int8](1 << 8)
var awaitPB = NewAwait[*PonBean](1 << 7)
var awaitLog = NewAwait[*LogDataBean](1 << 7)
var awaitProxy = NewAwait[*TableParam](1 << 7)
var commitWare = NewLinkedMap[int64, *PonBean]()

// ////////////////////////////////
var logger = Logger
var logError = LoggerError

// ///////////////////////////////
func myRecovr() {
	if err := recover(); err != nil {
		logger.Error(string(debug.Stack()))
	}
}

// ///////////////////////////////////////////
func newTxId() (txid int64) {
	return _newTxId(100)
}

func _newTxId(num int) (txid int64) {
	txid = NewTxId()
	bs := Int64ToBytes(txid)
	if !reTxId.Contains(bs) {
		reTxId.Add(bs)
	} else {
		if num == 0 {
			return _newTxId(100)
		}
		return _newTxId(num - 1)
	}
	return
}

func EncodeBatchPacket(bc *BatchPacket) *Buffer {
	encode := NewEncodeStub()
	encode.MapFieldToBytes(bc.Kvmap)
	encode.StringArrayToBytes(bc.Dels)
	return encode.Buf
}

func DecodeBatchPacket(bs []byte) (bc *BatchPacket) {
	bc = &BatchPacket{}
	decode := &DecodeStub{Bytes: bs}
	bc.Kvmap = decode.ToMapField()
	bc.Dels = decode.ToStrArray()
	return
}

func EncodePonBean(pb *PonBean) (buf *Buffer) {
	ptype := Int16ToBytes(pb.Ptype)       //2
	opType := Int16ToBytes(pb.OpType)     //2
	txid := Int64ToBytes(pb.Txid)         //8
	fromuuid := Int64ToBytes(pb.Fromuuid) //8
	isPass := boolToByte(pb.IsPass)       //1
	isFirst := boolToByte(pb.IsFirst)     //1
	isExcu := boolToByte(pb.IsExcu)       //1
	isVerify := boolToByte(pb.IsVerify)   //1
	doCommit := boolToByte(pb.DoCommit)   //1
	mode := byte(pb.Mode)                 //1
	//////
	buf = BufPool.Get()
	buf.Write(ptype)
	buf.Write(opType)
	buf.Write(txid)
	buf.Write(fromuuid)
	buf.WriteByte(isPass)
	buf.WriteByte(isFirst)
	buf.WriteByte(isExcu)
	buf.WriteByte(isVerify)
	buf.WriteByte(doCommit)
	buf.WriteByte(mode)
	if pb.ExecNode != nil {
		execNode := IntArrayToBytes(pb.ExecNode)
		buf.WriteInt32(len(execNode))
		buf.Write(execNode)
	} else {
		buf.WriteInt32(0)
	}
	if pb.ClusterNode != nil {
		clusterNode := IntArrayToBytes(pb.ClusterNode)
		buf.WriteInt32(len(clusterNode))
		buf.Write(clusterNode)
	} else {
		buf.WriteInt32(0)
	}
	if pb.Batch != nil {
		batch := EncodeBatchPacket(pb.Batch)
		defer batch.Free()
		buf.WriteInt32(len(batch.Bytes()))
		buf.Write(batch.Bytes())
	} else {
		buf.WriteInt32(0)
	}
	if pb.SyncId != nil {
		syncId := Int64ToBytes(*pb.SyncId)
		buf.Write(syncId)
	} else {
		buf.Write(Int64ToBytes(0))
	}
	if pb.Key != nil {
		buf.WriteInt32(len(pb.Key))
		buf.Write(pb.Key)
	} else {
		buf.WriteInt32(0)
	}

	if pb.Value != nil {
		buf.WriteInt32(len(pb.Value))
		buf.Write(pb.Value)
	} else {
		buf.WriteInt32(0)
	}

	if pb.Pbtime != nil {
		buf.Write(pb.Pbtime)
	}
	return
}

func DecodePonBean(bs []byte) (pb *PonBean) {
	var c *int32 = new(int32)
	pb = &PonBean{}
	pb.Ptype = BytesToInt16(bs[*c:add(c, 2)])     //2
	pb.OpType = BytesToInt16(bs[*c:add(c, 2)])    //2
	pb.Txid = BytesToInt64(bs[*c:add(c, 8)])      //8
	pb.Fromuuid = BytesToInt64(bs[*c:add(c, 8)])  //8
	pb.IsPass = byteToBool(bs[*c:add(c, 1)][0])   //1
	pb.IsFirst = byteToBool(bs[*c:add(c, 1)][0])  //1
	pb.IsExcu = byteToBool(bs[*c:add(c, 1)][0])   //1
	pb.IsVerify = byteToBool(bs[*c:add(c, 1)][0]) //1
	pb.DoCommit = byteToBool(bs[*c:add(c, 1)][0]) //1
	pb.Mode = int8(bs[*c:add(c, 1)][0])           //1
	if ln := BytesToInt32(bs[*c:add(c, 4)]); ln > 0 {
		pb.ExecNode = BytesToIntArray(bs[*c:add(c, ln)])
	}
	if ln := BytesToInt32(bs[*c:add(c, 4)]); ln > 0 {
		pb.ClusterNode = BytesToIntArray(bs[*c:add(c, ln)])
	}
	if ln := BytesToInt32(bs[*c:add(c, 4)]); ln > 0 {
		pb.Batch = DecodeBatchPacket(bs[*c:add(c, ln)])
	}
	if i := BytesToInt64(bs[*c:add(c, 8)]); i > 0 {
		pb.SyncId = &i
	}
	if ln := BytesToInt32(bs[*c:add(c, 4)]); ln > 0 {
		pb.Key = bs[*c:add(c, ln)]
	}
	if ln := BytesToInt32(bs[*c:add(c, 4)]); ln > 0 {
		pb.Value = bs[*c:add(c, ln)]
	}
	pb.Pbtime = bs[*c:]
	return
}

func add(c *int32, v int32) int32 {
	return atomic.AddInt32(c, v)
}

func boolToByte(b bool) (_r byte) {
	if b {
		_r = 1
	}
	return
}
func byteToBool(b byte) (_r bool) {
	if b == 1 {
		_r = true
	}
	return
}

func EncodeDataBean(db *DataBeen) (buf *Buffer) {
	buf = NewBufferByPool()
	isLocal := boolToByte(db.IsLocal)
	buf.WriteByte(isLocal)
	if db.Uuids != nil {
		uuids := IntArrayToBytes(db.Uuids)
		buf.WriteInt32(len(uuids))
		buf.Write(uuids)
	} else {
		buf.WriteInt32(0)
	}
	if db.Value != nil {
		buf.Write(db.Value)
	}
	return
}

func DecodeDataBeen(bs []byte) (_r *DataBeen) {
	var c *int32 = new(int32)
	_r = &DataBeen{}
	_r.IsLocal = byteToBool(bs[*c:add(c, 1)][0])
	if ln := BytesToInt32(bs[*c:add(c, 4)]); ln > 0 {
		_r.Uuids = BytesToIntArray(bs[*c:add(c, ln)])
	}
	_r.Value = bs[*c:]
	return
}
