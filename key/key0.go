// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package key

const (
	KEY0_SYS = "0_0_"
)

type keyLevel0 struct {
}

var KeyLevel0 = &keyLevel0{}

func (this *keyLevel0) UUID() string {
	return concat(KEY0_SYS, "UUID")
}

func (this *keyLevel0) NAMESPACE() string {
	return concat(KEY0_SYS, "NAMESPACE")
}
