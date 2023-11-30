// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package tc

import (
	"fmt"
	"testing"

	"github.com/donnie4w/gofer/util"
)

func BenchmarkTrans(b *testing.B) {
	bs, _ := valueToBytes([]byte("9"), "13446951783211093631")
	fmt.Println(bs)
	fmt.Println(uint64(util.BytesToInt64(bs)))
}
