package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/FedericoAntoniazzi/chuck/core"
	"github.com/FedericoAntoniazzi/chuck/registry/dockerhub"
	"github.com/FedericoAntoniazzi/chuck/types"
	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogFileName  = "chuck.log"
	defaultDBFileName   = "chuck.db"
	defaultLoggingLevel = "warn"
)

// registryClient defines the capabilities of a generic client for container registries
type registryClient interface {
	// GetTags fetches all available tags for a given image from the registry
	GetTags(ctx context.Context, image types.Image) ([]string, error)
}

func main() {
	// --- CLI Flags Definition ---
	logLevel := flag.String("logLevel", defaultLoggingLevel, "Configure the logging level")
	logToFile := flag.Bool("logFile", false, "Enable logging output to a file")
	logFilePath := flag.String("logFilePath", defaultLogFileName, "Log file path")
	dbPath := flag.String("db-path", defaultDBFileName, "Path to the SQLite database file")

	flag.Parse()

	// --- Logging Setup ---
	atomicLevel := zap.NewAtomicLevel()

	level, err := zapcore.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("invalid log level: %v", err)
		return
	}
	atomicLevel.SetLevel(level)

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = atomicLevel
	loggerConfig.Encoding = "console"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

	if *logToFile {
		loggerConfig.OutputPaths = []string{*logFilePath}
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("failed to initalize zap logger: %v", err)
	}

	defer logger.Sync()
	log := logger.Sugar()
	log.Debug("Chuck started")

	// --- Database Path Handling ---
	// Resolve the absolute path for the database file
	resolvedDBPath, err := filepath.Abs(*dbPath)
	if err != nil {
		log.Fatal("could not resolve absolute path for database file", "path", *dbPath, "err", err)
	}
	log.Debugf("Using database file: %s", resolvedDBPath)

	// --- Core Logic Placeholder ---
	// This is where the core logic for Docker interaction, registry checks,
	// and SQLite operations will eventually go.

	// Create a background context for Docker API calls
	ctx := context.Background()

	registryClients := make(map[string]registryClient)
	registryClients["docker.io"] = dockerhub.NewClient()
	// Hint: registryClients["ghcr.io"] = github.NewClient()

	containers, err := core.GetRunningContainerImages(ctx, log)
	if err != nil {
		log.Fatalf("Failed to get running containers: %v", err)
	}

	if len(containers) == 0 {
		log.Info("No running containers found")
		return
	}

	log.Infof("Found %d running containers", len(containers))
	// Store the images which tags have already been queried
	uniqueImages := make(map[string][]string)

	var allUpdateStatuses []types.ImageUpdateStatus

	for _, cnt := range containers {
		containerName := ""

		if len(cnt.Names[0]) > 0 {
			containerName = strings.TrimPrefix(cnt.Names[0], "/")
		}

		log.Debug("processing container ", containerName)

		status := types.ImageUpdateStatus{
			ContainerID:   cnt.ID,
			ContainerName: containerName,
			OriginalTag:   "latest",
			StatusMessage: "Processing",
		}

		image, err := core.ParseImageName(cnt.Image)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error parsing image name")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Warn("skipping invalid image name", "image", cnt.Image, "error", err)
			continue
		}

		status.Image = image
		status.OriginalTag = image.Tag

		// Check if the registry is supported
		if _, ok := registryClients[image.Registry]; !ok {
			status.StatusMessage = fmt.Sprintln("Unsupported registry")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Warn("skipping unsupported registry", "registry", image.Registry)
			continue
		}

		// Check if the tag is valid SemVer
		_, err = semver.NewVersion(image.Tag)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error parsing image tag")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Warnw("skipping invalid semver tag", "image", image.Raw, "tag", image.Tag)
			continue
		}

		// Check if image tags have already been fetched
		imageKey := fmt.Sprintf("%s/%s/%s", image.Registry, image.Namespace, image.Name)
		availableTags, fetched := uniqueImages[imageKey]
		if !fetched {
			log.Debugf("fetching tags for image %s", imageKey)
			// Fetch image tags from registry
			regClient := registryClients[image.Registry]
			tags, err := regClient.GetTags(ctx, image)
			if err != nil {
				status.StatusMessage = fmt.Sprintln("Error fetching images")
				status.Error = err.Error()
				allUpdateStatuses = append(allUpdateStatuses, status)
				log.Errorf("error retrieving tags from registry", "image", imageKey, "registry", image.Registry, "error", err)
				continue
			}

			uniqueImages[imageKey] = tags
			availableTags = tags
			log.Debugf("Found %d tags for image %s\n", len(tags), imageKey)
		}

		latestUpdateTag, isUpdateAvailable, err := core.FindLatestUpdate(image.Tag, availableTags)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error comparing tags")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Error("unexpected error during semver checks", "error", err)
		}

		status.UpdateAvailable = isUpdateAvailable
		status.LatestAvailableTag = latestUpdateTag

		if isUpdateAvailable {
			status.StatusMessage = "Update available"
			allUpdateStatuses = append(allUpdateStatuses, status)
			log.Debugf("Container %s (%s) can be upgraded to %s", containerName, imageKey, latestUpdateTag)
		} else {
			status.StatusMessage = "No update available"
			log.Debugf("checked updates for %s (%s). No updates available", containerName, imageKey)
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
