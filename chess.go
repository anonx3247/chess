package main

import (
	"encoding/csv"
	"strconv"
)

const Black Color = false
const White Color = true

const Alive State = true
const Dead State = false

const boardSize int = 7

const Right Direction = true
const Left Direction = false
const Up Direction = true
const Down Direction = false

type Direction bool
type Color bool
type State bool

type Square struct {
	x int
	y int
}

type Move Square

type ChessPiece struct {
	color      Color
	state      State
	legalMoves []Move
	value      int
	symbol     string
}

type Board struct {
	squares    map[Square]*ChessPiece
	nowPlaying Color
	moves      []string
}

func (s Square) Name() string {
	name := ""
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	name += letters[s.x]
	name += strconv.Itoa(s.y + 1)
	return name
}

type Movable interface {
	GetLegalMoves(loc Square) []Move
	GetAvailableSquares() []Square
}

func IsOutOfBounds(m Move, loc Square) bool {
	dest := calcMove(m, loc)
	newX, newY := dest.x, dest.y

	if newX < 0 || newX > boardSize || newY < 0 || newY > boardSize {
		return true
	} else {
		return false
	}
}

func (p ChessPiece) LegalMoves(loc Square) []Move {
	legal := make([]Move, 0)
	for _, move := range p.legalMoves {
		if !IsOutOfBounds(move, loc) {
			legal = append(legal, move)
		}
	}
	return legal
}

func newPawn() (pawn ChessPiece) {
	pawn.symbol = "p"
	pawn.value = 1
	pawn.legalMoves = []Move{
		/*
			| |x| |
			| |p| |
		*/
		Move{0, 1},
		/*
			|x| | |
			| |p| |
		*/
		Move{-1, 1},
		/*
			| | |x|
			| |p| |
		*/
		Move{1, 1},
		/*
			| |p| |
			| |x| |
		*/
		Move{0, -1},
		/*
			| |p| |
			|x| | |
		*/
		Move{-1, -1},
		/*
			| |p| |
			| | |x|
		*/
		Move{1, -1},
	}
	return
}

func diagMov(length int, RightLeft Direction, UpDown Direction) Move {
	if RightLeft {
		if UpDown {
			/*
				| | |x|
				| |x| |
				|B| | |
			*/
			return Move{length, length}
		} else {
			/*
				|B| | |
				| |x| |
				| | |x|
			*/
			return Move{length, -length}
		}
	} else {
		if UpDown {
			/*
				|x| | |
				| |x| |
				| | |B|
			*/
			return Move{-length, length}
		} else {
			/*
				| | |B|
				| |x| |
				|x| | |
			*/
			return Move{-length, -length}
		}
	}
}

func diagMovRange(rl Direction, ud Direction) []Move {
	L := make([]Move, 7)
	for i := 1; i < 8; i++ {
		L[i-1] = diagMov(i, rl, ud)
	}
	return L
}

func orthoMov(length int, RightLeft Direction, UpDown Direction) Move {
	if RightLeft {
		if UpDown {
			/*
				| | | |
				|R|x|x|
				| | | |
			*/
			return Move{length, 0}
		} else {
			/*
				| | | |
				|x|x|R|
				| | | |
			*/
			return Move{-length, 0}
		}
	} else {
		if UpDown {
			/*
				| |R| |
				| |x| |
				| |x| |
			*/
			return Move{0, -length}
		} else {
			/*
				| |x| |
				| |x| |
				| |R| |
			*/
			return Move{0, length}
		}
	}
}

func orthoMovRange(rl Direction, ud Direction) []Move {
	L := make([]Move, 7)
	for i := 1; i < 8; i++ {
		L[i-1] = orthoMov(i, rl, ud)
	}
	return L
}

func newBishop() (bishop ChessPiece) {
	bishop.value = 3
	bishop.symbol = "B"
	bishop.legalMoves = append(diagMovRange(Right, Up), diagMovRange(Left, Up)...)
	bishop.legalMoves = append(bishop.legalMoves, diagMovRange(Left, Down)...)
	bishop.legalMoves = append(bishop.legalMoves, diagMovRange(Right, Down)...)
	return
}

