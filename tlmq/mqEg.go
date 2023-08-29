/**
 * Copyright 2023 tldb Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package tlmq

import (
	"bytes"

	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tlnet"
)

const (
	MQ_AUTH      byte = 1
	MQ_PUBBYTE   byte = 2
	MQ_PUBJSON   byte = 3
	MQ_SUB       byte = 4
	MQ_PULLBYTE  byte = 5
	MQ_PULLJSON  byte = 6
	MQ_PING      byte = 7
	MQ_ERROR     byte = 8
	MQ_PUBMEM    byte = 9
	MQ_RECVACK   byte = 10
	MQ_MERGE     byte = 11
	MQ_SUBCANCEL byte = 12
	MQ_CURRENTID byte = 13
	MQ_ZLIB      byte = 14
	MQ_ACK       byte = 0

	MQ_ERROR_PUBBYTE   int64 = 1201 // 发布pubByte错误
	MQ_ERROR_PUBJSON   int64 = 1202 // 发布pubJson错误
	MQ_ERROR_PULLBYTE  int64 = 1203 // 发布pullByte错误
	MQ_ERROR_PULLJSON  int64 = 1204 // 发布pullJson错误
	MQ_ERROR_PUBMEM    int64 = 1205 // 发布pubMem错误
	MQ_ERROR_CURRENTID int64 = 1206 // 发布CURRENTID错误
	MQ_ERROR_NOPASS    int64 = 1301 // 验证不通过
)

type MqEg interface {
	AddConn(topic string, ws *tlnet.Websocket)
	DelConn(id int64)
	SubCount(topic string) (_r int64)
	DelTopic(topic string)
	DelTopicWithID(topic string, id int64)
	ClusPub(mqType int8, bs []byte) (err error)
	PubByte(mb *MqBean) (err error)
	PubJson(mb *JMqBean) (err error)
	PubMem(mb *JMqBean)
	PullByte(mb *MqBean, id int64) (err error)
	PullJson(mb *JMqBean, id int64) (err error)
	PullId(id int64) (err error)
	Ping(buf *bytes.Buffer, id int64)
	Ack(buf *bytes.Buffer, id int64)
}
