// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package level1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/donnie4w/gothrift/thrift"

	// thrift "github.com/apache/thrift/lib/go/thrift"
	"regexp"
	"strings"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = errors.New
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal
// (needed by validator.)
var _ = strings.Contains
var _ = regexp.MatchString

// Attributes:
//  - Kvmap
//  - Dels
type BatchPacket struct {
  Kvmap map[string][]byte `thrift:"kvmap,1" db:"kvmap" json:"kvmap,omitempty"`
  Dels []string `thrift:"dels,2" db:"dels" json:"dels,omitempty"`
}

func NewBatchPacket() *BatchPacket {
  return &BatchPacket{}
}

var BatchPacket_Kvmap_DEFAULT map[string][]byte

func (p *BatchPacket) GetKvmap() map[string][]byte {
  return p.Kvmap
}
var BatchPacket_Dels_DEFAULT []string

func (p *BatchPacket) GetDels() []string {
  return p.Dels
}
func (p *BatchPacket) IsSetKvmap() bool {
  return p.Kvmap != nil
}

func (p *BatchPacket) IsSetDels() bool {
  return p.Dels != nil
}

func (p *BatchPacket) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.LIST {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *BatchPacket)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string][]byte, size)
  p.Kvmap =  tMap
  for i := 0; i < size; i ++ {
var _key0 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key0 = v
}
var _val1 []byte
    if v, err := iprot.ReadBinary(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val1 = v
}
    p.Kvmap[_key0] = _val1
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *BatchPacket)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]string, 0, size)
  p.Dels =  tSlice
  for i := 0; i < size; i ++ {
var _elem2 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _elem2 = v
}
    p.Dels = append(p.Dels, _elem2)
  }
  if err := iprot.ReadListEnd(ctx); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *BatchPacket) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "BatchPacket"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *BatchPacket) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetKvmap() {
    if err := oprot.WriteFieldBegin(ctx, "kvmap", thrift.MAP, 1); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:kvmap: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Kvmap)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Kvmap {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteBinary(ctx, v); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 1:kvmap: ", p), err) }
  }
  return err
}

func (p *BatchPacket) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetDels() {
    if err := oprot.WriteFieldBegin(ctx, "dels", thrift.LIST, 2); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:dels: ", p), err) }
    if err := oprot.WriteListBegin(ctx, thrift.STRING, len(p.Dels)); err != nil {
      return thrift.PrependError("error writing list begin: ", err)
    }
    for _, v := range p.Dels {
      if err := oprot.WriteString(ctx, string(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteListEnd(ctx); err != nil {
      return thrift.PrependError("error writing list end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 2:dels: ", p), err) }
  }
  return err
}

func (p *BatchPacket) Equals(other *BatchPacket) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if len(p.Kvmap) != len(other.Kvmap) { return false }
  for k, _tgt := range p.Kvmap {
    _src3 := other.Kvmap[k]
    if bytes.Compare(_tgt, _src3) != 0 { return false }
  }
  if len(p.Dels) != len(other.Dels) { return false }
  for i, _tgt := range p.Dels {
    _src4 := other.Dels[i]
    if _tgt != _src4 { return false }
  }
  return true
}

func (p *BatchPacket) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("BatchPacket(%+v)", *p)
}

func (p *BatchPacket) Validate() error {
  return nil
}
// Attributes:
//  - Addr
//  - UUID
//  - Stat
//  - MqAddr
//  - CliAddr
//  - AdminAddr
//  - Nodekv
//  - Ns
type Node struct {
  Addr string `thrift:"addr,1,required" db:"addr" json:"addr"`
  UUID int64 `thrift:"uuid,2,required" db:"uuid" json:"uuid"`
  Stat int8 `thrift:"stat,3,required" db:"stat" json:"stat"`
  MqAddr string `thrift:"mqAddr,4,required" db:"mqAddr" json:"mqAddr"`
  CliAddr string `thrift:"cliAddr,5,required" db:"cliAddr" json:"cliAddr"`
  AdminAddr string `thrift:"adminAddr,6,required" db:"adminAddr" json:"adminAddr"`
  Nodekv map[int64]string `thrift:"nodekv,7" db:"nodekv" json:"nodekv,omitempty"`
  Ns string `thrift:"ns,8,required" db:"ns" json:"ns"`
}

func NewNode() *Node {
  return &Node{}
}


func (p *Node) GetAddr() string {
  return p.Addr
}

func (p *Node) GetUUID() int64 {
  return p.UUID
}

func (p *Node) GetStat() int8 {
  return p.Stat
}

func (p *Node) GetMqAddr() string {
  return p.MqAddr
}

func (p *Node) GetCliAddr() string {
  return p.CliAddr
}

func (p *Node) GetAdminAddr() string {
  return p.AdminAddr
}
var Node_Nodekv_DEFAULT map[int64]string

func (p *Node) GetNodekv() map[int64]string {
  return p.Nodekv
}

func (p *Node) GetNs() string {
  return p.Ns
}
func (p *Node) IsSetNodekv() bool {
  return p.Nodekv != nil
}

func (p *Node) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetAddr bool = false;
  var issetUUID bool = false;
  var issetStat bool = false;
  var issetMqAddr bool = false;
  var issetCliAddr bool = false;
  var issetAdminAddr bool = false;
  var issetNs bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetAddr = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetUUID = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetStat = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
        issetMqAddr = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
        issetCliAddr = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField6(ctx, iprot); err != nil {
          return err
        }
        issetAdminAddr = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 7:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField7(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 8:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField8(ctx, iprot); err != nil {
          return err
        }
        issetNs = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetAddr{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Addr is not set"));
  }
  if !issetUUID{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field UUID is not set"));
  }
  if !issetStat{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Stat is not set"));
  }
  if !issetMqAddr{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field MqAddr is not set"));
  }
  if !issetCliAddr{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field CliAddr is not set"));
  }
  if !issetAdminAddr{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field AdminAddr is not set"));
  }
  if !issetNs{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Ns is not set"));
  }
  return nil
}

func (p *Node)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Addr = v
}
  return nil
}

func (p *Node)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.UUID = v
}
  return nil
}

func (p *Node)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  temp := int8(v)
  p.Stat = temp
}
  return nil
}

func (p *Node)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.MqAddr = v
}
  return nil
}

func (p *Node)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.CliAddr = v
}
  return nil
}

func (p *Node)  ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  p.AdminAddr = v
}
  return nil
}

func (p *Node)  ReadField7(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[int64]string, size)
  p.Nodekv =  tMap
  for i := 0; i < size; i ++ {
var _key5 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key5 = v
}
var _val6 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val6 = v
}
    p.Nodekv[_key5] = _val6
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *Node)  ReadField8(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 8: ", err)
} else {
  p.Ns = v
}
  return nil
}

func (p *Node) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Node"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
    if err := p.writeField6(ctx, oprot); err != nil { return err }
    if err := p.writeField7(ctx, oprot); err != nil { return err }
    if err := p.writeField8(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *Node) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "addr", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:addr: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Addr)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.addr (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:addr: ", p), err) }
  return err
}

func (p *Node) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "uuid", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:uuid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.UUID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.uuid (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:uuid: ", p), err) }
  return err
}

func (p *Node) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "stat", thrift.BYTE, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:stat: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Stat)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.stat (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:stat: ", p), err) }
  return err
}

func (p *Node) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "mqAddr", thrift.STRING, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:mqAddr: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.MqAddr)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.mqAddr (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:mqAddr: ", p), err) }
  return err
}

func (p *Node) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "cliAddr", thrift.STRING, 5); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:cliAddr: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.CliAddr)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.cliAddr (5) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 5:cliAddr: ", p), err) }
  return err
}

func (p *Node) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "adminAddr", thrift.STRING, 6); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:adminAddr: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.AdminAddr)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.adminAddr (6) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 6:adminAddr: ", p), err) }
  return err
}

func (p *Node) writeField7(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetNodekv() {
    if err := oprot.WriteFieldBegin(ctx, "nodekv", thrift.MAP, 7); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:nodekv: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.I64, thrift.STRING, len(p.Nodekv)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Nodekv {
      if err := oprot.WriteI64(ctx, int64(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteString(ctx, string(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 7:nodekv: ", p), err) }
  }
  return err
}

func (p *Node) writeField8(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ns", thrift.STRING, 8); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:ns: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Ns)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ns (8) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 8:ns: ", p), err) }
  return err
}

func (p *Node) Equals(other *Node) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Addr != other.Addr { return false }
  if p.UUID != other.UUID { return false }
  if p.Stat != other.Stat { return false }
  if p.MqAddr != other.MqAddr { return false }
  if p.CliAddr != other.CliAddr { return false }
  if p.AdminAddr != other.AdminAddr { return false }
  if len(p.Nodekv) != len(other.Nodekv) { return false }
  for k, _tgt := range p.Nodekv {
    _src7 := other.Nodekv[k]
    if _tgt != _src7 { return false }
  }
  if p.Ns != other.Ns { return false }
  return true
}

func (p *Node) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("Node(%+v)", *p)
}

func (p *Node) Validate() error {
  return nil
}
// Attributes:
//  - Seq
//  - Key
//  - Value
type LogKV struct {
  Seq int64 `thrift:"seq,1,required" db:"seq" json:"seq"`
  Key string `thrift:"key,2,required" db:"key" json:"key"`
  Value []byte `thrift:"value,3" db:"value" json:"value,omitempty"`
}

func NewLogKV() *LogKV {
  return &LogKV{}
}


func (p *LogKV) GetSeq() int64 {
  return p.Seq
}

func (p *LogKV) GetKey() string {
  return p.Key
}
var LogKV_Value_DEFAULT []byte

func (p *LogKV) GetValue() []byte {
  return p.Value
}
func (p *LogKV) IsSetValue() bool {
  return p.Value != nil
}

func (p *LogKV) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetSeq bool = false;
  var issetKey bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetSeq = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetKey = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetSeq{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Seq is not set"));
  }
  if !issetKey{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Key is not set"));
  }
  return nil
}

func (p *LogKV)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Seq = v
}
  return nil
}

func (p *LogKV)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Key = v
}
  return nil
}

func (p *LogKV)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Value = v
}
  return nil
}

func (p *LogKV) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "LogKV"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *LogKV) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "seq", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:seq: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Seq)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.seq (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:seq: ", p), err) }
  return err
}

func (p *LogKV) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "key", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:key: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Key)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.key (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:key: ", p), err) }
  return err
}

func (p *LogKV) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetValue() {
    if err := oprot.WriteFieldBegin(ctx, "value", thrift.STRING, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:value: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Value); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.value (3) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:value: ", p), err) }
  }
  return err
}

func (p *LogKV) Equals(other *LogKV) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Seq != other.Seq { return false }
  if p.Key != other.Key { return false }
  if bytes.Compare(p.Value, other.Value) != 0 { return false }
  return true
}

func (p *LogKV) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("LogKV(%+v)", *p)
}

