package main

import "github.com/fogleman/maps"

func main() {
	dc := maps.NewMap(2560, 1440)
	dc.Center = maps.Point{-98.35, 39.5}
	dc.Projection = maps.NewLambertAzimuthalEqualAreaProjection(dc.Center)
	dc.Zoom = 1.5
	dc.SetTransform()

	dc.SetHexColor("#374140")
	dc.Clear()

	dc.DrawShapefile("files/ne_10m_land.shp")
	dc.SetHexColor("#D9CB9E")
	dc.Fill()

	dc.DrawShapefile("files/ne_10m_lakes.shp")
	dc.SetHexColor("#374140")
	dc.Fill()

	dc.DrawShapefile("files/ne_10m_admin_1_states_provinces_lines_shp.shp")
	dc.SetHexColor("#1E1E20")
	dc.SetLineWidth(2)
	dc.Stroke()

	dc.DrawShapefile("files/ne_10m_admin_0_boundary_lines_land.shp")
	dc.SetHexColor("#1E1E20")
	dc.SetLineWidth(4)
	dc.Stroke()

	dc.SavePNG("out.png")
}
