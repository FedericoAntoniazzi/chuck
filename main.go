package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	fmt.Println("Hello from Chuck! Core logic and database interaction not yet implemented.")

	log.Println("Chuck finished checking. No updates found (yet!).")
}
