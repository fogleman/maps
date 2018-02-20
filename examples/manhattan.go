package main

import (
	"fmt"
	"strconv"

	"github.com/fogleman/maps"
	"github.com/qedus/osmpbf"
)

var buildings = map[string]bool{
	"Empire State Building":              true,
	"Chrysler Building":                  true,
	"Statue of Liberty":                  true,
	"Saint Patrick's Old Cathedral":      true,
	"Grand Central Terminal":             true,
	"The Metropolitan Museum of Art":     true,
	"Cathedral of Saint John the Divine": true,
	"Trinity Church":                     true,
	"Woolworth Building":                 true,
	"Waldorf-Astoria Hotel":              true,
	"New york public library":            true,
	"Hearst Tower":                       true,
	"Flatiron Building":                  true,
	"Guggenheim Museum":                  true,
	"Plaza Hotel":                        true,
	"Radio City Music Hall":              true,
	"Time Warner Center":                 true,
	"United Nations Headquarters":        true,
	"Museum of Modern Art":               true,
	"Bloomberg Tower":                    true,
	"One World Trade Center":             true,
	"McSorley's Old Ale House":           true,
	"One Times Square":                   true,
	"Bank of America Tower":              true,
	"Three World Trade Center":           true,
	"The New York Times Building":        true,
	"One57":                                     true,
	"Four World Trade Center":                   true,
	"70 Pine Street":                            true,
	"Metropolitan Life Insurance Company Tower": true,
	"MetLife Building":                          true,
	"Citigroup Center":                          true,
	"Condé Nast Building":                       true,
}

func main() {
	dc := maps.NewMap(1024*6, 2048*6)
	dc.Center = maps.Point{-73.965, 40.779}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 250
	dc.Heading = 29
	dc.SetTransform()

	dc.SetHexColor("#FFFFFF")
	dc.Clear()

	dc.DrawShapefileFiltered(
		"manhattan/nybb_16d/cleaned.shp",
		maps.NewShapeTagFilter("BoroName", "Manhattan"))
	// dc.DrawShapefile("manhattan/nycd_16d/cleaned.shp")
	dc.SetHexColor("#D7DADB")
	dc.FillPreserve()
	dc.Clip()

	points := make(map[maps.Point]string)

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
		if buildings[relation.Tags["name"]] {
			bounds := maps.BoundsForRelation(pbf, relation)
			points[bounds.Center()] = relation.Tags["name"]
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
		if buildings[way.Tags["name"]] {
			bounds := maps.BoundsForWay(pbf, way)
			points[bounds.Center()] = way.Tags["name"]
		}
	}
	dc.SetHexColor("#2C3E50")
	dc.SetFillRuleEvenOdd()
	dc.Fill()
	dc.SetFillRuleWinding()

	dc.ResetClip()

	// for _, point := range points {
	// 	point = dc.Project(point)
	// 	dc.DrawPoint(point.X, point.Y, 12)
	// }
	// dc.SetHexColor("#FC4349")
	// dc.Fill()

	dc.LoadFontFace("/Library/Fonts/Arial Rounded Bold.ttf", 12)
	i := 0
	for point, name := range points {
		i++
		fmt.Println(name)
		point = dc.Project(point)
		s := strconv.Itoa(i)
		w, h := dc.MeasureString(s)

		dc.Push()
		x, y := dc.TransformPoint(point.X, point.Y)
		dc.Identity()
		dc.DrawRoundedRectangle(x-w/2-4, y-h/2-4, w+8, h+8, 4)
		dc.SetRGBA(1, 1, 1, 1)
		dc.Fill()
		dc.DrawRoundedRectangle(1800-w, y-h/2-4, w+8, h+8, 4)
		dc.SetRGBA(0, 0, 0, 1)
		// dc.SetHexColor("#FC434980")
		dc.Fill()
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(s, 1800-w+w/2+4, y, 0.5, 0.5)
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(name, 1820, y, 0, 0.5)
		dc.Pop()

		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(s, point.X, point.Y, 0.5, 0.5)
	}

	dc.SavePNG("out.png")
}
