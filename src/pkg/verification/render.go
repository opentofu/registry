package verification

import "fmt"

func (r *Result) RenderMarkdown() string {
	var output string
	for _, step := range r.Steps {
		output += fmt.Sprintf("## %s\n", step.Name)
		for _, remark := range step.Remarks {
			output += fmt.Sprintf("> [!NOTE]\n")
			output += fmt.Sprintf("> %s\n\n", remark)
		}
		if step.Status == StatusSuccess {
			output += "✅ **Success**\n"
		} else if step.Status == StatusFailure {
			output += "❌ **Failure**\n"
		} else if step.Status == StatusNotRun {
			output += "⚠️ **Not Run**\n"
		} else if step.Status == StatusSkipped {
			output += "⚠️ **Skipped**\n"
		} else if step.Status == StatusWarning {
			output += "⚠️ **Warning**\n"
		}

		for _, err := range step.Errors {
			output += fmt.Sprintf("- %s\n", err)
		}
		for _, subStep := range step.SubSteps {
			output += fmt.Sprintf("### %s\n", subStep.Name)
			for _, remark := range subStep.Remarks {
				output += fmt.Sprintf("> [!NOTE]\n")
				output += fmt.Sprintf("> %s\n\n", remark)
			}
			if subStep.Status == StatusSuccess {
				output += "✅ **Success**\n"
			} else if subStep.Status == StatusFailure {
				output += "❌ **Failure**\n"
			} else if subStep.Status == StatusNotRun {
				output += "⚠️ **Not Run**\n"
			} else if subStep.Status == StatusSkipped {
				output += "⚠️ **Skipped**\n"
			} else if subStep.Status == StatusWarning {
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
