// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level2

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync/atomic"

	. "github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/level0"
	"github.com/donnie4w/tldb/level1"
	. "github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"
)

func init() {
	level1.AddInject(NewTxId(), loadtableInfo)
	level1.Insert = Level2.Insert
	level1.InsertMq = Level2.InsertMq
	level1.SelectId = Level2.SelectId
	level1.SelectIdByIdx = Level2.SelectIdByIdx
	level1.SelectById = Level2.SelectById
	level1.SelectsByIdLimit = Level2.SelectsByIdLimit
	level1.SelectByIdx = Level2.SelectByIdx
	level1.SelectsByIdx = Level2.SelectsByIdx
	level1.SelectsByIdxLimit = Level2.SelectsByIdxLimit
	level1.Update = Level2.Update
	level1.Delete = Level2.Delete
	level1.CreateTable = Level2.CreateTable
	level1.AlterTable = Level2.AlterTable
	level1.CreateTableMq = Level2.CreateTableMq
	level1.DropTable = Level2.DropTable
	level1.LoadTableInfo = Level2.LoadTableInfo
	level1.LoadMQTableInfo = Level2.LoadMQTableInfo
	sys.Export = Level2.export
	sys.CcGet = Level2.ccGet
	sys.CcPut = Level2.ccPut
	sys.CountGet = Level2.getCount
	sys.CountPut = Level2.putCount
}

func loadtableInfo() {
	Tldb.LoadTableInfo()
	TldbMq.LoadTableInfo()
}

func newIdxStub() (is *IdxStub) {
	is = NewIdxStub()
	is.IdxMap = make(map[string]string, 0)
	return
}

func decode2IdxStub(bs []byte) (_idx *IdxStub) {
	_idx, _ = TDecode(bs, NewIdxStub())
	return
}

var Level2 = newLevel2()

type level2 struct {
	pLock   *LimitLock
	gLock   *LimitLock
	strLock *Strlock
}

func newLevel2() *level2 {
	return &level2{pLock: NewLimitLock(int(sys.COCURRENT_PUT), sys.ReadTimeout), gLock: NewLimitLock(int(sys.COCURRENT_GET), sys.WaitTimeout), strLock: NewStrlock(1 << 14)}
}

func (this *level2) Insert(call int, db *TableStub) (seq int64, err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, isClus := false, int64(0), level1.IsClusRun()
	if isClus && !sys.IsRUN() {
		isProxy = true
	} //else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
	// 	isProxy=true
	// }
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: db}, proxyuuid, 1, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				seq = tp.Seq
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
		return
	}

	if err = checkOp(db); err != nil {
		return
	}
	if sys.DBMode == 1 {
		if err = Tldb.check(db); err != nil {
			return
		}
	}
	tablename := db.Tablename
	if seq, err = level1.Level1.GetAndIncrKey(KeyLevel2.MaxSeqForId(tablename)); err == nil && seq > 0 {
		db.ID = seq
		bacMap := make(map[string][]byte, 0)
		_id_key := KeyLevel2.SeqKey(tablename, seq)
		if err = saveIdxs(db, seq, bacMap); err == nil {
			bacMap[_id_key] = EncodeTableStub(db) //TEncode(db)
			err = level1.Level1.Batch(1, bacMap, nil)
		}
	} else if err == nil && seq == 0 {
		err = Errors(sys.ERR_INCR_SEQ)
	}
	return
}

