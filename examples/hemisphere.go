package main

import "github.com/fogleman/maps"

func main() {
	dc := maps.NewMap(2560, 1440)
	dc.SetHexColor("#A4CDFD")
	dc.Clear()

	dc.DrawShapefile("files/ne_10m_land/ne_10m_land.shp")
	dc.SetHexColor("#F3F1ED")
	dc.Fill()

	dc.DrawShapefile("files/ne_10m_lakes/ne_10m_lakes.shp")
	dc.SetHexColor("#A4CDFD")
	dc.Fill()

	// dc.DrawShapefile("files/ne_10m_rivers_lake_centerlines/ne_10m_rivers_lake_centerlines.shp")
	// dc.DrawShapefile("files/ne_10m_rivers_north_america/ne_10m_rivers_north_america.shp")
	// dc.SetHexColor("10222B")
	// dc.SetLineWidth(1)
	// dc.Stroke()

	dc.DrawShapefile("files/cb_2015_us_county_500k/cb_2015_us_county_500k.shp")
	dc.SetHexColor("#999999")
	dc.SetLineWidth(0.25)
	dc.Stroke()

	dc.DrawShapefile("files/cb_2015_us_state_500k/cb_2015_us_state_500k.shp")
	dc.SetHexColor("#333333")
	dc.SetLineWidth(2)
	dc.Stroke()

	dc.SavePNG("out.png")
}
