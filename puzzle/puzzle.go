package puzzle

import (
	"encoding/json"
	"godoku/config"
	"io"
	"os"
	"path/filepath"
)

func ImportPuzzle(filename string) ([9][9]int, error) {
	var grid [9][9]int

	// Determine file path - try config directory first, then fallback to resources
	puzzlePath := config.GetPuzzlePath(filename)
	
	// If the config directory doesn't exist or file is not there, try resources directory
	if _, err := os.Stat(puzzlePath); os.IsNotExist(err) {
		puzzlePath = filepath.Join("resources", filename)
	}

	// Open the file
	file, err := os.Open(puzzlePath)
	if err != nil {
		return grid, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Read the file content
	data, err := io.ReadAll(file)
	if err != nil {
		return grid, err
	}

	// Unmarshal JSON data into grid
	var jsonData struct {
		Grid [9][9]int `json:"grid"`
	}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return grid, err
	}
	//fmt.Print(jsonData.Grid)
	return jsonData.Grid, nil
}
