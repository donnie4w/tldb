// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"testing"

	. "github.com/donnie4w/tldb/container"
	"github.com/donnie4w/tldb/stub"
	"github.com/golang/snappy"
)

// func TestFilter(t *testing.T) {
// 	f := NewBloomFilter(10, 1<<10)
// 	for i := 0; i < 1<<10+1; i++ {
// 		f.Add([]byte("wu" + strconv.Itoa(i)))
// 	}
// 	fmt.Println(time.Now())
// 	for i := 0; i < (1<<10 + 2); i++ {
// 		if !f.Contains([]byte("wu" + strconv.Itoa(i))) {
// 			fmt.Println("wu" + strconv.Itoa(i))
// 		}
// 	}

// 	// fmt.Println(time.Now())
// 	// fmt.Println(f.Contains([]byte("wu99")))
// 	// fmt.Println(f.Contains([]byte("wu1000")))
// 	// fmt.Println(f.Contains([]byte("wu" + fmt.Sprint(1<<20-1))))
// 	// fmt.Println(f.Contains([]byte("wu" + fmt.Sprint(1<<20+9))))
// }

// func TestFifo(t *testing.T) {
// 	f := NewFiFo[string, int64](3)
// 	f.Add("wu", 111)
// 	fmt.Println(f.Get("wu"))
// 	f.Add("wu1", 222)
// 	f.Add("wu2", 333)
// 	f.Add("wu4", 444)
// 	fmt.Println(f.Get("wu"))
// 	fmt.Println(f.Get("wu1"))
// 	fmt.Println(f.Get("wu2"))
// 	fmt.Println(f.Get("wu4"))
// 	fmt.Println(time.Now())
// }

// func TestList(t *testing.T) {
// 	li := NewList[int64]()
// 	li.PushFront(1)
// 	li.PushFront(2)
// 	fmt.Println(li.Front())
// 	fmt.Println(li.Len())
// 	li.DelBack()
// 	li.DelBack()
// 	li.DelBack()
// 	fmt.Println(li.Len())
// 	fmt.Println(li.Front())
// 	for i := 0; i < 10; i++ {
// 		li.PushFront(int64(i))
// 	}

// 	li.DelBack()
// 	li.DelBack()

// 	li.BackPrev(func(t int64) bool {
// 		fmt.Println("--->", t)
// 		if t == 0 {
// 			li.PushFront(11)
// 			fmt.Println(li.Front())
// 		}
// 		return true
// 	})

// 	li.FrontNext(func(t int64) bool {
// 		fmt.Println("====>", t)
// 		if t == 0 {
// 			li.PushBack(11)
// 			fmt.Println(li.GetBack())
// 		}
// 		return true
// 	})
// }

// func TestSortMap(t *testing.T) {
// 	sm := NewSortMap[int, string]()
// 	sm.Put(19, "wwww19")
// 	sm.Put(20, "wwww20")
// 	sm.Put(17, "wwww17")
// 	sm.Put(22, "wwww22")
// 	sm.Put(14, "wwww14")
// 	sm.Put(66, "wwww66")
// 	sm.Put(1, "wwww1")
// 	fmt.Println(sm)
// 	for sm.Len() > 0 {
// 		fmt.Println(sm.DelAndLoadBack())
// 	}
// }

// func TestMap(t *testing.T) {
// 	m := NewMap[string, int8]()
// 	for k := 0; k < 100; k++ {
// 		go func() {
// 			for i := 0; i < 100; i++ {
// 				m.Put(fmt.Sprint("dong", i), 0)
// 			}
// 		}()
// 	}
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("i:", m.Len())

// 	j := 0
// 	for ; j < 100; j++ {
// 		m.Del(fmt.Sprint("dong", j/10))
// 	}
// 	fmt.Println("j:", j, "  ", m.Len())
// }

// func TestLinkedMap(t *testing.T) {
// 	lm := NewLinkedMap[int, string]()

// 	for j := 0; j < 10; j++ {
// 		go func() {
// 			for i := 0; i < 10; i++ {
// 				lm.Put(i, fmt.Sprint("dong", i))
// 			}
// 		}()
// 	}
// 	time.Sleep(10 * time.Second)
// 	lm.BackForEach(func(k int, v string) bool {
// 		fmt.Println(k, "===>", v)
// 		return true
// 	})
// 	lm.Put(5, "555")
// 	lm.FrontForEach(func(k int, v string) bool {
// 		fmt.Println(k, "=====>", v)
// 		return true
// 	})
// 	fmt.Println(lm.Prev(5))
// 	fmt.Println(lm.Next(5))
// }

func TestCompareArray(t *testing.T) {
	a := []int{1, 2, 3, 6, 5}
	b := []int{3, 1, 2, 5, 6}
	fmt.Println(ArraySub(a, b))
	fmt.Println(ArraySub(b, a))
	fmt.Println(ArrayEqual(a, b))
}

func TestRand(t *testing.T) {
	m := NewMapL[int64, int8]()
	i := 0
	for i = 0; i < 1<<19; i++ {
		if k, err := RandStrict(Time().UnixNano()); err == nil {
			m.Put(k, 1)
		}
	}
	fmt.Println(i)
	fmt.Println(m.Len())
	fmt.Println(RandStrict(Time().UnixNano()))
}

func TestDataToBytes(t *testing.T) {
	arr := []int64{11<<20, 22222222, 3, 4<<10}
	bs := IntArrayToBytes(arr)
	fmt.Println(bs)
	arr2 := BytesToIntArray(bs)
	fmt.Println(arr2)
}

