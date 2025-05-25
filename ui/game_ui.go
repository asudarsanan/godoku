package ui

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"godoku/config"
	"godoku/logger"
	"godoku/puzzle"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UI represents the game interface
type UI struct {
	app           *tview.Application
	table         *tview.Table
	game          [9][9]int
	originalGame  [9][9]int
	userEdited    [9][9]bool
	mainContainer *tview.Flex
	controlBar    *tview.Flex
	resetButton   *tview.Button
	saveButton    *tview.Button
	loadButton    *tview.Button
	exitButton    *tview.Button
	statusText    *tview.TextView
}

// Using puzzle.GameState for game state representation

// NewUI creates a new game UI instance
func NewUI(game [9][9]int) *UI {
	// Create a copy of the original puzzle
	var originalGame [9][9]int
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			originalGame[r][c] = game[r][c]
		}
	}
	
	ui := &UI{
		app:          tview.NewApplication(),
		table:        tview.NewTable().SetBorders(true),
		game:         game,
		originalGame: originalGame,
		userEdited:   [9][9]bool{},
		statusText:   tview.NewTextView().SetText("").SetTextAlign(tview.AlignCenter),
	}
	
	// Create control buttons
	ui.resetButton = tview.NewButton("Reset").
		SetSelectedFunc(func() {
			ui.resetPuzzle()
		})
	
	ui.saveButton = tview.NewButton("Save").
		SetSelectedFunc(func() {
			ui.savePuzzle()
		})
	
	ui.loadButton = tview.NewButton("Load").
		SetSelectedFunc(func() {
			ui.loadPuzzle()
		})
	
	ui.exitButton = tview.NewButton("Exit").
		SetSelectedFunc(func() {
			ui.app.Stop()
		})
	
	// Create control bar
	ui.controlBar = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(ui.resetButton, 8, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(ui.saveButton, 8, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(ui.loadButton, 8, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(ui.exitButton, 8, 0, true).
		AddItem(nil, 0, 1, false)
	
	// Create main container with puzzle grid and control bar
	ui.mainContainer = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(ui.table, 37, 0, true). // Width of puzzle grid
				AddItem(nil, 0, 1, false),
			19, 0, true). // Height of puzzle grid
		AddItem(nil, 1, 0, false).
		AddItem(ui.statusText, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(ui.controlBar, 1, 0, true).
		AddItem(nil, 0, 1, false)
	
	return ui
}

// Run starts the game UI
func (ui *UI) Run() error {
	// Set up logging - using project's logger package instead
	// The logger is already initialized in main.go, so we don't need to initialize it here
	logger.Info("Starting game UI")

	ui.initGrid()
	return ui.app.SetRoot(ui.mainContainer, true).EnableMouse(true).Run()
}

// initGrid initializes the Sudoku grid
func (ui *UI) initGrid() {
	for r, row := range ui.game {
		for c, col := range row {
			var text string
			if col == 0 {
				text = "   " // Empty cell
			} else {
				text = fmt.Sprintf(" %d ", col)
			}

			// Set cell color based on whether it's an original puzzle cell or user-edited
			color := tcell.ColorYellowGreen // Default for empty cells
			immutable := false
			
			if col != 0 {
				if ui.userEdited[r][c] {
					// User edited cell
					color = tcell.ColorRed
				} else {
					// Original puzzle cell
					color = tcell.ColorAqua
					immutable = true
				}
			}
			
			bgColor := tcell.ColorSilver
			if (r == 3 || r == 4 || r == 5) || (c == 3 || c == 4 || c == 5) {
				bgColor = tcell.ColorGray
			}

			ui.table.SetCell(r, c, tview.NewTableCell(text).
				SetTextColor(color).SetBackgroundColor(bgColor).
				SetAlign(tview.AlignCenter).
				SetSelectable(!immutable || ui.userEdited[r][c]))
		}
	}

	ui.table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyCtrlC {
			ui.app.Stop()
		}
		if key == tcell.KeyEnter {
			ui.table.SetSelectable(true, true)
		}
	})

	ui.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Add Tab key to switch focus between puzzle grid and control buttons
		if event.Key() == tcell.KeyTab {
			ui.app.SetFocus(ui.resetButton)
			return nil
		}
	
		row, col := ui.table.GetSelection()
		
		// Check if row and col are valid indices
		if row < 0 || row >= 9 || col < 0 || col >= 9 {
			logger.Info("Invalid cell selection (%d, %d)", row, col)
			return event
		}
		
		cell := ui.table.GetCell(row, col)

		if !ui.userEdited[row][col] && ui.game[row][col] != 0 {
			logger.Info("The cell is non-editable (%d, %d)", row, col)
			return event
		}

		switch event.Key() {
		case tcell.KeyRune:
			r := event.Rune()
			if r >= '1' && r <= '9' {
				newText := string(r)
				cell.SetText(fmt.Sprintf(" %s ", newText)).SetTextColor(tcell.ColorRed)
				ui.updateGrid(row, col, newText)
				logger.Info("Updated cell (%d, %d) with new value: %s", row, col, newText)
			} else if r == '0' {
				cell.SetText("   ").SetTextColor(tcell.ColorYellowGreen) // Use consistent color
				ui.updateGrid(row, col, "")
				logger.Info("Cleared cell (%d, %d)", row, col)
			}
		case tcell.KeyBackspace, tcell.KeyDelete:
			cell.SetText("   ").SetTextColor(tcell.ColorYellowGreen) // Use consistent color
			ui.updateGrid(row, col, "")
			logger.Info("Cleared cell (%d, %d)", row, col)
		}
		return event
	})

	// Set up input capture for the main container to handle global key events
	ui.mainContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Cycle focus between the table and the buttons
			if ui.app.GetFocus() == ui.table {
				ui.app.SetFocus(ui.resetButton)
				return nil
			} else if ui.app.GetFocus() == ui.resetButton || 
				ui.app.GetFocus() == ui.saveButton || 
				ui.app.GetFocus() == ui.loadButton || 
				ui.app.GetFocus() == ui.exitButton {
				ui.app.SetFocus(ui.table)
				return nil
			}
		}
		return event
	})

	// Set up focus handling for the buttons
	ui.resetButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			ui.app.SetFocus(ui.saveButton)
			return nil
		}
		return event
	})

	ui.saveButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			ui.app.SetFocus(ui.resetButton)
			return nil
		} else if event.Key() == tcell.KeyRight {
			ui.app.SetFocus(ui.loadButton)
			return nil
		}
		return event
	})
	
	ui.loadButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			ui.app.SetFocus(ui.saveButton)
			return nil
		} else if event.Key() == tcell.KeyRight {
			ui.app.SetFocus(ui.exitButton)
			return nil
		}
		return event
	})

	ui.exitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			ui.app.SetFocus(ui.loadButton)
			return nil
		}
		return event
	})

	// Set selection behavior
	ui.table.SetSelectable(true, true)
}

