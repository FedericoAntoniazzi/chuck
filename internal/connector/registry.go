package connector

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type RegistryConnector struct{}

// ListAllTags queries the registry APIs and returns all tags associated to the given image
func (rc RegistryConnector) ListAllTags(image string) ([]string, error) {
	repo, err := name.NewRepository(image)
	if err != nil {
		return nil, err
	}

	puller, err := remote.NewPuller()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	lister, err := puller.Lister(ctx, repo)
	if err != nil {
		return nil, err
	}

	tags := []string{}
	for lister.HasNext() {
		remoteTags, err := lister.Next(ctx)
		if err != nil {
			return nil, err
		}
		tags = append(tags, remoteTags.Tags...)
	}

	return tags, nil
}