func (p *LogKV) Validate() error {
  return nil
}
// Attributes:
//  - Ptype
//  - OpType
//  - Txid
//  - Fromuuid
//  - IsPass
//  - IsFirst
//  - IsExcu
//  - IsVerify
//  - DoCommit
//  - Mode
//  - ExecNode
//  - ClusterNode
//  - Key
//  - Value
//  - Batch
//  - Stat
//  - LoadBean
//  - TimeBean
//  - SeqBean
//  - RunUUID
//  - SyncId
//  - Pbtime
type PonBean struct {
  Ptype int16 `thrift:"ptype,1,required" db:"ptype" json:"ptype"`
  OpType int16 `thrift:"opType,2,required" db:"opType" json:"opType"`
  Txid int64 `thrift:"txid,3,required" db:"txid" json:"txid"`
  Fromuuid int64 `thrift:"fromuuid,4,required" db:"fromuuid" json:"fromuuid"`
  IsPass bool `thrift:"isPass,5,required" db:"isPass" json:"isPass"`
  IsFirst bool `thrift:"isFirst,6,required" db:"isFirst" json:"isFirst"`
  IsExcu bool `thrift:"isExcu,7,required" db:"isExcu" json:"isExcu"`
  IsVerify bool `thrift:"isVerify,8,required" db:"isVerify" json:"isVerify"`
  DoCommit bool `thrift:"doCommit,9,required" db:"doCommit" json:"doCommit"`
  Mode int8 `thrift:"mode,10,required" db:"mode" json:"mode"`
  ExecNode []int64 `thrift:"execNode,11" db:"execNode" json:"execNode,omitempty"`
  ClusterNode []int64 `thrift:"clusterNode,12" db:"clusterNode" json:"clusterNode,omitempty"`
  Key []byte `thrift:"key,13" db:"key" json:"key,omitempty"`
  Value []byte `thrift:"value,14" db:"value" json:"value,omitempty"`
  Batch *BatchPacket `thrift:"batch,15" db:"batch" json:"batch,omitempty"`
  Stat *Stat `thrift:"stat,16" db:"stat" json:"stat,omitempty"`
  LoadBean *LoadBean `thrift:"loadBean,17" db:"loadBean" json:"loadBean,omitempty"`
  TimeBean *TimeBean `thrift:"timeBean,18" db:"timeBean" json:"timeBean,omitempty"`
  SeqBean *LogSeqBean `thrift:"seqBean,19" db:"seqBean" json:"seqBean,omitempty"`
  RunUUID []int64 `thrift:"runUUID,20" db:"runUUID" json:"runUUID,omitempty"`
  SyncId *int64 `thrift:"syncId,21" db:"syncId" json:"syncId,omitempty"`
  Pbtime []byte `thrift:"pbtime,22" db:"pbtime" json:"pbtime,omitempty"`
}

func NewPonBean() *PonBean {
  return &PonBean{}
}


func (p *PonBean) GetPtype() int16 {
  return p.Ptype
}

func (p *PonBean) GetOpType() int16 {
  return p.OpType
}

func (p *PonBean) GetTxid() int64 {
  return p.Txid
}

func (p *PonBean) GetFromuuid() int64 {
  return p.Fromuuid
}

func (p *PonBean) GetIsPass() bool {
  return p.IsPass
}

func (p *PonBean) GetIsFirst() bool {
  return p.IsFirst
}

func (p *PonBean) GetIsExcu() bool {
  return p.IsExcu
}

func (p *PonBean) GetIsVerify() bool {
  return p.IsVerify
}

func (p *PonBean) GetDoCommit() bool {
  return p.DoCommit
}

func (p *PonBean) GetMode() int8 {
  return p.Mode
}
var PonBean_ExecNode_DEFAULT []int64

func (p *PonBean) GetExecNode() []int64 {
  return p.ExecNode
}
var PonBean_ClusterNode_DEFAULT []int64

func (p *PonBean) GetClusterNode() []int64 {
  return p.ClusterNode
}
var PonBean_Key_DEFAULT []byte

func (p *PonBean) GetKey() []byte {
  return p.Key
}
var PonBean_Value_DEFAULT []byte

func (p *PonBean) GetValue() []byte {
  return p.Value
}
var PonBean_Batch_DEFAULT *BatchPacket
func (p *PonBean) GetBatch() *BatchPacket {
  if !p.IsSetBatch() {
    return PonBean_Batch_DEFAULT
  }
return p.Batch
}
var PonBean_Stat_DEFAULT *Stat
func (p *PonBean) GetStat() *Stat {
  if !p.IsSetStat() {
    return PonBean_Stat_DEFAULT
  }
return p.Stat
}
var PonBean_LoadBean_DEFAULT *LoadBean
func (p *PonBean) GetLoadBean() *LoadBean {
  if !p.IsSetLoadBean() {
    return PonBean_LoadBean_DEFAULT
  }
return p.LoadBean
}
var PonBean_TimeBean_DEFAULT *TimeBean
func (p *PonBean) GetTimeBean() *TimeBean {
  if !p.IsSetTimeBean() {
    return PonBean_TimeBean_DEFAULT
  }
return p.TimeBean
}
var PonBean_SeqBean_DEFAULT *LogSeqBean
func (p *PonBean) GetSeqBean() *LogSeqBean {
  if !p.IsSetSeqBean() {
    return PonBean_SeqBean_DEFAULT
  }
return p.SeqBean
}
var PonBean_RunUUID_DEFAULT []int64

func (p *PonBean) GetRunUUID() []int64 {
  return p.RunUUID
}
var PonBean_SyncId_DEFAULT int64
func (p *PonBean) GetSyncId() int64 {
  if !p.IsSetSyncId() {
    return PonBean_SyncId_DEFAULT
  }
return *p.SyncId
}
var PonBean_Pbtime_DEFAULT []byte

func (p *PonBean) GetPbtime() []byte {
  return p.Pbtime
}
func (p *PonBean) IsSetExecNode() bool {
  return p.ExecNode != nil
}

func (p *PonBean) IsSetClusterNode() bool {
  return p.ClusterNode != nil
}

func (p *PonBean) IsSetKey() bool {
  return p.Key != nil
}

func (p *PonBean) IsSetValue() bool {
  return p.Value != nil
}

func (p *PonBean) IsSetBatch() bool {
  return p.Batch != nil
}

func (p *PonBean) IsSetStat() bool {
  return p.Stat != nil
}

func (p *PonBean) IsSetLoadBean() bool {
  return p.LoadBean != nil
}

func (p *PonBean) IsSetTimeBean() bool {
  return p.TimeBean != nil
}

func (p *PonBean) IsSetSeqBean() bool {
  return p.SeqBean != nil
}

func (p *PonBean) IsSetRunUUID() bool {
  return p.RunUUID != nil
}

func (p *PonBean) IsSetSyncId() bool {
  return p.SyncId != nil
}

func (p *PonBean) IsSetPbtime() bool {
  return p.Pbtime != nil
}

func (p *PonBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetPtype bool = false;
  var issetOpType bool = false;
  var issetTxid bool = false;
  var issetFromuuid bool = false;
  var issetIsPass bool = false;
  var issetIsFirst bool = false;
  var issetIsExcu bool = false;
  var issetIsVerify bool = false;
  var issetDoCommit bool = false;
  var issetMode bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I16 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetPtype = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I16 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetOpType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetTxid = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
        issetFromuuid = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
        issetIsPass = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField6(ctx, iprot); err != nil {
          return err
        }
        issetIsFirst = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 7:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField7(ctx, iprot); err != nil {
          return err
        }
        issetIsExcu = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 8:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField8(ctx, iprot); err != nil {
          return err
        }
        issetIsVerify = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 9:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField9(ctx, iprot); err != nil {
          return err
        }
        issetDoCommit = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 10:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField10(ctx, iprot); err != nil {
          return err
        }
        issetMode = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 11:
      if fieldTypeId == thrift.LIST {
        if err := p.ReadField11(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 12:
      if fieldTypeId == thrift.LIST {
        if err := p.ReadField12(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 13:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField13(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 14:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField14(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 15:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField15(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 16:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField16(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 17:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField17(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 18:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField18(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 19:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField19(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 20:
      if fieldTypeId == thrift.LIST {
        if err := p.ReadField20(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 21:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField21(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 22:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField22(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetPtype{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Ptype is not set"));
  }
  if !issetOpType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field OpType is not set"));
  }
  if !issetTxid{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Txid is not set"));
  }
  if !issetFromuuid{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Fromuuid is not set"));
  }
  if !issetIsPass{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsPass is not set"));
  }
  if !issetIsFirst{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsFirst is not set"));
  }
  if !issetIsExcu{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsExcu is not set"));
  }
  if !issetIsVerify{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsVerify is not set"));
  }
  if !issetDoCommit{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field DoCommit is not set"));
  }
  if !issetMode{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Mode is not set"));
  }
  return nil
}

func (p *PonBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI16(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Ptype = v
}
  return nil
}

func (p *PonBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI16(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.OpType = v
}
  return nil
}

func (p *PonBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Txid = v
}
  return nil
}

func (p *PonBean)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.Fromuuid = v
}
  return nil
}

func (p *PonBean)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.IsPass = v
}
  return nil
}

func (p *PonBean)  ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  p.IsFirst = v
}
  return nil
}

func (p *PonBean)  ReadField7(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 7: ", err)
} else {
  p.IsExcu = v
}
  return nil
}

func (p *PonBean)  ReadField8(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 8: ", err)
} else {
  p.IsVerify = v
}
  return nil
}

func (p *PonBean)  ReadField9(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 9: ", err)
} else {
  p.DoCommit = v
}
  return nil
}

func (p *PonBean)  ReadField10(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 10: ", err)
} else {
  temp := int8(v)
  p.Mode = temp
}
  return nil
}

func (p *PonBean)  ReadField11(ctx context.Context, iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]int64, 0, size)
  p.ExecNode =  tSlice
  for i := 0; i < size; i ++ {
var _elem8 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _elem8 = v
}
    p.ExecNode = append(p.ExecNode, _elem8)
  }
  if err := iprot.ReadListEnd(ctx); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *PonBean)  ReadField12(ctx context.Context, iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]int64, 0, size)
  p.ClusterNode =  tSlice
  for i := 0; i < size; i ++ {
var _elem9 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _elem9 = v
}
    p.ClusterNode = append(p.ClusterNode, _elem9)
  }
  if err := iprot.ReadListEnd(ctx); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *PonBean)  ReadField13(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 13: ", err)
} else {
  p.Key = v
}
  return nil
}

func (p *PonBean)  ReadField14(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 14: ", err)
} else {
  p.Value = v
}
  return nil
}

func (p *PonBean)  ReadField15(ctx context.Context, iprot thrift.TProtocol) error {
  p.Batch = &BatchPacket{}
  if err := p.Batch.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Batch), err)
  }
  return nil
}

func (p *PonBean)  ReadField16(ctx context.Context, iprot thrift.TProtocol) error {
  p.Stat = &Stat{}
  if err := p.Stat.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Stat), err)
  }
  return nil
}

func (p *PonBean)  ReadField17(ctx context.Context, iprot thrift.TProtocol) error {
  p.LoadBean = &LoadBean{}
  if err := p.LoadBean.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.LoadBean), err)
  }
  return nil
}

func (p *PonBean)  ReadField18(ctx context.Context, iprot thrift.TProtocol) error {
  p.TimeBean = &TimeBean{}
  if err := p.TimeBean.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.TimeBean), err)
  }
  return nil
}

func (p *PonBean)  ReadField19(ctx context.Context, iprot thrift.TProtocol) error {
  p.SeqBean = &LogSeqBean{}
  if err := p.SeqBean.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.SeqBean), err)
  }
  return nil
}

func (p *PonBean)  ReadField20(ctx context.Context, iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]int64, 0, size)
  p.RunUUID =  tSlice
  for i := 0; i < size; i ++ {
var _elem10 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _elem10 = v
}
    p.RunUUID = append(p.RunUUID, _elem10)
  }
  if err := iprot.ReadListEnd(ctx); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *PonBean)  ReadField21(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 21: ", err)
} else {
  p.SyncId = &v
}
  return nil
}

func (p *PonBean)  ReadField22(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 22: ", err)
} else {
  p.Pbtime = v
}
  return nil
}

