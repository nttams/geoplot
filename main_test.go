package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func Test_GenerateRandomLatLonCSV(t *testing.T) {
	t.Skip("")
	err := generateRandomLatLonCSV("latlon.csv", 1_000)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func generateRandomLatLonCSV(filename string, count int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for i := 0; i < count; i++ {
		lat := -90 + rand.Float64()*180
		lon := -180 + rand.Float64()*360
		if lon < 0 {
			i -= 1
			continue
		}
		record := []string{
			strconv.FormatFloat(lat, 'f', 6, 64),
			strconv.FormatFloat(lon, 'f', 6, 64),
		}
		writer.Write(record)
	}
	return nil
}
