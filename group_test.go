package resultgroup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	err1 = errors.New("Error 1")
	err2 = errors.New("Error 2")
	err3 = errors.New("Error 3")
)

func TestGroup(t *testing.T) {
	t.Parallel()

	t.Run("no errors", testGroupNoErrors)
	t.Run("with errors", testGroupWithErrors)
	t.Run("max errors reached", testGroupMaxErrorsReached)
	t.Run("no error limit", testGroupNoErrorLimit)
}

// testGroupNoErrors checks if the Group works correctly when there are no errors.
func testGroupNoErrors(t *testing.T) {
	group, _ := WithErrorsThreshold[int](context.Background(), 3)

	for i := 0; i < 5; i++ {
		i := i

		group.Go(func() ([]int, error) {
			time.Sleep(100 * time.Millisecond)
			return []int{i}, nil
		})
	}

	results, err := group.Wait()

	assert.Nil(t, err, "Expected no error, got: %v", err)
	assert.Len(t, results, 5, "Expected 5 results, got: %d", len(results))
}

// testGroupWithErrors checks if the Group works correctly when some tasks return errors.
func testGroupWithErrors(t *testing.T) {
	group, _ := WithErrorsThreshold[int](context.Background(), 3)

	group.Go(func() ([]int, error) {
		return nil, err1
	})

	group.Go(func() ([]int, error) {
		return []int{1}, nil
	})

	group.Go(func() ([]int, error) {
		return nil, err2
	})

	group.Go(func() ([]int, error) {
		return []int{2}, nil
	})

	results, err := group.Wait()

	assert.NotNil(t, err, "Expected an error, got nil")
	assert.Len(t, err.Unwrap(), 2, "Expected 2 errors, got: %d", len(err.Unwrap()))
	assert.True(t, errors.Is(err, err1), "Expected error to be: %v, got: %v", err1, err)
	assert.True(t, errors.Is(err, err2), "Expected error to be: %v, got: %v", err2, err)
	assert.Len(t, results, 2, "Expected 2 results, got: %d", len(results))
}

// testGroupMaxErrorsReached checks if the Group cancels the context and stops processing when the error threshold is reached.
func testGroupMaxErrorsReached(t *testing.T) {
	group, ctx := WithErrorsThreshold[int](context.Background(), 2)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	group.Go(func() ([]int, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, err1
	})

	group.Go(func() ([]int, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, err2
	})

	group.Go(func() ([]int, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil, err3
		}
	})

	group.Go(func() ([]int, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return []int{1}, nil
		}
	})

	results, err := group.Wait()

	assert.NotNil(t, err, "Expected an error, got nil")
	assert.Len(t, results, 1, "Expected 1 result, got: %d", len(results))
}

// testGroupNoErrorLimit checks if the Group works correctly without setting an error threshold.
func testGroupNoErrorLimit(t *testing.T) {
	group := Group[int]{}

	group.Go(func() ([]int, error) {
		return nil, err1
	})

	group.Go(func() ([]int, error) {
		return []int{1}, nil
	})

	group.Go(func() ([]int, error) {
		return nil, err2
	})

	group.Go(func() ([]int, error) {
		return []int{2}, nil
	})

	results, err := group.Wait()

	assert.NotNil(t, err, "Expected an error, got nil")
	assert.Len(t, results, 2, "Expected 2 results, got: %d", len(results))
}
