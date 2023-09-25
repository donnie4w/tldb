// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"bytes"
	"context"
	"sort"
	"sync/atomic"
	"time"

	"github.com/donnie4w/tldb/keystore"
	. "github.com/donnie4w/tldb/lock"
	. "github.com/donnie4w/tldb/stub"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

type ItnetServ struct {
	mux *Numlock
}

func ctx2TlContext(ctx context.Context) *tlContext {
	return ctx.Value(tlContextCtx).(*tlContext)
}

// //////////////////////////////////////
func (this *ItnetServ) Ping(ctx context.Context, pingstr int64) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.pongNum > 10 && !tc.isAuth {
			logError.Error("ping error:", tc.RemoteUuid)
			tc.Close()
		} else {
			buf := pongBytes()
			util.GoPool.Go(func() { tc.Conn.Pong(context.Background(), buf.Bytes()) })
		}
		if pingstr != pistr() {
			util.GoPool.Go(func() { tc.Conn.SyncNode(context.Background(), nodeWare.GetUUIDNode(), true) })
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) Pong(ctx context.Context, pongBs []byte) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		defer this.mux.Unlock(tc.Id)
		this.mux.Lock(tc.Id)
		atomic.AddInt64(&tc.pingNum, -1)
		atomic.AddInt64(&tc.pongNum, 1)
		pong := util.BytesToInt64(pongBs[:8])
		tc.CcPut = util.BytesToInt64(pongBs[8:16])
		tc.CcGet = util.BytesToInt64(pongBs[16:])
		if tc.pongFlag != pong || pongs(tc.RemoteUuid, tc.stat) != pong {
			if tc.pongFlag > 0 {
				go tc.Conn.SyncNode(context.Background(), nodeWare.GetUUIDNode(), true)
			}
			tc.pongFlag = pong
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) Auth(ctx context.Context, authKey []byte) (_err error) {
	defer myRecovr()
	tc := ctx2TlContext(ctx)
	if authTc(tc, authKey) {
		if at, err := getAuth(); err == nil {
			_err = tc.Conn.Auth2(context.Background(), at)
		}
	}
	return
}

func (this *ItnetServ) Auth2(ctx context.Context, authKey []byte) (_err error) {
	defer myRecovr()
	tc := ctx2TlContext(ctx)
	if authTc(tc, authKey) {
		tc.Conn.SyncNode(context.Background(), nodeWare.GetUUIDNode(), true)
	}
	return
}

func (this *ItnetServ) PonMerge(ctx context.Context, pblist []*PonBean, id int64) (_err error) {
	go func() {
		defer myRecovr()
	}()
	return
}

// Parameters:
//   - Pb
func (this *ItnetServ) Pon(ctx context.Context, ponBeanBytes []byte, id int64) (err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			pb := DecodePonBean(ponBeanBytes)
			reTxId.Add(util.Int64ToBytes(pb.Txid))
			if !rePon.Has(id) {
				rePon.Put(id, 0)
				if !trashTx.Has(pb.Txid) {
					if taskWare.has(pb.Txid) {
						_, err = taskWare.pubTaskPB(pb)
					} else {
						switch pb.Ptype {
						case 1:
							util.GoPool.Go(func() { pas_ponIncrkey(pb) })
						case 2:
							err = pas_ponBatch(pb)
						case 22:
							err = pas_ponBatch(pb)
						case 3:
							pas_ponGet(pb)
						case 31:
							pas_ponGetRemote(pb)
						case 4:
							pas_ponSyncKey(pb)
						}
					}
				} else {
					err = util.Errors(sys.ERR_TXOVER)
					logError.Error("TrashTx txid:", pb.Txid, " ,Ptype:", pb.Ptype)
				}
			}
			if err == nil {
				syncTxAck(tc, tc.RemoteUuid, id, 0)
			} else {
				syncTxAck(tc, tc.RemoteUuid, id, 1)
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) Pon2(ctx context.Context, pb *PonBean, id int64) (_err error) {
	go func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			syncTxAck(tc, tc.RemoteUuid, id, 0)
			reTxId.Add(util.Int64ToBytes(pb.Txid))
			if !rePon.Has(id) {
				rePon.Put(id, 0)
				if !trashTx.Has(pb.Txid) {
					if !trashTx.Has(pb.Txid) {
						if taskWare.has(pb.Txid) {
							taskWare.pubTaskPB(pb)
						} else {
							switch pb.Ptype {
							case 5:
								pas_ponStat(pb)
							case 6:
								pas_ponLoad(pb)
							case 7:
								pas_ponTime(pb)
							case 8:
								pas_ponSeq(pb)
							}
						}
					} else {
						logger.Warn("TrashTx txid:", pb.Txid, " ,Ptype:", pb.Ptype)
					}
				}
			}
		}
	}()
	return
}

