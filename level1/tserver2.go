// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"context"
	"sync"

	. "github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/tldb/log"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
	. "github.com/donnie4w/tsf/packet"
	. "github.com/donnie4w/tsf/server"
)

type tfsClientServer struct {
	isClose         bool
	serverTransport *TServerSocket
}

var tfsclientserver = &tfsClientServer{}

func (this *tfsClientServer) server(wg *sync.WaitGroup, _addr string, processor Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext)) (err error) {
	defer wg.Done()
	if this.serverTransport, err = NewTServerSocketTimeout(_addr, sys.ConnectTimeout); err == nil {
		if err = this.serverTransport.Listen(); err == nil {
			sys.ADDR = _addr
			for {
				if transport, err := this.serverTransport.Accept(); err == nil {
					transport.SetSnappyMergeData(true)
					go func() {
						twork := &sockWork{transport: transport, processor: processor, handler: handler, cliError: cliErrorHandler, isServer: true, host: remoteHost2(transport)}
						twork.work("")
						defer twork.final()
						ProcessMerge(transport, func(socket TsfSocket, pkt *Packet) error {
							return twork.processRequests(pkt, processor)
						})
					}()
				}
			}
		}
	}
	if !this.isClose && err != nil {
		log.LoggerSys.Error("server init failed:", err)
		panic("server init failed:" + err.Error())
	}
	return
}

func (this *tfsClientServer) close() {
	defer myRecovr()
	this.isClose = true
	this.serverTransport.Close()
}

type sockWork struct {
	transport *TSocket
	processor Itnet
	handler   func(tc *tlContext)
	cliError  func(tc *tlContext)
	isServer  bool
	host      string
	tlcontext *tlContext
}

func (this *sockWork) work(addr string) (err error) {
	defer myRecovr()
	tlcontext := newTlContext2(this.transport)
	tlcontext.RemoteAddr = addr
	tlcontext.isServer = this.isServer
	tlcontext.RemoteHost = this.host
	go this.handler(tlcontext)
	tlcontext.defaultCtx = context.WithValue(context.Background(), tlContextCtx, tlcontext)
	tlcontext.cancleChan = make(chan byte)

	this.tlcontext = tlcontext
	return
}

func (this *sockWork) final() {
	defer this.tlcontext.Close()
	defer this.cliError(this.tlcontext)
}

func (this *sockWork) processRequests(packet *Packet, processor Itnet) (err error) {
	defer myRecovr()
	bs := packet.ToBytes()
	ln := int(int8(bs[0]))
	method := string(bs[1 : ln+1])
	ctx := this.tlcontext.defaultCtx
	bs = bs[1+ln:]
	switch method {
	case string(_Ping):
		processor.Ping(ctx, util.BytesToInt64(bs))
	case string(_Pong):
		processor.Pong(ctx, bs)
	case string(_Auth):
		processor.Auth(ctx, bs)
	case string(_Auth2):
		processor.Auth2(ctx, bs)
	case string(_Pon):
		id := util.BytesToInt64(bs[:8])
		PonBeanBytes := bs[8:]
		processor.Pon(ctx, PonBeanBytes, id)
	case string(_Pon2):
		id := util.BytesToInt64(bs[:8])
		pb := &PonBean{}
		if err := util.Decode(bs[8:], pb); err == nil {
			processor.Pon2(ctx, pb, id)
		} else {
			logger.Error(err)
		}
	case string(_Pon3):
		id := util.BytesToInt64(bs[:8])
		ack := byteToBool(bs[8:9][0])
		pb := bs[9:]
		processor.Pon3(ctx, pb, id, ack)
	case string(_Time):
		pretimenano := util.BytesToInt64(bs[:8])
		timenano := util.BytesToInt64(bs[8 : 8+8])
		num := util.BytesToInt16(bs[8+8 : 8+8+2])
		ir := byteToBool(bs[8+8+2:][0])
		processor.Time(ctx, pretimenano, timenano, num, ir)
	case string(_SyncNode):
		ir := byteToBool(bs[0:1][0])
		node := &Node{}
		if err := util.Decode(bs[1:], node); err == nil {
			processor.SyncNode(ctx, node, ir)
		} else {
			logger.Error(err)
		}
	case string(_SyncTx):
		syncId := util.BytesToInt64(bs[:8])
		result := int8(bs[8:][0])
		processor.SyncTx(ctx, syncId, result)
	case string(_SyncTxMerge):
		syncList := bytesToMap(bs)
		processor.SyncTxMerge(ctx, syncList)
	case string(_PubMq):
		syncId := util.BytesToInt64(bs[:8])
		mqType := int8(bs[8:9][0])
		bs := bs[9:]
		processor.PubMq(ctx, syncId, mqType, bs)
	case string(_ReInit):
		syncId := util.BytesToInt64(bs[:8])
		sbean := &SysBean{}
		util.Decode(bs[8:], sbean)
		processor.ReInit(ctx, syncId, sbean)
	case string(_PullData):
		syncId := util.BytesToInt64(bs[:8])
		ldb := &LogDataBean{}
		if err := util.Decode(bs[8:], ldb); err == nil {
			processor.PullData(ctx, syncId, ldb)
		} else {
			logger.Error(err)
		}
	case string(_ProxyCall):
		syncId := util.BytesToInt64(bs[:8])
		pType := int8(bs[8:9][0])
		cType := int8(bs[9:10][0])
		paramData := bs[10:]
		processor.ProxyCall(ctx, syncId, paramData, pType, cType)
	case string(_BroadToken):
		syncId := util.BytesToInt64(bs[:8])
		ack := false
		if bs[8] == 1 {
			ack = true
		}
		tt, _ := util.TDecode(bs[9:], &TokenTrans{})
		processor.BroadToken(ctx, syncId, tt, ack)
	}
	return
}

