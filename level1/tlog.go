// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package level1

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/level0"
	"github.com/donnie4w/tldb/log"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

var tlog = newTlog()

type logAdmin struct {
	mux       *sync.Mutex
	_txid     int64
	syncCount int64
}

func newTlog() (_r *logAdmin) {
	_r = &logAdmin{mux: &sync.Mutex{}}
	return
}

func (this *logAdmin) getStat() (lastTime, txid int64) {
	lastTime, txid = log.BinLog.GetLastTime()
	return
}

func (this *logAdmin) getCacheStat() (_time, _txid int64) {
	_time, _txid = log.CacheLog.GetFristStat()
	return
}

func (this *logAdmin) sendBinFile(syncId, lastTime, txid, uuid int64) (err error) {
	logger.Info("send BinFile >>", syncId, " ,", lastTime, " ,", txid, " ,", uuid)
	if statKey, fileNum, er := log.LogStat.GetNum(lastTime + 1); er == nil {
		if fileNum > 0 {
			logger.Info("fileNum >>", fileNum)
			var bs []byte
			var fname string
			if bs, fname, err = log.BinLog.ReadGzipLogFile(fileNum); err == nil && bs != nil {
				fb := NewFileBean()
				fb.Size, fb.Filename, fb.Body, fb.StatKey = int64(len(bs)), fname, bs, statKey
				_, err = nodeWare.pullData(uuid, syncId, _newLogDataBean(2, 0, 0, fb, true, true), 0, true)
			} else {
				_, err = nodeWare.pullData(uuid, syncId, _newLogDataBean(2, 0, 0, nil, false, true), 0, true)
			}
			return
		} else {
			logger.Info("fileNum >>0")
			if localt, _txid := log.BinLog.GetLastTime(); _txid > 0 {
				logger.Info("local  Txid>> ", _txid, " | reqtxid >> ", txid)
				logger.Info("lo lasttime>> ", localt, " | reqlastTime >> ", lastTime)
				var bs []byte
				if _txid != txid && localt > lastTime {
					ldb := _newLogDataBean(2, 0, 0, nil, false, true)
					if bs, err = log.BinLog.ReadCurrentLog2GzipByte(); err == nil {
						ldb.FBean = NewFileBean()
						ldb.FBean.Size, ldb.FBean.Filename, ldb.FBean.Body, ldb.FBean.StatKey = int64(len(bs)), "", bs, 0
					} else {
						logError.Error("send BinFile error >> ", err)
					}
					_, err = nodeWare.pullData(uuid, syncId, ldb, 0, true)
					return
				}
			}
		}
		_, err = nodeWare.pullData(uuid, syncId, _newLogDataBean(2, 0, 0, nil, false, false), 0, true)
	} else {
		_, err = nodeWare.pullData(uuid, syncId, _newLogDataBean(2, 0, 0, nil, false, true), 0, true)
	}
	return
}

func (this *logAdmin) parseGzfile(size int64, fname string, body []byte, datetime string) (txid int64, err error) {
	lastTime, _ := this.getStat()
	endtime := int64(0)
	if t, err := util.TimeStrToTime(datetime); err == nil {
		endtime = t.Unix()
	}
	// err = util.UnGzipByBytes(body, func(bs []byte) bool {
	// 	bss := bytes.Split(append(buf.Bytes(), bs...), log.STEP)
	// 	buf.Reset()
	// 	for i, bs := range bss {
	// 		if i < len(bss)-1 {
	// 			if len(bs) > 12 {
	// 				if lastTime/1e9 > util.BytesToInt64(bs[4:12])/1e9+3 {
	// 					continue
	// 				}
	// 				if endtime > 0 && endtime < util.BytesToInt64(bs[4:12])/1e9 {
	// 					return false
	// 				}
	// 				if err = this.writeToDBbyBytes(bs); err != nil {
	// 					return false
	// 				} else {
	// 					atomic.AddInt64(&this.syncCount, 1)
	// 					txid = util.BytesToInt64(bs[:8])
	// 				}
	// 			}
	// 		} else if len(bs) > 12 && util.BytesToInt32(bs[:4]) == int32(len(bs)) {
	// 			if lastTime/1e9 > util.BytesToInt64(bs[4:12])/1e9+3 {
	// 				continue
	// 			}
	// 			if endtime > 0 && endtime < util.BytesToInt64(bs[4:12])/1e9 {
	// 				return false
	// 			}
	// 			if err = this.writeToDBbyBytes(bs); err != nil {
	// 				return false
	// 			} else {
	// 				atomic.AddInt64(&this.syncCount, 1)
	// 				txid = util.BytesToInt64(bs[:8])
	// 			}
	// 		} else {
	// 			buf.Write(append(log.STEP, bs...))
	// 		}
	// 	}
	// 	return true
	// })
	var buf bytes.Buffer
	if buf, err = util.UnGzip(body); err == nil {
		bss := bytes.Split(buf.Bytes(), log.STEP)
		for _, bs := range bss {
			if len(bs) == 0 {
				continue
			}
			if bl := parseToBatchLog(bs); bl != nil {
				if lastTime/1e9 > bl.timenano/1e9+3 {
					continue
				}
				if endtime > 0 && endtime < bl.timenano/1e9 {
					break
				}
				if err = this.writeToDBbyBytes(bs); err != nil {
					break
				} else {
					atomic.AddInt64(&this.syncCount, 1)
					txid = bl.txid
				}
			} else {
				err = util.Errors(sys.ERR_LOADLOG)
				logger.Error(err)
			}
		}
	}
	return
}

