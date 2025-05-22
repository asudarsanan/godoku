package ui

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"godoku/config"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Builder represents the puzzle builder interface
type Builder struct {
	app           *tview.Application
	table         *tview.Table
	grid          [9][9]int
	mainContainer *tview.Flex
	controlBar    *tview.Flex
	clearButton   *tview.Button
	saveButton    *tview.Button
	exitButton    *tview.Button
	statusText    *tview.TextView
	nameInput     *tview.InputField
}

// NewBuilder creates a new builder UI instance
func NewBuilder() *Builder {
	builder := &Builder{
		app:        tview.NewApplication(),
		table:      tview.NewTable().SetBorders(true),
		grid:       [9][9]int{},
		statusText: tview.NewTextView().SetText("").SetTextAlign(tview.AlignCenter),
		nameInput:  tview.NewInputField().SetLabel("Puzzle Name: ").SetFieldWidth(30),
	}

	// Create control buttons
	builder.clearButton = tview.NewButton("Clear All").
		SetSelectedFunc(func() {
			builder.clearPuzzle()
		})

	builder.saveButton = tview.NewButton("Save Puzzle").
		SetSelectedFunc(func() {
			builder.savePuzzle()
		})

	builder.exitButton = tview.NewButton("Exit").
		SetSelectedFunc(func() {
			builder.app.Stop()
		})

	// Create control bar
	builder.controlBar = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(builder.clearButton, 10, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(builder.saveButton, 12, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(builder.exitButton, 8, 0, true).
		AddItem(nil, 0, 1, false)

	// Create input container
	inputContainer := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(builder.nameInput, 50, 0, true).
		AddItem(nil, 0, 1, false)

	// Create main container with puzzle grid and control bar
	builder.mainContainer = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(builder.table, 37, 0, true).
				AddItem(nil, 0, 1, false),
			19, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(inputContainer, 3, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(builder.statusText, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(builder.controlBar, 1, 0, true).
		AddItem(nil, 0, 1, false)

	return builder
}

// Run starts the builder UI
func (b *Builder) Run() error {
	// Set up logging
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
		}
	}(logFile)
	log.SetOutput(logFile)

	b.initGrid()
	return b.app.SetRoot(b.mainContainer, true).EnableMouse(true).Run()
}

func (b *Builder) initGrid() {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			b.table.SetCell(r, c, tview.NewTableCell("   ").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetSelectable(true))
		}
	}

	b.table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyCtrlC {
			b.app.Stop()
		}
		if key == tcell.KeyEnter {
			b.table.SetSelectable(true, true)
		}
	})

	b.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Add Tab key to switch focus between puzzle grid and control buttons
		if event.Key() == tcell.KeyTab {
			if b.app.GetFocus() == b.table {
				b.app.SetFocus(b.nameInput)
			} else {
				b.app.SetFocus(b.clearButton)
			}
			return nil
		}

		row, col := b.table.GetSelection()
		
		// Check if row and col are valid indices
		if row < 0 || row >= 9 || col < 0 || col >= 9 {
			log.Printf("Invalid cell selection (%d, %d)", row, col)
			return event
		}
		
		cell := b.table.GetCell(row, col)

		switch event.Key() {
		case tcell.KeyRune:
			r := event.Rune()
			if r >= '1' && r <= '9' {
				num, _ := strconv.Atoi(string(r))
				b.grid[row][col] = num
				cell.SetText(fmt.Sprintf(" %d ", num)).SetTextColor(tcell.ColorAqua)
				log.Printf("Set cell (%d, %d) to value: %d", row, col, num)
			} else if r == '0' {
				b.grid[row][col] = 0
				cell.SetText("   ").SetTextColor(tcell.ColorWhite)
				log.Printf("Cleared cell (%d, %d)", row, col)
			}
		case tcell.KeyBackspace, tcell.KeyDelete:
			b.grid[row][col] = 0
			cell.SetText("   ").SetTextColor(tcell.ColorWhite)
			log.Printf("Cleared cell (%d, %d)", row, col)
		}
		return event
	})

	// Set up input capture for the main container to handle global key events
	b.mainContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Cycle focus between the table, input field, and buttons
			if b.app.GetFocus() == b.table {
				b.app.SetFocus(b.nameInput)
				return nil
			} else if b.app.GetFocus() == b.nameInput {
				b.app.SetFocus(b.clearButton)
				return nil
			} else {
				b.app.SetFocus(b.table)
				return nil
			}
		}
		return event
	})

	// Set up navigation between buttons
	b.clearButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			b.app.SetFocus(b.saveButton)
			return nil
		}
		return event
	})

	b.saveButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			b.app.SetFocus(b.clearButton)
			return nil
		} else if event.Key() == tcell.KeyRight {
			b.app.SetFocus(b.exitButton)
			return nil
		}
		return event
	})

	b.exitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			b.app.SetFocus(b.saveButton)
			return nil
		}
		return event
	})

	// Set selection behavior
	b.table.SetSelectable(true, true)
}

