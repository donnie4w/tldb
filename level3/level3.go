// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package level3

var Level3 = &_level3{}

type _level3 struct{}

func (this *_level3) PutMulti(ms map[string][]byte) (err error) {
	return
}

func (this *_level3) Put(key, value []byte) (err error) {
	return
}

func (this *_level3) Get(key []byte) (value []byte, err error) {
	return
}
