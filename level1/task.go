// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package level1

import (
	"sync"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/sys"
	. "github.com/donnie4w/tldb/util"
)

var taskWare = NewTaskWare()

func NewTaskWare() (_r *_taskWare) {
	_r = &_taskWare{NewNumLock(1 << 13), NewMap[int64, task](), NewLinkedMap[int64, int64]()}
	go _r.timerForExpired()
	return
}

type _taskWare struct {
	mux              *Numlock
	pool             *Map[int64, task]
	expiredPoolCache *LinkedMap[int64, int64]
}

func (this *_taskWare) pubTaskPB(pb *PonBean) (_r task, err error) {
	this.mux.Lock(pb.Txid)
	defer this.mux.Unlock(pb.Txid)
	return this._newTask(pb)
}

func (this *_taskWare) deldone(txid int64) {
	this.pool.Del(txid)
	trashTx.Put(txid, 0)
	GoPool.Go(func() {
		this.expiredPoolCache.Del(txid)
	})
}

func (this *_taskWare) delExpired(txid int64) {
	this.pool.Del(txid)
	GoPool.Go(func() {
		this.expiredPoolCache.Del(txid)
	})
}

func (this *_taskWare) has(txid int64) bool {
	return this.pool.Has(txid)
}

func (this *_taskWare) _newTask(pb *PonBean) (_r task, err error) {
	var ok bool
	if _r, ok = this.pool.Get(pb.Txid); !ok {
		switch pb.Ptype {
		case 1:
			_r = NewIncrTask(pb.Txid)
		case 2:
			_r = NewBatchTask(pb.Txid)
		case 22:
			_r = NewBatchProcessTask(pb.Txid)
		case 3:
			_r = NewGetTask(pb.Txid)
		case 31:
			_r = NewGetRemoteTask(pb.Txid)
		case 4:
			_r = NewSyncKeyTask(pb.Txid)
		case 5:
			_r = NewStatTask(pb.Txid)
		case 6:
			_r = NewLoadTask(pb.Txid)
		case 7:
			_r = NewTimeTask(pb.Txid)
		case 8:
			_r = NewSeqTask(pb.Txid)
		}
		this.pool.Put(pb.Txid, _r)
		this.expiredPoolCache.Put(pb.Txid, Time().UnixNano())
	}
	err = _r.addPB(pb)
	return
}

func (this *_taskWare) timerForExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			func() {
				defer myRecovr()
				this.expiredPoolCache.BackForEach(func(k, v int64) bool {
					if v+3600 < Time().Unix() {
						this.delExpired(k)
						return true
					}
					return false
				})
			}()
		}
	}
}

type task interface {
	addPB(pb *PonBean) (err error)
	event() (icomes, noMatch bool)
	task() (err error)
	ch() chan int8
	getPB(uuid int64) *PonBean
	getFirstPB() *PonBean
	getAllPB() *Map[int64, *PonBean]
	onError() error
	isDone() bool
	waitWithTimeout(t time.Duration) (err error)
}

/******************************************************/

func NewBatchTask(txid int64) task {
	return &batchTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type batchTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	verifyPB *PonBean
	_isDone  bool
	C        chan int8
	_onError error
	mode     int8
}

func (this *batchTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		if pb.IsVerify {
			this.verifyPB = pb
		} else {
			if pb.IsFirst {
				this.firstPB = pb
			}
			this.m.Put(pb.Fromuuid, pb)
		}
		if this.mode < pb.Mode {
			this.mode = pb.Mode
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
			err = this._onError
		} else if noMatch {
			this._isDone = true
			defer func() {
				this._onError = Errors(sys.ERR_CLUS_NOMATCH)
				this.C <- 2
			}()
		}
	} else {
		logError.Error("batchTask is finish:txid:", pb.Txid, ",pb.Fromuuid:", pb.Fromuuid)
	}
	return
}

func (this *batchTask) event() (icomes, noMatch bool) {
	if (this.firstPB != nil && this.verifyPB != nil) || sys.IsStandAlone() {
		if icomes = txCallbackFunc(this.m, nodeWare.GetRemoteRunUUID()); icomes {
			this.m.Range(func(_ int64, v *PonBean) bool {
				if !ArrayEqual(v.ClusterNode, nodeWare.GetAllRunUUID()) || !ArrayEqual(v.ExecNode, this.firstPB.ExecNode) {
					noMatch = true
					return false
				}
				return true
			})
		}
	}
	return
}

func (this *batchTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	if !sys.IsStandAlone() {
		if !this.verifyPB.DoCommit {
			return
		}
	}
	this._isDone = true
	pb := this.firstPB
	if err := tlog.writePonBean(TimeNano(), pb); err != nil {
		fatalError(err)
	} else if pb.Batch.Dels != nil {
		incrExcuPool.delMulti(pb.Batch.Dels)
	}
	return
}

func (this *batchTask) ch() chan int8 {
	return this.C
}

func (this *batchTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *batchTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}

