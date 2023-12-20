// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package key

import (
	"strconv"
	"strings"
)

const (
	KEY2_ID_      = "2_0_"
	KEY2_IDX_     = "2_1_"
	KEY2_SEQ      = "2_2_"
	KEY2_PTE      = "2_3_"
	KEY2_DB       = "2_4_"
	KEY2_TABLE    = "2_5_"
	KEY2_MQ_TABLE = "2_6_"
)

type keyLevel2 struct{}

var KeyLevel2 = &keyLevel2{}

func (this *keyLevel2) SeqName(tablename string) string {
	return concat(KEY2_ID_, tr(tablename), "_id_")
}

func (this *keyLevel2) SeqKey(tablename string, _seq int64) string {
	return concat(this.SeqName(tablename), itoa(_seq))
}

func (this *keyLevel2) IndexName(tablename, idx_name, idx_value string) string {
	return concat(KEY2_IDX_, tr(tablename), "_", idx_name, "_", idx_value, "_")
}

func (this *keyLevel2) IndexKey(tablename, idx_name, idx_value string, _seq int64) string {
	return concat(this.IndexName(tablename, idx_name, idx_value), itoa(_seq))
}

func (this *keyLevel2) GetIdxSeqByKeySubName(idx_key, idx_name string) (seq int64) {
	if len(idx_key) >= len(idx_name) {
		if strings.Contains(idx_key, idx_name) {
			s := idx_key[len(idx_name):]
			seq, _ = strconv.ParseInt(s, 10, 64)
		}
	}
	return
}

func (this *keyLevel2) MaxSeqForId(tablename string) string {
	return concat(KEY2_SEQ, tr(tablename), "_id")
}

func (this *keyLevel2) MaxSeqForIdx(tablename, idx_name, idx_value string) string {
	return concat(KEY2_SEQ, tr(tablename), "_", idx_name, "_", idx_value)
}

func (this *keyLevel2) PteToIdxStub(tablename string, id_value int64) string {
	return concat(KEY2_PTE, tr(tablename), "_id_", itoa(id_value))
}

func (this *keyLevel2) Tables(tablename string) (ss []string) {
	ss = make([]string, 0)
	tname := tr(tablename)
	ss = append(ss, concat(KEY2_ID_, tname))
	ss = append(ss, concat(KEY2_IDX_, tname))
	ss = append(ss, concat(KEY2_SEQ, tname))
	ss = append(ss, concat(KEY2_PTE, tname))
	ss = append(ss, concat(KEY2_TABLE, tname))
	ss = append(ss, concat(KEY2_MQ_TABLE, tname))
	return
}

func (this *keyLevel2) DataBase(dbname string) string {
	return concat(KEY2_DB, dbname)
}

func (this *keyLevel2) Table(tablename string) string {
	return concat(KEY2_TABLE, tr(tablename))
}

func (this *keyLevel2) TableMQ(tablename string) string {
	return concat(KEY2_MQ_TABLE, tr(tablename))
}

func tr(tablename string) (_r string) {
	_r = ui32toa(crc_32([]byte(tablename)))
	return
}

const salt = "a0e$$@=kel1385&%*&&${{}]^|||???"

func Topic(tablename string) (_r string) {
	_r = tablename + salt
	return
}
