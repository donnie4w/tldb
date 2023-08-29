// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
package main

import (
	_ "github.com/donnie4w/tldb/level0"
	_ "github.com/donnie4w/tldb/level1"
	. "github.com/donnie4w/tldb/sys"
	_ "github.com/donnie4w/tldb/tc"
	_ "github.com/donnie4w/tldb/tlcli"
)

func main() {
	Start()
}
