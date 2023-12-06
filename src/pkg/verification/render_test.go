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
	assert.Equal(t, "## Step 1\n✅ **Success**\n\n## Step 2\n❌ **Failure**\n- Error 1\n- Error 2\n\n## Step 3\n⚠️ **Not Run**\n\n## Step 4\n⚠️ **Skipped**\n### Sub Step 1\n✅ **Success**\n\n", rendered)
}

func TestRender_Remarks(t *testing.T) {
	result := Result{}
	s := result.AddStep("Step 1", StatusSuccess)
	s.Remarks = append(s.Remarks, "Remark 1", "Remark 2")

	rendered := result.RenderMarkdown()
	assert.Equal(t, "## Step 1\n> [!NOTE]\n> Remark 1\n\n> [!NOTE]\n> Remark 2\n\n✅ **Success**\n\n", rendered)
}
