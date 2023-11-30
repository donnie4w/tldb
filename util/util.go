// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tldb
package util

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"hash/maphash"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/sys"
	"golang.org/x/exp/mmap"
)

func Concat(ss ...string) string {
	var builder strings.Builder
	for _, v := range ss {
		builder.WriteString(v)
	}
	return builder.String()
}

func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

func SHA1(s string) string {
	m := sha1.New()
	m.Write([]byte(s))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

func CRC32(bs []byte) uint32 {
	return crc32.ChecksumIEEE(bs)
}

func CRC64(bs []byte) uint64 {
	return crc64.Checksum(bs, crc64.MakeTable(crc64.ECMA))
}

var seed = maphash.MakeSeed()

func Hash(key []byte) uint64 {
	return maphash.Bytes(seed, key)
}

func Rand(i int) (_r int) {
	if i > 0 {
		_r = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(i)
	}
	return
}

func RandStrict(i int64) (_r int64, _err error) {
	if r, err := crand.Int(crand.Reader, big.NewInt(i)); err == nil {
		_r = r.Int64()
	} else {
		_err = err
	}
	return
}

func MatchString(pattern string, s string) bool {
	b, err := regexp.MatchString(pattern, s)
	if err != nil {
		b = false
	}
	return b
}

func Recovr() {
	if err := recover(); err != nil {
	}
}

func Time() time.Time {
	return time.Now().Add(time.Duration(sys.TIME_DEVIATION))
}

func TimeNano() int64 {
	return Time().UnixNano()
}

func TimeStr() string {
	return Time().Format("2006-01-02 15:04:05")
}

func TimeStrToTime(s string) (t time.Time, err error) {
	l, _ := time.LoadLocation("Asia/Shanghai")
	t, err = time.ParseInLocation("2006-01-02 15:04:05", s, l)
	return
}

func MapSub[K int | int8 | int32 | int64 | string, V any](m1, m2 *MapL[K, V]) (_r *MapL[K, V]) {
	_r = NewMapL[K, V]()
	if m1 != nil && m2 != nil {
		m1.Range(func(k K, v V) bool {
			if !m2.Has(k) {
				_r.Put(k, v)
			}
			return true
		})
	}
	return
}

func ArraySub[K int | int8 | int32 | int64 | string](a1, a2 []K) (_r []K) {
	_r = make([]K, 0)
	if a1 != nil && a2 != nil {
		m := make(map[K]byte, 0)
		for _, a := range a2 {
			m[a] = 0
		}
		for _, a := range a1 {
			if _, ok := m[a]; !ok {
				_r = append(_r, a)
			}
		}
	} else if a2 == nil {
		return a1
	}
	return
}

func ArraySub2[K int | int8 | int32 | int64 | string](a []K, k K) (_r []K) {
	if a != nil {
		_r = make([]K, 0)
		for _, v := range a {
			if v != k {
				_r = append(_r, v)
			}
		}
	}
	return
}

func MapToArray[K int | int8 | int32 | int64 | string, V any](m *Map[K, V]) (_r []K) {
	if m != nil {
		_r = make([]K, 0)
		m.Range(func(k K, _ V) bool {
			_r = append(_r, k)
			return true
		})
	}
	return
}

func MapToArray2[K int | int8 | int32 | int64 | string, V any](m map[K]V) (_r []K) {
	if m != nil && len(m) > 0 {
		_r = make([]K, len(m))
		i := 0
		for k := range m {
			_r[i] = k
			i++
		}
	}
	return
}

func ArrayToMap[K int | int8 | int32 | int64 | string, V int8](a []K) (_r *Map[K, V]) {
	_r = NewMap[K, V]()
	for _, v := range a {
		_r.Put(v, 0)
	}
	return
}

func ArrayToMap2[K int | int8 | int32 | int64 | string, V any](a []K, v V) (_r map[K]V) {
	_r = make(map[K]V, 0)
	for _, _a := range a {
		_r[_a] = v
	}
	return
}

func ArraySubMap[K int | int8 | int32 | int64 | string, V any](a []K, m *MapL[K, V]) (_r []K) {
	_r = make([]K, 0)
	if a != nil && m != nil {
		for _, v := range a {
			if !m.Has(v) {
				_r = append(_r, v)
			}
		}
	}
	return
}

func ArrayEqual[T int | int8 | int32 | int64 | string](a, b []T) bool {
	if a == nil || b == nil || len(a) != len(b) {
		return false
	}
	r1 := ArraySliceAsc(a)
	r2 := ArraySliceAsc(b)
	for i, v := range r1 {
		if r2[i] != v {
			return false
		}
	}
	return true
}

func ArrayContain[T int | int8 | int32 | int64 | string](a []T, v T) (b bool) {
	if a != nil {
		for _, _v := range a {
			if v == _v {
				b = true
				break
			}
		}
	}
	return
}

func Map2SyncMap[K int | int8 | int32 | int64 | string, V any](m1 map[K]V) (m2 *Map[K, V]) {
	if m1 != nil {
		m2 = NewMap[K, V]()
		for k, v := range m1 {
			m2.Put(k, v)
		}
	}
	return
}

func SyncMap2Map[K int | int8 | int32 | int64 | string, V any](m1 *Map[K, V]) (m2 map[K]V) {
	if m1 != nil {
		m2 = make(map[K]V, 0)
		m1.Range(func(k K, v V) bool {
			m2[k] = v
			return true
		})
	}
	return
}

func ArraySliceAsc[T int | int8 | int32 | int64 | string](a []T) (_r []T) {
	if a != nil {
		_r = make([]T, len(a))
		copy(_r, a)
		sort.Slice(_r, func(i, j int) bool { return _r[i] < _r[j] })
	}
	return
}

func NewUUID() uint32 {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(Int64ToBytes(int64(os.Getpid())))
	for i := 0; i < 100; i++ {
		buf.Write(Int64ToBytes(Time().UnixNano()))
	}
	if _r, err := RandStrict(1 << 31); err == nil && _r > 0 {
		buf.Write(Int64ToBytes(_r))
	}
	return CRC32(buf.Bytes())
}

/***********************************************************/
var __inc uint64

func inc() uint64 {
	return atomic.AddUint64(&__inc, 1)
}
func NewTxId() (txid int64) {
	b := make([]byte, 16)
	copy(b[0:8], Int64ToBytes(sys.UUID))
	copy(b[8:], Int64ToBytes(Time().UnixNano()))
	txid = int64(CRC32(b))
	txid = txid<<32 | int64(int32(inc()))
	if txid < 0 {
		txid = -txid
	}
	return
}

func Errors(code sys.ErrCodeType) error {
	return errors.New(fmt.Sprint(code))
}

func IsFileExist(path string) (_r bool) {
	if path != "" {
		_, err := os.Stat(path)
		_r = err == nil || os.IsExist(err)
	}
	return
}

func ReadFile(path string) (bs []byte, err error) {
	if r, er := mmap.Open(path); er == nil {
		defer r.Close()
		bs = make([]byte, r.Len())
		_, err = r.ReadAt(bs, 0)
	} else {
		err = er
	}
	return
}
