// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package level1

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/donnie4w/gothrift/thrift"

	"github.com/donnie4w/tldb/sys"
)

var _protocolFactory = thrift.NewTCompactProtocolFactoryConf(&thrift.TConfiguration{})
var tlclientServer = &tlClientServer{}

type tlClientServer struct {
	isClose         bool
	serverTransport *thrift.TServerSocket
}

func (this *tlClientServer) server(_addr string, itnet Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext)) (err error) {
	var serverTransport *thrift.TServerSocket
	if serverTransport, err = thrift.NewTServerSocketTimeout(_addr, sys.SocketTimeout); err == nil {
		this.serverTransport = serverTransport
		processor := NewItnetProcessor(itnet)
		tserver4 := thrift.NewTSimpleServer4(processor, serverTransport, nil, nil)
		if err = tserver4.Listen(); err == nil {
			sys.ADDR = _addr
			for {
				var transport thrift.TTransport
				if transport, err = tserver4.ServerTransport().Accept(); err == nil {
					if sys.TZLIB {
						if useTransport, er := thrift.NewTZlibTransport(transport, sys.ZLV); er == nil {
							twork := &transportWork{useTransport, processor, handler, cliErrorHandler, true, remoteHost(transport)}
							go twork.work("")
						}
					} else {
						twork := &transportWork{transport, processor, handler, cliErrorHandler, true, remoteHost(transport)}
						go twork.work("")
					}
				}
			}
		}
	}
	if !this.isClose && err != nil {
		sys.FmtLog("server init failed:", err)
		os.Exit(0)
	}
	return
}

func (this *tlClientServer) close() {
	defer errRecover()
	this.isClose = true
	this.serverTransport.Close()
}

type transportWork struct {
	transport thrift.TTransport
	processor thrift.TProcessor
	Handler   func(tc *tlContext)
	cliError  func(tc *tlContext)
	isServer  bool
	host      string
}

func (this *transportWork) work(addr string) (err error) {
	defer errRecover()
	tlcontext := newTlContext(this.transport)
	tlcontext.RemoteAddr = addr
	tlcontext.isServer = this.isServer
	tlcontext.RemoteHost = this.host
	go this.Handler(tlcontext)
	defaultCtx := context.WithValue(context.Background(), tlContextCtx, tlcontext)
	tlcontext.defaultCtx = defaultCtx
	tlcontext.cancleChan = make(chan byte)
	defer tlcontext.Close()
	defer this.cliError(tlcontext)
	err = this.processRequests(tlcontext, this.transport, this.processor)
	return
}

func (this *transportWork) processRequests(tlcontext *tlContext, client thrift.TTransport, processor thrift.TProcessor) (err error) {
	inputProtocol := _protocolFactory.GetProtocol(client)
	for {
		select {
		case <-tlcontext.cancleChan:
			goto END
		default:
			var ok bool
			ok, err = processor.Process(tlcontext.defaultCtx, inputProtocol, inputProtocol)
			if errors.Is(err, thrift.ErrAbandonRequest) {
				err = client.Close()
				goto END
			}
			if errors.As(err, new(thrift.TTransportException)) && err != nil {
				goto END
			}
			if !ok {
				goto END
			}
		}

	}
END:
	return
}

/**************************************************************************/

func transport2Client(tt thrift.TTransport) *ItnetClient {
	protocol := _protocolFactory.GetProtocol(tt)
	return NewItnetClientProtocol(tt, protocol, protocol)
}

/**************************************************************************/

type tlServerClient struct {
	transport thrift.TTransport
}

var tlserverclient = &tlServerClient{}

func (this *tlServerClient) server(addr string, itnet Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext), async bool) (err error) {
	defer errRecover()
	clientLinkCache.Put(addr, 0)
	defer clientLinkCache.Del(addr)
	transport := thrift.NewTSocketConf(addr, &thrift.TConfiguration{ConnectTimeout: sys.ConnectTimeout, SocketTimeout: sys.SocketTimeout})
	if sys.TZLIB {
		this.transport, err = thrift.NewTZlibTransport(transport, sys.ZLV)
	} else {
		this.transport = transport
	}
	if err = this.transport.Open(); err == nil {
		twork := &transportWork{this.transport, NewItnetProcessor(itnet), handler, cliErrorHandler, false, remoteHost(transport)}
		if async {
			go twork.work(addr)
		} else {
			twork.work(addr)
		}
	} else {
		logger.Error("tlServerClient to [", addr, "] Error:", err)
	}
	return
}

