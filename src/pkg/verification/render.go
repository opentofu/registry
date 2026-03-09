package verification

import "fmt"

func (r *Result) RenderMarkdown() string {
	var output string
	for _, step := range r.Steps {
		output += fmt.Sprintf("## %s\n", step.Name)
		for _, remark := range step.Remarks {
			output += "> [!NOTE]\n"
			output += fmt.Sprintf("> %s\n\n", remark)
		}
		switch step.Status {
		case StatusSuccess:
			output += "✅ **Success**\n"
		case StatusFailure:
			output += "❌ **Failure**\n"
		case StatusNotRun:
			output += "⚠️ **Not Run**\n"
		case StatusSkipped:
			output += "⚠️ **Skipped**\n"
		case StatusWarning:
			output += "⚠️ **Warning**\n"
		}

		for _, err := range step.Errors {
			output += fmt.Sprintf("- %s\n", err)
		}
		for _, subStep := range step.SubSteps {
			output += fmt.Sprintf("### %s\n", subStep.Name)
			for _, remark := range subStep.Remarks {
				output += "> [!NOTE]\n"
				output += fmt.Sprintf("> %s\n\n", remark)
			}
			switch subStep.Status {
			case StatusSuccess:
				output += "✅ **Success**\n"
			case StatusFailure:
				output += "❌ **Failure**\n"
			case StatusNotRun:
				output += "⚠️ **Not Run**\n"
			case StatusSkipped:
				output += "⚠️ **Skipped**\n"
			case StatusWarning:
				output += "⚠️ **Warning**\n"
			}

			for _, err := range subStep.Errors {
				output += fmt.Sprintf("- %s\n", err)
			}
		}
		output += "\n"
	}
	if r.DidFail() {
		output += "\nAfter the issue is fixed, update the title or the description of the issue to retrigger the submission workflow."
		output += "\n"
	}
	return output
}
