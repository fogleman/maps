package main

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"

	"github.com/fogleman/maps"
)

func main() {
	dc := maps.NewMap(1024, 1024)
	dc.Center = maps.Point{-119.509444, 37.229722}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 5
	dc.SetTransform()
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// load census tract populations from csv file
	file, _ := os.Open("examples/california.csv")
	defer file.Close()
	records, _ := csv.NewReader(file).ReadAll()
	population := make(map[string]string)
	for _, record := range records {
		population[record[0]] = record[1]
	}

	// render census tracts
	shapes, _ := maps.LoadShapefile("files/cb_2015_06_tract_500k.shp")
	for _, shape := range shapes {
		a, _ := strconv.Atoi(shape.Tags["ALAND"])
		p, _ := strconv.Atoi(population[shape.Tags["TRACTCE"]])
		d := 2589975.2356 * float64(p) / float64(a)
		t := math.Pow(d/10000, 0.5)
		dc.DrawShape(shape)
		dc.SetColor(maps.Viridis.Color(t))
		dc.FillPreserve()
		dc.Stroke()
	}

	// render county lines
	dc.DrawShapefileFiltered(
		"files/cb_2015_us_county_500k.shp",
		maps.NewShapeTagFilter("STATEFP", "06"))
	dc.SetHexColor("#888888")
	dc.SetLineWidth(0.5)
	dc.Stroke()

	// save output
	dc.SavePNG("out.png")
}
