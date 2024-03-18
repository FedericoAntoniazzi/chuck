package sources

import "strings"

type Image struct {
	Name    string
	Version string
}

type ImageSource interface {
	ListImages() []Image
}

type FakeImageSource struct{}

func (fis *FakeImageSource) ListAllImages() []Image {
	return fis.ListImages("")
}

func (fis *FakeImageSource) ListImages(filter string) []Image {
	images := []Image{
		{
			Name: "nginx",
			Version: "1.18",
		},
		{
			Name: "pause",
			Version: "1",
		},
	}

	result := []Image{}

	for _, image := range images {
		if strings.HasPrefix(image.Name, filter) {
			result = append(result, image)
		}
	}

	return result
}