func (this *batchTask) onError() error {
	return this._onError
}

func (this *batchTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}

func (this *batchTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *batchTask) waitWithTimeout(t time.Duration) (err error) {
	return
}

/******************************************************************************/
func NewBatchProcessTask(txid int64) task {
	return &batchProcessTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type batchProcessTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	verifyPB *PonBean
	_isDone  bool
	C        chan int8
	_onError error
	mode     int8
}

func (this *batchProcessTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		if pb.IsVerify {
			this.verifyPB = pb
		} else {
			if pb.IsFirst {
				this.firstPB = pb
			}
			this.m.Put(pb.Fromuuid, pb)
		}
		if this.mode < pb.Mode {
			this.mode = pb.Mode
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
			err = this._onError
		} else if noMatch {
			this._isDone = true
			defer func() {
				this._onError = Errors(sys.ERR_CLUS_NOMATCH)
				this.C <- 2
			}()
		}
	} else {
		logError.Error("batchTask is finish:txid:", pb.Txid, ",pb.Fromuuid:", pb.Fromuuid)
	}
	return
}

func (this *batchProcessTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil && this.verifyPB != nil {
		icomes = true
	}
	return
}

func (this *batchProcessTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	if !this.verifyPB.DoCommit {
		return
	}
	this._isDone = true
	err = doCommit(this.firstPB)
	return
}

func doCommit(pb *PonBean) (err error) {
	if sys.IsRUN() {
		if err := tlog.writePonBean(TimeNano(), pb); err != nil {
			fatalError(err)
		} else if pb.Batch.Dels != nil {
			incrExcuPool.delMulti(pb.Batch.Dels)
		}
	} else if sys.IsREADY() {
		tlog.WriteCacheLog(pb)
	}
	return
}

func (this *batchProcessTask) ch() chan int8 {
	return this.C
}

func (this *batchProcessTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *batchProcessTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}

func (this *batchProcessTask) onError() error {
	return this._onError
}

func (this *batchProcessTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}

func (this *batchProcessTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *batchProcessTask) waitWithTimeout(t time.Duration) (err error) {
	return
}

