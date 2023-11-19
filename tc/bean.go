// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

package tc

import (
	"github.com/donnie4w/tldb/stub"
)

type SysVar struct {
	StartTime      string
	LocalTime      string
	Time           string
	UUID           int64
	RUNUUIDS       string
	STAT           string
	TIME_DEVIATION string
	ADDR           string
	STORENODENUM   string
	CLUSTER_NUM    string
	MQADDR         string
	CLIADDR        string
	ADMINADDR      string
	CCPUT          int64
	CCGET          int64
	COUNTPUT       int64
	COUNTGET       int64
	SyncCount      int64
}

type SysVarView struct {
	Stat bool
	Show string
	SYS  *SysVar
	RN   []*stub.RemoteNode
}

/**********************************************************************************/
type AdminView struct {
	Show       string
	AdminUser  map[string]string
	CliUser    []string
	MqUser     []string
	Stat       bool
	Init       bool
	ShowCreate bool
}

/**********************************************************************************/

type Tables struct {
	Name    string
	Columns []string
	Idxs    []string
	Seq     int64
	Sub     int64
}

type TData struct {
	Name    string
	Id      int64
	Columns map[string]string
}

type SelectBean struct {
	Name        string
	Id          string
	ColumnName  string
	ColumnValue string
	StartId     string
	Limit       string
}

type DataView struct {
	Tb      []*Tables
	Tds     []*TData
	ColName map[string][]byte
	Sb      *SelectBean
	Stat    bool
}

/**********************************************************************************/
type SysParam struct {
	DBFILEDIR         string
	MQTLS             bool
	ADMINTLS          bool
	CLITLS            bool
	CLICRT            string
	CLIKEY            string
	MQCRT             string
	MQKEY             string
	ADMINCRT          string
	ADMINKEY          string
	COCURRENT_PUT     int64
	COCURRENT_GET     int64
	DBMode            int
	NAMESPACE         string
	VERSION           string
	BINLOGSIZE        int64
	ADDR              string
	CLIADDR           string
	MQADDR            string
	WEBADMINADDR      string
	CLUSTER_NUM       int
	PWD               string
	PUBLICKEY         string
	PRIVATEKEY        string
	CLUSTER_NUM_FINAL bool
}

type SysParamView struct {
	SYS  *SysParam
	Stat bool
}

/**********************************************************************************/
type AlterTable struct {
	TableName   string
	ID          int64
	Columns     map[string]*FieldInfo
	ColumnValue map[string]string
}

type FieldInfo struct {
	Idx   bool
	Type  string
	Tname string
}
