// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package sys

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	. "github.com/donnie4w/gofer/util"
)

func init() {
	Flag()
}

const VERSION = "0.0.3" //program version

const (
	GB = 1 << 30
	MB = 1 << 20
)

var NAMESPACE string              //db namespace
var STARTTIME = time.Now()        //Node startup time (Local time)
var DBSTOREMODE int8 = 1          //storage mode 1hubbed 2sequencing
var DBMode = 1                    //1table mode  2loose mode
var BatchMode int8 = 1            //A ,B mode
var UUID int64                    //local UUID
var ADDR string                   //local address
var DBFILEDIR string              //db file dir
var BINLOGSIZE int64              //size of binlog file
var BINLOGNAME = "binlog.tdb"     //bin log file name
var STATLOGNAME = "stat"          //stat log file name
var BACKLOGNAME = "back"          //back log file name
var CACHELOGNAME = "cache"        //cache log file name
var TLDB_SYS_LOG = "tldb_sys.log" //system log
var ROOTPATHLOG string            //log
var LOGON bool                    //log on
var SYNC bool                     //log on
var DBBUFFER int64                //
var ZLV int                       // zlib compress level  0-9
var TZLIB bool = true
var MEGERCLUSACK = true
var DATAZLIB bool
var MERGETIME int64
var MEMLIMIT int64
var GOGC int
var FREELOCKTIME int64     //maximum retention time of an idle lock
var STATSEQ int64          // stat seq
var MAXDELSEQ int64        // max del seq
var MAXDELSEQCURSOR int64  // max del seq cursor
var STORENODENUM = 1       //number of storage data nodes in the cluster
var CLUSTER_NUM_FINAL bool //fixed number of nodes ,minimum number of nodes in a cluster
var CLUSTER_NUM int        //minimum number of nodes in a cluster
var TIME_DEVIATION int64   //time deviation
var TIME_DEVIATION_LIST = make([]int64, 0)
var PWD string
var PRIVATEKEY string
var PUBLICKEY string
var WSORIGIN string     //mq websocket origin
var CLIADDR string      //client服务地址
var MQADDR string       //mq服务地址
var WEBADMINADDR string //web管理平台
var DEBUGADDR string    //debug pprof
var MQTLS bool          //MQ是否使用tls
var ADMINTLS bool       //web admin是否使用tls
var CLITLS bool         //客户端传输是否使用tls
var CLICRT string       //cli crt文件地址
var CLIKEY string       //cli key文件地址
var MQCRT string        //mq crt文件地址
var MQKEY string        //mq  key文件地址
var ADMINCRT string     //admin crt文件地址
var ADMINKEY string     //admin key文件地址
var TLDBJSON string
var Conf *ConfBean
var defaultConf string

// ////////////////////////////////////////////
var WaitTimeout = 30 * time.Second
var TransTimeout = 30 * time.Second
var ConnectTimeout = 10 * time.Second
var SocketTimeout = 15 * time.Second
var ReadTimeout = 10 * time.Second

var REDOCONN = 1000
var GOMAXLIMIT int64

// ////////////////////////////////////////////
var COCURRENT_PUT int64
var COCURRENT_GET int64

// ////////////////////////////////////////////
var CMD bool
var INIT bool

