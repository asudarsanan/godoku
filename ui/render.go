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
				text = "   " // Empty cell
			} else {
				text = fmt.Sprintf(" %d ", col)
			}
			color := tcell.ColorRed
			table.SetCell(r, c, tview.NewTableCell(text).SetTextColor(color).SetAlign(tview.AlignCenter))
		}
	}
	table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()

		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorWhite)
		table.SetSelectable(true, true)
	})

	// Run the application
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
