// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package tc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/keystore"
	"github.com/donnie4w/tldb/level1"
	. "github.com/donnie4w/tldb/level2"
	"github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/tlmq"
	"github.com/donnie4w/tldb/util"
	"github.com/donnie4w/tlnet"
)

func init() {
	go wsOvertime()
}

var logger = log.Logger

type mqService struct {
	isClose bool
	tlnetMq *tlnet.Tlnet
}

var mqservice = &mqService{false, tlnet.NewTlnet()}

func (this *mqService) Serve() (err error) {
	if strings.TrimSpace(sys.MQADDR) != "" {
		err = this._serve(strings.TrimSpace(sys.MQADDR), sys.MQTLS, sys.MQCRT, sys.MQKEY)
	} else{
		sys.FmtLog("no mq service")
	}
	return
}

func (this *mqService) Close() (err error) {
	defer util.Recovr()
	if strings.TrimSpace(sys.MQADDR) != "" {
		this.isClose = true
		err = this.tlnetMq.Close()
	}
	return
}

func (this *mqService) _serve( addr string, TLS bool, serverCrt, serverKey string) (err error) {
	sys.MQADDR = addr
	this.tlnetMq.Handle("/mq2", mq2handler)
	this.tlnetMq.HandleWebSocketBindConfig("/mq", mqHandler, mqConfig())
	if TLS {
		sys.FmtLog(fmt.Sprint("Tldb Mq start tls[", sys.MQADDR, "/mq]"))
		if util.IsFileExist(serverCrt) && util.IsFileExist(serverKey) {
			if err = this.tlnetMq.HttpStartTLS(addr, serverCrt, serverKey); err != nil {
				err = errors.New("Tldb Mq tls start failed:" + err.Error())
			}
		} else {
			if err = this.tlnetMq.HttpStartTlsBytes(addr, []byte(keystore.ServerCrt), []byte(keystore.ServerKey)); err != nil {
				err = errors.New("Tldb Mq tls by bytes start failed:" + err.Error())
			}
		}
	}
	if !this.isClose {
		sys.FmtLog(fmt.Sprint("Tldb Mq start[", sys.MQADDR, "/mq]"))
		if err = this.tlnetMq.HttpStart(addr); err != nil {
			err = errors.New("Tldb Mq service init failed:" + err.Error())
		}
	}

	if !this.isClose && err != nil {
		sys.FmtLog(err.Error())
		os.Exit(0)
	}
	return
}

func mqConfig() *tlnet.WebsocketConfig {
	wc := &tlnet.WebsocketConfig{}
	wc.Origin = sys.WSORIGIN
	wc.OnError = func(self *tlnet.Websocket) {
		MqWare.DelConn(self.Id)
		WsWare.Del(self.Id)
	}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		tempOverMap.Put(hc.WS, time.Now().Unix())
	}
	return wc
}

func mq2handler(hc *tlnet.HttpContext) {
	defer util.Recovr()
	auth, _ := hc.GetCookie("auth")
	if hc.ReqInfo.Header.Get("Origin") != sys.WSORIGIN || !mqAuth(auth) {
		hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_NOPASS)...))
		return
	}
	var buf bytes.Buffer
	io.Copy(&buf, hc.Request().Body)
	bs := buf.Bytes()
	switch bs[0] {
	case MQ_PULLBYTE:
		if mb, err := pullByte(bs[1:]); err == nil {
			me := MQEncode(mb)
			buf := NewBuffer()
			buf.WriteByte(MQ_PULLBYTE)
			buf.Write(me)
			hc.ResponseBytes(0, buf.Bytes())
		} else {
			hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_PULLBYTE)...))
		}
	case MQ_PULLJSON:
		if mb, err := pullJson(bs[1:]); err == nil {
			je := JEncode(mb)
			buf := NewBuffer()
			buf.WriteByte(MQ_PULLJSON)
			buf.Write(je)
			hc.ResponseBytes(0, buf.Bytes())
		} else {
			hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_PULLJSON)...))
		}
	case MQ_CURRENTID:
		topic := string(bs[1:])
		if seq, err := Level2.SelectId(0, &TableStub{Tablename: key.Topic(topic)}); err == nil {
			buf := NewBuffer()
			buf.WriteByte(MQ_CURRENTID)
			buf.Write(util.Int64ToBytes(seq))
			hc.ResponseBytes(0, buf.Bytes())
		} else {
			hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_CURRENTID)...))
		}
	case MQ_LOCK:
		switch bs[1] {
		case 1:
			strId, extime, str := util.BytesToInt64(bs[2:10]), util.BytesToInt32(bs[10:14]), string(bs[14:])
			if c := tokenKeeper.new(strId, hc, str, extime, false); c != nil {
				<-c
			}
		case 2:
			strId, extime, str := util.BytesToInt64(bs[2:10]), util.BytesToInt32(bs[10:14]), string(bs[14:])
			if token, ok := sys.TryLock(str); ok {
				tokenKeeper.new(strId, hc, str, extime, true)
				reqToken(strId, token)
			}
		case 3:
			strId := util.BytesToInt64(bs[2:10])
			if tb, ok := tokenKeeper.get(strId); ok {
				token := string(bs[10:26])
				if sys.UnLock(tb.str, token) {
					tokenKeeper.del(strId)
				}
			}
			hc.ResponseBytes(0, []byte{1})
		}
	}
}

