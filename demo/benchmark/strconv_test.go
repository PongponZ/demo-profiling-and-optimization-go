package benchmark_test

import (
	"fmt"
	"strconv"
	"testing"
)

// BenchmarkStrconvItoa benchmarks strconv.Itoa for converting int to string
func BenchmarkStrconvItoa(b *testing.B) {
	value := 12345
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(value)
	}
}

// BenchmarkFmtSprintf benchmarks fmt.Sprintf for converting int to string
func BenchmarkFmtSprintf(b *testing.B) {
	value := 12345
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%d", value)
	}
}
