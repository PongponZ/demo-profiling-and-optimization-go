package benchmark_test

import (
	"regexp"
	"testing"
)

var (
	// Pre-compiled regex pattern for email validation
	// This pattern matches most common email formats
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// BenchmarkRegexOnTheFly benchmarks regex compilation and matching on-the-fly
// This approach compiles the regex pattern every time it's used, which is inefficient
// because regex compilation is expensive. This is useful when the pattern changes
// frequently, but for static patterns like email validation, it's wasteful.
func BenchmarkRegexOnTheFly(b *testing.B) {
	testEmails := []string{
		"user@example.com",
		"test.email@domain.co.uk",
		"invalid.email",
		"another@test.org",
		"not-an-email",
		"valid123@subdomain.example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := testEmails[i%len(testEmails)]
		// Compile regex on every iteration (inefficient)
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
		_ = matched
	}
}

// BenchmarkRegexPreCompile benchmarks pre-compiled regex matching
// This approach compiles the regex pattern once (at package initialization) and reuses it.
// This is much more efficient because regex compilation is expensive, and for static patterns
// like email validation, we only need to compile once. The compiled regex can be reused
// across multiple calls, making it significantly faster than compiling on every use.
func BenchmarkRegexPreCompile(b *testing.B) {
	testEmails := []string{
		"user@example.com",
		"test.email@domain.co.uk",
		"invalid.email",
		"another@test.org",
		"not-an-email",
		"valid123@subdomain.example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := testEmails[i%len(testEmails)]
		// Use pre-compiled regex (efficient)
		matched := emailRegex.MatchString(email)
		_ = matched
	}
}

// BenchmarkRegexPreCompileParallel benchmarks pre-compiled regex with parallel execution
// This demonstrates that pre-compiled regexes are safe for concurrent use and can
// benefit from parallel execution, making them even more efficient in concurrent scenarios.
func BenchmarkRegexPreCompileParallel(b *testing.B) {
	testEmails := []string{
		"user@example.com",
		"test.email@domain.co.uk",
		"invalid.email",
		"another@test.org",
		"not-an-email",
		"valid123@subdomain.example.com",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			// Use local counter to cycle through test emails in parallel
			email := testEmails[counter%len(testEmails)]
			counter++
			// Use pre-compiled regex in parallel (thread-safe and efficient)
			matched := emailRegex.MatchString(email)
			_ = matched
		}
	})
}

