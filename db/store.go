// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package db

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	. "github.com/donnie4w/tldb/util"
)

var rwlock = new(sync.RWMutex)

type idxStub struct {
	IdxMap map[string]string
}

func newIdxStub() *idxStub {
	return &idxStub{make(map[string]string, 0)}
}
func (this *idxStub) encode() (bs []byte) {
	bs, _ = Encode(this)
	return
}
func (this *idxStub) put(idxName, idxKey string) {
	this.IdxMap[idxName] = idxKey
}
func (this *idxStub) has(idxName string) (b bool) {
	_, b = this.IdxMap[idxName]
	return
}
func (this *idxStub) get(idxName string) (s string) {
	s, _ = this.IdxMap[idxName]
	return
}
func decodeIdx(bs []byte) (_idx *idxStub) {
	_idx = new(idxStub)
	Decode(bs, _idx)
	return
}

type _Table[T any] struct {
	db *ldb
}

func (this _Table[T]) AddObject(_tableKey string, e any) {
	bys, _ := Encode(e)
	this.db.Put([]byte(_tableKey), bys)
}

func getObjectByLike[T any](db *ldb, prefix string) (ts []*T) {
	m, err := db.GetLike([]byte(prefix))
	if err == nil {
		ts = make([]*T, 0)
		for _, v := range m {
			t := new(T)
			Decode(v, t)
			ts = append(ts, t)
		}
	}
	return
}

// get and set Id,id increment 1
func (this _Table[T]) GetAndSetId(_idx_seq string) (id int64) {
	rwlock.Lock()
	defer rwlock.Unlock()
	ids, err := this.db.Get([]byte(_idx_seq))
	if err == nil && ids != nil {
		id = BytesToInt64(ids)
	}
	atomic.AddInt64(&id, 1)
	this.db.Put([]byte(_idx_seq), Int64ToBytes(id))
	return
}

func (this _Table[T]) GetIdSeqValue() (id int64) {
	var a T
	tname := getObjectName(a)
	idxSeqName := idx_id_seq(tname)
	ids, err := this.db.Get([]byte(idxSeqName))
	if err == nil && ids != nil {
		id = BytesToInt64(ids)
	}
	return
}

func (this _Table[T]) GetObjectByOrder(_tablename, _idx_id_name string, startId, count int64) (ts []*T) {
	ts = make([]*T, 0)
	ids, err := this.db.Get([]byte(_idx_id_name))
	var id int64
	if err == nil && ids != nil {
		id = BytesToInt64(ids)
	}
	for i := startId; i < count; i++ {
		if i <= id {
			v, err := this.db.Get([]byte(idx_id_key(_tablename, i)))
			if err == nil && v != nil {
				t := new(T)
				Decode(v, t)
				ts = append(ts, t)
			} else {
				count++
			}
		}
	}
	return
}

func (this _Table[T]) AddValue(key string, value []byte) error {
	return this.db.Put([]byte(key), value)
}

func (this _Table[T]) GetValue(key string) (value []byte, err error) {
	return this.db.Get([]byte(key))
}

func (this _Table[T]) DelKey(key string) (err error) {
	return this.db.Del([]byte(key))
}

func (this _Table[T]) hasKey(key string) bool {
	return this.db.Has([]byte(key))
}

func getObjectName(a any) (tname string) {
	t := reflect.TypeOf(a)
	if t.Kind() != reflect.Pointer {
		tname = strings.ToLower(t.Name())
	} else {
		tname = strings.ToLower(t.Elem().Name())
	}
	if tname == "" {
		panic("getObjectName error: table name is empty")
	}
	return
}

// func setId(a any, id_value int64) {
// 	v := reflect.ValueOf(a).Elem()
// 	fmt.Println("name:", reflect.TypeOf(a).Name())
// 	v.FieldByNameFunc(func(s string) bool {
// 		return strings.ToLower(s) == "id"
// 	}).SetInt(id_value)
// 	fmt.Println(a)
// }

