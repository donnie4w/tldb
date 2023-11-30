// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package log

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	. "github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/tldb/key"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

var BinLog *_binLog

var SUFGZIP = ".gz"

type lastStatBean struct {
	lastTime int64
	lastTxid int64
}

type fristStatBean struct {
	fristTime int64
	fristTxid int64
}

type _binLog struct {
	mux4write *sync.RWMutex
	mux4gzip  *sync.Mutex
	prename   string
	sufname   string
	statBean  *lastStatBean
	binDir    string
}

func NewBinLog() (blog *_binLog) {
	ss := strings.Split(sys.BINLOGNAME, ".")
	blog = &_binLog{&sync.RWMutex{}, &sync.Mutex{}, ss[0], ss[1], &lastStatBean{}, sys.DBFILEDIR + "/bin"}
	blog.statBean.lastTime, blog.statBean.lastTxid = blog.getLastTime()
	return
}

func (this *_binLog) gzipLog() {
	this.mux4gzip.Lock()
	defer this.mux4gzip.Unlock()
	if fs, err := os.Open(this.binDir); err == nil {
		defer fs.Close()
		if dirs, err := fs.ReadDir(-1); err == nil {
			for _, f := range dirs {
				fname := f.Name()
				if !f.IsDir() && strings.HasSuffix(fname, "."+this.sufname) && strings.HasPrefix(fname, this.prename+"_") {
					num := fname[len(this.prename)+1 : len(fname)-len(this.sufname)-1]
					if util.MatchString("^[0-9]{1,}$", num) {
						gzipfname := fmt.Sprint(this.binDir, "/", fname, SUFGZIP)
						if !util.IsFileExist(gzipfname) {
							if util.Gzip(gzipfname, fname, this.binDir) == nil {
								os.Remove(this.binDir + "/" + fname)
							}
						}
					}
				}
			}
		} else {
			LoggerError.Error(err)
		}
	} else {
		LoggerError.Error(err)
	}
}

func (this *_binLog) readGzip(name string) (bs []byte) {
	buf := bytes.NewBuffer(make([]byte, 0))
	util.UnGzipByFile(name, func(bs []byte) bool {
		buf.Write(bs)
		return true
	})
	return buf.Bytes()
}

func (this *_binLog) ReadCurrentLog2GzipByte() (bs []byte, _err error) {
	this.mux4write.Lock()
	defer this.mux4write.Unlock()
	var _bs []byte
	if _bs, _err = util.ReadFile(this.binDir + "/" + sys.BINLOGNAME); _err == nil {
		if buf, err := util.GzipWrite(_bs); err == nil {
			bs = buf.Bytes()
		} else {
			_err = err
		}
	}
	return
}

func (this *_binLog) WriteBytes(buf *Buffer, t int64) (err error) {
	this.mux4write.RLock()
	defer this.mux4write.RUnlock()
	defer buf.Free()
	bs := buf.Bytes()
	if bs != nil && len(bs) > 0 {
		bakfn := ""
		if err, bakfn = Binlog.Write(bs); err == nil && bakfn != "" {
			if fileNum, err := strconv.Atoi(bakfn[len(this.prename)+1 : len(bakfn)-len(this.sufname)-1]); err == nil {
				LogStat.saveLogStat(int32(fileNum), t)
			}
		}
		if this.statBean.lastTime < t {
			this.statBean.lastTime, this.statBean.lastTxid = t, util.BytesToInt64(bs[16:24])
		}
	}
	return
}

func (this *_binLog) ReadLog(fileNum int32) (bs []byte) {
	filename := fmt.Sprint(this.binDir, "/", this.prename, "_", fileNum, ".", this.sufname)
	filegzip := fmt.Sprint(filename, SUFGZIP)
	if util.IsFileExist(filegzip) {
		bs = this.readGzip(filegzip)
	}
	return
}

func (this *_binLog) ReadGzipLogFile(fileNum int32) (bs []byte, filegzip string, err error) {
	filename := fmt.Sprint(this.binDir, "/", this.prename, "_", fileNum, ".", this.sufname)
	filegzip = fmt.Sprint(filename, SUFGZIP)
	if util.IsFileExist(filegzip) {
		bs, err = util.ReadFile(filegzip)
	} else {
		err = util.Errors(sys.ERR_FILENOTEXIST)
	}
	return
}

func (this *_binLog) GetLastTime() (_time, _txid int64) {
	_time, _txid = this.statBean.lastTime, this.statBean.lastTxid
	return
}

func (this *_binLog) getLastTime() (_r, _txid int64) {
	if r, bs := this._getTime(1); r > 0 && bs != nil {
		_r, _txid = r, util.BytesToInt64(bs[:8])
	}
	return
}