// ////////////////////////////////////////////////////
func saveIdxs(db *TableStub, id_value int64, bacMap map[string][]byte) (err error) {
	idxstub := newIdxStub()
	for k := range db.Idx {
		if bs, ok := db.Field[k]; ok {
			if err = _saveIdxs(db.Tablename, id_value, k, bs, bacMap, idxstub); err != nil {
				break
			}
			// keyvalue := idxValue(bs)
			// var _i_v int64
			// if _i_v, err = level1.Level1.GetAndIncrKey(KeyLevel2.MaxSeqForIdx(db.Tablename, k, keyvalue)); _i_v > 0 {
			// 	_idx_key := KeyLevel2.IndexKey(db.Tablename, k, keyvalue, _i_v)
			// 	bacMap[_idx_key] = Int64ToBytes(id_value)
			// 	putPteKey(db.Tablename, k, _idx_key, id_value, bacMap, idxstub)
			// } else {
			// 	err = Errors(sys.ERR_INCR_SEQ)
			// 	logger.Error("saveIdexs err:" + KeyLevel2.MaxSeqForIdx(db.Tablename, k, keyvalue))
			// 	break
			// }
		}
	}
	return
}

func _saveIdxs(tablename string, table_id int64, fieldname string, fieldvalue []byte, bacMap map[string][]byte, idxstub *IdxStub) (err error) {
	keyvalue := idxValue(fieldvalue)
	if _i_v, _ := level1.Level1.GetAndIncrKey(KeyLevel2.MaxSeqForIdx(tablename, fieldname, keyvalue)); _i_v > 0 {
		_idx_key := KeyLevel2.IndexKey(tablename, fieldname, keyvalue, _i_v)
		bacMap[_idx_key] = Int64ToBytes(table_id)
		putPteKey(tablename, fieldname, _idx_key, table_id, bacMap, idxstub)
	} else {
		err = Errors(sys.ERR_INCR_SEQ)
	}
	return
}

func putPteKey(table_name, idx_name, _idx_key string, _table_id_value int64, bacMap map[string][]byte, idxstub *IdxStub) {
	_pte_key := KeyLevel2.PteToIdxStub(table_name, _table_id_value)
	idxstub.IdxMap[idx_name] = _idx_key
	bacMap[_pte_key] = TEncode(idxstub)
	return
}

/****************************************************************************/
func (this *level2) InsertMq(call int, db *TableStub) (seq int64, err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: db}, 0, 2, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				seq = tp.Seq
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
		return
	}
	if err = checkOp(db); err != nil {
		return
	}

	// if err = TlMQ.check(db); err != nil {
	// 	return
	// }
	tablename := db.Tablename
	if seq, err = level1.Level1.GetAndIncrKey(KeyLevel2.MaxSeqForId(tablename)); err == nil && seq > 0 {
		db.ID = seq
		bacMap := make(map[string][]byte, 0)
		_id_key := KeyLevel2.SeqKey(tablename, seq)
		bacMap[_id_key] = EncodeTableStub(db) //TEncode(db)
		if err = saveIdxs(db, seq, bacMap); err == nil {
			err = level1.Level1.Batch(1, bacMap, nil)
		}
	} else if err == nil && seq == 0 {
		err = Errors(sys.ERR_INCR_SEQ)
	}
	return
}

/*************************************************************************/
func (this *level2) SelectId(call int, db *TableStub) (_r int64, err error) {
	defer myRecovr()
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: db}, proxyuuid, 3, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.Seq
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		var bs []byte
		if bs, err = level1.Level1.GetLocal(KeyLevel2.MaxSeqForId(db.Tablename)); err == nil {
			_r = BytesToInt64(bs)
		}
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

