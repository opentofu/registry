package main

import (
	"os"
	"path"
	"testing"
)

func Test_validateModuleFile(t *testing.T) {
	tests := map[string]struct {
		fileContent []byte
		wantErr     bool
	}{
		"valid-file": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "0.0.2"
    },
	{
      "version": "0.0.1"
    }
  ]
}`),
			wantErr: false,
		},
		"invalid-json": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "0.0.1"
    },
  ]
}`),
			wantErr: true,
		},
		"invalid-version": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "0.0.1"
    },
	{
      "version": "foo"
    }
  ]
}`),
			wantErr: true,
		},
		"empty-versions-list": {
			fileContent: []byte(`{"versions": []}`),
			wantErr:     true,
		},
		"no-versions-attr-found": {
			fileContent: []byte(`{}`),
			wantErr:     true,
		},
	}

	t.Parallel()

	for name, tt := range tests {
		dir := t.TempDir()
		p := path.Join(dir, "module.json")
		if err := os.WriteFile(p, tt.fileContent, 0774); err != nil {
			t.Fatalf("cannot create temp file: %v", err)
		}

		t.Run(name, func(t *testing.T) {
			if err := validateModuleFile(p); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateProviderFile(t *testing.T) {
	tests := map[string]struct {
		fileContent []byte
		wantErr     bool
	}{
		"valid": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "0.0.2",
      "protocols": [
        "6.0"
      ],
      "shasums_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
      "shasums_signature_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
      "targets": [
        {
          "os": "darwin",
          "arch": "amd64",
          "filename": "terraform-provider-metal_0.0.2_darwin_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
          "shasum": "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a"
        },
        {
          "os": "darwin",
          "arch": "arm64",
          "filename": "terraform-provider-metal_0.0.2_darwin_arm64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_arm64.zip",
          "shasum": "2da0c86252ff399567852412a8b750ca7ab0ffa48bd8e335edc0e8cec4e12195"
        },
        {
          "os": "freebsd",
          "arch": "386",
          "filename": "terraform-provider-metal_0.0.2_freebsd_386.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_freebsd_386.zip",
          "shasum": "50ad273e5f51c9b63d2e458b5266d9f50727a78d1db2d34940e3845744f0a141"
        },
        {
          "os": "freebsd",
          "arch": "amd64",
          "filename": "terraform-provider-metal_0.0.2_freebsd_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_freebsd_amd64.zip",
          "shasum": "83522ae20f82fddb11cbc12cb3f60829b27bcb71ddcb1b6c07728731afda5059"
        },
        {
          "os": "freebsd",
          "arch": "arm",
          "filename": "terraform-provider-metal_0.0.2_freebsd_arm.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_freebsd_arm.zip",
          "shasum": "7b181fa369d508aacb05db1b23441925b24fb4e3a6b68998d4d47e4f541651db"
        },
        {
          "os": "freebsd",
          "arch": "arm64",
          "filename": "terraform-provider-metal_0.0.2_freebsd_arm64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_freebsd_arm64.zip",
          "shasum": "a213e541b786f4990af9d8a596877450f5a5e3deb593f566b5b977c3df74ae8f"
        },
        {
          "os": "linux",
          "arch": "386",
          "filename": "terraform-provider-metal_0.0.2_linux_386.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_386.zip",
          "shasum": "d2eea09e5c1607691e003db950ffc79274d791bf190c3b246fa041dbe7c15a78"
        },
        {
          "os": "linux",
          "arch": "amd64",
          "filename": "terraform-provider-metal_0.0.2_linux_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_amd64.zip",
          "shasum": "8452e035c92bf319a76e1dac8df11147f550affc85017237ae5cc22c9254bd1f"
        },
        {
          "os": "linux",
          "arch": "arm",
          "filename": "terraform-provider-metal_0.0.2_linux_arm.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm.zip",
          "shasum": "facb0f30c4b9171c42ac0995d6e99b4cf0dfa788ada353b06436afcecf72e412"
        },
        {
          "os": "linux",
          "arch": "arm64",
          "filename": "terraform-provider-metal_0.0.2_linux_arm64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_linux_arm64.zip",
          "shasum": "6c17811307361f07919514efc0b2eb7637a9c9f580d7d325871002ae23ea8897"
        },
        {
          "os": "windows",
          "arch": "386",
          "filename": "terraform-provider-metal_0.0.2_windows_386.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_windows_386.zip",
          "shasum": "e31b3495ac451450b1ba3ff96b0e3bde1f79d2d2b0ae6f1166cbe2418b2ac70a"
        },
        {
          "os": "windows",
          "arch": "amd64",
          "filename": "terraform-provider-metal_0.0.2_windows_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_windows_amd64.zip",
          "shasum": "329e8fcc52e716050f2a153a90f4da7e42e40ddb71a8d593d364db181f7f3cac"
        },
        {
          "os": "windows",
          "arch": "arm",
          "filename": "terraform-provider-metal_0.0.2_windows_arm.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_windows_arm.zip",
          "shasum": "d3866da93888cf2a6f38b4bf04f8a6eef8bb40fb639f9ef275afa38e39fac466"
        },
        {
          "os": "windows",
          "arch": "arm64",
          "filename": "terraform-provider-metal_0.0.2_windows_arm64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_windows_arm64.zip",
          "shasum": "92b428195b230b7dd0ad4300138a6e5a65dc6e87376f0799436158f7943270cf"
        }
      ]
    }
  ]
}`),
			wantErr: false,
		},
		"empty-versions-list": {
			fileContent: []byte(`{"versions": []}`),
			wantErr:     true,
		},
		"no-versions-attr-found": {
			fileContent: []byte(`{}`),
			wantErr:     true,
		},
		"invalid-version": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "foo",
      "protocols": [
        "6.0"
      ],
      "shasums_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/vfoo/terraform-provider-metal_foo_SHA256SUMS",
      "shasums_signature_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/vfoo/terraform-provider-metal_foo_SHA256SUMS.sig",
      "targets": [
        {
          "os": "darwin",
          "arch": "amd64",
          "filename": "terraform-provider-metal_foo_darwin_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/vfoo/terraform-provider-metal_foo_darwin_amd64.zip",
          "shasum": "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a"
        }
	  ]
    }
  ]
}
`),
			wantErr: true,
		},
		"empty-protocols-list": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "0.0.2",
      "protocols": [],
      "shasums_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
      "shasums_signature_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
      "targets": [
        {
          "os": "darwin",
          "arch": "amd64",
          "filename": "terraform-provider-metal_0.0.2_darwin_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
          "shasum": "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a"
        }
	  ]
    }
  ]
}
`),
			wantErr: true,
		},
		"no-protocols-attr-found": {
			fileContent: []byte(`{
  "versions": [
    {
      "version": "0.0.2",
      "shasums_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS",
      "shasums_signature_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_SHA256SUMS.sig",
      "targets": [
        {
          "os": "darwin",
          "arch": "amd64",
          "filename": "terraform-provider-metal_0.0.2_darwin_amd64.zip",
          "download_url": "https://github.com/metal-stack-cloud/terraform-provider-metal/releases/download/v0.0.2/terraform-provider-metal_0.0.2_darwin_amd64.zip",
          "shasum": "0192019c49306991cf2413d937315a5a5aa83e8a3e01ec0891d3a8460422683a"
        }
	  ]
    }
  ]
}
`),
			wantErr: true,
		},
	}

	t.Parallel()

	for name, tt := range tests {
		dir := t.TempDir()
		p := path.Join(dir, "module.json")
		if err := os.WriteFile(p, tt.fileContent, 0774); err != nil {
			t.Fatalf("cannot create temp file: %v", err)
		}

		t.Run(name, func(t *testing.T) {
			if err := validateProviderFile(p); (err != nil) != tt.wantErr {
				t.Errorf("validateProviderFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
