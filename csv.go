package maps

import (
	"encoding/csv"
	"os"
)

func LoadCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}