func (this *logAdmin) ForcedCoverageParseGzfile(body []byte, datetime string, limit int64) (txid int64, err error) {
	endtime := int64(0)
	if t, err := util.TimeStrToTime(datetime); err == nil {
		endtime = t.Unix()
	}
	count := int64(0)
	// err = util.UnGzipByBytes(body, func(bs []byte) bool {
	// 	bss := bytes.Split(append(buf.Bytes(), bs...), log.STEP)
	// 	buf.Reset()
	// 	for i, bs := range bss {
	// 		if i < len(bss)-1 {
	// 			if len(bs) > 12 {
	// 				if (endtime > 0 && endtime < util.BytesToInt64(bs[4:12])/1e9) || (limit > 0 && limit < count) {
	// 					return false
	// 				}
	// 				if err = this.writeToDBbyBytes(bs); err != nil {
	// 					return false
	// 				} else {
	// 					count++
	// 					atomic.AddInt64(&this.syncCount, 1)
	// 					txid = util.BytesToInt64(bs[:8])
	// 				}
	// 			}
	// 		} else if len(bs) > 12 && util.BytesToInt32(bs[:4]) == int32(len(bs)) {
	// 			if (endtime > 0 && endtime < util.BytesToInt64(bs[4:12])/1e9) || (limit > 0 && limit < count) {
	// 				return false
	// 			}
	// 			if err = this.writeToDBbyBytes(bs); err != nil {
	// 				return false
	// 			} else {
	// 				count++
	// 				atomic.AddInt64(&this.syncCount, 1)
	// 				txid = util.BytesToInt64(bs[:8])
	// 			}
	// 		} else {
	// 			buf.Write(append(log.STEP, bs...))
	// 		}
	// 	}
	// 	return true
	// })

	var buf bytes.Buffer
	if buf, err = util.UnGzip(body); err == nil {
		bss := bytes.Split(buf.Bytes(), log.STEP)
		for _, bs := range bss {
			if len(bs) == 0 {
				continue
			}
			if bl := parseToBatchLog(bs); bl != nil {
				if (endtime > 0 && endtime < bl.timenano/1e9) || (limit > 0 && limit < count) {
					break
				}
				if err = this.writeToDBbyBytes(bs); err != nil {
					break
				} else {
					count++
					atomic.AddInt64(&this.syncCount, 1)
					txid = bl.txid
				}
			} else {
				err = util.Errors(sys.ERR_LOADLOG)
				logger.Error(err)
			}
		}
	}
	return
}

func (this *logAdmin) writeKV(key string, txid int64, value []byte) (err error) {
	batch := &BatchPacket{Kvmap: map[string][]byte{key: value}}
	if err = this.writeBatch(util.TimeNano(), txid, batch); err == nil {
		err = level0.Level0.Put(key, value)
	}
	return
}

func (this *logAdmin) writePonBean(timenano int64, pb *PonBean) (err error) {
	if err = this.writeBatch(timenano, pb.Txid, pb.Batch); err == nil {
		err = this._writeToDB(pb)
	}
	return
}

func (this *logAdmin) writeBatch(timenano int64, txid int64, batch *BatchPacket) (err error) {
	batchlog := &batchLog{batch, txid, timenano}
	err = log.BinLog.WriteBytes(batchlog.toBuffer(), timenano)
	return
}

func (this *logAdmin) writeToDBbyBytes(bs []byte) (err error) {
	bpl := parseToBatchLog(bs)
	if err = log.BinLog.WriteBytes(bpl.toBuffer(), bpl.timenano); err == nil {
		err = this._writeToDB(&PonBean{Batch: bpl.batch})
	}
	return
}

func (this *logAdmin) _writeToDB(pb *PonBean) (err error) {
	islocal := true
	if pb.ExecNode != nil {
		islocal = isLocal(pb.ExecNode)
	}
	bp := pb.Batch
	transfer(bp.Kvmap, islocal, pb.ExecNode)
	err = level0.Level0.Batch(bp.Kvmap, bp.Dels)
	return
}

func (this *logAdmin) pullData(uuid int64) (txid int64) {
	defer this.mux.Unlock()
	this.mux.Lock()
	var err error
	if txid, err = this._pullData(uuid, 0); err == nil {
		statAdmin.put(uuid)
		this._txid = txid
	} else {
		logger.Error(err)
	}
	return
}

