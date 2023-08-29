// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

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
