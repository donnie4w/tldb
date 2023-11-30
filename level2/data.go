// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level2

import (
	"sync"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

var Tldb = &tldb{mux: &sync.RWMutex{}, tables: NewMapL[string, *TableStruct](), _t: 0}
var TldbMq = &tldb{mux: &sync.RWMutex{}, tables: NewMapL[string, *TableStruct](), _t: 1}

type tldb struct {
	mux    *sync.RWMutex
	tables *MapL[string, *TableStruct]
	_t     int8
}

func (this *tldb) GetTable(tablename string) (t *TableStruct, ok bool) {
	defer this.mux.RUnlock()
	this.mux.RLock()
	return this.tables.Get(tablename)
}

func (this *tldb) LoadTableInfo() {
	defer this.mux.Unlock()
	this.mux.Lock()
	var dbs []*TableStub
	if this._t == 1 {
		dbs = Level2.LoadMQTableInfo()
	} else {
		dbs = Level2.LoadTableInfo()
	}
	this.tables = NewMapL[string, *TableStruct]()
	if dbs != nil {
		for _, db := range dbs {
			this.tables.Put(db.Tablename, NewTableStruct(db))
		}
	}
}

func (this *tldb) GetTableList() *MapL[string, *TableStruct] {
	defer this.mux.RUnlock()
	this.mux.RLock()
	return this.tables
}

func (this *tldb) addTable(db *TableStub) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.tables.Put(db.Tablename, NewTableStruct(db))
}

func (this *tldb) check(ts *TableStub) (err error) {
	defer this.mux.RUnlock()
	this.mux.RLock()
	if t, ok := this.tables.Get(ts.Tablename); ok {
		if ts.Field == nil {
			err = util.Errors(sys.ERR_NO_MATCH_PARAM)
			return
		}
		for k := range ts.Field {
			if !t.Fields.Has(k) {
				err = util.Errors(sys.ERR_COLUMN_NOEXIST)
				break
			}
		}
		if ts.Idx != nil {
			for k := range ts.Idx {
				if !t.Idx.Has(k) {
					err = util.Errors(sys.ERR_IDX_NOEXIST)
					break
				}
			}
		}
		ts.Idx = util.SyncMap2Map(t.Idx)
	} else {
		err = util.Errors(sys.ERR_TABLE_NOEXIST)
	}
	return
}

func (this *tldb) checkGet(db *TableStub) (err error) {
	defer this.mux.RUnlock()
	this.mux.RLock()
	if t, ok := this.tables.Get(db.Tablename); ok {
		if db.Field != nil {
			for k := range db.Field {
				if !t.Fields.Has(k) {
					err = util.Errors(sys.ERR_COLUMN_NOEXIST)
					break
				}
			}
		}
		if db.Idx != nil {
			for k := range db.Idx {
				if !t.Idx.Has(k) {
					err = util.Errors(sys.ERR_IDX_NOEXIST)
					break
				}
			}
		}
	} else {
		err = util.Errors(sys.ERR_TABLE_NOEXIST)
	}
	return
}

type TableStruct struct {
	IdMax     int64
	IdxCount  int64
	KeyCount  int64
	Tablename string
	Fields    *Map[string, []byte]
	Idx       *Map[string, int8]
}

func NewTableStruct(db *TableStub) *TableStruct {
	f := util.Map2SyncMap(db.Field)
	i := util.Map2SyncMap(db.Idx)
	return &TableStruct{Tablename: db.Tablename, Fields: f, Idx: i}
}

/*********************************************************************/
// type timeLock struct {
// 	m   *LinkedMap[int64, chan bool]
// 	mux *sync.Mutex
// }

// func newTimeLock() *timeLock {
// 	return &timeLock{NewLinkedMap[int64, chan bool](), &sync.Mutex{}}
// }

// func (this *timeLock) lock(seq int64) (_r bool) {
// 	ch := make(chan bool, 1)
// 	this.m.Put(seq, ch)
// 	select {
// 	case <-ch:
// 		_r = true
// 	case <-time.After(sys.WaitTimeout):
// 		_r = false
// 		this.m.Del(seq)
// 		close(ch)
// 	}
// 	return
// }

// func (this *timeLock) unLockBack() {
// 	defer myRecovr()
// 	defer this.mux.Unlock()
// 	this.mux.Lock()
// 	k := this.m.Back()
// 	if v, ok := this.m.Get(k); ok {
// 		this.m.Del(k)
// 		close(v)
// 	}
// }

/****************************************************************/
