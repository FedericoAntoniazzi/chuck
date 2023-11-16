package connector

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type RegistryConnector struct{}

// ListNewerSemverTags queries the registry APIs and returns newers tags matching the semver format
func (rc RegistryConnector) ListNewerSemverTags(image, tag string) ([]string, error) {
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

		newerConstraint, err := semver.NewConstraint("> " + tag)
		if err != nil {
			fmt.Println("WARN", "skipping image with invalid semver version.", "image="+image, "tag="+tag, fmt.Sprintf("error=\"%s\"", err))
			break
		}

		for _, rt := range remoteTags.Tags {
			v, err := semver.NewVersion(rt)
			if err != nil {
				// fmt.Println("DEBUG", "skipping non-semver tag", "image="+image, "tag="+rt)
				continue
			}

			if newerConstraint.Check(v) {
				tags = append(tags, rt)
			}
		}
	}

	return tags, nil
}
