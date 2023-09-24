/**
 * Copyright 2023 tldb Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package tlmq

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/donnie4w/gothrift/thrift"
	// "github.com/apache/thrift/lib/go/thrift"
)

func JEncode(mb *JMqBean) (bs []byte) {
	bs, _ = json.Marshal(mb)
	return
}

func JDecode(bs []byte) (mb *JMqBean, err error) {
	err = json.Unmarshal(bs, &mb)
	return
}

type JMqBean struct {
	Id    int64
	Topic string
	Msg   string
}

/******************************************/
func MQEncode(ts thrift.TStruct) (_r []byte) {
	buf := thrift.NewTMemoryBuffer()
	tcf := thrift.NewTCompactProtocolFactory()
	tp := tcf.GetProtocol(buf)
	ts.Write(context.Background(), tp)
	_r = buf.Bytes()
	return
}

func MQDecode[T thrift.TStruct](bs []byte, ts T) (_r T, err error) {
	buf := thrift.NewTMemoryBuffer()
	buf.Buffer = bytes.NewBuffer(bs)
	tcf := thrift.NewTCompactProtocolFactory()
	tp := tcf.GetProtocol(buf)
	err = ts.Read(context.Background(), tp)
	return ts, err
}
