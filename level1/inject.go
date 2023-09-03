// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import . "github.com/donnie4w/tldb/stub"

var Insert func(call int, db *TableStub) (seq int64, err error)
var InsertMq func(call int, db *TableStub) (seq int64, err error)
var SelectId func(call int, db *TableStub) (_r int64, err error)
var SelectIdByIdx func(call int, table_name, idx_name string, _idx_value []byte) (_r int64, err error)
var SelectById func(call int, db *TableStub) (_r *TableStub, err error)
var SelectsByIdLimit func(call int, db *TableStub, start, limit int64) (_r []*TableStub, err error)
var SelectByIdx func(call int, table_name, idx_name string, _idx_value []byte) (_r *TableStub, err error)
var SelectsByIdx func(call int, table_name, idx_name string, _idx_value []byte) (_r []*TableStub, err error)
var SelectsByIdxLimit func(call int, table_name, idx_name string, idxValues [][]byte, startId, limit int64) (_r []*TableStub, err error)
var Update func(call int, ts *TableStub) (err error)
var Delete func(call int, ts *TableStub) (err error)
var CreateTable func(ts *TableStub) (err error)
var AlterTable func(ts *TableStub) (err error)
var CreateTableMq func(ts *TableStub) (err error)
var DropTable func(ts *TableStub) (err error)
var LoadTableInfo func() (tss []*TableStub)
var LoadMQTableInfo func() (tss []*TableStub)
var DeleteBatches func(call int, tablename string, fromId, toId int64) (err error)
var ClusPub func(mqType int8, bs []byte) (err error)
