// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

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
