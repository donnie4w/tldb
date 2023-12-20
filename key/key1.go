// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package key

const (
	TRANS_KEY = "1_0_"
	TRANS_DEL = "1_1_"
	BACK_KEY  = "1_2_"
	BATCH_KEY = "1_3_"
	STAT_KEY  = "1_4_"
)

type keyLevel1 struct {
}

var KeyLevel1 = &keyLevel1{}

func (this *keyLevel1) TransKey(key string) string {
	return concat(this.PreTransKey(), key)
}

func (this *keyLevel1) PreTransKey() string {
	return TRANS_KEY
}

func (this *keyLevel1) TransDelKey(key string) string {
	return concat(this.PreTransDelKey(), key)
}

func (this *keyLevel1) PreTransDelKey() string {
	return TRANS_DEL
}

func (this *keyLevel1) BackKey(key string) string {
	return concat(BACK_KEY, key)
}

func (this *keyLevel1) BackPrefix() string {
	return BACK_KEY
}

func (this *keyLevel1) BatchKey(key string) string {
	return concat(BATCH_KEY, key)
}

func (this *keyLevel1) StatSeq() string {
	return concat(STAT_KEY, "seq_")
}

func (this *keyLevel1) StatKey(seq int64) string {
	return concat(STAT_KEY, "id_", itoa(seq))
}
