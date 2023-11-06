package chuck

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/FedericoAntoniazzi/chuck/internal/client"
	"github.com/FedericoAntoniazzi/chuck/internal/registry"
)

type ContainerEngine interface {
	ListContainers() ([]client.Container, error)
}

func Job(ce ContainerEngine) error {
	// Get running containers with the label for chuck
	listedContainers, err := ce.ListContainers()
	if err != nil {
		return err
	}

	// Map container structs
	containers := make([]Container, len(listedContainers))
	for i, lc := range listedContainers {
		containers[i] = mapContainer(lc)
	}

	// Retrieve updates for running containers
	for i, cnt := range containers {
		image := fmt.Sprintf("%s/%s", cnt.Image.Registry, cnt.Image.Name)
		remoteTags, err := registry.ListRemoteTags(image)
		if err != nil {
			slog.Error("Error listing remote tags", "error", err)
		}

		containers[i].ImageUpdates = remoteTags
	}

	for _, c := range containers {
		fmt.Printf("Container %s (%s) can be updated: %v\n", c.Name, c.Image, c.ImageUpdates)
	}

	return nil
}

func mapContainer(clientContainer client.Container) Container {
	return Container{
		Name:  clientContainer.Name,
		Image: parseImage(clientContainer.Image),
	}
}

func parseImage(imageName string) Image {
	image := Image{
		Registry: "docker.io",
		Tag:      "latest",
	}

	splitNameForTag := strings.Split(imageName, ":")
	if len(splitNameForTag) > 0 {
		image.Tag = splitNameForTag[1]
	}

	splitNameForRegistry := strings.Split(imageName, "/")
	if len(splitNameForRegistry) == 3 {
		image.Registry = splitNameForRegistry[0]

		splitNameForName := strings.Join(splitNameForRegistry[1:3], "/")
		image.Name = strings.Split(splitNameForName, ":")[0]
	} else {
		image.Name = splitNameForTag[0]
	}

	slog.Info("parsed image", "registry", image.Registry, "name", image.Name, "tag", image.Tag)

	return image
}
