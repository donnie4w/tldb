// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file

package sys

import (
	"fmt"
)

func FmtLog(v ...any) {
	info := fmt.Sprint(v...)
	a, b := "", ""
	ll := 80
	if ll >= len(info) {
		for i := 0; i < (ll-len(info))/2; i++ {
			a = a + "="
		}
		b = a
		if ll > len(info)+len(a)*2 {
			b = a + "="
		}
	}
	Log().Info(a, info, b)
}

func BlankLine() {
	Log().Write([]byte("\n"))
}

func timlogo() {
	_r := `
              =============================================================================
              ===        =============   ===           =======        =======           ===
              ===             ===        ===           ===   ===      ===   ===         ===
              ===             ===        ===           ===    ===     ===    ===        ===
              ===             ===        ===           ===     ===    ===   ===         ===
              ===             ===        ===           ===     ===    =======           ===
              ===             ===        ===           ===     ===    ===   ===         ===
              ===             ===        ===           ===    ===     ===    ===        ===
              ===             ===        ===           ===   ===      ===   ===         ===
              ===             ===        ===========   =======        =======           ===
              =============================================================================
	`
	Log().Info(_r)
}
