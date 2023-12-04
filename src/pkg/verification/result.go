package verification

import "fmt"

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusNotRun  Status = "not_run"
	StatusSkipped Status = "skipped"
)

type Step struct {
	Name   string   `json:"name"`
	Status Status   `json:"status"`
	Errors []string `json:"errors"`

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

func (r *Result) DidFail() bool {
	for _, step := range r.Steps {
		if step.DidFail() {
			return true
		}
	}
	return false
}

func (r *Result) RenderMarkdown() string {
	var output string
	for _, step := range r.Steps {
		output += fmt.Sprintf("## %s\n", step.Name)
		if step.Status == StatusSuccess {
			output += "✅ **Success**\n"
		} else if step.Status == StatusFailure {
			output += "❌ **Failure**\n"
		} else if step.Status == StatusNotRun {
			output += "⚠️ **Not Run**\n"
		} else if step.Status == StatusSkipped {
			output += "⚠️ **Skipped**\n"
		}
		for _, err := range step.Errors {
			output += fmt.Sprintf("- %s\n", err)
		}
		for _, subStep := range step.SubSteps {
			output += fmt.Sprintf("### %s\n", subStep.Name)
			if subStep.Status == StatusSuccess {
				output += "✅ **Success**\n"
			} else if subStep.Status == StatusFailure {
				output += "❌ **Failure**\n"
			} else if subStep.Status == StatusNotRun {
				output += "⚠️ **Not Run**\n"
			} else if subStep.Status == StatusSkipped {
				output += "⚠️ **Skipped**\n"
			}
			for _, err := range subStep.Errors {
				output += fmt.Sprintf("- %s\n", err)
			}
		}
		output += "\n"
	}
	return output
}
