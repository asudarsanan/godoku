package main

import (
	"flag"
	"fmt"
	"godoku/builder"
	"godoku/config"
	"godoku/puzzle"
	"godoku/ui"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Set up logging
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %s\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Parse command line flags
	initFlag := flag.Bool("init", false, "Initialize configuration directories")
	builderFlag := flag.Bool("builder", false, "Launch the puzzle builder")
	flag.Parse()

	// Initialize configuration if requested or on first run
	if *initFlag {
		fmt.Println("Initializing configuration directories...")
		if err := config.InitConfig(); err != nil {
			fmt.Printf("Error initializing configuration: %v\n", err)
			log.Fatalf("Error initializing configuration: %v", err)
			return
		}
		fmt.Println("Configuration initialized successfully!")
		fmt.Printf("Puzzles directory: %s\n", config.PuzzlesDir)
		fmt.Printf("Saves directory: %s\n", config.SavesDir)
		
		// Exit if only initialization was requested
		if len(os.Args) == 2 && os.Args[1] == "--init" {
			return
		}
	}

	// Ensure config is initialized even if --init flag wasn't provided
	if _, err := os.Stat(config.ConfigDir); os.IsNotExist(err) {
		fmt.Println("First run detected, initializing configuration...")
		if err := config.InitConfig(); err != nil {
			fmt.Printf("Error initializing configuration: %v\n", err)
			log.Fatalf("Error initializing configuration: %v", err)
			return
		}
	}

	// Launch the builder if requested
	if *builderFlag {
		fmt.Println("Launching puzzle builder...")
		puzzleBuilder := builder.NewBuilder()
		if err := puzzleBuilder.Run(); err != nil {
			log.Fatalf("Error running puzzle builder: %v", err)
		}
		return
	}

	// Launch the game UI
	fmt.Println("Loading game...")
	
	// Get list of available puzzles
	puzzlesDir := config.PuzzlesDir
	puzzleFile := "puzzles.json"  // Default puzzle
	
	// Check if puzzles directory exists and has puzzle files
	if _, err := os.Stat(puzzlesDir); !os.IsNotExist(err) {
		puzzleFiles, err := os.ReadDir(puzzlesDir)
		if err == nil && len(puzzleFiles) > 0 {
			// If we have puzzle files, prompt the user to select one
			if len(puzzleFiles) > 1 {
				fmt.Println("Available puzzles:")
				for i, file := range puzzleFiles {
					if filepath.Ext(file.Name()) == ".json" {
						fmt.Printf("%d. %s\n", i+1, file.Name())
					}
				}
				
				fmt.Print("Select a puzzle (enter number or press Enter for default): ")
				var choice string
				fmt.Scanln(&choice)
				
				if choice != "" {
					var index int
					if _, err := fmt.Sscanf(choice, "%d", &index); err == nil && index > 0 && index <= len(puzzleFiles) {
						puzzleFile = puzzleFiles[index-1].Name()
					}
				}
			} else if filepath.Ext(puzzleFiles[0].Name()) == ".json" {
				// If only one puzzle file exists, use it
				puzzleFile = puzzleFiles[0].Name()
			}
		}
	}
	
	grid, err := puzzle.ImportPuzzle(puzzleFile)
	if err != nil {
		fmt.Printf("Error loading puzzle: %v\n", err)
		log.Fatalf("Error loading puzzle: %v", err)
		return
	}
	
	newGame := ui.NewUI(grid)
	if err := newGame.Run(); err != nil {
		log.Fatalf("Error running UI: %v", err)
	}
}
