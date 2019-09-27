

package main


import (
	"math/rand"
	"strings"
	"unicode"
	"fmt"
)


// Position ...
type Position struct {
	pieces          map[Coordinate]Piece
	rightsToCastle  [4]bool
	squareEnPassant Coordinate
}


func (position Position) pieceAtCoordinate(coordinate Coordinate) Piece {
	return position.pieces[coordinate]
}


func (position Position) children(forWhite bool) []Position {

	var positions []Position

	for coordinateInitial, piece := range position.pieces {

		if !piece.isNotNil() {
			continue
		}

		if piece.isWhite() != forWhite {
			continue
		}

		var vectorsPiece []Vector
		if piece.isKing() {
			vectorsPiece = vectorsKing()

		} else if piece.isQueen() {
			vectorsPiece = vectorsQueen()

		} else if piece.isRook() {
			vectorsPiece = vectorsRook()

		} else if piece.isBishop() {
			vectorsPiece = vectorsBishop()

		} else if piece.isKnight() {
			vectorsPiece = vectorsKnight()
		}

		for _, vector := range vectorsPiece {

			coordinatesInVector := vector.coordinatesFromCoordinate(coordinateInitial)

			for _, coordinateInVector := range coordinatesInVector {

				pieceInVector := position.pieceAtCoordinate(coordinateInVector)

				isPiece := pieceInVector.isNotNil()
				isPieceOfSameColor := pieceInVector.isNotNil() && pieceInVector.isWhite() == piece.isWhite()

				if isPieceOfSameColor {
					break
				}

				piecesPositionChild := make(map[Coordinate]Piece)
				for key, value := range position.pieces {
					piecesPositionChild[key] = value
				}

				delete(piecesPositionChild, coordinateInitial)
				piecesPositionChild[coordinateInVector] = piece

				rightsToCastle := position.rightsToCastle

				rightsToCastle[0] = rightsToCastle[0] && coordinateInitialOfKing(true) != coordinateInitial && coordinateInitialOfRook(true, true) != coordinateInitial
				rightsToCastle[1] = rightsToCastle[1] && coordinateInitialOfKing(true) != coordinateInitial && coordinateInitialOfRook(true, false) != coordinateInitial
				rightsToCastle[2] = rightsToCastle[2] && coordinateInitialOfKing(false) != coordinateInitial && coordinateInitialOfRook(false, true) != coordinateInitial
				rightsToCastle[3] = rightsToCastle[3] && coordinateInitialOfKing(false) != coordinateInitial && coordinateInitialOfRook(false, false) != coordinateInitial

				positionChild := Position{piecesPositionChild, rightsToCastle, Coordinate{-1, -1}}
				positions = append(positions, positionChild)

				if isPiece {
					break
				}
			}
		}

		if piece.isKing() && !position.coordinateIsChecked(coordinateInitial, !piece.isWhite()) {

			var vectorsCastle []Vector
			if (piece.isWhite() && position.rightsToCastle[1]) || (!piece.isWhite() && position.rightsToCastle[3]) {
				vectorsCastle = append(vectorsCastle, vectorCastle(true))
			}
			if (piece.isWhite() && position.rightsToCastle[0]) || (!piece.isWhite() && position.rightsToCastle[2]) {
				vectorsCastle = append(vectorsCastle, vectorCastle(false))
			}

			for _, vector := range vectorsCastle {

				canCastle := true
				coordinatesInVector := vector.coordinatesFromCoordinate(coordinateInitial)
				
				for _, coordinateInVector := range coordinatesInVector {
					pieceAtSquareOfVector := position.pieceAtCoordinate(coordinateInVector)
					if pieceAtSquareOfVector.isNotNil() {
						canCastle = false
						break
					}
					if position.coordinateIsChecked(coordinateInVector, !piece.isWhite()) {
						canCastle = false
						break
					}
				}

				if canCastle {

					piecesPositionChild := make(map[Coordinate]Piece)
					for key, value := range position.pieces {
						piecesPositionChild[key] = value
					}

					delete(piecesPositionChild, coordinateInitial)
					piecesPositionChild[coordinatesInVector[1]] = piece

					coordinateOfRook := coordinateInitialOfRook(piece.isWhite(), vector.dx > 0)

					pieceRook := position.pieces[coordinateOfRook]
					delete(piecesPositionChild, coordinateOfRook)
					piecesPositionChild[coordinatesInVector[0]] = pieceRook

					rightsToCastle := position.rightsToCastle

					rightsToCastle[0] = !piece.isWhite()
					rightsToCastle[1] = !piece.isWhite()
					rightsToCastle[2] = piece.isWhite()
					rightsToCastle[3] = piece.isWhite()

					positionChild := Position{piecesPositionChild, rightsToCastle, Coordinate{-1, -1}}
					positions = append(positions, positionChild)
				}
			}
		}

		if piece.isPawn() {

			yInitialPawn := yInitialPawn(forWhite)
			yFinalPawn := yFinalPawn(forWhite)

			vectorsPawnCapture := vectorsPawnCapture(forWhite)

			for _, vector := range vectorsPawnCapture {

				coordinatesOfVector := vector.coordinatesFromCoordinate(coordinateInitial)

				for _, coordinateOfVector := range coordinatesOfVector {

					pieceInVector := position.pieceAtCoordinate(coordinateOfVector)
					isWhitePieceAtSquareOfVector := pieceInVector.isWhite()
					isPiece := pieceInVector.isNotNil()
					isPieceOfOppositeColor := isPiece && isWhitePieceAtSquareOfVector != piece.isWhite()
					isEnPassant := position.squareEnPassant == coordinateOfVector
					if !isPieceOfOppositeColor && !isEnPassant {
						break
					}

					piecesPositionChild := make(map[Coordinate]Piece)
					for key, value := range position.pieces {
						piecesPositionChild[key] = value
					}

					delete(piecesPositionChild, coordinateInitial)
					piecesPositionChild[coordinateOfVector] = piece

					if isEnPassant {
						coordinateOfCapture := Coordinate{coordinateOfVector.x, coordinateOfVector.y - dyPawn(forWhite)}
						delete(piecesPositionChild, coordinateOfCapture)
					}

					if yFinalPawn == coordinateOfVector.y {
						for _, piece := range piecesPromote(forWhite) {
							piecesPositionChildPromote := make(map[Coordinate]Piece)
							for key, value := range piecesPositionChild {
								piecesPositionChildPromote[key] = value
							}
							piecesPositionChildPromote[coordinateOfVector] = piece
							positionChild := Position{piecesPositionChildPromote, position.rightsToCastle, Coordinate{-1, -1}}
							positions = append(positions, positionChild)
						}
					} else {
						positionChild := Position{piecesPositionChild, position.rightsToCastle, Coordinate{-1, -1}}
						positions = append(positions, positionChild)
					}
				}
			}

			vectorPawnAdvance := vectorPawnAdvance(forWhite)
			coordinatesOfVector := vectorPawnAdvance.coordinatesFromCoordinate(coordinateInitial)

			for index, coordinateOfVector := range coordinatesOfVector {

				pieceAtVector := position.pieceAtCoordinate(coordinateOfVector)

				isPiece := pieceAtVector.isNotNil()

				if isPiece {
					break
				}

				if index == 1 && coordinateInitial.y != yInitialPawn {
					break
				}

				piecesPositionChild := make(map[Coordinate]Piece)
				for key, value := range position.pieces {
					piecesPositionChild[key] = value
				}

				delete(piecesPositionChild, coordinateInitial)
				piecesPositionChild[coordinateOfVector] = piece

				var coordinateEnPassant Coordinate
				if index == 1 {
					coordinateEnPassant = coordinatesOfVector[0]
				} else {
					coordinateEnPassant = Coordinate{-1, -1}
				}

				if yFinalPawn == coordinateOfVector.y {
					for _, piece := range piecesPromote(forWhite) {
						piecesPositionChildPromote := make(map[Coordinate]Piece)
						for key, value := range piecesPositionChild {
							piecesPositionChildPromote[key] = value
						}
						piecesPositionChildPromote[coordinateOfVector] = piece
						positionChild := Position{piecesPositionChildPromote, position.rightsToCastle, coordinateEnPassant}
						positions = append(positions, positionChild)
					}
				} else {
					positionChild := Position{piecesPositionChild, position.rightsToCastle, coordinateEnPassant}
					positions = append(positions, positionChild)
				}
			}
		}
	}

	var positionsLegal []Position

	for _, position := range positions {
		king := pieceKing(forWhite)
		coordinateOfKing := position.coordinateOfPiece(king)
		if position.coordinateIsChecked(coordinateOfKing, !forWhite) {
			continue
		}
		positionsLegal = append(positionsLegal, position)
	}

	return positionsLegal
}


