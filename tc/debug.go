// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package tc

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strings"

	"github.com/donnie4w/tldb/log"
	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

func tlDebug() {
	if sys.DEBUGADDR != "" {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		if !strings.Contains(sys.DEBUGADDR, ":") && util.MatchString("^[0-9]{4,5}$", sys.DEBUGADDR) {
			sys.DEBUGADDR = fmt.Sprint(":", sys.DEBUGADDR)
		}
		log.LoggerSys.Info(sys.SysLog(fmt.Sprint("Debug start[", sys.DEBUGADDR, "]")))
		if err := http.ListenAndServe(sys.DEBUGADDR, nil); err != nil {
			panic("tl debug failed:" + err.Error())
		}
	}
}
