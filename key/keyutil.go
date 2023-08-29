// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package key

import (
	"hash/crc32"
	"strings"
)

func concat(ss ...string) string {
	var builder strings.Builder
	for _, v := range ss {
		builder.WriteString(v)
	}
	return builder.String()
}

func crc_32(bs []byte) uint32 {
	return crc32.ChecksumIEEE(bs)
}
