package sources

import (
	"reflect"
	"testing"
)

func TestNewImage(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Image
	}{
		{
			name:  "Short image name and no tag",
			input: "nginx",
			expected: Image{
				Name:    "nginx",
				Version: "latest",
			},
		},
		{
			name:  "Full image name and no tag",
			input: "docker.io/library/nginx",
			expected: Image{
				Name:    "docker.io/library/nginx",
				Version: "latest",
			},
		},
		{
			name:  "Short image name and latest tag",
			input: "nginx:latest",
			expected: Image{
				Name:    "nginx",
				Version: "latest",
			},
		},
		{
			name:  "Full image name and latest tag",
			input: "docker.io/library/nginx:latest",
			expected: Image{
				Name:    "docker.io/library/nginx",
				Version: "latest",
			},
		},
		{
			name:  "Short image name and semver tag",
			input: "nginx:v1.2.3",
			expected: Image{
				Name:    "nginx",
				Version: "v1.2.3",
			},
		},
		{
			name:  "Full image name and semver tag",
			input: "docker.io/library/nginx:v1.2.3",
			expected: Image{
				Name:    "docker.io/library/nginx",
				Version: "v1.2.3",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			got := NewImage(c.input)
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("expected %v, got %v", c.expected, got)
			}
		})
	}
}
