package engine

import (
	"context"
	"fmt"
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

	imgUpdates := []models.ImageUpdate{}

	// For each image, get the respective tags that match the semver format
	for _, cntr := range containers {
		if cntr.ImageTag() == "latest" {
			continue
		}

		allTags, err := job.RegistryRef.ListNewerSemverTags(cntr.ImageName(), cntr.ImageTag())
		if err != nil {
			return err
		}

		// Ignore images with no updates
		if len(allTags) == 0 {
			continue
		}

		imgUpdates = append(imgUpdates, models.ImageUpdate{
			Name:        cntr.ImageName(),
			CurrentTag:  cntr.ImageTag(),
			UpdatedTags: allTags,
		})
	}

	for _, imgUpd := range imgUpdates {
		title := fmt.Sprintf("Updates for %s:%s", imgUpd.Name, imgUpd.CurrentTag)
		message := strings.Join(imgUpd.UpdatedTags, ", ")
		_ = job.NotifyEndpoint.Send(title, message)
	}

	return nil
}
