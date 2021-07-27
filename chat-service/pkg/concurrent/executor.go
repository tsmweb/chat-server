package concurrent

import (
	"context"
)

// ExecutorService contains logic for goroutine execution.
type ExecutorService interface {

	// Schedule schedules the job to be performed on the executors.
	Schedule(job func(ctx context.Context)) error

	// Shutdown closes the executor.
	Shutdown()
}