func (position Position) coordinateIsChecked(coordinate Coordinate, byWhite bool) bool {

	for _, vector := range vectorsKing() {
		coordinatesOfVector := vector.coordinatesFromCoordinate(coordinate)
		for _, coordinateOfVector := range coordinatesOfVector {
			pieceAtVector := position.pieceAtCoordinate(coordinateOfVector)
			if pieceAtVector.isNotNil() && pieceAtVector.isWhite() == byWhite && pieceAtVector.isKing() {
				return true
			}
		}
	}

	for _, vector := range vectorsQueen() {
		coordinatesOfVector := vector.coordinatesFromCoordinate(coordinate)
		for _, coordinateOfVector := range coordinatesOfVector {
			pieceAtVector := position.pieceAtCoordinate(coordinateOfVector)
			if pieceAtVector.isNotNil() && pieceAtVector.isWhite() == byWhite && pieceAtVector.isQueen() {
				return true
			}
			if pieceAtVector.isNotNil() {
				break
			}
		}
	}

	for _, vector := range vectorsRook() {
		coordinatesOfVector := vector.coordinatesFromCoordinate(coordinate)
		for _, coordinateOfVector := range coordinatesOfVector {
			pieceAtVector := position.pieceAtCoordinate(coordinateOfVector)
			if pieceAtVector.isNotNil() && pieceAtVector.isWhite() == byWhite && pieceAtVector.isRook() {
				return true
			}
			if pieceAtVector.isNotNil() {
				break
			}
		}
	}

	for _, vector := range vectorsBishop() {
		coordinatesOfVector := vector.coordinatesFromCoordinate(coordinate)
		for _, coordinateOfVector := range coordinatesOfVector {
			pieceAtVector := position.pieceAtCoordinate(coordinateOfVector)
			if pieceAtVector.isNotNil() && pieceAtVector.isWhite() == byWhite && pieceAtVector.isBishop() {
				return true
			}
			if pieceAtVector.isNotNil() {
				break
			}
		}
	}

	for _, vector := range vectorsKnight() {
		coordinatesOfVector := vector.coordinatesFromCoordinate(coordinate)
		for _, coordinateOfVector := range coordinatesOfVector {
			pieceAtSquareOfVector := position.pieceAtCoordinate(coordinateOfVector)
			if pieceAtSquareOfVector.isNotNil() && pieceAtSquareOfVector.isWhite() == byWhite && pieceAtSquareOfVector.isKnight() {
				return true
			}
		}
	}

	for _, vector := range vectorsPawnCapture(!byWhite) {
		squaresOfVector := vector.coordinatesFromCoordinate(coordinate)
		for _, coordinateOfVector := range squaresOfVector {
			pieceAtVector := position.pieceAtCoordinate(coordinateOfVector)
			if pieceAtVector.isNotNil() && pieceAtVector.isWhite() == byWhite && pieceAtVector.isPawn() {
				return true
			}
		}
	}

	return false
}


