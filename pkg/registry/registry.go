package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
	"github.com/hashicorp/go-version"
)

func parseImageTag(image string) string {
	split := strings.Split(image, ":")
	if len(split) > 1 {
		return split[1]
	}
	return "default"
}

func ListNewerTags(image string) ([]string, error) {
	sys := &types.SystemContext{
		DockerInsecureSkipTLSVerify: types.OptionalBoolFalse,
	}

	ref, err := docker.ParseReference(fmt.Sprintf("//%s", image))
	if err != nil {
		return nil, err
	}

	tags, err := docker.GetRepositoryTags(context.Background(), sys, ref)
	if err != nil {
		return nil, err
	}

	newerTags := make([]string, 0)
	imageVersion := parseImageTag(image)
	if imageVersion == "latest" {
		return newerTags, errors.New("latest tag is not supported")
	}

	refVersion, err := version.NewVersion(imageVersion)
	if err != nil {
		return newerTags, err
	}

	for _, tag := range tags {
		v, err := version.NewVersion(tag)
		// Ignore non-semver tags
		if err != nil {
			continue
		}

		// Ignore variant tags (e.g. -alpine, -apache, -0)
		if strings.Contains(tag, "-") || strings.Contains(tag, "+") {
			continue
		}

		if refVersion.LessThan(v) {
			newerTags = append(newerTags, tag)
		}
	}

	return newerTags, nil
}