func (this *tlServerClient) Close() (err error) {
	defer errRecover()
	return this.transport.Close()
}

/////////////////////////////////////////////////

type ItnetWrite struct {
	mux  *sync.Mutex
	conn Itnet
}

func (this *ItnetWrite) Ping(ctx context.Context, pingstr int64) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Ping(ctx, pingstr)
}
func (this *ItnetWrite) Pong(ctx context.Context, pongBs []byte) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Pong(ctx, pongBs)
}

// Parameters:
//   - AuthKey
func (this *ItnetWrite) Auth(ctx context.Context, authKey []byte) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Auth(ctx, authKey)
}

// Parameters:
//   - AuthKey
func (this *ItnetWrite) Auth2(ctx context.Context, authKey []byte) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Auth2(ctx, authKey)
}

func (this *ItnetWrite) PonMerge(ctx context.Context, pblist []*PonBean, id int64) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.PonMerge(ctx, pblist, id)
}

// Parameters:
//   - Pb
func (this *ItnetWrite) Pon(ctx context.Context, ponBeanBytes []byte, id int64) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Pon(ctx, ponBeanBytes, id)
}

// Parameters:
//   - Pb
func (this *ItnetWrite) Pon2(ctx context.Context, pb *PonBean, id int64) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Pon2(ctx, pb, id)
}

// Parameters:
//   - Pb
func (this *ItnetWrite) Pon3(ctx context.Context, ponBeanBytes []byte, id int64, ack bool) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Pon3(ctx, ponBeanBytes, id, ack)
}

// Parameters:
//   - Pretimenano
//   - Timenano
//   - Num
//   - Ir
func (this *ItnetWrite) Time(ctx context.Context, pretimenano int64, timenano int64, num int16, ir bool) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.Time(ctx, pretimenano, timenano, num, ir)
}

// Parameters:
//   - Node
//   - Ir
func (this *ItnetWrite) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.SyncNode(ctx, node, ir)
}

func (this *ItnetWrite) SyncTx(ctx context.Context, syncId int64, result int8) (_err error) {
	if syncId > 0 {
		defer errRecover()
		this.mux.Lock()
		defer this.mux.Unlock()
		_err = this.conn.SyncTx(ctx, syncId, result)
	}
	return
}

func (this *ItnetWrite) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.SyncTxMerge(ctx, syncList)
}

// Parameters:
//   - SyncId
//   - Txid
func (this *ItnetWrite) CommitTx(ctx context.Context, syncId int64, txid int64) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.CommitTx(ctx, syncId, txid)
}

// Parameters:
//   - SyncId
//   - Commit
func (this *ItnetWrite) CommitTx2(ctx context.Context, syncId int64, txid int64, commit bool) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.CommitTx2(ctx, syncId, txid, commit)
}

func (this *ItnetWrite) PubMq(ctx context.Context, syncId int64, mqType int8, bs []byte) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.PubMq(ctx, syncId, mqType, bs)
}

func (this *ItnetWrite) PullData(ctx context.Context, syncId int64, ldb *LogDataBean) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.PullData(ctx, syncId, ldb)
}

func (this *ItnetWrite) ReInit(ctx context.Context, syncId int64, sbean *SysBean) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.ReInit(ctx, syncId, sbean)
}

func (this *ItnetWrite) ProxyCall(ctx context.Context, syncId int64, paramData []byte, pType int8, ctype int8) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.ProxyCall(ctx, syncId, paramData, pType, ctype)
}

func (this *ItnetWrite) BroadToken(ctx context.Context, syncId int64, tt *TokenTrans, ack bool) (_err error) {
	defer errRecover()
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.conn.BroadToken(ctx, syncId, tt, ack)
}
