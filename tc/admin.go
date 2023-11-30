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
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/key"
	. "github.com/donnie4w/tldb/keystore"
	. "github.com/donnie4w/tldb/level2"
	"github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/tlmq"
	"github.com/donnie4w/tldb/tnet"
	"github.com/donnie4w/tldb/util"
	"github.com/donnie4w/tlnet"
)

type adminService struct {
	isClose bool
	tlAdmin *tlnet.Tlnet
}

var adminservice = &adminService{false, tlnet.NewTlnet()}

func (this *adminService) Serve() (err error) {
	if strings.TrimSpace(sys.DEBUGADDR) != "" {
		go tlDebug()
		<-time.After(500 * time.Millisecond)
	}
	if sys.INIT {
		initAccount()
	}
	if strings.TrimSpace(sys.WEBADMINADDR) != "" {
		err = this._serve(strings.TrimSpace(sys.WEBADMINADDR), sys.ADMINTLS, sys.ADMINCRT, sys.ADMINKEY)
	} else{
		sys.FmtLog("no webAdmin service")
	}
	return
}

func (this *adminService) Close() (err error) {
	defer util.Recovr()
	if strings.TrimSpace(sys.WEBADMINADDR) != "" {
		this.isClose = true
		err = this.tlAdmin.Close()
	}
	return
}

func (this *adminService) _serve(addr string, TLS bool, serverCrt, serverKey string) (err error) {
	sys.WEBADMINADDR = addr
	StoreAdmin.PutOther("admin", sys.WEBADMINADDR)
	StoreAdmin.PutOther("admintauth", fmt.Sprint(uint(util.NewTxId())))
	this.tlAdmin.Handle("/login", loginHandler)
	this.tlAdmin.Handle("/init", initHandler)
	this.tlAdmin.Handle("/lang", langHandler)
	this.tlAdmin.Handle("/", initHandler)
	this.tlAdmin.Handle("/bootstrap.css", cssHandler)
	this.tlAdmin.HandleWithFilter("/sysvar", loginFilter(), sysVarHtml)
	this.tlAdmin.HandleWithFilter("/data", loginFilter(), dataHtml)
	this.tlAdmin.HandleWithFilter("/export", loginFilter(), exportHandler)

	this.tlAdmin.HandleWithFilter("/mq", loginFilter(), mqHtml)
	this.tlAdmin.HandleWithFilter("/sys", loginFilter(), sysParamHtml)
	this.tlAdmin.HandleWithFilter("/log", loginFilter(), func(hc *tlnet.HttpContext) { hc.ResponseString(log.SysLog()) })

	this.tlAdmin.HandleWithFilter("/create", loginFilter(), createHtml)
	this.tlAdmin.HandleWithFilter("/alter", loginFilter(), alterHtml)
	this.tlAdmin.HandleWithFilter("/insert", loginFilter(), insertHtml)
	this.tlAdmin.HandleWithFilter("/delete", loginFilter(), deleteHtml)
	this.tlAdmin.HandleWithFilter("/update", loginFilter(), updateHtml)
	this.tlAdmin.HandleWithFilter("/drop", loginFilter(), dropHtml)
	this.tlAdmin.HandleWithFilter("/load", loginFilter(), loadHtml)
	this.tlAdmin.HandleWebSocketBindConfig("/loadData", nil, wsConfig())

	this.tlAdmin.HandleWithFilter("/monitor", loginFilter(), monitorHtml)
	this.tlAdmin.HandleWebSocketBindConfig("/monitorData", mntHandler, mntConfig())

	this.tlAdmin.HandleWebSocketBindConfig("/local", localHandler, localwsConfig())
	if TLS {
		StoreAdmin.PutOther("admintls", "1")
		if util.IsFileExist(serverCrt) && util.IsFileExist(serverKey) {
			sys.FmtLog(fmt.Sprint("webAdmin start tls [", addr, "]"))
			if err = this.tlAdmin.HttpStartTLS(addr, serverCrt, serverKey); err != nil {
				err = errors.New("webAdmin start tls failed:" + err.Error())
			}
		} else {
			sys.FmtLog(fmt.Sprint("webAdmin start tls by bytes [", addr, "]"))
			if err = this.tlAdmin.HttpStartTlsBytes(addr, []byte(ServerCrt), []byte(ServerKey)); err != nil {
				err = errors.New("webAdmin start tls by bytes failed:" + err.Error())
			}
		}
	}
	if !this.isClose {
		sys.FmtLog(fmt.Sprint("webAdmin start [", addr, "]"))
		StoreAdmin.PutOther("admintls", "0")
		if err = this.tlAdmin.HttpStart(addr); err != nil {
			err = errors.New("webAdmin start failed:" + err.Error())
		}
	}
	if !this.isClose && err != nil {
		sys.FmtLog(err.Error())
		os.Exit(0)
	}
	return
}

