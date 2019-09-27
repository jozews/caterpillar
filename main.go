

package main


import (
	"fmt"
)


func main() {

	lichess := Lichess{}
	lichess.setup()

	forWhite := true
	pieces := lichess.pieces()
	position := Position{pieces, [4]bool{true, true, true, true}, Coordinate{-1, -1}}

	go func() {
		selection := Coordinate{-1, -1}
		for {
			selection = lichess.observeSelection(selection)
			if !selection.isInBounds() {
				continue
			}				
			children := position.childrenForCoordinate(selection, forWhite)
			for _, child := range children {
				fmt.Println(child)
			}
			fmt.Printf("Number of moves %v", len(children))
		}
	}()
	
	for {
		// manual move 
		piecesNew := lichess.observePieces(position.pieces)
		position = position.childForPieces(piecesNew, forWhite)
		fmt.Println(position)
		forWhite = !forWhite

		// computer move 
		// child := position.childBest(forWhite)
		// coordinate1, coordinate2, promotion := position.coordinatesAndPromotionForChild(child)
		// lichess.makeMove(coordinate1, coordinate2, promotion, true)
		// position = child
		// forWhite = !forWhite
		// fmt.Println(position)
		
	}
}


