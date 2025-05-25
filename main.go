package main

import (
	"flag"
	"fmt"
	"godoku/config"
	"godoku/logger"
	"godoku/puzzle"
	"godoku/ui"
	"os"
)

func main() {
	// Parse command line flags
	initFlag := flag.Bool("init", false, "Initialize configuration directories")
	builderFlag := flag.Bool("builder", false, "Launch the puzzle builder")
	flag.Parse()

	// Initialize configuration
	config.Init()
	
	// Initialize configuration if requested or on first run
	if *initFlag || !isConfigInitialized() {
		if *initFlag {
			fmt.Println("Initializing configuration directories...")
		} else {
			fmt.Println("First run detected, initializing configuration...")
		}
		
		if err := config.InitConfig(); err != nil {
			fmt.Printf("Error initializing configuration: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Configuration initialized successfully!")
		fmt.Printf("Puzzles directory: %s\n", config.PuzzlesDir)
		fmt.Printf("Saves directory: %s\n", config.SavesDir)
		fmt.Printf("Logs directory: %s\n", config.LogDir)
		
		// Exit if only initialization was requested
		if *initFlag && len(os.Args) == 2 {
			return
		}
	}

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %s\n", err)
	}
	defer logger.Close()
	
	logger.Info("Application started")
	// Launch the builder if requested
	if *builderFlag {
		logger.Info("Launching puzzle builder")
		fmt.Println("Launching puzzle builder...")
		
		// Create and run builder UI
		puzzleBuilder := ui.NewBuilder()
		if err := puzzleBuilder.Run(); err != nil {
			logger.Error("Error running puzzle builder: %v", err)
			fmt.Printf("Error running puzzle builder: %v\n", err)
		}
		return
	}

	// Launch the game UI
	logger.Info("Loading game")
	fmt.Println("Loading game...")
	
	// Try to load the most recent save game first
	gameState, success := puzzle.LoadLastGame()
	if success {
		fmt.Println("Found unfinished game, loading it...")
		// Create UI with loaded game state
		loadedGame := ui.NewUI(gameState.OriginalPuzzle)
		
		// Update with saved state
		loadedGame.SetGameState(gameState.CurrentState, gameState.UserEdited)
		loadedGame.ShowStatus("Last saved game loaded automatically")
		
		if err := loadedGame.Run(); err != nil {
			logger.Error("Error running UI for loaded game: %v", err)
			fmt.Printf("Error running UI for loaded game: %v\n", err)
		}
		return
	}
	
	// If no save game or failed to load, proceed with puzzle selection
	fmt.Println("No unfinished game found, selecting new puzzle...")
	
	// Show puzzle selector UI
	var grid [9][9]int
	selected := ui.ShowPuzzleSelector(func(puzzlePath string) bool {
		var err error
		grid, err = puzzle.ImportPuzzle(puzzlePath)
		if err != nil {
			logger.Error("Error loading puzzle: %v", err)
			fmt.Printf("Error loading puzzle: %v\n", err)
			return false
		}
		return true
	})
	
	if !selected {
		logger.Info("Puzzle selection canceled")
		fmt.Println("Puzzle selection canceled")
		return
	}
	
	// Create and run game UI
	newGame := ui.NewUI(grid)
	if err := newGame.Run(); err != nil {
		logger.Error("Error running UI: %v", err)
		fmt.Printf("Error running UI: %v\n", err)
	}
}

// isConfigInitialized checks if the configuration directories exist
func isConfigInitialized() bool {
	// Check if config directory exists
	if _, err := os.Stat(config.ConfigDir); os.IsNotExist(err) {
		return false
	}
	// Check if puzzles directory exists
	if _, err := os.Stat(config.PuzzlesDir); os.IsNotExist(err) {
		return false
	}
	// Check if saves directory exists
	if _, err := os.Stat(config.SavesDir); os.IsNotExist(err) {
		return false
	}
	return true
}

// This function is now replaced by the ShowPuzzleSelector function in ui/puzzle_selector.go