/*************************************************************************/
func (this *level2) SelectIdByIdx(call int, table_name, idx_name string, _idx_value []byte) (_r int64, err error) {
	defer myRecovr()
	if table_name == "" || idx_name == "" || _idx_value == nil {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
		return
	}
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{TableName: table_name, IdxName: idx_name, IdxValue: _idx_value}, proxyuuid, 4, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.Seq
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		var bs []byte
		if bs, err = level1.Level1.GetLocal(KeyLevel2.MaxSeqForIdx(table_name, idx_name, idxValue(_idx_value))); err == nil {
			_r = BytesToInt64(bs)
		}
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

/***********************************************************************************/
func (this *level2) SelectById(call int, db *TableStub) (_r *TableStub, err error) {
	defer myRecovr()
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: db}, proxyuuid, 5, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.Stub
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		if !checkSelectBeen(db) || db.ID <= 0 {
			err = Errors(sys.ERR_NO_MATCH_PARAM)
			return
		}
		// if sys.DBMode == 1 {
		// 	if err = Tldb.checkGet(db); err != nil {
		// 		return
		// 	}
		// }
		_r, err = this._selectById(db.Tablename, db.ID)
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

func (this *level2) SelectsByIdLimit(call int, db *TableStub, start, limit int64) (_r []*TableStub, err error) {
	defer myRecovr()
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: db, Start: start, Limit: limit}, proxyuuid, 6, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.StubArray
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		if !checkSelectBeen(db) || limit == 0 {
			err = Errors(sys.ERR_NO_MATCH_PARAM)
			return
		}
		_r = make([]*TableStub, 0)
		ms := KeyLevel2.MaxSeqForId(db.Tablename)
		if bs, err := level1.Level1.GetLocal(ms); err == nil {
			end := BytesToInt64(bs)
			c, i := int64(0), int64(0)
			for c < limit && start+i <= end {
				if t, err := this._selectById(db.Tablename, start+i); err == nil && t != nil {
					_r = append(_r, t)
					c++
				}
				i++
			}
		}
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

func (this *level2) _selectById(tablename string, id_value int64) (_r *TableStub, err error) {
	_id_key := KeyLevel2.SeqKey(tablename, id_value)
	this.strLock.RLock(_id_key)
	defer this.strLock.RUnlock(_id_key)
	if bs, _err := level1.Level1.GetValue(_id_key); _err == nil {
		_r = DecodeTableStub(bs)
	} else {
		err = Errors(sys.ERR_DATA_NOEXIST)
	}
	return
}

func (this *level2) SelectByIdx(call int, table_name, idx_name string, _idx_value []byte) (_r *TableStub, err error) {
	defer myRecovr()
	if table_name == "" || idx_name == "" || _idx_value == nil {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
		return
	}
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{TableName: table_name, IdxName: idx_name, IdxValue: _idx_value}, proxyuuid, 7, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.Stub
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		seqName := KeyLevel2.MaxSeqForIdx(table_name, idx_name, idxValue(_idx_value))
		var ids []byte
		if ids, err = level1.Level1.GetLocal(seqName); err == nil {
			idSeq := BytesToInt64(ids)
			for j := int64(1); j <= idSeq; j++ {
				_idx_key := KeyLevel2.IndexKey(table_name, idx_name, idxValue(_idx_value), j)
				if idbuf, _ := level1.Level1.GetLocal(_idx_key); idbuf != nil {
					tid := BytesToInt64(idbuf)
					if _r, err = this._selectById(table_name, tid); err == nil {
						return
					}
				}
			}
		}
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

func (this *level2) SelectsByIdx(call int, table_name, idx_name string, _idx_value []byte) (_r []*TableStub, err error) {
	defer myRecovr()
	if table_name == "" || idx_name == "" || _idx_value == nil {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
		return
	}
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{TableName: table_name, IdxName: idx_name, IdxValue: _idx_value}, proxyuuid, 8, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.StubArray
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		_r = make([]*TableStub, 0)
		seqName := KeyLevel2.MaxSeqForIdx(table_name, idx_name, idxValue(_idx_value))
		if ids, err := level1.Level1.GetLocal(seqName); err == nil {
			idSeq := BytesToInt64(ids)
			for j := int64(1); j <= idSeq; j++ {
				_idx_key := KeyLevel2.IndexKey(table_name, idx_name, idxValue(_idx_value), j)
				if idbuf, _ := level1.Level1.GetLocal(_idx_key); idbuf != nil {
					tid := BytesToInt64(idbuf)
					if t, err := this._selectById(table_name, tid); err == nil {
						if t != nil {
							_r = append(_r, t)
						}
					}
				}
			}
		}
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

func (this *level2) SelectsByIdxLimit(call int, table_name, idx_name string, idxValues [][]byte, startId, limit int64) (_r []*TableStub, err error) {
	defer myRecovr()
	if table_name == "" || idx_name == "" || idxValues == nil || startId < 0 || limit == 0 {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
		return
	}
	defer this.gLock.Unlock()
	if err = this.gLock.Lock(); err != nil {
		return
	}
	isProxy, proxyuuid, hasRemoteRun := false, int64(0), len(level1.GetRemoteRunUUID()) > 0
	if hasRemoteRun && !sys.IsRUN() {
		isProxy = true
	} else if proxyuuid = level1.LB(); proxyuuid > 0 && call == 0 {
		isProxy = true
	}
	if isProxy {
		if tp, exec := level1.BroadcastProxy(&TableParam{TableName: table_name, IdxName: idx_name, IdxValues: idxValues, Start: startId, Limit: limit}, proxyuuid, 9, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			} else {
				_r = tp.StubArray
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
	} else if sys.IsRUN() {
		_r = make([]*TableStub, 0)
		i, count := int64(0), limit
		smap := make(map[string]byte, 0)
		for _, v := range idxValues {
			if count <= 0 {
				return
			}
			idxvalue := idxValue(v)
			if _, ok := smap[idxvalue]; ok {
				continue
			} else {
				smap[idxvalue] = 0
			}
			seqName := KeyLevel2.MaxSeqForIdx(table_name, idx_name, idxvalue)
			if ids, er := level1.Level1.GetLocal(seqName); er == nil {
				id := BytesToInt64(ids)
				for j := int64(1); j <= id; j++ {
					if count <= 0 {
						return
					}
					_idx_key := KeyLevel2.IndexKey(table_name, idx_name, idxvalue, j)
					if level1.Level1.Has(_idx_key) {
						if i < startId {
							i++
						} else {
							if idbuf, _ := level1.Level1.GetLocal(_idx_key); idbuf != nil {
								tid := BytesToInt64(idbuf)
								if t, er := this._selectById(table_name, tid); er == nil {
									_r = append(_r, t)
									count--
								}
							}
						}
					}
				}
			}
		}
	} else {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	return
}

func (this *level2) Update(call int, ts *TableStub) (err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: ts}, 0, 10, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
		return
	}
	if err = checkOp(ts); err != nil {
		return
	}
	t, ok := Tldb.GetTable(ts.Tablename)
	if !ok {
		err = Errors(sys.ERR_TABLE_NOEXIST)
		return
	}
	if sys.DBMode == 1 {
		if err = Tldb.check(ts); err != nil {
			return
		}
	}
	_id_key := KeyLevel2.SeqKey(ts.Tablename, ts.ID)
	this.strLock.Lock(_id_key)
	defer this.strLock.Unlock(_id_key)
	if _bs, _err := level1.Level1.GetValue(_id_key); _err == nil {
		bacMap := make(map[string][]byte, 0)
		delArr := make([]string, 0)
		idxDelMap := make(map[string][]byte, 0)
		// if BytesEqual(bs, _bs) {
		// 	return
		// }
		_ts := DecodeTableStub(_bs)
		change := false
		for k, v := range ts.Field {
			if t.Fields.Has(k) {
				_ts.Field[k] = v
				change = true
			}
		}
		if !change {
			return
		}
		bacMap[_id_key] = EncodeTableStub(_ts)
		_pte_key := KeyLevel2.PteToIdxStub(ts.Tablename, ts.ID)
		var bs []byte
		if bs, err = level1.Level1.GetLocal(_pte_key); err == nil {
			is := decode2IdxStub(bs)
			reset := false

			for fieldname := range _ts.Idx {
				if _, ok := is.IdxMap[fieldname]; !ok {
					if v, ok := _ts.Field[fieldname]; ok {
						if err = _saveIdxs(ts.Tablename, ts.ID, fieldname, v, bacMap, is); err != nil {
							return
						}
					}
				}
			}

			for idx_name, _idx_key := range is.IdxMap {
				if _, ok := ts.Field[idx_name]; !ok {
					continue
				}
				if new_idx_value_bs, ok := _ts.Field[idx_name]; ok {
					new_idx_value := idxValue(new_idx_value_bs)
					new_pre_idx_key := KeyLevel2.IndexName(_ts.Tablename, idx_name, new_idx_value)
					if !strings.HasPrefix(_idx_key, new_pre_idx_key) {
						delArr = append(delArr, _idx_key)
						var _idx_seq_value int64
						if _idx_seq_value, err = level1.Level1.GetAndIncrKey(KeyLevel2.MaxSeqForIdx(_ts.Tablename, idx_name, new_idx_value)); err == nil {
							new_idx_key := KeyLevel2.IndexKey(_ts.Tablename, idx_name, new_idx_value, _idx_seq_value)
							bacMap[new_idx_key] = Int64ToBytes(_ts.ID)
							is.IdxMap[idx_name] = new_idx_key
							reset = true
						}
						idxDelMap[KeyLevel3.SeqForDel(fmt.Sprint(atomic.AddInt64(&sys.MAXDELSEQ, 1)))] = []byte(_idx_key)
						idxDelMap[KeyLevel3.KeyMaxDelSeq()] = Int64ToBytes(*&sys.MAXDELSEQ)
					}
				} else {
					delete(is.IdxMap, idx_name)
					delArr = append(delArr, _idx_key)
				}
			}
			if err == nil && reset {
				bacMap[_pte_key] = TEncode(is)
			}
		} else {
			err = saveIdxs(ts, ts.ID, bacMap)
		}
		if err == nil {
			if err = level1.Level1.Batch(2, bacMap, delArr); err == nil {
				level0.Level0.Batch(idxDelMap, nil)
			}
		}
	} else {
		err = Errors(sys.ERR_DATA_NOEXIST)
	}
	return
}

// delete
func (this *level2) Delete(call int, ts *TableStub) (err error) {
	if ts == nil || ts.ID <= 0 {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
		return
	}
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: ts}, 0, 11, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
		return
	}
	idxDelMap := make(map[string][]byte, 0)
	if ts.ID > 0 {
		_id_key := KeyLevel2.SeqKey(ts.Tablename, ts.ID)
		this.strLock.Lock(_id_key)
		defer this.strLock.Unlock(_id_key)
		if level1.Level1.Has(_id_key) {
			dels := make([]string, 0)
			dels = append(dels, _id_key)
			_pte_key := KeyLevel2.PteToIdxStub(ts.Tablename, ts.ID)
			if bs, err := level1.Level1.GetLocal(_pte_key); err == nil && bs != nil {
				is := decode2IdxStub(bs)
				for _, _idx_key := range is.IdxMap {
					dels = append(dels, _idx_key)
					idxDelMap[KeyLevel3.SeqForDel(fmt.Sprint(atomic.AddInt64(&sys.MAXDELSEQ, 1)))] = []byte(_idx_key)
					idxDelMap[KeyLevel3.KeyMaxDelSeq()] = Int64ToBytes(*&sys.MAXDELSEQ)
				}
				dels = append(dels, _pte_key)
			}
			if err == nil {
				if err = level1.Level1.Batch(3, nil, dels); err == nil {
					go level0.Level0.Batch(idxDelMap, nil)
				}
			}
		} else {
			err = Errors(sys.ERR_DATA_NOEXIST)
		}
	} else {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
	}
	return
}

// create table
func (this *level2) CreateTable(ts *TableStub) (err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: ts}, 0, 12, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
		return
	}
	if err = checkOp(ts); err != nil {
		return
	}
	defer this.strLock.Unlock(ts.Tablename)
	this.strLock.Lock(ts.Tablename)
	tablenameKey := KeyLevel2.Table(ts.Tablename)
	if _, er := level1.Level1.GetLocal(tablenameKey); er == nil {
		err = Errors(sys.ERR_TABLE_EXIST)
	} else {
		err = this._AlterTable(ts, 0)
	}
	return
}

