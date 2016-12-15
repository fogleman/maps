package main

import (
	"github.com/fogleman/maps"
	"github.com/qedus/osmpbf"
)

func main() {
	dc := maps.NewMap(6000, 18000)
	dc.Center = maps.LatLng(40.785, -73.9585)
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 290
	dc.Heading = 29
	dc.SetTransform()

	dc.SetHexColor("#FFFFFF")
	dc.Clear()

	dc.DrawShapefileFiltered(
		"manhattan/nybb_16d/cleaned.shp",
		maps.NewShapeTagFilter("BoroName", "Manhattan"))
	dc.SetHexColor("#E3E3E3")
	dc.Fill()

	dc.DrawShapefileFiltered(
		"manhattan/nybbwi_16d/cleaned.shp",
		maps.NewShapeTagFilter("BoroName", "Manhattan"))
	dc.Clip()

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
		building, ok := way.Tags["building"]
		if !ok || building == "no" {
			continue
		}
		if way.Tags["name"] == "The High Line" {
			continue
		}
		dc.DrawWay(pbf, way)
	}
	dc.SetRGB(0, 0, 0)
	dc.SetFillRuleEvenOdd()
	dc.Fill()
	dc.SetFillRuleWinding()
	dc.ResetClip()

	dc.SavePNG("out.png")
}
