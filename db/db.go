/**
 * Copyright 2023 tldb Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 *
 * github.com/donnie4w/tldb
 */
package db

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var dbMap map[string]*ldb = make(map[string]*ldb, 0)

var mux = new(sync.Mutex)

type ldb struct {
	db       *leveldb.DB
	dir_name string
}

func New(_dir string) (_r DBEngine) {
	mux.Lock()
	defer mux.Unlock()
	if db, ok := dbMap[_dir]; !ok {
		db = &ldb{dir_name: _dir}
		if err := db.Open(); err == nil {
			dbMap[_dir] = db
		} else {
			panic("init db error:" + err.Error())
		}
		_r = db
	}
	return
}

// //////////////////////////////////////////
func (this *ldb) Open() (err error) {
	options := &opt.Options{
		Filter:                 filter.NewBloomFilter(10),
		OpenFilesCacheCapacity: 1 << 10,
		BlockCacheCapacity:     int(sys.DBBUFFER / 2),
		WriteBuffer:            int(sys.DBBUFFER / 4),
	}
	this.db, err = leveldb.OpenFile(this.dir_name, options)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		this.db, err = leveldb.RecoverFile(this.dir_name, nil)
	}
	return
}

func (this *ldb) Close() (err error) {
	return this.db.Close()
}

func (this *ldb) Has(key []byte) (b bool) {
	b, _ = this.db.Has(key, nil)
	return
}

func (this *ldb) Put(key, value []byte) (err error) {
	return this.db.Put(key, value, &opt.WriteOptions{Sync: sys.SYNC})
}

func (this *ldb) Get(key []byte) (value []byte, err error) {
	value, err = this.db.Get(key, nil)
	return
}

func (this *ldb) GetString(key []byte) (value string, err error) {
	v, er := this.db.Get(key, nil)
	return string(v), er
}

func (this *ldb) Del(key []byte) (err error) {
	return this.db.Delete(key, &opt.WriteOptions{Sync: sys.SYNC})
}

func (this *ldb) BatchPut(kvmap map[string][]byte) (err error) {
	batch := new(leveldb.Batch)
	for k, v := range kvmap {
		batch.Put([]byte(k), v)
	}
	err = this.db.Write(batch, &opt.WriteOptions{Sync: sys.SYNC})
	return
}

func (this *ldb) Batch(put map[string][]byte, del []string) (err error) {
	batch := new(leveldb.Batch)
	if put != nil {
		for k, v := range put {
			batch.Put([]byte(k), v)
		}
	}
	if del != nil {
		for _, v := range del {
			batch.Delete([]byte(v))
		}
	}
	err = this.db.Write(batch, &opt.WriteOptions{Sync: sys.SYNC})
	return
}

func (this *ldb) GetLike(prefix []byte) (datamap map[string][]byte, err error) {
	iter := this.db.NewIterator(util.BytesPrefix(prefix), nil)
	if iter != nil {
		defer iter.Release()
		datamap = make(map[string][]byte, 0)
		for iter.Next() {
			bs := make([]byte, len(iter.Value()))
			copy(bs, iter.Value())
			datamap[string(iter.Key())] = bs
		}
	}
	err = iter.Error()
	return
}

// all key
func (this *ldb) GetKeys() (bys []string, err error) {
	iter := this.db.NewIterator(nil, nil)
	defer iter.Release()
	bys = make([]string, 0)
	for iter.Next() {
		bys = append(bys, string(iter.Key()))
	}
	err = iter.Error()
	return
}

// all key by prefix
func (this *ldb) GetKeysPrefix(prefix []byte) (bys []string, err error) {
	iter := this.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([]string, 0)
	for iter.Next() {
		bys = append(bys, string(iter.Key()))
	}
	err = iter.Error()
	return
}

// all key by prefix
func (this *ldb) GetKeysPrefixLimit(prefix []byte, limit int) (bys []string, err error) {
	iter := this.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([]string, 0)
	i := 0
	for iter.Next() {
		if i > limit {
			break
		}
		bys = append(bys, string(iter.Key()))
		i++
	}
	err = iter.Error()
	return
}

/*
*
Start of the key range, include in the range.
Limit of the key range, not include in the range.
*/
func (this *ldb) GetIterLimit(prefix string, limit string) (datamap map[string][]byte, err error) {
	iter := this.db.NewIterator(&util.Range{Start: []byte(prefix), Limit: []byte(limit)}, nil)
	if iter != nil {
		defer iter.Release()
		datamap = make(map[string][]byte, 0)
		for iter.Next() {
			data, er := this.db.Get(iter.Key(), nil)
			if er == nil {
				datamap[string(iter.Key())] = data
			}
		}
		err = iter.Error()
	}
	return
}

func (this *ldb) Snapshot() (*leveldb.Snapshot, error) {
	return this.db.GetSnapshot()
}

// ////////////////////////////////////////////////////////
type BakStub struct {
	Key   []byte
	Value []byte
}

func (this *BakStub) copy(k, v []byte) {
	this.Key, this.Value = make([]byte, len(k)), make([]byte, len(v))
	copy(this.Key, k)
	copy(this.Value, v)
}

func (this *ldb) BackupToDisk(filename string, prefix []byte) error {
	defer myRecover()
	snap, err := this.Snapshot()
	if err != nil {
		return err
	}
	defer snap.Release()
	bs := _TraverseSnap(snap, prefix)
	b, e := Encode(bs)
	if e != nil {
		return e
	}
	f, er := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if er != nil {
		return er
	}
	defer f.Close()
	_, er = f.Write(b)
	return er
}

func RecoverBackup(filename string) (bs []*BakStub) {
	defer myRecover()
	f, er := os.Open(filename)
	if er == nil {
		defer f.Close()
	} else {
		return
	}
	var buf bytes.Buffer
	_, err := io.Copy(&buf, f)
	if err == nil {
		Decode(buf.Bytes(), &bs)
	}
	return
}

func (this *ldb) LoadDataFile(filename string) (err error) {
	bs := RecoverBackup(filename)
	for _, v := range bs {
		err = this.Put(v.Key, v.Value)
	}
	return
}

func (this *ldb) LoadBytes(buf []byte) (err error) {
	var bs []*BakStub
	err = Decode(buf, &bs)
	if err == nil {
		for _, v := range bs {
			err = this.Put(v.Key, v.Value)
		}
	}
	return
}

func _TraverseSnap(snap *leveldb.Snapshot, prefix []byte) (bs []*BakStub) {
	ran := new(util.Range)
	if prefix != nil {
		ran = util.BytesPrefix(prefix)
	} else {
		ran = nil
	}
	iter := snap.NewIterator(ran, nil)
	defer iter.Release()
	bs = make([]*BakStub, 0)
	for iter.Next() {
		ss := new(BakStub)
		ss.copy(iter.Key(), iter.Value())
		bs = append(bs, ss)
	}
	return
}

func myRecover() {
	if err := recover(); err != nil {
	}
}
