// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
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
	lm := NewSortMap[int64, int]()
	for i := 0; i < 100; i++ {
		lm.Put(time.Now().UnixNano()+int64(i), 0)
		// fmt.Println(lm.GetFrontKey())
	}
	lm.DelAndLoadBack()
	lm.BackForEach(func(k int64, v int) bool {
		fmt.Println(k)
		return true
	})
	// time.Sleep(4 * time.Second)
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
	m := NewLinkedMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func Benchmark_LimitMap(b *testing.B) {
	m := NewLimitMap[int, int](1 << 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkParallel_LimitMap(b *testing.B) {
	m := NewLimitMap[int, int](1 << 16)
	i := 0
	j := 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			j++
			m.Put(i, i)
			// m.Put(j, j)
		}
	})
	fmt.Println(i, " >> ", j, " >> ", m.count)
}

func Benchmark_LimitUnRepeatMap(b *testing.B) {
	m := NewLimitMap[int, int](1 << 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkParallel_LimitUnRepeatMap(b *testing.B) {
	m := NewLimitMap[int, int](1 << 16)
	i := 0
	j := 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			j++
			m.Put(i, i)
			m.Put(j, j)
		}
	})
	fmt.Println(i, " >> ", j, " >> ", m.count)
	k := int64(0)
	m.m.Range(func(key, value any) bool {
		k++
		return true
	})
	fmt.Println(k == m.count)
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
			fm.Put(k, k)
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

func Benchmark_LRU(b *testing.B) {
	m := NewLruCache[int64](1 << 20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Add(i, int64(i))
		m.Get(i)
	}
}

func BenchmarkParallel_LRU(b *testing.B) {
	m := NewLruCache[int64](1 << 20)
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			m.Add(i, int64(i))
			// m.Get(i)
		}
	})
}

func BenchmarkParallel_LinkedMap(b *testing.B) {
	m := NewLinkedMap[int, int]()
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			m.Put(i, i)
			// m.Get(i)
			// m.Del(i)
		}
	})
}

func BenchmarkParallel_MapL(b *testing.B) {
	m := NewMapL[int, int]()
	i := 0
	j := 1
	z := 10
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			j++
			z++
			m.Put(i, i)
			m.Put(j, j)
			m.Del(z)
		}
	})
	fmt.Println(i, " >> ", j, " >> ", m.Len())
	k := 0
	m.Range(func(_ int, _ int) bool {
		k++
		return true
	})
	fmt.Println(k == int(m.Len()))
}

func BenchmarkSortMap(b *testing.B) {
	lm := NewSortMap[int64, int]()
	for i := 0; i < b.N; i++ {
		lm.Put(time.Now().UnixNano()+int64(i), 0)
		lm.BackForEach(func(k int64, v int) bool {
			return true
		})
	}
}

func BenchmarkParallelSortMap(b *testing.B) {
	lm := NewSortMap[int64, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			lm.Put(time.Now().UnixNano()+int64(i), 0)
			// lm.DelAndLoadBack()
			// lm.GetFrontKey()
			lm.BackForEach(func(k int64, v int) bool {
				return true
			})
		}
	})
}

func Test_consistenthash(t *testing.T) {
	h := NewConsistenthash(1000)
	h.Add(1,2,3,4)
	h.Del(1)
	h.Add(4,5)
	m := map[int64]int{}
	for i := 0; i < 1000; i++ {
		node := h.Get(int64(i))
		if v, ok := m[node]; ok {
			m[node] = v + 1
		} else {
			m[node] = 1
		}
	}
	i:=0
	h.m.Range(func(k uint64, v int64) bool {
		i++
		return true
	})
	fmt.Println(i)
	fmt.Println(m)
}

func BenchmarkParallelconsistenthash(b *testing.B) {
	h := NewConsistenthash(1 << 14)
	h.Add(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			h.Get(int64(i))
		}
	})
}