var tempCliMap = NewLimitMap[int64, int](1 << 10)
var tempOverMap = NewMap[*tlnet.Websocket, int64]()

func mqHandler(hc *tlnet.HttpContext) {
	defer util.Recovr()
	bs := hc.WS.Read()
	if bs == nil || (bs[0] != MQ_AUTH && !isAuth(hc.WS)) {
		nopass := false
		if bs == nil {
			nopass = true
		} else if bs[0] == MQ_PING {
			if v, ok := tempCliMap.Get(hc.WS.Id); ok {
				if v > 3 {
					nopass = true
				} else {
					v++
					tempCliMap.Put(hc.WS.Id, v)
				}
			} else {
				tempCliMap.Put(hc.WS.Id, 1)
			}
		}
		if nopass {
			logger.Error("nopass>>", hc.WS.Id)
			hc.WS.Close()
			return
		}
	}
	var err error
	switch bs[0] & 0x7f {
	case MQ_AUTH:
		if !WsWare.Has(hc.WS.Id) {
			if mqAuth(string(bs[9:])) {
				cliId := int64(0)
				if len(bs) >= 17 {
					cliId = util.BytesToInt64(bs[9:16])
				}
				err = AddConn(hc.WS, cliId)
				tempOverMap.Del(hc.WS)
			} else {
				sendErr(bs[1:9], hc.WS)
				hc.WS.Close()
				return
			}
		}
	case MQ_RECVACK:
		sec := int8(10)
		if len(bs) >= 10 {
			if s := int8(bs[9]); s > 0 {
				sec = s
			}
		}
		SetRecvAck(hc.WS.Id, sec)
	case MQ_PUBBYTE:
		err, _ = pubByte(bs[9:], hc.WS.Id, bs[0]|0x80 == 0x80)
	case MQ_PUBJSON:
		err, _ = pubJson(bs[9:], hc.WS.Id, bs[0]|0x80 == 0x80)
	case MQ_SUB:
		s := string(bs[9:])
		err = sub(s, 0, hc.WS)
		if bs[0]&0x80 == 0x80 {
			SetJsonOn(hc.WS.Id, s, true)
		} else {
			SetJsonOn(hc.WS.Id, s, false)
		}
	case MQ_SUBCANCEL:
		err = sub(string(bs[9:]), 1, hc.WS)
	case MQ_PULLBYTE:
		var mb *MqBean
		if mb, err = pullByte(bs[9:]); err == nil {
			MqWare.PullByte(mb, hc.WS.Id)
		}
	case MQ_PULLJSON:
		var mb *JMqBean
		if mb, err = pullJson(bs[9:]); err == nil {
			err = MqWare.PullJson(mb, hc.WS.Id)
		}
	case MQ_PING:
		MqWare.Ping(bytes.NewBuffer(bs), hc.WS.Id)
	case MQ_PUBMEM:
		err = pubMem(bs[9:], hc.WS.Id, bs[0]|0x80 == 0x80)
	case MQ_MERGE:
		size := int8(60)
		if len(bs) >= 10 {
			if s := int8(bs[9]); s > 0 {
				size = s
			}
		}
		SetMergeOn(hc.WS.Id, size)
	case MQ_ZLIB:
		on := false
		if bs[9] == 1 {
			on = true
		}
		SetZlib(hc.WS.Id, on)
	case MQ_CURRENTID:
		var seq int64
		if seq, err = Level2.SelectId(0, &TableStub{Tablename: key.Topic(string(bs[9:]))}); err == nil {
			MqWare.PullId(seq)
		}
	case MQ_ACK:
		if IsRecvAckOn(hc.WS.Id) {
			ackId := util.BytesToInt64(bs[1:9])
			DelAckId(hc.WS.Id, ackId)
		}
	}
	if bs[0] != MQ_PING && bs[0] != MQ_ACK {
		if err == nil {
			mqAck(bs[1:9], hc.WS)
		} else {
			sendErr(bs[1:9], hc.WS)
		}
	}
}