func TestArraySliceAsc(t *testing.T) {
	arr := []int64{1, 1, 21, 13, 4}
	v2 := ArraySliceAsc(arr)
	fmt.Println(v2)
}

func TestArrayEqual(t *testing.T) {
	a1 := []int64{1, 1, 21, 13, 4}
	a2 := []int64{1, 1, 21, 13, 1}
	a3 := []int64{1, 1, 21, 13, 4}
	fmt.Println(ArrayEqual(a1, a3))
	fmt.Println(ArrayEqual(a1, a2))
}

func TestTDecode(t *testing.T) {
	mb := new(stub.MqBean)
	mb.ID = 111111921123
	bs := TEncode(mb)
	fmt.Println(bs)
	mb2, _ := TDecode(bs, &stub.MqBean{})
	fmt.Println(mb2.ID)
}

func TestZlib(t *testing.T) {
	in := []byte("123456789")
	bs, err := ZlibCz(in)
	fmt.Println(err)
	fmt.Println(string(bs))
	bs, err = ZlibUnCz(bs)
	fmt.Println(err)
	fmt.Println(string(bs))
}

func Benchmark_md5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5("1234567890qwertyuiop")
	}
}

func Benchmark_crc32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CRC32([]byte("1234567890qwertyuiop1234567890qwertyuiop1234567890qwertyuiop1234567890qwertyuiop"))
	}
}

func Benchmark_crc64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CRC64([]byte("1234567890qwertyuiop1234567890qwertyuiop1234567890qwertyuiop1234567890qwertyuiop"))
	}
}

func Benchmark_bs2int(b *testing.B) {
	b.StopTimer()
	bs := Int64ToBytes(99)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BytesToInt64(bs)
	}
}

func Benchmark_int64Tobs(b *testing.B) {
	b.StopTimer()
	t := Time().UnixNano()
	bs := Int64ToBytes(t)
	b.StartTimer()
	var r int64
	for i := 0; i < b.N; i++ {
		r = BytesToInt64(bs)
	}
	fmt.Println(t == r)
}

func Test_int64Tobs(t *testing.T) {
	for i := int64(1 << 1); i < 1000; i++ {
		if bs := Int64ToBytes(i); i != BytesToInt64(bs) {
			panic("err >>" + fmt.Sprint(i))
		}
	}
	fmt.Println(BytesToInt64([]byte{0, 1}))
	fmt.Println("ok")
}

func Test_int32Tbs(t *testing.T) {
	for i := int32(0); i < 1<<30; i++ {
		if bs := Int32ToBytes(i); i != BytesToInt32(bs) {
			panic("err >>" + fmt.Sprint(i))
		}
	}
	fmt.Println("ok")
}

func Benchmark_maphash(b *testing.B) {
	ib := Hash([]byte("1234567789qwertyuiop"))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if Hash([]byte("1234567789qwertyuiop")) != ib {
				panic("err")
			}
		}
	})
}

func Benchmark_TEncode(b *testing.B) {
	b.StopTimer()
	mb := new(stub.MqBean)
	mb.ID = 111111921123
	mb.Msg = []byte("aaaaaaaabbbbbbbbbbbbbbaaaaaaaa")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TEncode(mb)
	}
}

func Benchmark_TDecode(b *testing.B) {
	b.StopTimer()
	ts := &stub.TableStub{Tablename: "tom", ID: 11, Idx: map[string]int8{"age": 12, "age1": 13}, Field: map[string][]byte{"name": []byte("haheeeeeeeeeeeeeeeeeeeeeeeaha"), "naqqqqqqqqqqqqqqqqqqqqme2": []byte("hahaqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqha123")}}
	// bs := TEncode(ts)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TEncode(ts)
		// TDecode(bs, &stub.TableStub{})
	}
}

func Benchmark_czlib(b *testing.B) {
	b.StopTimer()
	bs, _ := ReadFile("gob.go")
	fmt.Println("len(bs)>>", len(bs))
	var r []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r, _ = ZlibCz(bs)
	}
	fmt.Println("len(r)>>", len(r))
}

func Benchmark_snappy(b *testing.B) {
	b.StopTimer()
	bs, _ := ReadFile("gob.go")
	fmt.Println("len(bs)>>", len(bs))
	var dst []byte
	dst = snappy.Encode(nil, bs)
	var bs2 []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bs2, _ = snappy.Decode(nil, dst)
	}
	fmt.Println("len(r)>>", len(bs2))
}

func TestIntByte(t *testing.T) {
	b1 := Int16ToBytes(1<<15 - 1)
	fmt.Println("b >>", BytesToInt16(b1))
	b2 := Int32ToBytes(1<<31 - 1)
	fmt.Println("b >>", BytesToInt32(b2))
	b3 := Int64ToBytes(1<<63 - 1)
	fmt.Println("b >>", BytesToInt64(b3))
	var i float32 = 0.01
	var byt bytes.Buffer
	binary.Write(&byt, binary.BigEndian, i)
	binary.Read(&byt, binary.BigEndian, &i)
	fmt.Println(i)
	j, _ := strconv.ParseFloat("0.11", 32)
	fmt.Println(float32(j))
}

func Test_Gzip(t *testing.T) {
	buf, err := GzipWrite([]byte("hello123"))
	buf, err = UnGzip(buf.Bytes())
	fmt.Println(err)
	fmt.Println(string(buf.Bytes()))
}
