/**
 * Copyright 2023 tldb Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 * github.com/donnie4w/tldb
 */
package tc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/keystore"
	"github.com/donnie4w/tldb/level1"
	. "github.com/donnie4w/tldb/level2"
	"github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/tlmq"
	. "github.com/donnie4w/tldb/tlmq"
	"github.com/donnie4w/tldb/util"
	"github.com/donnie4w/tlnet"
)

var logger = log.Logger

type mqService struct {
	isClose bool
	tlnetMq *tlnet.Tlnet
}

var mqservice = &mqService{false, tlnet.NewTlnet()}

func (this *mqService) Serve(wg *sync.WaitGroup) (err error) {
	if strings.TrimSpace(sys.MQADDR) != "" {
		err = this._serve(wg, strings.TrimSpace(sys.MQADDR), sys.MQTLS, sys.MQCRT, sys.MQKEY)
	} else {
		wg.Done()
	}
	return
}

func (this *mqService) Close() (err error) {
	defer util.ErrRecovr()
	if strings.TrimSpace(sys.MQADDR) != "" {
		this.isClose = true
		err = this.tlnetMq.Close()
	}
	return
}

func (this *mqService) _serve(wg *sync.WaitGroup, addr string, TLS bool, serverCrt, serverKey string) (err error) {
	defer wg.Done()
	sys.MQADDR = addr
	this.tlnetMq.Handle("/mq2", mq2handler)
	this.tlnetMq.HandleWebSocketBindConfig("/mq", mqHandler, mqConfig())
	if TLS {
		log.LoggerSys.Info(sys.SysLog(fmt.Sprint("tlMq start tls[", sys.MQADDR, "/mq]")))
		if util.IsFileExist(serverCrt) && util.IsFileExist(serverKey) {
			if err = this.tlnetMq.HttpStartTLS(addr, serverCrt, serverKey); err != nil {
				err = errors.New("tlMq tls start failed:" + err.Error())
			}
		} else {
			if err = this.tlnetMq.HttpStartTlsBytes(addr, []byte(keystore.ServerCrt), []byte(keystore.ServerKey)); err != nil {
				err = errors.New("tlMq tls by bytes start failed:" + err.Error())
			}
		}
	}
	if !this.isClose {
		log.LoggerSys.Info(sys.SysLog(fmt.Sprint("tlMq start[", sys.MQADDR, "/mq]")))
		if err = this.tlnetMq.HttpStart(addr); err != nil {
			err = errors.New("tlMq start failed:" + err.Error())
		}
	}

	if !this.isClose && err != nil {
		log.LoggerSys.Error(err.Error())
		panic("tlMq start failed:" + err.Error())
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
	return wc
}

func mq2handler(hc *tlnet.HttpContext) {
	defer util.ErrRecovr()
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
		if mb, err := pullByte(bs[9:]); err == nil {
			me := MQEncode(mb)
			buf := util.BufferPool.Get(1 + len(me))
			buf.WriteByte(MQ_PULLBYTE)
			buf.Write(me)
			hc.ResponseBytes(0, buf.Bytes())
		} else {
			hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_PULLBYTE)...))
		}
	case MQ_PULLJSON:
		if mb, err := pullJson(bs[9:]); err == nil {
			je := JEncode(mb)
			buf := util.BufferPool.Get(1 + len(je))
			buf.WriteByte(MQ_PULLJSON)
			buf.Write(je)
			hc.ResponseBytes(0, buf.Bytes())
		} else {
			hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_PULLJSON)...))
		}
	case MQ_CURRENTID:
		topic := string(bs[9:])
		if seq, err := Level2.SelectId(0, &TableStub{Tablename: key.Topic(topic)}); err == nil {
			buf := util.BufferPool.Get(9)
			buf.WriteByte(MQ_CURRENTID)
			buf.Write(util.Int64ToBytes(seq))
			hc.ResponseBytes(0, buf.Bytes())
		} else {
			hc.ResponseBytes(0, append([]byte{MQ_ERROR}, util.Int64ToBytes(MQ_ERROR_CURRENTID)...))
		}
	}
}