// updateGrid updates the grid data with the new value
func (ui *UI) updateGrid(row int, col int, text string) {
	if len(text) == 0 {
		ui.game[row][col] = 0
	} else {
		num, err := strconv.Atoi(text)
		if err == nil {
			ui.game[row][col] = num
		}
	}
	ui.userEdited[row][col] = true
}

// resetPuzzle resets the puzzle to its original state by clearing all user inputs
func (ui *UI) resetPuzzle() {
	// Create a modal for confirmation
	modal := tview.NewModal().
		SetText("Are you sure you want to reset the puzzle? All progress will be lost.").
		AddButtons([]string{"Reset", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Reset" {
				// Clear all user inputs
				ui.userEdited = [9][9]bool{}
				
				// Reset game state to original
				for r := 0; r < 9; r++ {
					for c := 0; c < 9; c++ {
						ui.game[r][c] = ui.originalGame[r][c]
						
						// Update cell display
						cell := ui.table.GetCell(r, c)
						if ui.game[r][c] == 0 {
							cell.SetText("   ").SetTextColor(tcell.ColorYellowGreen)
						} else {
							// After reset, all non-empty cells are original puzzle cells
							cell.SetText(fmt.Sprintf(" %d ", ui.game[r][c]))
							cell.SetTextColor(tcell.ColorAqua)
							cell.SetSelectable(false)
						}
					}
				}
				
				logger.Info("Puzzle has been reset to its original state")
				ui.ShowStatus("Puzzle reset successfully")
			}
			
			// Return to the main UI
			ui.app.SetRoot(ui.mainContainer, true)
			ui.app.SetFocus(ui.table)
		})
	
	// Show the modal
	ui.app.SetRoot(modal, true)
}

// savePuzzle saves the current puzzle state to a file
func (ui *UI) savePuzzle() {
	// Create saves directory if it doesn't exist
	savesDir := config.SavesDir
	if _, err := os.Stat(savesDir); os.IsNotExist(err) {
		err := os.MkdirAll(savesDir, 0755)
		if err != nil {
			logger.Error("Error creating saves directory: %v", err)
			ui.ShowStatus("Error: Could not create saves directory")
			return
		}
	}
	
	// Generate filename with timestamp
	timestamp := time.Now()
	filename := filepath.Join(savesDir, fmt.Sprintf("sudoku_save_%s.json", 
		timestamp.Format("2006-01-02_15-04-05")))
	
	// Create game state
	gameState := puzzle.GameState{
		OriginalPuzzle: ui.originalGame,
		CurrentState:   ui.game,
		UserEdited:     ui.userEdited,
		Timestamp:      timestamp,
	}
	
	// Marshal to JSON
	jsonData, err := json.MarshalIndent(gameState, "", "  ")
	if err != nil {
		logger.Error("Error marshalling game state: %v", err)
		ui.ShowStatus("Error: Could not save game")
		return
	}
	
	// Write to file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		logger.Error("Error writing save file: %v", err)
		ui.ShowStatus("Error: Could not write save file")
		return
	}
	
	logger.Info("Game saved to %s", filename)
	ui.ShowStatus(fmt.Sprintf("Game saved to %s", filepath.Base(filename)))
	ui.app.SetFocus(ui.table)
}