func (position Position) coordinateOfPiece(piece Piece) Coordinate {
	for coordinate, p := range position.pieces {
		if p == piece {
			return coordinate
		}
	}
	return Coordinate{-1, -1}
}


func (position Position) String() string {
	str := ""
	for _, y := range ysAll(true) {
		for _, x := range xsAll(false) {
			coordinate := Coordinate{x, y}
			piece := position.pieces[coordinate]
			var print string
			if piece.code == "" {
				print = "-"
			} else {
				print = piece.code
			}
			str += print
			if x == 7 {
				str += "\n"
			}
		}
	}
	str += "________\n"
	return str
}


func (position Position) countPiecesAll() int {
	count := 0
	for _, piece := range position.pieces {
		if piece.isNotNil() {
			count++
		}
	}
	return count
}

func (position Position) countPieces(piece Piece) int {
	count := 0
	for _, p := range position.pieces {
		if p == piece {
			count++
		}
	}
	return count
}


func (position Position) evaluation() float64 {

	evaluation := 0.0

	for _, piece := range position.pieces {
		var value float64
		if piece.isQueen() {
			value += 9.0
		} else if piece.isRook() {
			value += 5.0
		} else if piece.isBishop() {
			value += 3.5
		} else if piece.isKnight() {
			value += 3.0
		} else if piece.isPawn() {
			value += 1.0
		}
		var factor float64
		if piece.isWhite() {
			factor = 1.0
		} else if piece.isBlack() {
			factor = -1.0
		}
		evaluation += (factor * value)
	}

	return evaluation
}


