package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"godoku/config"
	"godoku/logger"
	"os"
	"path/filepath"
)

// ShowPuzzleSelector displays a UI for selecting a puzzle
// The callback function is called with the selected puzzle file path
// Returns true if a puzzle was selected, false if canceled
func ShowPuzzleSelector(loadCallback func(string) bool) bool {
	puzzlesDir := config.PuzzlesDir
	defaultPuzzle := filepath.Join(puzzlesDir, "puzzles.json")  // Default puzzle name
	
	// Check if puzzles directory exists
	if _, err := os.Stat(puzzlesDir); os.IsNotExist(err) {
		logger.Info("Puzzles directory not found, using default puzzle")
		return loadCallback(defaultPuzzle)
	}
	
	// Read puzzle files
	puzzleFiles, err := os.ReadDir(puzzlesDir)
	if err != nil {
		logger.Error("Error reading puzzles directory: %v", err)
		return loadCallback(defaultPuzzle)
	}
	
	// Filter for JSON files
	var jsonFiles []string
	for _, file := range puzzleFiles {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}
	
	if len(jsonFiles) == 0 {
		logger.Info("No puzzle files found, using default")
		return loadCallback(defaultPuzzle)
	}
	
	// If only one puzzle file exists, use it
	if len(jsonFiles) == 1 {
		logger.Info("Using single puzzle file: %s", jsonFiles[0])
		return loadCallback(filepath.Join(puzzlesDir, jsonFiles[0]))
	}
	
	// Create UI for puzzle selection
	app := tview.NewApplication()
	
	// Create list to select puzzle file
	listBox := tview.NewList()
	listBox.SetTitle(" Select a Puzzle ")
	listBox.SetTitleAlign(tview.AlignCenter)
	listBox.SetBorder(true)
	
	// Add default puzzle option first
	listBox.AddItem("Default Puzzle", "Standard puzzle collection", 0, func() {
		app.Stop()
		loadCallback(defaultPuzzle)
	})
	
	// Add puzzle files
	for _, file := range jsonFiles {
		// Display file without extension
		displayName := filepath.Base(file)
		ext := filepath.Ext(displayName)
		if ext != "" {
			displayName = displayName[:len(displayName)-len(ext)]
		}
		
		// Format display name nicely
		displayName = fmt.Sprintf("%s Puzzle", displayName)
		
		// Local copy for closure
		puzzlePath := filepath.Join(puzzlesDir, file)
		
		listBox.AddItem(displayName, file, 0, func() {
			app.Stop()
			loadCallback(puzzlePath)
		})
	}
	
	// Add exit option
	listBox.AddItem("Cancel", "Return to main menu", 0, func() {
		app.Stop()
	})
	
	// Set up key handling
	listBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		}
		return event
	})
	
	// Show the list in a flex container to center it
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(listBox, 15, 1, true).
				AddItem(nil, 0, 1, false),
			50, 1, true).
		AddItem(nil, 0, 1, false)
	
	// Run the app
	selected := true
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		logger.Error("Error running puzzle selector UI: %v", err)
		selected = false
	}
	
	// Check if the list still has focus after app.Stop()
	// This indicates the user canceled the selection
	if listBox.HasFocus() {
		selected = false
	}
	
	return selected
}