func (this *_binLog) _getTime(ty int8) (_r int64, _bs []byte) {
	logfile := fmt.Sprint(this.binDir, "/", sys.BINLOGNAME)
	if ff, err := os.Stat(logfile); err == nil {
		size := ff.Size()
		var block int64 = 1 << 10
		var off int64
		if size > block {
			off = size - block
		}
		i := int64(1)
		if f, err := os.Open(logfile); err == nil {
			defer f.Close()
			for off >= 0 {
				if ty == 0 {
					bs := make([]byte, len(STEP)+8+4)
					f.ReadAt(bs, 0)
					_r = util.BytesToInt64(bs[len(STEP)+4 : len(STEP)+4+8])
					break
				} else if ty == 1 {
					bs := make([]byte, size-off)
					f.ReadAt(bs, off)
					if site := bytes.LastIndex(bs, STEP); site != -1 {
						_r = util.BytesToInt64(bs[site+len(STEP)+4 : site+len(STEP)+8+4])
						_bs = bs[site+len(STEP)+8+4:]
						break
					} else {
						off = size - block*i
					}
				}
				i++
			}
		}
	}
	return
}

/****************************************************************************/
var LogStat *_statLog

type _statLog struct {
	mux *sync.Mutex
}

func NewStatLog(dir string, name string) (statlog *_statLog, err error) {
	return &_statLog{&sync.Mutex{}}, nil
}

func (this *_statLog) saveLogStat(fileNum int32, t int64) {
	this.mux.Lock()
	defer this.mux.Unlock()
	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(t))
	buf.Write(util.Int32ToBytes(fileNum))
	seq := atomic.AddInt64(&sys.STATSEQ, 1)
	sys.Level0Put(key.KeyLevel1.StatSeq(), util.Int64ToBytes(seq))
	sys.Level0Put(key.KeyLevel1.StatKey(seq), buf.Bytes())
}

func (this *_statLog) GetNum(t int64) (_t int64, _n int32, _e error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if sys.STATSEQ > 0 {
		for i := int64(1); i <= sys.STATSEQ; i++ {
			if v, err := sys.Level0Get(key.KeyLevel1.StatKey(i)); err == nil {
				ti := util.BytesToInt64(v[:8])
				if t < ti {
					_t, _n = ti, util.BytesToInt32(v[8:])
					break
				}
			} else {
				_e = err
				break
			}
		}
	}
	return
}

/*****************************************************************/

var BackLog *LogUtil

type LogUtil struct {
	mux          *sync.Mutex
	fname        string
	_fileHandler *os.File
}

func NewLogUtil(dir string, name string) (_LogUtil *LogUtil, err error) {
	fname := fmt.Sprint(dir, "/", name)
	var _fileHandler *os.File
	if _fileHandler, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		_LogUtil = &LogUtil{&sync.Mutex{}, fname, _fileHandler}
	}
	return
}

func (this *LogUtil) Write(bs []byte) (n int64, length int, e error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if ff, err := this._fileHandler.Stat(); err == nil {
		n = ff.Size()
		length, e = this._fileHandler.Write(bs)
	} else {
		e = err
	}
	return
}

func (this *LogUtil) ReadAt(at int64, length int32) (bs []byte, err error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	var f *os.File
	if f, err = os.Open(this.fname); err == nil {
		defer f.Close()
		bs = make([]byte, length)
		_, err = f.ReadAt(bs, at)
	}
	return
}

/***************************************************************************/
var CacheLog *_cacheLog

type _cacheLog struct {
	mux          *sync.Mutex
	fname        string
	_fileHandler *os.File
	fristStat    *fristStatBean
}

func NewCacheLog(dir string, name string) (_Log *_cacheLog, err error) {
	fname := fmt.Sprint(dir, "/", name)
	var _fileHandler *os.File
	if _fileHandler, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		_Log = &_cacheLog{&sync.Mutex{}, fname, _fileHandler, &fristStatBean{}}
	}
	return
}

func (this *_cacheLog) Write(bb *bytes.Buffer, txid, time int64) (e error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if this.fristStat.fristTime == 0 {
		this.fristStat.fristTime, this.fristStat.fristTxid = time, txid
	}
	bs := bb.Bytes()
	buf := util.BufferPool.Get(4 + len(bs))
	defer util.BufferPool.Put(buf)
	defer util.BufferPool.Put(bb)
	buf.Write(STEP) //4
	buf.Write(bs)   //
	_, e = this._fileHandler.Write(buf.Bytes())
	return
}

func (this *_cacheLog) Read() (bs []byte, err error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	fi, err := this._fileHandler.Stat()
	if fi.Size() > 0 {
		if bs, err = util.ReadFile(this.fname); err == nil {
			this._fileHandler.Close()
			os.Remove(this.fname)
			this._fileHandler, err = os.OpenFile(this.fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if this.fristStat == nil {
				this.fristStat = &fristStatBean{}
			}
		}
	}
	return
}
func (this *_cacheLog) GetFristStat() (_time, _txid int64) {
	return this.fristStat.fristTime, this.fristStat.fristTxid
}

/*********************************************************************/
func SysLog() (_r string) {
	syspath := sys.TLDB_SYS_LOG
	if sys.ROOTPATHLOG != "" {
		syspath = sys.ROOTPATHLOG + "/" + syspath
	}
	if bs, err := util.ReadFile(syspath); err == nil {
		_r = string(bs)
	}
	return
}