func (this *level2) AlterTable(ts *TableStub) (err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: ts}, 0, 13, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			}
		} else {
			err = Errors(sys.ERR_PROXY)
			LoggerError.Error("altertable error:", err)
		}
		return
	}
	defer this.strLock.Unlock(ts.Tablename)
	this.strLock.Lock(ts.Tablename)
	return this._AlterTable(ts, 0)
}

func (this *level2) _AlterTable(ts *TableStub, t int8) (err error) {
	if err = checkTableStruct(ts); err != nil {
		return
	}
	for k := range ts.Field {
		if strings.TrimSpace(k) == "" {
			delete(ts.Field, k)
		}
	}
	if ts.Idx != nil {
		im := make(map[string]int8, 0)
		for k := range ts.Idx {
			if _, ok := ts.Field[k]; ok {
				im[k] = 0
			}
		}
		ts.Idx = im
	}
	var tablenameKey string
	if t == 0 {
		tablenameKey = KeyLevel2.Table(ts.Tablename)
	} else {
		tablenameKey = KeyLevel2.TableMQ(Topic(ts.Tablename))
	}
	bacMap := make(map[string][]byte, 0)
	bacMap[tablenameKey] = EncodeTableStub(ts)
	if err = level1.Level1.Batch(1, bacMap, nil); err == nil {
		if level1.BroadcastReinit(nil).Error() == nil {
			if t == 0 {
				Tldb.LoadTableInfo()
			} else {
				TldbMq.LoadTableInfo()
			}
		}
	}
	return
}

