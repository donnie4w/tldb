// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package util

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"

	"github.com/donnie4w/gothrift/thrift"

	// "github.com/apache/thrift/lib/go/thrift"
	"github.com/golang/snappy"
)

func Encode(e any) (by []byte, err error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(e)
	by = buf.Bytes()
	return
}

func Decode(buf []byte, e any) (err error) {
	decoder := gob.NewDecoder(bytes.NewReader(buf))
	err = decoder.Decode(e)
	return
}

func Int64ToBytes(n int64) (bs []byte) {
	bs = make([]byte, 8)
	for i := 0; i < 8; i++ {
		bs[i] = byte(n >> (8 * (7 - i)))
	}
	return
}

func BytesToInt64(bs []byte) (_r int64) {
	if len(bs) >= 8 {
		for i := 0; i < 8; i++ {
			_r = _r | int64(bs[i])<<(8*(7-i))
		}
	} else {
		bs8 := make([]byte, 8)
		for i, b := range bs {
			bs8[i+8-len(bs)] = b
		}
		_r = BytesToInt64(bs8)
	}
	return
}

func Int32ToBytes(n int32) (bs []byte) {
	bs = make([]byte, 4)
	for i := 0; i < 4; i++ {
		bs[i] = byte(n >> (8 * (3 - i)))
	}
	return
}

func Int16ToBytes(n int16) (bs []byte) {
	bs = make([]byte, 2)
	for i := 0; i < 2; i++ {
		bs[i] = byte(n >> (8 * (1 - i)))
	}
	return
}

func BytesToInt32(bs []byte) (_r int32) {
	if len(bs) >= 4 {
		for i := 0; i < 4; i++ {
			_r = _r | int32(bs[i])<<(8*(3-i))
		}
	} else {
		bs4 := make([]byte, 4)
		for i, b := range bs {
			bs4[i+4-len(bs)] = b
		}
		_r = BytesToInt32(bs4)
	}
	return
}

func BytesToInt16(bs []byte) (_r int16) {
	if len(bs) >= 2 {
		for i := 0; i < 2; i++ {
			_r = _r | int16(bs[i])<<(8*(1-i))
		}
	} else {
		bs2 := make([]byte, 2)
		for i, b := range bs {
			bs2[i+2-len(bs)] = b
		}
		_r = BytesToInt16(bs2)
	}
	return
}

func IntArrayToBytes(n []int64) []byte {
	bytesBuffer := BufferPool.Get(8 * len(n))
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}

func BytesToIntArray(bs []byte) (data []int64) {
	bytesBuffer := bytes.NewBuffer(bs)
	data = make([]int64, len(bs)/8)
	binary.Read(bytesBuffer, binary.BigEndian, data)
	return
}

var tconf = &thrift.TConfiguration{}

func TEncode(ts thrift.TStruct) (_r []byte) {
	buf := &thrift.TMemoryBuffer{Buffer: bytes.NewBuffer([]byte{})}
	protocol := thrift.NewTCompactProtocolConf(buf, tconf)
	ts.Write(context.Background(), protocol)
	protocol.Flush(context.Background())
	_r = buf.Bytes()
	return
}

func TDecode[T thrift.TStruct](bs []byte, ts T) (_r T, err error) {
	buf := &thrift.TMemoryBuffer{Buffer: bytes.NewBuffer(bs)}
	protocol := thrift.NewTCompactProtocolConf(buf, tconf)
	err = ts.Read(context.Background(), protocol)
	return ts, err
}

func ZlibCz(bs []byte) (_r []byte, err error) {
	var buf bytes.Buffer
	var compressor *zlib.Writer
	if compressor, err = zlib.NewWriterLevel(&buf, zlib.BestCompression); err == nil {
		defer compressor.Close()
		compressor.Write(bs)
		compressor.Flush()
		_r = buf.Bytes()
	} else {
		_r = bs
	}
	return
}

func ZlibUnCz(bs []byte) (_r []byte, err error) {
	var obuf bytes.Buffer
	var read io.ReadCloser
	if read, err = zlib.NewReader(bytes.NewReader(bs)); err == nil {
		defer read.Close()
		io.Copy(&obuf, read)
		_r = obuf.Bytes()
	} else {
		_r = bs
	}
	return
}

func Gzip(gzfname, srcfname, dir string) (err error) {
	var gf *os.File
	if gf, err = os.Create(gzfname); err == nil {
		defer gf.Close()
		var f1 *os.File
		if f1, err = os.Open(dir + "/" + srcfname); err == nil {
			defer f1.Close()
			gw := gzip.NewWriter(gf)
			defer gw.Close()
			gw.Header.Name = srcfname
			var buf bytes.Buffer
			io.Copy(&buf, f1)
			_, err = gw.Write(buf.Bytes())
		}
	}
	return
}

func GzipWrite(bs []byte) (buf bytes.Buffer, err error) {
	gw := gzip.NewWriter(&buf)
	defer gw.Close()
	_, err = gw.Write(bs)
	return
}

func UnGzip(bs []byte) (_bb bytes.Buffer, err error) {
	if gz, er := gzip.NewReader(bytes.NewBuffer(bs)); er == nil {
		defer gz.Close()
		var bs = make([]byte, 1024)
		for {
			var n int
			n, err := gz.Read(bs)
			if (err != nil && err != io.EOF) || n == 0 {
				break
			}
			_bb.Write(bs[:n])
		}
	} else {
		err = er
	}
	return
}

func UnGzipByFile(gzfname string, f func(bs []byte) bool) {
	if fgzip, err := os.Open(gzfname); err == nil {
		defer fgzip.Close()
		if gz, err := gzip.NewReader(fgzip); err == nil {
			defer gz.Close()
			var bs = make([]byte, 1024)
			for {
				n, err := gz.Read(bs)
				if (err != nil && err != io.EOF) || n == 0 {
					break
				}
				if !f(bs[:n]) {
					break
				}
			}
		}
	}
}

func UnGzipByBytes(bs []byte, f func(bs []byte) bool) (err error) {
	if buf := bytes.NewReader(bs); buf != nil {
		var gz *gzip.Reader
		if gz, err = gzip.NewReader(buf); err == nil {
			defer gz.Close()
			var bs = make([]byte, 1024)
			for {
				n, err := gz.Read(bs)
				if (err != nil && err != io.EOF) || n == 0 {
					break
				}
				if !f(bs[:n]) {
					break
				}
			}
		}
	}
	return
}

func CheckGzip(bs []byte) (err error) {
	if buf := bytes.NewReader(bs); buf != nil {
		_, err = gzip.NewReader(buf)
	}
	return
}

func SnappyEncode(bs []byte) (_r []byte) {
	return snappy.Encode(nil, bs)
}

func SnappyDecode(bs []byte) (_r []byte) {
	_r, _ = snappy.Decode(nil, bs)
	return
}
