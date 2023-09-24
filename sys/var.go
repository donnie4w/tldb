// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package sys

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

func init() {
	Flag()
}

const VERSION = "0.0.2" //program version

const GB = 1 << 30
const MB = 1 << 20

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

// /////////////////////////////////////////////
var STATSEQ int64                          // stat seq
var MAXDELSEQ int64                        // max del seq
var MAXDELSEQCURSOR int64                  // max del seq cursor
var STORENODENUM = 1                       //number of storage data nodes in the cluster
var CLUSTER_NUM_FINAL bool                 //fixed number of nodes ,minimum number of nodes in a cluster
var CLUSTER_NUM int                        //minimum number of nodes in a cluster
var TIME_DEVIATION int64                   //time deviation
var TIME_DEVIATION_LIST = make([]int64, 0) //
var PWD string
var PRIVATEKEY string
var PUBLICKEY string

// /////////////////////////////////////////////
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

	flag.BoolVar(&CLUSTER_NUM_FINAL, "clus_final", false, "if true,'clus' cannot be assigned automatically")
	flag.IntVar(&CLUSTER_NUM, "clus", 3, "minimum number of cluster nodes")
	flag.IntVar(&STORENODENUM, "store", 0, "number of store data node,If the value is 0, all nodes store data")
	flag.Int64Var(&MEMLIMIT, "memlimit", 1280, "memory limit(unit:MB)")
	flag.IntVar(&GOGC, "gogc", -1, "a collection is triggered when the ratio of freshly allocated data")

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

	flag.StringVar(&CLIADDR, "cli", ":7001", "client address")
	flag.StringVar(&MQADDR, "mq", ":5001", "mq address")
	flag.StringVar(&ADDR, "cs", ":6001", "cluster service address")
	flag.StringVar(&WEBADMINADDR, "admin", ":4001", "web admin platform address")

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
	if ZLV > 9 {
		ZLV = 9
	} else if ZLV < 0 {
		ZLV = -1
	}

	debug.SetMemoryLimit(MEMLIMIT * MB)
	debug.SetGCPercent(GOGC)
}

func usage() {
	exename := "tldb"
	if runtime.GOOS == "windows" {
		exename = "tldb.exe"
	}
	fmt.Fprintln(os.Stderr, `tldb version: tldb/0.0.1
Server Mode  Usage: `+exename+` [-dir data path] [-cs cluster service addr] [-cli client service addr] [-admin admin addr] [-mq mq service addr]
Command Mode Usage: `+exename+` [-dir data path] [-cmd]`)
}
