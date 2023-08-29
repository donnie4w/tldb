// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package sys

import "fmt"

func SysLog(info string) (_r string) {
	a, b := "", ""
	ll := 80
	if ll >= len(info) {
		for i := 0; i < (ll-len(info))/2; i++ {
			a = fmt.Sprint(a, "=")
		}
		b = a
		if ll > len(info)+len(a)*2 {
			b = a + "="
		}
	}
	_r = fmt.Sprint(a, info, b)
	return
}
