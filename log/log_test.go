// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.

package log

import (
	"fmt"
	"testing"

	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/tldb/sys"
)

func TestLog(t *testing.T) {
	sys.DBFILEDIR = "_dataTest"
	BinLog = NewBinLog()
	var err error
	if LogBIN, err = logging.NewLogger().SetRollingFile(sys.DBFILEDIR, sys.BINLOGNAME, 2, logging.MB); err != nil {
		panic("bin log init failed:" + err.Error())
	}
	if LogStat, err = NewStatLog(sys.DBFILEDIR, sys.STATLOGNAME); err != nil {
		panic("stat log init failed:" + err.Error())
	}
	// for i := 0; i < 50000; i++ {
	// 	BinLog.Write([]byte("1234567890-=qwertyuiopasdfghjklzxcvbnm[],./!@#$%^&*(){}<>"))
	// }
	BinLog.gzipLog()
	// LogStat.msort.BackForEach(func(k int64, v int32) bool {
	// 	fmt.Println("k:", k, "  v:", v)
	// 	return true
	// })
	bs := BinLog.ReadLog(2)
	fmt.Println(len(bs))
	fmt.Println(LogStat.GetNum(1682874648141243601))
	// fmt.Println(BinLog.GetFirstTime())
	_r, id := BinLog.GetLastTime()
	fmt.Println(_r)
	// fmt.Println(LogStat.GetNextNum(0))
	// fmt.Println(LogStat.GetNextNum(1682927170124950900))
	fmt.Println(id)
	// bs, _ = BinLog.ReadCurrentLog2GzipByte()
	// fmt.Println(string(bs))
}
