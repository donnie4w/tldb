// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package tlcli

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/donnie4w/gothrift/thrift"
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
	isClose         bool
	serverTransport thrift.TServerTransport
	connNum         int32
}

func (this *cliService) _server(_addr string, processor thrift.TProcessor, TLS bool, serverCrt, serverKey string) (err error) {
	// var serverTransport thrift.TServerTransport
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
			this.serverTransport, err = thrift.NewTSSLServerSocketTimeout(_addr, cfg, sys.SocketTimeout)
		}
	} else {
		this.serverTransport, err = thrift.NewTServerSocketTimeout(_addr, sys.SocketTimeout)
	}
	if err == nil && this.serverTransport != nil {
		server := thrift.NewTSimpleServer4(processor, this.serverTransport, nil, nil)
		if err = server.Listen(); err == nil {
			s := fmt.Sprint("Tldb Client service start[", _addr, "]")
			if TLS {
				s = fmt.Sprint("Tldb Client service start tls[", _addr, "]")
			}
			sys.FmtLog(s)
			for {
				if _transport, err := server.ServerTransport().Accept(); err == nil {
					atomic.AddInt32(&this.connNum, 1)
					go func() {
						defer func() { atomic.AddInt32(&this.connNum, -1) }()
						defer myRecover()
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
		sys.FmtLog("Tldb Client service init failed:", err)
		os.Exit(0)
	}
	return
}

func (this *cliService) Close() (err error) {
	defer util.Recovr()
	if strings.TrimSpace(sys.CLIADDR) != "" {
		this.isClose = true
		err = this.serverTransport.Close()
	}
	return
}

func (this *cliService) Serve() (err error) {
	if strings.TrimSpace(sys.CLIADDR) != "" {
		err = this._server(strings.TrimSpace(sys.CLIADDR), NewIcliProcessor(cliProcessor), sys.CLITLS, sys.CLICRT, sys.CLIKEY)
	} else {
		sys.FmtLog("no client service")
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
	defer myRecover()
	defer this.mux.Unlock()
	this.mux.Lock()
	if !this._isClose {
		this._isClose = true
		this.tt.Close()
	}
}

func myRecover() {
	if err := recover(); err != nil {
		logger.Error(err)
	}
}
