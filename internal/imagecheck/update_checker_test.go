package imagecheck

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRegistryClient struct {
	mock.Mock
}

func (m *MockRegistryClient) ListTags(ctx context.Context, repository string) ([]string, error) {
	args := m.Called(ctx, repository)
	return args.Get(0).([]string), args.Error(1)
}

func TestGetRemoteLatestTag(t *testing.T) {
	testCases := []struct {
		name           string
		mockTags       []string
		imageRef       string
		expectedLatest string
		expectError    bool
	}{
		{
			name:           "Semantic version tags",
			mockTags:       []string{"1.0.0", "1.1.0", "1.0.1", "2.0.0"},
			imageRef:       "nginx:latest",
			expectedLatest: "2.0.0",
			expectError:    false,
		},
		{
			name:           "Prefixed version tags",
			mockTags:       []string{"v1.0.0", "v1.1.0", "v0.9.9"},
			imageRef:       "alpine:latest",
			expectedLatest: "v1.1.0",
			expectError:    false,
		},
		{
			name:           "Single tag",
			mockTags:       []string{"1.0.0"},
			imageRef:       "busybox:latest",
			expectedLatest: "1.0.0",
			expectError:    false,
		},
		{
			name:           "Full image reference",
			mockTags:       []string{"1.0.0", "1.1.0"},
			imageRef:       "my.registry.com/library/busybox:1.0.0",
			expectedLatest: "1.1.0",
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock registry client
			mockClient := new(MockRegistryClient)
			mockClient.On("ListTags", mock.Anything, mock.Anything).Return(tc.mockTags, nil)

			// Create update checker with mock client
			uc := NewUpdateChecker(mockClient)

			// Call method
			latestTag, err := uc.GetRemoteLatestTag(context.Background(), tc.imageRef)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLatest, latestTag)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestParseImageReference(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedTag   string
		expectedError bool
	}{
		{
			name:          "Valid image with tag",
			input:         "nginx:1.19",
			expectedTag:   "1.19",
			expectedError: false,
		},
		{
			name:          "Image without tag",
			input:         "nginx",
			expectedTag:   "latest",
			expectedError: false,
		},
		{
			name:          "Image with tag with v prefix",
			input:         "nginx:v1.21",
			expectedTag:   "v1.21",
			expectedError: false,
		},
		{
			name:          "Full image reference without tag",
			input:         "docker.io/library/nginx",
			expectedTag:   "latest",
			expectedError: false,
		},
		{
			name:          "Full image reference with tag",
			input:         "docker.io/library/nginx:1.20",
			expectedTag:   "1.20",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := parseImageReference(tc.input)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTag, extractImageTag(ref.String()))
			}
		})
	}
}

func TestExtractImageRepository(t *testing.T) {
	tests := []struct {
		name      string
		imageName string
		want      string
	}{
		{
			name:      "Image name from Docker library",
			imageName: "nginx",
			want:      "nginx",
		},
		{
			name:      "Full image name from Docker library",
			imageName: "docker.io/library/nginx",
			want:      "docker.io/library/nginx",
		},
		{
			name:      "Image name with repository",
			imageName: "mylibrary/nginx",
			want:      "mylibrary/nginx",
		},
		{
			name:      "Image name with repository and registry",
			imageName: "my.registry.com/mylibrary/nginx",
			want:      "my.registry.com/mylibrary/nginx",
		},
		{
			name:      "Image name with tag",
			imageName: "nginx:v1.21",
			want:      "nginx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractImageRepository(tt.imageName); got != tt.want {
				t.Errorf("extractImageRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortSemanticVersionTags(t *testing.T) {
	tests := []struct {
		name    string
		rawTags []string
		want    []string
		wantErr bool
	}{
		{
			name:    "List of tags with the v prefix",
			rawTags: []string{"v2.1.0", "v2.0.0", "v1.2.3"},
			want:    []string{"v1.2.3", "v2.0.0", "v2.1.0"},
			wantErr: false,
		},
		{
			name:    "List of tags without the v prefix",
			rawTags: []string{"2.1.0", "2.0.0", "1.2.3"},
			want:    []string{"1.2.3", "2.0.0", "2.1.0"},
			wantErr: false,
		},
		{
			name:    "Mixed tags",
			rawTags: []string{"v2.1.0", "2.0.0", "v1.2.3"},
			want:    []string{"v1.2.3", "2.0.0", "v2.1.0"},
			wantErr: false,
		},
		{
			name:    "List of tags with suffix",
			rawTags: []string{"v1.0.0-alpha2", "v1.0.0-alpha3", "v1.0.0-beta0", "v1.0.0-alpha0"},
			want:    []string{"v1.0.0-alpha0", "v1.0.0-alpha2", "v1.0.0-alpha3", "v1.0.0-beta0"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortSemanticVersionTags(tt.rawTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortSemanticVersionTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortSemanticVersionTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
