package benchmark_test

import (
	"testing"
)

const (
	sliceSize = 1000
	mapSize   = 1000
)

// BenchmarkSliceNoCapacity benchmarks slice append without pre-allocated capacity
// When you append to a slice without pre-allocated capacity, Go needs to:
// 1. Check if there's enough capacity
// 2. If not, allocate a new underlying array (usually 2x the current size)
// 3. Copy all existing elements to the new array
// 4. Add the new element
// This results in multiple allocations and copies as the slice grows, making it slower.
func BenchmarkSliceNoCapacity(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var slice []int
		for j := 0; j < sliceSize; j++ {
			slice = append(slice, j) // May trigger reallocation multiple times
		}
		_ = slice
	}
}

// BenchmarkSliceWithCapacity benchmarks slice append with pre-allocated capacity
// Pre-allocating capacity using make([]int, 0, capacity) tells Go to allocate
// the underlying array with the specified capacity upfront. This means:
// 1. No reallocations needed (as long as you don't exceed capacity)
// 2. No copying of existing elements
// 3. Much faster append operations
// This is the recommended approach when you know the approximate size beforehand.
func BenchmarkSliceWithCapacity(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := make([]int, 0, sliceSize) // Pre-allocate capacity
		for j := 0; j < sliceSize; j++ {
			slice = append(slice, j) // No reallocation needed
		}
		_ = slice
	}
}

// BenchmarkMapNoCapacity benchmarks map operations without pre-allocated capacity
// Maps in Go are implemented as hash tables. When you add elements to a map without
// pre-allocated capacity, Go needs to:
// 1. Check if the map needs to grow (load factor threshold)
// 2. Allocate a new, larger bucket array
// 3. Rehash all existing key-value pairs into the new buckets
// This rehashing process can be expensive, especially as the map grows.
func BenchmarkMapNoCapacity(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[int]int) // No capacity hint
		for j := 0; j < mapSize; j++ {
			m[j] = j // May trigger rehashing multiple times
		}
		_ = m
	}
}

// BenchmarkMapWithCapacity benchmarks map operations with pre-allocated capacity
// Pre-allocating capacity using make(map[int]int, capacity) tells Go to allocate
// the initial bucket array with enough space for the expected number of elements.
// This reduces the number of rehashing operations needed as elements are added.
// Note: The actual capacity may be rounded up to the nearest power of 2 for efficiency.
func BenchmarkMapWithCapacity(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[int]int, mapSize) // Pre-allocate capacity
		for j := 0; j < mapSize; j++ {
			m[j] = j // Fewer or no rehashing operations
		}
		_ = m
	}
}
