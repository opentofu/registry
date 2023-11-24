package provider

import "testing"

const validShasumsFile = `7c4828b800cbc598c8e12fc7c812543317fb6782676012bbb4476e5f36048976  terraform-provider-aws_5.26.0_darwin_amd64.zip
34edd4e60f84ea5f3ad450416518bb47bfae5fc6666ccdbd13c6e2fbb6c9c4fe  terraform-provider-aws_5.26.0_darwin_arm64.zip
baceae64657a5acc3ae17cd8344934b019dd340e757606b07c1281e5f7ea6914  terraform-provider-aws_5.26.0_linux_386.zip
489ee1dc1ec1393060d6e84e0d732610bbcb5448d522585b2de1d0893b979ec7  terraform-provider-aws_5.26.0_linux_amd64.zip
02afd760c2b407d6dc207dccd5e6ae0bbda986c9ba5c540dedabfcd57e13bcf1  terraform-provider-aws_5.26.0_linux_arm.zip
3a36d9f48bd0d29487c7d8e8ce8b9366d61c69533b3608503fd201e20b0f2362  terraform-provider-aws_5.26.0_linux_arm64.zip
3513f3120cc4c3442836e70898a4ec51c42c192cc96547020d8cbc92a41d1e85  terraform-provider-aws_5.26.0_windows_386.zip
950dbc98636a69c3f9b6b713e3a58a2255046ac840264e571b32f4eb47be0763  terraform-provider-aws_5.26.0_windows_amd64.zip
3169efa20240f4c76bf9b33be8f74e73afd0de7b0bea6074b930ca685ab75010  terraform-provider-aws_5.26.0_windows_arm.zip
205a21409309a086a82bee4ec055e805c56e47f69d272ef4bf6a8c5f643e8260  terraform-provider-aws_5.26.0_windows_arm64.zip
`

func TestShaFileToMap(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  map[string]string
	}{
		{
			name:  "empty file",
			input: []byte(``),
			want:  map[string]string{},
		},
		{
			name:  "valid file",
			input: []byte(validShasumsFile),
			want: map[string]string{
				"terraform-provider-aws_5.26.0_darwin_amd64.zip":  "7c4828b800cbc598c8e12fc7c812543317fb6782676012bbb4476e5f36048976",
				"terraform-provider-aws_5.26.0_darwin_arm64.zip":  "34edd4e60f84ea5f3ad450416518bb47bfae5fc6666ccdbd13c6e2fbb6c9c4fe",
				"terraform-provider-aws_5.26.0_linux_386.zip":     "baceae64657a5acc3ae17cd8344934b019dd340e757606b07c1281e5f7ea6914",
				"terraform-provider-aws_5.26.0_linux_amd64.zip":   "489ee1dc1ec1393060d6e84e0d732610bbcb5448d522585b2de1d0893b979ec7",
				"terraform-provider-aws_5.26.0_linux_arm.zip":     "02afd760c2b407d6dc207dccd5e6ae0bbda986c9ba5c540dedabfcd57e13bcf1",
				"terraform-provider-aws_5.26.0_linux_arm64.zip":   "3a36d9f48bd0d29487c7d8e8ce8b9366d61c69533b3608503fd201e20b0f2362",
				"terraform-provider-aws_5.26.0_windows_386.zip":   "3513f3120cc4c3442836e70898a4ec51c42c192cc96547020d8cbc92a41d1e85",
				"terraform-provider-aws_5.26.0_windows_amd64.zip": "950dbc98636a69c3f9b6b713e3a58a2255046ac840264e571b32f4eb47be0763",
				"terraform-provider-aws_5.26.0_windows_arm.zip":   "3169efa20240f4c76bf9b33be8f74e73afd0de7b0bea6074b930ca685ab75010",
				"terraform-provider-aws_5.26.0_windows_arm64.zip": "205a21409309a086a82bee4ec055e805c56e47f69d272ef4bf6a8c5f643e8260",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shaFileToMap(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("shaFileToMap() = %v, want %v", got, tt.want)
			}
			for k, v := range got {
				if tt.want[k] != v {
					t.Errorf("shaFileToMap() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
