

package main


import (
	"github.com/fedesog/webdriver"
	"strconv"
	"strings"
	"fmt"
)


// Lichess ...
type Lichess struct {
	driver  *webdriver.ChromeDriver
	session *webdriver.Session
	sizeBoard webdriver.Size
}


func (lichess *Lichess) setup() {

	lichess.driver = webdriver.NewChromeDriver("/home/jozews/chromedriver")
	lichess.driver.Start()
	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	lichess.session, _ = lichess.driver.NewSession(desired, required)
	lichess.session.Url("https://lichess.org/0u3G8ViY/black")
}


func (lichess *Lichess) end() {
	lichess.session.Delete()
	lichess.driver.Stop()
}


func (lichess *Lichess) offsetForCoordinate(coordinate Coordinate) (int, int) {

	width := lichess.sizeBoard.Width
	widthSquare := width / 8

	if lichess.isOrientationWhite() {
		return width - (widthSquare*coordinate.x+widthSquare/2), 0 + (widthSquare*coordinate.y+widthSquare/2)
	}

	return 0 + (widthSquare*coordinate.x+widthSquare/2), width - (widthSquare*coordinate.y+widthSquare/2)
}


func (lichess *Lichess) coordinateForOffset(xPosition int, yPosition int) (int, int) {
	
	widthSquare := lichess.sizeBoard.Width / 8

	if lichess.isOrientationWhite() {
		return 7 - (xPosition/widthSquare), 0 + (yPosition/widthSquare)
	}
	return 0 + (xPosition/widthSquare), 7 - (yPosition/widthSquare)
}


func (lichess *Lichess) isOrientationWhite() bool {
	_, error := lichess.session.FindElement(webdriver.XPath, "//div[@class=\"cg-wrap cgv1 manipulable orientation-white\"]")
	return error != nil
}


func (lichess *Lichess) pieceForSquare(square webdriver.WebElement) Piece {

	class, _ := square.GetAttribute("class")
	style, _ := square.GetAttribute("style")

	if strings.Contains(class, "dragging") {
		return Piece{}
	}
	if strings.Contains(class, "ghost") && strings.Contains(style, "hidden") {
		return Piece{}
	}

	piece := pieceWithName(class)

	return piece
}


func (lichess *Lichess) coordinateForSquare(square webdriver.WebElement) Coordinate {

	style, _ := square.GetAttribute("style")

	if !strings.Contains(style, "transform:") {
		return Coordinate{-1, -1}
	}

	splits := strings.Split(style, "px")

	splitX := splits[0]
	splitY := splits[1]

	xPosition := 0
	yPosition := 0

	for index := range splitX {
		indexReversed := len(splitX) - 1 - index
		substring := splitX[indexReversed:len(splitX)]
		value, error := strconv.Atoi(substring)
		if error == nil {
			xPosition = value
		} else {
			break
		}
	}

	for index := range splitY {
		indexReversed := len(splitY) - 1 - index
		substring := splitY[indexReversed:len(splitY)]
		value, error := strconv.Atoi(substring)
		if error == nil {
			yPosition = value
		} else {
			break
		}
	}

	x, y := lichess.coordinateForOffset(xPosition, yPosition)
	coordinate := Coordinate{x, y}

	return coordinate
}


func (lichess *Lichess) pieces() map[Coordinate]Piece {

	pieces := make(map[Coordinate]Piece)

	board, _ := lichess.session.FindElement(webdriver.XPath, "//cg-board")
	lichess.sizeBoard, _ = board.Size()

	elements, _ := lichess.session.FindElements(webdriver.XPath, "//piece")

	for _, element := range elements {

		piece := lichess.pieceForSquare(element)
		coordinate := lichess.coordinateForSquare(element)

		if !piece.isNotNil() || !coordinate.isInBounds() {
			continue
		}

		pieces[coordinate] = piece
	}	

	return pieces
}


func (lichess *Lichess) observePieces(pieces map[Coordinate]Piece) map[Coordinate]Piece {

	for {
	
		sleep(0.5)
		
		board, _ := lichess.session.FindElement(webdriver.XPath, "//cg-board")
		lichess.sizeBoard, _ = board.Size()

		elements, _ := lichess.session.FindElements(webdriver.XPath, "//piece")
		
		for _, element := range elements {
			
			piece := lichess.pieceForSquare(element)
			coordinate := lichess.coordinateForSquare(element)
			
			if !piece.isNotNil() {
				continue
			}
			
			pieceBefore := pieces[coordinate]
			if pieceBefore.code != piece.code {
				pieces := lichess.pieces()
				return pieces
			}
		}
	}
}


func (lichess *Lichess) observeSelection(selection Coordinate) Coordinate {

	for {
		sleep(1)
		selectionNew := lichess.coordinateSelection()
		if selectionNew != selection {
			return selectionNew
		}
	}
}


func (lichess *Lichess) makeMove(coordinate1 Coordinate, coordinate2 Coordinate, promotion Piece, smooth bool) {

	board, _ := lichess.session.FindElement(webdriver.XPath, "//cg-board")
	lichess.sizeBoard, _ = board.Size()

	x1, y1 := lichess.offsetForCoordinate(coordinate1)
	x2, y2 := lichess.offsetForCoordinate(coordinate2)

	lichess.session.MoveTo(board, x1, y1)
	lichess.session.ButtonDown(webdriver.LeftButton)
	
	if smooth {
		sleep(0.5)
	}

	lichess.session.MoveTo(board, x2, y2)
	lichess.session.ButtonUp(webdriver.LeftButton)

	if smooth {
		sleep(0.5)
	}
	
	if promotion.isNotNil() {
		xpath := fmt.Sprintf("//piece[starts-with(@class, \"%v\")]/..", promotion.name())
		square, _ := lichess.session.FindElement(webdriver.XPath, xpath)
		square.Click()
	}
}


func (lichess Lichess) selectCoordinate(coordinate Coordinate) {	
	
	board, _ := lichess.session.FindElement(webdriver.XPath, "//cg-board")
	lichess.sizeBoard, _ = board.Size()

	x, y := lichess.offsetForCoordinate(coordinate)

	lichess.session.MoveTo(board, x, y)
	lichess.session.ButtonDown(webdriver.LeftButton)
}


func (lichess Lichess) coordinateSelection() Coordinate {	

	square, error := lichess.session.FindElement(webdriver.XPath, "//square[contains(@class, \"selected\")]")

	if error != nil {
		return Coordinate{-1, -1}
	}

	coordinate := lichess.coordinateForSquare(square)
	return coordinate
}


func (lichess Lichess) coordinatesForDestinations() []Coordinate {

	var coordinates []Coordinate

	squares, _ := lichess.session.FindElements(webdriver.XPath, "//square[@class=\"move-dest\"]")

	for _, square := range squares {
		coordinate := lichess.coordinateForSquare(square)
		coordinates = append(coordinates, coordinate)
	}

	return coordinates
}