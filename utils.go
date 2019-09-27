
package main

import (
	"strings"
	"math"
	"time"
)

// Constructors
// ...

// position

func positionInitial() Position {
	strPosition := "RNBQKBNRPPPPPPPP--------------------------------pppppppprnbqkbnr"
	return positionFromString(strPosition)
}

func positionFromString(str string) Position {
	var pieces map[Coordinate]Piece
	pieces = make(map[Coordinate]Piece)
	for index, code := range strings.Split(str, "") {
		if code == "-" {
			continue
		}
		x := index % 8
		y := index / 8
		coordinate := Coordinate{x, y}
		pieces[coordinate] = Piece{code}
	}
	return Position{pieces, [4]bool{true, true, true, true}, Coordinate{-1, -1}}
}

// code

func codeKing() string {
	return "k"
}

func codeQueen() string {
	return "q"
}

func codeRook() string {
	return "r"
}

func codeBishop() string {
	return "b"
}

func codeKnight() string {
	return "n"
}

func codePawn() string {
	return "p"
}

// piece

func pieceKing(isWhite bool) Piece {
	code := "k"
	if isWhite {
		code = strings.ToUpper(code)
	}
	return Piece{code}
}

func pieceQueen(isWhite bool) Piece {
	code := "q"
	if isWhite {
		code = strings.ToUpper(code)
	}
	return Piece{code}
}

func pieceRook(isWhite bool) Piece {
	code := "r"
	if isWhite {
		code = strings.ToUpper(code)
	}
	return Piece{code}
}

func pieceBishop(isWhite bool) Piece {
	code := "b"
	if isWhite {
		code = strings.ToUpper(code)
	}
	return Piece{code}
}

func pieceKnight(isWhite bool) Piece {
	code := "n"
	if isWhite {
		code = strings.ToUpper(code)
	}
	return Piece{code}
}

func piecePawn(isWhite bool) Piece {
	code := "p"
	if isWhite {
		code = strings.ToUpper(code)
	}
	return Piece{code}
}

func pieceForNotation(notation string, forWhite bool) Piece {
	if notation == "k" {
		return pieceKing(forWhite)
	} else if notation == "q" {
		return pieceQueen(forWhite)
	} else if notation == "r" {
		return pieceRook(forWhite)
	} else if notation == "b" {
		return pieceBishop(forWhite)
	} else if notation == "n" {
		return pieceKnight(forWhite)
	} else if notation == "p" {
		return piecePawn(forWhite)
	}
	return Piece{"-"}
}

func piecesPromote(forWhite bool) [3]Piece {
	return [3]Piece{pieceQueen(forWhite), pieceRook(forWhite), pieceBishop(forWhite)}

}

func pieceWithName(name string) Piece {

	code := ""

	if strings.Contains(name, "king") {
		code = codeKing()
	} else if strings.Contains(name, "queen") {
		code = codeQueen()
	} else if strings.Contains(name, "rook") {
		code = codeRook()
	} else if strings.Contains(name, "bishop") {
		code = codeBishop()
	} else if strings.Contains(name, "knight") {
		code = codeKnight()
	} else if strings.Contains(name, "pawn") {
		code = codePawn()
	}

	if strings.Contains(name, "white") {
		code = strings.ToUpper(code)
	}

	return Piece{code}
}

// notations

func notationsPieces() string {
	return "kqrbn"
}

func notationsX() string {
	return "abcdefgh"
}

func notationsY() string {
	return "12345678"
}

func notation2Castle(isWhite bool, isKingside bool) string {
	if isWhite {
		if isKingside {
			return "kg1"
		}
		return "kc1"
	}
	if isKingside {
		return "kg8"
	}
	return "kc8"
}

func notationCastle(isKingside bool) string {
	if isKingside {
		return "o-o"
	}
	return "o-o-o"
}

// coordinate

