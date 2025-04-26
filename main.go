package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func main() {
	http.HandleFunc("/geo/", func(w http.ResponseWriter, r *http.Request) {
		files := []string{}
		for i := range 100 {
			if file := r.URL.Query().Get(fmt.Sprintf("geofile%d", i)); file != "" {
				files = append(files, file)
			} else {
				break
			}
		}
		renderGeo(w, files...)
	})

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderGeo(w http.ResponseWriter, files ...string) {
	pointGroup := make(map[string][]opts.GeoData)
	for _, file := range files {
		points := loadPoints(file)
		if len(points) <= 0 {
			continue
		}
		pointGroup[file] = points
		slog.Info("Loaded", "file", file, "count", len(points))
	}

	geo := charts.NewGeo()
	geo.SetGlobalOptions(
		charts.WithGeoComponentOpts(opts.GeoComponent{
			Map:    "world",
			Silent: opts.Bool(true),
		}),
	)

	colors := []string{"red", "blue", "black"}
	for i, file := range files {
		geo.AddSeries(file, types.ChartScatter, pointGroup[file], func(s *charts.SingleSeries) {
			s.SymbolSize = 1
			s.Color = colors[i%len(colors)]
		})
	}
	page := components.NewPage()
	page.AddCharts(geo)
	page.Render(w)
}

// format: ...,latitude,longitude
// the first field that can be parsed as a valid latitude will be treated as the latitude,
// and the next one will be longitude
func loadPoints(filepath string) []opts.GeoData {
	points := []opts.GeoData{}
	file, _ := os.Open(filepath)
	scanner := bufio.NewScanner(file)
	latIndex := -1

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		// Detect which column is lat
		if latIndex == -1 {
			for i := range parts {
				if lat, err := strconv.ParseFloat(parts[i], 64); err == nil && lat >= -90 && lat <= 90 {
					latIndex = i
					break
				}
			}
			if latIndex < 0 || latIndex >= len(parts)-1 {
				continue
			}
		}
		lat, err := strconv.ParseFloat(parts[latIndex], 64)
		if err != nil {
			continue
		}
		lng, err := strconv.ParseFloat(parts[latIndex+1], 64)
		if err != nil {
			continue
		}
		points = append(points, opts.GeoData{Value: []float64{lng, lat}})
	}

	sort.Slice(points, func(i, j int) bool {
		v1, _ := points[i].Value.([]float64)
		v2, _ := points[j].Value.([]float64)
		return v1[1] < v2[1]
	})
	return points
}
