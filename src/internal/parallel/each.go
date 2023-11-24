package parallel

// Action is the function to be consumed by ForEach. It should return an error if the action failed
// and should not be retried
type Action func() error

// ForEach runs the given actions in parallel
// If an error is returned from an action, it is added to the slice of errors and the action is not retried
// any actions that are still running will be allowed to complete
func ForEach(actions []Action, maxConcurrency int) []error {
	// Populate tokens
	tokens := make(chan int, maxConcurrency)
	for i := 0; i < maxConcurrency; i++ {
		tokens <- i
	}

	errChan := make(chan error, len(actions))
	for _, a := range actions {
		a := a
		token := <-tokens
		go func() {
			defer func() { tokens <- token }()
			errChan <- a()
		}()
	}

	var errs []error
	for range actions {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
