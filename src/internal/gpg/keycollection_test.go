package gpg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeKeyFile(t *testing.T, dir string, filename string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0755))
	key, err := generateGPGKey()
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, filename), []byte(key), 0600))
}

func TestListKeys(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, root string)
		namespace     string
		providerName  string
		expectedCount int
	}{
		{
			name:          "no keys directory returns empty",
			setup:         func(_ *testing.T, _ string) {},
			namespace:     "acme",
			providerName:  "widget",
			expectedCount: 0,
		},
		{
			name: "namespace-level key is found",
			setup: func(t *testing.T, root string) {
				writeKeyFile(t, filepath.Join(root, "a", "acme"), "provider.asc")
			},
			namespace:     "acme",
			providerName:  "widget",
			expectedCount: 1,
		},
		{
			name: "provider-level key is found",
			setup: func(t *testing.T, root string) {
				writeKeyFile(t, filepath.Join(root, "a", "acme", "widget"), "provider.asc")
			},
			namespace:     "acme",
			providerName:  "widget",
			expectedCount: 1,
		},
		{
			name: "namespace and provider keys are both found",
			setup: func(t *testing.T, root string) {
				writeKeyFile(t, filepath.Join(root, "a", "acme"), "provider.asc")
				writeKeyFile(t, filepath.Join(root, "a", "acme", "widget"), "provider.asc")
			},
			namespace:     "acme",
			providerName:  "widget",
			expectedCount: 2,
		},
		{
			name: "keys for different namespace are not returned",
			setup: func(t *testing.T, root string) {
				writeKeyFile(t, filepath.Join(root, "o", "other"), "provider.asc")
			},
			namespace:     "acme",
			providerName:  "widget",
			expectedCount: 0,
		},
		{
			name: "key for different provider in same namespace is not returned",
			setup: func(t *testing.T, root string) {
				writeKeyFile(t, filepath.Join(root, "a", "acme", "gadget"), "provider.asc")
			},
			namespace:     "acme",
			providerName:  "widget",
			expectedCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			tc.setup(t, root)

			collection := KeyCollection{
				Namespace:    tc.namespace,
				ProviderName: tc.providerName,
				Directory:    root,
			}

			keys, err := collection.ListKeys()
			require.NoError(t, err)
			assert.Len(t, keys, tc.expectedCount)
		})
	}
}
