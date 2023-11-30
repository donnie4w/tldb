// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package tc

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/donnie4w/tldb/sys"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Monitor struct {
	Alloc        uint64
	TotalAlloc   uint64
	NumGC        uint32
	LastGC       uint64
	NextGC       uint64
	NumGoroutine int
	NumCPU       int
	CountPut     int64
	CcPut        int64
	CountGet     int64
	CcGet        int64
	RamUsage     float64
	DiskFree     uint64
	CpuUsage     float64
}

func monitorToJson() (_r string, err error) {
	var bs []byte
	if bs, err = json.Marshal(getMonitor()); err == nil {
		_r = string(bs)
	}
	return
}

func getMonitor() (_r *Monitor) {
	_r = &Monitor{}
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	_r.NumGoroutine = runtime.NumGoroutine()
	_r.NumCPU = runtime.NumCPU()

	_r.Alloc = rtm.Alloc
	_r.TotalAlloc = rtm.TotalAlloc
	_r.NumGC = rtm.NumGC
	_r.NextGC = rtm.NextGC
	_r.LastGC = rtm.LastGC
	_r.CcGet = sys.CcGet()
	_r.CcPut = sys.CcPut()
	_r.CountGet = sys.CountGet()
	_r.CountPut = sys.CountPut()

	if ram, err := getRAM(); err == nil {
		_r.RamUsage = float64(ram.UsedMB) / float64(ram.TotalMB)
	}

	if d, err := getDisk(); err == nil {
		_r.DiskFree = d.TotalGB - d.UsedGB
	}

	if c, err := getCPU(); err == nil {
		s := float64(0)
		for _, v := range c.Cpus {
			s += v
		}
		_r.CpuUsage = s
	}

	return
}

type Cpu struct {
	Cpus  []float64
	Cores int
}

type Ram struct {
	UsedMB  uint64
	TotalMB uint64
}

type Disk struct {
	UsedGB  uint64
	TotalGB uint64
}

func getRAM() (r Ram, err error) {
	if u, err := mem.VirtualMemory(); err == nil {
		r.UsedMB = u.Used / sys.MB
		r.TotalMB = u.Total / sys.MB
	}
	return r, nil
}

func getDisk() (d Disk, err error) {
	if u, err := disk.Usage("/"); err == nil {
		d.UsedGB = u.Used / sys.GB
		d.TotalGB = u.Total / sys.GB
	}
	return d, nil
}

func getCPU() (_r Cpu, err error) {
	_r.Cores, err = cpu.Counts(false)
	_r.Cpus, err = cpu.Percent(100*time.Millisecond, true)
	return
}
