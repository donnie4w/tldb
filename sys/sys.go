// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package sys

import (
	"bytes"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/donnie4w/gofer/gosignal"
	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/stub"
)

// ///////////////////////////////////////////
type STATTYPE int8

var SYS_STAT = READY

const (
	READY STATTYPE = 0
	PROXY STATTYPE = 1
	RUN   STATTYPE = 2
)

func IsREADY() bool {
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

var _stat_seq = int64(1)

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
		timlogo()
		addStopEvent()
		Service.BackForEach(func(_ int, s stub.Server) bool {
			go func() {
				defer func() { recover() }()
				s.Serve()
			}()
			<-time.After(500 * time.Millisecond)
			return true
		})
		select {}
	}
}

func Stop() {
	defer os.Exit(0)
	Service.FrontForEach(func(_ int, s stub.Server) bool {
		s.Close()
		return true
	})
}

func addStopEvent() {
	gosignal.ListenSignalEvent(func(sig os.Signal) {
		Stop()
	}, syscall.SIGTERM, syscall.SIGINT)
}

var (
	IsClusRun               func() bool
	GetRemoteNode           func() []*stub.RemoteNode
	GetRunUUID              func() []int64
	SyncCount               func() int64
	TryRunToProxy           func() error
	TryProxyToReady         func() error
	ReSetStoreNodeNumber    func(int32) error
	LoadData2TLDB           func([]byte, string) error
	ForcedCoverageData2TLDB func([]byte, string, int64) error
	Export                  func(string) (bytes.Buffer, error)
	Level0Put               func(string, []byte) error
	Level0Get               func(string) ([]byte, error)
	Client2Serve            func(string) error
	BroadRmNode             func() error
	ForceUnLock             func(string)
	Lock                    func(int64, string)
	ReqToken                func(int64, string) error
	TryLock                 func(string) (string, bool)
	UnLock                  func(string, string) bool
	CcGet                   func() int64
	CcPut                   func() int64
	CountPut                func() int64
	CountGet                func() int64
	Cmd                     func()
	Log                     func() *logging.Logging
)