// ////////////////////////////////////////////
func Flag() {
	flag.StringVar(&PWD, "pwd", "gotldb2023", "password of cluster service node")
	flag.StringVar(&PUBLICKEY, "publickey", "", "file of RSA PUBLIC KEY")
	flag.StringVar(&PRIVATEKEY, "privatekey", "", "file of RSA PRIVATE KEY")
	flag.StringVar(&NAMESPACE, "ns", "tldb", "NAMESPACE (default 'tldb')")
	flag.StringVar(&DBFILEDIR, "dir", "_data", "directory path of database files")
	flag.StringVar(&ROOTPATHLOG, "logdir", "", "directory path of log files")
	flag.Int64Var(&BINLOGSIZE, "binsize", 1<<10, "file size of binlog,unit(MB)")
	flag.BoolVar(&LOGON, "log", false, "debug log on or off")
	flag.BoolVar(&SYNC, "sync", false, "sync data to disk")
	flag.Int64Var(&DBBUFFER, "buffer", 1<<6, "allot buffer (MB) for db")
	flag.Int64Var(&MERGETIME, "mt", 1000, "dwell time of mq merge(default 1000 microsecond)")
	flag.IntVar(&ZLV, "zlv", 9, "zlib compress level 0-9 default 9")
	flag.BoolVar(&TZLIB, "tzlib", true, "transmission use compression")
	flag.BoolVar(&DATAZLIB, "dz", false, "compress data")
	flag.StringVar(&DEBUGADDR, "debug", "", "debug addr")
	flag.DurationVar(&WaitTimeout, "wait", 30, "waitTimeout unit(second)")
	flag.DurationVar(&TransTimeout, "trans", 30, "transTimeout unit(second)")
	flag.DurationVar(&ConnectTimeout, "connect", 10, "connecttimeout unit(second)")
	flag.DurationVar(&SocketTimeout, "socket", 15, "sockettimeout unit(second)")
	flag.DurationVar(&ReadTimeout, "read", 10, "readtimeout unit(second)")

	flag.Int64Var(&COCURRENT_PUT, "put", 500, "maximum number of cocurrent put")
	flag.Int64Var(&COCURRENT_GET, "get", 200, "maximum number of cocurrent get")
	flag.Int64Var(&GOMAXLIMIT, "go", 1<<9, "maximum number of goroutine")
	flag.Int64Var(&FREELOCKTIME, "freelock", 1<<16, "maximum retention time of an idle lock")

	flag.BoolVar(&CLUSTER_NUM_FINAL, "clus_final", false, "if true,'clus' cannot be assigned automatically")
	flag.IntVar(&CLUSTER_NUM, "clus", 0, "minimum number of cluster nodes")
	flag.IntVar(&STORENODENUM, "store", 0, "number of store data node,If the value is 0, all nodes store data")
	flag.Int64Var(&MEMLIMIT, "memlimit", 1580, "memory limit(unit:MB)")
	flag.IntVar(&GOGC, "gc", -1, "a collection is triggered when the ratio of freshly allocated data")

	flag.BoolVar(&CLITLS, "clitls", false, "use the TLS secure transport protocol for client")
	flag.BoolVar(&ADMINTLS, "admintls", false, "use the TLS secure transport protocol for web admin")
	flag.BoolVar(&MQTLS, "mqtls", false, "use the TLS secure transport protocol for mq")

	flag.StringVar(&WSORIGIN, "origin", "http://tldb-mq", "mq websocket origin")
	flag.StringVar(&CLICRT, "clicrt", "", "path of client tls crt file")
	flag.StringVar(&CLIKEY, "clikey", "", "path of client tls key file")
	flag.StringVar(&MQCRT, "mqcrt", "", "path of mq tls crt file")
	flag.StringVar(&MQKEY, "mqkey", "", "path of mq  tls key file")
	flag.StringVar(&ADMINCRT, "admincrt", "", "path of admin tls crt file")
	flag.StringVar(&ADMINKEY, "adminkey", "", "path of admin tls key file")

	flag.StringVar(&CLIADDR, "cli", "", "client address")
	flag.StringVar(&MQADDR, "mq", "", "mq address")
	flag.StringVar(&ADDR, "cs", "", "cluster service address")
	flag.StringVar(&WEBADMINADDR, "admin", ":4001", "web admin platform address")
	flag.StringVar(&TLDBJSON, "c", "tldb.json", "config file of tldb")

	flag.BoolVar(&CMD, "cmd", false, "command line interaction mode")
	flag.BoolVar(&INIT, "init", false, "init and create default account")
	flag.Usage = usage
	flag.Parse()
	flag.Usage()

	WaitTimeout = WaitTimeout * time.Second
	TransTimeout = TransTimeout * time.Second
	ConnectTimeout = ConnectTimeout * time.Second
	SocketTimeout = SocketTimeout * time.Second
	ReadTimeout = ReadTimeout * time.Second
	DBBUFFER = DBBUFFER * MB
	parsec()
	if ZLV > 9 {
		ZLV = 9
	} else if ZLV < 0 {
		ZLV = -1
	}
	if Conf.PWD != nil {
		PWD = *Conf.PWD
	}
	if Conf.PUBLICKEY != nil && Conf.PRIVATEKEY != nil {
		PUBLICKEY = *Conf.PUBLICKEY
		PRIVATEKEY = *Conf.PRIVATEKEY
	}
	if Conf.NAMESPACE != nil {
		NAMESPACE = *Conf.NAMESPACE
	}
	if Conf.DBFILEDIR != nil {
		DBFILEDIR = *Conf.DBFILEDIR
	}
	if Conf.ROOTPATHLOG != nil {
		ROOTPATHLOG = *Conf.ROOTPATHLOG
	}
	if Conf.BINLOGSIZE > 0 {
		BINLOGSIZE = Conf.BINLOGSIZE
	}
	if Conf.CLUSTER_NUM_FINAL {
		CLUSTER_NUM_FINAL = Conf.CLUSTER_NUM_FINAL
	}
	if Conf.CLUSTER_NUM > 0 {
		CLUSTER_NUM = Conf.CLUSTER_NUM
	}
	if Conf.STORENODENUM > 0 {
		STORENODENUM = Conf.STORENODENUM
	}
	if Conf.MEMLIMIT > 0 {
		MEMLIMIT = Conf.MEMLIMIT
	}
	if Conf.CLITLS {
		CLITLS = Conf.CLITLS
	}
	if Conf.GOGC != nil {
		GOGC = *Conf.GOGC
	}
	if Conf.ADMINTLS {
		ADMINTLS = Conf.ADMINTLS
	}
	if Conf.MQTLS {
		MQTLS = Conf.MQTLS
	}
	if Conf.WSORIGIN != nil {
		WSORIGIN = *Conf.WSORIGIN
	}
	if Conf.CLICRT != nil {
		CLICRT = *Conf.CLICRT
	}
	if Conf.CLIKEY != nil {
		CLIKEY = *Conf.CLIKEY
	}
	if Conf.MQCRT != nil {
		MQCRT = *Conf.MQCRT
	}
	if Conf.MQKEY != nil {
		MQKEY = *Conf.MQKEY
	}
	if Conf.ADMINCRT != nil {
		ADMINCRT = *Conf.ADMINCRT
	}
	if Conf.CLIADDR != nil {
		CLIADDR = *Conf.CLIADDR
	}
	if Conf.MQADDR != nil {
		MQADDR = *Conf.MQADDR
	}
	if Conf.ADDR != nil {
		ADDR = *Conf.ADDR
	}
	if Conf.WEBADMINADDR != nil {
		WEBADMINADDR = *Conf.WEBADMINADDR
	}
	debug.SetMemoryLimit(MEMLIMIT * MB)
	debug.SetGCPercent(GOGC)
}

