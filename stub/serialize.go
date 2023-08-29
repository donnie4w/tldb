// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
package stub

import . "github.com/donnie4w/gofer/buffer"

func EncodeTableStub(ts *TableStub) []byte {
	encode := NewEncodeStub()
	encode.StringToBytes(ts.Tablename)
	encode.Int64ToBytes(ts.ID)
	encode.MapFieldToBytes(ts.Field)
	encode.MapIdxToBytes(ts.Idx)
	return encode.ToBytes()
}

func DecodeTableStub(bs []byte) (ts *TableStub) {
	ts = &TableStub{}
	decode := &DecodeStub{Bytes: bs}
	ts.Tablename = decode.ToString()
	ts.ID = decode.ToInt64()
	ts.Field = decode.ToMapField()
	ts.Idx = decode.ToIdx()
	return
}

type EncodeStub struct {
	Buf *Buffer
}

func NewEncodeStub() *EncodeStub {
	return &EncodeStub{NewBufferByPool()}
}

func (this *EncodeStub) Int64ToBytes(v int64) {
	this.Buf.Write(int64ToBytes(v))
}

func (this *EncodeStub) StringToBytes(s string) {
	this.Buf.Write(int16ToBytes(int16(len(s))))
	this.Buf.WriteString(s)
}

func (this *EncodeStub) StringArrayToBytes(arr []string) {
	length := this.StringArrayLen(arr)
	this.Buf.Write(int32ToBytes(length))
	if arr != nil {
		for _, v := range arr {
			this.Buf.WriteByte(byte(len(v))) //key maxlength:127
			this.Buf.WriteString(v)
		}
	}
}

func (this *EncodeStub) MapIdxToBytes(m map[string]int8) {
	length := this.IdxLen(m)
	this.Buf.Write(int32ToBytes(length))
	if m != nil {
		for k := range m {
			this.Buf.WriteByte(byte(len(k)))
			this.Buf.WriteString(k)
		}
	}
}

func (this *EncodeStub) MapFieldToBytes(m map[string][]byte) {
	this.Buf.Write(this.MapFieldLen(m))
	if m != nil {
		for k, v := range m {
			this.Buf.WriteByte(byte(len(k)))
			this.Buf.WriteString(k)
			this.Buf.Write(int32ToBytes(int32(len(v))))
			this.Buf.Write(v)
		}
	}
}

func (this *EncodeStub) ToBytes() []byte {
	return this.Buf.Bytes()
}

type DecodeStub struct {
	Bytes  []byte
	custor int
}

func (this *DecodeStub) ToString() (_r string) {
	bs := this.Bytes[this.custor:]
	tnLen := bytesToInt16(bs[:2])
	_r = string(bs[2 : 2+tnLen])
	this.custor += 2 + int(tnLen)
	return
}

func (this *DecodeStub) ToInt64() (_r int64) {
	bs := this.Bytes[this.custor:]
	_r = bytesToInt64(bs[:8])
	this.custor += 8
	return
}

func (this *DecodeStub) ToMapField() map[string][]byte {
	bs := this.Bytes[this.custor:]
	fLen := bytesToInt32(bs[:4])
	bs = bs[4 : 4+fLen]
	i := 0
	m := make(map[string][]byte, 0)
	for i < len(bs) {
		kLen := int(int8(bs[i]))
		vlen := int(bytesToInt32(bs[i+1+kLen : i+1+kLen+4]))
		m[string(bs[i+1:i+1+kLen])] = bs[i+1+kLen+4 : i+1+kLen+4+vlen]
		i = i + 1 + kLen + 4 + vlen
	}
	this.custor += 4 + int(fLen)
	return m
}

func (this *EncodeStub) MapFieldLen(m map[string][]byte) (_r []byte) {
	i := int32(0)
	if m != nil {
		for k, v := range m {
			i += 1
			i += int32(len([]byte(k)))
			i += 4
			i += int32(len(v))
		}
	}
	_r = int32ToBytes(i)
	return
}

func (this *EncodeStub) StringArrayLen(arr []string) (_r int32) {
	if arr != nil {
		for _, v := range arr {
			_r += int32(len([]byte(v))) + 1
		}
	}
	return
}

func (this *EncodeStub) IdxLen(m map[string]int8) (_r int32) {
	if m != nil {
		for v := range m {
			_r += int32(len([]byte(v)))
		}
	}
	return
}

func (this *DecodeStub) ToIdx() map[string]int8 {
	bs := this.Bytes[this.custor:]
	length := bytesToInt32(bs[:4])
	bs = bs[4 : 4+length]
	i := 0
	m := make(map[string]int8, 0)
	for i < len(bs) {
		kLen := int(int8(bs[i]))
		m[string(bs[i+1:i+1+kLen])] = 0
		i = i + 1 + kLen
	}
	this.custor += 4 + int(length)
	return m
}

func (this *DecodeStub) ToStrArray() []string {
	bs := this.Bytes[this.custor:]
	length := bytesToInt32(bs[:4])
	bs = bs[4 : 4+length]
	i := 0
	m := make([]string, 0)
	for i < len(bs) {
		kLen := int(int8(bs[i]))
		m = append(m, string(bs[i+1:i+1+kLen]))
		i += 1 + kLen
	}
	this.custor += 4 + int(length)
	return m
}

/**********************************************************/
func int64ToBytes(n int64) (bs []byte) {
	bs = make([]byte, 8)
	for i := 0; i < 8; i++ {
		bs[i] = byte(n >> (8 * (7 - i)))
	}
	return
}

func bytesToInt64(bs []byte) (_r int64) {
	if len(bs) >= 8 {
		for i := 0; i < 8; i++ {
			_r = _r | int64(bs[i])<<(8*(7-i))
		}
	} else {
		bs8 := make([]byte, 8)
		for i, b := range bs {
			bs8[7-i] = b
		}
		_r = bytesToInt64(bs8)
	}
	return
}

func int32ToBytes(n int32) (bs []byte) {
	bs = make([]byte, 4)
	for i := 0; i < 4; i++ {
		bs[i] = byte(n >> (8 * (3 - i)))
	}
	return
}

func bytesToInt32(bs []byte) (_r int32) {
	for i := 0; i < 4; i++ {
		_r = _r | int32(bs[i])<<(8*(3-i))
	}
	return
}

func int16ToBytes(n int16) (bs []byte) {
	bs = make([]byte, 2)
	for i := 0; i < 2; i++ {
		bs[i] = byte(n >> (8 * (1 - i)))
	}
	return
}

func bytesToInt16(bs []byte) (_r int16) {
	for i := 0; i < 2; i++ {
		_r = _r | int16(bs[i])<<(8*(1-i))
	}
	return
}

/****************************************************/
