// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package db

var default_db *ldb

func AddObject(_tableKey string, a any) {
	Table[any](default_db.dir_name).AddObject(_tableKey, a)
}

func GetAndSetId[T any](_idx_seq string) (id int64) {
	return Table[T](default_db.dir_name).GetAndSetId(_idx_seq)
}
func GetIdSeqValue[T any]() (id int64) {
	return Table[T](default_db.dir_name).GetIdSeqValue()
}
func GetObjectByOrder[T any](_tablename, _idx_id_name string, startId, count int64) (ts []*T) {
	return Table[T](default_db.dir_name).GetObjectByOrder(_tablename, _idx_id_name, startId, count)
}

func AddValue(key string, value []byte) error {
	return Table[any](default_db.dir_name).AddValue(key, value)
}

func GetValue(key string) (value []byte, err error) {
	return Table[any](default_db.dir_name).GetValue(key)
}

func DelKey(key string) (err error) {
	return Table[any](default_db.dir_name).DelKey(key)
}

func Insert(a any) (err error) {
	return Table[any](default_db.dir_name).Insert(a)
}
func Update(a any) (err error) {
	return Table[any](default_db.dir_name).Update(a)
}
func Delete(a any) (err error) {
	return Table[any](default_db.dir_name).Delete(a)
}
func DeleteWithId[T any](id int64) (err error) {
	return Table[T](default_db.dir_name).DeleteWithId(id)
}
func DeleteWithKey(table_key string) (err error) {
	return Table[any](default_db.dir_name).DeleteWithKey(table_key)
}

/*
start  :  table  start id
end    :  table  end id
*/
func Selects[T any](start, end int64) (_r []*T) {
	return Table[T](default_db.dir_name).Selects(start, end)
}

/*
_id :  table id

	one return
*/
func SelectOne[T any](_id int64) (_r *T) {
	return Table[T](default_db.dir_name).SelectOne(_id)
}

/*
idx_name :  index name
_idx_value:  index value

	one return
*/
func SelectOneByIdxName[T any](idx_name, _idx_value string) (_r *T) {
	return Table[T](default_db.dir_name).SelectOneByIdxName(idx_name, _idx_value)
}

/*
idx_name :  index name
_idx_value:  index value

	multiple return
*/
func SelectByIdxName[T any](idx_name, _idx_value string) (_r []*T) {
	return Table[T](default_db.dir_name).SelectByIdxName(idx_name, _idx_value)
}

/*
idx_name :  index name
idxValues:  index value array
startId  :  start number
limit    :  maximum return number
*/
func SelectByIdxNameLimit[T any](idx_name string, idxValues []string, startId, limit int64) (_r []*T) {
	return Table[T](default_db.dir_name).SelectByIdxNameLimit(idx_name, idxValues, startId, limit)
}

func BuildIndex[T any]() (err error, _r string) {
	return Table[T](default_db.dir_name).BuildIndex()
}
