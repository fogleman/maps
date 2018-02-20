package maps

import "math"

type Projection interface {
	Project(point Point) Point
}

// Mercator

type MercatorProjection struct {
	InvertY bool
}

func NewMercatorProjection() Projection {
	return &MercatorProjection{false}
}

func (p *MercatorProjection) Project(point Point) Point {
	x := Radians(point.X)
	y := math.Asinh(math.Tan(Radians(point.Y)))
	if p.InvertY {
		y = -y
	}
	return Point{x, y}
}

// Lambert Azimuthal Equal Area

type LambertAzimuthalEqualAreaProjection struct {
	Center Point
}

func NewLambertAzimuthalEqualAreaProjection(center Point) Projection {
	return &LambertAzimuthalEqualAreaProjection{center}
}

func (p *LambertAzimuthalEqualAreaProjection) Project(point Point) Point {
	lng, lat := Radians(point.X), Radians(point.Y)
	clng, clat := Radians(p.Center.X), Radians(p.Center.Y)
	k := math.Sqrt(2 / (1 + math.Sin(clat)*math.Sin(lat) + math.Cos(clat)*math.Cos(lat)*math.Cos(lng-clng)))
	x := k * math.Cos(lat) * math.Sin(lng-clng)
	y := k * (math.Cos(clat)*math.Sin(lat) - math.Sin(clat)*math.Cos(lat)*math.Cos(lng-clng))
	return Point{x, y}
}
