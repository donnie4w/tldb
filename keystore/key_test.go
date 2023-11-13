// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package keystore

import (
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	msg, err := RsaEncrypt([]byte("123456789"), "")
	fmt.Println(err)
	ct, err := RsaDecrypt(msg, "")
	fmt.Println(err)
	fmt.Println(string(ct))
}
