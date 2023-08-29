/**
 * Copyright 2023 tldb Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package tlmq

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/level1"
	"github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
	"github.com/donnie4w/tlnet"
)

var logger = log.Logger

func init() {
	level1.ClusPub = MqWare.ClusPub
}

var MqWare = NewMqEngine()

var WsWare = NewMapL[int64, *WsSock]()
var CliWare = NewMapL[int64, int64]()

func AddConn(ws *tlnet.Websocket, cliId int64) {
	WsWare.Put(ws.Id, NewWsSock(ws, cliId))
	if cliId > 0 {
		CliWare.Put(cliId, ws.Id)
	}
}

func DelAckId(wsid, ackid int64) {
	if wss, ok := WsWare.Get(wsid); ok {
		wss.DelAckId(ackid)
	}
}

func SetRecvAck(id int64, sec int8) {
	if wss, ok := WsWare.Get(id); ok {
		wss.SetRecvAckOn(sec)
	}
}

func SetMergeOn(id int64, size int8) {
	if wss, ok := WsWare.Get(id); ok {
		wss.SetMergeOn(size)
	}
}

func SetZlib(id int64, on bool) {
	if wss, ok := WsWare.Get(id); ok {
		wss.SetZlib(on)
	}
}

// ///////////////////////////////////////////
func NewMqEngine() MqEg {
	return &mqWare{NewMapL[string, *MapL[int64, byte]](), NewMapL[int64, *MapL[string, int8]](), &sync.Mutex{}}
}

type mqWare struct {
	subTopicMap *MapL[string, *MapL[int64, byte]]
	wsTop       *MapL[int64, *MapL[string, int8]]
	mux         *sync.Mutex
}

func (this *mqWare) AddConn(topic string, ws *tlnet.Websocket) {
	defer this.mux.Unlock()
	this.mux.Lock()
	var tpm *MapL[int64, byte]
	var ok bool
	if tpm, ok = this.subTopicMap.Get(topic); !ok {
		tpm = NewMapL[int64, byte]()
		this.subTopicMap.Put(topic, tpm)
	}
	tpm.Put(ws.Id, 0)
	///////////////////////////////
	var wsm *MapL[string, int8]
	if wsm, ok = this.wsTop.Get(ws.Id); !ok {
		wsm = NewMapL[string, int8]()
		this.wsTop.Put(ws.Id, wsm)
	}
	wsm.Put(topic, 0)
}

// count of sub topic
func (this *mqWare) SubCount(topic string) (_r int64) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if m, ok := this.subTopicMap.Get(topic); ok {
		_r = m.Len()
	}
	return
}

func (this *mqWare) DelTopic(topic string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.subTopicMap.Del(topic)
	this.wsTop.Range(func(_ int64, v *MapL[string, int8]) bool {
		v.Del(topic)
		return true
	})
}

func (this *mqWare) DelTopicWithID(topic string, id int64) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if wm, ok := this.subTopicMap.Get(topic); ok {
		wm.Del(id)
	}
	if ss, ok := this.wsTop.Get(id); ok {
		ss.Del(topic)
	}
}

func (this *mqWare) DelConn(id int64) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.wsTop.Del(id)
	this.subTopicMap.Range(func(k string, v *MapL[int64, byte]) bool {
		v.Del(id)
		if v.Len() == 0 {
			this.subTopicMap.Del(k)
		}
		return true
	})
}

func (this *mqWare) ClusPub(mqType int8, bs []byte) (err error) {
	switch byte(mqType) {
	case MQ_PUBBYTE:
		if mb, err := MQDecode(bs, &MqBean{}); err == nil {
			err = this.PubByte(mb)
		}
	case MQ_PUBJSON:
		if mb, err := JDecode(bs); err == nil {
			err = this.PubJson(mb)
		}
	case MQ_PUBMEM:
		if mb, err := JDecode(bs); err == nil {
			this.PubMem(mb)
		}
	}
	return
}

func (this *mqWare) PubByte(mb *MqBean) (err error) {
	bs := MQEncode(mb)
	buf := util.BufferPool.Get(1 + len(bs))
	buf.WriteByte(MQ_PUBBYTE)
	buf.Write(bs)
	if m, ok := this.subTopicMap.Get(mb.Topic); ok {
		m.Range(func(k int64, _ byte) bool {
			if wss, ok := WsWare.Get(k); ok {
				util.GoPool.Go(func() {
					wss.SendWithAck(buf)
				})
			}
			return true
		})
	}
	return
}

func (this *mqWare) PullByte(mb *MqBean, id int64) (err error) {
	bs := MQEncode(mb)
	buf := util.BufferPool.Get(1 + len(bs))
	buf.WriteByte(MQ_PULLBYTE)
	buf.Write(bs)
	if wss, ok := WsWare.Get(id); ok {
		util.GoPool.Go(func() {
			wss.SendWithAck(buf)
		})
	}
	return
}

func (this *mqWare) PubJson(mb *JMqBean) (err error) {
	bs := JEncode(mb)
	buf := util.BufferPool.Get(1 + len(bs))
	buf.WriteByte(MQ_PUBJSON)
	buf.Write(bs)
	if m, ok := this.subTopicMap.Get(mb.Topic); ok {
		m.Range(func(k int64, _ byte) bool {
			if wss, ok := WsWare.Get(k); ok {
				util.GoPool.Go(func() { wss.SendWithAck(buf) })
			}
			return true
		})
	}
	return
}

func (this *mqWare) PullJson(mb *JMqBean, id int64) (err error) {
	bs := JEncode(mb)
	buf := util.BufferPool.Get(1 + len(bs))
	buf.WriteByte(MQ_PULLJSON)
	buf.Write(bs)
	if wss, ok := WsWare.Get(id); ok {
		util.GoPool.Go(func() {
			wss.SendWithAck(buf)
		})
	}
	return
}

func (this *mqWare) PullId(id int64) (err error) {
	bs := util.Int64ToBytes(id)
	buf := util.BufferPool.Get(1 + len(bs))
	buf.WriteByte(MQ_CURRENTID)
	buf.Write(bs)
	if wss, ok := WsWare.Get(id); ok {
		util.GoPool.Go(func() {
			wss.SendWithAck(buf)
		})
	}
	return
}

func (this *mqWare) PubMem(mb *JMqBean) {
	bs := JEncode(mb)
	buf := util.BufferPool.Get(1 + len(bs))
	buf.WriteByte(MQ_PUBMEM)
	buf.Write(bs)
	if m, ok := this.subTopicMap.Get(mb.Topic); ok {
		m.Range(func(k int64, _ byte) bool {
			if wss, ok := WsWare.Get(k); ok {
				util.GoPool.Go(func() { wss.Send(buf) })
			}
			return true
		})
	}
}

func (this *mqWare) Ping(buf *bytes.Buffer, id int64) {
	if wss, ok := WsWare.Get(id); ok {
		util.GoPool.Go(func() {
			wss.Send(buf)
		})
	}
}

func (this *mqWare) Ack(buf *bytes.Buffer, id int64) {
	if wss, ok := WsWare.Get(id); ok {
		util.GoPool.Go(func() {
			wss.Send(buf)
		})
	}
}

/*********************************************************/

