// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

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
