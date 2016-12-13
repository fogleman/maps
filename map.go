package maps

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/qedus/osmpbf"
)

type Map struct {
	*gg.Context
	Projection Projection
	Center     Point
	Zoom       float64
	Heading    float64
}

func NewMap(w, h int) *Map {
	m := Map{}
	m.Context = gg.NewContext(w, h)
	m.Projection = NewMercatorProjection()
	m.Center = Point{}
	m.Zoom = 1
	m.Heading = 0
	m.SetTransform()
	return &m
}

func (m *Map) SetTransform() {
	center := m.Project(m.Center)
	m.Identity()
	m.Translate(-center.X, -center.Y)
	m.Translate(float64(m.Width())/2, float64(m.Height())/2)
	m.Rotate(Radians(-m.Heading))
}

func (m *Map) FitBounds(bounds Bounds, margin float64) {
	m.Center = bounds.Center()
	m.Zoom = 1
	size := m.Project(bounds.Max).Sub(m.Project(bounds.Min))
	zx := (float64(m.Width()) - margin*2) / size.X
	zy := (float64(m.Height()) - margin*2) / size.Y
	m.Zoom = math.Abs(math.Min(zx, zy))
	m.SetTransform()
}

func (m *Map) Project(point Point) Point {
	scale := float64(m.Height()) * m.Zoom
	point = m.Projection.Project(point)
	point = Point{point.X * scale, point.Y * -scale}
	return point
}

func (m *Map) DrawShapefile(path string) error {
	shapes, err := LoadShapefile(path)
	if err != nil {
		return err
	}
	m.DrawShapes(shapes)
	return nil
}

func (m *Map) DrawShapefileFiltered(path string, filter ShapeFilterFunc) error {
	shapes, err := LoadShapefile(path)
	if err != nil {
		return err
	}
	for _, shape := range shapes {
		if filter(shape) {
			m.DrawShape(shape)
		}
	}
	return nil
}

func (m *Map) DrawShapes(shapes []Shape) {
	for _, shape := range shapes {
		m.DrawShape(shape)
	}
}

func (m *Map) DrawShape(shape Shape) {
	for _, line := range shape.Lines {
		m.NewSubPath()
		for _, pt := range line.Points {
			point := m.Project(Point{pt.X, pt.Y})
			m.LineTo(point.X, point.Y)
		}
	}
}

func (m *Map) DrawWay(pbf *PBF, way *osmpbf.Way) {
	m.NewSubPath()
	for _, id := range way.NodeIDs {
		node := pbf.Nodes[id]
		point := m.Project(Point{node.Lon, node.Lat})
		m.LineTo(point.X, point.Y)
	}
}

func (m *Map) DrawMultiPolygon(pbf *PBF, relation *osmpbf.Relation) {
	for _, member := range relation.Members {
		if member.Type != osmpbf.WayType {
			continue
		}
		if way, ok := pbf.Ways[member.ID]; ok {
			m.DrawWay(pbf, way)
		}
	}
}
