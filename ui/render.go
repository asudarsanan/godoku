package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
			} else {
				text = fmt.Sprintf(" %d ", col)
			}
			color := tcell.ColorRed
			ui.table.SetCell(r, c, tview.NewTableCell(text).SetTextColor(color).SetAlign(tview.AlignCenter))
		}
	}
	ui.table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			ui.app.Stop()

		}
		if key == tcell.KeyEnter {
			ui.table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		ui.table.GetCell(row, column).SetTextColor(tcell.ColorWhite)
		ui.table.SetSelectable(true, true)
	})
}

func (ui *UI) updateGrid(row int, col int, text string) {
	if len(text) == 0 {
		ui.game[row][col] = 0
	} else {
		ui.game[row][col] = int(text[0] - '0')
	}
}
