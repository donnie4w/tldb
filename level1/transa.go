// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package level1

import (
	. "github.com/donnie4w/tldb/key"
	. "github.com/donnie4w/tldb/level0"
	. "github.com/donnie4w/tldb/lock"
)

var tsLockPool = NewTsLockPool()

type trans struct {
	key4Get string
	kvmap   map[string][]byte
	dels    []string
	del_key map[string]int8
	ts      *TsLock
}

func NewTrans(key string, kvmap map[string][]byte, dels []string) (t *trans) {
	t = &trans{key4Get: key, kvmap: kvmap, dels: dels, del_key: make(map[string]int8, 0)}
	strs := make([]string, 0)
	if key != "" {
		strs = append(strs, key)
	}
	if kvmap != nil {
		for k := range kvmap {
			strs = append(strs, k)
		}
	}
	if dels != nil {
		for _, k := range dels {
			strs = append(strs, k)
		}
	}
	t.ts = tsLockPool.GetLock(strs)
	return
}

func (this *trans) begin() {
	this.ts.Lock()
}

func (this *trans) end() {
	this.ts.UnLock()
	tsLockPool.RmLock(this.ts.GetKey())
}

func (this *trans) batch() (err error) {
	bacmap := make(map[string][]byte, 0)
	if this.kvmap != nil {
		for k, v := range this.kvmap {
			if value0, err := Level0.Get(k); err == nil {
				bacmap[KeyLevel1.TransKey(k)] = value0
			} else {
				this.del_key[k] = 1
				bacmap[KeyLevel1.TransDelKey(k)] = []byte{0}
			}
			bacmap[k] = v
		}
	}
	if this.dels != nil {
		for _, k := range this.dels {
			if value0, err := Level0.Get(k); err == nil {
				bacmap[KeyLevel1.TransKey(k)] = value0
			}
		}
	}
	err = Level0.Batch(bacmap, this.dels)
	return
}

// func (this *trans) get() (value []byte, err error) {
// 	return Level0.Get(KeyLevel1.TransKey(this.key4Get))
// }

func (this *trans) commit() (err error) {
	delarr := make([]string, 0)
	if this.kvmap != nil {
		for k := range this.kvmap {
			delarr = append(delarr, KeyLevel1.TransKey(k))
		}
	}
	if this.dels != nil {
		delarr = append(delarr, this.dels...)
	}
	if len(this.del_key) > 0 {
		for k := range this.del_key {
			delarr = append(delarr, KeyLevel1.TransDelKey(k))
		}
	}
	err = Level0.Batch(nil, delarr)
	return
}

func (this *trans) rollback() (err error) {
	bacmap := make(map[string][]byte, 0)
	delarr := make([]string, 0)
	if this.kvmap != nil {
		for k := range this.kvmap {
			if value0, err := Level0.Get(KeyLevel1.TransKey(k)); err == nil {
				bacmap[k] = value0
			}
			delarr = append(delarr, KeyLevel1.TransKey(k))
		}
	}
	if this.dels != nil {
		for _, k := range this.dels {
			if value0, err := Level0.Get(KeyLevel1.TransKey(k)); err == nil {
				bacmap[k] = value0
			}
			delarr = append(delarr, KeyLevel1.TransKey(k))
		}
	}
	if len(this.del_key) > 0 {
		for k := range this.del_key {
			delarr = append(delarr, KeyLevel1.TransDelKey(k))
			delarr = append(delarr, k)
		}
	}
	err = Level0.Batch(bacmap, delarr)
	return
}