func getTableIdValue(a any) (_r int64) {
	v := reflect.ValueOf(a)
	if reflect.TypeOf(a).Kind() == reflect.Pointer {
		v = v.Elem()
	}
	id_v := v.FieldByNameFunc(func(s string) bool {
		return strings.ToLower(s) == "id"
	})
	if id_v.Kind() == reflect.Pointer {
		_r = *(*int64)(id_v.UnsafePointer())
	} else {
		_r = id_v.Int()
	}
	return
}

func (this _Table[T]) Insert(a any) (err error) {
	if isPointer(a) {
		table_name := getObjectName(a)
		_table_id_value := this.GetAndSetId(idx_id_seq(table_name))
		_idx_key := idx_id_key(table_name, _table_id_value)
		v := reflect.ValueOf(a).Elem()
		id_v := v.FieldByNameFunc(func(s string) bool {
			return strings.ToLower(s) == "id"
		})
		if id_v.Kind() == reflect.Pointer {
			id_v.Set(reflect.ValueOf(&_table_id_value))
		} else {
			id_v.SetInt(_table_id_value)
		}
		this.AddObject(_idx_key, a)
		go this._saveIdx_(a, table_name, _table_id_value)
	} else {
		err = errors.New("insert object must be pointer")
	}
	return
}

func (this _Table[T]) _saveIdx_(a any, tablename string, _table_id_value int64) {
	t := reflect.TypeOf(a).Elem()
	v := reflect.ValueOf(a).Elem()
	for i := 0; i < t.NumField(); i++ {
		idxName := t.Field(i).Name
		if checkIndexField(idxName, t.Field(i).Tag) {
			f := v.FieldByName(idxName)
			idx_value, _ := getValueFromkind(f)
			if idx_value != nil {
				this._insertWithTableId(tablename, strings.ToLower(idxName), *idx_value, _table_id_value)
			}
		}
	}
}

func getValueFromkind(f reflect.Value) (_v *string, e error) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	var v string
	isSet := false
	if f.CanInt() {
		v, isSet = fmt.Sprint(f.Int()), true
	} else if f.CanFloat() {
		v, isSet = fmt.Sprint(f.Float()), true
	} else if f.CanUint() {
		v, isSet = fmt.Sprint(f.Uint()), true
	} else if f.Kind() == reflect.String {
		v, isSet = f.String(), true
	} else if f.Kind() == reflect.Pointer {
		switch f.Interface().(type) {
		case *int:
			v, isSet = fmt.Sprint(*(*int)(f.UnsafePointer())), true
		case *int8:
			v, isSet = fmt.Sprint(*(*int8)(f.UnsafePointer())), true
		case *int16:
			v, isSet = fmt.Sprint(*(*int16)(f.UnsafePointer())), true
		case *int32:
			v, isSet = fmt.Sprint(*(*int32)(f.UnsafePointer())), true
		case *int64:
			v, isSet = fmt.Sprint(*(*int64)(f.UnsafePointer())), true
		case *uint:
			v, isSet = fmt.Sprint(*(*uint)(f.UnsafePointer())), true
		case *uint16:
			v, isSet = fmt.Sprint(*(*uint16)(f.UnsafePointer())), true
		case *uint32:
			v, isSet = fmt.Sprint(*(*uint32)(f.UnsafePointer())), true
		case *uint64:
			v, isSet = fmt.Sprint(*(*uint64)(f.UnsafePointer())), true
		case *float32:
			v, isSet = fmt.Sprint(*(*float32)(f.UnsafePointer())), true
		case *float64:
			v, isSet = fmt.Sprint(*(*float64)(f.UnsafePointer())), true
		case *string:
			v, isSet = *(*string)(f.UnsafePointer()), true
		}
	}
	if isSet {
		_v = &v
	} else {
		e = errors.New(fmt.Sprint(f.String(), " build index  error"))
	}
	return
}

// key: tablename_idxName_idxValue_idSeq: idvalue
func (this _Table[T]) _insertWithTableId(table_name, idx_name, idx_value string, _table_id_value int64) (err error) {
	_idx_seq_value := this.GetAndSetId(idx_seq(table_name, idx_name, idx_value))
	_idx_key := idx_key(table_name, idx_name, idx_value, _idx_seq_value)
	err = this.db.Put([]byte(_idx_key), Int64ToBytes(_table_id_value))
	this.putPteKey(table_name, idx_name, _idx_key, _table_id_value)
	return
}