func (p *PonBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "PonBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
    if err := p.writeField6(ctx, oprot); err != nil { return err }
    if err := p.writeField7(ctx, oprot); err != nil { return err }
    if err := p.writeField8(ctx, oprot); err != nil { return err }
    if err := p.writeField9(ctx, oprot); err != nil { return err }
    if err := p.writeField10(ctx, oprot); err != nil { return err }
    if err := p.writeField11(ctx, oprot); err != nil { return err }
    if err := p.writeField12(ctx, oprot); err != nil { return err }
    if err := p.writeField13(ctx, oprot); err != nil { return err }
    if err := p.writeField14(ctx, oprot); err != nil { return err }
    if err := p.writeField15(ctx, oprot); err != nil { return err }
    if err := p.writeField16(ctx, oprot); err != nil { return err }
    if err := p.writeField17(ctx, oprot); err != nil { return err }
    if err := p.writeField18(ctx, oprot); err != nil { return err }
    if err := p.writeField19(ctx, oprot); err != nil { return err }
    if err := p.writeField20(ctx, oprot); err != nil { return err }
    if err := p.writeField21(ctx, oprot); err != nil { return err }
    if err := p.writeField22(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *PonBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ptype", thrift.I16, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ptype: ", p), err) }
  if err := oprot.WriteI16(ctx, int16(p.Ptype)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ptype (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ptype: ", p), err) }
  return err
}

func (p *PonBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "opType", thrift.I16, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:opType: ", p), err) }
  if err := oprot.WriteI16(ctx, int16(p.OpType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.opType (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:opType: ", p), err) }
  return err
}

func (p *PonBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "txid", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:txid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Txid)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.txid (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:txid: ", p), err) }
  return err
}

func (p *PonBean) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "fromuuid", thrift.I64, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:fromuuid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Fromuuid)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.fromuuid (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:fromuuid: ", p), err) }
  return err
}

func (p *PonBean) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "isPass", thrift.BOOL, 5); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:isPass: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.IsPass)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.isPass (5) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 5:isPass: ", p), err) }
  return err
}

func (p *PonBean) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "isFirst", thrift.BOOL, 6); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:isFirst: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.IsFirst)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.isFirst (6) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 6:isFirst: ", p), err) }
  return err
}

func (p *PonBean) writeField7(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "isExcu", thrift.BOOL, 7); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:isExcu: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.IsExcu)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.isExcu (7) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 7:isExcu: ", p), err) }
  return err
}

func (p *PonBean) writeField8(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "isVerify", thrift.BOOL, 8); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:isVerify: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.IsVerify)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.isVerify (8) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 8:isVerify: ", p), err) }
  return err
}

func (p *PonBean) writeField9(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "doCommit", thrift.BOOL, 9); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 9:doCommit: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.DoCommit)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.doCommit (9) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 9:doCommit: ", p), err) }
  return err
}

func (p *PonBean) writeField10(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "mode", thrift.BYTE, 10); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:mode: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Mode)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.mode (10) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 10:mode: ", p), err) }
  return err
}

func (p *PonBean) writeField11(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetExecNode() {
    if err := oprot.WriteFieldBegin(ctx, "execNode", thrift.LIST, 11); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 11:execNode: ", p), err) }
    if err := oprot.WriteListBegin(ctx, thrift.I64, len(p.ExecNode)); err != nil {
      return thrift.PrependError("error writing list begin: ", err)
    }
    for _, v := range p.ExecNode {
      if err := oprot.WriteI64(ctx, int64(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteListEnd(ctx); err != nil {
      return thrift.PrependError("error writing list end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 11:execNode: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField12(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetClusterNode() {
    if err := oprot.WriteFieldBegin(ctx, "clusterNode", thrift.LIST, 12); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 12:clusterNode: ", p), err) }
    if err := oprot.WriteListBegin(ctx, thrift.I64, len(p.ClusterNode)); err != nil {
      return thrift.PrependError("error writing list begin: ", err)
    }
    for _, v := range p.ClusterNode {
      if err := oprot.WriteI64(ctx, int64(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteListEnd(ctx); err != nil {
      return thrift.PrependError("error writing list end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 12:clusterNode: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField13(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetKey() {
    if err := oprot.WriteFieldBegin(ctx, "key", thrift.STRING, 13); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 13:key: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Key); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.key (13) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 13:key: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField14(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetValue() {
    if err := oprot.WriteFieldBegin(ctx, "value", thrift.STRING, 14); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 14:value: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Value); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.value (14) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 14:value: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField15(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBatch() {
    if err := oprot.WriteFieldBegin(ctx, "batch", thrift.STRUCT, 15); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 15:batch: ", p), err) }
    if err := p.Batch.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Batch), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 15:batch: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField16(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetStat() {
    if err := oprot.WriteFieldBegin(ctx, "stat", thrift.STRUCT, 16); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 16:stat: ", p), err) }
    if err := p.Stat.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Stat), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 16:stat: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField17(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetLoadBean() {
    if err := oprot.WriteFieldBegin(ctx, "loadBean", thrift.STRUCT, 17); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 17:loadBean: ", p), err) }
    if err := p.LoadBean.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.LoadBean), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 17:loadBean: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField18(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetTimeBean() {
    if err := oprot.WriteFieldBegin(ctx, "timeBean", thrift.STRUCT, 18); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 18:timeBean: ", p), err) }
    if err := p.TimeBean.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.TimeBean), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 18:timeBean: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField19(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetSeqBean() {
    if err := oprot.WriteFieldBegin(ctx, "seqBean", thrift.STRUCT, 19); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 19:seqBean: ", p), err) }
    if err := p.SeqBean.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.SeqBean), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 19:seqBean: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField20(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetRunUUID() {
    if err := oprot.WriteFieldBegin(ctx, "runUUID", thrift.LIST, 20); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 20:runUUID: ", p), err) }
    if err := oprot.WriteListBegin(ctx, thrift.I64, len(p.RunUUID)); err != nil {
      return thrift.PrependError("error writing list begin: ", err)
    }
    for _, v := range p.RunUUID {
      if err := oprot.WriteI64(ctx, int64(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteListEnd(ctx); err != nil {
      return thrift.PrependError("error writing list end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 20:runUUID: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField21(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetSyncId() {
    if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 21); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 21:syncId: ", p), err) }
    if err := oprot.WriteI64(ctx, int64(*p.SyncId)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.syncId (21) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 21:syncId: ", p), err) }
  }
  return err
}

func (p *PonBean) writeField22(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetPbtime() {
    if err := oprot.WriteFieldBegin(ctx, "pbtime", thrift.STRING, 22); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 22:pbtime: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Pbtime); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.pbtime (22) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 22:pbtime: ", p), err) }
  }
  return err
}

func (p *PonBean) Equals(other *PonBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Ptype != other.Ptype { return false }
  if p.OpType != other.OpType { return false }
  if p.Txid != other.Txid { return false }
  if p.Fromuuid != other.Fromuuid { return false }
  if p.IsPass != other.IsPass { return false }
  if p.IsFirst != other.IsFirst { return false }
  if p.IsExcu != other.IsExcu { return false }
  if p.IsVerify != other.IsVerify { return false }
  if p.DoCommit != other.DoCommit { return false }
  if p.Mode != other.Mode { return false }
  if len(p.ExecNode) != len(other.ExecNode) { return false }
  for i, _tgt := range p.ExecNode {
    _src11 := other.ExecNode[i]
    if _tgt != _src11 { return false }
  }
  if len(p.ClusterNode) != len(other.ClusterNode) { return false }
  for i, _tgt := range p.ClusterNode {
    _src12 := other.ClusterNode[i]
    if _tgt != _src12 { return false }
  }
  if bytes.Compare(p.Key, other.Key) != 0 { return false }
  if bytes.Compare(p.Value, other.Value) != 0 { return false }
  if !p.Batch.Equals(other.Batch) { return false }
  if !p.Stat.Equals(other.Stat) { return false }
  if !p.LoadBean.Equals(other.LoadBean) { return false }
  if !p.TimeBean.Equals(other.TimeBean) { return false }
  if !p.SeqBean.Equals(other.SeqBean) { return false }
  if len(p.RunUUID) != len(other.RunUUID) { return false }
  for i, _tgt := range p.RunUUID {
    _src13 := other.RunUUID[i]
    if _tgt != _src13 { return false }
  }
  if p.SyncId != other.SyncId {
    if p.SyncId == nil || other.SyncId == nil {
      return false
    }
    if (*p.SyncId) != (*other.SyncId) { return false }
  }
  if bytes.Compare(p.Pbtime, other.Pbtime) != 0 { return false }
  return true
}

func (p *PonBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("PonBean(%+v)", *p)
}

func (p *PonBean) Validate() error {
  return nil
}
// Attributes:
//  - Seq
type LogSeqBean struct {
  Seq int64 `thrift:"seq,1,required" db:"seq" json:"seq"`
}

func NewLogSeqBean() *LogSeqBean {
  return &LogSeqBean{}
}


func (p *LogSeqBean) GetSeq() int64 {
  return p.Seq
}
func (p *LogSeqBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetSeq bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetSeq = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetSeq{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Seq is not set"));
  }
  return nil
}

func (p *LogSeqBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Seq = v
}
  return nil
}

func (p *LogSeqBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "LogSeqBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *LogSeqBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "seq", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:seq: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Seq)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.seq (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:seq: ", p), err) }
  return err
}

func (p *LogSeqBean) Equals(other *LogSeqBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Seq != other.Seq { return false }
  return true
}

func (p *LogSeqBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("LogSeqBean(%+v)", *p)
}

func (p *LogSeqBean) Validate() error {
  return nil
}
// Attributes:
//  - IsPass
//  - RunUUID
type TimeBean struct {
  IsPass bool `thrift:"isPass,1,required" db:"isPass" json:"isPass"`
  RunUUID int64 `thrift:"runUUID,2,required" db:"runUUID" json:"runUUID"`
}

func NewTimeBean() *TimeBean {
  return &TimeBean{}
}


func (p *TimeBean) GetIsPass() bool {
  return p.IsPass
}

func (p *TimeBean) GetRunUUID() int64 {
  return p.RunUUID
}
func (p *TimeBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetIsPass bool = false;
  var issetRunUUID bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetIsPass = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetRunUUID = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetIsPass{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsPass is not set"));
  }
  if !issetRunUUID{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RunUUID is not set"));
  }
  return nil
}

func (p *TimeBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.IsPass = v
}
  return nil
}

func (p *TimeBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.RunUUID = v
}
  return nil
}

func (p *TimeBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "TimeBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *TimeBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "isPass", thrift.BOOL, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:isPass: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.IsPass)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.isPass (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:isPass: ", p), err) }
  return err
}

func (p *TimeBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "runUUID", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:runUUID: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.RunUUID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.runUUID (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:runUUID: ", p), err) }
  return err
}

func (p *TimeBean) Equals(other *TimeBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.IsPass != other.IsPass { return false }
  if p.RunUUID != other.RunUUID { return false }
  return true
}

func (p *TimeBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TimeBean(%+v)", *p)
}

func (p *TimeBean) Validate() error {
  return nil
}
// Attributes:
//  - Stat
//  - Timenano
type Stat struct {
  Stat int8 `thrift:"stat,1,required" db:"stat" json:"stat"`
  Timenano int64 `thrift:"timenano,2,required" db:"timenano" json:"timenano"`
}

func NewStat() *Stat {
  return &Stat{}
}


func (p *Stat) GetStat() int8 {
  return p.Stat
}

func (p *Stat) GetTimenano() int64 {
  return p.Timenano
}
func (p *Stat) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetStat bool = false;
  var issetTimenano bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetStat = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetTimenano = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetStat{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Stat is not set"));
  }
  if !issetTimenano{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Timenano is not set"));
  }
  return nil
}

func (p *Stat)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.Stat = temp
}
  return nil
}

func (p *Stat)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Timenano = v
}
  return nil
}

func (p *Stat) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Stat"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *Stat) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "stat", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:stat: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Stat)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.stat (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:stat: ", p), err) }
  return err
}

func (p *Stat) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "timenano", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:timenano: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Timenano)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.timenano (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:timenano: ", p), err) }
  return err
}

func (p *Stat) Equals(other *Stat) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Stat != other.Stat { return false }
  if p.Timenano != other.Timenano { return false }
  return true
}

func (p *Stat) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("Stat(%+v)", *p)
}

func (p *Stat) Validate() error {
  return nil
}
// Attributes:
//  - S
//  - T
//  - Logm
type LoadBean struct {
  S int64 `thrift:"s,1,required" db:"s" json:"s"`
  T int64 `thrift:"t,2,required" db:"t" json:"t"`
  Logm map[int64]*LogKV `thrift:"logm,3" db:"logm" json:"logm,omitempty"`
}

func NewLoadBean() *LoadBean {
  return &LoadBean{}
}


func (p *LoadBean) GetS() int64 {
  return p.S
}

func (p *LoadBean) GetT() int64 {
  return p.T
}
var LoadBean_Logm_DEFAULT map[int64]*LogKV

func (p *LoadBean) GetLogm() map[int64]*LogKV {
  return p.Logm
}
func (p *LoadBean) IsSetLogm() bool {
  return p.Logm != nil
}

func (p *LoadBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetS bool = false;
  var issetT bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetS = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetT = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetS{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field S is not set"));
  }
  if !issetT{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field T is not set"));
  }
  return nil
}

func (p *LoadBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.S = v
}
  return nil
}

