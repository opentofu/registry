package parallel

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestForEach(t *testing.T) {
	// Note: these actions are not thread-safe, but that's fine for this test
	// The idea of this test is to ensure that the concurrency is working as expected
	// and that the errors are returned as expected

	successAction := func() error {
		time.Sleep(1 * time.Millisecond) // simulate some work
		return nil
	}

	errorAction := func() error {
		time.Sleep(1 * time.Millisecond) // simulate some work
		return errors.New("error occurred")
	}

	tests := []struct {
		name           string
		actions        []Action
		maxConcurrency int
		expectedErrs   int
	}{
		{
			name:           "All successful actions",
			actions:        []Action{successAction, successAction, successAction},
			maxConcurrency: 2,
			expectedErrs:   0,
		},
		{
			name:           "Some actions with errors",
			actions:        []Action{errorAction, successAction, errorAction},
			maxConcurrency: 2,
			expectedErrs:   2,
		},
		{
			name:           "Some actions with errors, more actions than concurrency",
			actions:        []Action{errorAction, successAction, errorAction, errorAction, errorAction, errorAction},
			maxConcurrency: 2,
			expectedErrs:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ForEach(tt.actions, tt.maxConcurrency)
			assert.Len(t, errs, tt.expectedErrs, "The number of errors should match expected")
		})
	}
}
