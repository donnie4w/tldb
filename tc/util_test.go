// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

package tc

import (
	"fmt"
	"testing"
	"tlnetTop/util"
)

func BenchmarkTrans(b *testing.B) {
	bs, _ := valueToBytes([]byte("9"), "13446951783211093631")
	fmt.Println(bs)
	fmt.Println(uint64(util.BytesToInt64(bs)))
}