func (this *ItnetServ) Pon3(ctx context.Context, ponBeanBytes []byte, id int64, ack bool) (_err error) {
	f := func() {
		defer myRecovr()
		pb := DecodePonBean(ponBeanBytes)
		if ack {
			awaitPB.DelAndPut(id, pb)
			tc := ctx2TlContext(ctx)
			syncTxAck(tc, tc.RemoteUuid, id, 0)
		} else {
			ackPb := newPonBean(pb.Txid, pb.Ptype, pb.OpType, nil, nil)
			key := string(pb.Key)
			switch pb.Ptype {
			case 1:
				value, _ := incrNetKeyExcu(key, pb)
				ackPb.Key, ackPb.Value = pb.GetKey(), util.Int64ToBytes(value)
				responeSingleTx(ackPb, pb.Fromuuid, *pb.SyncId)
			case 3:
				if value, err := Level1.GetLocal(key); err == nil {
					if data := decodeDataBeen(value); data != nil && data.IsLocal {
						ackPb.Value = data.Value
					}
				}
				responeSingleTx(ackPb, pb.Fromuuid, *pb.SyncId)
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
	go func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if tc.NameSpace != node.Ns {
				logError.Error("namespace mismatch:", node.Ns)
				nodeWare.del(tc)
				return
			}
			if tc.RemoteUuid == 0 || node.UUID == 0 || node.Addr == "" || tc.RemoteUuid != node.UUID {
				logError.Error("Information matching error:", node.UUID)
				return
			}
			if nodeWare.has(tc) {
				tc.stat, tc.RemoteAdminAddr, tc.RemoteMqAddr = sys.STATTYPE(node.Stat), node.AdminAddr, node.MqAddr
				nodeWare.setStat(tc.RemoteUuid, sys.STATTYPE(node.Stat), 0)
			} else {
				tc.RemoteUuid, tc.RemoteAddr, tc.stat, tc.RemoteAdminAddr, tc.RemoteMqAddr, tc.RemoteCliAddr = node.UUID, node.Addr, sys.STATTYPE(node.Stat), node.AdminAddr, node.MqAddr, node.CliAddr
			}
			err := nodeWare.add(tc)
			if ir {
				tc.Conn.SyncNode(context.Background(), nodeWare.GetUUIDNode(), !ir)
			}
			for k, v := range node.Nodekv {
				if k != sys.UUID && !nodeWare.hasUUID(k) && !clientLinkCache.Has(v) {
					if !nodeWare.hasUUID(k) {
						<-time.After(time.Duration(util.Rand(6)) * time.Second)
						tnetservice.Connect(v, true)
					}
				}
			}
			if err != nil {
				return
			}
			if sys.IsREADY() && tc.stat == sys.RUN && sys.TIME_DEVIATION == 0 {
				syncTime(tc)
			}
		}
	}()
	return
}

// Parameters:
//   - Timenano
//   - Num
//   - Ir
func (this *ItnetServ) Time(ctx context.Context, pretimenano, timenano int64, num int16, ir bool) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			defer this.mux.Unlock(tc.Id)
			this.mux.Lock(tc.Id)
			if (ir && sys.IsREADY()) || (!sys.IsREADY() && len(sys.TIME_DEVIATION_LIST) > 0) {
				deviation := timenano + (time.Now().UnixNano()-pretimenano)/2 - time.Now().UnixNano()
				sys.TIME_DEVIATION_LIST = append(sys.TIME_DEVIATION_LIST, deviation)
				pretimenano = time.Now().UnixNano()
				if num <= 1 {
					var i int64
					for _, v := range sys.TIME_DEVIATION_LIST {
						i += v
					}
					sys.TIME_DEVIATION = i / int64(len(sys.TIME_DEVIATION_LIST))
					sys.TIME_DEVIATION_LIST = make([]int64, 0)
				}
			}
			if num > 0 {
				tc.Conn.Time(context.Background(), pretimenano, util.Time().UnixNano(), num-1, !ir)
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) SyncTx(ctx context.Context, syncId int64, result int8) (_err error) {
	f := func() {
		defer myRecovr()
		if ctx2TlContext(ctx).isAuth {
			this._SyncTx(syncId, result)
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
	f := func() {
		defer myRecovr()
		if ctx2TlContext(ctx).isAuth {
			if syncList != nil {
				for k, v := range syncList {
					this._SyncTx(k, v)
				}
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) _SyncTx(syncId int64, result int8) {
	defer myRecovr()
	if result == 0 {
		await.DelAndClose(syncId)
	} else {
		await.DelAndPut(syncId, result)
	}
}

// Parameters:
//   - SyncId
//   - Txid
func (this *ItnetServ) CommitTx(ctx context.Context, syncId int64, txid int64) (_err error) {
	// f := func() {
	// 	defer myRecovr()
	// 	tc := ctx2TlContext(ctx)
	// 	if tc.isAuth {
	// 		defer this.mux.Unlock(txid)
	// 		this.mux.Lock(txid)
	// 		uuid := tc.RemoteUuid
	// 		if tlog.HasBackLog(txid) {
	// 			if err := nodeWare.commitTx2(uuid, syncId, txid, true, sys.ReadTimeout); err == nil {
	// 				tlog.RemoveBackLogUUid(txid, uuid)
	// 			}
	// 		} else {
	// 			nodeWare.commitTx2(uuid, syncId, txid, false, sys.ReadTimeout)
	// 		}
	// 	}
	// }
	// util.GoPool.Go(f)
	return
}

// Parameters:
//   - SyncId
//   - Commit
func (this *ItnetServ) CommitTx2(ctx context.Context, syncId int64, txid int64, commit bool) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if syncId > 0 {
				syncTxAck(tc, tc.RemoteUuid, syncId, 1)
			}
			if commit {
				await.DelAndPut(syncId, 2)
			} else {
				await.DelAndPut(syncId, 1)
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) PubMq(ctx context.Context, syncId int64, mqType int8, bs []byte) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if syncId > 0 {
				go syncTxAck(tc, tc.RemoteUuid, syncId, 1)
			}
			if !reMq.Has(syncId) {
				reMq.Put(syncId, 0)
				ClusPub(mqType, util.SnappyDecode(bs))
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) PullData(ctx context.Context, syncId int64, ldb *LogDataBean) (_err error) {
	go func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if ldb.Type == 2 && syncId > 0 {
				awaitLog.DelAndPut(syncId, ldb)
			} else if ldb.Type == 1 {
				tlog.sendBinFile(syncId, ldb.LastTime, ldb.Txid, tc.RemoteUuid)
			}
		}
	}()
	return
}

func (this *ItnetServ) ReInit(ctx context.Context, syncId int64, sbean *SysBean) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if sbean.DBMode > 0 {
				sys.DBMode = int(sbean.DBMode)
			}
			if sbean.STORENODENUM > 0 {
				sys.STORENODENUM = int(sbean.STORENODENUM)
			}
			if sbean.REMOVEUUID > 0 {
				removeNode(sbean.REMOVEUUID)
			}
			statAdmin.syncRunInject()
			if syncId > 0 {
				syncTxAck(tc, tc.RemoteUuid, syncId, 0)
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) ProxyCall(ctx context.Context, syncId int64, paramData []byte, pType int8, ctype int8) (_err error) {
	f := func() {
		defer myRecovr()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if ctype == 1 {
				if tp, err := util.TDecode(paramData, &TableParam{}); err == nil {
					switch pType {
					case 1:
						seq, err := Insert(1, tp.Stub)
						tp = &TableParam{Seq: seq}
						if err != nil {
							tp.Err = err.Error()
						}
					case 2:
						seq, err := InsertMq(1, tp.Stub)
						tp = &TableParam{Seq: seq}
						if err != nil {
							tp.Err = err.Error()
						}
					case 3:
						_r, err := SelectId(1, tp.Stub)
						tp = &TableParam{Seq: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 4:
						_r, err := SelectIdByIdx(1, tp.TableName, tp.IdxName, tp.IdxValue)
						tp = &TableParam{Seq: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 5:
						_r, err := SelectById(1, tp.Stub)
						tp = &TableParam{Stub: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 6:
						_r, err := SelectsByIdLimit(1, tp.Stub, tp.Start, tp.Limit)
						tp = &TableParam{StubArray: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 7:
						_r, err := SelectByIdx(1, tp.TableName, tp.IdxName, tp.IdxValue)
						tp = &TableParam{Stub: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 8:
						_r, err := SelectsByIdx(1, tp.TableName, tp.IdxName, tp.IdxValue)
						tp = &TableParam{StubArray: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 9:
						_r, err := SelectsByIdxLimit(1, tp.TableName, tp.IdxName, tp.IdxValues, tp.Start, tp.Limit)
						tp = &TableParam{StubArray: _r}
						if err != nil {
							tp.Err = err.Error()
						}
					case 10:
						err := Update(1, tp.Stub)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					case 11:
						err := Delete(1, tp.Stub)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					case 12:
						err := CreateTable(tp.Stub)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					case 13:
						err := AlterTable(tp.Stub)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					case 14:
						err := CreateTableMq(tp.Stub)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					case 15:
						err := DropTable(tp.Stub)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					case 16:
						tss := LoadTableInfo()
						tp = &TableParam{StubArray: tss}
					case 17:
						tss := LoadMQTableInfo()
						tp = &TableParam{StubArray: tss}
					case 18:
						err := DeleteBatches(1, tp.TableName, tp.Start, tp.Limit+tp.Start)
						tp = &TableParam{}
						if err != nil {
							tp.Err = err.Error()
						}
					}
					bs := util.TEncode(tp)
					ackProxy(tc, tc.RemoteUuid, syncId, util.SnappyEncode(bs), pType)
				} else {
					logger.Error(err)
				}
			} else if ctype == 2 {
				bs := util.SnappyDecode(paramData)
				if tp, err := util.TDecode(bs, &TableParam{}); err == nil {
					awaitProxy.DelAndPut(syncId, tp)
				}
			}
		}
	}
	util.GoPool.Go(f)
	return
}

func (this *ItnetServ) BroadToken(ctx context.Context, syncId int64, tt *TokenTrans, ack bool) (_err error) {
	defer myRecovr()
	tc := ctx2TlContext(ctx)
	if tc.isAuth {
		if ack {
			awaitToken.DelAndPut(syncId, tt)
			syncTxAck(tc, tc.RemoteUuid, syncId, 0)
		} else {
			go tokenroute.processTokenTrans(tt, tc.RemoteUuid)
		}
	}
	return
}

// //////////////////////////////////////////////////////////////////

func syncTxAck(tc *tlContext, uuid, syncId int64, result int8) (err error) {
	defer myRecovr()
	if sys.MEGERCLUSACK {
		// var uuidMg = UUID(uuid)
		// uuidMg.syncTxAdd(&syncBean{syncId, result}, func() bool { return false }, uuidMg.syncTxMerge)
		tc.mergeChan <- &syncBean{syncId, result}
		atomic.AddInt64(&tc.mergeCount, 1)
		if tc.mergemux.TryLock() {
			go syncMerge(tc)
		}
	} else {
		if err = tc.Conn.SyncTx(context.Background(), syncId, result); err != nil {
			if _tc := nodeWare.GetTlContext(uuid); _tc != nil {
				err = _tc.Conn.SyncTx(context.Background(), syncId, result)
			}
		}
	}
	return
}

func ackProxy(tc *tlContext, uuid, syncId int64, paramData []byte, pType int8) (err error) {
	defer myRecovr()
	if err = tc.Conn.ProxyCall(context.Background(), syncId, paramData, pType, 2); err != nil {
		if _tc := nodeWare.GetTlContext(uuid); _tc != nil {
			err = _tc.Conn.ProxyCall(context.Background(), syncId, paramData, pType, 2)
		}
	}
	return
}

func authTc(tc *tlContext, authKey []byte) (b bool) {
	if authKey != nil && len(authKey) > 0 {
		if bs, err := keystore.RsaDecrypt(authKey, sys.PRIVATEKEY); err == nil {
			if r := util.BytesToInt64(bs[:8]); r > 0 && r != sys.UUID && !nodeWare.hasUUID(r) && string(bs[8:]) == sys.PWD {
				tc.RemoteUuid = r
				tc.isAuth = true
				b = true
			}
		}
	}
	return
}

func pistr() (_r int64) {
	all := nodeWare.GetALLUUID()
	if len(all) > 0 {
		sort.Slice(all, func(i, j int) bool {
			return all[i] < all[j]
		})
		var buf bytes.Buffer
		for _, a := range all {
			buf.Write(util.Int64ToBytes(a))
		}
		_r = int64(util.CRC64(buf.Bytes()))
	}
	return
}

func pongBytes() (buf *bytes.Buffer) {
	buf = util.BufferPool.Get(24)
	buf.Write(util.Int64ToBytes(pongs(sys.UUID, sys.SYS_STAT)))
	buf.Write(util.Int64ToBytes(sys.CcPut()))
	buf.Write(util.Int64ToBytes(sys.CcGet()))
	return
}

func pongs(uuid int64, stat sys.STATTYPE) (_r int64) {
	buf := util.BufferPool.Get(12)
	buf.Write(util.Int64ToBytes(uuid))
	buf.Write(util.Int32ToBytes(int32(stat)))
	_r = int64(util.CRC64(buf.Bytes()))
	util.BufferPool.Put(buf)
	return
}

func getAuth() (bs []byte, err error) {
	buf := util.BufferPool.Get(8 + len([]byte(sys.PWD)))
	buf.Write(util.Int64ToBytes(sys.UUID))
	buf.WriteString(sys.PWD)
	bs, err = keystore.RsaEncrypt(buf.Bytes(), sys.PUBLICKEY)
	util.BufferPool.Put(buf)
	return
}
