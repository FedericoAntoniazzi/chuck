package containerruntime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerClientInterface defines the methods we use from the Docker client
type DockerClientInterface interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
}

type DockerRuntime struct {
	client DockerClientInterface
}

// NewDockerRuntime creates a new Docker runtime with a real client
func NewDockerRuntime() *DockerRuntime {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil
	}
	return &DockerRuntime{client: cli}
}

// NewDockerRuntime creates a new DockerRuntime instance
func NewDockerRuntimeWithClient(client DockerClientInterface) *DockerRuntime {
	return &DockerRuntime{client: client}
}

func (d *DockerRuntime) ListContainers() ([]Container, error) {
	containers, err := d.client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []Container
	for _, c := range containers {
		result = append(result, Container{
			ID:          c.ID,
			Name:        c.Names[0],
			Image:       c.Image,
			RuntimeType: "docker",
		})
	}

	return result, nil
}
