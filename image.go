package main

import "strings"

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

	// Extract the tag if present, else use latest
	nameAndTag := strings.Split(fullImageName, ":")
	safeName := nameAndTag[0]

	// If the image name contains ":" for the port of the registry
	if len(nameAndTag) == 3 {
		image.Tag = nameAndTag[2]
		safeName = nameAndTag[0] + ":" + nameAndTag[1]

	} else if len(nameAndTag) == 2 {

		// If it is Registry:port/Image
		if strings.ContainsRune(nameAndTag[1], '/') {
			safeName = fullImageName
			image.Tag = "latest"

			// If it is Image:Tag
		} else {
			image.Tag = nameAndTag[1]
		}

	} else {
		image.Tag = "latest"
	}

	// Extract registry and namespace from image name
	imageNameParts := strings.Split(safeName, "/")

	// If just the name is present, use docker's default values
	if len(imageNameParts) == 1 {
		image.Registry = "docker.io"
		image.Namespace = "library"
		image.Name = imageNameParts[0]

		// If the image is composed by namespace/image
	} else if len(imageNameParts) == 2 {
		image.Registry = "docker.io"
		image.Namespace = imageNameParts[0]
		image.Name = imageNameParts[1]

		// If the image is composed by registry/namespace/image
	} else if len(imageNameParts) == 3 {
		image.Registry = imageNameParts[0]
		image.Namespace = imageNameParts[1]
		image.Name = imageNameParts[2]
	}

	return image, nil
}
