package runtime

import (
	"context"
	"testing"
)

func TestDockerRuntime_NewDockerRuntime(t *testing.T) {
	runtime, err := NewDockerRuntime()
	if err != nil {
		t.Fatalf("Failed to create Docker runtime: %v", err)
	}
	if runtime.client == nil {
		t.Fatal("Docker runtime is nil")
	}
}

func TestDockerRuntime_ListRunningContainers(t *testing.T) {
	ctx := context.Background()
	runtime, _ := NewDockerRuntime()

	containers, err := runtime.ListRunningContainers(ctx)
	if err != nil {
		t.Fatalf("Unable to list containers %v", err)
	}

	if len(containers) == 0 {
		t.Fatalf("At least 1 container is required")
	}
}
