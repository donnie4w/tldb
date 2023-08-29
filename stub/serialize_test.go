// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
package stub

import (
	"fmt"
	"testing"
)

func Test_serialize(t *testing.T) {
	ts := &TableStub{Idx: map[string]int8{"age": 0, "age1": 0}, Field: map[string][]byte{"name": []byte("hahaha"), "name2": []byte("hahaha123")}}
	ets := EncodeTableStub(ts)
	dts := DecodeTableStub(ets)
	fmt.Println(dts)
}

func Benchmark_encode(b *testing.B) {
	b.StopTimer()
	ts := &TableStub{"tom", 11, map[string]int8{"age": 0, "age1": 0}, map[string][]byte{"name": []byte("haheeeeeeeeeeeeeeeeeeeeeeeaha"), "naqqqqqqqqqqqqqqqqqqqqme2": []byte("hahaqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqha123")}}
	ets := EncodeTableStub(ts)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DecodeTableStub(ets)
	}
}
