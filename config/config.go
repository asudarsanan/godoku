package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Directories and paths
var (
	// ConfigDir is the main configuration directory
	ConfigDir string

	// PuzzlesDir is where puzzle files are stored
	PuzzlesDir string

	// SavesDir is where saved games are stored
	SavesDir string

	// LogDir is where log files are stored
	LogDir string

	// CurrentLogFile is the path to the current log file
	CurrentLogFile string
)

// Init initializes config paths without creating directories
func Init() {
	// Initialize directory paths
	homeDir, err := os.UserHomeDir()
	if err == nil {
		if runtime.GOOS == "windows" {
			ConfigDir = filepath.Join(homeDir, "AppData", "Local", "godoku")
		} else {
			// For Linux, macOS, and other Unix-like systems
			ConfigDir = filepath.Join(homeDir, ".config", "godoku")
		}
		
		PuzzlesDir = filepath.Join(ConfigDir, "puzzles")
		SavesDir = filepath.Join(ConfigDir, "saves")
	} else {
		// Fallback to current directory if home directory can't be determined
		ConfigDir = filepath.Join(".", ".config", "godoku")
		PuzzlesDir = filepath.Join(ConfigDir, "puzzles")
		SavesDir = filepath.Join(ConfigDir, "saves")
	}

	// Set up log directory
	initLogDirectory()
}

// InitConfig initializes the configuration directories
func InitConfig() error {
	// Initialize paths
	Init()

	// Create directories if they don't exist
	dirs := []string{ConfigDir, PuzzlesDir, SavesDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Initialize log directory
	if err := InitLogDirectory(); err != nil {
		return err
	}

	// Copy default puzzles if they don't exist in the new location
	return copyDefaultPuzzles()
}

// copyDefaultPuzzles copies the default puzzles from the resources directory to the config directory
func copyDefaultPuzzles() error {
	// Check if default puzzles file exists in the config directory
	defaultPuzzlesPath := filepath.Join(PuzzlesDir, "puzzles.json")
	if _, err := os.Stat(defaultPuzzlesPath); err == nil {
		// File already exists, no need to copy
		return nil
	}

	// Possible source locations for puzzle file
	sourcePaths := []string{
		filepath.Join("resources", "puzzles.json"),
	}

	// Try to find executable directory
	if execDir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		sourcePaths = append(sourcePaths, filepath.Join(execDir, "resources", "puzzles.json"))
	}

	// Try each possible source path
	var data []byte
	var readErr error
	for _, sourcePath := range sourcePaths {
		data, readErr = os.ReadFile(sourcePath)
		if readErr == nil {
			break // Successfully read the file
		}
	}

	// If we couldn't find or read the source file, create an empty puzzle
	if readErr != nil {
		emptyPuzzles := []byte(`{"grid":[[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0]]}`)
		return os.WriteFile(defaultPuzzlesPath, emptyPuzzles, 0644)
	}

	// Write to destination
	return os.WriteFile(defaultPuzzlesPath, data, 0644)
}

// initLogDirectory sets up the log directory based on OS
func initLogDirectory() {
	if runtime.GOOS == "windows" {
		// On Windows, use %TEMP%\godoku for logs
		LogDir = filepath.Join(os.TempDir(), "godoku")
	} else {
		// On Linux/Unix systems, use /tmp/godoku
		LogDir = filepath.Join("/tmp", "godoku")
	}

	// Create a unique log file for this session
	timestamp := time.Now().Format("20060102-150405")
	CurrentLogFile = filepath.Join(LogDir, "godoku-"+timestamp+".log")
}

// InitLogDirectory creates the log directory if it doesn't exist
func InitLogDirectory() error {
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		return err
	}
	return nil
}

// GetLogPath returns the path to the current log file
func GetLogPath() string {
	return CurrentLogFile
}

// GetPuzzlePath returns the path to a puzzle file
func GetPuzzlePath(filename string) string {
	return filepath.Join(PuzzlesDir, filename)
}

// GetSavePath returns the path to a save file
func GetSavePath(filename string) string {
	return filepath.Join(SavesDir, filename)
}