// clearPuzzle resets the grid to empty
func (b *Builder) clearPuzzle() {
	// Create a modal for confirmation
	modal := tview.NewModal().
		SetText("Are you sure you want to clear the puzzle? All cells will be emptied.").
		AddButtons([]string{"Clear", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Clear" {
				// Clear the grid
				b.grid = [9][9]int{}
				
				// Clear all cell displays
				for r := 0; r < 9; r++ {
					for c := 0; c < 9; c++ {
						cell := b.table.GetCell(r, c)
						cell.SetText("   ").SetTextColor(tcell.ColorWhite)
					}
				}
				
				log.Printf("Puzzle has been cleared")
				b.showStatus("Puzzle cleared")
			}
			
			// Return to the main UI
			b.app.SetRoot(b.mainContainer, true)
			b.app.SetFocus(b.table)
		})
	
	// Show the modal
	b.app.SetRoot(modal, true)
}

// savePuzzle saves the current puzzle to a file
func (b *Builder) savePuzzle() {
	// Check if any cells are filled
	isEmpty := true
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.grid[r][c] != 0 {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			break
		}
	}
	
	if isEmpty {
		b.showStatus("Cannot save an empty puzzle")
		return
	}
	
	// Get puzzle name
	puzzleName := b.nameInput.GetText()
	if puzzleName == "" {
		puzzleName = fmt.Sprintf("puzzle_%s", time.Now().Format("2006-01-02_15-04-05"))
	}
	
	// Add .json extension if not present
	if filepath.Ext(puzzleName) != ".json" {
		puzzleName += ".json"
	}
	
	// Create puzzles directory if it doesn't exist
	puzzlesDir := config.PuzzlesDir
	if _, err := os.Stat(puzzlesDir); os.IsNotExist(err) {
		err := os.MkdirAll(puzzlesDir, 0755)
		if err != nil {
			log.Printf("Error creating puzzles directory: %v", err)
			b.showStatus("Error: Could not create puzzles directory")
			return
		}
	}
	
	// Construct puzzle data
	puzzleData := struct {
		Grid [9][9]int `json:"grid"`
	}{
		Grid: b.grid,
	}
	
	// Marshal to JSON
	jsonData, err := json.MarshalIndent(puzzleData, "", "  ")
	if err != nil {
		log.Printf("Error marshalling puzzle data: %v", err)
		b.showStatus("Error: Could not save puzzle")
		return
	}
	
	// Save to file
	filePath := filepath.Join(puzzlesDir, puzzleName)
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing puzzle file: %v", err)
		b.showStatus("Error: Could not write puzzle file")
		return
	}
	
	log.Printf("Puzzle saved to %s", filePath)
	b.showStatus(fmt.Sprintf("Puzzle saved to %s", filepath.Base(filePath)))
}

// showStatus displays a status message and clears it after a delay
func (b *Builder) showStatus(message string) {
	b.statusText.SetText(message).SetTextColor(tcell.ColorYellow)
	
	// Schedule clearing the message after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		b.app.QueueUpdateDraw(func() {
			b.statusText.SetText("")
		})
	}()
}
