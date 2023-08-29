// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
package container

import (
	"container/list"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := NewMapL[string, int8]()
	for k := 0; k < 10; k++ {
		go func() {
			for i := 0; i < 10000; i++ {
				m.Put(fmt.Sprint("dong", i), 0)
			}
		}()
	}
	time.Sleep(10 * time.Second)
	fmt.Println("i:", m.Len())

	j := 0
	for ; j < 100; j++ {
		m.Del(fmt.Sprint("dong", j/10))
	}
	fmt.Println("j:", j, "  ", m.Len())
}

func TestList(t *testing.T) {
	l := list.New()
	l.PushBack(1)
	fmt.Println(l.Front().Value) //1
	value := l.Remove(l.Front())
	fmt.Println(value) //1
	value1 := l.Remove(l.Front())
	fmt.Println(value1)
}

func TestLinkedMap(t *testing.T) {
	lm := NewLinkedMap[int64, int]()
	c := int64(1)
	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				lm.Put(atomic.AddInt64(&c, 1), 0)
				lm.Put(int64(j), 0)
			}
		}()
		go lm.FrontForEach(func(k int64, v int) bool {
			return true
		})
		lm.Back()
		lm.Front()

		lm.BackForEach(func(k int64, v int) bool {
			return true
		})
		fmt.Println(lm.Len())
	}

	for j := 0; j < int(c)-2; j++ {
		lm.Del(lm.Back())
	}
	time.Sleep(10 * time.Second)
	fmt.Println("--------------------------------")
	fmt.Println(c)
	fmt.Println(lm.Len())
	lm.FrontForEach(func(k int64, v int) bool {
		fmt.Println(k)
		return true
	})
}

func TestSortMap(t *testing.T) {
	// lm := NewSortMap[int64, int]()
	// for i := 0; i < 100; i++ {
	// 	go func() {
	// 		for j := 0; j < 10000; j++ {
	// 			lm.Put(time.Now().UnixMicro(), 0)
	// 		}
	// 	}()
	// 	lm.DelAndLoadBack()
	// 	lm.GetFrontKey()
	// 	lm.BackForEach(func(k int64, v int) bool {
	// 		return true
	// 	})
	// }
	// time.Sleep(10 * time.Second)
}

func TestLimitMap(t *testing.T) {
	lm := NewLimitMap[string, int](1 << 20)
	for i := 0; i < 100; i++ {
		// go func() {
		for j := 0; j < 10000; j++ {
			lm.Put(fmt.Sprint(i, ",", j), 0)
		}
		// }()
	}
	for j := 0; j < 10000; j++ {
		if !lm.Has(fmt.Sprint(0, ",", j)) {
			fmt.Println(fmt.Sprint(0, ",", j))
		}
	}
	fmt.Println(lm.Has(fmt.Sprint(0, ",", 9999)))
	fmt.Println(lm.Has(fmt.Sprint(0, ",", 10000)))
	fmt.Println(lm.count)
	// time.Sleep(10 * time.Second)
}

func TestBloomFilter(t *testing.T) {
	f := NewBloomFilter(1 << 10)
	// for i := 0; i < 1<<10; i++ {
	// 	go func() {
	for i := 1; i < 1<<11-1; i++ {
		f.Add([]byte("wu" + strconv.Itoa(i)))
	}
	// 	}()
	// }
	// time.Sleep(5 * time.Second)

	for i := 1; i < 1<<10; i++ {
		if !f.Contains([]byte("wu" + strconv.Itoa(i))) {
			fmt.Println("wu" + strconv.Itoa(i))
		}
	}
	// time.Sleep(10 * time.Second)
}

func Benchmark_Bloom(b *testing.B) {
	b.StopTimer()
	f := NewBloomFilter(1 << 15)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f.Add([]byte("wu" + strconv.Itoa(i)))
	}
}

func Benchmark_Map(b *testing.B) {
	b.StopTimer()
	m := NewMapL[int, int]()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
		// m.Get(i)
	}
}

func Benchmark_LinkedMap(b *testing.B) {
	b.StopTimer()
	m := NewLinkedMap[int, int]()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func Benchmark_LimitMap(b *testing.B) {
	b.StopTimer()
	m := NewLimitMap[int, int](1 << 20)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func TestSwap(t *testing.T) {
	var custor int64 = 0
	for i := int64(0); i <= 2000000; i++ {
		go func(i int64) {
			if custor < i {
				custor = i
			}
			// swapbig(&custor, i)
		}(i)
	}
	<-time.After(5 * time.Second)
	fmt.Println(custor)
}

func Test_fastLinkMap(t *testing.T) {
	fm := NewFastLinkedMap[int64, int64]()
	for i := int64(0); i <= 1000000; i++ {
		go func(k int64) {
			fm.Put(k, i)
		}(i)
	}
	<-time.After(5 * time.Second)
	fm.RangeAndDelete(func(id int64, k int64, v int64) bool {
		if k < 999990 {
			return true
		} else {
			return false
		}
	})
	fm.m.Range(func(k int64, v *kv[int64, int64]) bool {
		fmt.Println(k)
		return true
	})
	<-time.After(10 * time.Second)
}