func newKnight() (knight ChessPiece) {
	knight.value = 3
	knight.symbol = "N"
	knight.legalMoves = []Move{
		/*
			|x| | |
			| | | |
			| |N| |
		*/
		Move{-1, 2},
		/*
			| | |x|
			| | | |
			| |N| |
		*/
		Move{1, 2},
		/*
			| |N| |
			| | | |
			| | |x|
		*/
		Move{1, -2},
		/*
			| |N| |
			| | | |
			|x| | |
		*/
		Move{-1, -2},
		/*
			|x| | |
			| | |N|
			| | | |
		*/
		Move{-2, 1},
		/*
			| | | |
			| | |N|
			|x| | |
		*/
		Move{-2, -1},
		/*
			| | |x|
			|N| | |
			| | | |
		*/
		Move{2, 1},
		/*
			| | | |
			|N| | |
			| | |x|
		*/
		Move{2, -1},
	}
	return
}

func newRook() (rook ChessPiece) {
	rook.value = 5
	rook.symbol = "R"
	rook.legalMoves = append(orthoMovRange(Right, Up), orthoMovRange(Left, Up)...)
	rook.legalMoves = append(rook.legalMoves, orthoMovRange(Left, Down)...)
	rook.legalMoves = append(rook.legalMoves, orthoMovRange(Right, Down)...)
	return
}

func newQueen() (queen ChessPiece) {
	queen.value = 9
	queen.symbol = "Q"
	// queens can move like bishop or rook
	queen.legalMoves = append(newBishop().legalMoves, newRook().legalMoves...)
	return
}

func newKing() (king ChessPiece) {
	king.value = 100 // completely arbitrary
	king.symbol = "N"
	king.legalMoves = []Move{
		diagMov(1, Right, Up),
		diagMov(1, Right, Down),
		diagMov(1, Left, Up),
		diagMov(1, Left, Down),
		orthoMov(1, Right, Up),
		orthoMov(1, Right, Down),
		orthoMov(1, Left, Up),
		orthoMov(1, Left, Down),
		// Castling
		/*
			| | | | | |... -> | | | | | |...
			|R| | | |K|... -> | | |N|K| |...
		*/
		Move{-2, 0},
		/*
			...| | | | | | -> ...| | | | | |
			...| |K| | |R| -> ...| | |R|K| |
		*/
		Move{2, 0},
	}
	return
}

func newEmpty() (empty ChessPiece) {
	empty.value = 0
	empty.symbol = "_"
	empty.legalMoves = []Move{}
	return
}

func newGame(r csv.Reader) (b Board) {

	pawn := newPawn()
	bishop := newBishop()
	knight := newKnight()
	rook := newRook()
	queen := newQueen()
	king := newKing()
	empty := newEmpty()

	// helper function for setting squares
	set := func(sq Square, col Color, c *ChessPiece) {
		m := new(ChessPiece)
		*m = *c
		m.state = Alive
		m.color = col
		b.squares[sq] = m
	}

	squares, err := r.ReadAll()
	check(err)

	piece := map[byte]*ChessPiece{
		'p': &pawn,
		'B': &bishop,
		'N': &knight,
		'K': &king,
		'Q': &queen,
		'R': &rook,
		'_': &empty,
	}

	colorer := map[byte]Color{
		'w': White,
		'b': Black,
	}

	for i := 0; i <= boardSize; i++ {
		for j := 0; j <= boardSize; j++ {
			// we read from top to bottom, but we insert from bootom to top
			sq := Square{boardSize - i, boardSize - j}
			val := squares[i][j]
			if j < (boardSize+1)/2 {
				set(sq, colorer[val[0]], piece[val[1]])
			} else {
				set(sq, colorer[val[0]], piece[val[1]])
			}
		}
	}

	return
}

func getSquare(name string) Square {
	if len(name) >= 3 {
		panic("name too long!")
	}

	alphabet := map[byte]int{
		'a': 0,
		'b': 1,
		'c': 2,
		'd': 3,
		'e': 4,
		'f': 5,
		'g': 6,
		'h': 7,
	}

	x := alphabet[name[0]]
	y, e := strconv.Atoi(name[1:])
	check(e)
	y -= 1
	return Square{x, y}

}

func (m Move) write() string {
	sq := Square{m.x, m.y}
	return sq.Name()
}
