package gpg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKey(t *testing.T) {
	stringPtr := func(s string) *string {
		return &s
	}

	privateKey := `-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC/rbBBBC0iJkW3\nbRnKICpWrVhEsjmgSPs/cF7KihVYhriGspcgpaF8RmVw+k98Wj0SmO+sGc81pfDv\nxpl4PrN9lTsvvkxF4CfNXZr21vLaaC32jD8HnZtEuWsA+bByd7lEKw0N0652cWjm\nRD8/CYUSMFWe4VM1tjPMD1l9/oL8T6/z2nGHgMH+hXkES2aYOclfyqwHkHnLObMf\nZGvYxp0CZ+fbSW+KRfL6qHaM03/1ZsofNS2phJhk341R09suxn1fWpkhyWPtUSZF\n1TVUEjgCkUlppTvoCtE/UAsM17ey5AZKlhFxgh0IvN21sFeXFtRjzQrcZTUA5Dj5\n7lYaicmJAgMBAAECggEBALcloq+86dMjdqHZITc8nLfNUfXxxZYdpdPr7ubgIZ1A\nvLgXlMeg+zffm7Xjtmc/YfOPJhLvZkoAkMLKpIF8h8yK9s6bqg1qLR3RPux0Xf/K\nY4CcaO1B7sYv1MpNygbV1rQH3qVDigOqQW0j8LquwfOrM2RoMDW2Lq/gSsZUlZu2\nb6OpJJWRTDEbHsJausOU8Yp/u2ZOv0fIGE65cU/oyX4pcJzSnXXrrB6bdEJfkgJ2\nejzK3rvdz97iE7N8i/9o3mBNTTZgnO6sqGH4vIa5dZpEuAgUpYPanbO27ELRkA2K\nMsyukonBdkGhNAAEIHfYErFw+a8C/29dXwCe1C8ipDECgYEA8ZMpccSnIssnV2xB\nRP5AXTkyQMrBqtqGzmcs1BT3D1TsOywe5U3OaDZlFcdEJq3DIgYheNZsalQNy55C\nPQO6Mh4McXHlAQ7BA4ukYg+ww3E5LbXo+5TkPpHgq6iIxSNusOk3boloMEhlT56G\nIwY2YHgesiCQWvc2hC92+2QGgpcCgYEAyx/I46PSThkoRCxzIHKERJA6XghbNI+y\nnXW0N0M87iVeT261FhTbOBx8uzGzTnqTXdgBZLO2Kdw4r4stR8cMXlFESciX7Jw6\n+aOJrNMIXXI/lhTqhcni6ncPYgipM77a9quvx/oGEpx/0yitkFY772z3UMUiH8Wh\nzuXZbxi7ON8CgYBMuc/M+YeYDmwVYSWt0w8ATN1AJOWz7Sopvi1HwszhSrio5o99\nhuPKx5P9gceMfV3fnZDd/0R51O54wHALTvbBWjfbhDAW0OfOx3hTSOZ8fKaLdR5l\nYVnI4a449xNRgbpzZ+8aJXw48ZVz30Z9M0jsBNrC+oK+0Yu4GhcxKwjCSwKBgHjN\nLmwzwZ8w1wG0bcOeV4tvO0cxMQzRaSi8F7HGCzaWgsA61veK79UvG/84T6scuwfU\nrv904aGDlzLPUt6dQn3VVweKhM/zGh/dYsOlvhPVHnvjdJacupc2t69V90sO9qo8\n8Q29ZF8tM9ghGRf+MSbzZyJiGylKIDEsAWRREQeBAoGAPfqkmisNWKuHpcGsyYjT\nAkPVev7SZrUW9YwMcOa2+aA56CI/zKCO7U0ZvVf6TEWnPzT/UGedkYFP+/6vQxUh\n2TxMmjU7l4C829k5ZZ0OlG6IHvtaK2vkJmEKWI7sGE2GMewwXEsfjJs0y+Szxb6U\nzPRa6IiVxnJAiqbjs+W/fC4=\n-----END PRIVATE KEY-----\n`

	tests := []struct {
		name        string
		data        string
		expectedErr *string
	}{
		{
			name:        "private key should fail",
			data:        privateKey,
			expectedErr: stringPtr("bleh"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseKey(test.data)

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, "could not build public key from ascii armor")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
