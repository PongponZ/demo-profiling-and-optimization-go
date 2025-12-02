package benchmark_test

import (
	"sync"
	"testing"
)

// BenchmarkMutexReadHeavy benchmarks sync.Mutex with read-heavy workload
// Mutex provides exclusive locking - only one goroutine can access the resource at a time,
// even for read operations. This makes it slower in read-heavy scenarios because
// all read operations must wait for each other sequentially.
func BenchmarkMutexReadHeavy(b *testing.B) {
	var mu sync.Mutex
	var data int
	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		defer wg.Done()
		for pb.Next() {
			mu.Lock()
			_ = data // read operation
			mu.Unlock()
		}
	})
	wg.Wait()
}

// BenchmarkRWMutexReadHeavy benchmarks sync.RWMutex with read-heavy workload
// RWMutex allows multiple readers to access the resource concurrently using RLock(),
// while only writers need exclusive access using Lock(). In read-heavy scenarios,
// this parallel read capability makes RWMutex significantly faster than Mutex
// because multiple goroutines can read simultaneously without blocking each other.
func BenchmarkRWMutexReadHeavy(b *testing.B) {
	var mu sync.RWMutex
	var data int
	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		defer wg.Done()
		for pb.Next() {
			mu.RLock()
			_ = data // read operation
			mu.RUnlock()
		}
	})
	wg.Wait()
}

// BenchmarkMutexWriteHeavy benchmarks sync.Mutex with write-heavy workload
// In write-heavy scenarios, both Mutex and RWMutex need exclusive locks for writes.
// Mutex might be slightly faster here because it has less overhead (no need to track
// reader count), but the difference is usually minimal since writes require exclusive access anyway.
func BenchmarkMutexWriteHeavy(b *testing.B) {
	var mu sync.Mutex
	var data int
	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		defer wg.Done()
		for pb.Next() {
			mu.Lock()
			data++ // write operation
			mu.Unlock()
		}
	})
	wg.Wait()
}

// BenchmarkRWMutexWriteHeavy benchmarks sync.RWMutex with write-heavy workload
// In write-heavy scenarios, RWMutex's Lock() behaves similarly to Mutex's Lock(),
// requiring exclusive access. However, RWMutex has slightly more overhead due to
// reader count tracking, making it marginally slower than Mutex for write-only workloads.
func BenchmarkRWMutexWriteHeavy(b *testing.B) {
	var mu sync.RWMutex
	var data int
	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		defer wg.Done()
		for pb.Next() {
			mu.Lock()
			data++ // write operation
			mu.Unlock()
		}
	})
	wg.Wait()
}

// BenchmarkMutexMixed benchmarks sync.Mutex with mixed read/write workload
// Even with mixed workloads, Mutex forces all operations (reads and writes) to be sequential,
// which can create contention when there are many concurrent readers that could run in parallel.
func BenchmarkMutexMixed(b *testing.B) {
	var mu sync.Mutex
	var data int
	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		defer wg.Done()
		for pb.Next() {
			mu.Lock()
			if data%2 == 0 {
				_ = data // read
			} else {
				data++ // write
			}
			mu.Unlock()
		}
	})
	wg.Wait()
}

// BenchmarkRWMutexMixed benchmarks sync.RWMutex with mixed read/write workload
// RWMutex excels in mixed workloads because reads can happen concurrently (RLock),
// while writes still get exclusive access (Lock). This allows better parallelism
// when there are multiple readers, making it faster than Mutex in most real-world scenarios
// where reads typically outnumber writes.
func BenchmarkRWMutexMixed(b *testing.B) {
	var mu sync.RWMutex
	var data int
	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		defer wg.Done()
		for pb.Next() {
			if data%2 == 0 {
				mu.RLock()
				_ = data // read
				mu.RUnlock()
			} else {
				mu.Lock()
				data++ // write
				mu.Unlock()
			}
		}
	})
	wg.Wait()
}
