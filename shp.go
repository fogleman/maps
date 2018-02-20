package maps

import (
	"strings"

	"github.com/jonas-p/go-shp"
)

type Shape struct {
	Bounds Bounds
	Lines  []*Polyline
	Tags   map[string]string
}

func LoadShapefile(path string, filters ...string) ([]Shape, error) {
	file, err := shp.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fields := file.Fields()
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = strings.Trim(field.String(), "\x00")
	}

	var shapes []Shape
	for file.Next() {
		n, shape := file.Shape()
		tags := make(map[string]string)
		for i, name := range names {
			tags[name] = file.ReadAttribute(n, i)
		}
		lines := getPolylines(shape)
		bounds := BoundsForPolylines(lines...)
		shapes = append(shapes, Shape{bounds, lines, tags})
	}

	if len(filters) > 0 {
		filteredShapes := shapes[:0]
		for _, shape := range shapes {
			ok := true
			for i := 0; i < len(filters); i += 2 {
				key := filters[i]
				value := filters[i+1]
				if shape.Tags[key] != value {
					ok = false
					break
				}
			}
			if ok {
				filteredShapes = append(filteredShapes, shape)
			}
		}
		shapes = filteredShapes
	}

	return shapes, nil
}

func getPolylines(shape shp.Shape) []*Polyline {
	var line *shp.PolyLine
	switch v := shape.(type) {
	case *shp.PolyLine:
		line = v
	case *shp.Polygon:
		l := shp.PolyLine(*v)
		line = &l
	default:
		return nil
	}
	var result []*Polyline
	parts := append(line.Parts, line.NumPoints)
	for part := 0; part < len(parts)-1; part++ {
		var points []Point
		a := parts[part]
		b := parts[part+1]
		for i := a; i < b; i++ {
			point := line.Points[i]
			points = append(points, Point{point.X, point.Y})
		}
		result = append(result, NewPolyline(points))
	}
	return result
}