func (position Position) childForNotation(notation string, forWhite bool) Position {

	if notation == notationCastle(true) {
		notation = notation2Castle(forWhite, true)
	} else if notation == notationCastle(false) {
		notation = notation2Castle(forWhite, false)
	}

	notationsPieces := notationsPieces()
	notationsX := notationsX()
	notationsY := notationsY()

	notationPiece := "-"
	notationPromotion := "-"

	if strings.Index(notationsPieces, string(notation[0])) != -1 {
		notationPiece = string(notation[0])
		notation = notation[1:]
	} else {
		notationPiece = "p"
	}

	if strings.Index(notationsPieces, string(notation[len(notation)-1])) != -1 {
		notationPromotion = string(notation[len(notation)-1])
		notation = notation[:len(notation)-1]
	}

	var xs []int
	var ys []int

	for _, rune := range notation {
		str := string(rune)
		if unicode.IsLetter(rune) {
			x := strings.Index(notationsX, str)
			if x != -1 {
				xs = append(xs, x)
			}
		} else if unicode.IsNumber(rune) {
			y := strings.Index(notationsY, str)
			ys = append(ys, y)
		}
	}

	pieceMove := pieceForNotation(notationPiece, forWhite)
	piecePromote := pieceForNotation(notationPromotion, forWhite)

	coordinateFinal := Coordinate{-1, -1}
	if len(xs) > 0 && len(ys) > 0 {
		coordinateFinal = Coordinate{xs[len(xs)-1], ys[len(ys)-1]}
	}

	xInitial := -1
	yInitial := -1

	if len(xs) > 1 {
		xInitial = xs[0]
	}
	if len(ys) > 1 {
		yInitial = ys[0]
	}

	var positionsChildren = position.children(forWhite)

	for _, positionChild := range positionsChildren {

		pieceFinal := positionChild.pieces[coordinateFinal]

		if pieceFinal != pieceMove && pieceFinal != piecePromote {
			continue
		}

		if xInitial != -1 {
			isPieceMissingInX := false
			for _, y := range ysAll(false) {
				coordinate := Coordinate{xInitial, y}
				pieceBefore := position.pieces[coordinate]
				pieceAfter := positionChild.pieces[coordinate]
				if pieceBefore.isNotNil() && !pieceAfter.isNotNil() {
					isPieceMissingInX = true
					break
				}
			}
			if !isPieceMissingInX {
				continue
			}
		}

		if yInitial != -1 {
			isPieceMissingInY := false
			for _, x := range xsAll(false) {
				coordinate := Coordinate{x, yInitial}
				pieceBefore := position.pieces[coordinate]
				pieceAfter := positionChild.pieces[coordinate]
				if pieceBefore.isNotNil() && !pieceAfter.isNotNil() {
					isPieceMissingInY = true
					break
				}
			}
			if !isPieceMissingInY {
				continue
			}
		}

		if piecePromote.isNotNil() {
			pawn := piecePawn(forWhite)
			if positionChild.countPieces(pawn) == position.countPieces(pawn) {
				continue
			}
		}

		if pieceMove.isRook() {
			king := pieceKing(forWhite)
			if position.coordinateOfPiece(king) != positionChild.coordinateOfPiece(king) {
				continue
			}
		}

		return positionChild
	}

	message := fmt.Sprintf("No child found for notation %v", notation)
	panic(message)
}


