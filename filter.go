package maps

type ShapeFilterFunc func(Shape) bool

func NewShapeTagFilter(tag, value string) ShapeFilterFunc {
	return func(shape Shape) bool {
		return shape.Tags[tag] == value
	}
}
