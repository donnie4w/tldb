// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package tc

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/donnie4w/tldb/container"
	. "github.com/donnie4w/tldb/keystore"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/tnet"
	"github.com/donnie4w/tldb/util"
	"github.com/donnie4w/tlnet"
)

var cmdMap = NewMap[int64, *tln]()

func localHandler(hc *tlnet.HttpContext) {
	auth := hc.WS.Read()
	if v, ok := StoreAdmin.GetOther("admintauth"); ok && v == string(auth) {
		defer recover()
		tl := newTln()
		cmdMap.Put(hc.WS.Id, tl)
		hc.WS.Send(tl.port)
	}
}

func localwsConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.Origin = "http://tldb-admin"
	wc.OnError = func(self *tlnet.Websocket) {
		fmt.Println("command over:", util.TimeStr())
		defer recover()
		if tln, ok := cmdMap.Get(self.Id); ok {
			cmdMap.Del(self.Id)
			tln.close()
		}
	}
	return
}

type tln struct {
	portc   chan string
	runPort chan string
	tt      *tlnet.Tlnet
	port    string
	isClose bool
}

func newTln() (tl *tln) {
	tl = &tln{portc: make(chan string, 0), runPort: make(chan string, 0), tt: tlnet.NewTlnet()}
	go tl.newTlnet()
	go tl.ping()
	tl.port = <-tl.runPort
	return
}

func (this *tln) newTlnet() {
	port := fmt.Sprint(util.Rand(15535) + 50000)
	go func() {
		<-time.After(time.Second)
		this.portc <- port
	}()
	this.tt.GET("/test", testHandler)
	this.tt.HandleWebSocketBindConfig("/", cmdHandler, cmdWsConfig())
	if err := this.tt.HttpStart(fmt.Sprint("127.0.0.1:", port)); err != nil && !this.isClose {
		this.newTlnet()
	}
}

func (this *tln) ping() {
	var port string
	for {
		port = <-this.portc
		if port != "" {
			if resp, err := http.Get(fmt.Sprint("http://127.0.0.1:", port, "/test")); err == nil {
				defer resp.Body.Close()
				bs := make([]byte, 1)
				resp.Body.Read(bs)
				if bs[0] == 1 {
					this.runPort <- port
					break
				}
			}
		}
		<-time.After(time.Second)
	}
	return
}

func (this *tln) close() {
	defer recover()
	if !this.isClose {
		this.isClose = true
		this.tt.Close()
	}
}

func testHandler(hc *tlnet.HttpContext) {
	hc.ResponseBytes(0, []byte{1})
}

func cmdWsConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.Origin = "http://tldb-cmd"
	wc.OnOpen = func(_ *tlnet.HttpContext) {
		fmt.Println("command:", util.TimeStr())
	}
	return
}

func cmdHandler(hc *tlnet.HttpContext) {
	cmdstr := string(hc.WS.Read())
	hc.WS.Send(parseCmd(cmdstr, hc))
}

