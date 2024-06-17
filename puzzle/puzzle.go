package puzzle

import (
	"encoding/json"
	"io"
	"os"
)

func LoadFromJSON(filename string) ([9][9]int, error) {
	var grid [9][9]int

	// Open the file
	file, err := os.Open("resources/" + filename)
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