func (position Position) childForCoordinates(coordinate1 Coordinate, coordinate2 Coordinate, forWhite bool) Position {

	for _, child := range position.children(forWhite) {
		
		piece1Before := position.pieces[coordinate1]
		piece1After :=  child.pieces[coordinate1]
		piece2After :=  child.pieces[coordinate2]

		if !piece1Before.isNotNil() {
			continue
		}

		if piece1After.isNotNil() {
			continue
		}

		if piece1Before.code != piece2After.code {
			continue
		}

		return child
	}
	
	message := fmt.Sprintf("No child found for coordinates %v, %v", coordinate1, coordinate2)
	panic(message)
}


func (position Position) childForPieces(pieces map[Coordinate]Piece, forWhite bool) Position {

	children := position.children(forWhite)

	for _, child := range children {
		hasSamePieces := true
		for coordinate, piece := range child.pieces {
			pieceChild := pieces[coordinate]
			if pieceChild.code != piece.code {
				hasSamePieces = false
				break
			}
		}
		if hasSamePieces {
			return child
		}
	}
	
	// for _, child := range children {
	// 	fmt.Println(child)
	// }
	
	positionFailed := Position{pieces, [4]bool{true, true, true, true}, Coordinate{}}
	fmt.Println(positionFailed)
	panic("No child found for pieces")
}


func (position Position) childrenForCoordinate(coordinate1 Coordinate, forWhite bool) []Position {

	var children []Position

	for _, child := range position.children(forWhite) {
		
		piece1Before := position.pieces[coordinate1]
		piece1After :=  child.pieces[coordinate1]

		if !piece1Before.isNotNil() {
			continue
		}

		if piece1After.isNotNil() {
			continue
		}

		children = append(children, child)
	}
	
	return children
}


func (position Position) coordinatesAndPromotionForChild(child Position) (Coordinate, Coordinate, Piece) {

	coordinate1 := Coordinate{}

	for coordinate, piece := range position.pieces {
		pieceAfter := child.pieces[coordinate]
		if !pieceAfter.isNotNil() {
			coordinate1 = coordinate
			if !piece.isRook() {
				break
			}
		}
	}

	coordinate2 := Coordinate{}

	for coordinate, piece := range child.pieces {
		pieceBefore := position.pieces[coordinate]
		if !pieceBefore.isNotNil() || pieceBefore.isWhite() != piece.isWhite() {
			coordinate2 = coordinate
			if !piece.isRook() {
				break
			}
		}
	}

	promotion := Piece{}
	pieceBefore := position.pieces[coordinate1]
	pieceAfter := child.pieces[coordinate2]
	isPromotion := pieceBefore.isPawn() && !pieceAfter.isPawn()
	if isPromotion {
		promotion = pieceAfter
	}
	
	return coordinate1, coordinate2, promotion
}


func (position Position) childBest(forWhite bool) Position {

	children := position.children(forWhite)

	if len(children) == 0 {
		panic("No children")
	}

	random := rand.Intn(len(children))
	child := children[random]
	
	return child
}
