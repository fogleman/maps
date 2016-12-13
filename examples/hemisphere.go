package main

import "github.com/fogleman/maps"

func main() {
	dc := maps.NewMap(1024, 1024)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.DrawShapefile("files/ne_50m_land/ne_50m_land.shp")
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SavePNG("out.png")
}
