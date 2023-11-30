// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package stub

type Merge[T any, V any] interface {
	Add(t T)
	Del(id int64)
	Get(id int64) (_r T, ok bool)
	Len() int64
	MergeSize(size int8)
	SetZlib(zlib bool)
	CallBack(cancelfunc func() bool, f func(v V) error) (err error)
}

type RemoteNode struct {
	Addr      string
	UUID      int64
	Host      string
	AdminAddr string
	MQAddr    string
	CliAddr   string
	Stat      int8
	StatDesc  string
}

type Server interface {
	Serve() (err error)
	Close() (err error)
}

type Pool[T any] interface {
	Get(size int) T
	Put(t T) (ok bool)
}