func parseCmd(cmdstr string, hc *tlnet.HttpContext) (_r string) {
	defer func() {
		if err := recover(); err != nil {
			_r = failed("invalid command: " + cmdstr)
		}
	}()
	ss := strings.Split(cmdstr, " ")
	switch ss[0] {
	case "addnode":
		if err := tnet.AddNode(ss[1]); err == nil {
			_r = "addnode:[" + ss[1] + "]"
			_r = success(_r)
		} else {
			_r = "[" + ss[1] + "]" + err.Error()
			_r = failed(_r)
		}
	case "load", "loadforce":
		path := replaceFilePath(ss[1])
		datetime, limit := "", int64(0)
		if len(ss) >= 4 {
			datetime = ss[2] + " " + ss[3]
			if _, err := util.TimeStrToTime(datetime); err != nil {
				_r = failed(err.Error())
				return
			}
		}
		if len(ss) == 5 {
			if i, er := strconv.Atoi(ss[4]); er == nil {
				limit = int64(i)
			}
		}
		var err error
		var bs []byte
		if bs, err = util.ReadFile(path); err == nil {
			err = util.CheckGzip(bs)
			if ss[0] == "load" {
				go sys.LoadData2TLDB(bs, datetime)
			} else {
				go sys.ForcedCoverageData2TLDB(bs, datetime, limit)
			}
		}
		if err == nil {
			count := int64(0)
			i := 5
			c := 1
			for sys.SyncCount() > count || i > 0 {
				count = sys.SyncCount()
				hc.WS.Send(".")
				<-time.After(20 * time.Millisecond)
				c++
				if c%50 == 0 {
					i--
					hc.WS.Send(fmt.Sprint("\n", sys.SyncCount()))
				}
			}
			hc.WS.Send(fmt.Sprint("\n", sys.SyncCount()))
			hc.WS.Send("\n")
			_r = success("load:" + ss[1])
		} else {
			_r = failed(err.Error())
		}
	case "pwd", "add":
		ok := false
		if len(ss) == 3 {
			StoreAdmin.PutAdmin(ss[1], ss[2], 1)
			ok = true
		} else if len(ss) == 4 {
			if _ty, err := strconv.Atoi(ss[3]); err == nil && (_ty == 1 || _ty == 2) {
				StoreAdmin.PutAdmin(ss[1], ss[2], int8(_ty))
				ok = true
			} else {
				_r = err.Error()
			}
		}
		if ok {
			_r = success("pwd :" + ss[1])
		} else {
			_r = failed("pwd :" + ss[1])
		}
	case "cli", "pwdcli", "addcli":
		if len(ss) == 3 {
			StoreAdmin.PutClient(ss[1], ss[2], 1)
			_r = success("addcli :" + ss[1])
		} else {
			_r = failed("addcli :" + ss[1])
		}
	case "mq", "pwdmq", "addmq":
		if len(ss) == 3 {
			StoreAdmin.PutMq(ss[1], ss[2], 1)
			_r = success("addmq :" + ss[1])
		} else {
			_r = failed("addmq :" + ss[1])
		}

	case "delmq":
		if len(ss) == 2 {
			StoreAdmin.DelMq(ss[1])
			_r = success("delmq :" + ss[1])
		} else {
			_r = failed("delmq")
		}
	case "delcli":
		if len(ss) == 2 {
			StoreAdmin.DelClient(ss[1])
			_r = success("delcli :" + ss[1])
		} else {
			_r = failed("delcli")
		}
	case "del":
		if len(ss) == 2 {
			StoreAdmin.DelAdmin(ss[1])
			sessionMap.Range(func(k string, v *UserBean) bool {
				if v.Name == ss[1] {
					sessionMap.Del(k)
				}
				return true
			})
			_r = success("del :" + ss[1])
		} else {
			_r = failed("del")
		}
	case "close":
		go hc.WS.Send("exit")
		tnet.CloseSelf()
		sys.Stop()
	case "monitor":
		interval := 3
		count := 10
		if len(ss) >= 2 {
			if i, err := strconv.Atoi(ss[1]); err == nil {
				interval = i
			}
		}
		if len(ss) >= 3 {
			if i, err := strconv.Atoi(ss[2]); err == nil {
				count = i
			}
		}
		hc.WS.Send(fmt.Sprintln(">> >> >>[monitor  ", interval, "  ", count, "] Start"))
		s := `Alloc(MB)  TotalAlloc(MB)  NumGC  NumGoroutine  ConcurrencyCUD  ConcurrencyR`
		hc.WS.Send(fmt.Sprintln(s))
		for count > 0 {
			count--
			m := getMonitor()
			s := fmt.Sprintln(m.Alloc/sys.MB, "		", m.TotalAlloc/sys.MB, "	   ", m.NumGC, "	   ", m.NumGoroutine, "		", m.CcPut, "		", m.CcGet)
			hc.WS.Send(s)
			<-time.After(time.Duration(interval) * time.Second)
		}
		hc.WS.Send(fmt.Sprintln(">>>End"))
	case "init":
		initAccount()
		_r = success("init")
	case "export":
		var err error
		if len(ss) >= 2 {
			name := ss[1]
			var buf bytes.Buffer
			if buf, err = sys.Export(name); err == nil {
				hc.WS.Send(buf.Bytes())
			}
		} else {
			err = errors.New("param error")
		}
		if err == nil {
			_r = success("export:" + ss[1])
		} else {
			_r = failed("export: " + err.Error())
		}
	default:
		_r = failed("invalid command: " + cmdstr)
	}
	return
}

func success(s string) (_r string) {
	_r = fmt.Sprintln(">>>success! ", s)
	return
}

func failed(s string) (_r string) {
	_r = fmt.Sprintln(">>>failed! ", s)
	return
}

func replaceFilePath(path string) (_r string) {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.ReplaceAll(path, `\`, "/")
	path = strings.ReplaceAll(path, "\t", "")
	path = strings.ReplaceAll(path, "\n", "")
	path = strings.ReplaceAll(path, "\r", "")
	path = strings.ReplaceAll(path, `"`, "")
	return path
}
