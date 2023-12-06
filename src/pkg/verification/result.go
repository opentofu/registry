package verification

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusNotRun  Status = "not_run"
	StatusSkipped Status = "skipped"
)

type Result struct {
	Steps []*Step `json:"steps"`
}

func (r *Result) AddStep(name string, status Status, errors ...string) *Step {
	step := Step{
		Name:   name,
		Status: status,
		Errors: errors,
	}
	r.Steps = append(r.Steps, &step)
	return &step
}

func (r *Result) DidFail() bool {
	for _, step := range r.Steps {
		if step.DidFail() {
			return true
		}
	}
	return false
}
