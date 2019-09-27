
package main

// Vector ...
type Vector struct {
	dx     int
	dy     int
	length int
}


func (vector Vector) coordinatesFromCoordinate(coordinateInitial Coordinate) []Coordinate {

	var coordinates []Coordinate
	coordinate := coordinateInitial

	for i := 0; i < vector.length; i++ {
		coordinateNew := Coordinate{coordinate.x + vector.dx, coordinate.y + vector.dy} 
		if !coordinateNew.isInBounds() {
			break
		}
		coordinates = append(coordinates, coordinateNew)
		coordinate = coordinateNew
	}

	return coordinates
}