package runtime

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerRuntime represents a docker runtime
type DockerRuntime struct {
	client *client.Client
}

// NewDockerRuntime creates a client for the Docker Runtime
// It supports configuration via default Docker environment variables
func NewDockerRuntime() (DockerRuntime, error) {
	cl, err := client.NewClientWithOpts(client.FromEnv)

	return DockerRuntime{
		client: cl,
	}, err
}

// ListRunningContainers returns information about running containers
func (d *DockerRuntime) ListRunningContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := d.client.ContainerList(ctx, container.ListOptions{
		All: false,
	})
	if err != nil {
		return []ContainerInfo{}, err
	}

	result := make([]ContainerInfo, len(containers))
	for i, cnt := range containers {
		cntName := "unknown"
		if len(cnt.Names) > 0 {
			cntName = cnt.Names[0]
		}

		info := ContainerInfo{
			ID:    cnt.ID,
			Name:  strings.Replace(cntName, "/", "", -1),
			Image: cnt.Image,
		}
		result[i] = info
	}

	return result, nil
}
