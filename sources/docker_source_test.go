package sources

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerMock struct {
	client.APIClient
}

func mockedContainersInfo() []types.Container {
	return []types.Container{
		{
			ID:      "12345",
			Names:   []string{"container1"},
			Image:   "image:v1.2.3",
			ImageID: "123456789",
			Labels:  map[string]string{"chuck.enable": "true"},
		},
		{
			ID:      "12345",
			Names:   []string{"container1"},
			Image:   "image:v1.2.3",
			ImageID: "123456789",
			Labels:  map[string]string{"chuck.enable": "false"},
		},
	}
}

func (d DockerMock) ContainerList(ctx context.Context, opts container.ListOptions) ([]types.Container, error) {
	mockedData := mockedContainersInfo()
	result := []types.Container{}

	if opts.All {
		// "only running containers must be present")
		return result, errors.New("expected returing only running containers")
	}

	for _, d := range mockedData {
		if len(opts.Filters.Get("label")) == 1 && opts.Filters.Get("label")[0] == "chuck.enable=true" {
			if d.Labels["chuck.enable"] == "true" {
				result = append(result, d)
			}
		}
	}

	return result, nil
}

func TestMockDockerImageListContainersWithLabels(t *testing.T) {
	mockedClient := DockerMock{}
	mockedImageSource := NewDockerImageSource(mockedClient)

	images, err := mockedImageSource.ListRunningContainerImages()
	if err != nil {
		t.Error(err)
	}

	t.Run("each image must come from a container with the label chuck.enable=true", func(t *testing.T) {
		if len(mockedContainersInfo()) == len(images) {
			t.Errorf("expected %d results, got %d", len(mockedContainersInfo()), len(images))
		}
	})
}
