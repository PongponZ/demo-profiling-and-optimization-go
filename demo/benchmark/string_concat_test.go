package benchmark_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var (
	testStrings = []string{
		"Hello",
		"World",
		"Go",
		"Programming",
		"Language",
		"Benchmark",
		"Testing",
		"Performance",
		"Hello",
		"World",
		"Go",
		"Programming",
		"Language",
		"Benchmark",
		"Testing",
		"Performance",
		"Hello",
		"World",
		"Go",
		"Programming",
		"Language",
		"Benchmark",
		"Testing",
		"Performance",
		"Hello",
		"World",
		"Go",
		"Programming",
		"Language",
		"Benchmark",
		"Testing",
		"Performance",
		"Hello",
		"World",
		"Go",
		"Programming",
		"Language",
		"Benchmark",
		"Testing",
		"Performance",
		"Hello",
		"World",
		"Go",
		"Programming",
		"Language",
		"Benchmark",
		"Testing",
		"Performance",
	}
)

// BenchmarkStringConcatPlus benchmarks string concatenation using the + operator
// This approach creates a new string for each concatenation, which means it allocates
// new memory and copies all previous characters. For many concatenations, this results
// in O(n²) time complexity and excessive memory allocations, making it inefficient.
func BenchmarkStringConcatPlus(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := ""
		for _, s := range testStrings {
			result += s + " " // Inefficient: creates new string each time
		}
		_ = result
	}
}

// BenchmarkStringConcatFmtSprintf benchmarks string concatenation using fmt.Sprintf
// Similar to the + operator, fmt.Sprintf creates new strings for each concatenation.
// It's convenient but has overhead from parsing the format string and is generally
// slower than direct concatenation methods for simple string building.
func BenchmarkStringConcatFmtSprintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := ""
		for _, s := range testStrings {
			result = fmt.Sprintf("%s %s", result, s)
		}
		_ = result
	}
}

// BenchmarkStringConcatStringsBuilder benchmarks string concatenation using strings.Builder
// strings.Builder is the recommended approach for building strings incrementally.
// It uses an internal buffer that grows efficiently, minimizing allocations.
// It's much faster than + operator for multiple concatenations because it only
// allocates memory when the buffer needs to grow, not on every concatenation.
func BenchmarkStringConcatStringsBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		for _, s := range testStrings {
			builder.WriteString(s)
			builder.WriteString(" ")
		}
		result := builder.String()
		_ = result
	}
}

// BenchmarkStringConcatBytesBuffer benchmarks string concatenation using bytes.Buffer
// bytes.Buffer works similarly to strings.Builder but is more general-purpose.
// It can handle both string and byte operations. For string-only concatenation,
// strings.Builder is slightly more efficient, but bytes.Buffer is still much
// better than the + operator for multiple concatenations.
func BenchmarkStringConcatBytesBuffer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		for _, s := range testStrings {
			buffer.WriteString(s)
			buffer.WriteString(" ")
		}
		result := buffer.String()
		_ = result
	}
}

// BenchmarkStringConcatStringsBuilderPreAllocated benchmarks strings.Builder with pre-allocated capacity
// Pre-allocating capacity can improve performance by reducing the number of buffer
// reallocations. If you know approximately how large the final string will be,
// pre-allocating can make strings.Builder even more efficient.
func BenchmarkStringConcatStringsBuilderPreAllocated(b *testing.B) {
	// Estimate total capacity needed (sum of all string lengths + spaces)
	estimatedCapacity := 0
	for _, s := range testStrings {
		estimatedCapacity += len(s) + 1 // +1 for space
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		builder.Grow(estimatedCapacity) // Pre-allocate capacity
		for _, s := range testStrings {
			builder.WriteString(s)
			builder.WriteString(" ")
		}
		result := builder.String()
		_ = result
	}
}
