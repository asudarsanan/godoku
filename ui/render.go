package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RenderPuzzle(puzzle [9][9]int) {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	for r, row := range puzzle {
		for c, col := range row {
			var text string
			if col == 0 {
				text = " " // Empty cell
			} else {
				text = fmt.Sprintf("%d", col)
			}
			color := tcell.ColorWhite
			table.SetCell(r, c, tview.NewTableCell(text).SetTextColor(color).SetAlign(tview.AlignCenter))
		}
	}

	// Run the application
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
