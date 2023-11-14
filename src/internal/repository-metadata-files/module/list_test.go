package module

import (
	"testing"

	"registry-stable/internal/module"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractModuleDetailsFromPath(t *testing.T) {
	type TestCase struct {
		name           string
		input          string
		expectedOutput *module.Module
	}

	testCases := []TestCase{
		{
			name:  "Valid module path",
			input: "modules/t/terraform-aws-modules/lambda/aws.json",
			expectedOutput: &module.Module{
				Namespace:    "terraform-aws-modules",
				Name:         "lambda",
				TargetSystem: "aws",
			},
		},
		{
			name:  "Valid module path with mixed case",
			input: "modules/t/terraform-aws-modules/lambda/AWS.json",
			expectedOutput: &module.Module{
				Namespace:    "terraform-aws-modules",
				Name:         "lambda",
				TargetSystem: "AWS",
			},
		},
		{
			name:  "Valid module path with numbers",
			input: "modules/t/terraform-aws-modules1234/lambda1234/aws1234.json",
			expectedOutput: &module.Module{
				Namespace:    "terraform-aws-modules1234",
				Name:         "lambda1234",
				TargetSystem: "aws1234",
			},
		},
		{
			name:  "Valid module path with hypens and underscores",
			input: "modules/t/terraform-aws-modules_test-abcd/lambda_test-abcd/aws_test-abcd.json",
			expectedOutput: &module.Module{
				Namespace:    "terraform-aws-modules_test-abcd",
				Name:         "lambda_test-abcd",
				TargetSystem: "aws_test-abcd",
			},
		},
		{
			name:           "Invalid module path (no .json)",
			input:          "modules/t/terraform-aws-modules/lambda/aws",
			expectedOutput: nil,
		},
		{
			name:           "Invalid module path (empty)",
			input:          "",
			expectedOutput: nil,
		},
		{
			name:           "Invalid module path (missing namespace)",
			input:          "modules/t//lambda/aws.json",
			expectedOutput: nil,
		},
		{
			name:           "Invalid module path (missing name)",
			input:          "modules/t/terraform-aws-modules//aws.json",
			expectedOutput: nil,
		},
		{
			name:           "Invalid module path (missing target system)",
			input:          "modules/t/terraform-aws-modules/lambda/.json",
			expectedOutput: nil,
		},
		{
			name:           "Invalid module path (missing namespace, name, and target system)",
			input:          "modules/t///.json",
			expectedOutput: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, extractModuleDetailsFromPath(tc.input), "Extracted module details do not match expected output")
		})
	}
}
