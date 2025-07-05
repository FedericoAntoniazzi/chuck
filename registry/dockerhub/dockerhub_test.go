package dockerhub

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FedericoAntoniazzi/chuck/types"
	"github.com/stretchr/testify/assert"
)

// TestNewClient verifies that NewClient initializes the client correctly
func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 15*time.Second, client.httpClient.Timeout)
}

// TestGetTags_UnsupportedRegistry tests the case where the image registry is not docker.io
func TestGetTags_UnsupportedRegistry(t *testing.T) {
	client := NewClient()
	image := types.Image{
		Registry:  "gcr.io",
		Namespace: "google-samples",
		Name:      "hello-app",
	}

	tags, err := client.GetTags(context.Background(), image)
	assert.Nil(t, tags)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported url for Docker Hub")
}

// TestGetTags_Success tests a successful API call to Docker Hub
func TestGetTags_Success(t *testing.T) {
	expectedTags := []string{"latest", "v1.0.0", "v1.1.0"}
	mockResponse := dockerHubTagsResponse{
		Results: []dockerHubListTagResult{
			{Name: "latest"},
			{Name: "v1.0.0"},
			{Name: "v1.1.0"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure it is a GET request
		assert.Equal(t, http.MethodGet, r.Method)
		// Check the URL
		assert.Contains(t, r.URL.Path, "/namespaces/library/repositories/alpine/tags")

		// Return okay
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Override the dockerHubBaseURL for testing
	originalDockerHubBaseURL := dockerHubBaseURL
	dockerHubBaseURL = server.URL
	// Restore original URL
	defer func() { dockerHubBaseURL = originalDockerHubBaseURL }()

	client := NewClient()
	image := types.Image{
		Registry:  "docker.io",
		Namespace: "library",
		Name:      "alpine",
	}

	tags, err := client.GetTags(context.Background(), image)
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.ElementsMatch(t, expectedTags, tags)
}

// TestGetTags_NonOKStatus tests when Docker Hub returns a non-200 status code
func TestGetTags_NonOKStatus(t *testing.T) {
	// Create a mock HTTP server that returns a 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	originalDockerHubBaseURL := dockerHubBaseURL
	dockerHubBaseURL = server.URL
	defer func() { dockerHubBaseURL = originalDockerHubBaseURL }()

	client := NewClient()
	image := types.Image{
		Registry:  "docker.io",
		Namespace: "library",
		Name:      "alpine",
	}

	tags, err := client.GetTags(context.Background(), image)
	assert.Nil(t, tags)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "received non-OK status code from Docker Hub (500): 500 Internal Server Error (Body: internal server error)")
}

// TestGetTags_Integration tests the actual Docker Hub API to ensure contract
// NOTE: This test should be run selectively (e.g., in CI/CD nightly builds)
// and not as part of every unit test run, as it relies on external services.
func TestGetTags_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	client := NewClient()
	image := types.Image{
		Raw:       "docker.io/library/nginx:1.25",
		Registry:  "docker.io",
		Namespace: "library",
		Name:      "nginx",
		Tag:       "1.25",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Longer timeout for real API
	defer cancel()

	tags, err := client.GetTags(ctx, image)
	assert.NoError(t, err, "Integration test failed to fetch tags from real Docker Hub")
	assert.NotNil(t, tags, "Integration test: Tags should not be nil")
	assert.True(t, len(tags) > 0, "Integration test: Expected to find at least one tag for nginx")

	// You can add more specific assertions here if you have known tags,
	// but the primary goal is to ensure the API response structure is still compatible.
	assert.Contains(t, tags, "latest", "Integration test: 'latest' tag should be present for nginx")
}
