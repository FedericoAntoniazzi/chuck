package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/FedericoAntoniazzi/chuck/core"
	"github.com/FedericoAntoniazzi/chuck/registry/dockerhub"
	"github.com/FedericoAntoniazzi/chuck/types"
	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultDBFileName    = "chuck.db"
	defaultLoggingLevel  = "warn"
	defaultLoggingFormat = "text"
)

// registryClient defines the capabilities of a generic client for container registries
type registryClient interface {
	// GetTags fetches all available tags for a given image from the registry
	GetTags(ctx context.Context, image types.Image) ([]string, error)
}

func defineLogger(logLevel string, logFormat string) (*zap.SugaredLogger, error) {
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder

	switch strings.ToLower(logFormat) {
	// Machine-readable JSON format for console
	case "json":
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "text":
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	parsedLevel := zap.InfoLevel
	if err := parsedLevel.UnmarshalText([]byte(strings.ToLower(logLevel))); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level %s. Defaulting to 'info'. Error: %s", logLevel, err)
	}
	atomicLevel := zap.NewAtomicLevelAt(parsedLevel)

	// Lock the output to allow safe concurrent writes
	outputSyncer := zapcore.Lock(os.Stdout)

	core := zapcore.NewCore(encoder, outputSyncer, atomicLevel)
	baseLogger := zap.New(core, zap.AddCaller())

	return baseLogger.Sugar(), nil
}

func main() {
	// --- CLI Flags Definition ---
	logFormat := flag.String("logFormat", defaultLoggingFormat, "Log format (text, json)")
	logLevel := flag.String("logLevel", defaultLoggingLevel, "Configure the logging level (debug, info, warn, error)")
	dbPath := flag.String("db-path", defaultDBFileName, "Path to the SQLite database file")

	flag.Parse()

	// --- Logging Setup ---
	logger, err := defineLogger(*logLevel, *logFormat)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}
	defer logger.Sync()

	// --- Database Path Handling ---
	// Resolve the absolute path for the database file
	resolvedDBPath, err := filepath.Abs(*dbPath)
	if err != nil {
		logger.Fatal("could not resolve absolute path for database file", "path", *dbPath, "err", err)
	}
	logger.Debugf("Using database file: %s", resolvedDBPath)

	// --- Core Logic Placeholder ---
	// This is where the core logic for Docker interaction, registry checks,
	// and SQLite operations will eventually go.

	// Create a background context for Docker API calls
	ctx := context.Background()

	registryClients := make(map[string]registryClient)
	registryClients["docker.io"] = dockerhub.NewClient()
	// Hint: registryClients["ghcr.io"] = github.NewClient()

	containers, err := core.GetRunningContainerImages(ctx, logger)
	if err != nil {
		logger.Fatalf("Failed to get running containers: %v", err)
	}

	if len(containers) == 0 {
		logger.Info("No running containers found")
		return
	}

	logger.Infof("Found %d running containers", len(containers))
	// Store the images which tags have already been queried
	uniqueImages := make(map[string][]string)

	var allUpdateStatuses []types.ImageUpdateStatus

	for _, cnt := range containers {
		containerName := ""

		if len(cnt.Names[0]) > 0 {
			containerName = strings.TrimPrefix(cnt.Names[0], "/")
		}

		logger.Debug("processing container ", containerName)

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
			logger.Warn("skipping invalid image name", "image", cnt.Image, "error", err)
			continue
		}

		status.Image = image
		status.OriginalTag = image.Tag

		// Check if the registry is supported
		if _, ok := registryClients[image.Registry]; !ok {
			status.StatusMessage = fmt.Sprintln("Unsupported registry")
			status.Error = "Unsupported registry"
			allUpdateStatuses = append(allUpdateStatuses, status)
			logger.Warn("skipping unsupported registry (", image.Registry, ") for image ", image.Raw)
			continue
		}

		// Check if the tag is valid SemVer
		_, err = semver.NewVersion(image.Tag)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error parsing image tag")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			logger.Warnw("skipping invalid semver tag", "image", image.Raw, "tag", image.Tag)
			continue
		}

		// Check if image tags have already been fetched
		imageKey := fmt.Sprintf("%s/%s/%s", image.Registry, image.Namespace, image.Name)
		availableTags, fetched := uniqueImages[imageKey]
		if !fetched {
			logger.Debugf("fetching tags for image %s", imageKey)
			// Fetch image tags from registry
			regClient := registryClients[image.Registry]
			tags, err := regClient.GetTags(ctx, image)
			if err != nil {
				status.StatusMessage = fmt.Sprintln("Error fetching images")
				status.Error = err.Error()
				allUpdateStatuses = append(allUpdateStatuses, status)
				logger.Errorf("error retrieving tags from registry", "image", imageKey, "registry", image.Registry, "error", err)
				continue
			}

			uniqueImages[imageKey] = tags
			availableTags = tags
			logger.Debugf("Found %d tags for image %s", len(tags), imageKey)
		}

		latestUpdateTag, isUpdateAvailable, err := core.FindLatestUpdate(image.Tag, availableTags)
		if err != nil {
			status.StatusMessage = fmt.Sprintln("Error comparing tags")
			status.Error = err.Error()
			allUpdateStatuses = append(allUpdateStatuses, status)
			logger.Error("unexpected error during semver checks", "error", err)
		}

		status.UpdateAvailable = isUpdateAvailable
		status.LatestAvailableTag = latestUpdateTag

		if isUpdateAvailable {
			status.StatusMessage = "Update available"
			allUpdateStatuses = append(allUpdateStatuses, status)
			logger.Debugf("Container %s (%s) can be upgraded to %s", containerName, imageKey, latestUpdateTag)
		} else {
			status.StatusMessage = "No update available"
			logger.Debugf("checked updates for %s (%s). No updates available", containerName, imageKey)
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