/******************************************************/
func NewGetTask(txid int64) task {
	return &getTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type getTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *getTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		this.m.Put(pb.Fromuuid, pb)
		if this.firstPB == nil && pb.IsFirst {
			this.firstPB = pb
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("getTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *getTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil {
		icomes = this.m.Has(this.firstPB.ExecNode[0])
	}
	return
}

func (this *getTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *getTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *getTask) ch() chan int8 {
	return this.C
}

func (this *getTask) onError() error {
	return this._onError
}
func (this *getTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *getTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *getTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}
func (this *getTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewGetRemoteTask(txid int64) task {
	return &getRemoteTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type getRemoteTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	validPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *getRemoteTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		if pb.IsFirst {
			this.firstPB = pb
		}
		this.m.Put(pb.Fromuuid, pb)
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("getAllTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *getRemoteTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil {
		this.m.Range(func(_ int64, pb *PonBean) bool {
			if pb.Value != nil {
				this.firstPB.Value = pb.Value
				return false
			}
			return true
		})
		icomes = this.firstPB.Value != nil
	}
	return
}

func (this *getRemoteTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *getRemoteTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *getRemoteTask) ch() chan int8 {
	return this.C
}

func (this *getRemoteTask) onError() error {
	return this._onError
}
func (this *getRemoteTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *getRemoteTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *getRemoteTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}
func (this *getRemoteTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewIncrTask(txid int64) task {
	return &incrTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type incrTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *incrTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		this.m.Put(pb.Fromuuid, pb)
		if this.firstPB == nil && pb.IsFirst {
			this.firstPB = pb
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("incrTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *incrTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil && this.m.Has(this.firstPB.ExecNode[0]) {
		icomes = true
	}
	return
}

func (this *incrTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *incrTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}
func (this *incrTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *incrTask) ch() chan int8 {
	return this.C
}

func (this *incrTask) onError() error {
	return this._onError
}
func (this *incrTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}

func (this *incrTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *incrTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewSyncKeyTask(txid int64) task {
	return &syncKeyTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type syncKeyTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *syncKeyTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		this.m.Put(pb.Fromuuid, pb)
		if pb.IsFirst {
			this.firstPB = pb
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("syncKeyTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *syncKeyTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil {
		icomes = txCallbackFunc(this.m, nodeWare.GetRemoteRunUUID())
	}
	return
}

func (this *syncKeyTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *syncKeyTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *syncKeyTask) ch() chan int8 {
	return this.C
}

func (this *syncKeyTask) onError() error {
	return this._onError
}

func (this *syncKeyTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *syncKeyTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *syncKeyTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *syncKeyTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewStatTask(txid int64) task {
	return &statTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type statTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
	verifyPB *PonBean
}

func (this *statTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		if pb.IsVerify {
			this.verifyPB = pb
		} else {
			if pb.IsFirst {
				this.firstPB = pb
			}
			this.m.Put(pb.Fromuuid, pb)
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("statTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *statTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil && this.verifyPB != nil {
		icomes = true
	}
	return
}

func (this *statTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	if this.firstPB.Fromuuid == sys.UUID {
		setStat(sys.STATTYPE(this.firstPB.Stat.Stat), time.Duration(this.firstPB.Stat.Timenano-Time().UnixNano()))
	} else {
		nodeWare.setStat(this.firstPB.Fromuuid, sys.STATTYPE(this.firstPB.Stat.Stat), this.firstPB.Stat.Timenano-Time().UnixNano())
		if sys.STATTYPE(this.firstPB.Stat.Stat) != sys.RUN {
			checkAndSetRunToReady()
		}
	}
	if !sys.CLUSTER_NUM_FINAL {
		if sys.STATTYPE(this.firstPB.Stat.Stat) == sys.RUN {
			if n := (len(nodeWare.GetAllRunUUID())+1)/2 + 1; n > sys.CLUSTER_NUM && sys.CLUSTER_NUM > 0 {
				sys.CLUSTER_NUM = n
			}
		}
		if sys.STATTYPE(this.firstPB.Stat.Stat) == sys.PROXY && nodeWare.IsClusRun() {
			sys.CLUSTER_NUM = (len(nodeWare.GetAllRunUUID())+1)/2 + 1
		}
	}
	return
}

func (this *statTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *statTask) ch() chan int8 {
	return this.C
}

func (this *statTask) onError() error {
	return this._onError
}

func (this *statTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *statTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *statTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *statTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewLoadTask(txid int64) task {
	return &loadTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type loadTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *loadTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		this.m.Put(pb.Fromuuid, pb)
		if pb.IsFirst {
			this.firstPB = pb
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("loadTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *loadTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil && this.m.Has(this.firstPB.ExecNode[0]) {
		icomes = true
	}
	return
}

func (this *loadTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *loadTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *loadTask) ch() chan int8 {
	return this.C
}

func (this *loadTask) onError() error {
	return this._onError
}

func (this *loadTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *loadTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *loadTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *loadTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewTimeTask(txid int64) task {
	return &timeTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type timeTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *timeTask) addPB(pb *PonBean) (err error) {
	go this._addPB(pb)
	return
}
func (this *timeTask) _addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		this.m.Put(pb.Fromuuid, pb)
		if pb.IsFirst {
			this.firstPB = pb
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("timeTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *timeTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil {
		icomes = txCallbackFunc(this.m, nodeWare.GetRemoteUUIDS())
	}
	return
}

func (this *timeTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *timeTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *timeTask) ch() chan int8 {
	return this.C
}

func (this *timeTask) onError() error {
	return this._onError
}

func (this *timeTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *timeTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *timeTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *timeTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}

/******************************************************/
func NewSeqTask(txid int64) task {
	return &seqTask{txid: txid, m: NewMap[int64, *PonBean](), mux: &sync.Mutex{}, C: make(chan int8, 1)}
}

type seqTask struct {
	txid     int64
	m        *Map[int64, *PonBean]
	mux      *sync.Mutex
	firstPB  *PonBean
	_isDone  bool
	C        chan int8
	_onError error
}

func (this *seqTask) addPB(pb *PonBean) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this._isDone {
		this.m.Put(pb.Fromuuid, pb)
		if pb.IsFirst {
			this.firstPB = pb
		}
		if icomes, noMatch := this.event(); icomes && !noMatch {
			this._isDone = true
			defer func() { this.C <- 1 }()
			this._onError = this.task()
		} else if noMatch {
			this._isDone = true
			defer func() { this.C <- 2 }()
		}
	} else {
		logError.Error("seqTask is finish:txid:", pb.Txid, ",pb.Ptype:", pb.Ptype)
	}
	return
}

func (this *seqTask) event() (icomes, noMatch bool) {
	if this.firstPB != nil {
		icomes = txCallbackFunc(this.m, nodeWare.GetRemoteUUIDS())
	}
	return
}

func (this *seqTask) task() (err error) {
	defer myRecovr()
	defer taskWare.deldone(this.txid)
	return
}

func (this *seqTask) getPB(uuid int64) (_r *PonBean) {
	_r, _ = this.m.Get(uuid)
	return
}

func (this *seqTask) ch() chan int8 {
	return this.C
}

func (this *seqTask) onError() error {
	return this._onError
}

func (this *seqTask) getAllPB() *Map[int64, *PonBean] {
	return this.m
}
func (this *seqTask) getFirstPB() (_r *PonBean) {
	_r = this.firstPB
	return
}
func (this *seqTask) isDone() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this._isDone
}

func (this *seqTask) waitWithTimeout(t time.Duration) (err error) {
	if !this.isDone() {
		select {
		case <-this.ch():
		case <-time.After(t):
			err = Errors(sys.ERR_TIMEOUT)
		}
	}
	return
}