func mqHandler(hc *tlnet.HttpContext) {
	defer util.ErrRecovr()
	bs := hc.WS.Read()
	if bs == nil || (bs[0] != MQ_AUTH && !isAuth(hc.WS)) {
		hc.WS.Close()
		return
	}
	switch bs[0] {
	case MQ_AUTH:
		if mqAuth(string(bs[9:])) {
			cliId := int64(0)
			if len(bs) >= 17 {
				cliId = util.BytesToInt64(bs[9:17])
			}
			AddConn(hc.WS, cliId)
		} else {
			sendErr(MQ_ERROR_NOPASS, hc.WS)
			hc.WS.Close()
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
		if err := pubByte(bs[9:]); err != nil {
			sendErr(MQ_ERROR_PUBBYTE, hc.WS)
		}
	case MQ_PUBJSON:
		if err := pubJson(bs[9:]); err != nil {
			sendErr(MQ_ERROR_PUBJSON, hc.WS)
		}
	case MQ_SUB:
		sub(string(bs[9:]), 0, hc.WS)
	case MQ_SUBCANCEL:
		sub(string(bs[9:]), 1, hc.WS)
	case MQ_PULLBYTE:
		if mb, err := pullByte(bs[9:]); err == nil {
			MqWare.PullByte(mb, hc.WS.Id)
		} else {
			sendErr(MQ_ERROR_PULLBYTE, hc.WS)
		}
	case MQ_PULLJSON:
		if mb, err := pullJson(bs[9:]); err == nil {
			MqWare.PullJson(mb, hc.WS.Id)
		} else {
			sendErr(MQ_ERROR_PULLJSON, hc.WS)
		}
	case MQ_PING:
		MqWare.Ping(bytes.NewBuffer(bs), hc.WS.Id)
	case MQ_PUBMEM:
		if err := pubMem(bs[9:]); err != nil {
			sendErr(MQ_ERROR_PUBMEM, hc.WS)
		}
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
		topic := string(bs[9:])
		if seq, err := Level2.SelectId(0, &TableStub{Tablename: key.Topic(topic)}); err == nil {
			MqWare.PullId(seq)
		} else {
			sendErr(MQ_ERROR_CURRENTID, hc.WS)
		}
	case MQ_ACK:
		if IsRecvAckOn(hc.WS.Id) {
			ackId := util.BytesToInt64(bs[1:9])
			tlmq.DelAckId(hc.WS.Id, ackId)
		}
	}
	if bs[0] != MQ_PING && bs[0] != MQ_ACK {
		ack := bs[1:9]
		buf := util.BufferPool.Get(9)
		buf.WriteByte(MQ_ACK)
		buf.Write(ack)
		MqWare.Ack(buf, hc.WS.Id)
	}
}

func isAuth(ws *tlnet.Websocket) bool {
	return WsWare.Has(ws.Id)
}

/******************************************************/
func pubByte(bs []byte) (err error) {
	mb := &MqBean{}
	if _, err = MQDecode(bs, mb); err == nil {
		var seq int64
		if seq, err = saveTopic(mb.Topic, mb.Msg); err == nil && seq > 0 {
			mb.ID = seq
			MqWare.PubByte(mb)
			level1.BroadcastMQ(int8(MQ_PUBBYTE), MQEncode(mb))
		}
	}
	return
}

func pubJson(bs []byte) (err error) {
	var mb *JMqBean
	if mb, err = JDecode(bs); err == nil {
		var seq int64
		if seq, err = saveTopic(mb.Topic, []byte(mb.Msg)); err == nil && seq > 0 {
			mb.Id = seq
			MqWare.PubJson(mb)
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

func pubMem(bs []byte) (err error) {
	var mb *JMqBean
	if mb, err = JDecode(bs); err == nil {
		MqWare.PubMem(mb)
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

func sendErr(errCode int64, ws *tlnet.Websocket) {
	errbs := util.Int64ToBytes(errCode)
	ws.Send(append([]byte{MQ_ERROR}, errbs...))
}
