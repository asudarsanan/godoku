package config

import (
	"os"
	"path/filepath"
)

// Directories and paths
var (
	// ConfigDir is the main configuration directory
	ConfigDir string

	// PuzzlesDir is where puzzle files are stored
	PuzzlesDir string

	// SavesDir is where saved games are stored
	SavesDir string
)

func init() {
	// Initialize directory paths
	homeDir, err := os.UserHomeDir()
	if err == nil {
		ConfigDir = filepath.Join(homeDir, ".config", "godoku")
		PuzzlesDir = filepath.Join(ConfigDir, "puzzles")
		SavesDir = filepath.Join(ConfigDir, "saves")
	} else {
		// Fallback to current directory if home directory can't be determined
		ConfigDir = filepath.Join(".", ".config", "godoku")
		PuzzlesDir = filepath.Join(ConfigDir, "puzzles")
		SavesDir = filepath.Join(ConfigDir, "saves")
	}
}

// InitConfig initializes the configuration directories
func InitConfig() error {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Set up config paths
	ConfigDir = filepath.Join(homeDir, ".config", "godoku")
	PuzzlesDir = filepath.Join(ConfigDir, "puzzles")
	SavesDir = filepath.Join(ConfigDir, "saves")

	// Create directories if they don't exist
	dirs := []string{ConfigDir, PuzzlesDir, SavesDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Copy default puzzles if they don't exist in the new location
	err = copyDefaultPuzzles()
	if err != nil {
		return err
	}

	return nil
}

// copyDefaultPuzzles copies the default puzzles from the resources directory to the config directory
func copyDefaultPuzzles() error {
	// Check if default puzzles file exists in the config directory
	defaultPuzzlesPath := filepath.Join(PuzzlesDir, "puzzles.json")
	if _, err := os.Stat(defaultPuzzlesPath); err == nil {
		// File already exists, no need to copy
		return nil
	}

	// Get executable directory
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// Source puzzles file
	sourcePath := filepath.Join(execDir, "resources", "puzzles.json")
	
	// If source doesn't exist in the expected location, try current directory
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		sourcePath = filepath.Join("resources", "puzzles.json")
	}

	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		// Create an empty puzzles file as fallback
		emptyPuzzles := []byte(`{"grid":[[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0]]}`)
		return os.WriteFile(defaultPuzzlesPath, emptyPuzzles, 0644)
	}

	// Read source file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Write to destination
	return os.WriteFile(defaultPuzzlesPath, data, 0644)
}

// GetPuzzlePath returns the path to a puzzle file
func GetPuzzlePath(filename string) string {
	return filepath.Join(PuzzlesDir, filename)
}

// GetSavePath returns the path to a save file
func GetSavePath(filename string) string {
	return filepath.Join(SavesDir, filename)
}
