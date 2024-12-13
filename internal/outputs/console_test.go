package outputs

import (
	"bytes"
	"testing"

	"github.com/FedericoAntoniazzi/chuck/internal/containerruntime"
	"github.com/FedericoAntoniazzi/chuck/internal/imagecheck"
	"github.com/stretchr/testify/assert"
)

func TestConsoleOutput_Submit(t *testing.T) {
	testCases := []struct {
		name           string
		inputResults   []imagecheck.UpdateResult
		expectedOutput string
	}{
		{
			name: "Single container with update",
			inputResults: []imagecheck.UpdateResult{
				{
					Container: containerruntime.Container{
						Name:  "test-container",
						Image: "nginx",
					},
					HasUpdate:     true,
					LatestVersion: "1.25.0",
				},
			},
			expectedOutput: "Container test-container can be updated to nginx:1.25.0\n",
		},
		{
			name: "Multiple containers with mixed update status",
			inputResults: []imagecheck.UpdateResult{
				{
					Container: containerruntime.Container{
						Name:  "updated-container",
						Image: "nginx",
					},
					HasUpdate:     true,
					LatestVersion: "1.25.0",
				},
				{
					Container: containerruntime.Container{
						Name:  "up-to-date-container",
						Image: "alpine",
					},
					HasUpdate: false,
				},
			},
			expectedOutput: "Container updated-container can be updated to nginx:1.25.0\n",
		},
		{
			name: "No containers with updates",
			inputResults: []imagecheck.UpdateResult{
				{
					Container: containerruntime.Container{
						Name:  "up-to-date-container",
						Image: "alpine",
					},
					HasUpdate: false,
				},
			},
			expectedOutput: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer

			// Create ConsoleOutput with the buffer
			consoleOutput := NewConsoleOutput(&buf)

			// Submit results
			consoleOutput.Submit(tc.inputResults)

			// Check the output
			assert.Equal(t, tc.expectedOutput, buf.String())
		})
	}
}
