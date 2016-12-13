package maps

import (
	"strings"

	"github.com/jonas-p/go-shp"
)

type Shape struct {
	shp.Shape
	Tags map[string]string
}

func LoadSHP(path string) ([]Shape, error) {
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

	var result []Shape
	for file.Next() {
		n, shape := file.Shape()
		tags := make(map[string]string)
		for i, name := range names {
			tags[name] = file.ReadAttribute(n, i)
		}
		result = append(result, Shape{shape, tags})
	}
	return result, nil
}

func (shape Shape) GetPoints() [][]Point {
	switch v := shape.Shape.(type) {
	case *shp.PolyLine:
		return getPoints(v)
	case *shp.Polygon:
		line := shp.PolyLine(*v)
		return getPoints(&line)
	}
	return nil
}

func getPoints(line *shp.PolyLine) [][]Point {
	var result [][]Point
	parts := append(line.Parts, line.NumPoints)
	for part := 0; part < len(parts)-1; part++ {
		var points []Point
		a := parts[part]
		b := parts[part+1]
		for i := a; i < b; i++ {
			point := line.Points[i]
			points = append(points, Point{point.X, point.Y})
		}
		result = append(result, points)
	}
	return result
}
