// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package key

const (
	KEY3_SEQ    = "3_1_"
	KEY3_MAXSEQ = "3_2_"
)

type keyLevel3 struct{}

var KeyLevel3 = &keyLevel3{}

func (this *keyLevel3) KeyMaxDelSeq() string {
	return concat(KEY3_SEQ, "del")
}

func (this *keyLevel3) KeyMaxDelSeqCursor() string {
	return concat(KEY3_SEQ, "del_cursor")
}

func (this *keyLevel3) SeqForDel(seq string) string {
	return concat(KEY3_MAXSEQ, seq)
}
