package parallel

type Action func() error

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
