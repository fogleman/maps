package maps

import (
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
	m.Center = Point{-98.35, 39.5}
	m.Projection = NewLambertAzimuthalEqualAreaProjection(m.Center)
	m.Zoom = 0
	m.Heading = 0
	m.setTransform()
	return &m
}

func (m *Map) setTransform() {
	m.Push()
	center := m.Project(m.Center)
	m.Translate(-center.X, -center.Y)
	m.Translate(float64(m.Width())/2, float64(m.Height())/2)
}

func (m *Map) Project(point Point) Point {
	const scale = 600
	point = m.Projection.Project(point)
	point = Point{point.X * scale, point.Y * -scale}
	return point
}

func (m *Map) DrawShapefile(path string) error {
	shapes, err := LoadSHP(path)
	if err != nil {
		return err
	}
	m.DrawShapes(shapes)
	return nil
}

func (m *Map) DrawShapes(shapes []Shape) {
	for _, shape := range shapes {
		m.DrawShape(shape)
	}
}

func (m *Map) DrawShape(shape Shape) {
	groups := shape.GetPoints()
	for _, group := range groups {
		m.NewSubPath()
		for _, pt := range group {
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