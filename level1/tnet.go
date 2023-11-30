// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package level1

import (
	"errors"
	"strings"
	"time"

	. "github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/tldb/level0"
	. "github.com/donnie4w/tldb/lock"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

func init() {
	sys.IsClusRun = IsClusRun
	sys.GetRemoteNode = getRemoteNode
	sys.SyncCount = syncCount
	sys.GetRunUUID = getRunUUID
	sys.TryRunToProxy = tryRunToProxy
	sys.TryProxyToReady = tryProxyToReady
	sys.ReSetStoreNodeNumber = reSetStoreNodeNumber
	sys.LoadData2TLDB = loadData2TLDB
	sys.ForcedCoverageData2TLDB = forcedCoverageData2TLDB
	sys.Client2Serve = client2Serve
	sys.Level0Put = level0.Level0.Put
	sys.Level0Get = level0.Level0.Get
	sys.BroadRmNode = broadRmNode
	sys.Service.Put(5, tnetservice)
}

var tnetservice = &tnetService{}

type tnetService struct{}

func (this *tnetService) Serve() (err error) {
	synctx.init()
	if !sys.IsStandAlone() {
		this._serve(sys.ADDR)
	} 
	return
}

func (this *tnetService) _serve(addr string) {
	sys.FmtLog("service start[", addr, "] ns[", sys.NAMESPACE, "]")
	tnetserver := new(tnetServer)
	tnetserver.handle(&ItnetServ{NewNumLock(64)}, myServer2ClientHandler, mySecvErrorHandler)
	go heardbeat()
	tnetserver.Serve(addr)
}

func (this *tnetService) Connect(addr string, async bool) (err1, err2 error) {
	if nodeWare.LenByAddr(addr) > 0 {
		err1 = errors.New("addr:" + addr + " already existed")
	}
	if addr == sys.ADDR {
		err1 = errors.New("addr:" + addr + " is local addr")
	}
	if addr = strings.Trim(addr, " "); addr == "" {
		err1 = errors.New("addr is bad")
	}
	if err1 != nil {
		return
	}
	logger.Info("connect >> ", addr)
	tnetserver := new(tnetServer)
	tnetserver.handle(&ItnetServ{NewNumLock(64)}, myClient2ServerHandler, myCliErrorHandler)
	err2 = tnetserver.Connect(addr, async)
	return
}

func (this *tnetService) Close() (err error) {
	if !sys.IsStandAlone() {
		// tlclientServer.close()
		tfsclientserver.close()
		nodeWare.close()
	}
	return
}

func pubTxSyncWithTimeOut(pb *PonBean, to time.Duration) *sys.Exception {
	defer errRecover()
	return nodeWare.broadcastSync(pb, PON, to, false)
}

func pubTxSyncAlltime(pb *PonBean) *sys.Exception {
	defer errRecover()
	return nodeWare.broadcastSync(pb, PON, sys.WaitTimeout, true)
}

func pubTxSync(pb *PonBean) *sys.Exception {
	defer errRecover()
	return nodeWare.broadcastSync(pb, PON, sys.WaitTimeout, false)
}

func pubTxAsync(pb *PonBean) {
	defer errRecover()
	nodeWare.broadcast(pb, PON, sys.WaitTimeout, false)
}

func pubTxByRunAsync(pb *PonBean, timeout time.Duration) {
	defer errRecover()
	nodeWare.broadcast(pb, PON, timeout, true)
}

func pubSingleTxAlltime(pb *PonBean, uuid int64, isBack bool) {
	defer errRecover()
	nodeWare.singleBroadAlltime(pb, uuid, PON, 240*time.Hour, isBack)
}

func pubSingleTx(pb *PonBean, uuid int64) error {
	defer errRecover()
	return nodeWare.singleBroad(pb, uuid, PON, sys.WaitTimeout, false)
}

func requestSingleTx(pb *PonBean, uuid int64) (*PonBean, error) {
	defer errRecover()
	return nodeWare.pon3(uuid, 0, false, pb, sys.ReadTimeout)
}

func responeSingleTx(pb *PonBean, uuid, syncId int64) (*PonBean, error) {
	defer errRecover()
	return nodeWare.pon3(uuid, syncId, true, pb, sys.ReadTimeout)
}

func pub2TxSync(pb *PonBean) *sys.Exception {
	defer errRecover()
	return nodeWare.broadcastSync(pb, PON2, sys.WaitTimeout, false)
}

func pub2SingleTx(pb *PonBean, uuid int64) error {
	defer errRecover()
	return nodeWare.singleBroad(pb, uuid, PON2, sys.WaitTimeout, false)
}

func syncCount() int64 {
	return tlog.syncCount
}

func getRunUUID() []int64 {
	return nodeWare.GetAllRunUUID()
}

func removeNode(uuid int64) {
	if addr := nodeWare.GetTlContext(uuid).RemoteAddr; addr != "" {
		delNodeMap.Put(addr, 0)
	}
}

/********************************************************/

func BroadcastMQ(mqType int8, bs []byte) {
	defer errRecover()
	nodeWare.broadcastMQ(mqType, util.SnappyEncode(bs))
}

func BroadcastReinit(sbean *SysBean) (exce *sys.Exception) {
	defer errRecover()
	if sbean == nil {
		sbean = &SysBean{int8(sys.DBMode), int32(sys.STORENODENUM), 0}
	}
	return nodeWare.broadcastReinit(sbean)
}

func BroadcastProxy(tp *TableParam, uuid int64, pType int8, ctype int8) (_r *TableParam, exce *sys.Exception) {
	defer errRecover()
	return nodeWare.broadcastRunNodeProxy(util.TEncode(tp), uuid, pType, ctype)
}

func WriteBatchLogBuffer(bp *BatchPacket) (buf *Buffer) {
	// return tlog.writeBatchToBuffer(util.TimeNano(), newTxId(), bp)
	bl := &batchLog{bp, newTxId(), util.TimeNano()}
	buf = bl.toBuffer()
	defer buf.Free()
	return
}

func reSetStoreNodeNumber(num int32) (err error) {
	defer errRecover()
	if err = BroadcastReinit(&SysBean{STORENODENUM: num}).Error(); err == nil {
		sys.STORENODENUM = int(num)
	}
	return
}

func broadRmNode() (err error) {
	defer errRecover()
	err = BroadcastReinit(&SysBean{0, 0, sys.UUID}).Error()
	return
}

func GetRemoteRunUUID() []int64 {
	return nodeWare.GetRemoteRunUUID()
}

func IsClusRun() bool {
	return nodeWare.IsClusRun()
}

func LB() (_r int64) {
	if sys.CcPut() < sys.COCURRENT_PUT/10 && sys.CcGet() < sys.COCURRENT_GET/10 {
		return
	}
	tcs := nodeWare.GetAllTlContext()
	for _, tc := range tcs {
		if tc.stat == sys.RUN && tc.CcPut < sys.CcPut()/2 && tc.CcGet < sys.CcGet()/2 {
			_r = tc.RemoteUuid
		}
	}
	return
}

/************************************************************/

func client2Serve(addr string) (err error) {
	if sys.IsStandAlone() {
		return errors.New("node is stand-alone")
	}
	if nodeWare.LenByAddr(addr) > 0 {
		return errors.New("addr:" + addr + " already existed")
	}
	if addr == sys.ADDR {
		return errors.New("addr:" + addr + " is local addr")
	}
	if addr = strings.Trim(addr, " "); addr == "" {
		return errors.New("addr is bad")
	}
	delNodeMap.Del(addr)
	err1, err2 := tnetservice.Connect(addr, true)
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return
}

func AddInject(id int64, f func()) {
	statAdmin.addInject(id, f)
}

/************************************************************/
func loadData2TLDB(bs []byte, datetime string) (err error) {
	_, err = tlog.parseGzfile(0, "", bs, datetime)
	return
}

func forcedCoverageData2TLDB(body []byte, datetime string, limit int64) (err error) {
	_, err = tlog.ForcedCoverageParseGzfile(body, datetime, limit)
	return
}

/************************************************************/

func getRemoteNode() []*RemoteNode {
	return nodeWare.getRemoteNodes()
}
