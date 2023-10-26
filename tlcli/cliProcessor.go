// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tlcli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/donnie4w/tldb/keystore"
	. "github.com/donnie4w/tldb/level2"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

type processor struct{}

var cliProcessor = &processor{}

func ctx2CliContext(ctx context.Context) *cliContext {
	return ctx.Value("CliContext").(*cliContext)
}

// Parameters:
//   - I
func (this *processor) Ping(ctx context.Context, i int64) (_r *Ack, _err error) {
	defer myRecovr()
	mux := ctx2CliContext(ctx).mux
	defer mux.Unlock()
	mux.Lock()
	_r = newAck(true, 0, "")
	return
}

// Parameters:
//   - S
func (this *processor) Auth(ctx context.Context, s string) (_r *Ack, _err error) {
	defer myRecovr()
	mux := ctx2CliContext(ctx).mux
	defer mux.Unlock()
	mux.Lock()
	if Auth(s) {
		ctx2CliContext(ctx).isAuth = true
		_r = newAck(true, 0, "")
	} else {
		_r = newAck(false, ERR_AUTH_NOPASS, "no pass")
	}
	return
}

// Parameters:
//   - Tb
func (this *processor) Create(ctx context.Context, tb *TableBean) (_r *Ack, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAck(true, 0, "")
	if noAuthAndClose(cc) {
		_r = newErrAck(int64(sys.ERR_AUTH_NOPASS), "")
		return
	}
	if tb.Name == "" {
		_r = newErrAck(ERR_NO_MATCH_PARAM, "")
		return
	}
	if err := Level2.CreateTable(&TableStub{Tablename: tb.Name, Field: tb.Columns, Idx: tb.Idx}); err != nil {
		_r = newErrAck(0, err.Error())
		return
	}
	return
}

// Parameters:
//   - Tb
func (this *processor) Alter(ctx context.Context, tb *TableBean) (_r *Ack, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAck(true, 0, "")
	if noAuthAndClose(cc) {
		_r = newErrAck(int64(sys.ERR_AUTH_NOPASS), "")
		return
	}
	if tb.Name == "" {
		_r = newErrAck(ERR_NO_MATCH_PARAM, "")
		return
	}
	if err := Level2.AlterTable(&TableStub{Tablename: tb.Name, Field: tb.Columns, Idx: tb.Idx}); err != nil {
		_r = newErrAck(0, err.Error())
		return
	}
	return
}

// Parameters:
//   - Name
func (this *processor) Drop(ctx context.Context, name string) (_r *Ack, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAck(true, 0, "")
	if noAuthAndClose(cc) {
		_r = newErrAck(int64(sys.ERR_AUTH_NOPASS), "")
		return
	}
	if name == "" {
		_r = newErrAck(ERR_NO_MATCH_PARAM, "")
		return
	}
	if err := Level2.DropTable(&TableStub{Tablename: name}); err != nil {
		_r = newErrAck(0, err.Error())
		return
	}
	return
}

// Parameters:
//   - Name
func (this *processor) SelectId(ctx context.Context, name string) (_r int64, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	_r, _err = Level2.SelectId(0, &TableStub{Tablename: name})
	return
}

// Parameters:
//   - Name
func (this *processor) SelectIdByIdx(ctx context.Context, name string, column string, value []byte) (_r int64, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	_r, _err = Level2.SelectIdByIdx(0, name, column, value)
	return
}

// Parameters:
//   - Name
//   - ID
func (this *processor) SelectById(ctx context.Context, name string, id int64) (_r *DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = &DataBean{ID: 0}
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if db, err := Level2.SelectById(0, &TableStub{Tablename: name, ID: id}); err == nil && db != nil {
		_r.ID, _r.TBean = db.ID, db.Field
	}
	return
}

