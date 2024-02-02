package provider

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type TestCase struct {
		name       string
		input      Metadata
		wantErrStr string
	}
	tests := []TestCase{
		{
			name: "valid",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"6.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "",
		},
		{
			name:       "no versions data",
			input:      Metadata{},
			wantErrStr: "found empty list of versions\n",
		},
		{
			name: "invalid version",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "foo",
						Protocols:           []string{"6.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
					{
						Version:             "0.0.2",
						Protocols:           []string{"6.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "vfoo: found semver-incompatible version: foo\n",
		},
		{
			name: "no protocols",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "v0.0.2: empty protocols list\n",
		},
		{
			name: "invalid protocol",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0", "foo"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "v0.0.2: unsupported protocol found: foo\n",
		},
		{
			name: "no targets data",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
					},
				},
			},
			wantErrStr: "v0.0.2: empty targets list\n",
		},
		{
			name: "invalid target os",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "foo",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "v0.0.2: target foo-amd64: unsupported OS: foo\n",
		},
		{
			name: "invalid target arch",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "foo",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "v0.0.2: target darwin-foo: unsupported ARCH: foo\n",
		},
		{
			name: "filename does not match url",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "amd64",
								Filename:    "foobar.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897",
							},
						},
					},
				},
			},
			wantErrStr: "v0.0.2: target darwin-amd64: 'filename' is not consistent with 'download_url'\n",
		},
		{
			name: "target shasum length is wrong",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "darwin",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea889",
							},
						},
					},
				},
			},
			wantErrStr: "v0.0.2: target linux-arm64: SHASum length is wrong\n",
		},
		{
			name: "multiple errors",
			input: Metadata{
				Versions: []Version{
					{
						Version:             "0.0.1",
						Protocols:           []string{"5.0", "xxx"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.1/terraform-provider-metal_0.0.1_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.1/terraform-provider-metal_0.0.1_SHA256SUMS.sig",
					},
					{
						Version:             "0.0.2",
						Protocols:           []string{"5.0"},
						SHASumsURL:          "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
						SHASumsSignatureURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
						Targets: []Target{
							{
								OS:          "foo",
								Arch:        "amd64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_amd64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
								SHASum:      "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a",
							},
							{
								OS:          "darwin",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_darwin_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
								SHASum:      "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195",
							},
							{
								OS:          "linux",
								Arch:        "arm64",
								Filename:    "terraform-provider-metal_0.0.2_linux_arm64.zip",
								DownloadURL: "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
								SHASum:      "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea889",
							},
						},
					},
				},
			},
			wantErrStr: `v0.0.1: unsupported protocol found: xxx
v0.0.1: empty targets list
v0.0.2: target foo-amd64: unsupported OS: foo
v0.0.2: target linux-arm64: SHASum length is wrong
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			switch tt.wantErrStr != "" {
			case true:
				if err == nil || tt.wantErrStr != err.Error() {
					t.Fatalf("unexpected error message, want = %s, got = %v", tt.wantErrStr, err)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error message: %v", err)
				}
			}
		})
	}
}
