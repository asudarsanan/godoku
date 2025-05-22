package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"strconv"
)

type UI struct {
	app           *tview.Application
	table         *tview.Table
	game          [9][9]int
	userEdited    [9][9]bool
	mainContainer *tview.Flex
	controlBar    *tview.Flex
	resetButton   *tview.Button
	saveButton    *tview.Button
	exitButton    *tview.Button
}

func NewUI(game [9][9]int) *UI {
	ui := &UI{
		app:        tview.NewApplication(),
		table:      tview.NewTable().SetBorders(true),
		game:       game,
		userEdited: [9][9]bool{},
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
	
	ui.exitButton = tview.NewButton("Exit").
		SetSelectedFunc(func() {
			ui.app.Stop()
		})
	
	// Create control bar
	ui.controlBar = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(ui.resetButton, 10, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(ui.saveButton, 10, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(ui.exitButton, 10, 0, true).
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
		AddItem(ui.controlBar, 1, 0, true).
		AddItem(nil, 0, 1, false)
	
	return ui
}

func (ui *UI) Run() error {
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

	ui.initGrid()
	return ui.app.SetRoot(ui.mainContainer, true).EnableMouse(true).Run()
}

func (ui *UI) initGrid() {
	for r, row := range ui.game {
		for c, col := range row {
			var text string
			if col == 0 {
				text = "   " // Empty cell
			} else {
				text = fmt.Sprintf(" %d ", col)
			}

			color := tcell.ColorYellowGreen
			immutable := col != 0
			if immutable {
				color = tcell.ColorAqua
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
			log.Printf("Invalid cell selection (%d, %d)", row, col)
			return event
		}
		
		cell := ui.table.GetCell(row, col)

		if !ui.userEdited[row][col] && ui.game[row][col] != 0 {
			log.Printf("The cell is non-editable (%d, %d)", row, col)
			return event
		}

		switch event.Key() {
		case tcell.KeyRune:
			r := event.Rune()
			if r >= '1' && r <= '9' {
				newText := string(r)
				cell.SetText(fmt.Sprintf(" %s ", newText)).SetTextColor(tcell.ColorRed)
				ui.updateGrid(row, col, newText)
				log.Printf("Updated cell (%d, %d) with new value: %s", row, col, newText)
			} else if r == '0' {
				cell.SetText("   ").SetTextColor(tcell.Color20)
				ui.updateGrid(row, col, "")
				log.Printf("Cleared cell (%d, %d)", row, col)
			}
		case tcell.KeyBackspace, tcell.KeyDelete:
			cell.SetText("   ").SetTextColor(tcell.Color20)
			ui.updateGrid(row, col, "")
			log.Printf("Cleared cell (%d, %d)", row, col)
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
			ui.app.SetFocus(ui.exitButton)
			return nil
		}
		return event
	})

	ui.exitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			ui.app.SetFocus(ui.saveButton)
			return nil
		}
		return event
	})

	// Set selection behavior
	ui.table.SetSelectable(true, true)
}

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
	// Clear all user inputs
	ui.userEdited = [9][9]bool{}
	
	// Refresh all cells to their original state
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			cell := ui.table.GetCell(r, c)
			if ui.game[r][c] == 0 {
				cell.SetText("   ").SetTextColor(tcell.ColorYellowGreen)
			} else {
				cell.SetText(fmt.Sprintf(" %d ", ui.game[r][c])).
					SetTextColor(tcell.ColorAqua).
					SetSelectable(false)
			}
		}
	}
	
	log.Printf("Puzzle has been reset to its original state")
	ui.app.SetFocus(ui.table)
}

// savePuzzle saves the current puzzle state (placeholder for future implementation)
func (ui *UI) savePuzzle() {
	// Here we would implement actual saving logic in future iterations
	log.Printf("Puzzle state saved (placeholder)")
	ui.app.SetFocus(ui.table)
}
