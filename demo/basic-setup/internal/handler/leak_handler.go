package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LeakHandler struct{}

func NewLeakHandler() *LeakHandler {
	return &LeakHandler{}
}

func (h *LeakHandler) GoroutineLeak(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	limit := 100000

	for i := range limit {
		wg.Add(1)
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return c.SendString("context done")
		default:
			go func() {
				defer wg.Done()

				timer := time.NewTimer(time.Duration(10*i) * time.Second)
				defer timer.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-timer.C:
						return
					default:
						fmt.Println("working", i)
					}
				}

			}()
		}
	}

	wg.Wait()

	return c.SendString("leak done")
}

func (h *LeakHandler) Block(c *fiber.Ctx) error {
	var mu sync.Mutex
	var wg sync.WaitGroup
	numGoroutines := 100

	// Shared resource that will cause contention
	sharedData := make(map[int]string)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Heavy mutex contention - all goroutines trying to access the same mutex
			for j := 0; j < 1000; j++ {
				mu.Lock()
				// Simulate some work while holding the lock
				sharedData[id] = fmt.Sprintf("data-%d-%d", id, j)
				time.Sleep(1 * time.Millisecond) // Block while holding lock
				_ = sharedData[id]
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return c.SendString(fmt.Sprintf("block done - processed %d goroutines", numGoroutines))
}

func (h *LeakHandler) Alloc(c *fiber.Ctx) error {
	// Create many allocations: strings, slices, maps, structs
	var result []string
	iterations := 100000

	for i := 0; i < iterations; i++ {
		// Allocate new string
		str := fmt.Sprintf("allocation-%d-data-%d-more-data-%d", i, i*2, i*3)

		// Allocate new slice
		slice := make([]byte, 1024)
		for j := range slice {
			slice[j] = byte(i + j)
		}

		// Allocate new map
		m := make(map[string]int)
		for k := 0; k < 100; k++ {
			m[fmt.Sprintf("key-%d", k)] = k + i
		}

		// Allocate new struct
		type Data struct {
			ID    int
			Name  string
			Value []byte
			Meta  map[string]interface{}
		}
		data := Data{
			ID:    i,
			Name:  str,
			Value: slice,
			Meta:  make(map[string]interface{}),
		}
		for k := 0; k < 50; k++ {
			data.Meta[fmt.Sprintf("meta-%d", k)] = k * i
		}

		result = append(result, data.Name)
	}

	return c.SendString(fmt.Sprintf("alloc done - created %d allocations, result size: %d", iterations, len(result)))
}

func (h *LeakHandler) CPUIntensive(c *fiber.Ctx) error {
	// CPU-intensive task: Calculate prime numbers up to a large number
	limit := 100000
	primes := make([]int, 0, limit/10)

	// Sieve of Eratosthenes - CPU intensive
	sieve := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		sieve[i] = true
	}

	// Mark non-primes
	for i := 2; i*i <= limit; i++ {
		if sieve[i] {
			for j := i * i; j <= limit; j += i {
				sieve[j] = false
			}
		}
	}

	// Collect primes
	for i := 2; i <= limit; i++ {
		if sieve[i] {
			primes = append(primes, i)
		}
	}

	// Additional CPU work: Calculate sum of squares of primes
	sum := 0
	for _, prime := range primes {
		sum += prime * prime
	}

	// More CPU work: Matrix multiplication simulation
	size := 500
	matrixA := make([][]int, size)
	matrixB := make([][]int, size)
	matrixC := make([][]int, size)

	for i := 0; i < size; i++ {
		matrixA[i] = make([]int, size)
		matrixB[i] = make([]int, size)
		matrixC[i] = make([]int, size)
		for j := 0; j < size; j++ {
			matrixA[i][j] = i + j
			matrixB[i][j] = i - j
		}
	}

	// Matrix multiplication - very CPU intensive
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				matrixC[i][j] += matrixA[i][k] * matrixB[k][j]
			}
		}
	}

	return c.SendString(fmt.Sprintf("cpu intensive done - found %d primes, sum of squares: %d, matrix size: %dx%d", len(primes), sum, size, size))
}