func (this _Table[T]) putPteKey(table_name, idx_name, _idx_key string, _table_id_value int64) {
	_pte_key := pte_key(table_name, _table_id_value)
	if bs, err := this.db.Get([]byte(_pte_key)); err == nil {
		is := decodeIdx(bs)
		is.put(idx_name, _idx_key)
		this.db.Put([]byte(_pte_key), is.encode())
	} else {
		is := newIdxStub()
		is.put(idx_name, _idx_key)
		this.db.Put([]byte(_pte_key), is.encode())
	}
}

func (this _Table[T]) updatePteKey(a any, table_name string, _table_id_value int64) {
	_pte_key := pte_key(table_name, _table_id_value)
	if bs, err := this.db.Get([]byte(_pte_key)); err == nil {
		is := decodeIdx(bs)
		rv := reflect.ValueOf(a).Elem()
		reset := false
		for idx_name, _idx_key := range is.IdxMap {
			// f := rv.FieldByName(idx_name)
			f := rv.FieldByNameFunc(func(s string) bool {
				return strings.ToLower(s) == idx_name
			})
			new_idx_value, e := getValueFromkind(f)
			if e == nil {
				new_pre_idx_key := idx_key_prefix(table_name, idx_name, new_idx_value)
				if !strings.Contains(_idx_key, new_pre_idx_key) {
					this.db.Del([]byte(_idx_key))
					_idx_seq_value := this.GetAndSetId(idx_seq(table_name, idx_name, new_idx_value))
					new_idx_key := idx_key(table_name, idx_name, new_idx_value, _idx_seq_value)
					err = this.db.Put([]byte(new_idx_key), Int64ToBytes(_table_id_value))
					is.put(idx_name, new_idx_key)
					reset = true
				}
			}
		}
		if reset {
			this.db.Put([]byte(_pte_key), is.encode())
		}
	}
}

func (this _Table[T]) Update(a any) (err error) {
	if !isPointer(a) {
		return errors.New("update object must be pointer")
	}
	table_name := getObjectName(a)
	_table_id_value := getTableIdValue(a)
	_idx_key := idx_id_key(table_name, _table_id_value)
	if this.hasKey(_idx_key) {
		this.AddObject(_idx_key, a)
		go this.updatePteKey(a, table_name, _table_id_value)
	} else {
		err = errors.New(fmt.Sprint("key[", _idx_key, "] is not exist"))
	}
	return
}

func (this _Table[T]) Delete(a any) (err error) {
	table_name := getObjectName(a)
	_table_id_value := getTableIdValue(a)
	return this._delete(table_name, _table_id_value)
}

func (this _Table[T]) DeleteWithId(id int64) (err error) {
	var a T
	table_name := getObjectName(a)
	return this._delete(table_name, id)
}

func (this _Table[T]) DeleteWithKey(table_key string) (err error) {
	table_name, id := parse_idx_id_key(table_key)
	return this._delete(table_name, id)
}

func (this _Table[T]) _delete(table_name string, _table_id_value int64) (err error) {
	if _table_id_value == 0 {
		return errors.New("The ID value for deletion is not set")
	}
	_idx_key := idx_id_key(table_name, _table_id_value)
	this.DelKey(_idx_key)
	_pte_key := pte_key(table_name, _table_id_value)
	if bs, err := this.db.Get([]byte(_pte_key)); err == nil {
		is := decodeIdx(bs)
		for _, pte_idx_key := range is.IdxMap {
			this.db.Del([]byte(pte_idx_key))
		}
	}
	this.db.Del([]byte(_pte_key))
	return
}

/*
start  :  table  start id
end    :  table  end id
*/
func (this _Table[T]) Selects(start, end int64) (_r []*T) {
	var a T
	if !isStruct(a) {
		panic("type of genericity must be struct")
	}
	tname := getObjectName(a)
	idxSeqName := idx_id_seq(tname)
	_r = this.GetObjectByOrder(tname, idxSeqName, start, end)
	return
}

