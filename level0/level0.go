// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level0

import (
	"fmt"
	"sync"

	"github.com/donnie4w/tldb/db"
	. "github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/keystore"
	"github.com/donnie4w/tldb/log"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"
)

func init() {
	sys.Service.Put(1, Level0)
}

var Level0 = &level0{}

type level0 struct {
	_db db.DBEngine
}

func (this *level0) Serve(wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	this._db = db.New(sys.DBFILEDIR + "/db")
	e1, e2 := Level0.initUUid(), Level0.transRollback()
	if e1 != nil || e2 != nil {
		log.LoggerSys.Error("init failed:", e1)
		panic("init failed")
	}
	log.LogInit()
	keystore.Init()
	log.LoggerSys.Write([]byte{'\n', '\n'})
	log.LoggerSys.Info(sys.SysLog(fmt.Sprint("UUID:", sys.UUID)))
	return
}

func (this *level0) Put(key string, value []byte) (err error) {
	err = this._db.Put([]byte(key), value)
	return
}

func (this *level0) Get(key string) (value []byte, err error) {
	return this._db.Get([]byte(key))
}

func (this *level0) Del(key string) (err error) {
	return this._db.Del([]byte(key))
}

func (this *level0) Has(key string) (_r bool) {
	return this._db.Has([]byte(key))
}

func (this *level0) Batch(kvmap map[string][]byte, dels []string) (err error) {
	return this._db.Batch(kvmap, dels)
}

func (this *level0) GetByPrefix(prefixKey string) (datamap map[string][]byte, err error) {
	return this._db.GetLike([]byte(prefixKey))
}

func (this *level0) GetKeys() ([]string, error) {
	return this._db.GetKeys()
}

func (this *level0) GetKeysPrefix(pre string) (bys []string, err error) {
	return this._db.GetKeysPrefix([]byte(pre))
}

func (this *level0) GetKeysPrefixLimit(pre string, limit int) (bys []string, err error) {
	return this._db.GetKeysPrefixLimit([]byte(pre), limit)
}

func (this *level0) Close() (err error) {
	err = this._db.Close()
	log.LoggerSys.Info(sys.SysLog(fmt.Sprint("close node:", sys.UUID)))
	return
}

// ///////////////////////////////////////////
func (this *level0) initUUid() (err error) {
	if bs, er := this.Get(KeyLevel0.UUID()); er == nil {
		sys.UUID = BytesToInt64(bs)
	} else {
		sys.UUID = int64(NewUUID())
		err = this.Put(KeyLevel0.UUID(), Int64ToBytes(sys.UUID))
	}
	if sys.NAMESPACE != "" {
		Level0.Put(KeyLevel0.NAMESPACE(), []byte(sys.NAMESPACE))
	} else if bs, err := Level0.Get(KeyLevel0.NAMESPACE()); err == nil {
		sys.NAMESPACE = string(bs)
	}
	if bs, er := this.Get(KeyLevel3.KeyMaxDelSeq()); er == nil {
		sys.MAXDELSEQ = BytesToInt64(bs)
	}
	if bs, er := this.Get(KeyLevel3.KeyMaxDelSeqCursor()); er == nil {
		sys.MAXDELSEQCURSOR = BytesToInt64(bs)
	}

	if bs, er := this.Get(KeyLevel1.StatSeq()); er == nil {
		sys.STATSEQ = BytesToInt64(bs)
	}
	return
}

func (this *level0) transRollback() (err error) {
	if ss, er := this.GetKeysPrefix(KeyLevel1.PreTransDelKey()); er == nil {
		ds := make([]string, 0)
		for _, s := range ss {
			ds = append(ds, s, s[len(KeyLevel1.PreTransDelKey()):])
		}
		if err = this.Batch(nil, ds); err != nil {
			return
		}
	}
	if ss, er := this.GetKeysPrefix(KeyLevel1.PreTransKey()); er == nil {
		km := make(map[string][]byte, 0)
		for _, s := range ss {
			if vs, err := this.Get(s[len(KeyLevel1.PreTransKey()):]); err == nil {
				km[s[len(KeyLevel1.PreTransKey()):]] = vs
			}
		}
		if err = this.Batch(km, ss); err != nil {
			return
		}
	}
	return
}
