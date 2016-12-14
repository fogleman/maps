package main

import "github.com/fogleman/maps"

func main() {
	dc := maps.NewMap(1024*2, 2048*2)
	dc.Center = maps.Point{-73.965, 40.779}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 250
	dc.Heading = 29
	dc.SetTransform()

	dc.SetHexColor("#374140")
	dc.Clear()

	dc.DrawShapefileFiltered(
		"manhattan/nybb_15b/cleaned.shp",
		maps.NewShapeTagFilter("BoroName", "Manhattan"))
	dc.SetHexColor("#D9CB9E")
	dc.FillPreserve()
	dc.Clip()

	pbf, _ := maps.LoadPBF("manhattan/manhattan.osm.pbf")
	for _, way := range pbf.Ways {
		building := way.Tags["building"]
		if building == "" || building == "no" {
			continue
		}
		dc.DrawWay(pbf, way)
	}
	dc.SetHexColor("#000000")
	dc.Fill()

	dc.SavePNG("out.png")
}
