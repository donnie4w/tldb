// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/tldb/log"
	. "github.com/donnie4w/tldb/util"
)

type tnetServer struct {
	servermux *serverMux
}

func (this *tnetServer) handle(processor Itnet, handler func(tc *tlContext), cliError func(tc *tlContext)) {
	if this.servermux == nil {
		this.servermux = &serverMux{}
		this.servermux.Handle(processor, handler, cliError)
	}
}

func (this *tnetServer) Serve(wg *sync.WaitGroup, _addr string) (err error) {
	if !strings.Contains(_addr, ":") {
		if MatchString("^[0-9]{4,5}$", _addr) {
			_addr = ":" + _addr
		} else {
			log.LoggerSys.Error("address error:" + _addr)
			err = errors.New("address format error:" + _addr + "| e.g. :4646 | localhost:4646 | 4646")
		}
	}
	// err = tlclientServer.server(wg, _addr, this.servermux.processor, this.servermux.handler, this.servermux.cliError)
	err = tfsclientserver.server(wg, _addr, this.servermux.processor, this.servermux.handler, this.servermux.cliError)
	return
}

func (this *tnetServer) Connect(_addr string, async bool) (err error) {
	if !strings.Contains(_addr, ":") && MatchString("^[0-9]{4,5}$", _addr) {
		_addr = ":" + _addr
	}
	// return tlserverclient.server(_addr, this.servermux.processor, this.servermux.handler, this.servermux.cliError, async)
	return tfsserverclient.server(_addr, this.servermux.processor, this.servermux.handler, this.servermux.cliError, async)
}

type serverMux struct {
	processor Itnet
	handler   func(tc *tlContext)
	cliError  func(tc *tlContext)
}

func (this *serverMux) Handle(processor Itnet, handler func(tc *tlContext), cliError func(tc *tlContext)) {
	this.processor = processor
	this.handler = handler
	this.cliError = cliError
}

func tlProcessor() (processor thrift.TProcessor) {
	handler := &ItnetServ{NewNumLock(64)}
	processor = NewItnetProcessor(handler)
	return
}

func myServer2ClientHandler(tc *tlContext) {
}

func myClient2ServerHandler(tc *tlContext) {
	defer myRecovr()
	logger.Info("tc.isClose>>", tc.isClose)
	if at, err := getAuth(); err == nil && !tc.isClose {
		if err := tc.Conn.Auth(context.Background(), at); err != nil {
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}
}

func mySecvErrorHandler(tc *tlContext) {
	go reconn(tc)
}

func myCliErrorHandler(tc *tlContext) {
	clientLinkCache.Del(tc.RemoteAddr)
	go reconn(tc)
}

/***********************************************************************************/
func reconn(tc *tlContext) {
	if tc == nil || tc.RemoteAddr == "" || tc.RemoteUuid == 0 {
		return
	}
	defer myRecovr()
	// defer tc.mux.Unlock()
	// tc.mux.Lock()
	logger.Warn("reconn[", tc.RemoteUuid, "][", tc.RemoteAddr, "]")
	nodeWare.del(tc)
	checkAndResetStat()
	if !tc.isServer && !tc._do_reconn {
		tc._do_reconn = true
		i := 0
		for !delNodeMap.Has(tc.RemoteAddr) && !nodeWare.hasUUID(tc.RemoteUuid) {
			if clientLinkCache.Has(tc.RemoteAddr) {
				<-time.After(time.Duration(Rand(6)) * time.Second)
			} else {
				if err1, err2 := tnetservice.Connect(tc.RemoteAddr, false); err1 != nil {
					break
				} else if err2 != nil {
					if i < 100 {
						i++
					}
					<-time.After(time.Duration(Rand(6+i)) * time.Second)
				}
			}
		}
	}
}

func heardbeat() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			<-time.After(time.Duration(Rand(5)) * time.Second)
			_heardbeat(nodeWare.GetAllTlContext())
		}
	}
}

func _heardbeat(tcs []*tlContext) {
	for _, tc := range tcs {
		func(tc *tlContext) {
			defer myRecovr()
			tc.Conn.Ping(context.Background(), pistr())
			if atomic.AddInt64(&tc.pingNum, 1) > 5 {
				logError.Error("ping failed:[", tc.RemoteUuid, "][", tc.RemoteAddr, "] ping number:", tc.pingNum)
				go reconn(tc)
			}
		}(tc)
	}
}
