package connector

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
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

// ListAllSemverTags queries the registry APIs and returns all tags associated to the given image that match the semver format
func (rc RegistryConnector) ListAllSemverTags(image string) ([]string, error) {
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

	semverRegexp, _ := regexp.Compile(`[0-9]\.[0-9]+\.[0.9]+.*`)

	for lister.HasNext() {
		remoteTags, err := lister.Next(ctx)
		if err != nil {
			return nil, err
		}

		for _, tag := range remoteTags.Tags {
			if strings.HasPrefix(tag, "sha256-") {
				continue
			}
			if !semverRegexp.MatchString(tag) {
				continue
			}

			greaterSemverConstraint, err := semver.NewConstraint(fmt.Sprintf("> %s", tag))
			if err != nil {
				fmt.Printf("ERROR - WTF Constraint - %s", err)
			}
			v, err := semver.NewVersion(tag)
			if err != nil {
				fmt.Printf("ERROR - Skipping invalid tag - %s", err)
				continue
			}
			if ok, _ := greaterSemverConstraint.Validate(v); ok {
				tags = append(tags, tag)
			}
		}
	}

	fmt.Println(tags)

	return tags, nil
}
