package maps

type Polyline struct {
	Points []Point
}

func NewPolyline(points []Point) *Polyline {
	return &Polyline{points}
}

func (p *Polyline) Bounds() Bounds {
	return BoundsForPoints(p.Points...)
}

func (p *Polyline) Length() float64 {
	var length float64
	for i := 1; i < len(p.Points); i++ {
		a := p.Points[i-1]
		b := p.Points[i]
		length += a.Distance(b)
	}
	return length
}
