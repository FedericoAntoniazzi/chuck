package containerruntime

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDockerClient struct {
	mock.Mock
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.Container), args.Error(1)
}

func TestDockerRuntimeListContainersUnit(t *testing.T) {
	mockClient := new(MockDockerClient)

	mockContainers := []types.Container{
		{
			ID:    "unit-test-container-1",
			Names: []string{"/test-container-1"},
			Image: "nginx:latest",
		},
	}

	mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(mockContainers, nil)

	dockerRuntime := &DockerRuntime{client: mockClient}
	containers, err := dockerRuntime.ListContainers()

	assert.NoError(t, err)
	assert.Len(t, containers, 1)
	assert.Equal(t, "unit-test-container-1", containers[0].ID)
	assert.Equal(t, "nginx:latest", containers[0].Image)
	assert.Equal(t, "docker", containers[0].RuntimeType)

	mockClient.AssertExpectations(t)
}
