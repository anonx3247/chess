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

func (p ChessPiece) GetLegalMoves(loc Square) []Move {
	legal := make([]Move, 0)
	for _, move := range p.legalMoves {
		if !IsOutOfBounds(move, loc) {
			legal = append(legal, move)
		}
	}
	return legal
}

func newGame(r csv.Reader) (b Board) {

	// Pawn
	pawn := ChessPiece{}
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

	// Bishop

	diag := func(n int, RightLeft Direction, UpDown Direction) Move {
		if RightLeft {
			if UpDown {
				/*
					| | |x|
					| |x| |
					|B| | |
				*/
				return Move{n, n}
			} else {
				/*
					|B| | |
					| |x| |
					| | |x|
				*/
				return Move{n, -n}
			}
		} else {
			if UpDown {
				/*
					|x| | |
					| |x| |
					| | |B|
				*/
				return Move{-n, n}
			} else {
				/*
					| | |B|
					| |x| |
					|x| | |
				*/
				return Move{-n, -n}
			}
		}
	}

	rangeDiag := func(rl Direction, ud Direction) []Move {
		L := make([]Move, 7)
		for i := 1; i < 8; i++ {
			L[i-1] = diag(i, rl, ud)
		}
		return L
	}

	bishop := ChessPiece{}
	bishop.value = 3
	bishop.symbol = "B"
	bishop.legalMoves = make([]Move, 0)
	bishop.legalMoves = append(rangeDiag(Right, Up), rangeDiag(Left, Up)...)
	bishop.legalMoves = append(bishop.legalMoves, rangeDiag(Left, Down)...)
	bishop.legalMoves = append(bishop.legalMoves, rangeDiag(Right, Down)...)

	// Knight

	knight := ChessPiece{}
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

	// Rook

	rook := ChessPiece{}
	rook.value = 5
	rook.symbol = "R"

	horiz := func(n int, RightLeft Direction, UpDown Direction) Move {
		if RightLeft {
			if UpDown {
				/*
					| | | |
					|R|x|x|
					| | | |
				*/
				return Move{n, 0}
			} else {
				/*
					| | | |
					|x|x|R|
					| | | |
				*/
				return Move{-n, 0}
			}
		} else {
			if UpDown {
				/*
					| |R| |
					| |x| |
					| |x| |
				*/
				return Move{0, -n}
			} else {
				/*
					| |x| |
					| |x| |
					| |R| |
				*/
				return Move{0, n}
			}
		}
	}

	rangeHoriz := func(rl Direction, ud Direction) []Move {
		L := make([]Move, 7)
		for i := 1; i < 8; i++ {
			L[i-1] = horiz(i, rl, ud)
		}
		return L
	}

	rook.legalMoves = make([]Move, 0)
	rook.legalMoves = append(rangeHoriz(Right, Up), rangeHoriz(Left, Up)...)
	rook.legalMoves = append(rook.legalMoves, rangeHoriz(Left, Down)...)
	rook.legalMoves = append(rook.legalMoves, rangeHoriz(Right, Down)...)

	// Queen

	queen := ChessPiece{}
	queen.value = 9
	queen.symbol = "Q"
	// queens can move like bishop or rook
	queen.legalMoves = append(bishop.legalMoves, rook.legalMoves...)

	// King

	king := ChessPiece{}
	king.value = 100 // completely arbitrary
	king.symbol = "N"
	king.legalMoves = []Move{
		diag(1, Right, Up),
		diag(1, Right, Down),
		diag(1, Left, Up),
		diag(1, Left, Down),
		horiz(1, Right, Up),
		horiz(1, Right, Down),
		horiz(1, Left, Up),
		horiz(1, Left, Down),
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

	// Empty Square
	empty := ChessPiece{}
	empty.value = 0
	empty.symbol = "_"
	empty.legalMoves = []Move{}

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
