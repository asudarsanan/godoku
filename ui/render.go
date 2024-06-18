package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
)

type UI struct {
	app   *tview.Application
	table *tview.Table
	game  [9][9]int
}

func NewUI(game [9][9]int) *UI {

	return &UI{
		app:   tview.NewApplication(),
		table: tview.NewTable().SetBorders(true),
		game:  game,
	}
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
	//fmt.Sprintf("this is before return")
	return ui.app.SetRoot(ui.table, true).Run()
}

func (ui *UI) initGrid() {
	for r, row := range ui.game {
		for c, col := range row {
			var text string
			if col == 0 {
				text = "   " // Empty cell
				color := tcell.ColorWhite
				ui.table.SetCell(r, c, tview.NewTableCell(text).SetTextColor(color).SetAlign(tview.AlignCenter))

			} else {
				text = fmt.Sprintf(" %d ", col)
				color := tcell.ColorRed
				ui.table.SetCell(r, c, tview.NewTableCell(text).SetTextColor(color).SetAlign(tview.AlignCenter))
			}

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
	// Function to update cell content
	updateCell := func(r, c int) {
		if ui.table.GetCell(r, c).Text == " " {
			log.Printf("The cell is non-editable (%d, %d)", r, c)

		}
		inputField := tview.NewInputField()
		inputField.SetFieldStyle(tcell.Style{})
		inputField.
			SetLabel("New value: ").
			SetDoneFunc(func(key tcell.Key) {
				newText := inputField.GetText()
				if key == tcell.KeyEnter {
					ui.table.GetCell(r, c).SetText(newText)
					ui.app.SetRoot(ui.table, true).SetFocus(ui.table)
					log.Printf("Updated cell (%d, %d) with new value: %s", r, c, newText)
				} else if key == tcell.KeyEscape {
					ui.app.SetRoot(ui.table, true).SetFocus(ui.table)
					log.Printf("Update canceled for cell (%d, %d)", r, c)
				}
			})

		ui.app.SetRoot(inputField, true).SetFocus(inputField)
		log.Printf("Editing cell (%d, %d)", r, c)
	}

	// Set selection behavior
	ui.table.SetSelectable(true, true).
		SetSelectedFunc(func(r, c int) {
			updateCell(r, c)
		})
}

func (ui *UI) updateGrid(row int, col int, text string) {
	if len(text) == 0 {
		ui.game[row][col] = 0
	} else {
		ui.game[row][col] = int(text[0] - '0')
	}
}
