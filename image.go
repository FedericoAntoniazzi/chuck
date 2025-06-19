package main

import (
	"fmt"
	"path"

	"github.com/distribution/reference"
)

// Image describes a container image.
// A container image is composed by Registry/Namespace/Name:Tag (e.g docker.io/library/nginx:1.25)
type Image struct {
	Raw       string // Unparsed image reference
	Registry  string // Registry's URL
	Namespace string // Image namespace (In case of Docker Hub may be library, or the username)
	Name      string // Name of the image (e.g nginx, redis)
	Tag       string // Tag assigned to the image (e.g. latest, 1.15, v2.1.0)
}

// ParseImageName parses a full image string (e.g., "nginx:latest", "myregistry.com/user/image:tag")
// into its components (Registry, Namespace, Name, Tag).
// It handles default Docker Hub implicit values (docker.io, library).
func ParseImageName(fullImageName string) (Image, error) {
	image := Image{
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
