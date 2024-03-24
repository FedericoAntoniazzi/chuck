package sources

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
)

type DockerClient interface {
	NewClientWithOpts(...dockerClient.Opt) (DockerClient, error)
	ContainerList(context.Context, container.ListOptions) ([]types.Container, error)
}

type DockerImageSource struct {
	client dockerClient.APIClient
}

// NewDockerClientFromEnv creates a new client for Docker APIs using the default Docker Environment variables
func NewDockerClientFromEnv() (*dockerClient.Client, error) {
	return dockerClient.NewClientWithOpts(dockerClient.FromEnv)
}

// NewDockerImageSource returns the client for retrieving running images from containers managed by docker
func NewDockerImageSource(client dockerClient.APIClient) *DockerImageSource {
	return &DockerImageSource{
		client: client,
	}
}

// listContainersWithLabel list all running containers with the given label
func (dis DockerImageSource) listContainersWithLabel(ctx context.Context, label string) ([]types.Container, error) {
	containerListOpts := container.ListOptions{
		All: false,
		Filters: filters.NewArgs(
			filters.Arg("label", label),
		),
	}
	return dis.client.ContainerList(ctx, containerListOpts)
}

// ListRunningContainerImages queries all images from running containers
func (dis DockerImageSource) ListRunningContainerImages() ([]Image, error) {
	ctx := context.Background()

	containers, err := dis.listContainersWithLabel(ctx, "chuck.enable=true")
	if err != nil {
		return nil, err
	}

	containerImages := make([]Image, len(containers))
	for i, c := range containers {
		image := NewImage(c.Image)
		containerImages[i] = image
	}

	return containerImages, nil
}