func (p *LoadBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.T = v
}
  return nil
}

func (p *LoadBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[int64]*LogKV, size)
  p.Logm =  tMap
  for i := 0; i < size; i ++ {
var _key14 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key14 = v
}
    _val15 := &LogKV{}
    if err := _val15.Read(ctx, iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _val15), err)
    }
    p.Logm[_key14] = _val15
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *LoadBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "LoadBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *LoadBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "s", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:s: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.S)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.s (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:s: ", p), err) }
  return err
}

func (p *LoadBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "t", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:t: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.T)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.t (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:t: ", p), err) }
  return err
}

func (p *LoadBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetLogm() {
    if err := oprot.WriteFieldBegin(ctx, "logm", thrift.MAP, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:logm: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.I64, thrift.STRUCT, len(p.Logm)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Logm {
      if err := oprot.WriteI64(ctx, int64(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := v.Write(ctx, oprot); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
      }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:logm: ", p), err) }
  }
  return err
}

func (p *LoadBean) Equals(other *LoadBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.S != other.S { return false }
  if p.T != other.T { return false }
  if len(p.Logm) != len(other.Logm) { return false }
  for k, _tgt := range p.Logm {
    _src16 := other.Logm[k]
    if !_tgt.Equals(_src16) { return false }
  }
  return true
}

func (p *LoadBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("LoadBean(%+v)", *p)
}

func (p *LoadBean) Validate() error {
  return nil
}
// Attributes:
//  - Filename
//  - Size
//  - Body
//  - StatKey
type FileBean struct {
  Filename string `thrift:"filename,1,required" db:"filename" json:"filename"`
  Size int64 `thrift:"size,2,required" db:"size" json:"size"`
  Body []byte `thrift:"body,3,required" db:"body" json:"body"`
  StatKey int64 `thrift:"statKey,4,required" db:"statKey" json:"statKey"`
}

func NewFileBean() *FileBean {
  return &FileBean{}
}


func (p *FileBean) GetFilename() string {
  return p.Filename
}

func (p *FileBean) GetSize() int64 {
  return p.Size
}

func (p *FileBean) GetBody() []byte {
  return p.Body
}

func (p *FileBean) GetStatKey() int64 {
  return p.StatKey
}
func (p *FileBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetFilename bool = false;
  var issetSize bool = false;
  var issetBody bool = false;
  var issetStatKey bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetFilename = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetSize = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetBody = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
        issetStatKey = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetFilename{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Filename is not set"));
  }
  if !issetSize{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Size is not set"));
  }
  if !issetBody{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Body is not set"));
  }
  if !issetStatKey{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field StatKey is not set"));
  }
  return nil
}

func (p *FileBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Filename = v
}
  return nil
}

func (p *FileBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Size = v
}
  return nil
}

func (p *FileBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Body = v
}
  return nil
}

func (p *FileBean)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.StatKey = v
}
  return nil
}

func (p *FileBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "FileBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *FileBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "filename", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:filename: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Filename)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.filename (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:filename: ", p), err) }
  return err
}

func (p *FileBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "size", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:size: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Size)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.size (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:size: ", p), err) }
  return err
}

func (p *FileBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "body", thrift.STRING, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:body: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.Body); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.body (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:body: ", p), err) }
  return err
}

func (p *FileBean) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "statKey", thrift.I64, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:statKey: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.StatKey)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.statKey (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:statKey: ", p), err) }
  return err
}

func (p *FileBean) Equals(other *FileBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Filename != other.Filename { return false }
  if p.Size != other.Size { return false }
  if bytes.Compare(p.Body, other.Body) != 0 { return false }
  if p.StatKey != other.StatKey { return false }
  return true
}

func (p *FileBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("FileBean(%+v)", *p)
}

func (p *FileBean) Validate() error {
  return nil
}
// Attributes:
//  - Type
//  - LastTime
//  - Txid
//  - FBean
//  - HasNext
//  - NotNull
type LogDataBean struct {
  Type int8 `thrift:"type,1,required" db:"type" json:"type"`
  LastTime int64 `thrift:"lastTime,2,required" db:"lastTime" json:"lastTime"`
  Txid int64 `thrift:"txid,3,required" db:"txid" json:"txid"`
  FBean *FileBean `thrift:"fBean,4" db:"fBean" json:"fBean,omitempty"`
  HasNext bool `thrift:"hasNext,5,required" db:"hasNext" json:"hasNext"`
  NotNull bool `thrift:"notNull,6,required" db:"notNull" json:"notNull"`
}

func NewLogDataBean() *LogDataBean {
  return &LogDataBean{}
}


func (p *LogDataBean) GetType() int8 {
  return p.Type
}

func (p *LogDataBean) GetLastTime() int64 {
  return p.LastTime
}

func (p *LogDataBean) GetTxid() int64 {
  return p.Txid
}
var LogDataBean_FBean_DEFAULT *FileBean
func (p *LogDataBean) GetFBean() *FileBean {
  if !p.IsSetFBean() {
    return LogDataBean_FBean_DEFAULT
  }
return p.FBean
}

func (p *LogDataBean) GetHasNext() bool {
  return p.HasNext
}

func (p *LogDataBean) GetNotNull() bool {
  return p.NotNull
}
func (p *LogDataBean) IsSetFBean() bool {
  return p.FBean != nil
}

func (p *LogDataBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetType bool = false;
  var issetLastTime bool = false;
  var issetTxid bool = false;
  var issetHasNext bool = false;
  var issetNotNull bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetLastTime = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetTxid = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
        issetHasNext = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField6(ctx, iprot); err != nil {
          return err
        }
        issetNotNull = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Type is not set"));
  }
  if !issetLastTime{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field LastTime is not set"));
  }
  if !issetTxid{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Txid is not set"));
  }
  if !issetHasNext{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field HasNext is not set"));
  }
  if !issetNotNull{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NotNull is not set"));
  }
  return nil
}

func (p *LogDataBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.Type = temp
}
  return nil
}

func (p *LogDataBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.LastTime = v
}
  return nil
}

func (p *LogDataBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Txid = v
}
  return nil
}

func (p *LogDataBean)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  p.FBean = &FileBean{}
  if err := p.FBean.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.FBean), err)
  }
  return nil
}

func (p *LogDataBean)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.HasNext = v
}
  return nil
}

func (p *LogDataBean)  ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  p.NotNull = v
}
  return nil
}

func (p *LogDataBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "LogDataBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
    if err := p.writeField6(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *LogDataBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "type", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:type: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Type)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.type (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:type: ", p), err) }
  return err
}

func (p *LogDataBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "lastTime", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:lastTime: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.LastTime)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.lastTime (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:lastTime: ", p), err) }
  return err
}

func (p *LogDataBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "txid", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:txid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Txid)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.txid (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:txid: ", p), err) }
  return err
}

func (p *LogDataBean) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetFBean() {
    if err := oprot.WriteFieldBegin(ctx, "fBean", thrift.STRUCT, 4); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:fBean: ", p), err) }
    if err := p.FBean.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.FBean), err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 4:fBean: ", p), err) }
  }
  return err
}

func (p *LogDataBean) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "hasNext", thrift.BOOL, 5); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:hasNext: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.HasNext)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.hasNext (5) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 5:hasNext: ", p), err) }
  return err
}

func (p *LogDataBean) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "notNull", thrift.BOOL, 6); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:notNull: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.NotNull)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.notNull (6) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 6:notNull: ", p), err) }
  return err
}

func (p *LogDataBean) Equals(other *LogDataBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Type != other.Type { return false }
  if p.LastTime != other.LastTime { return false }
  if p.Txid != other.Txid { return false }
  if !p.FBean.Equals(other.FBean) { return false }
  if p.HasNext != other.HasNext { return false }
  if p.NotNull != other.NotNull { return false }
  return true
}

func (p *LogDataBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("LogDataBean(%+v)", *p)
}

func (p *LogDataBean) Validate() error {
  return nil
}
// Attributes:
//  - DBMode
//  - STORENODENUM
//  - REMOVEUUID
type SysBean struct {
  DBMode int8 `thrift:"DBMode,1,required" db:"DBMode" json:"DBMode"`
  STORENODENUM int32 `thrift:"STORENODENUM,2,required" db:"STORENODENUM" json:"STORENODENUM"`
  REMOVEUUID int64 `thrift:"REMOVEUUID,3,required" db:"REMOVEUUID" json:"REMOVEUUID"`
}

func NewSysBean() *SysBean {
  return &SysBean{}
}


func (p *SysBean) GetDBMode() int8 {
  return p.DBMode
}

func (p *SysBean) GetSTORENODENUM() int32 {
  return p.STORENODENUM
}

func (p *SysBean) GetREMOVEUUID() int64 {
  return p.REMOVEUUID
}
func (p *SysBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetDBMode bool = false;
  var issetSTORENODENUM bool = false;
  var issetREMOVEUUID bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetDBMode = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetSTORENODENUM = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetREMOVEUUID = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetDBMode{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field DBMode is not set"));
  }
  if !issetSTORENODENUM{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field STORENODENUM is not set"));
  }
  if !issetREMOVEUUID{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field REMOVEUUID is not set"));
  }
  return nil
}

func (p *SysBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.DBMode = temp
}
  return nil
}

func (p *SysBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.STORENODENUM = v
}
  return nil
}

func (p *SysBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.REMOVEUUID = v
}
  return nil
}

func (p *SysBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SysBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *SysBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "DBMode", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:DBMode: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.DBMode)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.DBMode (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:DBMode: ", p), err) }
  return err
}

