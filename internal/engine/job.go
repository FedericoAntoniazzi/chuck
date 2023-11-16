package engine

import (
	"context"
	"strings"

	"github.com/FedericoAntoniazzi/chuck/internal/connector"
	"github.com/FedericoAntoniazzi/chuck/internal/models"
)

type NotificationEndpoint interface {
	Send(title, message string) error
}

type Job struct {
	ContainerEngine *connector.DockerConnector
	RegistryRef     *connector.RegistryConnector
	NotifyEndpoint  NotificationEndpoint
}

func NewJob() (*Job, error) {
	ce, err := connector.NewDockerConnector()
	if err != nil {
		return nil, err
	}

	return &Job{
		ContainerEngine: ce,
		RegistryRef:     &connector.RegistryConnector{},
		NotifyEndpoint:  connector.ConsoleConnector{},
	}, nil
}

func (job *Job) Run() error {
	ctx := context.Background()

	// Get all containers that should be checked for updates
	containers, err := job.ContainerEngine.ListContainers(ctx)
	if err != nil {
		return err
	}

	imgUpdates := make([]models.ImageUpdate, len(containers))

	// For each image, get the respective tags that match the semver format
	for i, cntr := range containers {
		if cntr.ImageTag() == "latest" {
			continue
		}

		allTags, err := job.RegistryRef.ListAllTags(cntr.ImageName())
		if err != nil {
			return err
		}

		imgUpdates[i] = models.ImageUpdate{
			Name:        cntr.ImageName(),
			CurrentTag:  cntr.ImageTag(),
			UpdatedTags: allTags,
		}
	}

	for _, imgUpd := range imgUpdates {
		_ = job.NotifyEndpoint.Send("Updates", strings.Join(imgUpd.UpdatedTags, ", "))
	}

	return nil
}
