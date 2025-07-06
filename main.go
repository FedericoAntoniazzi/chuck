package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/FedericoAntoniazzi/chuck/registry/dockerhub"
	"github.com/FedericoAntoniazzi/chuck/types"
	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	defaultLogFileName = "chuck.log"
	defaultDBFileName  = "chuck.db"
)

// registryClient defines the capabilities of a generic client for container registries
type registryClient interface {
	// GetTags fetches all available tags for a given image from the registry
	GetTags(ctx context.Context, image types.Image) ([]string, error)
}

func main() {
	// --- CLI Flags Definition ---
	logToFile := flag.Bool("log-file", false, "Enable logging output to a file")
	dbPath := flag.String("db-path", defaultDBFileName, "Path to the SQLite database file")

	flag.Parse()

	// --- Logging Setup ---
	var logOutput *os.File
	if *logToFile {
		// For simplicity, we'll keep log file fixed relative to CWD for now
		logFilePath := defaultLogFileName
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			// If we can't open the log file, we should still try to log to stderr
			log.Printf("ERROR: Failed to open log file '%s': %v. Logging to stderr instead.", logFilePath, err)
			log.SetOutput(os.Stderr)
		} else {
			logOutput = file
			// Ensure the log file is closed when main exits
			defer logOutput.Close()
			log.SetOutput(logOutput)
		}
	} else {
		log.SetOutput(os.Stdout) // Default to console output
	}

	log.Println("Chuck: Starting container image update check...")

	// --- Database Path Handling ---
	// Resolve the absolute path for the database file
	resolvedDBPath, err := filepath.Abs(*dbPath)
	if err != nil {
		log.Fatalf("ERROR: Could not resolve absolute path for database file '%s': %v", *dbPath, err)
	}
	log.Printf("Using database file: %s", resolvedDBPath)

	// --- Core Logic Placeholder ---
	// This is where the core logic for Docker interaction, registry checks,
	// and SQLite operations will eventually go.

	// Create a background context for Docker API calls
	ctx := context.Background()

	registryClients := make(map[string]registryClient)
	registryClients["docker.io"] = dockerhub.NewClient()
	// Hint: registryClients["ghcr.io"] = github.NewClient()

	containers, err := getRunningContainerImages(ctx)
	if err != nil {
		log.Fatalf("Failed to get running containers: %v", err)
	}

	if len(containers) == 0 {
		log.Println("No running containers found")
		return
	}

	log.Printf("Found %d running containers: \n", len(containers))
	// Store the images which tags have already been queried
	uniqueImages := make(map[string][]string)

	var allUpdateStatuses []types.ImageUpdateStatus

	for _, cnt := range containers {
		containerName := ""

		if len(cnt.Names[0]) > 0 {
			containerName = strings.TrimPrefix(cnt.Names[0], "/")
		}

		log.Printf("DEBUG: processing container %s", containerName)

		status := types.ImageUpdateStatus{
			ContainerID:   cnt.ID,
			ContainerName: containerName,
			OriginalTag:   "latest",
			StatusMessage: "Processing",
		}

		image, err := ParseImageName(cnt.Image)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error parsing image name")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Printf("WARN: Cannot parse container image name: %v\n", err)
			continue
		}

		status.Image = image
		status.OriginalTag = image.Tag

		// Check if the registry is supported
		if _, ok := registryClients[image.Registry]; !ok {
			status.StatusMessage = fmt.Sprintln("Unsupported registry")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Printf("INFO: Registry %s is not yet supported. Skipping\n", image.Registry)
			continue
		}

		// Check if the tag is valid SemVer
		_, err = semver.NewVersion(image.Tag)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error parsing image tag")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Printf("INFO: Image %s is not a valid semver tag. Skipping\n", image.Tag)
			continue
		}

		// Check if image tags have already been fetched
		imageKey := fmt.Sprintf("%s/%s/%s", image.Registry, image.Namespace, image.Name)
		availableTags, fetched := uniqueImages[imageKey]
		if !fetched {
			log.Printf("INFO: Fetching tags for image %s\n", imageKey)
			// Fetch image tags from registry
			regClient := registryClients[image.Registry]
			tags, err := regClient.GetTags(ctx, image)
			if err != nil {
				status.StatusMessage = fmt.Sprintln("Error fetching images")
				status.Error = err.Error()
				allUpdateStatuses = append(allUpdateStatuses, status)
				log.Printf("ERROR: Failed to get tags for %s from registry '%s': %v. Skipping.", imageKey, image.Registry, err)
				continue
			}

			uniqueImages[imageKey] = tags
			availableTags = tags
			log.Printf("Found %d tags for image %s\n", len(tags), imageKey)
		}

		latestUpdateTag, isUpdateAvailable, err := FindLatestUpdate(image.Tag, availableTags)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error comparing tags")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Printf("Unexpected error during semver update checks: %v\n", err)
		}

		status.UpdateAvailable = isUpdateAvailable
		status.LatestAvailableTag = latestUpdateTag

		if isUpdateAvailable {
			status.StatusMessage = "Update available"
			allUpdateStatuses = append(allUpdateStatuses, status)
			// log.Printf("DEBUG: Container %s (%s) can be upgraded from %s to %s\n", cnt.Names[0], imageKey, image.Tag, latestUpdateTag)
		} else {
			status.StatusMessage = "No update available"
			log.Printf("No update for image %s", imageKey)
		}
	}

	for _, update := range allUpdateStatuses {
		if update.UpdateAvailable {
			fmt.Printf("Container %s (%s) can be upgraded to %s\n",
				update.ContainerName,
				update.Image.Raw,
				update.LatestAvailableTag,
			)
		}
	}
}

// getRunningContainerImages connects to the Docker daemon and returns a list of running containers.
func getRunningContainerImages(ctx context.Context) ([]container.Summary, error) {
	log.Println("Connecting to Docker daemon...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}
	defer cli.Close()

	log.Println("Listing running containers...")
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %w", err)
	}
	return containers, nil
}