func coordinateInitialOfKing(forWhite bool) Coordinate {
	if forWhite {
		return Coordinate{4, 0}
	}
	return Coordinate{4, 7}
}

func coordinateInitialOfRook(forWhite bool, forKingSide bool) Coordinate {
	if forWhite {
		if forKingSide {
			return Coordinate{7, 0}
		}
		return Coordinate{0, 0}
	}
	if forKingSide {
		return Coordinate{7, 7}
	}
	return Coordinate{0, 7}
}

// vector

func vectorsKing() []Vector {
	vectors := []Vector{}
	deltas := [3]int{1, 0, -1}
	for _, delta1 := range deltas {
		for _, delta2 := range deltas {
			if delta1 == 0 && delta2 == 0 {
				continue
			}
			vector := Vector{delta1, delta2, 1}
			vectors = append(vectors, vector)
		}
	}
	return vectors
}

func vectorsQueen() []Vector {
	vectors := []Vector{}
	deltas := [3]int{1, 0, -1}
	for _, delta1 := range deltas {
		for _, delta2 := range deltas {
			if delta1 == 0 && delta2 == 0 {
				continue
			}
			vector := Vector{delta1, delta2, 7}
			vectors = append(vectors, vector)
		}
	}
	return vectors
}

func vectorsRook() []Vector {
	vectors := []Vector{}
	deltas := [3]int{1, 0, -1}
	for _, delta1 := range deltas {
		for _, delta2 := range deltas {
			if math.Abs(float64(delta1)) == math.Abs(float64(delta2)) {
				continue
			}
			vector := Vector{delta1, delta2, 7}
			vectors = append(vectors, vector)
		}
	}
	return vectors
}

func vectorsBishop() []Vector {
	vectors := []Vector{}
	deltas := [3]int{1, 0, -1}
	for _, delta1 := range deltas {
		for _, delta2 := range deltas {
			if math.Abs(float64(delta1)) != math.Abs(float64(delta2)) || delta1 == 0 || delta2 == 0 {
				continue
			}
			vector := Vector{delta1, delta2, 7}
			vectors = append(vectors, vector)
		}
	}
	return vectors
}

func vectorsKnight() []Vector {
	vectors := []Vector{}
	deltas := [4]int{2, 1, -1, -2}
	for _, delta1 := range deltas {
		for _, delta2 := range deltas {
			if math.Abs(float64(delta1)) == math.Abs(float64(delta2)) {
				continue
			}
			vector := Vector{delta1, delta2, 1}
			vectors = append(vectors, vector)
		}
	}
	return vectors
}

func vectorCastle(king bool) Vector {
	if king {
		return Vector{1, 0, 2}
	}
	return Vector{-1, 0, 2}
}

func vectorsPawnCapture(forWhite bool) []Vector {
	dy := dyPawn(forWhite)
	return []Vector{Vector{-1, dy, 1}, Vector{1, dy, 1}}
}

func vectorPawnAdvance(forWhite bool) Vector {
	dy := dyPawn(forWhite)
	return Vector{0, dy, 2}
}

// x, y

func dyPawn(forWhite bool) int {
	if forWhite {
		return 1
	}
	return -1
}

func yInitialPawn(forWhite bool) int {
	if forWhite {
		return 1
	}
	return 6
}

func yFinalPawn(forWhite bool) int {
	if forWhite {
		return 7
	}
	return 0
}

func xsAll(reversed bool) [8]int {
	if reversed {
		return [8]int{7, 6, 5, 4, 3, 2, 1, 0}
	}
	return [8]int{0, 1, 2, 3, 4, 5, 6, 7}
}

func ysAll(reversed bool) [8]int {
	if reversed {
		return [8]int{7, 6, 5, 4, 3, 2, 1, 0}
	}
	return [8]int{0, 1, 2, 3, 4, 5, 6, 7}
}

// time

func sleep(seconds float32) {
	durationDelay := time.Duration(seconds * 1000000000)
	time.Sleep(durationDelay)
}