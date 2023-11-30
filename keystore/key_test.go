// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
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
