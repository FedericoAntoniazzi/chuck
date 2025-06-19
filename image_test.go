package main

import (
	"reflect"
	"testing"
)

func TestParseImageName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Image
		wantErr  bool
	}{
		{
			name:  "Short Docker Hub image - explicit latest tag",
			input: "nginx:latest",
			expected: Image{
				Raw:       "nginx:latest",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Short Docker Hub image - implicit latest tag",
			input: "nginx",
			expected: Image{
				Raw:       "nginx",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Short Docker Hub image - versioned tag",
			input: "nginx:1.25",
			expected: Image{
				Raw:       "nginx:1.25",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "1.25",
			},
			wantErr: false,
		},
		{
			name:  "Standard Docker Hub image with namespace - implicit latest tag",
			input: "library/nginx",
			expected: Image{
				Raw:       "library/nginx",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Standard Docker Hub image with namespace - explicit latest tag",
			input: "library/nginx:latest",
			expected: Image{
				Raw:       "library/nginx:latest",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Standard Docker Hub image with namespace - versioned tag",
			input: "library/nginx:1.25",
			expected: Image{
				Raw:       "library/nginx:1.25",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "1.25",
			},
			wantErr: false,
		},
		{
			name:  "Extended Docker Hub image with namespace - implicit latest tag",
			input: "docker.io/library/nginx",
			expected: Image{
				Raw:       "docker.io/library/nginx",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Extended Docker Hub image with namespace - explicit latest tag",
			input: "docker.io/library/nginx:latest",
			expected: Image{
				Raw:       "docker.io/library/nginx:latest",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Extended Docker Hub image with namespace - versioned tag",
			input: "docker.io/library/nginx:1.25",
			expected: Image{
				Raw:       "docker.io/library/nginx:1.25",
				Registry:  "docker.io",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "1.25",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry - implicit latest tag",
			input: "myregistry.com/library/nginx",
			expected: Image{
				Raw:       "myregistry.com/library/nginx",
				Registry:  "myregistry.com",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry - explicit latest tag",
			input: "myregistry.com/library/nginx:latest",
			expected: Image{
				Raw:       "myregistry.com/library/nginx:latest",
				Registry:  "myregistry.com",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry - versioned tag",
			input: "myregistry.com/library/nginx:1.25",
			expected: Image{
				Raw:       "myregistry.com/library/nginx:1.25",
				Registry:  "myregistry.com",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "1.25",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry with port - implicit latest tag",
			input: "myregistry.com:8080/library/nginx",
			expected: Image{
				Raw:       "myregistry.com:8080/library/nginx",
				Registry:  "myregistry.com:8080",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry with port - explicit latest tag",
			input: "myregistry.com:8080/library/nginx:latest",
			expected: Image{
				Raw:       "myregistry.com:8080/library/nginx:latest",
				Registry:  "myregistry.com:8080",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "latest",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry with port - versioned tag",
			input: "myregistry.com:8080/library/nginx:1.25",
			expected: Image{
				Raw:       "myregistry.com:8080/library/nginx:1.25",
				Registry:  "myregistry.com:8080",
				Namespace: "library",
				Name:      "nginx",
				Tag:       "1.25",
			},
			wantErr: false,
		},
		{
			name:  "Custom registry with port - versioned tag",
			input: "myregistry.com:8080/myorg/myproject/image:1.25",
			expected: Image{
				Raw:       "myregistry.com:8080/myorg/myproject/image:1.25",
				Registry:  "myregistry.com:8080",
				Namespace: "myorg/myproject",
				Name:      "image",
				Tag:       "1.25",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseImageName(tc.input)
			if tc.wantErr && err == nil {
				t.Errorf("ParseImageName() expected error but got nil")
			}
			if tc.wantErr && err != nil {
				t.Errorf("ParseImageName() unexpected error = %v", err)
			}
			if !tc.wantErr && !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("ParseImageName() got = %+v, expected %+v", got, tc.expected)
			}
		})
	}
}
