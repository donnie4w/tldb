// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package log

import (
	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/tldb/sys"
)

var STEP = []byte{'|', '|', '|', '|'}
var logFormat = logging.FORMAT_SHORTFILENAME | logging.FORMAT_DATE | logging.FORMAT_TIME | logging.FORMAT_MICROSECNDS
var Logger = logging.NewLogger().SetLevel(logging.LEVEL_INFO).SetFormat(logFormat)
var LoggerError = logging.NewLogger().SetLevel(logging.LEVEL_ERROR).SetFormat(logFormat)
var LoggerSys = logging.NewLogger().SetFormat(logging.FORMAT_DATE | logging.FORMAT_TIME | logging.FORMAT_MICROSECNDS)

var LogBIN = logging.NewLogger()

func LogInit() {
	var err error
	LogBIN.SetGzipOn(true)
	if LogBIN, err = LogBIN.SetRollingFile(sys.DBFILEDIR+"/bin", sys.BINLOGNAME, sys.BINLOGSIZE, logging.MB); err != nil {
		LoggerSys.Error("bin log init failed:", err)
		panic("bin log init failed:" + err.Error())
	}

	if LogStat, err = NewStatLog(sys.DBFILEDIR, sys.STATLOGNAME); err != nil {
		LoggerSys.Error("stat log init failed:", err)
		panic("stat log init failed:" + err.Error())
	}

	BinLog = NewBinLog()

	if BackLog, err = NewLogUtil(sys.DBFILEDIR, sys.BACKLOGNAME); err != nil {
		LoggerSys.Error("back log init failed:", err)
		panic("back log init failed:" + err.Error())
	}

	if CacheLog, err = NewCacheLog(sys.DBFILEDIR, sys.CACHELOGNAME); err != nil {
		LoggerSys.Error("cache log init failed:", err)
		panic("cache log init failed:" + err.Error())
	}

	if !sys.LOGON {
		Logger.SetLevel(logging.LEVEL_OFF)
	}

	if _, err = Logger.SetRollingFile(sys.ROOTPATHLOG, "tldb.log", 1, logging.GB); err != nil {
		panic("log init failed:" + err.Error())
	}

	if _, err = LoggerError.SetRollingFile(sys.ROOTPATHLOG, "tldb_error.log", 1, logging.GB); err != nil {
		panic("log init failed:" + err.Error())
	}

	if _, err = LoggerSys.SetRollingFile(sys.ROOTPATHLOG, sys.TLDB_SYS_LOG, 1, logging.MB); err != nil {
		panic("log init failed:" + err.Error())
	}
}
