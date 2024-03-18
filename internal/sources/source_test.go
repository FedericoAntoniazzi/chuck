package sources

import (
	"reflect"
	"testing"
)

func TestListImages(t *testing.T) {
	fis := FakeImageSource{}

	got := fis.ListImages("nginx")
	want := []Image{
		{
			Name: "nginx",
			Version: "1.18",
		},
	}

	if len(got) != len(want) {
		t.Errorf("got %d, want %d", len(got), len(want))
	}

	if got[0].Name != want[0].Name {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestListAllImages(t *testing.T) {
	fis := FakeImageSource{}

	got := fis.ListAllImages()
	want := []Image{
		{
			Name: "nginx",
			Version: "1.18",
		},
		{
			Name: "pause",
			Version: "1",
		},
	}

	if len(got) != len(want) {
		t.Errorf("got %d, want %d", len(got), len(want))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
