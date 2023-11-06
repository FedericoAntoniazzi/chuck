package registry

import (
	"context"
	"log/slog"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ListRemoteTags lists all tags from the container registry that hosts the given image.
// The image must be in the format <registry>/<user/group>/<name> without the tag.
func ListRemoteTags(image string) ([]string, error) {
	repo, err := name.NewRepository(image)
	if err != nil {
		slog.Error("error parsing repo", "image", image, "error", err)
	}

	puller, err := remote.NewPuller()
	if err != nil {
		slog.Error("error creating remote puller", "error", err)
	}

	ctx := context.Background()
	lister, err := puller.Lister(ctx, repo)
	if err != nil {
		slog.Error("error listing all tags", "image", image, "error", err)
	}

	remoteTags := []string{}
	for lister.HasNext() {
		tags, err := lister.Next(ctx)
		if err != nil {
			slog.Error("error listing tags in page", "image", image, "error", err)
		}

		remoteTags = append(remoteTags, tags.Tags...)
	}

	return remoteTags, nil
}
