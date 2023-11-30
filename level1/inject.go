// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

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
var DeleteBatch func(call int, tablename string, ids []int64) (err error)
var SelectByIdxDescLimit func(call int, table_name string, idx_name string, _idx_value []byte, startId int64, limit int64) (_r []*TableStub, err error)
var SelectByIdxAscLimit func(call int, table_name string, idx_name string, _idx_value []byte, startId int64, limit int64) (_r []*TableStub, err error)
var SelectIdByIdxSeq func(call int, table_name string, idx_name string, _idx_value []byte, seq int64) (_r int64, err error)
var ClusPub func(mqType int8, bs []byte) (err error)
