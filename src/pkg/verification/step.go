package verification

type Step struct {
	Name    string   `json:"name"`
	Status  Status   `json:"status"`
	Errors  []string `json:"errors"`
	Remarks []string `json:"remarks"`

	SubSteps []*Step `json:"sub_steps"`
}

func (s *Step) AddStep(name string, status Status, errors ...string) *Step {
	step := Step{
		Name:   name,
		Status: status,
		Errors: errors,
	}
	s.SubSteps = append(s.SubSteps, &step)
	return &step
}

func (s *Step) FailureToWarning() {
	if s.Status == StatusFailure {
		s.Status = StatusWarning
	}
}

func (s *Step) RunStep(name string, fn func() error) *Step {
	step := s.AddStep(name, StatusNotRun)
	err := fn()
	if err != nil {
		step.AddError(err)
		step.Status = StatusFailure
	} else {
		step.Status = StatusSuccess
	}
	return step
}

func (s *Step) AddError(err error) {
	s.Errors = append(s.Errors, err.Error())
}

func (s *Step) DidFail() bool {
	if s.Status == StatusFailure {
		return true
	}

	for _, step := range s.SubSteps {
		if step.DidFail() {
			return true
		}
	}
	return false
}
