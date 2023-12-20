// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
//

package key

import (
	"hash/crc32"
	"strconv"
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

func itoa(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ui32toa(i uint32)string{
	return strconv.FormatUint(uint64(i),10)
}
