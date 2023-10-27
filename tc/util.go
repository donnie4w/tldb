// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/donnie4w/tldb/util"
)

func type2Name(bs []byte) (_r string) {
	switch string(bs) {
	case "0":
		_r = "String"
	case "1":
		_r = "INT64"
	case "2":
		_r = "INT32"
	case "3":
		_r = "INT16"
	case "4":
		_r = "INT8"
	case "5":
		_r = "FLOAT64"
	case "6":
		_r = "FLOAT32"
	case "7":
		_r = "BINARY"
	case "8":
		_r = "Byte"
	case "9":
		_r = "Unsigned INT64"
	case "10":
		_r = "Unsigned INT32"
	case "11":
		_r = "Unsigned INT16"
	case "12":
		_r = "Unsigned INT8"
	default:
	}
	return
}

func type2value(_type, value []byte) (_r string) {
	var err error
	switch string(_type) {
	case "0":
		_r = string(value)
	case "1":
		if value != nil && len(value) > 0 {
			_r = fmt.Sprint(util.BytesToInt64(value))
		}
	case "2":
		if value != nil && len(value) > 0 {
			_r = fmt.Sprint(util.BytesToInt32(value))
		}
	case "3":
		if value != nil && len(value) > 0 {
			_r = fmt.Sprint(util.BytesToInt16(value))
		}
	case "4":
		var v int8
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &v); err == nil {
			_r = fmt.Sprint(v)
		}
	case "5":
		var f float64
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &f); err == nil {
			_r = fmt.Sprint(f)
		}
	case "6":
		var f float32
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &f); err == nil {
			_r = fmt.Sprint(f)
		}
	case "7":
		_r = string(value)
	case "8":
		if value != nil && len(value) > 0 {
			_r = string(value[0])
		}
	case "9":
		var v uint64
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &v); err == nil {
			_r = fmt.Sprint(v)
		}
	case "10":
		var v uint32
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &v); err == nil {
			_r = fmt.Sprint(v)
		}
	case "11":
		var v uint16
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &v); err == nil {
			_r = fmt.Sprint(v)
		}
	case "12":
		var v uint8
		if err = binary.Read(bytes.NewBuffer(value), binary.BigEndian, &v); err == nil {
			_r = fmt.Sprint(v)
		}
	default:
		_r = string(value)
	}
	if err != nil {
		_r = string(value)
	}
	return
}

func valueToBytes(_type []byte, value string) (_r []byte, err error) {
	switch string(_type) {
	case "0":
		_r = []byte(value)
	case "1": //int64
		var v int64
		if v, err = strconv.ParseInt(value, 10, 64); err == nil {
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &v)
			_r = buf.Bytes()
		}
	case "2": //int32
		var v int64
		if v, err = strconv.ParseInt(value, 10, 64); err == nil {
			i := int32(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	case "3": //int16
		var v int64
		if v, err = strconv.ParseInt(value, 10, 64); err == nil {
			i := int16(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	case "4": //int8
		var v int64
		if v, err = strconv.ParseInt(value, 10, 64); err == nil {
			i := int8(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	case "5": //float64
		var v float64
		if v, err = strconv.ParseFloat(value, 64); err == nil {
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &v)
			_r = buf.Bytes()
		}
	case "6": //float32
		var v float64
		if v, err = strconv.ParseFloat(value, 32); err == nil {
			i := float32(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	case "7": //BINARY
		_r = []byte(value)
	case "8": //byte
		if value != "" {
			_r = []byte{value[0]}
		}
	case "9": //uint64
		var v uint64
		if v, err = strconv.ParseUint(value, 10, 64); err == nil {
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &v)
			_r = buf.Bytes()
		}
	case "10": //uint32
		var v uint64
		if v, err = strconv.ParseUint(value, 10, 64); err == nil {
			i := uint32(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	case "11": //uint16
		var v uint64
		if v, err = strconv.ParseUint(value, 10, 64); err == nil {
			i := uint16(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	case "12": //uint8
		var v uint64
		if v, err = strconv.ParseUint(value, 10, 64); err == nil {
			i := uint8(v)
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, &i)
			_r = buf.Bytes()
		}
	}
	if err != nil {
		err = errors.New(fmt.Sprint("type conversion error:", value))
	}
	if _r == nil && err == nil {
		_r = []byte(value)
	}
	return
}