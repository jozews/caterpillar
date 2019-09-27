
package main

// Coordinate ...
type Coordinate struct {
	x int
	y int
}

func (coordinate Coordinate) isInBounds() bool {
	return coordinate.x >= 0 && coordinate.x <= 7 && coordinate.y >= 0 && coordinate.y <= 7
}
