package parallel

// Action is the function to be consumed by ForEach. It should return an error if the action failed
// and should not be retried
type Action func() error

// This is similar to sync.ErrorGroup, but returns all errors and does not cancel.
type ErrorGroup []Action

func (eg ErrorGroup) Errors() []error {
	errChan := make(chan error, len(eg))

	for _, a := range eg {
		a := a
		go func() {
			errChan <- a()
		}()
	}

	errs := make([]error, 0)
	for range eg {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
