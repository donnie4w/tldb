// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
// github.com/donnie4w/tldb

package util

import (
	gobuffer "github.com/donnie4w/gofer/pool/buffer"
	"github.com/donnie4w/gofer/pool/gopool"
	"github.com/donnie4w/tldb/sys"
)

var BufferPool = gobuffer.NewBufferPool(8)

var GoPool = gopool.NewPool(200, sys.GOMAXLIMIT)