package connector

import (
	"context"
	"fmt"

	"github.com/FedericoAntoniazzi/chuck/internal/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

type DockerConnector struct {
	client *docker.Client
}

// NewDockerConnector returns the connector for a docker engine instance
func NewDockerConnector() (*DockerConnector, error) {
	dockerClient, err := defaultDockerClient()
	if err != nil {
		return nil, err
	}

	return &DockerConnector{
		client: dockerClient,
	}, nil
}

// listContainersWithLabel queries the docker APIs to fetch running containers with the given label
func (dc *DockerConnector) listContainersWithLabel(ctx context.Context, label string) ([]types.Container, error) {
	containerListOpts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("%s=true", ChuckEnableLabel)),
		),
	}
	return dc.client.ContainerList(ctx, containerListOpts)
}

func (dc *DockerConnector) ListContainers(ctx context.Context) ([]models.Container, error) {
	labeledContainers, err := dc.listContainersWithLabel(ctx, ChuckEnableLabel)
	if err != nil {
		return nil, err
	}

	containers := make([]models.Container, len(labeledContainers))
	for i, lc := range labeledContainers {
		containers[i] = newFromDockerContainer(lc)
	}
	return containers, err
}
