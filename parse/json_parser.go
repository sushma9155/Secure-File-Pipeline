package parser

import (
	"encoding/json"
	"fmt"
	// "io"
	"os"
	"sort"
)

// ParseJSON parses a JSON array of objects into [][]string
func ParseJSON(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON: %w", err)
	}
	defer file.Close()

	var objects []map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&objects); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(objects) == 0 {
		return nil, nil
	}

	// Use sorted keys from the first object as headers
	headers := getSortedKeys(objects[0])
	data := [][]string{headers}

	for _, obj := range objects {
		row := []string{}
		for _, key := range headers {
			val := fmt.Sprintf("%v", obj[key])
			row = append(row, val)
		}
		data = append(data, row)
	}

	return data, nil
}

func getSortedKeys(m map[string]interface{}) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
