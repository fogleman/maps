package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/fogleman/maps"
)

func main() {
	dc := maps.NewMap(2560, 1440)
	dc.Center = maps.Point{-98.35, 38.5}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 2
	dc.SetTransform()
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// load census tract populations from csv file
	type TractKey struct {
		State string
		Tract string
	}

	population := make(map[TractKey]string)

	file, _ := os.Open("files/B01003_001E.csv")
	defer file.Close()
	records, _ := csv.NewReader(file).ReadAll()
	for _, record := range records {
		key := TractKey{record[1], record[3]}
		population[key] = record[0]
	}

	// render census tracts
	states := []string{"01", "02", "04", "05", "06", "08", "09", "10", "11", "12", "13", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "44", "45", "46", "47", "48", "49", "50", "51", "53", "54", "55", "56", "60", "66", "69", "72", "78"}
	for _, state := range states {
		path := fmt.Sprintf("tract/cb_2015_%s_tract_500k.shp", state)
		fmt.Println(path)
		shapes, _ := maps.LoadShapefile(path)
		for _, shape := range shapes {
			key := TractKey{shape.Tags["STATEFP"], shape.Tags["TRACTCE"]}
			a, _ := strconv.Atoi(shape.Tags["ALAND"])
			p, _ := strconv.Atoi(population[key])
			d := 2589975.2356 * float64(p) / float64(a)
			t := math.Pow(d/5000, 0.5)
			dc.DrawShape(shape)
			dc.SetColor(maps.Viridis.Color(t))
			dc.FillPreserve()
			dc.Stroke()
		}
	}

	// render county lines
	dc.DrawShapefile("files/cb_2015_us_state_500k.shp")
	dc.SetHexColor("#888888")
	dc.SetLineWidth(1)
	dc.Stroke()

	// save output
	dc.SavePNG("out.png")
}
