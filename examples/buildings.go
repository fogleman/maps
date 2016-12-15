package main

import (
	"github.com/fogleman/maps"
	"github.com/qedus/osmpbf"
)

func main() {
	dc := maps.NewMap(6000, 18000)
	dc.Center = maps.Point{-73.958, 40.785}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 280
	dc.Heading = 29
	dc.SetTransform()

	dc.SetHexColor("#FFFFFF")
	dc.Clear()

	dc.DrawShapefileFiltered(
		"manhattan/nybb_16d/cleaned.shp",
		maps.NewShapeTagFilter("BoroName", "Manhattan"))
	dc.SetHexColor("#D7DADB")
	// dc.FillPreserve()
	dc.Clip()

	// dc.DrawShapefile("manhattan/nycd_16d/cleaned.shp")
	// dc.SetHexColor("#000000")
	// dc.Stroke()

	pbf, _ := maps.LoadPBF("manhattan/manhattan.osm.pbf")
	seen := make(map[int64]bool)
	for _, relation := range pbf.Relations {
		building, ok := relation.Tags["building"]
		if !ok || building == "no" {
			continue
		}
		dc.DrawMultiPolygon(pbf, relation)
		for _, member := range relation.Members {
			if member.Type == osmpbf.WayType {
				seen[member.ID] = true
			}
		}
	}
	for _, way := range pbf.Ways {
		if seen[way.ID] {
			continue
		}
		if way.Tags["name"] == "The High Line" {
			continue
		}
		building, ok := way.Tags["building"]
		if !ok || building == "no" {
			continue
		}
		dc.DrawWay(pbf, way)
	}
	// dc.SetHexColor("#2C3E50")
	dc.SetRGB(0, 0, 0)
	dc.SetFillRuleEvenOdd()
	dc.Fill()
	dc.SetFillRuleWinding()
	dc.ResetClip()

	dc.SavePNG("out.png")
}
