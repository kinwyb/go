package heldiamgo

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewMs15LenIDGenerator(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU()) //多核计算
	lock := &sync.RWMutex{}
	m := map[int64]int64{}
	count := int64(0)
	w := &sync.WaitGroup{}
	convey.Convey("ID生成器", t, func() {
		for a := 0; a < 32; a++ {
			w.Add(10)
			go func(a int) {
				idWorker := NewMs15LenIDGenerator(int64(a))
				for j := 0; j < 10; j++ { //每个生成器10个协程同时生成id
					go func(idWorker *Ms15LenIDGenerator) {
						var id int64
						for i := 0; i < 10000; i++ {
							id, _ = idWorker.NextID()
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
			}(a)
		}
		w.Wait()
		fmt.Printf("一共生成[%d]个,重复数[%d]", len(m), count)
		convey.So(count, convey.ShouldEqual, 0)
	})
}

func BenchmarkMs15LenIDGenerator_NextString(b *testing.B) {
	idWorker := NewMs15LenIDGenerator(10)
	for i := 0; i < b.N; i++ {
		idWorker.NextID()
	}
}

func BenchmarkMs15LenIDGenerator_NextIdRunParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		idWorker := NewMs15LenIDGenerator(1)
		for pb.Next() {
			idWorker.NextID()
		}
	})
}