type tfsServerClient struct {
	transport *TSocket
}

var tfsserverclient = &tfsServerClient{}

func (this *tfsServerClient) server(addr string, processor Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext), async bool) (err error) {
	defer myRecovr()
	clientLinkCache.Put(addr, 0)
	defer clientLinkCache.Del(addr)
	this.transport = NewTSocketConf(addr, &TConfiguration{ConnectTimeout: sys.ConnectTimeout})
	if err = this.transport.Open(); err == nil {
		this.transport.SetSnappyMergeData(true)
		twork := &sockWork{transport: this.transport, processor: processor, handler: handler, cliError: cliErrorHandler, isServer: false, host: remoteHost2(this.transport)}
		twork.work(addr)
		if async {
			go func() {
				defer twork.final()
				ProcessMerge(this.transport, func(socket TsfSocket, pkt *Packet) error {
					return twork.processRequests(pkt, processor)
				})
			}()
		} else {
			defer twork.final()
			ProcessMerge(this.transport, func(socket TsfSocket, pkt *Packet) error {
				return twork.processRequests(pkt, processor)
			})
		}
	} else {
		logger.Error("serverClient to [", addr, "] Error:", err)
	}
	return
}

func (this *tfsServerClient) Close() (err error) {
	defer myRecovr()
	return this.transport.Close()
}

func process(packet *Packet) {
	packet.ToBytes()
}

type ItnetImpl struct {
	socket *TSocket
}

func method(name []byte) (buf *Buffer) {
	buf = NewBuffer()
	i := len(name)
	buf.WriteByte(byte(i))
	buf.Write(name)
	return
}

