package worker

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TasksProcessed tracks the total number of tasks processed
	TasksProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_tasks_processed_total",
			Help: "Total number of tasks processed by the worker",
		},
		[]string{"worker_type", "operation"},
	)

	// TasksInQueue tracks the current number of tasks in the queue
	TasksInQueue = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_tasks_in_queue",
			Help: "Current number of tasks in the worker queue",
		},
		[]string{"worker_type"},
	)

	// TaskProcessingDuration tracks task processing duration
	TaskProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "worker_task_processing_duration_seconds",
			Help:    "Duration of task processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"worker_type", "operation"},
	)

	// ActiveGoroutines tracks the number of active goroutines
	ActiveGoroutines = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_active_goroutines",
			Help: "Number of active goroutines in the worker",
		},
		[]string{"worker_type"},
	)

	// TaskErrors tracks the number of task processing errors
	TaskErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_task_errors_total",
			Help: "Total number of task processing errors",
		},
		[]string{"worker_type", "error_type"},
	)

	// AllocatedMemory tracks memory allocation
	AllocatedMemory = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_allocated_memory_bytes",
			Help: "Current allocated memory in bytes",
		},
		[]string{"worker_type"},
	)

	// TotalAllocations tracks total allocations
	TotalAllocations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_total_allocations_bytes",
			Help: "Total memory allocated in bytes",
		},
		[]string{"worker_type"},
	)

	// GCRuns tracks garbage collection runs
	GCRuns = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_gc_runs_total",
			Help: "Total number of garbage collection runs",
		},
		[]string{"worker_type"},
	)
)