func (this *logAdmin) _pullData(uuid int64, statKey int64) (_txid int64, err error) {
	logger.Info("pullData:", sys.UUID, " ,uuid>> ", uuid)
	lastTime, txid := this.getStat()
	if statKey > 0 {
		lastTime = statKey
	}
	var ldb *LogDataBean
	if ldb, err = nodeWare.pullData(uuid, 0, _newLogDataBean(1, lastTime, txid, nil, false, false), sys.WaitTimeout, false); err == nil {
		if ldb.FBean != nil && ldb.FBean.Body != nil {
			logger.Info("ldb.FBean.Body.len >>", len(ldb.FBean.Body))
			if _txid, err = this.parseGzfile(0, "", ldb.FBean.Body, ""); err == nil {
				if ldb.HasNext {
					this._pullData(uuid, ldb.FBean.StatKey)
				}
			}
		} else if ldb.NotNull {
			err = util.Errors(sys.ERR_SYNCDATA)
		}
	}
	logger.Warn("pullData end << ", uuid)
	return
}

/**********************************************************************************/

func (this *logAdmin) WriteBackLog(pb *PonBean, uuids []int64) {
	bs := util.TEncode(pb)
	uuidsbs := util.IntArrayToBytes(uuids)
	buf := util.BufferPool.Get(8 + 4 + len(uuidsbs) + len(bs))
	defer util.BufferPool.Put(buf)
	buf.Write(util.Int64ToBytes(util.TimeNano()))     //8
	buf.Write(util.Int32ToBytes(int32(len(uuidsbs)))) //4
	buf.Write(uuidsbs)                                //
	buf.Write(bs)                                     //
	level0.Level0.Put(KeyLevel1.BackKey(fmt.Sprint(pb.Txid)), buf.Bytes())
	return
}

func (this *logAdmin) GetPonBeanByBack(key string) (pb *PonBean, uuids []int64, err error) {
	var bs []byte
	if bs, err = level0.Level0.Get(key); err == nil {
		nano := util.BytesToInt64(bs[0:8])
		if int64(time.Duration(nano)+240*time.Hour) < util.TimeNano() {
			err = util.Errors(sys.ERR_TXOVER)
		} else {
			uuidlen := util.BytesToInt32(bs[8:12])
			uuids = util.BytesToIntArray(bs[12 : 12+uuidlen])
			pb, err = util.TDecode(bs[12+uuidlen:], &PonBean{})
		}
	}
	return
}

func (this *logAdmin) RemoveBackLogUUid(txid int64, uuid int64) (err error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	key := KeyLevel1.BackKey(fmt.Sprint(txid))
	if bs, err := level0.Level0.Get(key); err == nil {
		arrlen := util.BytesToInt32(bs[8:12])
		arr := util.BytesToIntArray(bs[12 : 12+arrlen])
		arr2 := util.ArraySub2(arr, uuid)
		if len(arr2) > 0 {
			nbs := util.IntArrayToBytes(arr2)
			buf := bytes.NewBuffer([]byte{})
			buf.Write(bs[0:8])
			buf.Write(util.Int32ToBytes(int32(len(nbs))))
			buf.Write(nbs)
			buf.Write(bs[12+arrlen:])
			level0.Level0.Put(key, buf.Bytes())
		} else {
			level0.Level0.Del(key)
		}
	}
	return
}

/**********************************************************************************/
func (this *logAdmin) WriteCacheLog(pb *PonBean) (err error) {
	bpbuf := EncodeBatchPacket(pb.Batch)
	defer bpbuf.Free()
	bs := util.SnappyEncode(bpbuf.Bytes()) // util.TEncode(pb.Batch)
	buf := util.BufferPool.Get(8 + 8 + len(bs))
	t := util.TimeNano()
	buf.Write(util.Int64ToBytes(t))
	buf.Write(util.Int64ToBytes(pb.Txid))
	buf.Write(bs)
	return log.CacheLog.Write(buf, pb.Txid, t)
}

func (this *logAdmin) LoadCacheLog() (err error) {
	var bs []byte
	if bs, err = log.CacheLog.Read(); err == nil && bs != nil {
		for _, _bs := range bytes.Split(bs, log.STEP) {
			if len(_bs) > 16 {
				t := util.BytesToInt64(_bs[:8])
				txid := util.BytesToInt64(_bs[8:16])
				bp := DecodeBatchPacket(util.SnappyDecode(_bs[16:]))
				pb := &PonBean{Batch: bp, Txid: txid}
				if err = tlog.writePonBean(t, pb); err != nil {
					fatalError(err)
				}
			}
		}
	}
	return
}

/**********************************************************************************/
func _newLogDataBean(_type int8, lastTime, txid int64, fb *FileBean, hasNext bool, notnull bool) (_r *LogDataBean) {
	_r = NewLogDataBean()
	_r.Type = _type
	_r.LastTime = lastTime
	_r.Txid = txid
	_r.FBean = fb
	_r.HasNext = hasNext
	_r.NotNull = notnull
	return
}