/*
_id :  table id

	one return
*/
func (this _Table[T]) SelectOne(_id int64) (_r *T) {
	var a T
	if !isStruct(a) {
		panic("type of genericity must be struct")
	}
	tname := getObjectName(a)
	_r = this._selectoneFromId(tname, _id)
	return
}

func (this _Table[T]) _selectoneFromId(tablename string, _id int64) (_r *T) {
	v, err := this.db.Get([]byte(idx_id_key(tablename, _id)))
	if err == nil && v != nil {
		_r = new(T)
		Decode(v, _r)
	}
	return
}

/*
idx_name :  index name
_idx_value:  index value

	one return
*/
func (this _Table[T]) SelectOneByIdxName(idx_name, _idx_value string) (_r *T) {
	var a T
	if !isStruct(a) {
		panic("type of genericity must be struct")
	}
	idx_name = parseIdxName[T](idx_name)
	tname := getObjectName(a)
	idxSeqName := idx_seq(tname, idx_name, _idx_value)
	ids, err := this.db.Get([]byte(idxSeqName))
	if err == nil && ids != nil {
		id := BytesToInt64(ids)
		for j := int64(1); j <= id; j++ {
			_idx_key := idx_key(tname, idx_name, _idx_value, j)
			idbuf, _ := this.db.Get([]byte(_idx_key))
			tid := BytesToInt64(idbuf)
			_r = this._selectoneFromId(tname, tid)
			if _r != nil {
				return
			}
		}
	}
	return
}

/*
idx_name :  index name
_idx_value:  index value

	multiple return
*/
func (this _Table[T]) SelectByIdxName(idx_name, _idx_value string) (_r []*T) {
	var a T
	if !isStruct(a) {
		panic("type of genericity must be struct")
	}
	idx_name = parseIdxName[T](idx_name)
	tname := getObjectName(a)
	_r = make([]*T, 0)
	idxSeqName := idx_seq(tname, idx_name, _idx_value)
	ids, err := this.db.Get([]byte(idxSeqName))
	if err == nil && ids != nil {
		id := BytesToInt64(ids)
		for j := int64(1); j <= id; j++ {
			_idx_key := idx_key(tname, idx_name, _idx_value, j)
			idbuf, _ := this.db.Get([]byte(_idx_key))
			tid := BytesToInt64(idbuf)
			t := this._selectoneFromId(tname, tid)
			if t != nil {
				_r = append(_r, t)
			}
		}
	}
	return
}

/*
idx_name :  index name
idxValues:  index value array
startId  :  start number
limit    :  maximum return number
*/
func (this _Table[T]) SelectByIdxNameLimit(idx_name string, idxValues []string, startId, limit int64) (_r []*T) {
	var a T
	if !isStruct(a) {
		panic("type of genericity must be struct")
	}
	idx_name = parseIdxName[T](idx_name)
	tname := getObjectName(a)
	_r = make([]*T, 0)
	i, count := int64(0), limit
	for _, v := range idxValues {
		if count <= 0 {
			return
		}
		idxSeqName := idx_seq(tname, idx_name, v)
		ids, err := this.db.Get([]byte(idxSeqName))
		if err == nil && ids != nil {
			id := BytesToInt64(ids)
			for j := int64(1); j <= id; j++ {
				if count <= 0 {
					return
				}
				_idx_key := idx_key(tname, idx_name, v, j)
				if this.db.Has([]byte(_idx_key)) {
					if i < startId {
						i++
					} else {
						idbuf, _ := this.db.Get([]byte(_idx_key))
						tid := BytesToInt64(idbuf)
						t := this._selectoneFromId(tname, tid)
						if t != nil {
							_r = append(_r, t)
							count--
						}
					}
				}
			}
		}
	}
	return
}