func (this *level2) CreateTableMq(ts *TableStub) (err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: ts}, 0, 14, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			}
		} else {
			err = Errors(sys.ERR_PROXY)
			LoggerError.Error("createtablemq error:", err)
		}
		return
	}
	defer this.strLock.Unlock(ts.Tablename)
	this.strLock.Lock(ts.Tablename)
	tablenameKey := KeyLevel2.TableMQ(Topic(ts.Tablename))
	if _, er := level1.Level1.GetLocal(tablenameKey); er == nil {
		err = Errors(sys.ERR_TABLE_EXIST)
	} else {
		err = this._AlterTable(ts, 1)
	}
	return
}

func (this *level2) DropTable(ts *TableStub) (err error) {
	defer myRecovr()
	defer this.pLock.Unlock()
	if err = this.pLock.Lock(); err != nil {
		return
	}
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{Stub: ts}, 0, 15, 1); tp != nil && exec.Error() == nil {
			if tp.Err != "" {
				err = errors.New(tp.Err)
			}
		} else {
			err = Errors(sys.ERR_PROXY)
		}
		return
	}
	if err = checkStat(nil); err != nil {
		return
	}
	defer this.strLock.Unlock(ts.Tablename)
	this.strLock.Lock(ts.Tablename)
	ss := KeyLevel2.Tables(ts.Tablename)
	limit := 1 << 15
	for {
		keys := make([]string, 0)
		for _, s := range ss {
			if klike, err := level1.Level1.GetLocalKeysLikeLimit(s, limit); err == nil {
				keys = append(keys, klike...)
			}
			if len(keys) >= limit {
				break
			}
		}
		if len(keys) > 0 {
			if err = level1.Level1.Batch(1, nil, keys); err == nil {
				Tldb.LoadTableInfo()
				TldbMq.LoadTableInfo()
				level1.BroadcastReinit(nil)
			}
		} else {
			break
		}
	}

	return
}