// Parameters:
//   - Name
//   - Column
//   - Value
func (this *processor) SelectByIdx(ctx context.Context, name string, column string, value []byte) (_r *DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = &DataBean{ID: 0}
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if db, err := Level2.SelectByIdx(0, name, column, value); err == nil && db != nil {
		_r.ID, _r.TBean = db.ID, db.Field
	}
	return
}

// Parameters:
//   - Name
//   - Column
//   - Value
func (this *processor) SelectsByIdLimit(ctx context.Context, name string, startId int64, limit int64) (_r []*DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = make([]*DataBean, 0)
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if dbs, err := Level2.SelectsByIdLimit(0, &TableStub{Tablename: name}, startId, limit); err == nil && dbs != nil {
		for _, db := range dbs {
			_d := &DataBean{ID: db.ID}
			_d.TBean = db.Field
			_r = append(_r, _d)
		}
	}
	return
}

// Parameters:
//   - Name
//   - Column
//   - Value
func (this *processor) SelectAllByIdx(ctx context.Context, name string, column string, value []byte) (_r []*DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = make([]*DataBean, 0)
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if dbs, err := Level2.SelectsByIdx(0, name, column, value); err == nil && dbs != nil {
		for _, db := range dbs {
			_d := &DataBean{ID: db.ID}
			_d.TBean = db.Field
			_r = append(_r, _d)
		}
	}
	return
}

// Parameters:
//   - Name
//   - Column
//   - Value
//   - StartId
//   - Limit
func (this *processor) SelectByIdxLimit(ctx context.Context, name string, column string, value [][]byte, startId int64, limit int64) (_r []*DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = make([]*DataBean, 0)
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if dbs, err := Level2.SelectsByIdxLimit(0, name, column, value, startId, limit); err == nil && dbs != nil {
		for _, db := range dbs {
			_d := &DataBean{ID: db.ID}
			_d.TBean = db.Field
			_r = append(_r, _d)
		}
	}
	return
}

// Parameters:
//   - Tb
func (this *processor) Update(ctx context.Context, tb *TableBean) (_r *AckBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAckBean()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if tb != nil && tb.Name != "" && tb.GetID() > 0 && tb.Columns != nil && len(tb.Columns) > 0 {
		ts := &TableStub{Tablename: tb.Name, ID: tb.GetID(), Field: tb.Columns}
		if err := Level2.Update(0, ts); err != nil {
			_r.Ack = newErrAck(0, err.Error())
		}
	} else {
		_r.Ack = newErrAck(ERR_NO_MATCH_PARAM, "")
	}
	return
}

// Parameters:
//   - Tb
func (this *processor) Delete(ctx context.Context, tb *TableBean) (_r *AckBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAckBean()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if tb != nil && tb.Name != "" && tb.GetID() > 0 {
		ts := &TableStub{Tablename: tb.Name, ID: tb.GetID()}
		if err := Level2.Delete(0, ts); err != nil {
			_r.Ack = newErrAck(0, err.Error())
		}
	} else {
		_r.Ack = newErrAck(ERR_NO_MATCH_PARAM, "")
	}
	return
}

// Parameters:
//   - Tb
func (this *processor) Insert(ctx context.Context, tb *TableBean) (_r *AckBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAckBean()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if tb != nil && tb.Name != "" && tb.Columns != nil {
		ts := &TableStub{Tablename: tb.Name, Field: tb.Columns, Idx: tb.Idx}
		if seq, err := Level2.Insert(0, ts); seq > 0 {
			_r.Seq = seq
		} else {
			s := ""
			if err != nil {
				s = err.Error()
			} else {
				s = strconv.Itoa(int(sys.ERR_UNDEFINED))
			}
			_r.Ack = newErrAck(0, s)
		}
	} else {
		_r.Ack = newErrAck(ERR_NO_MATCH_PARAM, "")
	}
	return
}

