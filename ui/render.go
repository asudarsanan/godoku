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
	app        *tview.Application
	table      *tview.Table
	game       [9][9]int
	userEdited [9][9]bool
}

func NewUI(game [9][9]int) *UI {
	return &UI{
		app:        tview.NewApplication(),
		table:      tview.NewTable().SetBorders(true),
		game:       game,
		userEdited: [9][9]bool{},
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
	return ui.app.SetRoot(ui.table, true).Run()
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
				color = tcell.ColorBlack
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
		row, col := ui.table.GetSelection()
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
				cell.SetText(fmt.Sprintf(" %s ", newText)).SetTextColor(tcell.Color20)
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
