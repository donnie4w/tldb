// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

package keystore

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/donnie4w/tldb/sys"
	"github.com/donnie4w/tldb/util"
)

func Init() {
	if err := InitKey(sys.DBFILEDIR); err != nil {
		panic("keyStore init failed")
	}

	if sys.PRIVATEKEY != "" || sys.PUBLICKEY != "" {
		a := fmt.Sprint(time.Now().Nanosecond())
		var err error
		var bs []byte
		var ok bool
		if bs, err = RsaEncrypt([]byte(a), sys.PUBLICKEY); err == nil {
			if bs, err = RsaDecrypt(bs, sys.PRIVATEKEY); err == nil {
				ok = a == string(bs)
			}
		}
		if err != nil || !ok {
			panic("publickey and privatekey authFailed")
		}
	}
}

func InitKey(dir string) (err error) {
	if KeyStore, err = NewKeyStore(dir, "keystore"); err == nil {
		StoreAdmin.Load()
	}
	return
}

var KeyStore *_keyStore

type _keyStore struct {
	mux          *sync.Mutex
	fname        string
	_fileHandler *os.File
}

func NewKeyStore(dir string, name string) (ks *_keyStore, err error) {
	fname := fmt.Sprint(dir, "/", name)
	var _fileHandler *os.File
	if _fileHandler, err = os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666); err == nil {
		ks = &_keyStore{&sync.Mutex{}, fname, _fileHandler}
	}
	return
}

func (this *_keyStore) Write(bs []byte) (err error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this._fileHandler.Seek(0, io.SeekStart)
	this._fileHandler.Truncate(0)
	var obs []byte
	if obs, err = util.ZlibCz(bs); err != nil {
		obs = bs
	}
	_, err = this._fileHandler.Write(obs)
	return
}

func (this *_keyStore) Read() (bs []byte, err error) {
	defer this.mux.Unlock()
	this.mux.Lock()
	fi, err := this._fileHandler.Stat()
	if fi.Size() > 0 {
		if bs, err = util.ReadFile(this.fname); err == nil && bs != nil {
			if obs, er := util.ZlibUnCz(bs); er == nil {
				return obs, er
			}
		}
	} else {
		err = errors.New("empty file")
	}
	return
}
