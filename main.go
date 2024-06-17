package main

import (
	"goduko/puzzle"
	"goduko/ui"
)

func main() {
	grid, _ := puzzle.LoadFromJSON("puzzles.json")
	ui.RenderPuzzle(grid)
}
