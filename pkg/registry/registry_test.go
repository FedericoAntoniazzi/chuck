package registry

import (
	"testing"
)

func TestRegistry_parseImageTag(t *testing.T) {
	inputs := []struct {
		inputImage  string
		expectedTag string
	}{
		{
			inputImage:  "nginx:1.21",
			expectedTag: "1.21",
		},
		{
			inputImage:  "nginx:v1.21",
			expectedTag: "v1.21",
		},
		{
			inputImage:  "nginx:v1.21-alpine",
			expectedTag: "v1.21-alpine",
		},
	}

	for _, input := range inputs {
		tag := parseImageTag(input.inputImage)
		if tag != input.expectedTag {
			t.Fatalf("Expected %s, got %s", input.expectedTag, tag)
		}
	}
}

func TestRegistry_ListNewerTags(t *testing.T) {
	tags, err := ListNewerTags("nginx:1.21")
	if err != nil {
		t.Fatal(err)
	}

	if len(tags) == 0 {
		t.Fatal("Result cannot be 0 length")
	}
}
