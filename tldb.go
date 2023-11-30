// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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


//tldb database
//https://github.com/donnie4w/tldb

//database client
//java:  https://github.com/donnie4w/tlcli-j  
//go:    https://github.com/donnie4w/tlcli-go
//python https://github.com/donnie4w/tlcli-py

//orm 
//java:  https://github.com/donnie4w/tlorm-java  
//go:    https://github.com/donnie4w/tlorm-go

//MQ client
//java:  https://github.com/donnie4w/tlmq-j  
//go:    https://github.com/donnie4w/tlmq-go
//python https://github.com/donnie4w/tlmq-py
//js     https://github.com/donnie4w/tlmq-js

//website
//https://tlnet.top
//Email:donnie4w@gmail.com