func (this *level2) LoadTableInfo() (tss []*TableStub) {
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{}, 0, 16, 1); tp != nil && exec.Error() == nil {
			tss = tp.StubArray
		}
		return
	}
	if sys.IsRUN() {
		if datamap, err := level1.Level1.GetLocalLike(KEY2_TABLE); err == nil {
			tss = make([]*TableStub, 0)
			for _, v := range datamap {
				ts := DecodeTableStub(v)
				tss = append(tss, ts)
			}
		}
	}
	return
}

func (this *level2) LoadMQTableInfo() (tss []*TableStub) {
	if level1.IsClusRun() && !sys.IsRUN() {
		if tp, exec := level1.BroadcastProxy(&TableParam{}, 0, 17, 1); tp != nil && exec.Error() == nil {
			tss = tp.StubArray
		}
		return
	}
	if sys.IsRUN() {
		if datamap, err := level1.Level1.GetLocalLike(KEY2_MQ_TABLE); err == nil {
			tss = make([]*TableStub, 0)
			for _, v := range datamap {
				if ts := DecodeTableStub(v); ts != nil {
					tss = append(tss, ts)
				}
			}
		}
	}
	return
}

func (this *level2) ccPut() int64 {
	return this.pLock.Cc()
}

func (this *level2) ccGet() int64 {
	return this.gLock.Cc()
}

