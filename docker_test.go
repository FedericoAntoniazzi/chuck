package main

import (
	"context"
	"io"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// TestGetRunningContainerImages is an integration test for the getRunningContainerImages function.
// It requires a running Docker daemon and will attempt to create/remove containers.
func TestGetRunningContainerImages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if Docker daemon is accessible
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Docker daemon not reachable: %v", err)
	}
	defer cli.Close()

	// Check if Docker daemon is responsive
	_, err = cli.Ping(ctx)
	if err != nil {
		t.Skipf("Docker daemon not responsive: %v", err)
	}

	// Use dedicated output for tests
	// originalLogOutput := log.Writer()
	// log.SetOutput(os.Stderr)

	createAndStartContainer := func(t *testing.T, ctx context.Context, cli *client.Client, imageName, containerName string) string {
		t.Helper()

		// Check if image is already present in the host
		images, err := cli.ImageList(ctx, image.ListOptions{})
		if err != nil {
			t.Fatalf("Failed to list images in the host: %v", err)
		}

		imagePresent := false
		for _, image := range images {
			imagePresent = slices.Contains(image.RepoTags, imageName)
			if imagePresent {
				break
			}
		}

		if !imagePresent {
			// Ensure image is pulled
			log.Printf("Pulling image %s", imageName)
			pullResponse, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				t.Fatalf("Failed to pull image %s: %v", imageName, err)
			}
			_, err = io.ReadAll(pullResponse)
			if err != nil {
				t.Fatalf("Failed to read pull image result: %v", err)
			}
			log.Printf("Image pulled")
		}

		// Create container
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Cmd:   []string{"sleep", "infinity"}, // Keep container running
		}, nil, nil, nil, containerName)
		if err != nil {
			t.Fatalf("Failed to create container %s with image %s, %v", containerName, imageName, err)
		}

		// Start container
		err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
		if err != nil {
			t.Fatalf("Failed to start container %s: %v", resp.ID, err)
		}

		// Defer cleanup after the test
		t.Cleanup(func() {
			t.Logf("Cleaning up container %s (%s)", containerName, resp.ID[:12])
			// Not used in favor of ContainerKill which is faster
			// Keeping there in case some ContainerKill will have undesired side effects
			//
			// if err := cli.ContainerStop(ctx, resp.ID, container.StopOptions{}); err != nil {
			// 	t.Logf("Warning: Failed to stop container %s during cleanup: %v", resp.ID, err)
			// }
			if err := cli.ContainerKill(ctx, resp.ID, "SIGKILL"); err != nil {
				t.Logf("Warning: Failed to kill container %s during cleanup: %v", resp.ID, err)
			}
			if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{}); err != nil {
				t.Logf("Warning: Failed to remove container %s during cleanup: %v", resp.ID, err)
			}
		})

		// Wait some time for the container to start
		time.Sleep(500 * time.Millisecond)

		return resp.ID
	}

	// Test without running containers
	// The user is responsible to run this code on a clean docker host with no running containers (yet)
	t.Run("NoRunningContainers", func(t *testing.T) {
		containers, err := getRunningContainerImages(ctx)
		if err != nil {
			t.Fatalf("Expected no error when listing 0 containers, got: %v", err)
		}
		if len(containers) != 0 {
			t.Errorf("Expected 0 running containers, got %d", len(containers))
		}
	})

	// Test with 1 running container
	t.Run("WithOneRunningContainer", func(t *testing.T) {
		containerName := "chuck-tests-0"
		imageName := "nginx:1.25"
		_ = createAndStartContainer(t, ctx, cli, imageName, containerName)

		containers, err := getRunningContainerImages(ctx)
		if err != nil {
			t.Fatalf("Expected no error when listing containers, got %v", err)
		}

		if len(containers) != 1 {
			t.Errorf("Expected 1 running container, got %d", len(containers))
		}
	})

	// Test with 2 running container
	t.Run("WithManyRunningContainers", func(t *testing.T) {
		testContainers := []struct {
			container string
			image     string
		}{
			{
				container: "chuck-tests-1",
				image:     "nginx:1.25",
			},
			{
				container: "chuck-tests-2",
				image:     "nginx:1.25",
			},
			{
				container: "chuck-tests-3",
				image:     "nginx:1.25",
			},
		}

		for _, testContainer := range testContainers {
			_ = createAndStartContainer(t, ctx, cli, testContainer.image, testContainer.container)
		}

		containers, err := getRunningContainerImages(ctx)
		if err != nil {
			t.Fatalf("Expected no error when listing containers, got %v", err)
		}

		if len(containers) != len(testContainers) {
			t.Errorf("Expected %d running container, got %d", len(testContainers), len(containers))
		}
	})

	// Test invalid docker client configuration
	t.Run("DockerClientInvalidConfig", func(t *testing.T) {
		// Temporarily set an invalid DOCKER_HOST reference
		originalDockerHost := os.Getenv("DOCKER_HOST")
		_ = os.Setenv("DOCKER_HOST", "tcp://127.0.0.1:63001")
		defer os.Setenv("DOCKER_HOST", originalDockerHost)

		_, err := getRunningContainerImages(ctx)
		if err == nil {
			t.Error("Expected an error in case of docker client creation failure, but got none")
		}
	})
}
