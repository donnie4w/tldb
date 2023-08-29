// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
package level1

import (
	"fmt"
	"testing"
)

func Test_batch(t *testing.T) {
	bc := &BatchPacket{Kvmap: map[string][]byte{"name": []byte("hahaha"), "name2": []byte("hahaha123")}, Dels: []string{"aaa", "bbb"}}
	ets := EncodeBatchPacket(bc)
	fmt.Println(len(ets.Bytes()))
	dts := DecodeBatchPacket(ets.Bytes())
	fmt.Println(dts)

}

func Benchmark_batch(b *testing.B) {
	b.StopTimer()
	bc := &BatchPacket{Kvmap: map[string][]byte{"name": {0, 1}, "name1": []byte("hahakkkkeifjeifaiefniinenfiefaha"), "name2": []byte("haha88888888888888888ha123")}, Dels: []string{"aaa", "bbb"}}
	ets := EncodeBatchPacket(bc)
	DecodeBatchPacket(ets.Bytes())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EncodeBatchPacket(bc)
		DecodeBatchPacket(ets.Bytes())
	}
}

func Test_PonBean(t *testing.T) {
	pb := newPonBean(1, 2, 3, []int64{1, 2, 3}, nil)
	bc := &BatchPacket{Kvmap: map[string][]byte{"name": []byte("hahaha"), "name2": []byte("hahaha123")}, Dels: []string{"aaa", "bbb"}}
	pb.Batch = bc
	i := int64(1234567)
	pb.SyncId = &i
	pb.Pbtime = []byte{1, 2, 3, 4, 5}
	pb.Key = []byte("key")
	epb := EncodePonBean(pb)
	fmt.Println(len(epb.Bytes()))
	dts := DecodePonBean(epb.Bytes())
	fmt.Println(dts, "  >>>", *dts.SyncId)
}

func Test_DataBeen(t *testing.T) {
	db := newDataBeen(true, []int64{1, 2, 3}, []byte{4, 5, 6})
	edb := EncodeDataBean(db)
	fmt.Println(len(edb.Bytes()))
	ddb := DecodeDataBeen(edb.Bytes())
	fmt.Println(ddb)
}

func Benchmark_PonBean(b *testing.B) {
	b.StopTimer()
	pb := newPonBean(1, 2, 3, []int64{1, 2, 3}, []int64{5, 6, 7})
	bc := &BatchPacket{Kvmap: map[string][]byte{"name": []byte("hahaha"), "name2": []byte("hahaha123")}, Dels: []string{"aaa", "bbb"}}
	pb.Batch = bc
	bs := EncodePonBean(pb)
	// bs2 := util.TEncode(pb)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//
		DecodePonBean(bs.Bytes())
		// util.TDecode(bs2, &PonBean{})
	}
}

func Test_map2byte(t *testing.T) {
	m := map[int64]int8{22: 1, 33: 2, 44: 3}
	buf := mapTobytes(m)
	m2 := bytesToMap(buf.Bytes())
	fmt.Println(m2)
}