func isAuth(ws *tlnet.Websocket) bool {
	return WsWare.Has(ws.Id)
}

func mqAck(bs []byte, ws *tlnet.Websocket) {
	var buf bytes.Buffer
	buf.WriteByte(MQ_ACK)
	buf.Write(bs)
	MqWare.Ack(&buf, ws.Id)
}

/******************************************************/
func pubByte(bs []byte, wsId int64, self bool) (err error, seq int64) {
	mb := &MqBean{}
	if _, err = MQDecode(bs, mb); err == nil {
		if seq, err = saveTopic(mb.Topic, mb.Msg); err == nil && seq > 0 {
			mb.ID = seq
			MqWare.PubByte(mb, wsId, self)
			level1.BroadcastMQ(int8(MQ_PUBBYTE), MQEncode(mb))
		}
	}
	return
}

func pubJson(bs []byte, wsId int64, self bool) (err error, seq int64) {
	var mb *JMqBean
	if mb, err = JDecode(bs); err == nil {
		if seq, err = saveTopic(mb.Topic, []byte(mb.Msg)); err == nil && seq > 0 {
			mb.Id = seq
			MqWare.PubJson(mb, wsId, self)
			level1.BroadcastMQ(int8(MQ_PUBJSON), JEncode(mb))
		}
	}
	return
}

func pullByte(bs []byte) (mb *MqBean, err error) {
	mb = &MqBean{}
	if _, err = MQDecode(bs, mb); err == nil {
		var _bs []byte
		if _bs, err = getMsgFromId(mb.Topic, mb.ID); err == nil {
			mb.Msg = _bs
		}
	}
	return
}

func pullJson(bs []byte) (mb *JMqBean, err error) {
	if mb, err = JDecode(bs); err == nil {
		var _bs []byte
		if _bs, err = getMsgFromId(mb.Topic, mb.Id); err == nil {
			mb.Msg = string(_bs)
		}
	}
	return
}

func pullId(bs []byte) (mb *JMqBean, err error) {
	if mb, err = JDecode(bs); err == nil {
		var _bs []byte
		if _bs, err = getMsgFromId(mb.Topic, mb.Id); err == nil {
			mb.Msg = string(_bs)
		}
	}
	return
}

func pubMem(bs []byte, wsId int64, self bool) (err error) {
	var mb *JMqBean
	if mb, err = JDecode(bs); err == nil {
		MqWare.PubMem(mb, wsId, self)
		level1.BroadcastMQ(int8(MQ_PUBMEM), JEncode(mb))
	}
	return
}

func saveTopic(topic string, msg []byte) (id int64, err error) {
	if _, ok := TldbMq.GetTable(topic); !ok {
		Level2.CreateTableMq(&TableStub{Tablename: topic, Field: map[string][]byte{"msg": msg}})
	}
	id, err = Level2.InsertMq(0, &TableStub{Tablename: key.Topic(topic), Field: map[string][]byte{"msg": msg}})
	return
}

func getMsgFromId(topic string, id int64) (bs []byte, err error) {
	if _r, er := Level2.SelectById(0, &TableStub{Tablename: key.Topic(topic), ID: id}); er == nil {
		bs = _r.Field["msg"]
	} else {
		err = er
	}
	return
}

func sub(topic string, t byte, ws *tlnet.Websocket) (err error) {
	if t == 0 {
		MqWare.AddConn(topic, ws)
	} else {
		MqWare.DelTopicWithID(topic, ws.Id)
	}
	return
}

func mqAuth(msg string) (isAuth bool) {
	if msg != "" {
		if ss := strings.Split(string(msg), "="); len(ss) == 2 {
			if _r, ok := keystore.StoreAdmin.GetMq(ss[0]); ok {
				isAuth = _r.Pwd == util.MD5(ss[1])
			}
		}
	}
	return
}

func sendErr(errId []byte, ws *tlnet.Websocket) {
	ws.Send(append([]byte{MQ_ERROR}, errId...))
}

func wsOvertime() {
	tk := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tk.C:
			tempOverMap.Range(func(k *tlnet.Websocket, v int64) bool {
				defer recover()
				if v+60 < time.Now().Unix() {
					logger.Error("overtime>>", k.Id)
					k.Close()
					tempOverMap.Del(k)
				}
				return true
			})
		}
	}
}
