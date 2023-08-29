// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package sys

import (
	"bytes"
	"os"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/stub"
)

// ///////////////////////////////////////////
type STATTYPE int8   //
var SYS_STAT = READY //
const (
	READY STATTYPE = 0 //就绪状态
	PROXY STATTYPE = 1 //代理状态
	RUN   STATTYPE = 2 //运行状态
) //
func IsREADY() bool { //
	return SYS_STAT == READY
}
func IsPROXY() bool {
	return SYS_STAT == PROXY
}
func IsRUN() bool {
	return SYS_STAT == RUN
}

func IsStandAlone() bool {
	return CLUSTER_NUM == 0
}

var _stat_seq = int64(1) //
func SetStat(stat STATTYPE, timenano time.Duration) {
	if timenano > 0 {
		seq := atomic.AddInt64(&_stat_seq, 1)
		go func() {
			<-time.NewTimer(timenano).C
			if _stat_seq == seq {
				SYS_STAT = stat
			}
		}()
	} else {
		SYS_STAT = stat
	}
}

// ///////////////////////////////////////////
var Service = NewSortMap[int, stub.Server]()

func Start() {
	if CMD {
		Cmd()
	} else {
		var wg sync.WaitGroup
		Service.BackForEach(func(_ int, s stub.Server) bool {
			wg.Add(1)
			go func() {
				s.Serve(&wg)
			}()
			<-time.After(500 * time.Millisecond)
			return true
		})
		wg.Wait()
	}
}

func Stop() {
	Service.FrontForEach(func(_ int, s stub.Server) bool {
		s.Close()
		return true
	})
	os.Exit(0)
}

/*****************************************************/
var IsClusRun func() bool
var GetRemoteNode func() []*stub.RemoteNode
var GetRunUUID func() []int64
var SyncCount func() int64
var TryRunToProxy func() (err error)
var TryProxyToReady func() (err error)
var ReSetStoreNodeNumber func(num int32) (err error)
var LoadData2TLDB func(bs []byte, datetime string) (err error)
var ForcedCoverageData2TLDB func(bs []byte, datetime string, limit int64) (err error)
var Export func(tablename string) (_r bytes.Buffer, err error)
var Level0Put func(key string, value []byte) (err error)
var Level0Get func(key string) (value []byte, err error)
var Client2Serve func(addr string) (err error)
var BroadRmNode func() (err error)
var CcGet func() int64
var CcPut func() int64
var CountPut func() int64
var CountGet func() int64
var Cmd func()

/*****************************************************/
