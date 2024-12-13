package imagecheck

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/FedericoAntoniazzi/chuck/internal/containerruntime"
	"github.com/Masterminds/semver/v3"
	"github.com/distribution/reference"
)

type UpdateResult struct {
	Container      containerruntime.Container
	HasUpdate      bool
	CurrentVersion string
	LatestVersion  string
}

// RegistryClient interface for fetching tags
type RegistryClient interface {
	ListTags(ctx context.Context, repository string) ([]string, error)
}

type UpdateChecker struct {
	registryClient RegistryClient
}

func NewUpdateChecker(client RegistryClient) *UpdateChecker {
	return &UpdateChecker{
		registryClient: client,
	}
}

// parseImageReference parses the image name into components
// It also ensures the image name contains the default `latest` tag if not set
func parseImageReference(imageRef string) (reference.Named, error) {
	parsedRef, err := reference.ParseNormalizedNamed(imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image ref: %v", err)
	}

	parsedRef = reference.TagNameOnly(parsedRef)

	return parsedRef, nil
}

func extractImageTag(imageName string) string {
	// Split the image name by ':'
	parts := strings.Split(imageName, ":")

	// If no tag is specified, return 'latest' as default
	if len(parts) == 1 {
		return "latest"
	}

	// Return the last part as the tag
	return parts[len(parts)-1]
}

func extractImageRepository(imageName string) string {
	// Split the image name by ':'
	parts := strings.Split(imageName, ":")

	// Return the first part as the repository name
	return parts[0]
}

func (uc *UpdateChecker) GetRemoteLatestTag(ctx context.Context, imageRef string) (string, error) {
	repository := extractImageRepository(imageRef)

	// Retrieve image tags from registry
	tags, err := uc.registryClient.ListTags(ctx, repository)
	if err != nil {
		return "", fmt.Errorf("failed to list tags: %v", err)
	}

	// sort tags and return the last (most recent) tag
	sortedTags, err := sortSemanticVersionTags(tags)
	if err != nil {
		return "", fmt.Errorf("failed to sort tags: %v", err)
	}

	return sortedTags[len(tags)-1], nil
}

func sortSemanticVersionTags(rawTags []string) ([]string, error) {
	// Convert tags to semantic versions
	sortedVersions := make([]*semver.Version, len(rawTags))
	for i, raw := range rawTags {
		ver, err := semver.NewVersion(raw)
		if err != nil {
			return []string{}, fmt.Errorf("error sorting versions: %v", err)
		}

		sortedVersions[i] = ver
	}

	// Sort semantic versions
	sort.Sort(semver.Collection(sortedVersions))

	// Convert back tags from semantic versions
	resultTags := make([]string, len(rawTags))
	for i, version := range sortedVersions {
		resultTags[i] = version.Original()
	}

	return resultTags, nil
}

func (uc *UpdateChecker) CheckImageUpdate(container containerruntime.Container) (UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	currentTag := extractImageTag(container.Image)

	latestTag, err := uc.GetRemoteLatestTag(ctx, container.Image)
	if err != nil {
		return UpdateResult{}, fmt.Errorf("failed to retrieve latest tag from registry: %v", err)
	}

	hasUpdate := currentTag != latestTag

	return UpdateResult{
		Container:      container,
		HasUpdate:      hasUpdate,
		CurrentVersion: currentTag,
		LatestVersion:  latestTag,
	}, nil
}