var sessionMap = NewMapL[string, *UserBean]()

func loginFilter() (f *tlnet.Filter) {
	defer recover()
	f = tlnet.NewFilter()
	f.AddIntercept(".*?", func(hc *tlnet.HttpContext) bool {
		if len(StoreAdmin.AdminList()) > 0 {
			if !isLogin(hc) {
				hc.Redirect("/login")
				return true
			}
		} else {
			hc.Redirect("/init")
			return true
		}
		return false
	})

	f.AddIntercept(`[^\s]+`, func(hc *tlnet.HttpContext) bool {
		if hc.PostParamTrimSpace("atype") != "" && !isAdmin(hc) {
			hc.ResponseString(resultHtml("Permission Denied"))
			return true
		}
		return false
	})
	return
}

func getSessionid() string {
	return fmt.Sprint("t", util.CRC32(util.Int64ToBytes(sys.UUID)))
}

func getLangId() string {
	return fmt.Sprint("l", util.CRC32(util.Int64ToBytes(sys.UUID)))
}

func isLogin(hc *tlnet.HttpContext) (isLogin bool) {
	if len(StoreAdmin.AdminList()) > 0 {
		if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
			isLogin = true
		}
	}
	return
}

func isAdmin(hc *tlnet.HttpContext) (_r bool) {
	if c, err := hc.GetCookie(getSessionid()); err == nil {
		if u, ok := sessionMap.Get(c); ok {
			_r = u.Type == 1
		}
	}
	return
}

func langHandler(hc *tlnet.HttpContext) {
	defer recover()
	lang := hc.GetParamTrimSpace("lang")
	if lang == "en" || lang == "zh" {
		hc.SetCookie(getLangId(), lang, "/", 86400)
	}
	hc.Redirect("/")
}

func getLang(hc *tlnet.HttpContext) LANG {
	if lang, err := hc.GetCookie(getLangId()); err == nil {
		if lang == "zh" {
			return ZH
		} else if lang == "en" {
			return EN
		}
	}
	return ZH
}

func cssHandler(hc *tlnet.HttpContext) {
	hc.Writer().Header().Add("Content-Type", "text/html")
	textTplByText(cssContent(), nil, hc)
}

/***********************************************************************/
func initHandler(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if len(StoreAdmin.AdminList()) > 0 && !isLogin(hc) {
		hc.Redirect("/login")
		return
	}
	if _type := hc.GetParam("type"); _type != "" {
		isadmin := isAdmin(hc)
		if _type == "1" {
			if name, pwd, _type := hc.PostParamTrimSpace("adminName"), hc.PostParamTrimSpace("adminPwd"), hc.PostParamTrimSpace("adminType"); name != "" && pwd != "" {
				if n := len(StoreAdmin.AdminList()); (n > 0 && isadmin) || n == 0 {
					alterType := false
					if t, err := strconv.Atoi(_type); err == nil {
						if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
							if u, ok := sessionMap.Get(_r); ok && u.Name == name && t != int(u.Type) {
								alterType = true
							}
						}
						if !alterType {
							StoreAdmin.PutAdmin(name, pwd, int8(t))
						}
					}
				} else {
					goto DENIED
				}
			}
			if name, pwd := hc.PostParamTrimSpace("mqName"), hc.PostParamTrimSpace("mqPwd"); name != "" && pwd != "" {
				if isadmin {
					StoreAdmin.PutMq(name, pwd, 1)
				} else {
					goto DENIED
				}
			}
			if name, pwd := hc.PostParamTrimSpace("cliName"), hc.PostParamTrimSpace("cliPwd"); name != "" && pwd != "" {
				if isadmin {
					StoreAdmin.PutClient(name, pwd, 1)
				} else {
					goto DENIED
				}
			}
		} else if _type == "2" && isLogin(hc) {
			if isadmin {
				if name := hc.PostParamTrimSpace("adminName"); name != "" {
					if u, ok := StoreAdmin.GetAdmin(name); ok && u.Type == 1 {
						i, j := 0, 0
						for _, s := range StoreAdmin.AdminList() {
							if _u, _ := StoreAdmin.GetAdmin(s); _u.Type == 1 {
								i++
							} else if _u.Type == 2 {
								j++
							}
						}
						if j > 0 && i == 1 {
							hc.ResponseString(resultHtml("failed,There cannot be only Data-only users"))
							return
						}
					}
					StoreAdmin.DelAdmin(name)
					sessionMap.Range(func(k string, v *UserBean) bool {
						if v.Name == name {
							sessionMap.Del(k)
						}
						return true
					})
				}
				if name := hc.PostParamTrimSpace("mqName"); name != "" {
					StoreAdmin.DelMq(name)
				}
				if name := hc.PostParamTrimSpace("cliName"); name != "" {
					StoreAdmin.DelClient(name)
				}
			} else {
				goto DENIED
			}
		}
		hc.Redirect("/init")
		return
	} else {
		initHtml(hc)
		return
	}
