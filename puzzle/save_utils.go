package puzzle

import (
	"encoding/json"
	"godoku/config"
	"godoku/logger"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// GameState represents the state of a Sudoku game
type GameState struct {
	OriginalPuzzle [9][9]int   `json:"originalPuzzle"`
	CurrentState   [9][9]int   `json:"currentState"`
	UserEdited     [9][9]bool  `json:"userEdited"`
	Timestamp      time.Time   `json:"timestamp"`
}

// GetMostRecentSaveFile returns the most recent saved game file path if exists
// Returns the file path, and a boolean indicating if a save was found
func GetMostRecentSaveFile() (string, bool) {
	// Check if saves directory exists
	savesDir := config.SavesDir
	if _, err := os.Stat(savesDir); os.IsNotExist(err) {
		logger.Info("No saves directory found")
		return "", false
	}
	
	// Read saved game files
	files, err := os.ReadDir(savesDir)
	if err != nil {
		logger.Error("Error reading saves directory: %v", err)
		return "", false
	}
	
	// Filter for JSON files and get their full stats
	var saveFiles []fs.FileInfo
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			fullPath := filepath.Join(savesDir, file.Name())
			fileInfo, err := os.Stat(fullPath)
			if err != nil {
				logger.Error("Error reading file info for %s: %v", file.Name(), err)
				continue
			}
			saveFiles = append(saveFiles, fileInfo)
		}
	}
	
	if len(saveFiles) == 0 {
		logger.Info("No save files found")
		return "", false
	}
	
	// Sort by modification time, most recent first
	sort.Slice(saveFiles, func(i, j int) bool {
		return saveFiles[i].ModTime().After(saveFiles[j].ModTime())
	})
	
	// Return the most recent save file
	mostRecentSave := filepath.Join(savesDir, saveFiles[0].Name())
	logger.Info("Found most recent save: %s (from %s)", 
		saveFiles[0].Name(), 
		saveFiles[0].ModTime().Format(time.RFC3339))
	
	return mostRecentSave, true
}

// LoadLastGame attempts to load the most recent saved game
// Returns the puzzle state and true if successful, false otherwise
func LoadLastGame() (*GameState, bool) {
	saveFilePath, found := GetMostRecentSaveFile()
	if !found {
		logger.Info("No save file found to auto-load")
		return nil, false
	}
	
	// Read file content
	data, err := os.ReadFile(saveFilePath)
	if err != nil {
		logger.Error("Error reading save file: %v", err)
		return nil, false
	}
	
	// Unmarshal JSON
	var gameState GameState
	err = json.Unmarshal(data, &gameState)
	if err != nil {
		logger.Error("Error parsing save file: %v", err)
		return nil, false
	}
	
	logger.Info("Successfully loaded last game from %s", saveFilePath)
	return &gameState, true
}