func NewWsSock(ws *tlnet.Websocket, cliId int64) (_r *WsSock) {
	_r = &WsSock{}
	_r.cliId = cliId
	_r.connId = ws.Id
	_r.ws = ws
	_r.mux = &sync.Mutex{}
	_r.ackMap = NewLinkedMap[int64, []any]()
	_r.mergeChan = make(chan *bytes.Buffer, 1<<15)
	_r.mergemux = &sync.Mutex{}
	// _r.merge = &mergeMq[*bytes.Buffer, *bytes.Buffer]{mmap: NewMapL[int64, *bytes.Buffer](), mux: &sync.Mutex{}}
	return
}

type WsSock struct {
	cliId     int64
	connId    int64
	ws        *tlnet.Websocket
	recvAckOn bool
	recvSec   time.Duration
	mergeSize int8
	zlibOn    bool
	mux       *sync.Mutex
	ackMap    *LinkedMap[int64, []any]
	// merge      Merge[*bytes.Buffer, *bytes.Buffer]

	mergemux   *sync.Mutex
	mergeChan  chan *bytes.Buffer
	mergeCount int64
}

func (this *WsSock) SetRecvAckOn(sec int8) {
	this.recvSec = time.Duration(sec) * time.Second
	defer this.mux.Unlock()
	this.mux.Lock()
	if !this.recvAckOn {
		this.recvAckOn = true
		go this.AckTimer()
	}
}

func (this *WsSock) SetMergeOn(size int8) {
	this.mergeSize = size
	// this.merge.MergeSize(size)
}

func (this *WsSock) SetZlib(zlib bool) {
	this.zlibOn = zlib
	// this.merge.SetZlib(zlib)
}

func (this *WsSock) AckTimer() {
	ticker := time.NewTicker(3 * time.Second)
	for this.ws.Error == nil {
		select {
		case <-ticker.C:
			func() {
				defer util.ErrRecovr()
				this.ackMap.BackForEach(func(_ int64, bs []interface{}) bool {
					t := bs[0].(int64)
					if t+int64(this.recvSec) < time.Now().UnixNano() {
						if err := this._send(bs[1].(*bytes.Buffer)); err == nil {
							return true
						}
					}
					return false
				})
			}()
		}
	}
}

func (this *WsSock) DelAckId(id int64) {
	this.ackMap.Del(id)
}

func (this *WsSock) _send(buf *bytes.Buffer) (err error) {
	defer util.ErrRecovr()
	return this.ws.Send(buf.Bytes())
}

