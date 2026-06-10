package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Save traceroute results to a JSON file.
func ExportJSON(filename string, data any) error {

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Saved results to %s\n", filename)

	return nil
}
