package client

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const ChuckEnableLabel = "chuck.enable=true"

// DockerClient describes a client for the Docker engine
type DockerClient struct {
	client *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	slog.Debug("Creating docker client")

	// Initialise the default client from env vars
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		slog.Error("error creating docker client", "error", err)
	}

	return &DockerClient{
		client: cli,
	}, err
}

// ListLabeledContainers retrieve all the containers with the label chuck.enabled.
func (dc *DockerClient) listLabeledContainers() ([]types.Container, error) {
	ctx := context.Background()
	return dc.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", ChuckEnableLabel),
		),
	})
}

func (dc *DockerClient) ListContainers() ([]Container, error) {
	labeledContainers, err := dc.listLabeledContainers()
	if err != nil {
		slog.Error("error listing containers.", "error", err)
		return nil, err
	}

	slog.Info("fetched containers", "count", len(labeledContainers))

	containers := make([]Container, len(labeledContainers))
	for i, labeledContainer := range labeledContainers {
		if len(labeledContainer.Names) == 0 {
			slog.Warn("found container with no name. skipping.", "containerId", labeledContainer.ID)
			continue
		}

		containers[i] = Container{
			Name:  labeledContainer.Names[0],
			Image: labeledContainer.Image,
		}
	}

	return containers, nil
}
