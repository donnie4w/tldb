// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package tnet

import (
	"github.com/donnie4w/tldb/sys"
)

func AddNode(addr string) (err error) {
	return sys.Client2Serve(addr)
}

func CloseSelf() {
	sys.BroadRmNode()
}
