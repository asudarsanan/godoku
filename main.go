package main

import (
	"goduko/puzzle"
	"goduko/ui"
	"log"
)

func main() {
	grid, _ := puzzle.ImportPuzzle("puzzles.json")
	newGame := ui.NewUI(grid)
	// Run the UI

	if err := newGame.Run(); err != nil {
		log.Fatalf("Error running UI: %v", err)
	}
}
