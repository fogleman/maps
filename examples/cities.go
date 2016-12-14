package main

import (
	"strconv"

	"github.com/fogleman/maps"
)

func main() {
	dc := maps.NewMap(2560, 1440)
	dc.Center = maps.Point{-98.35, 38.5}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 2
	dc.SetTransform()

	dc.SetHexColor("#2A2C2B")
	dc.Clear()

	dc.Push()
	dc.Translate(-8, 8)
	dc.DrawShapefile("files/cb_2015_us_state_500k.shp")
	dc.SetHexColor("#000000")
	dc.Fill()
	dc.Pop()

	dc.DrawShapefile("files/cb_2015_us_state_500k.shp")
	dc.SetHexColor("#D9CB9E")
	dc.Fill()

	records, err := maps.LoadCSV("files/cities.csv")
	if err != nil {
		panic(err)
	}

	dc.SetHexColor("#1E1E2080")
	for _, record := range records[1:] {
		if record[0] != "US" {
			continue
		}
		lat, _ := strconv.ParseFloat(record[7], 64)
		lng, _ := strconv.ParseFloat(record[8], 64)
		p := dc.Project(maps.Point{lng, lat})
		dc.DrawCircle(p.X, p.Y, 1)
		dc.Fill()
	}

	dc.SavePNG("out.png")
}
