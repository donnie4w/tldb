// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tlcli

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/donnie4w/gothrift/thrift"

	// thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/tldb/keystore"
	"github.com/donnie4w/tldb/log"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

func init() {
	sys.Service.Put(4, cliservice)
}

var logger = log.Logger
var _transportFactory = thrift.NewTBufferedTransportFactory(1 << 13)
var _tcompactProtocolFactory = thrift.NewTCompactProtocolFactoryConf(&thrift.TConfiguration{})

var cliservice = &cliService{}

type cliService struct {
	isClose bool
	server  thrift.TServerTransport
	connNum int32
}

func (this *cliService) _server(wg *sync.WaitGroup, _addr string, processor thrift.TProcessor, TLS bool, serverCrt, serverKey string) (err error) {
	defer wg.Done()
	var serverTransport thrift.TServerTransport
	if TLS {
		cfg := &tls.Config{}
		var cert tls.Certificate
		if util.IsFileExist(serverCrt) && util.IsFileExist(serverKey) {
			cert, err = tls.LoadX509KeyPair(serverCrt, serverKey)
		} else {
			cert, err = tls.X509KeyPair([]byte(keystore.ServerCrt), []byte(keystore.ServerKey))
		}
		if err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
			serverTransport, err = thrift.NewTSSLServerSocketTimeout(_addr, cfg, sys.SocketTimeout)
		}
	} else if serverTransport, err = thrift.NewTServerSocketTimeout(_addr, sys.SocketTimeout); err == nil {
		this.server = serverTransport
	}
	if err == nil && serverTransport != nil {
		server := thrift.NewTSimpleServer4(processor, serverTransport, nil, nil)
		if err = server.Listen(); err == nil {
			s := fmt.Sprint("tlClient start[", _addr, "]")
			if TLS {
				s = fmt.Sprint("tlClient start tls[", _addr, "]")
			}
			log.LoggerSys.Info(sys.SysLog(s))
			for {
				if _transport, err := server.ServerTransport().Accept(); err == nil {
					atomic.AddInt32(&this.connNum, 1)
					go func() {
						defer func() { atomic.AddInt32(&this.connNum, -1) }()
						defer myRecovr()
						cc := newCliContext(_transport)
						defer cc.close()
						defaultCtx := context.WithValue(context.Background(), "CliContext", cc)
						if inputTransport, err := _transportFactory.GetTransport(_transport); err == nil {
							inputProtocol := _tcompactProtocolFactory.GetProtocol(inputTransport)
							for {
								ok, err := processor.Process(defaultCtx, inputProtocol, inputProtocol)
								if errors.Is(err, thrift.ErrAbandonRequest) {
									break
								}
								if errors.As(err, new(thrift.TTransportException)) && err != nil {
									break
								}
								if !ok {
									break
								}
							}
						}
					}()
				}
			}
		}
	}
	if !this.isClose && err != nil {
		log.LoggerSys.Error("tlClient start failed:", err)
		panic("tlClient start failed:" + err.Error())
	}
	return
}

func (this *cliService) Close() (err error) {
	defer util.Recovr()
	if strings.TrimSpace(sys.CLIADDR) != "" {
		this.isClose = true
		err = this.server.Close()
	}
	return
}

func (this *cliService) Serve(wg *sync.WaitGroup) (err error) {
	if strings.TrimSpace(sys.CLIADDR) != "" {
		err = this._server(wg, strings.TrimSpace(sys.CLIADDR), NewIcliProcessor(cliProcessor), sys.CLITLS, sys.CLICRT, sys.CLIKEY)
	} else {
		wg.Done()
	}
	return
}

type cliContext struct {
	Id       int64
	isAuth   bool
	tt       thrift.TTransport
	mux      *sync.Mutex
	_isClose bool
}

func newCliContext(tt thrift.TTransport) (cc *cliContext) {
	cc = &cliContext{int64(util.NewTxId()), false, tt, &sync.Mutex{}, false}
	return
}

func (this *cliContext) close() {
	defer myRecovr()
	defer this.mux.Unlock()
	this.mux.Lock()
	if !this._isClose {
		this._isClose = true
		this.tt.Close()
	}
}

func myRecovr() {
	if err := recover(); err != nil {
		logger.Error(string(debug.Stack()))
	}
}
