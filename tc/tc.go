// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tc

import (
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tlnet"
)

func init() {
	tlnet.SetLogOFF()
	sys.Service.Put(2, adminservice)
	sys.Service.Put(3, mqservice)
	sys.Cmd = cmd.Connect
}
