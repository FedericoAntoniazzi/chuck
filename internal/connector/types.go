package connector

import (
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

const ChuckEnableLabel = "chuck.enable"

// defaultDockerClient returns a docker client with automatic configuration from env
func defaultDockerClient() (*docker.Client, error) {
	return docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
}

// Container describes common parameters of a container
type Container struct {
	Image  string
	Labels map[string]string
	Name   string
}

// newFromDockerContainer returns a Container instance from a docker Container type
func newFromDockerContainer(dockerContainer types.Container) Container {
	return Container{
		Name:   dockerContainer.Names[0],
		Image:  dockerContainer.Image,
		Labels: dockerContainer.Labels,
	}
}

// nameSortedContainers describe a list of containers sorted by name.
// It's there for testing purposes
type nameSortedContainers []Container

func (sc nameSortedContainers) Len() int {
	return len(sc)
}

func (sc nameSortedContainers) Less(i, j int) bool {
	return sc[i].Name < sc[j].Name
}

func (sc nameSortedContainers) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}
