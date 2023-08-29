// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"fmt"
	"strings"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

func newDataBeen(isLocal bool, uuids []int64, value []byte) (db *DataBeen) {
	db = NewDataBeen()
	db.IsLocal, db.Uuids, db.Value = isLocal, uuids, value
	return
}

func decodeDataBeen(bs []byte) (_data *DataBeen) {
	switch _DTYPE(bs[0]) {
	case T:
		// _data, _ = util.TDecode(bs[1:], NewDataBeen())
		_data = DecodeDataBeen(bs[1:])
	case Z:
		if _r, err := util.ZlibUnCz(bs[1:]); err == nil {
			// _data, _ = util.TDecode(_r, NewDataBeen())
			_data = DecodeDataBeen(_r)
		}
	}
	return
}

func PrintBeen(bs []byte) string {
	d := decodeDataBeen(bs)
	return fmt.Sprint("[isLocal:]", d.IsLocal, "[uuids:]", d.Uuids, "[value:]", string(d.Value))
}

func transfer(kvmap map[string][]byte, isLocal bool, excuUUids []int64) map[string][]byte {
	if isLocal {
		for k, v := range kvmap {
			if strings.HasPrefix(k, KEY2_ID_) {
				var t = T
				// bs := util.TEncode(newDataBeen(isLocal, nil, v))
				edb := EncodeDataBean(newDataBeen(isLocal, nil, v))
				defer edb.Free()
				bs := edb.Bytes()
				if sys.DATAZLIB {
					if _r, err := util.ZlibCz(bs); err == nil {
						bs = _r
						t = Z
					}
				}
				kvmap[k] = append([]byte{byte(t)}, bs...)
			} else {
				kvmap[k] = v
			}
		}
	} else {
		for k, v := range kvmap {
			if strings.HasPrefix(k, KEY2_ID_) {
				var t = T
				// bs := util.TEncode(newDataBeen(isLocal, excuUUids, nil))
				edb := EncodeDataBean(newDataBeen(isLocal, nil, v))
				defer edb.Free()
				bs := edb.Bytes()
				if _r, err := util.ZlibCz(bs); err == nil {
					bs = _r
					t = Z
				}
				kvmap[k] = append([]byte{byte(t)}, bs...)
			} else {
				kvmap[k] = v
			}
		}
	}
	return kvmap
}

// ///////////////////////////////////////////////////////////////

var incrExcuPool = newIncrPool()

type _incrPool struct {
	pool *LruCache[*int64]
}

func newIncrPool() *_incrPool {
	return &_incrPool{pool: (NewLruCache[*int64](1 << 20))}
}

func (this *_incrPool) add(key string, value *int64) {
	this.pool.Add(key, value)
}

func (this *_incrPool) del(key string) {
	this.pool.Remove(key)
}

func (this *_incrPool) delMulti(keys []string) {
	this.pool.RemoveMulti(keys)
}

func (this *_incrPool) get(key string) (v *int64, b bool) {
	// defer netLock.RUnlock(key)
	// netLock.RLock(key)
	return this.pool.Get(key)
}

func (this *_incrPool) getWithoutLock(key string) (v *int64, b bool) {
	return this.pool.Get(key)
}

// ///////////////////////////////////////////////////////////////
func newPonBean(txid int64, ptype int16, optype int16, uuids, clusUUids []int64) (pb *PonBean) {
	pb = NewPonBean()
	pb.Txid = txid
	pb.ExecNode = uuids
	pb.ClusterNode = clusUUids
	pb.OpType = optype
	pb.Ptype = ptype
	pb.Fromuuid = sys.UUID
	pb.IsFirst = false
	pb.IsPass = true
	pb.IsExcu = false
	pb.IsVerify = false
	pb.DoCommit = true
	pb.Mode = 0
	return
}

type _DTYPE byte

var T _DTYPE = 0
var Z _DTYPE = 1

// ////////////////////////////////////////////////////////////////////////////////

type batchLog struct {
	batch    *BatchPacket
	txid     int64
	timenano int64
}

// STEP+batchLog
func (this *batchLog) toBuffer() (write *Buffer) {
	bpbuf := EncodeBatchPacket(this.batch)
	defer bpbuf.Free()
	batchBys := util.SnappyEncode(bpbuf.Bytes())
	write = NewBufferByPool()
	write.Write(log.STEP)
	write.Write(util.Int32ToBytes(4 + 8 + 8 + int32(len(batchBys)))) //4
	write.Write(util.Int64ToBytes(this.timenano))                    //8
	write.Write(util.Int64ToBytes(this.txid))                        //8
	write.Write(batchBys)
	return
}

func parseToBatchLog(bs []byte) (bl *batchLog) {
	defer myRecovr()
	timenano := util.BytesToInt64(bs[4:12])
	txid := util.BytesToInt64(bs[12:20])
	batch := DecodeBatchPacket(util.SnappyDecode(bs[20:]))
	bl = &batchLog{batch, txid, timenano}
	return
}
