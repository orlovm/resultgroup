package resultgroup

import (
	"context"
	"sync"
)

// Group is a generic struct that holds errors and results from concurrent tasks.
// To create a Group without a context and error threshold, use the struct directly:
// group := resultgroup.Group[T]{}
type Group[T any] struct {
	mutex     sync.Mutex
	errs      []error
	wg        sync.WaitGroup
	cancel    func()
	threshold int
	results   []T
}

// WithErrorsThreshold creates a new Group with the provided context
// and a threshold for the maximum number of errors.
// If the threshold is reached, the context will be canceled.
// Threshold must be greater than or equal to 1.
func WithErrorsThreshold[T any](ctx context.Context, threshold int) (Group[T], context.Context) {
	if threshold < 1 {
		panic("threshold must be greater than or equal to 1")
	}

	ctx, cancel := context.WithCancel(ctx)

	return Group[T]{cancel: cancel, threshold: threshold}, ctx
}

// Go runs the provided function in a new goroutine and append the results
// to aggregated slice that will be returned by Wait.
// If the function returns an error, it will be appended to the aggregated
// slice of errors if the threshold is not reached.
func (g *Group[T]) Go(f func() ([]T, error)) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		res, err := f()
		g.processResult(res, err)
	}()
}

func (g *Group[T]) processResult(res []T, err error) {
	if err != nil {
		g.handleErrors(err)
	}

	g.appendResults(res)
}

func (g *Group[T]) handleErrors(err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.threshold == 0 || len(g.errs) < g.threshold {
		g.errs = append(g.errs, err)
	}

	if len(g.errs) == g.threshold {
		if g.cancel != nil {
			g.cancel()
		}
	}
}

func (g *Group[T]) appendResults(res []T) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.results = append(g.results, res...)
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the concatenated results and a multiError containing all errors that
// are below the threshold.
func (g *Group[T]) Wait() ([]T, errorWithUnwrap) {
	g.wg.Wait()
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.cancel != nil {
		g.cancel()
	}

	if len(g.errs) == 0 {
		return g.results, nil
	}

	return g.results, &multiError{errs: g.errs}
}