DENIED:
	hc.ResponseString(resultHtml("Permission Denied"))
}

func loginHandler(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if hc.PostParamTrimSpace("type") == "1" {
		name, pwd := hc.PostParamTrimSpace("name"), hc.PostParamTrimSpace("pwd")
		if _r, ok := StoreAdmin.GetAdmin(name); ok {
			if _r.Pwd == util.MD5(pwd) {
				sid := util.MD5(fmt.Sprint(time.Now().UnixNano()))
				sessionMap.Put(sid, _r)
				hc.SetCookie(getSessionid(), sid, "/", 86400)
				hc.Redirect("/")
				return
			}
		}
		hc.ResponseString(resultHtml("Login Failed"))
		return
	}
	loginHtml(hc)
}

/*****************************************************************************/
func initHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	_isAdmin := isAdmin(hc)
	show, init, stat, sc := "", false, sys.IsClusRun(), _isAdmin
	if len(StoreAdmin.AdminList()) == 0 {
		show, init, sc = "no user is created for admin, create a management user first", true, true
	}
	av := &AdminView{Show: show, Stat: stat, Init: init, ShowCreate: sc}
	if isLogin(hc) {
		m := make(map[string]string, 0)
		for _, s := range StoreAdmin.AdminList() {
			if u, ok := StoreAdmin.GetAdmin(s); ok {
				if _isAdmin && u.Type == 1 {
					m[s] = "Admin"
				} else if u.Type == 2 {
					m[s] = "Data-only"
				}
			}
		}
		av.AdminUser, av.CliUser, av.MqUser = m, StoreAdmin.ClientList(), StoreAdmin.MqList()
	}
	tplToHtml(getLang(hc), INIT, av, hc)
}

func loginHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	tplToHtml(getLang(hc), LOGIN, []byte{}, hc)
}

func sysVarHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	rn := sys.GetRemoteNode()
	sort.Slice(rn, func(i, j int) bool { return rn[i].UUID > rn[j].UUID })
	svv := &SysVarView{Show: "", RN: rn}
	if _type := hc.PostParamTrimSpace("atype"); _type != "" {
		if _type == "1" {
			_addr := hc.PostParamTrimSpace("addr")
			if addr := strings.Trim(_addr, " "); addr != "" {
				if err := tnet.AddNode(addr); err != nil {
					hc.ResponseString(resultHtml("Failed :", err.Error()))
					return
				} else {
					<-time.After(1000 * time.Millisecond)
					svv.Show = "ADD NODE [ " + _addr + " ]"
				}
			}
		} else if _type == "2" {
			if time_deviation, err := strconv.Atoi(hc.PostParamTrimSpace("time_deviation")); err == nil {
				if time.Now().Add(time.Duration(time_deviation)).After(util.Time()) {
					sys.TIME_DEVIATION = int64(time_deviation)
					svv.Show = "alter time successful."
				} else {
					svv.Show = "failed,The restoration time cannot be before the cluster time"
				}
			}
		} else if _type == "3" {
			if !sys.IsRUN() {
				hc.ResponseString(resultHtml("must be run state"))
				return
			} else {
				if sys.TryRunToProxy() != nil {
					hc.ResponseString(resultHtml("reset stat failed"))
					return
				}
			}
		} else if _type == "4" {
			if !sys.IsPROXY() {
				hc.ResponseString(resultHtml("must be proxy state"))
				return
			} else {
				if sys.TryProxyToReady() != nil {
					hc.ResponseString(resultHtml("reset stat failed"))
					return
				}
			}
		} else if _type == "5" {
			if storeNumForm, err := strconv.Atoi(hc.PostParamTrimSpace("storeNum")); err == nil {
				if err = sys.ReSetStoreNodeNumber(int32(storeNumForm)); err == nil {
					svv.Show = "successful,reset store number to:" + hc.PostParamTrimSpace("storeNum")
				}
			}
		}
	}
	svv.SYS, svv.Stat = sysvar(), sys.IsClusRun()
	for _, cn := range svv.RN {
		cn.StatDesc = statStr(sys.STATTYPE(cn.Stat))
	}
	tplToHtml(getLang(hc), SYSVAR, svv, hc)
}

func sysParamHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if _type := hc.PostParamTrimSpace("atype"); _type != "" {
		if _type == "1" {
			if !sys.IsPROXY() && !sys.IsStandAlone() {
				hc.ResponseString(resultHtml("failed, must be proxy stat"))
				return
			}
			f, _, e := hc.FormFile("loadfile1")
			if e == nil {
				var buf bytes.Buffer
				io.Copy(&buf, f)
				if err := util.CheckGzip(buf.Bytes()); err == nil {
					go sys.LoadData2TLDB(buf.Bytes(), "")
					hc.Redirect("/load")
				} else {
					hc.ResponseString(resultHtml("Failed :", err))
					return
				}
			}
		} else if _type == "2" {
			if !sys.IsPROXY() && !sys.IsStandAlone() {
				hc.ResponseString(resultHtml("failed, must be proxy stat"))
				return
			}
			f, _, e := hc.FormFile("loadfile2")
			if e == nil {
				var buf bytes.Buffer
				io.Copy(&buf, f)
				if err := util.CheckGzip(buf.Bytes()); err == nil {
					go sys.ForcedCoverageData2TLDB(buf.Bytes(), "", 0)
					hc.Redirect("/load")
				} else {
					hc.ResponseString(resultHtml("Failed :", err))
					return
				}
			}
		} else if _type == "3" {
			tnet.CloseSelf()
			sys.Stop()
		}
	}
	spv := &SysParamView{SYS: &SysParam{}, Stat: sys.IsClusRun()}
	spv.SYS.ADMINCRT = sys.ADMINCRT
	spv.SYS.ADMINKEY = sys.ADMINKEY
	spv.SYS.ADMINTLS = sys.ADMINTLS
	spv.SYS.CLICRT = sys.CLICRT
	spv.SYS.CLIKEY = sys.CLIKEY
	spv.SYS.CLITLS = sys.CLITLS
	spv.SYS.COCURRENT_GET = sys.COCURRENT_GET
	spv.SYS.COCURRENT_PUT = sys.COCURRENT_PUT
	spv.SYS.DBFILEDIR = sys.DBFILEDIR
	spv.SYS.DBMode = sys.DBMode
	spv.SYS.MQCRT = sys.MQCRT
	spv.SYS.MQKEY = sys.MQKEY
	spv.SYS.MQTLS = sys.MQTLS
	spv.SYS.CLUSTER_NUM_FINAL = sys.CLUSTER_NUM_FINAL
	spv.SYS.NAMESPACE = sys.NAMESPACE
	spv.SYS.VERSION = sys.VERSION
	spv.SYS.BINLOGSIZE = sys.BINLOGSIZE
	spv.SYS.ADDR = sys.ADDR
	spv.SYS.CLIADDR = sys.CLIADDR
	spv.SYS.MQADDR = sys.MQADDR
	spv.SYS.WEBADMINADDR = sys.WEBADMINADDR
	spv.SYS.CLUSTER_NUM = sys.CLUSTER_NUM
	spv.SYS.PUBLICKEY = sys.PUBLICKEY
	spv.SYS.PRIVATEKEY = sys.PRIVATEKEY
	spv.SYS.PWD = sys.PWD
	tplToHtml(getLang(hc), SYS, spv, hc)
}

func dataHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	dv := &DataView{make([]*Tables, 0), make([]*TData, 0), make(map[string][]byte, 0), &SelectBean{}, sys.IsClusRun()}
	Tldb.LoadTableInfo()
	if gtl := Tldb.GetTableList(); gtl != nil {
		gtl.Range(func(k string, v *TableStruct) bool {
			t := &Tables{Name: k}
			t.Columns = util.MapToArray(v.Fields)
			t.Idxs = util.MapToArray(v.Idx)
			t.Seq, _ = Level2.SelectId(0, &TableStub{Tablename: k})
			dv.Tb = append(dv.Tb, t)
			return true
		})
	}
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		if _type == "1" {
			tableName, id := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("tableId")
			if tst, ok := Tldb.GetTable(tableName); ok {
				dv.ColName = util.SyncMap2Map(tst.Fields)
				dv.Sb.Name, dv.Sb.Id = tableName, id
				if _id, err := strconv.Atoi(id); err == nil {
					if ts, err := Level2.SelectById(0, &TableStub{Tablename: tableName, ID: int64(_id)}); err == nil {
						td := &TData{}
						td.Name, td.Id = ts.Tablename, ts.ID
						td.Columns = make(map[string]string, 0)
						for k, v := range dv.ColName {
							td.Columns[k] = type2value(v, ts.Field[k])
						}
						dv.Tds = append(dv.Tds, td)
					}
				}
			}
		} else if _type == "2" {
			tableName, cloName, cloValue, start, limit := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("cloName"), hc.PostParamTrimSpace("cloValue"), hc.PostParamTrimSpace("start"), hc.PostParamTrimSpace("limit")
			if tst, ok := Tldb.GetTable(tableName); ok {
				dv.ColName = util.SyncMap2Map(tst.Fields)
				dv.Sb.Name, dv.Sb.ColumnName, dv.Sb.ColumnValue, dv.Sb.StartId, dv.Sb.Limit = tableName, cloName, cloValue, start, limit
				if _start, err := strconv.Atoi(start); err == nil {
					if _limit, err := strconv.Atoi(limit); err == nil {
						bs, _ := tst.Fields.Get(cloName)
						clobyte, _ := valueToBytes(bs, cloValue)
						if tss, err := Level2.SelectsByIdxLimit(0, tableName, cloName, [][]byte{clobyte}, int64(_start), int64(_limit)); err == nil {
							for _, ts := range tss {
								td := &TData{}
								td.Name, td.Id = ts.Tablename, ts.ID
								td.Columns = make(map[string]string, 0)
								for k, v := range dv.ColName {
									td.Columns[k] = type2value(v, ts.Field[k])
								}
								dv.Tds = append(dv.Tds, td)
							}
						}
					}
				}
			}
		} else if _type == "3" {
			tableName, start, limit := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("start"), hc.PostParamTrimSpace("limit")
			if tst, ok := Tldb.GetTable(tableName); ok {
				dv.ColName = util.SyncMap2Map(tst.Fields)
				dv.Sb.Name, dv.Sb.StartId, dv.Sb.Limit = tableName, start, limit
				if _start, err := strconv.Atoi(start); err == nil {
					if _limit, err := strconv.Atoi(limit); err == nil {
						if tss, err := Level2.SelectsByIdLimit(0, &TableStub{Tablename: tableName}, int64(_start), int64(_limit)); err == nil {
							for _, ts := range tss {
								td := &TData{}
								td.Name, td.Id = tableName, ts.ID
								td.Columns = make(map[string]string, 0)
								for k, v := range dv.ColName {
									td.Columns[k] = type2value(v, ts.Field[k])
								}
								dv.Tds = append(dv.Tds, td)
							}
						}
					}
				}
			}
		}
	}
	sort.Slice(dv.Tb, func(i, j int) bool { return dv.Tb[i].Name < dv.Tb[j].Name })
	for _, cn := range dv.Tb {
		sort.Strings(cn.Columns)
		sort.Strings(cn.Idxs)
	}
	tplToHtml(getLang(hc), DATA, dv, hc)
}

func mqHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	dv := &DataView{make([]*Tables, 0), make([]*TData, 0), make(map[string][]byte, 0), &SelectBean{}, sys.IsClusRun()}
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		if _type == "2" {
			tableName, id := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("tableId")
			if tst, ok := TldbMq.GetTable(tableName); ok {
				dv.ColName = util.SyncMap2Map(tst.Fields)
				dv.Sb.Name, dv.Sb.Id = tableName, id
				if _id, err := strconv.Atoi(id); err == nil {
					if ts, err := Level2.SelectById(0, &TableStub{Tablename: key.Topic(tableName), ID: int64(_id)}); err == nil {
						td := &TData{}
						td.Name, td.Id = tableName, ts.ID
						td.Columns = make(map[string]string, 0)
						for k := range dv.ColName {
							td.Columns[k] = string(ts.Field[k])
						}
						dv.Tds = append(dv.Tds, td)
					}
				}
			}
		} else if _type == "3" {
			tableName, start, limit := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("start"), hc.PostParamTrimSpace("limit")
			if tst, ok := TldbMq.GetTable(tableName); ok {
				dv.ColName = util.SyncMap2Map(tst.Fields)
				dv.Sb.Name, dv.Sb.StartId, dv.Sb.Limit = tableName, start, limit
				if _start, err := strconv.Atoi(start); err == nil {
					if _limit, err := strconv.Atoi(limit); err == nil {
						if tss, err := Level2.SelectsByIdLimit(0, &TableStub{Tablename: key.Topic(tableName)}, int64(_start), int64(_limit)); err == nil {
							for _, _ts := range tss {
								td := &TData{}
								td.Name, td.Id = tableName, _ts.ID
								td.Columns = make(map[string]string, 0)
								for k := range dv.ColName {
									td.Columns[k] = string(_ts.Field[k])
								}
								dv.Tds = append(dv.Tds, td)
							}
						}
					}
				}
			}
		}
	} else if _type := hc.PostParamTrimSpace("atype"); _type != "" {
		if _type == "1" {
			tableName := hc.PostParamTrimSpace("tableName")
			Level2.DropTable(&TableStub{Tablename: key.Topic(tableName)})
			tlmq.MqWare.DelTopic(tableName)
		} else if _type == "2" {
			tableName, fromId, limit := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("fromId"), hc.PostParamTrimSpace("limit")
			fromId64, _ := strconv.ParseInt(fromId, 10, 64)
			limit64, _ := strconv.ParseInt(limit, 10, 64)
			if tableName != "" && fromId64 >= 0 && limit64 > 0 {
				Level2.DeleteBatches(0, key.Topic(tableName), fromId64, limit64)
			}
		}
	}

	TldbMq.LoadTableInfo()
	if gtl := TldbMq.GetTableList(); gtl != nil {
		gtl.Range(func(_ string, v *TableStruct) bool {
			t := &Tables{Name: v.Tablename}
			t.Columns = util.MapToArray(v.Fields)
			t.Idxs = util.MapToArray(v.Idx)
			t.Seq, _ = Level2.SelectId(0, &TableStub{Tablename: key.Topic(v.Tablename)})
			t.Sub = tlmq.MqWare.SubCount(v.Tablename)
			dv.Tb = append(dv.Tb, t)
			return true
		})
	}

	sort.Slice(dv.Tb, func(i, j int) bool { return dv.Tb[i].Name < dv.Tb[j].Name })
	tplToHtml(getLang(hc), MQ, dv, hc)
}

func createHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	var err error
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		var tableName string
		if tableName = hc.PostParamTrimSpace("tableName"); tableName != "" {
			colums := hc.PostParams("colum")
			if tableName == "" || colums == nil || len(colums) == 0 {
				err = util.Errors(sys.ERR_NO_MATCH_PARAM)
			} else {
				idxs := hc.PostParams("index")
				ftype := hc.PostParams("ftype")
				ts := &TableStub{}
				ts.Tablename = tableName
				ts.Field = make(map[string][]byte, 0)
				ids := make([]string, 0)
				for i := 0; i < len(colums); i++ {
					if idxs[i] == "true" {
						ids = append(ids, colums[i])
					}
					ts.Field[colums[i]] = []byte(ftype[i])
				}
				if len(ids) > 0 {
					ts.Idx = util.ArrayToMap2(ids, int8(0))
				}
				if _type == "1" {
					err = Level2.CreateTable(ts)
				} else if _type == "2" {
					err = Level2.AlterTable(ts)
				}
			}
		} else {
			err = util.Errors(sys.ERR_NO_MATCH_PARAM)
		}
		if err != nil {
			hc.ResponseString(resultHtml("Failed :", err))
		} else {
			hc.ResponseString(resultHtmlAndClose("successful :", tableName))
		}
		return
	}
	tplToHtml(getLang(hc), CREATE, nil, hc)
}

func alterHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	var t = &AlterTable{}
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		if _type == "1" {
			var err error
			if name := hc.PostParamTrimSpace("tableName"); name != "" {
				if tb, ok := Tldb.GetTable(name); ok {
					t.TableName = tb.Tablename
					t.Columns = make(map[string]*FieldInfo, 0)
					tb.Fields.Range(func(k string, v []byte) bool {
						t.Columns[k] = &FieldInfo{Idx: tb.Idx.Has(k), Type: string(v), Tname: type2Name(v)}
						return true
					})
				} else {
					err = util.Errors(sys.ERR_TABLE_NOEXIST)
				}
			}
			if err != nil {
				hc.ResponseString(resultHtml("Failed :", err))
				return
			}
		}
	}
	tplToHtml(getLang(hc), ALTER, t, hc)
}

func dropHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		if _type == "1" {
			if tn := hc.PostParamTrimSpace("tableName"); tn != "" {
				if err := Level2.DropTable(&TableStub{Tablename: tn}); err != nil {
					hc.ResponseString(resultHtml("Failed :", err))
					return
				}
			}
		}
	}
	dv := &DataView{Tb: make([]*Tables, 0)}
	Tldb.LoadTableInfo()
	if gtl := Tldb.GetTableList(); gtl != nil {
		gtl.Range(func(k string, v *TableStruct) bool {
			t := &Tables{Name: k}
			t.Columns = util.MapToArray(v.Fields)
			t.Idxs = util.MapToArray(v.Idx)
			t.Seq, _ = Level2.SelectId(0, &TableStub{Tablename: k})
			dv.Tb = append(dv.Tb, t)
			return true
		})
	}
	tplToHtml(getLang(hc), DROP, dv, hc)
}

func insertHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	var t = &AlterTable{}
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		var err error
		if _type == "1" {
			if name := hc.PostParamTrimSpace("tableName"); name != "" {
				if tst, ok := Tldb.GetTable(name); ok {
					t.TableName = tst.Tablename
					t.Columns = make(map[string]*FieldInfo, 0)
					tst.Fields.Range(func(k string, _ []byte) bool {
						t.Columns[k] = nil
						return true
					})
				} else {
					err = util.Errors(sys.ERR_TABLE_NOEXIST)
				}
			}
		} else if _type == "2" {
			if tableName := hc.PostParamTrimSpace("tableName"); tableName != "" {
				if tst, ok := Tldb.GetTable(tableName); ok {
					colums := hc.PostParams("colums")
					values := hc.PostParams("values")
					ts := &TableStub{Tablename: tableName}
					ts.Field = make(map[string][]byte, 0)
					for i := 0; i < len(colums); i++ {
						if bs, ok := tst.Fields.Get(colums[i]); ok {
							if ts.Field[colums[i]], err = valueToBytes(bs, values[i]); err != nil {
								goto END
							}
						}
					}
					if _, err = Level2.Insert(0, ts); err == nil {
						hc.ResponseString(resultHtmlAndClose("successful!"))
						return
					}
				} else {
					err = util.Errors(sys.ERR_TABLE_NOEXIST)
				}
			}
		}
	END:
		if err != nil {
			hc.ResponseString(resultHtml("Failed :", err))
			return
		}
	}
	tplToHtml(getLang(hc), INSERT, t, hc)
}

func updateHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	var t = &AlterTable{}
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		var err error
		if _type == "1" {
			tableName, id := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("tableId")
			if tst, _ := Tldb.GetTable(tableName); tst != nil {
				var _id int
				if _id, err = strconv.Atoi(id); err == nil {
					var ts *TableStub
					if ts, err = Level2.SelectById(0, &TableStub{Tablename: tableName, ID: int64(_id)}); err == nil {
						t.TableName = tableName
						t.ColumnValue = make(map[string]string, 0)
						t.ID = int64(_id)
						for k, v := range ts.Field {
							if bs, ok := tst.Fields.Get(k); ok {
								t.ColumnValue[k] = type2value(bs, v)
							}
						}
						tst.Fields.Range(func(k string, _ []byte) bool {
							if _, ok := t.ColumnValue[k]; !ok {
								t.ColumnValue[k] = ""
							}
							return true
						})
					}
				} else {
					err = util.Errors(sys.ERR_NO_MATCH_PARAM)
				}
			} else {
				err = util.Errors(sys.ERR_TABLE_NOEXIST)
			}
		} else {
			if tableName := hc.PostParamTrimSpace("tableName"); tableName != "" {
				if tst, _ := Tldb.GetTable(tableName); tst != nil {
					id := hc.PostParamTrimSpace("tableId")
					var _id int
					if _id, err = strconv.Atoi(id); err == nil {
						colums := hc.PostParams("colums")
						values := hc.PostParams("values")
						ts := &TableStub{Tablename: tableName, ID: int64(_id)}
						ts.Field = make(map[string][]byte, 0)
						for i := 0; i < len(colums); i++ {
							if bs, ok := tst.Fields.Get(colums[i]); ok {
								if ts.Field[colums[i]], err = valueToBytes(bs, values[i]); err != nil {
									goto END
								}
							}
						}
						if err = Level2.Update(0, ts); err == nil {
							hc.ResponseString(resultHtmlAndClose("successful!"))
							return
						}
					} else {
						err = util.Errors(sys.ERR_NO_MATCH_PARAM)
					}
				} else {
					err = util.Errors(sys.ERR_TABLE_NOEXIST)
				}
			}
		}
	END:
		if err != nil {
			hc.ResponseString(resultHtml("Failed :", err))
			return
		}
	}
	tplToHtml(getLang(hc), UPDATE, t, hc)
}

func deleteHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	var t = &AlterTable{}
	if _type := hc.PostParamTrimSpace("type"); _type != "" {
		var err error
		if _type == "1" {
			tableName, id := hc.PostParamTrimSpace("tableName"), hc.PostParamTrimSpace("tableId")
			var _id int
			if _id, err = strconv.Atoi(id); err == nil {
				if tst, ok := Tldb.GetTable(tableName); ok {
					var ts *TableStub
					if ts, err = Level2.SelectById(0, &TableStub{Tablename: tableName, ID: int64(_id)}); err == nil {
						t.TableName = tableName
						t.ColumnValue = make(map[string]string, 0)
						t.ID = int64(_id)
						for k, v := range ts.Field {
							if bs, ok := tst.Fields.Get(k); ok {
								t.ColumnValue[k] = type2value(bs, v)
							}
						}
					}
					tst.Fields.Range(func(k string, _ []byte) bool {
						if _, ok := t.ColumnValue[k]; !ok {
							t.ColumnValue[k] = ""
						}
						return true
					})
				}
			}
		} else {
			if name := hc.PostParamTrimSpace("tableName"); name != "" {
				id := hc.PostParamTrimSpace("tableId")
				var _id int
				if _id, err = strconv.Atoi(id); err == nil {
					ts := &TableStub{Tablename: name, ID: int64(_id)}
					if err = Level2.Delete(0, ts); err == nil {
						hc.ResponseString(resultHtmlAndClose("successful!"))
						return
					}
				}
			}
		}
		if err != nil {
			hc.ResponseString(resultHtml("Failed :", err))
			return
		}
	}
	tplToHtml(getLang(hc), DELETE, t, hc)
}

func loadHtml(hc *tlnet.HttpContext) {
	tplToHtml(getLang(hc), LOAD, nil, hc)
}

func monitorHtml(hc *tlnet.HttpContext) {
	tplToHtml(getLang(hc), MONITOR, sys.IsClusRun(), hc)
}

func exportHandler(hc *tlnet.HttpContext) {
	if name := hc.PostParamTrimSpace("exportName"); name != "" {
		if buf, err := sys.Export(name); err == nil {
			hc.Writer().Header().Set("Content-Disposition", "attachment; filename="+name+".gz")
			hc.Writer().Header().Set("Content-Type", "application/octet-stream")
			hc.Writer().Header().Set("Content-Length", fmt.Sprint(buf.Len()))
			hc.Writer().Write(buf.Bytes())
		}
	}
}

/*********************************************************************************/

// 系统变量
func sysvar() (s *SysVar) {
	s = &SysVar{}
	s.StartTime = fmt.Sprint(sys.STARTTIME)
	s.LocalTime = fmt.Sprint(time.Now())
	s.Time = fmt.Sprint(util.Time())
	s.UUID = sys.UUID
	s.RUNUUIDS = fmt.Sprint(sys.GetRunUUID())
	s.STAT = fmt.Sprint(sys.SYS_STAT)
	s.TIME_DEVIATION = fmt.Sprint(sys.TIME_DEVIATION)
	s.ADDR = fmt.Sprint(sys.ADDR)
	s.STORENODENUM = fmt.Sprint(sys.STORENODENUM)
	s.CLUSTER_NUM = fmt.Sprint(sys.CLUSTER_NUM)
	s.MQADDR = sys.MQADDR
	s.CLIADDR = sys.CLIADDR
	s.ADMINADDR = sys.WEBADMINADDR
	s.CCPUT = sys.CcPut()
	s.CCGET = sys.CcGet()
	s.COUNTPUT = sys.CountPut()
	s.COUNTGET = sys.CountGet()
	s.SyncCount = sys.SyncCount()
	return
}

func statStr(stat sys.STATTYPE) (_r string) {
	switch stat {
	case sys.READY:
		_r = "&#23601;&#32490; &#9200;"
	case sys.PROXY:
		_r = "&#20195;&#29702; &#128274;"
	case sys.RUN:
		_r = "&#36816;&#34892; &#9989;"
	}
	return
}

func wsConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !isLogin(hc) {
			hc.WS.Close()
			return
		}
		count := int64(0)
		hc.WS.Send("0")
		i := 5
		for sys.SyncCount() > count || i > 0 {
			count = sys.SyncCount()
			hc.WS.Send(fmt.Sprint(sys.SyncCount()))
			<-time.After(1 * time.Second)
			i--
		}
		hc.WS.Send("")
		hc.WS.Close()
	}
	return
}

func mntConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !isLogin(hc) {
			hc.WS.Close()
			return
		}
	}
	return
}

func mntHandler(hc *tlnet.HttpContext) {
	defer util.Recovr()
	s := string(hc.WS.Read())
	if t, err := strconv.Atoi(s); err == nil {
		if t < 1 {
			t = 1
		}
		for hc.WS.Error == nil {
			if j, err := monitorToJson(); err == nil {
				hc.WS.Send(j)
			}
			<-time.After(time.Duration(t) * time.Second)
		}
	}
}

func initAccount() {
	if len(StoreAdmin.AdminList()) == 0 {
		StoreAdmin.PutAdmin("admin", "123", 1)
	}
	if len(StoreAdmin.ClientList()) == 0 {
		StoreAdmin.PutClient("mycli", "123", 1)
	}
	if len(StoreAdmin.MqList()) == 0 {
		StoreAdmin.PutMq("mymq", "123", 1)
	}
}
