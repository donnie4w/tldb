// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level1

import (
	"sort"
	"strings"

	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"
)

type pon struct {
	store_num int
	uuids     []int64
}

func Pon() *pon {
	return &pon{sys.STORENODENUM, nodeWare.GetAllRunUUID()}
}

func PonUUIDS(uuids []int64) *pon {
	return &pon{sys.STORENODENUM, uuids}
}
func (this *pon) ExcuNode(key string) int64 {
	if sys.IsStandAlone() {
		return sys.UUID
	}
	umap := make(map[uint32]int64, 0)
	keylen := len([]byte(key))
	for _, uuid := range this.uuids {
		bs := make([]byte, keylen+8)
		copy(bs[0:keylen], []byte(key))
		copy(bs[keylen:], Int64ToBytes(uuid))
		umap[CRC32(bs)] = uuid
	}
	return max(umap)
}

func (this *pon) ExcuNodes(key string) (ids []int64) {
	length := this.store_num
	if len(this.uuids) < length || length == 0 {
		length = len(this.uuids)
	}
	umap := make(map[uint32]int64, 0)
	keybslen := len([]byte(key))
	for _, uuid := range this.uuids {
		bs := make([]byte, keybslen+8)
		copy(bs[0:keybslen], []byte(key))
		copy(bs[keybslen:], Int64ToBytes(uuid))
		umap[CRC32(bs)] = uuid
	}
	return sortArr(umap, length)
}

func (this *pon) ExcuBatchNodes(kvmap map[string][]byte, delarr []string) (ids []int64) {
	length := this.store_num
	if len(this.uuids) <= length || length == 0 {
		return this.uuids
	}
	keylen := 0
	if kvmap != nil {
		keylen = keylen + len(kvmap)
	}
	if delarr != nil {
		keylen = keylen + len(delarr)
	}
	keys := make([]string, keylen)
	i := 0
	if kvmap != nil {
		for k := range kvmap {
			keys[i] = k
			i++
		}
	}
	if delarr != nil {
		for _, k := range delarr {
			keys[i] = k
			i++
		}
	}
	sort.Strings(keys)
	key := strings.Join(keys, "")
	umap := make(map[uint32]int64, 0)
	keybslen := len([]byte(key))
	for _, uuid := range this.uuids {
		bs := make([]byte, keybslen+8)
		copy(bs[0:keybslen], []byte(key))
		copy(bs[keybslen:], Int64ToBytes(uuid))
		umap[CRC32(bs)] = uuid
	}
	return sortArr(umap, length)
}

// func (this *pon) checkuuid(key string, checkUUids []int64) bool {
// 	uuid := this.ExcuNode(key)
// 	for _, v := range checkUUids {
// 		if uuid == v {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (this *pon) compareArrHalf(u1, u2 []int64) bool {
// 	if len(u1) > len(u2) {
// 		return this.compareArrHalf(u2, u1)
// 	}
// 	f := 0
// 	_u2 := make([]int64, len(u2))
// 	copy(_u2, u2)
// 	for _, v := range u1 {
// 		for i, v2 := range _u2 {
// 			if v == v2 {
// 				_u2[i] = 17
// 				f++
// 				break
// 			}
// 		}
// 	}
// 	return f >= len(u1)/2+1
// }

// func (this *pon) compareArrALL(u1, u2 []int64) bool {
// 	if len(u1) != len(u2) {
// 		return false
// 	}
// 	_u2 := make([]int64, len(u2))
// 	copy(_u2, u2)
// 	for _, v := range u1 {
// 		f := false
// 		for i, v2 := range _u2 {
// 			if v == v2 {
// 				_u2[i] = 17
// 				f = true
// 				break
// 			}
// 		}
// 		if !f {
// 			return false
// 		}
// 	}
// 	return true
// }

func sortArr(umap map[uint32]int64, length int) (uuids []int64) {
	ustrs := make([]uint32, len(umap))
	i := 0
	for n := range umap {
		ustrs[i] = n
		i++
	}
	sort.Slice(ustrs, func(i, j int) bool { return ustrs[i] > ustrs[j] })
	uuids = make([]int64, length)
	for i := 0; i < length; i++ {
		uuids[i] = umap[ustrs[i]]
	}
	return
}

func max(umap map[uint32]int64) (uuid int64) {
	var _r uint32
	for u := range umap {
		if _r < u {
			_r = u
		}
	}
	uuid = umap[_r]
	return
}