func (p *SysBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "STORENODENUM", thrift.I32, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:STORENODENUM: ", p), err) }
  if err := oprot.WriteI32(ctx, int32(p.STORENODENUM)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.STORENODENUM (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:STORENODENUM: ", p), err) }
  return err
}

func (p *SysBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "REMOVEUUID", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:REMOVEUUID: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.REMOVEUUID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.REMOVEUUID (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:REMOVEUUID: ", p), err) }
  return err
}

func (p *SysBean) Equals(other *SysBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.DBMode != other.DBMode { return false }
  if p.STORENODENUM != other.STORENODENUM { return false }
  if p.REMOVEUUID != other.REMOVEUUID { return false }
  return true
}

func (p *SysBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("SysBean(%+v)", *p)
}

func (p *SysBean) Validate() error {
  return nil
}
// Attributes:
//  - ReqType
//  - Str
//  - TokenFlag
//  - Status
//  - Token
//  - SyncId
//  - Bs
//  - Extra
//  - List
//  - BackOk
//  - ForceBack
type TokenTrans struct {
  ReqType int8 `thrift:"ReqType,1,required" db:"ReqType" json:"ReqType"`
  Str string `thrift:"Str,2,required" db:"Str" json:"Str"`
  TokenFlag int64 `thrift:"TokenFlag,3,required" db:"TokenFlag" json:"TokenFlag"`
  Status int8 `thrift:"Status,4,required" db:"Status" json:"Status"`
  Token string `thrift:"Token,5,required" db:"Token" json:"Token"`
  SyncId *int64 `thrift:"SyncId,6" db:"SyncId" json:"SyncId,omitempty"`
  Bs []byte `thrift:"bs,7" db:"bs" json:"bs,omitempty"`
  Extra map[string][]byte `thrift:"Extra,8" db:"Extra" json:"Extra,omitempty"`
  List []*TokenTrans `thrift:"List,9" db:"List" json:"List,omitempty"`
  BackOk *bool `thrift:"BackOk,10" db:"BackOk" json:"BackOk,omitempty"`
  ForceBack *bool `thrift:"ForceBack,11" db:"ForceBack" json:"ForceBack,omitempty"`
}

func NewTokenTrans() *TokenTrans {
  return &TokenTrans{}
}


func (p *TokenTrans) GetReqType() int8 {
  return p.ReqType
}

func (p *TokenTrans) GetStr() string {
  return p.Str
}

func (p *TokenTrans) GetTokenFlag() int64 {
  return p.TokenFlag
}

func (p *TokenTrans) GetStatus() int8 {
  return p.Status
}

func (p *TokenTrans) GetToken() string {
  return p.Token
}
var TokenTrans_SyncId_DEFAULT int64
func (p *TokenTrans) GetSyncId() int64 {
  if !p.IsSetSyncId() {
    return TokenTrans_SyncId_DEFAULT
  }
return *p.SyncId
}
var TokenTrans_Bs_DEFAULT []byte

func (p *TokenTrans) GetBs() []byte {
  return p.Bs
}
var TokenTrans_Extra_DEFAULT map[string][]byte

func (p *TokenTrans) GetExtra() map[string][]byte {
  return p.Extra
}
var TokenTrans_List_DEFAULT []*TokenTrans

func (p *TokenTrans) GetList() []*TokenTrans {
  return p.List
}
var TokenTrans_BackOk_DEFAULT bool
func (p *TokenTrans) GetBackOk() bool {
  if !p.IsSetBackOk() {
    return TokenTrans_BackOk_DEFAULT
  }
return *p.BackOk
}
var TokenTrans_ForceBack_DEFAULT bool
func (p *TokenTrans) GetForceBack() bool {
  if !p.IsSetForceBack() {
    return TokenTrans_ForceBack_DEFAULT
  }
return *p.ForceBack
}
func (p *TokenTrans) IsSetSyncId() bool {
  return p.SyncId != nil
}

func (p *TokenTrans) IsSetBs() bool {
  return p.Bs != nil
}

func (p *TokenTrans) IsSetExtra() bool {
  return p.Extra != nil
}

func (p *TokenTrans) IsSetList() bool {
  return p.List != nil
}

func (p *TokenTrans) IsSetBackOk() bool {
  return p.BackOk != nil
}

func (p *TokenTrans) IsSetForceBack() bool {
  return p.ForceBack != nil
}

func (p *TokenTrans) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetReqType bool = false;
  var issetStr bool = false;
  var issetTokenFlag bool = false;
  var issetStatus bool = false;
  var issetToken bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetReqType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetStr = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetTokenFlag = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
        issetStatus = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
        issetToken = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField6(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 7:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField7(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 8:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField8(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 9:
      if fieldTypeId == thrift.LIST {
        if err := p.ReadField9(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 10:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField10(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 11:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField11(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetReqType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field ReqType is not set"));
  }
  if !issetStr{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Str is not set"));
  }
  if !issetTokenFlag{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TokenFlag is not set"));
  }
  if !issetStatus{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Status is not set"));
  }
  if !issetToken{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Token is not set"));
  }
  return nil
}

func (p *TokenTrans)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.ReqType = temp
}
  return nil
}

func (p *TokenTrans)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Str = v
}
  return nil
}

func (p *TokenTrans)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.TokenFlag = v
}
  return nil
}

func (p *TokenTrans)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  temp := int8(v)
  p.Status = temp
}
  return nil
}

func (p *TokenTrans)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.Token = v
}
  return nil
}

func (p *TokenTrans)  ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  p.SyncId = &v
}
  return nil
}

func (p *TokenTrans)  ReadField7(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 7: ", err)
} else {
  p.Bs = v
}
  return nil
}

func (p *TokenTrans)  ReadField8(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string][]byte, size)
  p.Extra =  tMap
  for i := 0; i < size; i ++ {
var _key17 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key17 = v
}
var _val18 []byte
    if v, err := iprot.ReadBinary(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val18 = v
}
    p.Extra[_key17] = _val18
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *TokenTrans)  ReadField9(ctx context.Context, iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]*TokenTrans, 0, size)
  p.List =  tSlice
  for i := 0; i < size; i ++ {
    _elem19 := &TokenTrans{}
    if err := _elem19.Read(ctx, iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem19), err)
    }
    p.List = append(p.List, _elem19)
  }
  if err := iprot.ReadListEnd(ctx); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *TokenTrans)  ReadField10(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 10: ", err)
} else {
  p.BackOk = &v
}
  return nil
}

func (p *TokenTrans)  ReadField11(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 11: ", err)
} else {
  p.ForceBack = &v
}
  return nil
}

func (p *TokenTrans) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "TokenTrans"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
    if err := p.writeField6(ctx, oprot); err != nil { return err }
    if err := p.writeField7(ctx, oprot); err != nil { return err }
    if err := p.writeField8(ctx, oprot); err != nil { return err }
    if err := p.writeField9(ctx, oprot); err != nil { return err }
    if err := p.writeField10(ctx, oprot); err != nil { return err }
    if err := p.writeField11(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *TokenTrans) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ReqType", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ReqType: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.ReqType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ReqType (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ReqType: ", p), err) }
  return err
}

func (p *TokenTrans) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "Str", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:Str: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Str)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.Str (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:Str: ", p), err) }
  return err
}

func (p *TokenTrans) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "TokenFlag", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:TokenFlag: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.TokenFlag)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.TokenFlag (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:TokenFlag: ", p), err) }
  return err
}

func (p *TokenTrans) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "Status", thrift.BYTE, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:Status: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Status)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.Status (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:Status: ", p), err) }
  return err
}

func (p *TokenTrans) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "Token", thrift.STRING, 5); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:Token: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Token)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.Token (5) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 5:Token: ", p), err) }
  return err
}

func (p *TokenTrans) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetSyncId() {
    if err := oprot.WriteFieldBegin(ctx, "SyncId", thrift.I64, 6); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:SyncId: ", p), err) }
    if err := oprot.WriteI64(ctx, int64(*p.SyncId)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.SyncId (6) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 6:SyncId: ", p), err) }
  }
  return err
}

func (p *TokenTrans) writeField7(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBs() {
    if err := oprot.WriteFieldBegin(ctx, "bs", thrift.STRING, 7); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:bs: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Bs); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.bs (7) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 7:bs: ", p), err) }
  }
  return err
}

func (p *TokenTrans) writeField8(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetExtra() {
    if err := oprot.WriteFieldBegin(ctx, "Extra", thrift.MAP, 8); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:Extra: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Extra)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Extra {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteBinary(ctx, v); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 8:Extra: ", p), err) }
  }
  return err
}

func (p *TokenTrans) writeField9(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetList() {
    if err := oprot.WriteFieldBegin(ctx, "List", thrift.LIST, 9); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 9:List: ", p), err) }
    if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.List)); err != nil {
      return thrift.PrependError("error writing list begin: ", err)
    }
    for _, v := range p.List {
      if err := v.Write(ctx, oprot); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
      }
    }
    if err := oprot.WriteListEnd(ctx); err != nil {
      return thrift.PrependError("error writing list end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 9:List: ", p), err) }
  }
  return err
}

func (p *TokenTrans) writeField10(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBackOk() {
    if err := oprot.WriteFieldBegin(ctx, "BackOk", thrift.BOOL, 10); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:BackOk: ", p), err) }
    if err := oprot.WriteBool(ctx, bool(*p.BackOk)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.BackOk (10) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 10:BackOk: ", p), err) }
  }
  return err
}

func (p *TokenTrans) writeField11(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetForceBack() {
    if err := oprot.WriteFieldBegin(ctx, "ForceBack", thrift.BOOL, 11); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 11:ForceBack: ", p), err) }
    if err := oprot.WriteBool(ctx, bool(*p.ForceBack)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.ForceBack (11) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 11:ForceBack: ", p), err) }
  }
  return err
}

func (p *TokenTrans) Equals(other *TokenTrans) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.ReqType != other.ReqType { return false }
  if p.Str != other.Str { return false }
  if p.TokenFlag != other.TokenFlag { return false }
  if p.Status != other.Status { return false }
  if p.Token != other.Token { return false }
  if p.SyncId != other.SyncId {
    if p.SyncId == nil || other.SyncId == nil {
      return false
    }
    if (*p.SyncId) != (*other.SyncId) { return false }
  }
  if bytes.Compare(p.Bs, other.Bs) != 0 { return false }
  if len(p.Extra) != len(other.Extra) { return false }
  for k, _tgt := range p.Extra {
    _src20 := other.Extra[k]
    if bytes.Compare(_tgt, _src20) != 0 { return false }
  }
  if len(p.List) != len(other.List) { return false }
  for i, _tgt := range p.List {
    _src21 := other.List[i]
    if !_tgt.Equals(_src21) { return false }
  }
  if p.BackOk != other.BackOk {
    if p.BackOk == nil || other.BackOk == nil {
      return false
    }
    if (*p.BackOk) != (*other.BackOk) { return false }
  }
  if p.ForceBack != other.ForceBack {
    if p.ForceBack == nil || other.ForceBack == nil {
      return false
    }
    if (*p.ForceBack) != (*other.ForceBack) { return false }
  }
  return true
}

func (p *TokenTrans) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TokenTrans(%+v)", *p)
}

func (p *TokenTrans) Validate() error {
  return nil
}
type Itnet interface {
  // Parameters:
  //  - Ping
  Ping(ctx context.Context, ping int64) (_err error)
  // Parameters:
  //  - PongBs
  Pong(ctx context.Context, pongBs []byte) (_err error)
  // Parameters:
  //  - AuthKey
  Auth(ctx context.Context, authKey []byte) (_err error)
  // Parameters:
  //  - AuthKey
  Auth2(ctx context.Context, authKey []byte) (_err error)
  // Parameters:
  //  - Pblist
  //  - ID
  PonMerge(ctx context.Context, pblist []*PonBean, id int64) (_err error)
  // Parameters:
  //  - PonBeanBytes
  //  - ID
  Pon(ctx context.Context, PonBeanBytes []byte, id int64) (_err error)
  // Parameters:
  //  - Pb
  //  - ID
  Pon2(ctx context.Context, pb *PonBean, id int64) (_err error)
  // Parameters:
  //  - PonBeanBytes
  //  - ID
  //  - Ack
  Pon3(ctx context.Context, PonBeanBytes []byte, id int64, ack bool) (_err error)
  // Parameters:
  //  - Pretimenano
  //  - Timenano
  //  - Num
  //  - Ir
  Time(ctx context.Context, pretimenano int64, timenano int64, num int16, ir bool) (_err error)
  // Parameters:
  //  - Node
  //  - Ir
  SyncNode(ctx context.Context, node *Node, ir bool) (_err error)
  // Parameters:
  //  - SyncId
  //  - Result_
  SyncTx(ctx context.Context, syncId int64, result int8) (_err error)
  // Parameters:
  //  - SyncList
  SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error)
  // Parameters:
  //  - SyncId
  //  - Txid
  CommitTx(ctx context.Context, syncId int64, txid int64) (_err error)
  // Parameters:
  //  - SyncId
  //  - Txid
  //  - Commit
  CommitTx2(ctx context.Context, syncId int64, txid int64, commit bool) (_err error)
  // Parameters:
  //  - SyncId
  //  - MqType
  //  - Bs
  PubMq(ctx context.Context, syncId int64, mqType int8, bs []byte) (_err error)
  // Parameters:
  //  - SyncId
  //  - Ldb
  PullData(ctx context.Context, syncId int64, ldb *LogDataBean) (_err error)
  // Parameters:
  //  - SyncId
  //  - Sbean
  ReInit(ctx context.Context, syncId int64, sbean *SysBean) (_err error)
  // Parameters:
  //  - SyncId
  //  - ParamData
  //  - PType
  //  - Ctype
  ProxyCall(ctx context.Context, syncId int64, paramData []byte, pType int8, ctype int8) (_err error)
  // Parameters:
  //  - SyncId
  //  - Tt
  //  - Ack
  BroadToken(ctx context.Context, syncId int64, tt *TokenTrans, ack bool) (_err error)
}

type ItnetClient struct {
  c thrift.TClient
  meta thrift.ResponseMeta
}

func NewItnetClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ItnetClient {
  return &ItnetClient{
    c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
  }
}

func NewItnetClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ItnetClient {
  return &ItnetClient{
    c: thrift.NewTStandardClient(iprot, oprot),
  }
}

func NewItnetClient(c thrift.TClient) *ItnetClient {
  return &ItnetClient{
    c: c,
  }
}

func (p *ItnetClient) Client_() thrift.TClient {
  return p.c
}

func (p *ItnetClient) LastResponseMeta_() thrift.ResponseMeta {
  return p.meta
}

func (p *ItnetClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
  p.meta = meta
}

// Parameters:
//  - Ping
func (p *ItnetClient) Ping(ctx context.Context, ping int64) (_err error) {
  var _args22 ItnetPingArgs
  _args22.Ping = ping
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Ping", &_args22, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - PongBs
func (p *ItnetClient) Pong(ctx context.Context, pongBs []byte) (_err error) {
  var _args23 ItnetPongArgs
  _args23.PongBs = pongBs
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Pong", &_args23, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - AuthKey
func (p *ItnetClient) Auth(ctx context.Context, authKey []byte) (_err error) {
  var _args24 ItnetAuthArgs
  _args24.AuthKey = authKey
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Auth", &_args24, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - AuthKey
func (p *ItnetClient) Auth2(ctx context.Context, authKey []byte) (_err error) {
  var _args25 ItnetAuth2Args
  _args25.AuthKey = authKey
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Auth2", &_args25, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Pblist
//  - ID
func (p *ItnetClient) PonMerge(ctx context.Context, pblist []*PonBean, id int64) (_err error) {
  var _args26 ItnetPonMergeArgs
  _args26.Pblist = pblist
  _args26.ID = id
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "PonMerge", &_args26, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - PonBeanBytes
//  - ID
func (p *ItnetClient) Pon(ctx context.Context, PonBeanBytes []byte, id int64) (_err error) {
  var _args27 ItnetPonArgs
  _args27.PonBeanBytes = PonBeanBytes
  _args27.ID = id
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Pon", &_args27, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Pb
//  - ID
func (p *ItnetClient) Pon2(ctx context.Context, pb *PonBean, id int64) (_err error) {
  var _args28 ItnetPon2Args
  _args28.Pb = pb
  _args28.ID = id
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Pon2", &_args28, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - PonBeanBytes
//  - ID
//  - Ack
func (p *ItnetClient) Pon3(ctx context.Context, PonBeanBytes []byte, id int64, ack bool) (_err error) {
  var _args29 ItnetPon3Args
  _args29.PonBeanBytes = PonBeanBytes
  _args29.ID = id
  _args29.Ack = ack
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Pon3", &_args29, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Pretimenano
//  - Timenano
//  - Num
//  - Ir
func (p *ItnetClient) Time(ctx context.Context, pretimenano int64, timenano int64, num int16, ir bool) (_err error) {
  var _args30 ItnetTimeArgs
  _args30.Pretimenano = pretimenano
  _args30.Timenano = timenano
  _args30.Num = num
  _args30.Ir = ir
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Time", &_args30, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Node
//  - Ir
func (p *ItnetClient) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
  var _args31 ItnetSyncNodeArgs
  _args31.Node = node
  _args31.Ir = ir
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "SyncNode", &_args31, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - Result_
func (p *ItnetClient) SyncTx(ctx context.Context, syncId int64, result int8) (_err error) {
  var _args32 ItnetSyncTxArgs
  _args32.SyncId = syncId
  _args32.Result_ = result
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "SyncTx", &_args32, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncList
func (p *ItnetClient) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
  var _args33 ItnetSyncTxMergeArgs
  _args33.SyncList = syncList
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "SyncTxMerge", &_args33, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - Txid
func (p *ItnetClient) CommitTx(ctx context.Context, syncId int64, txid int64) (_err error) {
  var _args34 ItnetCommitTxArgs
  _args34.SyncId = syncId
  _args34.Txid = txid
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "CommitTx", &_args34, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - Txid
//  - Commit
func (p *ItnetClient) CommitTx2(ctx context.Context, syncId int64, txid int64, commit bool) (_err error) {
  var _args35 ItnetCommitTx2Args
  _args35.SyncId = syncId
  _args35.Txid = txid
  _args35.Commit = commit
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "CommitTx2", &_args35, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - MqType
//  - Bs
func (p *ItnetClient) PubMq(ctx context.Context, syncId int64, mqType int8, bs []byte) (_err error) {
  var _args36 ItnetPubMqArgs
  _args36.SyncId = syncId
  _args36.MqType = mqType
  _args36.Bs = bs
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "PubMq", &_args36, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - Ldb
func (p *ItnetClient) PullData(ctx context.Context, syncId int64, ldb *LogDataBean) (_err error) {
  var _args37 ItnetPullDataArgs
  _args37.SyncId = syncId
  _args37.Ldb = ldb
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "PullData", &_args37, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - Sbean
func (p *ItnetClient) ReInit(ctx context.Context, syncId int64, sbean *SysBean) (_err error) {
  var _args38 ItnetReInitArgs
  _args38.SyncId = syncId
  _args38.Sbean = sbean
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "ReInit", &_args38, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - ParamData
//  - PType
//  - Ctype
func (p *ItnetClient) ProxyCall(ctx context.Context, syncId int64, paramData []byte, pType int8, ctype int8) (_err error) {
  var _args39 ItnetProxyCallArgs
  _args39.SyncId = syncId
  _args39.ParamData = paramData
  _args39.PType = pType
  _args39.Ctype = ctype
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "ProxyCall", &_args39, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncId
//  - Tt
//  - Ack
func (p *ItnetClient) BroadToken(ctx context.Context, syncId int64, tt *TokenTrans, ack bool) (_err error) {
  var _args40 ItnetBroadTokenArgs
  _args40.SyncId = syncId
  _args40.Tt = tt
  _args40.Ack = ack
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "BroadToken", &_args40, nil); err != nil {
    return err
  }
  return nil
}

type ItnetProcessor struct {
  processorMap map[string]thrift.TProcessorFunction
  handler Itnet
}

func (p *ItnetProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
  p.processorMap[key] = processor
}

func (p *ItnetProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
  processor, ok = p.processorMap[key]
  return processor, ok
}

func (p *ItnetProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
  return p.processorMap
}

func NewItnetProcessor(handler Itnet) *ItnetProcessor {

  self41 := &ItnetProcessor{handler:handler, processorMap:make(map[string]thrift.TProcessorFunction)}
  self41.processorMap["Ping"] = &itnetProcessorPing{handler:handler}
  self41.processorMap["Pong"] = &itnetProcessorPong{handler:handler}
  self41.processorMap["Auth"] = &itnetProcessorAuth{handler:handler}
  self41.processorMap["Auth2"] = &itnetProcessorAuth2{handler:handler}
  self41.processorMap["PonMerge"] = &itnetProcessorPonMerge{handler:handler}
  self41.processorMap["Pon"] = &itnetProcessorPon{handler:handler}
  self41.processorMap["Pon2"] = &itnetProcessorPon2{handler:handler}
  self41.processorMap["Pon3"] = &itnetProcessorPon3{handler:handler}
  self41.processorMap["Time"] = &itnetProcessorTime{handler:handler}
  self41.processorMap["SyncNode"] = &itnetProcessorSyncNode{handler:handler}
  self41.processorMap["SyncTx"] = &itnetProcessorSyncTx{handler:handler}
  self41.processorMap["SyncTxMerge"] = &itnetProcessorSyncTxMerge{handler:handler}
  self41.processorMap["CommitTx"] = &itnetProcessorCommitTx{handler:handler}
  self41.processorMap["CommitTx2"] = &itnetProcessorCommitTx2{handler:handler}
  self41.processorMap["PubMq"] = &itnetProcessorPubMq{handler:handler}
  self41.processorMap["PullData"] = &itnetProcessorPullData{handler:handler}
  self41.processorMap["ReInit"] = &itnetProcessorReInit{handler:handler}
  self41.processorMap["ProxyCall"] = &itnetProcessorProxyCall{handler:handler}
  self41.processorMap["BroadToken"] = &itnetProcessorBroadToken{handler:handler}
return self41
}

func (p *ItnetProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
  if err2 != nil { return false, thrift.WrapTException(err2) }
  if processor, ok := p.GetProcessorFunction(name); ok {
    return processor.Process(ctx, seqId, iprot, oprot)
  }
  iprot.Skip(ctx, thrift.STRUCT)
  iprot.ReadMessageEnd(ctx)
  x42 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function " + name)
  oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
  x42.Write(ctx, oprot)
  oprot.WriteMessageEnd(ctx)
  oprot.Flush(ctx)
  return false, x42

}

type itnetProcessorPing struct {
  handler Itnet
}

func (p *itnetProcessorPing) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPingArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Ping(ctx, args.Ping); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPong struct {
  handler Itnet
}

func (p *itnetProcessorPong) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPongArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Pong(ctx, args.PongBs); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorAuth struct {
  handler Itnet
}

func (p *itnetProcessorAuth) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetAuthArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Auth(ctx, args.AuthKey); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorAuth2 struct {
  handler Itnet
}

func (p *itnetProcessorAuth2) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetAuth2Args{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Auth2(ctx, args.AuthKey); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPonMerge struct {
  handler Itnet
}

func (p *itnetProcessorPonMerge) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPonMergeArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.PonMerge(ctx, args.Pblist, args.ID); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPon struct {
  handler Itnet
}

func (p *itnetProcessorPon) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPonArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Pon(ctx, args.PonBeanBytes, args.ID); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPon2 struct {
  handler Itnet
}

func (p *itnetProcessorPon2) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPon2Args{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Pon2(ctx, args.Pb, args.ID); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPon3 struct {
  handler Itnet
}

func (p *itnetProcessorPon3) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPon3Args{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Pon3(ctx, args.PonBeanBytes, args.ID, args.Ack); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorTime struct {
  handler Itnet
}

func (p *itnetProcessorTime) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetTimeArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Time(ctx, args.Pretimenano, args.Timenano, args.Num, args.Ir); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorSyncNode struct {
  handler Itnet
}

func (p *itnetProcessorSyncNode) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetSyncNodeArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.SyncNode(ctx, args.Node, args.Ir); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorSyncTx struct {
  handler Itnet
}

func (p *itnetProcessorSyncTx) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetSyncTxArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.SyncTx(ctx, args.SyncId, args.Result_); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorSyncTxMerge struct {
  handler Itnet
}

func (p *itnetProcessorSyncTxMerge) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetSyncTxMergeArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.SyncTxMerge(ctx, args.SyncList); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorCommitTx struct {
  handler Itnet
}

func (p *itnetProcessorCommitTx) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetCommitTxArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.CommitTx(ctx, args.SyncId, args.Txid); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorCommitTx2 struct {
  handler Itnet
}

func (p *itnetProcessorCommitTx2) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetCommitTx2Args{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.CommitTx2(ctx, args.SyncId, args.Txid, args.Commit); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPubMq struct {
  handler Itnet
}

func (p *itnetProcessorPubMq) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPubMqArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.PubMq(ctx, args.SyncId, args.MqType, args.Bs); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPullData struct {
  handler Itnet
}

func (p *itnetProcessorPullData) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPullDataArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.PullData(ctx, args.SyncId, args.Ldb); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorReInit struct {
  handler Itnet
}

func (p *itnetProcessorReInit) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetReInitArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.ReInit(ctx, args.SyncId, args.Sbean); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorProxyCall struct {
  handler Itnet
}

func (p *itnetProcessorProxyCall) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetProxyCallArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.ProxyCall(ctx, args.SyncId, args.ParamData, args.PType, args.Ctype); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorBroadToken struct {
  handler Itnet
}

func (p *itnetProcessorBroadToken) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetBroadTokenArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.BroadToken(ctx, args.SyncId, args.Tt, args.Ack); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}


// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Ping
type ItnetPingArgs struct {
  Ping int64 `thrift:"ping,1" db:"ping" json:"ping"`
}

func NewItnetPingArgs() *ItnetPingArgs {
  return &ItnetPingArgs{}
}


func (p *ItnetPingArgs) GetPing() int64 {
  return p.Ping
}
func (p *ItnetPingArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPingArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Ping = v
}
  return nil
}

func (p *ItnetPingArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Ping_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPingArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ping", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ping: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Ping)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ping (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ping: ", p), err) }
  return err
}

func (p *ItnetPingArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPingArgs(%+v)", *p)
}

// Attributes:
//  - PongBs
type ItnetPongArgs struct {
  PongBs []byte `thrift:"pongBs,1" db:"pongBs" json:"pongBs"`
}

func NewItnetPongArgs() *ItnetPongArgs {
  return &ItnetPongArgs{}
}


func (p *ItnetPongArgs) GetPongBs() []byte {
  return p.PongBs
}
func (p *ItnetPongArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPongArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.PongBs = v
}
  return nil
}

func (p *ItnetPongArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Pong_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPongArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pongBs", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:pongBs: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.PongBs); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.pongBs (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:pongBs: ", p), err) }
  return err
}

func (p *ItnetPongArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPongArgs(%+v)", *p)
}

// Attributes:
//  - AuthKey
type ItnetAuthArgs struct {
  AuthKey []byte `thrift:"authKey,1" db:"authKey" json:"authKey"`
}

func NewItnetAuthArgs() *ItnetAuthArgs {
  return &ItnetAuthArgs{}
}


func (p *ItnetAuthArgs) GetAuthKey() []byte {
  return p.AuthKey
}
func (p *ItnetAuthArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetAuthArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.AuthKey = v
}
  return nil
}

func (p *ItnetAuthArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Auth_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetAuthArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "authKey", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:authKey: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.AuthKey); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.authKey (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:authKey: ", p), err) }
  return err
}

func (p *ItnetAuthArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetAuthArgs(%+v)", *p)
}

// Attributes:
//  - AuthKey
type ItnetAuth2Args struct {
  AuthKey []byte `thrift:"authKey,1" db:"authKey" json:"authKey"`
}

func NewItnetAuth2Args() *ItnetAuth2Args {
  return &ItnetAuth2Args{}
}


func (p *ItnetAuth2Args) GetAuthKey() []byte {
  return p.AuthKey
}
func (p *ItnetAuth2Args) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetAuth2Args)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.AuthKey = v
}
  return nil
}

func (p *ItnetAuth2Args) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Auth2_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetAuth2Args) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "authKey", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:authKey: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.AuthKey); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.authKey (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:authKey: ", p), err) }
  return err
}

func (p *ItnetAuth2Args) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetAuth2Args(%+v)", *p)
}

// Attributes:
//  - Pblist
//  - ID
type ItnetPonMergeArgs struct {
  Pblist []*PonBean `thrift:"pblist,1" db:"pblist" json:"pblist"`
  ID int64 `thrift:"id,2" db:"id" json:"id"`
}

func NewItnetPonMergeArgs() *ItnetPonMergeArgs {
  return &ItnetPonMergeArgs{}
}


func (p *ItnetPonMergeArgs) GetPblist() []*PonBean {
  return p.Pblist
}

func (p *ItnetPonMergeArgs) GetID() int64 {
  return p.ID
}
func (p *ItnetPonMergeArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.LIST {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPonMergeArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]*PonBean, 0, size)
  p.Pblist =  tSlice
  for i := 0; i < size; i ++ {
    _elem43 := &PonBean{}
    if err := _elem43.Read(ctx, iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem43), err)
    }
    p.Pblist = append(p.Pblist, _elem43)
  }
  if err := iprot.ReadListEnd(ctx); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *ItnetPonMergeArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.ID = v
}
  return nil
}

