package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	defaultLogFileName = "chuck.log"
	defaultDBFileName  = "chuck.db"
)

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
	containers, err := getRunningContainerImages(ctx)
	if err != nil {
		log.Fatalf("Failed to get running containers: %v", err)
	}

	if len(containers) == 0 {
		log.Println("No running containers found")
	} else {
		log.Printf("Found %d running containers:\n", len(containers))
		for _, container := range containers {
			log.Printf("\tContainer ID: %s, Image: %s", container.ID[:12], container.Image)
		}
	}

	log.Println("Chuck finished checking. No updates found (yet!).")
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
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %w", err)
	}
	return containers, nil
}
