package worker

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"worker/internal/task"
)

// BadWorker demonstrates problematic patterns:
// 1. Goroutine leak - goroutines that never finish
// 2. Slice/Map capacity issues - no pre-allocation
// 3. Excessive allocations - creating objects in hot paths
type BadWorker struct {
	taskQueue chan *task.Task
	results   []*task.Result
	mu        sync.Mutex
	wg        sync.WaitGroup
}

// NewBadWorker creates a new problematic worker
func NewBadWorker() *BadWorker {
	return &BadWorker{
		taskQueue: make(chan *task.Task), // Unbuffered channel - can cause blocking
		results:   []*task.Result{},      // No capacity pre-allocation
	}
}

// Start begins processing tasks (with goroutine leaks)
func (w *BadWorker) Start(numWorkers int) {
	// PROBLEM 1: Goroutine Leak
	// Spawn workers without proper cleanup mechanism
	// These goroutines will run forever even after tasks are done
	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go w.worker(i) // No context cancellation, no way to stop gracefully
	}

	// Another goroutine leak: background monitor that never exits
	go w.monitor() // This goroutine runs forever
	
	// Start metrics collection
	go w.collectMetrics()
}

// worker processes tasks (demonstrates allocation problems)
func (w *BadWorker) worker(id int) {
	defer w.wg.Done()

	// PROBLEM 3: Allocation in hot path
	// Creating new string for each task processed
	workerName := fmt.Sprintf("Worker-%d", id) // String allocation

	for t := range w.taskQueue {
		// PROBLEM 3: Allocating new struct in hot path
		logEntry := &LogEntry{
			WorkerID: id,
			TaskID:   t.ID,
			Time:     time.Now(),
			Message:  workerName + " processing task", // String concatenation
		}

		// Process task
		start := time.Now()
		result := t.Process()
		duration := time.Since(start).Seconds()
		
		// Track metrics
		TasksProcessed.WithLabelValues("bad", fmt.Sprintf("%d", t.Operation)).Inc()
		TaskProcessingDuration.WithLabelValues("bad", fmt.Sprintf("%d", t.Operation)).Observe(duration)
		
		// PROBLEM 2: Slice capacity issue - appending without pre-allocation
		w.mu.Lock()
		w.results = append(w.results, result) // No capacity hint
		w.mu.Unlock()

		// PROBLEM 3: Unnecessary allocation - creating map for each task
		metadata := make(map[string]interface{}) // No pre-sizing
		metadata["worker"] = id
		metadata["task"] = t.ID
		metadata["processed"] = true

		// PROBLEM 3: String concatenation in loop
		logMessage := ""
		for k, v := range metadata {
			logMessage += fmt.Sprintf("%s:%v ", k, v) // Inefficient string concatenation
		}
		_ = logMessage // Use it to avoid compiler optimization
		_ = logEntry   // Use it to avoid compiler optimization

		// PROBLEM 1: Goroutine leak - spawn helper goroutine that may never finish
		go w.helperTask(t.ID) // No way to cancel this
	}
}

// helperTask demonstrates goroutine leak - runs indefinitely
func (w *BadWorker) helperTask(taskID int) {
	// This goroutine runs forever, checking something periodically
	// No context, no cancellation mechanism
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// PROBLEM 1: This loop runs forever with no way to exit
	// No context, no cancellation, no timeout
	for {
		<-ticker.C // Wait for ticker (blocks forever if ticker never fires)

		// PROBLEM 2: Creating new slice without capacity in loop
		checkData := []int{} // No capacity
		for i := 0; i < 10; i++ {
			checkData = append(checkData, i) // Multiple reallocations
		}
		_ = checkData

		// PROBLEM 3: Allocating new struct each tick
		status := &Status{
			TaskID:  taskID,
			Checked: time.Now(),
		}
		_ = status
	}
}

// monitor demonstrates another goroutine leak
func (w *BadWorker) monitor() {
	// This goroutine never exits
	// PROBLEM 2: Creating map without pre-sizing
	stats := make(map[string]int) // No capacity hint

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.mu.Lock()
			count := len(w.results)
			w.mu.Unlock()

			// PROBLEM 2: Adding to map without pre-sizing
			stats["total"] = count
			stats["timestamp"] = int(time.Now().Unix())

			// PROBLEM 3: String concatenation
			report := ""
			for k, v := range stats {
				report += fmt.Sprintf("%s=%d ", k, v) // Inefficient
			}
			_ = report
		}
	}
}

// ProcessTasks adds tasks to the queue
func (w *BadWorker) ProcessTasks(tasks []*task.Task) {
	// PROBLEM 1: Sending to unbuffered channel can block
	// If workers are slow, this will block the caller
	for _, t := range tasks {
		w.taskQueue <- t // Can block indefinitely
	}
}

// GetResults returns processed results
func (w *BadWorker) GetResults() []*task.Result {
	w.mu.Lock()
	defer w.mu.Unlock()

	// PROBLEM 2: Creating new slice without capacity
	resultCopy := []*task.Result{}
	for _, r := range w.results {
		resultCopy = append(resultCopy, r) // No capacity pre-allocation
	}
	return resultCopy
}

// Stop attempts to stop the worker (but doesn't properly clean up goroutines)
func (w *BadWorker) Stop() {
	close(w.taskQueue)
	w.wg.Wait() // Wait for workers, but monitor() and helperTask() goroutines still leak
}

// collectMetrics periodically collects and reports metrics
func (w *BadWorker) collectMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	var lastNumGC uint32
	var lastTotalAlloc uint64
	
	for {
		<-ticker.C
		
		// Collect runtime metrics
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// Update metrics
		ActiveGoroutines.WithLabelValues("bad").Set(float64(runtime.NumGoroutine()))
		AllocatedMemory.WithLabelValues("bad").Set(float64(m.Alloc))
		
		// Track incremental allocations
		if m.TotalAlloc > lastTotalAlloc {
			TotalAllocations.WithLabelValues("bad").Add(float64(m.TotalAlloc - lastTotalAlloc))
			lastTotalAlloc = m.TotalAlloc
		}
		
		// Track GC runs
		if m.NumGC > lastNumGC {
			GCRuns.WithLabelValues("bad").Add(float64(m.NumGC - lastNumGC))
			lastNumGC = m.NumGC
		}
		
		// Track queue size
		w.mu.Lock()
		queueSize := len(w.taskQueue)
		w.mu.Unlock()
		TasksInQueue.WithLabelValues("bad").Set(float64(queueSize))
	}
}

// Helper types for demonstration
type LogEntry struct {
	WorkerID int
	TaskID   int
	Time     time.Time
	Message  string
}

type Status struct {
	TaskID  int
	Checked time.Time
}
