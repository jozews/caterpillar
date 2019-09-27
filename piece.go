package main

import (
	"strings"
)

// Piece ...
type Piece struct {
	code string
}

func (piece Piece) isWhite() bool {
	return piece.code == strings.ToUpper(piece.code)
}

func (piece Piece) isBlack() bool {
	return piece.code == strings.ToLower(piece.code)
}

func (piece Piece) isKing() bool {
	return strings.EqualFold(piece.code, codeKing())
}

func (piece Piece) isQueen() bool {
	return strings.EqualFold(piece.code, codeQueen())
}

func (piece Piece) isRook() bool {
	return strings.EqualFold(piece.code, codeRook())
}

func (piece Piece) isBishop() bool {
	return strings.EqualFold(piece.code, codeBishop())
}

func (piece Piece) isKnight() bool {
	return strings.EqualFold(piece.code, codeKnight())
}

func (piece Piece) isPawn() bool {
	return strings.EqualFold(piece.code, codePawn())
}

func (piece Piece) isNotNil() bool {
	return piece.isKing() || piece.isQueen() || piece.isRook() || piece.isBishop() || piece.isKnight() || piece.isPawn()
}

func (piece Piece) name() string {
	if piece.isKing() {
		return "king"
	} else if piece.isQueen() {
		return "queen"
	} else if piece.isRook() {
		return "rook"
	} else if piece.isBishop() {
		return "bishop"
	} else if piece.isKnight() {
		return "knight"
	} else if piece.isPawn() {
		return "pawn"
	}
	return ""
}
