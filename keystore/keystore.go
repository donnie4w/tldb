// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package keystore

import (
	"sync"

	"github.com/donnie4w/tldb/util"
)

var StoreAdmin = NewStoreAdmin()

type storeAdmin struct {
	kb  *KeyBean
	mux *sync.Mutex
}

func NewStoreAdmin() *storeAdmin {
	return &storeAdmin{mux: &sync.Mutex{}}
}

func (this *storeAdmin) Load() {
	defer this.mux.Unlock()
	this.mux.Lock()
	var err error
	var bs []byte
	if bs, err = KeyStore.Read(); err == nil {
		this.kb, err = util.TDecode(bs, &KeyBean{})
	}
	if err != nil {
		this.kb = &KeyBean{make(map[string]*UserBean, 0), make(map[string]*UserBean, 0), make(map[string]*UserBean, 0), make(map[string]string, 0)}
	}
}

func (this *storeAdmin) PutAdmin(name, pwd string, _type int8) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.kb.Admin[name] = &UserBean{name, util.MD5(pwd), _type}
	KeyStore.Write(util.TEncode(this.kb))
}

func (this *storeAdmin) DelAdmin(name string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if _, ok := this.kb.Admin[name]; ok {
		delete(this.kb.Admin, name)
		KeyStore.Write(util.TEncode(this.kb))
	}
}

func (this *storeAdmin) GetAdmin(name string) (_r *UserBean, ok bool) {
	defer this.mux.Unlock()
	this.mux.Lock()
	_r, ok = this.kb.Admin[name]
	return
}

func (this *storeAdmin) AdminList() (ss []string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	ss = make([]string, 0)
	for k := range this.kb.Admin {
		ss = append(ss, k)
	}
	return
}

func (this *storeAdmin) PutClient(name, pwd string, _type int8) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.kb.Client[name] = &UserBean{name, util.MD5(pwd), _type}
	KeyStore.Write(util.TEncode(this.kb))
}

func (this *storeAdmin) DelClient(name string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if _, ok := this.kb.Client[name]; ok {
		delete(this.kb.Client, name)
		KeyStore.Write(util.TEncode(this.kb))
	}
}

func (this *storeAdmin) GetClient(name string) (_r *UserBean, ok bool) {
	defer this.mux.Unlock()
	this.mux.Lock()
	_r, ok = this.kb.Client[name]
	return
}

func (this *storeAdmin) ClientList() (ss []string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	ss = make([]string, 0)
	for k := range this.kb.Client {
		ss = append(ss, k)
	}
	return
}

func (this *storeAdmin) PutMq(name, pwd string, _type int8) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.kb.Mq[name] = &UserBean{name, util.MD5(pwd), _type}
	KeyStore.Write(util.TEncode(this.kb))
}

func (this *storeAdmin) DelMq(name string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if _, ok := this.kb.Mq[name]; ok {
		delete(this.kb.Mq, name)
		KeyStore.Write(util.TEncode(this.kb))
	}
}

func (this *storeAdmin) GetMq(name string) (_r *UserBean, ok bool) {
	defer this.mux.Unlock()
	this.mux.Lock()
	_r, ok = this.kb.Mq[name]
	return
}

func (this *storeAdmin) MqList() (ss []string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	ss = make([]string, 0)
	for k := range this.kb.Mq {
		ss = append(ss, k)
	}
	return
}

func (this *storeAdmin) PutOther(key, value string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	this.kb.Other[key] = value
	KeyStore.Write(util.TEncode(this.kb))
}

func (this *storeAdmin) DelOther(key string) {
	defer this.mux.Unlock()
	this.mux.Lock()
	if _, ok := this.kb.Other[key]; ok {
		delete(this.kb.Other, key)
		KeyStore.Write(util.TEncode(this.kb))
	}
}

func (this *storeAdmin) GetOther(key string) (value string, ok bool) {
	defer this.mux.Unlock()
	this.mux.Lock()
	value, ok = this.kb.Other[key]
	return
}
