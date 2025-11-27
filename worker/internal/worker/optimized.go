package worker

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
	"worker/internal/task"
)

// OptimizedWorker demonstrates best practices:
// 1. Proper goroutine management with context cancellation
// 2. Pre-allocated slice/map capacity
// 3. Reduced allocations with buffer reuse and object pools
type OptimizedWorker struct {
	taskQueue chan *task.Task
	results   []*task.Result
	mu        sync.Mutex
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc

	// OPTIMIZATION 3: Reusable buffers to reduce allocations
	stringBuilderPool sync.Pool
	resultBuffer      []*task.Result // Pre-allocated buffer
}

// NewOptimizedWorker creates a new optimized worker
func NewOptimizedWorker(expectedTasks int) *OptimizedWorker {
	ctx, cancel := context.WithCancel(context.Background())
	
	// OPTIMIZATION 2: Pre-allocate slice capacity
	results := make([]*task.Result, 0, expectedTasks)
	
	w := &OptimizedWorker{
		taskQueue:    make(chan *task.Task, expectedTasks), // Buffered channel
		results:      results,
		ctx:          ctx,
		cancel:       cancel,
		resultBuffer: make([]*task.Result, 0, 100), // Pre-allocated buffer
	}
	
	// OPTIMIZATION 3: Initialize object pool for string builders
	w.stringBuilderPool.New = func() interface{} {
		return &strings.Builder{}
	}
	
	return w
}

// Start begins processing tasks with proper context management
func (w *OptimizedWorker) Start(numWorkers int) {
	// OPTIMIZATION 1: All goroutines respect context cancellation
	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}
	
	// OPTIMIZATION 1: Monitor goroutine can be cancelled
	w.wg.Add(1)
	go w.monitor()
	
	// Start metrics collection
	w.wg.Add(1)
	go w.collectMetrics()
}

// worker processes tasks with reduced allocations
func (w *OptimizedWorker) worker(id int) {
	defer w.wg.Done()

	// OPTIMIZATION 3: Reuse string builder from pool
	sb := w.stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		w.stringBuilderPool.Put(sb)
	}()

	// OPTIMIZATION 3: Pre-allocate worker name (do once, not per task)
	workerName := fmt.Sprintf("Worker-%d", id)

	// OPTIMIZATION 3: Reuse log entry struct
	logEntry := &LogEntry{
		WorkerID: id,
		Time:     time.Now(),
	}

	// OPTIMIZATION 2: Pre-size metadata map
	metadata := make(map[string]interface{}, 3) // Pre-sized with expected capacity

	for {
		select {
		case <-w.ctx.Done():
			return // OPTIMIZATION 1: Proper exit on cancellation
		case t, ok := <-w.taskQueue:
			if !ok {
				return
			}

			// OPTIMIZATION 3: Reuse struct fields instead of allocating new struct
			logEntry.TaskID = t.ID
			logEntry.Time = time.Now()

			// OPTIMIZATION 3: Use string builder instead of concatenation
			sb.Reset()
			sb.WriteString(workerName)
			sb.WriteString(" processing task")
			_ = sb.String()
			
			// Process task
			start := time.Now()
			result := t.Process()
			duration := time.Since(start).Seconds()
			
			// Track metrics
			TasksProcessed.WithLabelValues("optimized", fmt.Sprintf("%d", t.Operation)).Inc()
			TaskProcessingDuration.WithLabelValues("optimized", fmt.Sprintf("%d", t.Operation)).Observe(duration)

			// OPTIMIZATION 2: Append to pre-allocated slice
			w.mu.Lock()
			w.results = append(w.results, result) // Capacity already allocated
			w.mu.Unlock()

			// OPTIMIZATION 2: Reuse map, clear and refill
			for k := range metadata {
				delete(metadata, k)
			}
			metadata["worker"] = id
			metadata["task"] = t.ID
			metadata["processed"] = true

			// OPTIMIZATION 3: Use string builder for efficient string building
			sb.Reset()
			for k, v := range metadata {
				sb.WriteString(k)
				sb.WriteString(":")
				sb.WriteString(fmt.Sprintf("%v", v))
				sb.WriteString(" ")
			}
			_ = sb.String()

			// OPTIMIZATION 1: Helper task with context cancellation
			w.wg.Add(1)
			go w.helperTask(t.ID)
		}
	}
}

