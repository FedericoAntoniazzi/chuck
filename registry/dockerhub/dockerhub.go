package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/FedericoAntoniazzi/chuck/types"
)

// dockerHubListTagsResult represents the details of each query for repository tags
type dockerHubListTagResult struct {
	Name string `json:"name,omitempty"`
}

// dockerHubTagsResponses represents the response of Docker Hub when querying tags
type dockerHubTagsResponse struct {
	Results []dockerHubListTagResult `json:"results,omitempty"`
}

var dockerHubBaseURL string = "https://registry.hub.docker.com/v2"

// Client is the DockerHub registry client
type Client struct {
	httpClient *http.Client
}

// NewClient creates and returns a new DockerHub client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetTags fetches all available tags for a given image from Docker Hub
func (c *Client) GetTags(ctx context.Context, image types.Image) ([]string, error) {
	if image.Registry != "docker.io" {
		return nil, fmt.Errorf("unsupported url for Docker Hub: %s", image.Registry)
	}

	url := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags?page_size=200", dockerHubBaseURL, image.Namespace, image.Name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request to Docker Hub: %w", err)
	}

	// TODO: Configure Auth for Docker Hub

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to Docker Hub: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK status code from Docker Hub (%d): %s (Body: %s)", resp.StatusCode, resp.Status, string(respBody))
	}

	var tagsResponse dockerHubTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Docker Hub API response: %w", err)
	}

	tags := make([]string, len(tagsResponse.Results))
	for i, tag := range tagsResponse.Results {
		tags[i] = tag.Name
	}

	return tags, nil
}
