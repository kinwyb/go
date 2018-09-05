package db

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

//并发测试ID生成是否唯一
func Test_IDGenerator(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU()) //多核计算
	lock := &sync.RWMutex{}
	m := map[int64]int64{}
	count := int64(0)
	w := &sync.WaitGroup{}
	convey.Convey("ID生成器", t, func() {
		for a := 0; a < 32; a++ {
			w.Add(10)
			go func(a int) {
				idWorker := NewIDGenerator(int64(a))
				for j := 0; j < 10; j++ { //每个生成器10个协程同时生成id
					go func(idWorker *IDGenerator) {
						var id int64
						for i := 0; i < 100000; i++ {
							id, _ = idWorker.NextId()
							if id < 0 {
								w.Done()
								panic("出现负数")
							}
							lock.Lock()
							if _, ok := m[id]; ok {
								atomic.AddInt64(&count, 1)
							} else {
								m[id] = id
							}
							lock.Unlock()
						}
						w.Done()
					}(idWorker)
				}
				w.Done()
			}(a)
		}
		w.Wait()
		convey.So(count, convey.ShouldEqual, 0)
	})
}

func BenchmarkIDGenerator_NextString(b *testing.B) {
	idWorker := NewIDGenerator(10)
	for i := 0; i < b.N; i++ {
		idWorker.NextString()
	}
}

func BenchmarkIDGenerator_Next(b *testing.B) {
	idWorker := NewIDGenerator(10)
	for i := 0; i < b.N; i++ {
		idWorker.NextId()
	}
}

func BenchmarkIDGenerator_NextStringRunParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		idWorker := NewIDGenerator(1)
		for pb.Next() {
			idWorker.NextString()
		}
	})
}