func (this *WsSock) MergeSend(buf *bytes.Buffer) {
	// this.merge.Add(buf)
	// this.merge.CallBack(func() bool { return this.ws.Error != nil }, this._sendWithAck)
	this.mergeChan <- buf
	atomic.AddInt64(&this.mergeCount, 1)
	if this.mergemux.TryLock() {
		go wssMerge(this)
	}
}

func (this *WsSock) Send(buf *bytes.Buffer) (err error) {
	if this.mergeSize > 0 {
		this.MergeSend(buf)
	} else {
		return this._send(buf)
	}
	return
}

func (this *WsSock) SendWithAck(buf *bytes.Buffer) (err error) {
	if this.mergeSize > 0 {
		this.MergeSend(buf)
	} else {
		return this._sendWithAck(buf)
	}
	return
}

func (this *WsSock) _sendWithAck(buf *bytes.Buffer) (err error) {
	if this.recvAckOn {
		util.GoPool.Go(func() {
			if this.ackMap.Len() < 1<<18 {
				this.ackMap.Put(int64(util.CRC32(buf.Bytes())), []any{time.Now().UnixNano(), buf})
			}
		})
	}
	err = this._send(buf)
	return
}

func IsRecvAckOn(id int64) (_r bool) {
	if wss, ok := WsWare.Get(id); ok {
		_r = wss.recvAckOn
	}
	return
}

/*****************************************/
type mergeMq[T *bytes.Buffer, V *bytes.Buffer] struct {
	mmap      *MapL[int64, T]
	seq       int64
	cursor    int64
	mux       *sync.Mutex
	mergeSize int8
	zlib      bool
}

func (this *mergeMq[T, V]) Add(buf T) {
	this.mmap.Put(atomic.AddInt64(&this.seq, 1), buf)
}

func (this *mergeMq[T, V]) Del(endId int64) {
	this.mmap.Del(endId)
}

func (this *mergeMq[T, V]) Get(id int64) (_r T, ok bool) {
	_r, ok = this.mmap.Get(id)
	return
}

func (this *mergeMq[T, V]) Len() int64 {
	return this.mmap.Len()
}

func (this *mergeMq[T, V]) SetZlib(zlib bool) {
	this.zlib = zlib
}

func (this *mergeMq[T, V]) MergeSize(size int8) {
	this.mergeSize = size
}

func (this *mergeMq[T, V]) CallBack(cancelfunc func() bool, f func(bs V) error) (err error) {
	if this.mux.TryLock() {
		defer this.mux.Unlock()
	} else {
		return
	}
	for this.Len() > 0 && this.seq > this.cursor && !cancelfunc() {
		if sys.MERGETIME > 0 {
			t := time.Duration(sys.MERGETIME) * time.Microsecond
			if t > time.Second {
				t = time.Second
			}
			<-time.After(t)
		}
		if buf, err := this.getMergeBean(); err == nil {
			f(buf)
		}
	}
	return
}

func (this *mergeMq[T, V]) getMergeBean() (_r *bytes.Buffer, err error) {
	defer util.ErrRecovr()
	size := int(this.mergeSize) * 1 << 20
	mb := &MergeBean{BeanList: make([][]byte, 0)}
	start, end := this.cursor+1, this.cursor+1
	for i := start; size > 0 && i <= this.seq; i++ {
		var bb *bytes.Buffer
		var ok bool
		for {
			if bb, ok = this.Get(i); ok {
				break
			} else {
				<-time.After(1 * time.Microsecond)
			}
		}
		size = size - bb.Len()
		mb.BeanList = append(mb.BeanList, bb.Bytes())
		end = i
	}
	var _bs []byte
	if this.zlib {
		_bs, err = util.ZlibCz(util.TEncode(mb))
	} else {
		_bs = util.TEncode(mb)
	}
	if err == nil {
		for i := start; i <= end; i++ {
			this.Del(i)
		}
		this.cursor = end
		_r = bytes.NewBuffer([]byte{})
		_r.WriteByte(MQ_MERGE)
		_r.Write(_bs)
	}
	return
}

func wssMerge(wss *WsSock) {
	defer recover()
	defer wss.mergemux.Unlock()
	_wssMerge(wss)
}

func _wssMerge(wss *WsSock) {
	mb := &MergeBean{BeanList: make([][]byte, 0)}
	size := int(wss.mergeSize) * 1 << 20
	for bb := range wss.mergeChan {
		mb.BeanList = append(mb.BeanList, bb.Bytes())
		if atomic.AddInt64(&wss.mergeCount, -1) <= 0 {
			break
		}
		if size = size - bb.Len(); size <= 0 {
			break
		}
	}
	var _bs []byte
	if wss.zlibOn {
		_bs, _ = util.ZlibCz(util.TEncode(mb))
	} else {
		_bs = util.TEncode(mb)
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(MQ_MERGE)
	buf.Write(_bs)

	wss._sendWithAck(buf)
	if wss.mergeCount > 0 {
		_wssMerge(wss)
	}
}
