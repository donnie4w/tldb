// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"fmt"
	"strings"
	"sync"
	"time"

	. "github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/level0"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

var synctx = newSyncTx()

type syncTx struct {
	mux *sync.Mutex
}

func newSyncTx() (_r *syncTx) {
	_r = &syncTx{&sync.Mutex{}}
	return
}

func (this *syncTx) init() {
	go this.backInit()
	go this.flow(8 * time.Second)
	go this.CleanIdxsCache(8 * time.Second)
}

/************************************************/
func (this *syncTx) flow(t time.Duration) {
	ticker := time.NewTicker(t)
	for {
		select {
		case <-ticker.C:
			func() {
				defer myRecovr()
				if !sys.IsPROXY() && isSyncOver() {
					tlog.LoadCacheLog()
				}
				if sys.IsREADY() && isSyncOver() && len(nodeWare.GetALLUUID()) >= sys.CLUSTER_NUM {
					if sys.TIME_DEVIATION == 0 {
						if pos_time() != nil {
							return
						}
					}
					if ready2Run() == nil {
						tlog.LoadCacheLog()
					}
				} else if !isSyncOver() && sys.IsREADY() {
					statAdmin.pullData()
					statAdmin.amendPull()
				}
			}()
		}
	}
}

func (this *syncTx) backInit() {
	if keys, err := level0.Level0.GetKeysPrefix(KeyLevel1.BackPrefix()); err == nil && keys != nil {
		for _, key := range keys {
			if pb, uuids, err := tlog.GetPonBeanByBack(key); err == nil {
				for _, uuid := range uuids {
					go pubSingleTxAlltime(pb, uuid, true)
				}
			} else {
				level0.Level0.Del(key)
			}
		}
	}
}

func ready2Run() (err error) {
	if sys.IsREADY() && isSyncOver() && len(nodeWare.GetALLUUID()) >= sys.CLUSTER_NUM {
		i := 10
		for i > 0 {
			i--
			if err = pos_stat(sys.RUN, 0); err == nil {
				break
			}
		}
		tlog.syncCount = 0
		statAdmin.syncRunInject()
	}
	return
}
func isSyncOver() (b bool) {
	if sys.IsRUN() {
		return true
	}
	b = statAdmin.isSyncOver()
	return
}

/*********************************************************************/

func (this *syncTx) CleanIdxsCache(t time.Duration) (err error) {
	ticker := time.NewTicker(t)
	for {
		select {
		case <-ticker.C:
			var limit int64 = 1000
			if sys.MAXDELSEQ-sys.MAXDELSEQCURSOR >= int64(limit) {
				c := int64(0)
				delkeys := make([]string, 0)
				for i := sys.MAXDELSEQCURSOR; i <= sys.MAXDELSEQ; i++ {
					if bs, err := level0.Level0.Get(KeyLevel3.SeqForDel(fmt.Sprint(i))); err == nil {
						_idx_key := string(bs)
						li := strings.LastIndex(_idx_key, "_")
						maxSeqKey := _idx_key[:li]
						idx := fmt.Sprint(KEY2_SEQ, maxSeqKey[len(KEY2_IDX_):])
						if ss, err := level0.Level0.GetKeysPrefixLimit(maxSeqKey+"_", 1); err != nil || len(ss) == 0 {
							delkeys = append(delkeys, idx)
						}
						delkeys = append(delkeys, KeyLevel3.SeqForDel(fmt.Sprint(i)))
					}
					c++
					if c >= limit {
						break
					}
				}
				if err := Level1.Batch(3, nil, delkeys); err == nil {
					sys.MAXDELSEQCURSOR += c
					level0.Level0.Put(KeyLevel3.KeyMaxDelSeqCursor(), util.Int64ToBytes(sys.MAXDELSEQCURSOR))
				}
			}
		}
	}
}
