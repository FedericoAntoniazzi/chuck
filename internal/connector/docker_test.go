package connector

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/FedericoAntoniazzi/chuck/internal/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
)

func runTestingContainers(ctx context.Context) ([]testcontainers.Container, []models.Container) {
	defaultLabels := map[string]string{
		ChuckEnableLabel: "true",
	}
	containers := []models.Container{
		{
			Name:   "nginx_125",
			Image:  "nginx:1.25",
			Labels: defaultLabels,
		},
		{
			Name:   "library_nginx_125",
			Image:  "library/nginx:1.25",
			Labels: defaultLabels,
		},
		{
			Name:   "dockerio_library_nginx_125",
			Image:  "docker.io/library/nginx:1.25",
			Labels: defaultLabels,
		},
	}

	dockerClient, err := defaultDockerClient()
	if err != nil {
		panic(err)
	}

	testContainers := make([]testcontainers.Container, len(containers))
	for i, cntr := range containers {
		req := testcontainers.ContainerRequest{
			Name:   cntr.Name,
			Image:  cntr.Image,
			Labels: cntr.Labels,
		}

		tc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			panic(err)
		}
		testContainers[i] = tc
	}

	// Get updated info from running containers
	// Workaround until go-containers allows retrieving the image from test containers.
	for i, tc := range testContainers {
		cont, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
			Filters: filters.NewArgs(
				filters.Arg("id", tc.GetContainerID()),
			),
		})
		if err != nil {
			panic(err)
		}
		containers[i] = newFromDockerContainer(cont[0])
	}

	sort.Sort(nameSortedContainers(containers))

	return testContainers, containers
}

func destroyTestingContainers(ctx context.Context, containers []testcontainers.Container) {
	for _, cntr := range containers {
		if err := cntr.Terminate(ctx); err != nil {
			panic(err)
		}
	}
}

// TestValidateDefaultDockerClient verifies the docker environment
// The default client should never return an error in a valid environment
func TestValidateDockerTestingEnvironment(t *testing.T) {
	cli, err := defaultDockerClient()
	if err != nil {
		t.Errorf("error creating default docker instance: %s", err)
	}

	ctx := context.Background()
	_, err = cli.ServerVersion(ctx)
	if err != nil {
		t.Errorf("error in docker environment: %s", err)
	}

}

// TestNewDockerConnector verifies the creation of the default Docker Connector implementation
func TestNewDockerConnector(t *testing.T) {
	got, gotError := NewDockerConnector()

	wantClient, _ := defaultDockerClient()
	want := &DockerConnector{
		client: wantClient,
	}
	var wantError error = nil

	if gotError != wantError {
		t.Errorf("want error %s, got error %s", wantError, gotError)
	}

	if reflect.DeepEqual(want, got) {
		t.Errorf("want %v, want %v", want, got)
	}
}

func TestListContainersWithLabel(t *testing.T) {
	dc, _ := NewDockerConnector()
	ctx := context.Background()

	tcs, want := runTestingContainers(ctx)
	defer destroyTestingContainers(ctx, tcs)

	got, err := dc.ListContainers(ctx)
	if err != nil {
		t.Errorf("error listing containers: %s", err)
	}

	if len(want) != len(got) {
		t.Fatalf("want %d containers, got %d", len(want), len(got))
	}

	for i := 0; i < len(want); i++ {
		if !cmp.Equal(want[i], got[i]) {
			t.Error(cmp.Diff(want[i], got[i]))
		}
	}
}