func (this *level2) putCount() int64 {
	return this.pLock.LockCount()
}

func (this *level2) getCount() int64 {
	return this.gLock.LockCount()
}

func (this *level2) export(tablename string) (_r bytes.Buffer, err error) {
	ss := KeyLevel2.Tables(tablename)
	buf := bytes.NewBuffer([]byte{})
	for _, s := range ss {
		if keys, err := level1.Level1.GetLocalKeysLike(s); err == nil && len(keys) > 0 {
			kvmap := make(map[string][]byte, 0)
			for _, key := range keys {
				if strings.HasPrefix(key, KEY2_ID_) {
					if v, err := level1.Level1.GetValue(key); err == nil {
						kvmap[key] = v
					}
				} else {
					if v, err := level1.Level1.GetLocal(key); err == nil {
						kvmap[key] = v
					}
				}
			}
			bp := &level1.BatchPacket{Kvmap: kvmap}
			buf.Write(level1.WriteBatchLogBuffer(bp).Bytes())
		}
	}
	if buf.Len() > 0 {
		_r, err = GzipWrite(buf.Bytes())
	} else {
		err = Errors(sys.ERR_TABLE_NOEXIST)
	}
	return
}

/***********************************************************/
func idxValue(bs []byte) (_r string) {
	if len(bs) < 20 {
		_r = string(bs)
	} else {
		_r = fmt.Sprint(CRC64(bs))
	}
	return
}

/***********************************************************/
func myRecovr() {
	if err := recover(); err != nil {
		LoggerError.Error(string(debug.Stack()))
	}
}

func checkSelectBeen(db *TableStub) bool {
	if db != nil && db.Tablename != "" {
		return true
	}
	return false
}

func checkDataBeen(db *TableStub) bool {
	if db != nil && db.Tablename != "" && db.Field != nil && len(db.Field) > 0 {
		return true
	}
	return false
}

func checkTableStruct(ts *TableStub) (err error) {
	if !checkDataBeen(ts) {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
	} else {
		for k := range ts.Field {
			if len(k) > 80 {
				err = Errors(sys.ERR_TABLE_FEILD_EXIST)
				break
			}
		}
	}
	return
}

func checkOp(db *TableStub) (err error) {
	if !sys.IsRUN() {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	if !level1.IsClusRun() {
		err = Errors(sys.ERR_NO_CLUSTER)
	}
	if !checkDataBeen(db) {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
	}
	return
}

func checkStat(db *TableStub) (err error) {
	if !sys.IsRUN() {
		err = Errors(sys.ERR_NO_RUNSTAT)
	}
	if !level1.IsClusRun() {
		err = Errors(sys.ERR_NO_CLUSTER)
	}
	if db != nil && db.Tablename == "" {
		err = Errors(sys.ERR_NO_MATCH_PARAM)
	}
	return
}