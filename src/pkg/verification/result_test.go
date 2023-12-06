package verification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
