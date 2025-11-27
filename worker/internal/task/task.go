package task

import (
	"fmt"
	"math"
)

// OperationType represents the type of operation to perform
type OperationType int

const (
	OperationSum OperationType = iota
	OperationMultiply
	OperationSquare
	OperationFilter
)

// Task represents a work item to be processed
type Task struct {
	ID        int
	Data      []int
	Operation OperationType
}

// Result represents the result of processing a task
type Result struct {
	TaskID int
	Value  interface{}
	Info   string
}

// Process performs the operation on the task data
func (t *Task) Process() *Result {
	var value interface{}
	var info string

	switch t.Operation {
	case OperationSum:
		sum := 0
		for _, v := range t.Data {
			sum += v
		}
		value = sum
		info = fmt.Sprintf("Sum of %d numbers", len(t.Data))

	case OperationMultiply:
		product := 1
		for _, v := range t.Data {
			product *= v
		}
		value = product
		info = fmt.Sprintf("Product of %d numbers", len(t.Data))

	case OperationSquare:
		squares := make([]int, len(t.Data))
		for i, v := range t.Data {
			squares[i] = v * v
		}
		value = squares
		info = fmt.Sprintf("Squared %d numbers", len(t.Data))

	case OperationFilter:
		filtered := []int{}
		for _, v := range t.Data {
			if v%2 == 0 {
				filtered = append(filtered, v)
			}
		}
		value = filtered
		info = fmt.Sprintf("Filtered %d even numbers from %d", len(filtered), len(t.Data))

	default:
		value = nil
		info = "Unknown operation"
	}

	return &Result{
		TaskID: t.ID,
		Value:  value,
		Info:   info,
	}
}

// NewTask creates a new task with random data
func NewTask(id int, dataSize int, op OperationType) *Task {
	data := make([]int, dataSize)
	for i := range data {
		data[i] = (i*7 + id*3) % 100 // Simple pseudo-random generation
	}
	return &Task{
		ID:        id,
		Data:      data,
		Operation: op,
	}
}

// CalculateComplexity performs some complex calculations (for profiling)
func (t *Task) CalculateComplexity() float64 {
	result := 0.0
	for _, v := range t.Data {
		result += math.Sqrt(float64(v)) * math.Log(float64(v+1))
	}
	return result
}
