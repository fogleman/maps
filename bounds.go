package maps

import "github.com/qedus/osmpbf"

type Bounds struct {
	Min, Max Point
}

func BoundsForPoints(points ...Point) Bounds {
	if len(points) == 0 {
		return Bounds{}
	}
	min := points[0]
	max := points[0]
	for _, v := range points {
		min = min.Min(v)
		max = max.Max(v)
	}
	return Bounds{min, max}
}

func BoundsForPolylines(lines ...*Polyline) Bounds {
	if len(lines) == 0 {
		return Bounds{}
	}
	result := lines[0].Bounds()
	for _, line := range lines {
		result = result.Extend(line.Bounds())
	}
	return result
}

func BoundsForShapes(shapes ...Shape) Bounds {
	if len(shapes) == 0 {
		return Bounds{}
	}
	result := shapes[0].Bounds
	for _, shape := range shapes {
		result = result.Extend(shape.Bounds)
	}
	return result
}

func BoundsForWay(pbf *PBF, way *osmpbf.Way) Bounds {
	var points []Point
	for _, id := range way.NodeIDs {
		node := pbf.Nodes[id]
		points = append(points, Point{node.Lon, node.Lat})
	}
	return BoundsForPoints(points...)
}

func BoundsForRelation(pbf *PBF, relation *osmpbf.Relation) Bounds {
	var points []Point
	for _, member := range relation.Members {
		if member.Type == osmpbf.WayType {
			if way, ok := pbf.Ways[member.ID]; ok {
				bounds := BoundsForWay(pbf, way)
				points = append(points, bounds.Min)
				points = append(points, bounds.Max)
			}
		}
		if member.Type == osmpbf.NodeType {
			if node, ok := pbf.Nodes[member.ID]; ok {
				points = append(points, Point{node.Lon, node.Lat})
			}
		}
	}
	return BoundsForPoints(points...)
}

func (a Bounds) Extend(b Bounds) Bounds {
	return Bounds{a.Min.Min(b.Min), a.Max.Max(b.Max)}
}

func (a Bounds) Offset(dx, dy float64) Bounds {
	min := a.Min.Sub(Point{dx, dy})
	max := a.Max.Add(Point{dx, dy})
	return Bounds{min, max}
}

func (a Bounds) Contains(b Point) bool {
	return a.Min.X <= b.X && a.Max.X >= b.X &&
		a.Min.Y <= b.Y && a.Max.Y >= b.Y
}

func (a Bounds) Intersects(b Bounds) bool {
	return !(a.Min.X > b.Max.X || a.Max.X < b.Min.X || a.Min.Y > b.Max.Y ||
		a.Max.Y < b.Min.Y)
}

func (a Bounds) Anchor(anchor Point) Point {
	return a.Min.Add(a.Size().Mul(anchor))
}

func (a Bounds) Center() Point {
	return a.Anchor(Point{0.5, 0.5})
}

func (a Bounds) Size() Point {
	return a.Max.Sub(a.Min)
}

func (a Bounds) Area() float64 {
	s := a.Size()
	return s.X * s.Y
}