func (this _Table[T]) BuildIndex() (err error, _r string) {
	var a T
	if !isStruct(a) {
		err = errors.New("type of genericity must be struct")
		return
	}
	table_name := getObjectName(a)
	t := reflect.TypeOf(a)
	mustBuild := false
	idx_array := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		idx_name := strings.ToLower(t.Field(i).Name)
		if checkIndexField(idx_name, t.Field(i).Tag) {
			_idx_seq := idx_seq(table_name, idx_name, "")
			is := getObjectByLike[int64](this.db, _idx_seq)
			if is == nil || len(is) == 0 {
				mustBuild = true
				idx_array = append(idx_array, t.Field(i).Name)
			}
		}
	}
	if mustBuild {
		_r = fmt.Sprintln("BuildIndex table[", table_name, "],field:", idx_array, "")
		idxSeqName := idx_id_seq(table_name)
		ids, err := this.db.Get([]byte(idxSeqName))
		var id int64
		if err == nil && ids != nil {
			id = BytesToInt64(ids)
			for i := int64(1); i <= id; i++ {
				s := this.SelectOne(i)
				if s != nil {
					v := reflect.ValueOf(s).Elem()
					for _, field_name := range idx_array {
						f := v.FieldByName(field_name)
						idx_value, e := getValueFromkind(f)
						if idx_value != nil {
							this._insertWithTableId(table_name, strings.ToLower(field_name), *idx_value, i)
						} else {
							err = e
						}
					}
				}
			}
		}
	} else {
		err = errors.New("no need build index")
	}
	return
}

func checkIndexField(field_name string, tag reflect.StructTag) (b bool) {
	return strings.HasSuffix(field_name, "_") || string(tag) == "idx" || tag.Get("idx") == "1"
}

func parseIdxName[T any](idx_name string) string {
	if !strings.HasSuffix(idx_name, "_") {
		var a T
		t := reflect.TypeOf(a)
		isTagIdx := false
		isOtherIdx := false
		for i := 0; i < t.NumField(); i++ {
			field_name := t.Field(i).Name
			if checkIndexField("", t.Field(i).Tag) && idx_name == field_name {
				isTagIdx = true
				break
			}
			if strings.ToLower(field_name) == fmt.Sprint(strings.ToLower(idx_name), "_") {
				isOtherIdx = true
			}
		}
		if !isTagIdx && isOtherIdx {
			idx_name = fmt.Sprint(idx_name, "_")
		}
	}
	return strings.ToLower(idx_name)
}

func isPointer(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Pointer
}

func isStruct(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Struct
}

/*index key:tablename idx_name seq_value: user_age_22_1*/
func idx_key(tablename, idx_name string, idx_value any, id_seq_value int64) string {
	return fmt.Sprint(idx_key_prefix(tablename, idx_name, idx_value), id_seq_value)
}

/*prefix index  key: tablename idx_name  : user_id_ or user_id_100_*/
func idx_key_prefix(tablename, idx_name string, idx_value any) string {
	return fmt.Sprint("1_", tablename, "_", idx_name, "_", idx_value, "_")
}

/*
table id:
index key:tablename idx_name seq_value: user_id_1
*/
func idx_id_key(tablename string, id_value int64) string {
	return fmt.Sprint(idx_id_key_prefix(tablename), id_value)
}

/*
table  id:
prefix index  key: tablename idx_name  : user_id_ or user_id_100_
*/
func idx_id_key_prefix(tablename string) string {
	return fmt.Sprint("0_", tablename, "_id_")
}

func parse_idx_id_key(table_key string) (table_name string, id int64) {
	id_site := strings.LastIndex(table_key, "_id_")
	table_name = table_key[2:id_site]
	_id := table_key[id_site+4:]
	id, _ = strconv.ParseInt(_id, 10, 0)
	return
}

/*seq index key idx_tablenaem idx_name: idx_user_id*/
func idx_seq(tablename, idx_name string, idx_value any) string {
	return fmt.Sprint("idx_", tablename, "_", idx_name, "_", idx_value)
}

/*table id*/
func idx_id_seq(tablename string) string {
	return fmt.Sprint("idx_", tablename, "_id")
}

/*id  to  indexs*/
func pte_key(tablename string, id_value int64) string {
	return fmt.Sprint("pte_", tablename, "_id_", id_value)
}

func Table[T any](dbname string) _Table[T] {
	db, _ := dbMap[dbname]
	return _Table[T]{db}
}