func (this *ItnetImpl) Ping(ctx context.Context, pingstr int64) (_err error) {
	buf := method(_Ping)
	buf.Write(util.Int64ToBytes(pingstr))

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) Pong(ctx context.Context, pongBs []byte) (_err error) {
	buf := method(_Pong)
	buf.Write(pongBs)
	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

// Parameters:
//   - AuthKey
func (this *ItnetImpl) Auth(ctx context.Context, authKey []byte) (_err error) {
	buf := method(_Auth)
	buf.Write(authKey)
	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

// Parameters:
//   - AuthKey
func (this *ItnetImpl) Auth2(ctx context.Context, authKey []byte) (_err error) {
	buf := method(_Auth2)
	buf.Write(authKey)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) PonMerge(ctx context.Context, pblist []*PonBean, id int64) (_err error) {
	return
}

// Parameters:
//   - Pb
func (this *ItnetImpl) Pon(ctx context.Context, ponBeanBytes []byte, id int64) (_err error) {
	buf := method(_Pon)
	buf.Write(util.Int64ToBytes(id))
	buf.Write(ponBeanBytes)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

// Parameters:
//   - Pb
func (this *ItnetImpl) Pon2(ctx context.Context, pb *PonBean, id int64) (_err error) {
	buf := method(_Pon2)

	buf.Write(util.Int64ToBytes(id))
	pbbytes, _ := util.Encode(pb)
	buf.Write(pbbytes)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

// Parameters:
//   - Pb
func (this *ItnetImpl) Pon3(ctx context.Context, ponBeanBytes []byte, id int64, ack bool) (_err error) {
	buf := method(_Pon3)

	buf.Write(util.Int64ToBytes(id))
	buf.WriteByte(boolToByte(ack))
	buf.Write(ponBeanBytes)
	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

// Parameters:
//   - Pretimenano
//   - Timenano
//   - Num
//   - Ir
func (this *ItnetImpl) Time(ctx context.Context, pretimenano int64, timenano int64, num int16, ir bool) (_err error) {
	buf := method(_Time)

	buf.Write(util.Int64ToBytes(pretimenano))
	buf.Write(util.Int64ToBytes(timenano))
	buf.Write(util.Int16ToBytes(num))
	buf.WriteByte(boolToByte(ir))
	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

// Parameters:
//   - Node
//   - Ir
func (this *ItnetImpl) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
	buf := method(_SyncNode)

	buf.WriteByte(boolToByte(ir))
	nodebs, _ := util.Encode(node)
	buf.Write(nodebs)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) SyncTx(ctx context.Context, syncId int64, result int8) (_err error) {
	buf := method(_SyncTx)

	buf.Write(util.Int64ToBytes(syncId))
	buf.WriteByte(byte(result))

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
	buf := method(_SyncTxMerge)

	syncListBuf := mapTobytes(syncList)
	buf.Write(syncListBuf.Bytes())

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func mapTobytes(syncList map[int64]int8) (buf *Buffer) {
	buf = NewBuffer()
	for k, v := range syncList {
		buf.Write(util.Int64ToBytes(k))
		buf.WriteByte(byte(v))
	}
	return
}

func bytesToMap(bs []byte) (syncList map[int64]int8) {
	syncList = make(map[int64]int8, 0)
	for i := 0; i < len(bs); i += 9 {
		k := util.BytesToInt64(bs[i : i+8])
		syncList[k] = int8(bs[i+8:][0])
	}
	return
}

// Parameters:
//   - SyncId
//   - Txid
func (this *ItnetImpl) CommitTx(ctx context.Context, syncId int64, txid int64) (_err error) {
	return
}

// Parameters:
//   - SyncId
//   - Commit
func (this *ItnetImpl) CommitTx2(ctx context.Context, syncId int64, txid int64, commit bool) (_err error) {
	return
}

func (this *ItnetImpl) PubMq(ctx context.Context, syncId int64, mqType int8, bs []byte) (_err error) {
	buf := method(_PubMq)

	buf.Write(util.Int64ToBytes(syncId))
	buf.WriteByte(byte(mqType))
	buf.Write(bs)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) PullData(ctx context.Context, syncId int64, ldb *LogDataBean) (_err error) {
	buf := method(_PullData)

	buf.Write(util.Int64ToBytes(syncId))
	ldbBys, _ := util.Encode(ldb)
	buf.Write(ldbBys)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) ReInit(ctx context.Context, syncId int64, sbean *SysBean) (_err error) {
	buf := method(_ReInit)
	buf.Write(util.Int64ToBytes(syncId))
	sbeanBys, _ := util.Encode(sbean)
	buf.Write(sbeanBys)

	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) ProxyCall(ctx context.Context, syncId int64, paramData []byte, pType int8, ctype int8) (_err error) {
	buf := method(_ProxyCall)

	buf.Write(util.Int64ToBytes(syncId))
	buf.WriteByte(byte(pType))
	buf.WriteByte(byte(ctype))

	buf.Write(paramData)
	go this.socket.WriteWithMerge(buf.Bytes())
	return
}

func (this *ItnetImpl) BroadToken(ctx context.Context, syncId int64, tt *TokenTrans, ack bool) (_err error) {
	buf := method(_BroadToken)
	buf.Write(util.Int64ToBytes(syncId))
	if ack {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	buf.Write(util.TEncode(tt))
	go this.socket.WriteWithMerge(buf.Bytes())
	return
}
