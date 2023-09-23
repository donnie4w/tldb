// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"context"
	"strings"
	"sync"
	"github.com/donnie4w/gothrift/thrift"
	// "github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tsf/server"
)

const tlContextCtx = "tlContext"

type tlContext struct {
	Id              int64
	Conn            Itnet
	RemoteAddr      string
	RemoteUuid      int64
	RemoteHost      string
	RemoteAdminAddr string
	RemoteMqAddr    string
	RemoteCliAddr   string
	NameSpace       string
	CcPut           int64
	CcGet           int64
	////////////////////////////////
	pongFlag   int64
	transport  thrift.TTransport
	mux        *sync.Mutex
	defaultCtx context.Context
	cancleChan chan byte
	mergeChan  chan *syncBean
	mergeCount int64
	mergemux   *sync.Mutex
	////////////////////////////////
	isServer bool
	isClose  bool
	pingNum  int64
	pongNum  int64
	isAuth   bool
	////////////////////////////////
	stat      sys.STATTYPE
	_stat_seq int64
	////////////////////////////////
	_do_reconn bool
}

func (this *tlContext) SetId(_id int64) {
	this.Id = _id
}

func (this *tlContext) Close() {
	defer recover()
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this.isClose {
		this.isClose = true
		this.transport.Close()
		close(this.cancleChan)
		nodeWare.del(this)
	}
}

func newTlContext(transport thrift.TTransport) (tc *tlContext) {
	tc = &tlContext{Id: newTxId(), mux: new(sync.Mutex), mergeChan: make(chan *syncBean, 1<<17), mergemux: &sync.Mutex{}}
	tc.transport = transport
	tc.Conn = &ItnetWrite{&sync.Mutex{}, transport2Client(transport)}
	return
}

func newTlContext2(socket *TSocket) (tc *tlContext) {
	tc = &tlContext{Id: newTxId(), mux: new(sync.Mutex), mergeChan: make(chan *syncBean, 1<<17), mergemux: &sync.Mutex{}}
	tc.transport = socket
	tc.Conn = &ItnetImpl{socket}
	return
}

func remoteHost(transport thrift.TTransport) (_r string) {
	defer recover()
	if addr := transport.(*thrift.TSocket).Conn().RemoteAddr(); addr != nil {
		if ss := strings.Split(addr.String(), ":"); len(ss) == 2 {
			_r = ss[0]
		}
	}
	return
}

func remoteHost2(tsocket *TSocket) (_r string) {
	defer recover()
	if addr := tsocket.Conn().RemoteAddr(); addr != nil {
		if ss := strings.Split(addr.String(), ":"); len(ss) == 2 {
			_r = ss[0]
		}
	}
	return
}