func (p *ItnetPonMergeArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "PonMerge_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPonMergeArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pblist", thrift.LIST, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:pblist: ", p), err) }
  if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Pblist)); err != nil {
    return thrift.PrependError("error writing list begin: ", err)
  }
  for _, v := range p.Pblist {
    if err := v.Write(ctx, oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
    }
  }
  if err := oprot.WriteListEnd(ctx); err != nil {
    return thrift.PrependError("error writing list end: ", err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:pblist: ", p), err) }
  return err
}

func (p *ItnetPonMergeArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "id", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:id: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.ID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.id (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:id: ", p), err) }
  return err
}

func (p *ItnetPonMergeArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPonMergeArgs(%+v)", *p)
}

// Attributes:
//  - PonBeanBytes
//  - ID
type ItnetPonArgs struct {
  PonBeanBytes []byte `thrift:"PonBeanBytes,1" db:"PonBeanBytes" json:"PonBeanBytes"`
  ID int64 `thrift:"id,2" db:"id" json:"id"`
}

func NewItnetPonArgs() *ItnetPonArgs {
  return &ItnetPonArgs{}
}


func (p *ItnetPonArgs) GetPonBeanBytes() []byte {
  return p.PonBeanBytes
}

func (p *ItnetPonArgs) GetID() int64 {
  return p.ID
}
func (p *ItnetPonArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPonArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.PonBeanBytes = v
}
  return nil
}

func (p *ItnetPonArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.ID = v
}
  return nil
}

func (p *ItnetPonArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Pon_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPonArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "PonBeanBytes", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:PonBeanBytes: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.PonBeanBytes); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.PonBeanBytes (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:PonBeanBytes: ", p), err) }
  return err
}

func (p *ItnetPonArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "id", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:id: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.ID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.id (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:id: ", p), err) }
  return err
}

func (p *ItnetPonArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPonArgs(%+v)", *p)
}

// Attributes:
//  - Pb
//  - ID
type ItnetPon2Args struct {
  Pb *PonBean `thrift:"pb,1" db:"pb" json:"pb"`
  ID int64 `thrift:"id,2" db:"id" json:"id"`
}

func NewItnetPon2Args() *ItnetPon2Args {
  return &ItnetPon2Args{}
}

var ItnetPon2Args_Pb_DEFAULT *PonBean
func (p *ItnetPon2Args) GetPb() *PonBean {
  if !p.IsSetPb() {
    return ItnetPon2Args_Pb_DEFAULT
  }
return p.Pb
}

func (p *ItnetPon2Args) GetID() int64 {
  return p.ID
}
func (p *ItnetPon2Args) IsSetPb() bool {
  return p.Pb != nil
}

func (p *ItnetPon2Args) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPon2Args)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  p.Pb = &PonBean{}
  if err := p.Pb.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Pb), err)
  }
  return nil
}

func (p *ItnetPon2Args)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.ID = v
}
  return nil
}

func (p *ItnetPon2Args) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Pon2_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPon2Args) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pb", thrift.STRUCT, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:pb: ", p), err) }
  if err := p.Pb.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Pb), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:pb: ", p), err) }
  return err
}

func (p *ItnetPon2Args) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "id", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:id: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.ID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.id (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:id: ", p), err) }
  return err
}

func (p *ItnetPon2Args) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPon2Args(%+v)", *p)
}

// Attributes:
//  - PonBeanBytes
//  - ID
//  - Ack
type ItnetPon3Args struct {
  PonBeanBytes []byte `thrift:"PonBeanBytes,1" db:"PonBeanBytes" json:"PonBeanBytes"`
  ID int64 `thrift:"id,2" db:"id" json:"id"`
  Ack bool `thrift:"ack,3" db:"ack" json:"ack"`
}

func NewItnetPon3Args() *ItnetPon3Args {
  return &ItnetPon3Args{}
}


func (p *ItnetPon3Args) GetPonBeanBytes() []byte {
  return p.PonBeanBytes
}

func (p *ItnetPon3Args) GetID() int64 {
  return p.ID
}

func (p *ItnetPon3Args) GetAck() bool {
  return p.Ack
}
func (p *ItnetPon3Args) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPon3Args)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.PonBeanBytes = v
}
  return nil
}

func (p *ItnetPon3Args)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.ID = v
}
  return nil
}

func (p *ItnetPon3Args)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Ack = v
}
  return nil
}

func (p *ItnetPon3Args) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Pon3_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPon3Args) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "PonBeanBytes", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:PonBeanBytes: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.PonBeanBytes); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.PonBeanBytes (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:PonBeanBytes: ", p), err) }
  return err
}

func (p *ItnetPon3Args) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "id", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:id: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.ID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.id (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:id: ", p), err) }
  return err
}

func (p *ItnetPon3Args) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ack", thrift.BOOL, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:ack: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ack)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ack (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:ack: ", p), err) }
  return err
}

func (p *ItnetPon3Args) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPon3Args(%+v)", *p)
}

// Attributes:
//  - Pretimenano
//  - Timenano
//  - Num
//  - Ir
type ItnetTimeArgs struct {
  Pretimenano int64 `thrift:"pretimenano,1" db:"pretimenano" json:"pretimenano"`
  Timenano int64 `thrift:"timenano,2" db:"timenano" json:"timenano"`
  Num int16 `thrift:"num,3" db:"num" json:"num"`
  Ir bool `thrift:"ir,4" db:"ir" json:"ir"`
}

func NewItnetTimeArgs() *ItnetTimeArgs {
  return &ItnetTimeArgs{}
}


func (p *ItnetTimeArgs) GetPretimenano() int64 {
  return p.Pretimenano
}

func (p *ItnetTimeArgs) GetTimenano() int64 {
  return p.Timenano
}

func (p *ItnetTimeArgs) GetNum() int16 {
  return p.Num
}

func (p *ItnetTimeArgs) GetIr() bool {
  return p.Ir
}
func (p *ItnetTimeArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I16 {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetTimeArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Pretimenano = v
}
  return nil
}

func (p *ItnetTimeArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Timenano = v
}
  return nil
}

func (p *ItnetTimeArgs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI16(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Num = v
}
  return nil
}

func (p *ItnetTimeArgs)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.Ir = v
}
  return nil
}

func (p *ItnetTimeArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Time_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetTimeArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pretimenano", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:pretimenano: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Pretimenano)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.pretimenano (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:pretimenano: ", p), err) }
  return err
}

func (p *ItnetTimeArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "timenano", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:timenano: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Timenano)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.timenano (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:timenano: ", p), err) }
  return err
}

func (p *ItnetTimeArgs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "num", thrift.I16, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:num: ", p), err) }
  if err := oprot.WriteI16(ctx, int16(p.Num)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.num (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:num: ", p), err) }
  return err
}

func (p *ItnetTimeArgs) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ir", thrift.BOOL, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:ir: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ir)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ir (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:ir: ", p), err) }
  return err
}

func (p *ItnetTimeArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetTimeArgs(%+v)", *p)
}

// Attributes:
//  - Node
//  - Ir
type ItnetSyncNodeArgs struct {
  Node *Node `thrift:"node,1" db:"node" json:"node"`
  Ir bool `thrift:"ir,2" db:"ir" json:"ir"`
}

func NewItnetSyncNodeArgs() *ItnetSyncNodeArgs {
  return &ItnetSyncNodeArgs{}
}

var ItnetSyncNodeArgs_Node_DEFAULT *Node
func (p *ItnetSyncNodeArgs) GetNode() *Node {
  if !p.IsSetNode() {
    return ItnetSyncNodeArgs_Node_DEFAULT
  }
return p.Node
}

func (p *ItnetSyncNodeArgs) GetIr() bool {
  return p.Ir
}
func (p *ItnetSyncNodeArgs) IsSetNode() bool {
  return p.Node != nil
}

func (p *ItnetSyncNodeArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetSyncNodeArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  p.Node = &Node{}
  if err := p.Node.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Node), err)
  }
  return nil
}

func (p *ItnetSyncNodeArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Ir = v
}
  return nil
}

func (p *ItnetSyncNodeArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SyncNode_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetSyncNodeArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "node", thrift.STRUCT, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:node: ", p), err) }
  if err := p.Node.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Node), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:node: ", p), err) }
  return err
}

func (p *ItnetSyncNodeArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ir", thrift.BOOL, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:ir: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ir)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ir (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:ir: ", p), err) }
  return err
}

func (p *ItnetSyncNodeArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetSyncNodeArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - Result_
type ItnetSyncTxArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  Result_ int8 `thrift:"result,2" db:"result" json:"result"`
}