// loadPuzzle loads a saved puzzle from a file
func (ui *UI) loadPuzzle() {
	// Check if saves directory exists
	savesDir := config.SavesDir
	if _, err := os.Stat(savesDir); os.IsNotExist(err) {
		ui.ShowStatus("No saved games found")
		return
	}
	
	// Read saved game files
	files, err := os.ReadDir(savesDir)
	if err != nil {
		logger.Error("Error reading saves directory: %v", err)
		ui.ShowStatus("Error reading saved games")
		return
	}
	
	// Filter for json files
	var saveFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			saveFiles = append(saveFiles, file.Name())
		}
	}
	
	if len(saveFiles) == 0 {
		ui.ShowStatus("No saved games found")
		return
	}
	
	// Create list to select save file
	list := tview.NewList()
	list.SetTitle("Select a saved game").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	for i, file := range saveFiles {
		// Display file without extension
		displayName := file
		if len(displayName) > 40 {
			displayName = displayName[:37] + "..."
		}
		
		// Local copy of i for the closure
		index := i
		list.AddItem(displayName, filepath.Join(savesDir, file), 0, func() {
			ui.loadGameFromFile(filepath.Join(savesDir, saveFiles[index]))
		})
	}
	
	// Add cancel option
	list.AddItem("Cancel", "", 0, func() {
		ui.app.SetRoot(ui.mainContainer, true)
		ui.app.SetFocus(ui.table)
	})
	
	// Show the list in a flex container to center it
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(list, 15, 1, true).
				AddItem(nil, 0, 1, false),
			50, 1, true).
		AddItem(nil, 0, 1, false)
	
	ui.app.SetRoot(flex, true)
}

// loadGameFromFile loads a game from the specified file
func (ui *UI) loadGameFromFile(filePath string) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Error reading save file: %v", err)
		ui.ShowStatus("Error reading save file")
		ui.app.SetRoot(ui.mainContainer, true)
		ui.app.SetFocus(ui.table)
		return
	}
	
	// Unmarshal JSON
	var gameState puzzle.GameState
	err = json.Unmarshal(data, &gameState)
	if err != nil {
		logger.Error("Error parsing save file: %v", err)
		ui.ShowStatus("Error parsing save file")
		ui.app.SetRoot(ui.mainContainer, true)
		ui.app.SetFocus(ui.table)
		return
	}
	
	// Update game state
	ui.originalGame = gameState.OriginalPuzzle
	ui.game = gameState.CurrentState
	ui.userEdited = gameState.UserEdited
	
	// Refresh grid
	ui.refreshGrid()
	
	ui.app.SetRoot(ui.mainContainer, true)
	ui.app.SetFocus(ui.table)
	
	ui.ShowStatus("Game loaded successfully")
	logger.Info("Game loaded from %s", filePath)
}

// refreshGrid updates the UI grid based on the current game state
func (ui *UI) refreshGrid() {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			cell := ui.table.GetCell(r, c)
			
			// Set cell content
			if ui.game[r][c] == 0 {
				cell.SetText("   ")
			} else {
				cell.SetText(fmt.Sprintf(" %d ", ui.game[r][c]))
			}
			
			// Set cell color based on whether it's an original puzzle cell or user-edited
			color := tcell.ColorYellowGreen // Default for empty cells
			isSelectable := true
			
			if ui.game[r][c] != 0 {
				if ui.userEdited[r][c] {
					// User edited cell
					color = tcell.ColorRed
				} else {
					// Original puzzle cell (not user-edited)
					color = tcell.ColorAqua
					isSelectable = false
				}
			}
			
			cell.SetTextColor(color)
			cell.SetSelectable(isSelectable)
			
			// Set background color based on position
			bgColor := tcell.ColorSilver
			if (r == 3 || r == 4 || r == 5) || (c == 3 || c == 4 || c == 5) {
				bgColor = tcell.ColorGray
			}
			cell.SetBackgroundColor(bgColor)
		}
	}
}

// Using ShowStatus method instead

// SetGameState updates the game state with provided values
func (ui *UI) SetGameState(gameState [9][9]int, userEdited [9][9]bool) {
	ui.game = gameState
	ui.userEdited = userEdited
	ui.refreshGrid()
}

// ShowStatus displays a status message and clears it after a delay
func (ui *UI) ShowStatus(message string) {
	ui.statusText.SetText(message).SetTextColor(tcell.ColorYellow)
	
	// Schedule clearing the message after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		ui.app.QueueUpdateDraw(func() {
			ui.statusText.SetText("")
		})
	}()
}
