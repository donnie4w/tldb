// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package sys

import (
	"errors"
	"fmt"

	. "github.com/donnie4w/tldb/container"
)

type ErrCodeType int16

// file not exist
var ERR_FILENOTEXIST ErrCodeType = 201

// undefined
var ERR_UNDEFINED ErrCodeType = 500

// no run stat
var ERR_NO_RUNSTAT ErrCodeType = 501

// no clus stat
var ERR_NO_CLUSTER ErrCodeType = 502

// incr seq error
var ERR_INCR_SEQ ErrCodeType = 503

// node not find
var ERR_NODE_NOFOUND ErrCodeType = 504

// the cluster nodes do not match
var ERR_CLUS_NOMATCH ErrCodeType = 505

// over time
var ERR_TIMEOUT ErrCodeType = 506

// batch fail
var ERR_BATCHFAIL ErrCodeType = 512

// repetition send
var ERR_RESEND ErrCodeType = 513

// tx is over
var ERR_TXOVER ErrCodeType = 514

// uuid re use
var ERR_UUID_REUSE ErrCodeType = 515

// sync data error
var ERR_SYNCDATA ErrCodeType = 516

// get data error
var ERR_GETDATA ErrCodeType = 517

// broadcast error
var ERR_BROADCAST ErrCodeType = 518

// re set stat failed
var ERR_SETSTAT ErrCodeType = 519

var ERR_PROXY ErrCodeType = 520

// Disallowed operation
var ERR_EPERM ErrCodeType = 521

// load log error
var ERR_LOADLOG ErrCodeType = 522

/*********************************************/
// the input parameters are incorrect
var ERR_NO_MATCH_PARAM ErrCodeType = 401

// table field error
var ERR_TABLE_FEILD_EXIST ErrCodeType = 408

// table exist
var ERR_TABLE_EXIST ErrCodeType = 409

// data no exist
var ERR_DATA_NOEXIST ErrCodeType = 410

// table not exist
var ERR_TABLE_NOEXIST ErrCodeType = 411

// column not exist
var ERR_COLUMN_NOEXIST ErrCodeType = 412

// index not exist
var ERR_IDX_NOEXIST ErrCodeType = 413

// select over time
var ERR_TIMEOUT_1 ErrCodeType = 414

// update over time
var ERR_TIMEOUT_2 ErrCodeType = 415

// insert over time
var ERR_TIMEOUT_3 ErrCodeType = 416

// create over time
var ERR_TIMEOUT_4 ErrCodeType = 417

// truncate over time
var ERR_TIMEOUT_5 ErrCodeType = 418

// column type error
var ERR_COLUMNTYPE ErrCodeType = 419

/*********************************************/
var ERR_AUTH_NOPASS ErrCodeType = 1301

/*********************************************/
type Exception struct {
	Code ErrCodeType
	m    *Map[int64, int8]
}

func NewException() *Exception {
	return &Exception{m: NewMap[int64, int8]()}
}

func (this *Exception) Error() (err error) {
	if this.Code > 0 {
		err = errors.New(fmt.Sprint(this.Code))
	}
	return
}

func (this *Exception) Put(uuid int64) {
	this.m.Put(uuid, 0)
}

func (this *Exception) Map() *Map[int64, int8] {
	return this.m
}
