// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

package keystore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	thrift "github.com/donnie4w/gothrift/thrift"
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
//  - Name
//  - Pwd
//  - Type
type UserBean struct {
  Name string `thrift:"name,1,required" db:"name" json:"name"`
  Pwd string `thrift:"pwd,2,required" db:"pwd" json:"pwd"`
  Type int8 `thrift:"type,3,required" db:"type" json:"type"`
}

func NewUserBean() *UserBean {
  return &UserBean{}
}


func (p *UserBean) GetName() string {
  return p.Name
}

func (p *UserBean) GetPwd() string {
  return p.Pwd
}

func (p *UserBean) GetType() int8 {
  return p.Type
}
func (p *UserBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetName bool = false;
  var issetPwd bool = false;
  var issetType bool = false;

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
        issetName = true
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
        issetPwd = true
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
        issetType = true
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
  if !issetName{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Name is not set"));
  }
  if !issetPwd{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Pwd is not set"));
  }
  if !issetType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Type is not set"));
  }
  return nil
}

func (p *UserBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Name = v
}
  return nil
}

func (p *UserBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Pwd = v
}
  return nil
}

func (p *UserBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  temp := int8(v)
  p.Type = temp
}
  return nil
}

func (p *UserBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "UserBean"); err != nil {
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

func (p *UserBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "name", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:name: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Name)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.name (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:name: ", p), err) }
  return err
}

func (p *UserBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pwd", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:pwd: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Pwd)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.pwd (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:pwd: ", p), err) }
  return err
}

func (p *UserBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "type", thrift.BYTE, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:type: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Type)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.type (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:type: ", p), err) }
  return err
}

func (p *UserBean) Equals(other *UserBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Name != other.Name { return false }
  if p.Pwd != other.Pwd { return false }
  if p.Type != other.Type { return false }
  return true
}

func (p *UserBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("UserBean(%+v)", *p)
}

func (p *UserBean) Validate() error {
  return nil
}
// Attributes:
//  - Admin
//  - Client
//  - Mq
//  - Other
type KeyBean struct {
  Admin map[string]*UserBean `thrift:"admin,1" db:"admin" json:"admin,omitempty"`
  Client map[string]*UserBean `thrift:"client,2" db:"client" json:"client,omitempty"`
  Mq map[string]*UserBean `thrift:"mq,3" db:"mq" json:"mq,omitempty"`
  Other map[string]string `thrift:"other,4" db:"other" json:"other,omitempty"`
}

func NewKeyBean() *KeyBean {
  return &KeyBean{}
}

var KeyBean_Admin_DEFAULT map[string]*UserBean

func (p *KeyBean) GetAdmin() map[string]*UserBean {
  return p.Admin
}
var KeyBean_Client_DEFAULT map[string]*UserBean

func (p *KeyBean) GetClient() map[string]*UserBean {
  return p.Client
}
var KeyBean_Mq_DEFAULT map[string]*UserBean

func (p *KeyBean) GetMq() map[string]*UserBean {
  return p.Mq
}
var KeyBean_Other_DEFAULT map[string]string

func (p *KeyBean) GetOther() map[string]string {
  return p.Other
}
func (p *KeyBean) IsSetAdmin() bool {
  return p.Admin != nil
}

func (p *KeyBean) IsSetClient() bool {
  return p.Client != nil
}

func (p *KeyBean) IsSetMq() bool {
  return p.Mq != nil
}

func (p *KeyBean) IsSetOther() bool {
  return p.Other != nil
}

func (p *KeyBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
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
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
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
    case 4:
      if fieldTypeId == thrift.MAP {
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

func (p *KeyBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]*UserBean, size)
  p.Admin =  tMap
  for i := 0; i < size; i ++ {
var _key0 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key0 = v
}
    _val1 := &UserBean{}
    if err := _val1.Read(ctx, iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _val1), err)
    }
    p.Admin[_key0] = _val1
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *KeyBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]*UserBean, size)
  p.Client =  tMap
  for i := 0; i < size; i ++ {
var _key2 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key2 = v
}
    _val3 := &UserBean{}
    if err := _val3.Read(ctx, iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _val3), err)
    }
    p.Client[_key2] = _val3
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *KeyBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]*UserBean, size)
  p.Mq =  tMap
  for i := 0; i < size; i ++ {
var _key4 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key4 = v
}
    _val5 := &UserBean{}
    if err := _val5.Read(ctx, iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _val5), err)
    }
    p.Mq[_key4] = _val5
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *KeyBean)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]string, size)
  p.Other =  tMap
  for i := 0; i < size; i ++ {
var _key6 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key6 = v
}
var _val7 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val7 = v
}
    p.Other[_key6] = _val7
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *KeyBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "KeyBean"); err != nil {
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

func (p *KeyBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetAdmin() {
    if err := oprot.WriteFieldBegin(ctx, "admin", thrift.MAP, 1); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:admin: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRUCT, len(p.Admin)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Admin {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := v.Write(ctx, oprot); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
      }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 1:admin: ", p), err) }
  }
  return err
}

func (p *KeyBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetClient() {
    if err := oprot.WriteFieldBegin(ctx, "client", thrift.MAP, 2); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:client: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRUCT, len(p.Client)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Client {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := v.Write(ctx, oprot); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
      }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 2:client: ", p), err) }
  }
  return err
}

func (p *KeyBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetMq() {
    if err := oprot.WriteFieldBegin(ctx, "mq", thrift.MAP, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:mq: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRUCT, len(p.Mq)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Mq {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := v.Write(ctx, oprot); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
      }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:mq: ", p), err) }
  }
  return err
}

func (p *KeyBean) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetOther() {
    if err := oprot.WriteFieldBegin(ctx, "other", thrift.MAP, 4); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:other: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Other)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Other {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteString(ctx, string(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 4:other: ", p), err) }
  }
  return err
}

func (p *KeyBean) Equals(other *KeyBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if len(p.Admin) != len(other.Admin) { return false }
  for k, _tgt := range p.Admin {
    _src8 := other.Admin[k]
    if !_tgt.Equals(_src8) { return false }
  }
  if len(p.Client) != len(other.Client) { return false }
  for k, _tgt := range p.Client {
    _src9 := other.Client[k]
    if !_tgt.Equals(_src9) { return false }
  }
  if len(p.Mq) != len(other.Mq) { return false }
  for k, _tgt := range p.Mq {
    _src10 := other.Mq[k]
    if !_tgt.Equals(_src10) { return false }
  }
  if len(p.Other) != len(other.Other) { return false }
  for k, _tgt := range p.Other {
    _src11 := other.Other[k]
    if _tgt != _src11 { return false }
  }
  return true
}

func (p *KeyBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("KeyBean(%+v)", *p)
}

func (p *KeyBean) Validate() error {
  return nil
}
