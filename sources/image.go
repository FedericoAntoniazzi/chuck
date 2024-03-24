package sources

import "strings"

type Image struct {
	Name    string
	Version string
}

func NewImage(imageRef string) Image {
	imageDetail := strings.Split(imageRef, ":")

	image := Image{
		Name:    imageDetail[0],
		Version: "latest",
	}

	if len(imageDetail) > 1 {
		image.Version = imageDetail[1]
	}

	return image
}
