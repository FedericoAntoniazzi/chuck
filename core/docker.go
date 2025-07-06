package core

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

// getRunningContainerImages connects to the Docker daemon and returns a list of running containers.
func GetRunningContainerImages(ctx context.Context, log *zap.SugaredLogger) ([]container.Summary, error) {
	log.Info("Connecting to Docker daemon")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}
	defer cli.Close()
	log.Info("Successfully connected to Docker daemon")

	log.Info("Listing running containers")
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %w", err)
	}
	return containers, nil
}
