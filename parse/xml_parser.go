package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
)

// FlatObject is a generic map for XML decoding
type FlatObject map[string]string

// ParseXML parses a list of XML nodes into [][]string
func ParseXML(filePath string, recordTag string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var data []FlatObject
	var current FlatObject
	var inRecord bool
	var currentKey string

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("XML read error: %w", err)
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == recordTag {
				inRecord = true
				current = FlatObject{}
			} else if inRecord {
				currentKey = tok.Name.Local
			}
		case xml.CharData:
			if inRecord && currentKey != "" {
				current[currentKey] = string(tok)
			}
		case xml.EndElement:
			if tok.Name.Local == recordTag {
				inRecord = false
				data = append(data, current)
			} else if inRecord {
				currentKey = ""
			}
		}
	}

	if len(data) == 0 {
		return nil, nil
	}

	headers := getSortedKeysFromXML(data[0])
	result := [][]string{headers}

	for _, obj := range data {
		row := []string{}
		for _, key := range headers {
			row = append(row, obj[key])
		}
		result = append(result, row)
	}

	return result, nil
}

func getSortedKeysFromXML(obj map[string]string) []string {
	keys := []string{}
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