type ConfBean struct {
	PWD               *string `json:"pwd"`
	PUBLICKEY         *string `json:"publickey"`
	PRIVATEKEY        *string `json:"privatekey"`
	NAMESPACE         *string `json:"namespace"`
	DBFILEDIR         *string `json:"dir"`
	ROOTPATHLOG       *string `json:"logdir"`
	BINLOGSIZE        int64   `json:"binsize"`
	CLUSTER_NUM_FINAL bool    `json:"clus_final"`
	CLUSTER_NUM       int     `json:"clus"`
	GOGC              *int    `json:"gc"`
	STORENODENUM      int     `json:"store"`
	MEMLIMIT          int64   `json:"memlimit"`
	CLITLS            bool    `json:"clitls"`
	ADMINTLS          bool    `json:"admintls"`
	MQTLS             bool    `json:"mqtls"`
	WSORIGIN          *string `json:"origin"`
	CLICRT            *string `json:"clicrt"`
	CLIKEY            *string `json:"clikey"`
	MQCRT             *string `json:"mqcrt"`
	MQKEY             *string `json:"mqkey"`
	ADMINCRT          *string `json:"admincrt"`
	CLIADDR           *string `json:"cli"`
	MQADDR            *string `json:"mq"`
	ADDR              *string `json:"cs"`
	WEBADMINADDR      *string `json:"admin"`
}

func parsec() {
	if defaultConf != "" {
		Conf, _ = JsonDecode[*ConfBean]([]byte(defaultConf))
	} else if bs, err := ReadFile(TLDBJSON); err == nil {
		Conf, _ = JsonDecode[*ConfBean](bs)
	}
	if Conf == nil {
		Conf = &ConfBean{}
	}
}

func usage() {
	exename := "tldb"
	if runtime.GOOS == "windows" {
		exename = "tldb.exe"
	}
	fmt.Fprintln(os.Stderr, `tldb version: `+VERSION+`
Server Mode  Usage: `+exename+` [-dir data path] [-cs cluster service addr] [-cli client service addr] [-admin admin addr] [-mq mq service addr]
Command Mode Usage: `+exename+` [-dir data path] [-cmd]`)
}
