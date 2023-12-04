package verification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	result := Result{}
	result.AddStep("Step 1", StatusSuccess)
	result.AddStep("Step 2", StatusFailure, "Error 1", "Error 2")
	result.AddStep("Step 3", StatusNotRun)
	s := result.AddStep("Step 4", StatusSkipped)

	s.AddStep("Sub Step 1", StatusSuccess)

	rendered := result.RenderMarkdown()
	assert.Equal(t, rendered, "## Step 1\n✅ **Success**\n\n## Step 2\n❌ **Failure**\n- Error 1\n- Error 2\n\n## Step 3\n⚠️ **Not Run**\n\n## Step 4\n⚠️ **Skipped**\n### Sub Step 1\n✅ **Success**\n\n")
}

func TestDidFail_MainStepFailed(t *testing.T) {
	result := Result{}
	result.AddStep("Step 1", StatusSuccess)
	result.AddStep("Step 2", StatusFailure, "Error 1", "Error 2")
	result.AddStep("Step 3", StatusNotRun)
	s := result.AddStep("Step 4", StatusSkipped)

	s.AddStep("Sub Step 1", StatusSuccess)

	assert.True(t, result.DidFail())
}

func TestDidFail_SubStepFailed(t *testing.T) {
	result := Result{}
	result.AddStep("Step 1", StatusSuccess)
	result.AddStep("Step 2", StatusSuccess)
	result.AddStep("Step 3", StatusNotRun)
	s := result.AddStep("Step 4", StatusSkipped)

	s.AddStep("Sub Step 1", StatusFailure)

	assert.True(t, result.DidFail())
}
