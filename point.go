package maps

import "math"

type Point struct {
	X, Y float64
}

func (a Point) Length() float64 {
	return math.Hypot(a.X, a.Y)
}

func (a Point) Distance(b Point) float64 {
	return a.Sub(b).Length()
}

func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

func (a Point) Sub(b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}

func (a Point) Mul(b Point) Point {
	return Point{a.X * b.X, a.Y * b.Y}
}

func (a Point) Div(b Point) Point {
	return Point{a.X / b.X, a.Y / b.Y}
}

func (a Point) Min(b Point) Point {
	return Point{math.Min(a.X, b.X), math.Min(a.Y, b.Y)}
}

func (a Point) Max(b Point) Point {
	return Point{math.Max(a.X, b.X), math.Max(a.Y, b.Y)}
}