// Parameters:
//   - Name
func (this *processor) ShowTable(ctx context.Context, name string) (_r *TableBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = NewTableBean()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if t, ok := Tldb.GetTable(name); ok {
		_r.Name = name
		_r.Columns = util.SyncMap2Map(t.Fields)
		_r.Idx = util.SyncMap2Map(t.Idx)
	} else {
		_err = util.Errors(sys.ERR_TABLE_NOEXIST)
	}
	return
}

func (this *processor) ShowAllTables(ctx context.Context) (_r []*TableBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = make([]*TableBean, 0)
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if ls := Tldb.GetTableList(); ls != nil {
		ls.Range(func(k string, v *TableStruct) bool {
			tb := NewTableBean()
			tb.Name = k
			tb.Columns = util.SyncMap2Map(v.Fields)
			tb.Idx = util.SyncMap2Map(v.Idx)
			_r = append(_r, tb)
			return true
		})
	}
	return
}

// Parameters:
//   - Name
//   - Ids
func (this *processor) DeleteBatch(ctx context.Context, name string, ids []int64) (_r *AckBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = newAckBean()
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if name != "" && ids != nil && len(ids) > 0 {
		if err := Level2.DeleteBatch(0, name, ids); err != nil {
			_r.Ack = newErrAck(0, err.Error())
		}
	} else {
		_r.Ack = newErrAck(ERR_NO_MATCH_PARAM, "")
	}
	return
}

// Parameters:
//   - Name
//   - Column
//   - Value
//   - StartId
//   - Limit
func (this *processor) SelectByIdxDescLimit(ctx context.Context, name string, column string, value []byte, startId int64, limit int64) (_r []*DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = make([]*DataBean, 0)
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if dbs, err := Level2.SelectByIdxDescLimit(0, name, column, value, startId, limit); err == nil && dbs != nil {
		for _, db := range dbs {
			_d := &DataBean{ID: db.ID}
			_d.TBean = db.Field
			_r = append(_r, _d)
		}
	}
	return
}

// Parameters:
//   - Name
//   - Column
//   - Value
//   - StartId
//   - Limit
func (this *processor) SelectByIdxAscLimit(ctx context.Context, name string, column string, value []byte, startId int64, limit int64) (_r []*DataBean, _err error) {
	defer myRecovr()
	cc := ctx2CliContext(ctx)
	defer cc.mux.Unlock()
	cc.mux.Lock()
	_r = make([]*DataBean, 0)
	if noAuthAndClose(cc) {
		_err = util.Errors(sys.ERR_AUTH_NOPASS)
		return
	}
	if dbs, err := Level2.SelectByIdxAscLimit(0, name, column, value, startId, limit); err == nil && dbs != nil {
		for _, db := range dbs {
			_d := &DataBean{ID: db.ID}
			_d.TBean = db.Field
			_r = append(_r, _d)
		}
	}
	return
}

/******************************************************************/
func Auth(s string) (_ok bool) {
	defer myRecovr()
	if ss := strings.Split(s, "="); len(ss) == 2 {
		if _r, ok := keystore.StoreAdmin.GetClient(ss[0]); ok {
			if _r.Pwd == util.MD5(ss[1]) {
				_ok = true
			}
		}
	}
	return
}

func newAck(ok bool, code int64, desc string) *Ack {
	return &Ack{ok, code, desc}
}

func newErrAck(code int64, desc string) *Ack {
	if code == 0 {
		code = int64(ERR_UNDEFINED)
		if desc != "" {
			if i, err := strconv.Atoi(desc); err == nil {
				code = int64(i)
			}
		}
	}
	if desc == "" {
		desc = fmt.Sprint(code)
	}
	return &Ack{false, code, desc}
}

func newAckBean() (ab *AckBean) {
	ab = &AckBean{Seq: 0, Ack: newAck(true, 0, "")}
	return
}

func noAuthAndClose(cc *cliContext) (b bool) {
	if !cc.isAuth {
		cc.tt.Close()
		b = true
	}
	return
}