// helperTask demonstrates proper goroutine lifecycle management
func (w *OptimizedWorker) helperTask(taskID int) {
	defer w.wg.Done()

	// OPTIMIZATION 1: Use context for cancellation
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// OPTIMIZATION 2: Pre-allocate slice with capacity
	checkData := make([]int, 0, 10) // Pre-allocated capacity

	// OPTIMIZATION 3: Reuse status struct
	status := &Status{
		TaskID: taskID,
	}

	for {
		select {
		case <-w.ctx.Done():
			return // OPTIMIZATION 1: Proper exit
		case <-ticker.C:
			// OPTIMIZATION 2: Reuse slice, reset length
			checkData = checkData[:0] // Reset without reallocating
			for i := 0; i < 10; i++ {
				checkData = append(checkData, i) // No reallocation needed
			}
			_ = checkData

			// OPTIMIZATION 3: Reuse struct, just update fields
			status.Checked = time.Now()
			_ = status
		}
	}
}

// monitor demonstrates proper goroutine lifecycle
func (w *OptimizedWorker) monitor() {
	defer w.wg.Done()

	// OPTIMIZATION 2: Pre-size map with expected capacity
	stats := make(map[string]int, 2) // Pre-sized

	// OPTIMIZATION 3: Reuse string builder
	sb := w.stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		w.stringBuilderPool.Put(sb)
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return // OPTIMIZATION 1: Proper exit
		case <-ticker.C:
			w.mu.Lock()
			count := len(w.results)
			w.mu.Unlock()

			// OPTIMIZATION 2: Reuse map, just update values
			stats["total"] = count
			stats["timestamp"] = int(time.Now().Unix())

			// OPTIMIZATION 3: Use string builder efficiently
			sb.Reset()
			for k, v := range stats {
				sb.WriteString(k)
				sb.WriteString("=")
				sb.WriteString(fmt.Sprintf("%d", v))
				sb.WriteString(" ")
			}
			_ = sb.String()
		}
	}
}

// ProcessTasks adds tasks to the queue
func (w *OptimizedWorker) ProcessTasks(tasks []*task.Task) {
	// OPTIMIZATION 1: Use buffered channel to avoid blocking
	// OPTIMIZATION 1: Check context before sending
	for _, t := range tasks {
		select {
		case <-w.ctx.Done():
			return
		case w.taskQueue <- t:
			// Task queued successfully
		}
	}
}

// GetResults returns processed results efficiently
func (w *OptimizedWorker) GetResults() []*task.Result {
	w.mu.Lock()
	defer w.mu.Unlock()

	// OPTIMIZATION 2: Pre-allocate result copy with known capacity
	resultCopy := make([]*task.Result, len(w.results))
	copy(resultCopy, w.results) // More efficient than append loop
	return resultCopy
}

// Stop properly stops all goroutines
func (w *OptimizedWorker) Stop() {
	// OPTIMIZATION 1: Cancel context to signal all goroutines to stop
	w.cancel()
	close(w.taskQueue)
	w.wg.Wait() // Wait for all goroutines to finish
}

// collectMetrics periodically collects and reports metrics
func (w *OptimizedWorker) collectMetrics() {
	defer w.wg.Done()
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	var lastNumGC uint32
	var lastTotalAlloc uint64
	
	for {
		select {
		case <-w.ctx.Done():
			return // OPTIMIZATION 1: Proper exit
		case <-ticker.C:
			// Collect runtime metrics
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			
			// Update metrics
			ActiveGoroutines.WithLabelValues("optimized").Set(float64(runtime.NumGoroutine()))
			AllocatedMemory.WithLabelValues("optimized").Set(float64(m.Alloc))
			
			// Track incremental allocations
			if m.TotalAlloc > lastTotalAlloc {
				TotalAllocations.WithLabelValues("optimized").Add(float64(m.TotalAlloc - lastTotalAlloc))
				lastTotalAlloc = m.TotalAlloc
			}
			
			// Track GC runs
			if m.NumGC > lastNumGC {
				GCRuns.WithLabelValues("optimized").Add(float64(m.NumGC - lastNumGC))
				lastNumGC = m.NumGC
			}
			
			// Track queue size
			w.mu.Lock()
			queueSize := len(w.taskQueue)
			w.mu.Unlock()
			TasksInQueue.WithLabelValues("optimized").Set(float64(queueSize))
		}
	}
}