func NewItnetSyncTxArgs() *ItnetSyncTxArgs {
  return &ItnetSyncTxArgs{}
}


func (p *ItnetSyncTxArgs) GetSyncId() int64 {
  return p.SyncId
}

func (p *ItnetSyncTxArgs) GetResult_() int8 {
  return p.Result_
}
func (p *ItnetSyncTxArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetSyncTxArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetSyncTxArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  temp := int8(v)
  p.Result_ = temp
}
  return nil
}

func (p *ItnetSyncTxArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SyncTx_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetSyncTxArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetSyncTxArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "result", thrift.BYTE, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:result: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Result_)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.result (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:result: ", p), err) }
  return err
}

func (p *ItnetSyncTxArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetSyncTxArgs(%+v)", *p)
}

// Attributes:
//  - SyncList
type ItnetSyncTxMergeArgs struct {
  SyncList map[int64]int8 `thrift:"syncList,1" db:"syncList" json:"syncList"`
}

func NewItnetSyncTxMergeArgs() *ItnetSyncTxMergeArgs {
  return &ItnetSyncTxMergeArgs{}
}


func (p *ItnetSyncTxMergeArgs) GetSyncList() map[int64]int8 {
  return p.SyncList
}
func (p *ItnetSyncTxMergeArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetSyncTxMergeArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[int64]int8, size)
  p.SyncList =  tMap
  for i := 0; i < size; i ++ {
var _key44 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key44 = v
}
var _val45 int8
    if v, err := iprot.ReadByte(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    temp := int8(v)
    _val45 = temp
}
    p.SyncList[_key44] = _val45
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *ItnetSyncTxMergeArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SyncTxMerge_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetSyncTxMergeArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncList", thrift.MAP, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncList: ", p), err) }
  if err := oprot.WriteMapBegin(ctx, thrift.I64, thrift.BYTE, len(p.SyncList)); err != nil {
    return thrift.PrependError("error writing map begin: ", err)
  }
  for k, v := range p.SyncList {
    if err := oprot.WriteI64(ctx, int64(k)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    if err := oprot.WriteByte(ctx, int8(v)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
  }
  if err := oprot.WriteMapEnd(ctx); err != nil {
    return thrift.PrependError("error writing map end: ", err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncList: ", p), err) }
  return err
}

func (p *ItnetSyncTxMergeArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetSyncTxMergeArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - Txid
type ItnetCommitTxArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  Txid int64 `thrift:"txid,2" db:"txid" json:"txid"`
}

func NewItnetCommitTxArgs() *ItnetCommitTxArgs {
  return &ItnetCommitTxArgs{}
}


func (p *ItnetCommitTxArgs) GetSyncId() int64 {
  return p.SyncId
}

func (p *ItnetCommitTxArgs) GetTxid() int64 {
  return p.Txid
}
func (p *ItnetCommitTxArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetCommitTxArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetCommitTxArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Txid = v
}
  return nil
}

func (p *ItnetCommitTxArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CommitTx_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetCommitTxArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetCommitTxArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "txid", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:txid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Txid)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.txid (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:txid: ", p), err) }
  return err
}

func (p *ItnetCommitTxArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetCommitTxArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - Txid
//  - Commit
type ItnetCommitTx2Args struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  Txid int64 `thrift:"txid,2" db:"txid" json:"txid"`
  Commit bool `thrift:"commit,3" db:"commit" json:"commit"`
}

func NewItnetCommitTx2Args() *ItnetCommitTx2Args {
  return &ItnetCommitTx2Args{}
}


func (p *ItnetCommitTx2Args) GetSyncId() int64 {
  return p.SyncId
}

func (p *ItnetCommitTx2Args) GetTxid() int64 {
  return p.Txid
}

func (p *ItnetCommitTx2Args) GetCommit() bool {
  return p.Commit
}
func (p *ItnetCommitTx2Args) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetCommitTx2Args)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetCommitTx2Args)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Txid = v
}
  return nil
}

func (p *ItnetCommitTx2Args)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Commit = v
}
  return nil
}

func (p *ItnetCommitTx2Args) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CommitTx2_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetCommitTx2Args) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetCommitTx2Args) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "txid", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:txid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.Txid)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.txid (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:txid: ", p), err) }
  return err
}

func (p *ItnetCommitTx2Args) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "commit", thrift.BOOL, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:commit: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Commit)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.commit (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:commit: ", p), err) }
  return err
}

func (p *ItnetCommitTx2Args) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetCommitTx2Args(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - MqType
//  - Bs
type ItnetPubMqArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  MqType int8 `thrift:"mqType,2" db:"mqType" json:"mqType"`
  Bs []byte `thrift:"bs,3" db:"bs" json:"bs"`
}

func NewItnetPubMqArgs() *ItnetPubMqArgs {
  return &ItnetPubMqArgs{}
}


func (p *ItnetPubMqArgs) GetSyncId() int64 {
  return p.SyncId
}

func (p *ItnetPubMqArgs) GetMqType() int8 {
  return p.MqType
}

func (p *ItnetPubMqArgs) GetBs() []byte {
  return p.Bs
}
func (p *ItnetPubMqArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPubMqArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetPubMqArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  temp := int8(v)
  p.MqType = temp
}
  return nil
}

func (p *ItnetPubMqArgs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Bs = v
}
  return nil
}

func (p *ItnetPubMqArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "PubMq_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPubMqArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetPubMqArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "mqType", thrift.BYTE, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:mqType: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.MqType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.mqType (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:mqType: ", p), err) }
  return err
}

func (p *ItnetPubMqArgs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "bs", thrift.STRING, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:bs: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.Bs); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.bs (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:bs: ", p), err) }
  return err
}

func (p *ItnetPubMqArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPubMqArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - Ldb
type ItnetPullDataArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  Ldb *LogDataBean `thrift:"ldb,2" db:"ldb" json:"ldb"`
}

func NewItnetPullDataArgs() *ItnetPullDataArgs {
  return &ItnetPullDataArgs{}
}


func (p *ItnetPullDataArgs) GetSyncId() int64 {
  return p.SyncId
}
var ItnetPullDataArgs_Ldb_DEFAULT *LogDataBean
func (p *ItnetPullDataArgs) GetLdb() *LogDataBean {
  if !p.IsSetLdb() {
    return ItnetPullDataArgs_Ldb_DEFAULT
  }
return p.Ldb
}
func (p *ItnetPullDataArgs) IsSetLdb() bool {
  return p.Ldb != nil
}

func (p *ItnetPullDataArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPullDataArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetPullDataArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  p.Ldb = &LogDataBean{}
  if err := p.Ldb.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Ldb), err)
  }
  return nil
}

func (p *ItnetPullDataArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "PullData_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPullDataArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetPullDataArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ldb", thrift.STRUCT, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:ldb: ", p), err) }
  if err := p.Ldb.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Ldb), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:ldb: ", p), err) }
  return err
}

func (p *ItnetPullDataArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPullDataArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - Sbean
type ItnetReInitArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  Sbean *SysBean `thrift:"sbean,2" db:"sbean" json:"sbean"`
}

func NewItnetReInitArgs() *ItnetReInitArgs {
  return &ItnetReInitArgs{}
}


func (p *ItnetReInitArgs) GetSyncId() int64 {
  return p.SyncId
}
var ItnetReInitArgs_Sbean_DEFAULT *SysBean
func (p *ItnetReInitArgs) GetSbean() *SysBean {
  if !p.IsSetSbean() {
    return ItnetReInitArgs_Sbean_DEFAULT
  }
return p.Sbean
}
func (p *ItnetReInitArgs) IsSetSbean() bool {
  return p.Sbean != nil
}

func (p *ItnetReInitArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetReInitArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetReInitArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  p.Sbean = &SysBean{}
  if err := p.Sbean.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Sbean), err)
  }
  return nil
}

func (p *ItnetReInitArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "ReInit_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetReInitArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetReInitArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "sbean", thrift.STRUCT, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:sbean: ", p), err) }
  if err := p.Sbean.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Sbean), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:sbean: ", p), err) }
  return err
}

func (p *ItnetReInitArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetReInitArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - ParamData
//  - PType
//  - Ctype
type ItnetProxyCallArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  ParamData []byte `thrift:"paramData,2" db:"paramData" json:"paramData"`
  PType int8 `thrift:"pType,3" db:"pType" json:"pType"`
  Ctype int8 `thrift:"ctype,4" db:"ctype" json:"ctype"`
}

func NewItnetProxyCallArgs() *ItnetProxyCallArgs {
  return &ItnetProxyCallArgs{}
}


func (p *ItnetProxyCallArgs) GetSyncId() int64 {
  return p.SyncId
}

func (p *ItnetProxyCallArgs) GetParamData() []byte {
  return p.ParamData
}

func (p *ItnetProxyCallArgs) GetPType() int8 {
  return p.PType
}

func (p *ItnetProxyCallArgs) GetCtype() int8 {
  return p.Ctype
}
func (p *ItnetProxyCallArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetProxyCallArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetProxyCallArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.ParamData = v
}
  return nil
}

func (p *ItnetProxyCallArgs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  temp := int8(v)
  p.PType = temp
}
  return nil
}

func (p *ItnetProxyCallArgs)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  temp := int8(v)
  p.Ctype = temp
}
  return nil
}

func (p *ItnetProxyCallArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "ProxyCall_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetProxyCallArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetProxyCallArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "paramData", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:paramData: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.ParamData); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.paramData (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:paramData: ", p), err) }
  return err
}

func (p *ItnetProxyCallArgs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pType", thrift.BYTE, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:pType: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.PType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.pType (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:pType: ", p), err) }
  return err
}

func (p *ItnetProxyCallArgs) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ctype", thrift.BYTE, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:ctype: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Ctype)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ctype (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:ctype: ", p), err) }
  return err
}

func (p *ItnetProxyCallArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetProxyCallArgs(%+v)", *p)
}

// Attributes:
//  - SyncId
//  - Tt
//  - Ack
type ItnetBroadTokenArgs struct {
  SyncId int64 `thrift:"syncId,1" db:"syncId" json:"syncId"`
  Tt *TokenTrans `thrift:"tt,2" db:"tt" json:"tt"`
  Ack bool `thrift:"ack,3" db:"ack" json:"ack"`
}

func NewItnetBroadTokenArgs() *ItnetBroadTokenArgs {
  return &ItnetBroadTokenArgs{}
}


func (p *ItnetBroadTokenArgs) GetSyncId() int64 {
  return p.SyncId
}
var ItnetBroadTokenArgs_Tt_DEFAULT *TokenTrans
func (p *ItnetBroadTokenArgs) GetTt() *TokenTrans {
  if !p.IsSetTt() {
    return ItnetBroadTokenArgs_Tt_DEFAULT
  }
return p.Tt
}

func (p *ItnetBroadTokenArgs) GetAck() bool {
  return p.Ack
}
func (p *ItnetBroadTokenArgs) IsSetTt() bool {
  return p.Tt != nil
}

func (p *ItnetBroadTokenArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetBroadTokenArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncId = v
}
  return nil
}

func (p *ItnetBroadTokenArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  p.Tt = &TokenTrans{}
  if err := p.Tt.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Tt), err)
  }
  return nil
}

func (p *ItnetBroadTokenArgs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Ack = v
}
  return nil
}

func (p *ItnetBroadTokenArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "BroadToken_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetBroadTokenArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SyncId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.syncId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncId: ", p), err) }
  return err
}

func (p *ItnetBroadTokenArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "tt", thrift.STRUCT, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:tt: ", p), err) }
  if err := p.Tt.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Tt), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:tt: ", p), err) }
  return err
}

func (p *ItnetBroadTokenArgs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ack", thrift.BOOL, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:ack: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ack)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ack (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:ack: ", p), err) }
  return err
}

func (p *ItnetBroadTokenArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetBroadTokenArgs(%+v)", *p)
}


