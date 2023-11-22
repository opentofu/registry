package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractProviderDetailsFromPath(t *testing.T) {
	type TestCase struct {
		name           string
		input          string
		expectedOutput *Provider
	}

	testCases := []TestCase{
		{
			name:  "Valid provider path",
			input: "providers/o/opentofu/aws.json",
			expectedOutput: &Provider{
				Namespace:    "opentofu",
				ProviderName: "aws",
			},
		},
		{
			name:  "Valid provider path with numbers",
			input: "providers/o/opentofu1234/aws1234.json",
			expectedOutput: &Provider{
				Namespace:    "opentofu1234",
				ProviderName: "aws1234",
			},
		},
		{
			name:  "Valid provider path with hypens and underscores",
			input: "providers/o/opentofu_test-abcd/aws_test-abcd.json",
			expectedOutput: &Provider{
				Namespace:    "opentofu_test-abcd",
				ProviderName: "aws_test-abcd",
			},
		},
		{
			name:           "Invalid provider path (no .json)",
			input:          "providers/o/opentofu/aws",
			expectedOutput: nil,
		},
		{
			name:           "Invalid provider path (empty)",
			input:          "",
			expectedOutput: nil,
		},
		{
			name:           "Invalid provider path (missing namespace)",
			input:          "providers/o//aws.json",
			expectedOutput: nil,
		},
		{
			name:           "Invalid provider path (missing provider name)",
			input:          "providers/o/opentofu/.json",
			expectedOutput: nil,
		},
		{
			name:           "Invalid provider path (missing namespace and name)",
			input:          "providers/o//.json",
			expectedOutput: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, extractProviderDetailsFromPath(tc.input), "Extracted provider details do not match expected output")
		})
	}
}
