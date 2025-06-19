package main

import (
	"fmt"
	"path"

	"github.com/FedericoAntoniazzi/chuck/types"
	"github.com/distribution/reference"
)

// ParseImageName parses a full image string (e.g., "nginx:latest", "myregistry.com/user/image:tag")
// into its components (Registry, Namespace, Name, Tag).
// It handles default Docker Hub implicit values (docker.io, library).
func ParseImageName(fullImageName string) (types.Image, error) {
	image := types.Image{
		Raw: fullImageName,
	}

	// Parse the full image name
	namedRef, err := reference.ParseNormalizedNamed(fullImageName)
	if err != nil {
		return image, fmt.Errorf("failed to parse image reference: %w", err)
	}

	// Extract registry
	image.Registry = reference.Domain(namedRef)

	// Extract image reference without the registry
	imagePath := reference.Path(namedRef)

	// Split path into image parts
	image.Name = path.Base(imagePath)
	image.Namespace = path.Dir(imagePath)

	// Extract tag
	taggedRef, isTagged := namedRef.(reference.Tagged)
	if isTagged {
		image.Tag = taggedRef.Tag()
	} else {
		image.Tag = "latest"
	}

	return image, nil
}
