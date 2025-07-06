package core

import (
	"fmt"
	"log"
	"slices"
	"sort"

	"github.com/Masterminds/semver/v3"
)

// FindLatestUpdate finds the latest semver-compatible update for a given current tag
// from a list of available tags. It returns the latest found tag and a boolean indicating if an update was found.
func FindLatestUpdate(currentTag string, availableTags []string) (string, bool, error) {
	currentVersion, err := semver.NewVersion(currentTag)
	if err != nil {
		// If the current tag is not a valid semver
		log.Printf("WARNING: Current tag '%s' is not a valid semantic version. Skipping strict SemVer comparison.", currentTag)
		// Fallback: If current tag is 'latest' and it's in the list, no update.
		if currentTag == "latest" {
			if slices.Contains(availableTags, currentTag) {
				return "", false, nil // "latest" found and is current, no update.
			}
		}
		return "", false, fmt.Errorf("current tag '%s' is not a valid semver: %w", currentTag, err)
	}

	var versions []*semver.Version
	for _, tag := range availableTags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			// Ignore tags that are not valid semantic versions
			continue
		}
		versions = append(versions, v)
	}

	if len(versions) == 0 {
		return "", false, nil // No valid semver tags found to compare against
	}

	// Sort versions in descending order
	sort.Sort(sort.Reverse(semver.Collection(versions)))

	// The latest version is the first one after sorting
	latestAvailableVersion := versions[0]

	// Check if the latest available version is newer than the current version
	if latestAvailableVersion.GreaterThan(currentVersion) {
		return latestAvailableVersion.String(), true, nil // Update found
	}

	return "", false, nil // No update found or current is already the latest
}
