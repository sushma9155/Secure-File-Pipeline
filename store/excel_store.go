package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelStore struct {
	file     *excelize.File
	filePath string
	sheet    string
	rowCount int
}

// NewExcelStore creates an Excel file and sheet if they don't exist
func NewExcelStore(filePath string, sheet string) (*ExcelStore, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	var f *excelize.File
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist → create new
		f = excelize.NewFile()
		f.DeleteSheet("Sheet1") // remove default sheet
		if _, err := f.NewSheet(sheet); err != nil {
			return nil, fmt.Errorf("failed to create sheet: %w", err)
		}

		// Add headers
		f.SetCellValue(sheet, "A1", "Source")
		f.SetCellValue(sheet, "B1", "Filename")
		f.SetCellValue(sheet, "C1", "Data")
		f.SetCellValue(sheet, "D1", "CreatedAt")

		if err := f.SaveAs(filePath); err != nil {
			return nil, fmt.Errorf("cannot save new Excel file: %w", err)
		}
	} else {
		// Open existing file
		var err error
		f, err = excelize.OpenFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open Excel file: %w", err)
		}
	}

	// Count existing rows
	rows, _ := f.GetRows(sheet)
	return &ExcelStore{
		file:     f,
		filePath: filePath,
		sheet:    sheet,
		rowCount: len(rows),
	}, nil
}

// SaveRecord appends a row
func (s *ExcelStore) SaveRecord(source, filename, data, createdAt string) error {
	s.rowCount++
	row := s.rowCount
	s.file.SetCellValue(s.sheet, fmt.Sprintf("A%d", row), source)
	s.file.SetCellValue(s.sheet, fmt.Sprintf("B%d", row), filename)
	s.file.SetCellValue(s.sheet, fmt.Sprintf("C%d", row), data)
	s.file.SetCellValue(s.sheet, fmt.Sprintf("D%d", row), createdAt)
	return s.file.SaveAs(s.filePath)
}

// Close saves the file (optional)
func (s *ExcelStore) Close() error {
	return s.file.SaveAs(s.filePath)
}

func main() {
	filePath := "records.xlsx" // absolute path
	sheetName := "Records"

	excel, err := NewExcelStore(filePath, sheetName)
	if err != nil {
		log.Fatal(err)
	}
	defer excel.Close()

	// Example record
	err = excel.SaveRecord("HTTP", "data.json", `{"user":"Aryan"}`, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Record saved to Excel successfully!")
}
