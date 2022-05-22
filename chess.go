package main

import (
	"strconv"
)

const Black Color = false
const White Color = true

const Alive State = true
const Dead State = false

const boardSize int = 7

type Color bool
type State bool

type Square struct {
	x int
	y int
}

type Move struct {
	destination Square
}

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
}

func IsOutOfBounds(m Move, loc Square) bool {
	newX, newY := loc.x+m.destination.x, loc.y+m.destination.y

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

func newGame() (b Board) {

	// Pawn
	pawn := ChessPiece{}
	pawn.symbol = "p"
	pawn.value = 1
	pawn.legalMoves = []Move{
		/*
			| |x| |
			| |p| |
		*/
		Move{Square{0, 1}},
		/*
			|x| | |
			| |p| |
		*/
		Move{Square{-1, 1}},
		/*
			| | |x|
			| |p| |
		*/
		Move{Square{1, 1}},
	}

	// Bishop

	diag := func(n int, RightLeft bool, UpDown bool) Move {
		if RightLeft {
			if UpDown {
				/*
					| | |x|
					| |x| |
					|B| | |
				*/
				return Move{Square{n, n}}
			} else {
				/*
					|B| | |
					| |x| |
					| | |x|
				*/
				return Move{Square{n, -n}}
			}
		} else {
			if UpDown {
				/*
					|x| | |
					| |x| |
					| | |B|
				*/
				return Move{Square{-n, n}}
			} else {
				/*
					| | |B|
					| |x| |
					|x| | |
				*/
				return Move{Square{-n, -n}}
			}
		}
	}

	right := true
	left := false
	up := true
	down := false

	bishop := ChessPiece{}
	bishop.value = 3
	bishop.symbol = "B"
	bishop.legalMoves = []Move{
		diag(1, right, up),
		diag(2, right, up),
		diag(3, right, up),
		diag(4, right, up),
		diag(5, right, up),
		diag(6, right, up),
		diag(7, right, up),
		diag(1, left, up),
		diag(2, left, up),
		diag(3, left, up),
		diag(4, left, up),
		diag(5, left, up),
		diag(6, left, up),
		diag(7, left, up),
		diag(1, left, down),
		diag(2, left, down),
		diag(3, left, down),
		diag(4, left, down),
		diag(5, left, down),
		diag(6, left, down),
		diag(7, left, down),
		diag(1, right, down),
		diag(2, right, down),
		diag(3, right, down),
		diag(4, right, down),
		diag(5, right, down),
		diag(6, right, down),
		diag(7, right, down),
	}

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
		Move{Square{-1, 2}},
		/*
			| | |x|
			| | | |
			| |N| |
		*/
		Move{Square{1, 2}},
		/*
			| |N| |
			| | | |
			| | |x|
		*/
		Move{Square{1, -2}},
		/*
			| |N| |
			| | | |
			|x| | |
		*/
		Move{Square{-1, -2}},
		/*
			|x| | |
			| | |N|
			| | | |
		*/
		Move{Square{-2, 1}},
		/*
			| | | |
			| | |N|
			|x| | |
		*/
		Move{Square{-2, -1}},
		/*
			| | |x|
			|N| | |
			| | | |
		*/
		Move{Square{2, 1}},
		/*
			| | | |
			|N| | |
			| | |x|
		*/
		Move{Square{2, -1}},
	}

	// Rook

	rook := ChessPiece{}
	rook.value = 5
	rook.symbol = "R"

	horiz := func(n int, RightLeft bool, UpDown bool) Move {
		if RightLeft {
			if UpDown {
				/*
					| | | |
					|R|x|x|
					| | | |
				*/
				return Move{Square{n, 0}}
			} else {
				/*
					| | | |
					|x|x|R|
					| | | |
				*/
				return Move{Square{-n, 0}}
			}
		} else {
			if UpDown {
				/*
					| |R| |
					| |x| |
					| |x| |
				*/
				return Move{Square{0, -n}}
			} else {
				/*
					| |x| |
					| |x| |
					| |R| |
				*/
				return Move{Square{0, n}}
			}
		}
	}

	rook.legalMoves = []Move{
		horiz(1, right, up),
		horiz(2, right, up),
		horiz(3, right, up),
		horiz(4, right, up),
		horiz(5, right, up),
		horiz(6, right, up),
		horiz(7, right, up),
		horiz(1, left, up),
		horiz(2, left, up),
		horiz(3, left, up),
		horiz(4, left, up),
		horiz(5, left, up),
		horiz(6, left, up),
		horiz(7, left, up),
		horiz(1, left, down),
		horiz(2, left, down),
		horiz(3, left, down),
		horiz(4, left, down),
		horiz(5, left, down),
		horiz(6, left, down),
		horiz(7, left, down),
		horiz(1, right, down),
		horiz(2, right, down),
		horiz(3, right, down),
		horiz(4, right, down),
		horiz(5, right, down),
		horiz(6, right, down),
		horiz(7, right, down),
	}

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
		diag(1, true, true),
		diag(1, true, false),
		diag(1, false, true),
		diag(1, false, false),
		horiz(1, true, true),
		horiz(1, true, false),
		horiz(1, false, true),
		horiz(1, false, false),
		// Castling
		/*
			| | | | | |... -> | | | | | |...
			|R| | | |K|... -> | | |N|K| |...
		*/
		Move{Square{-2, 0}},
		/*
			...| | | | | | -> ...| | | | | |
			...| |K| | |R| -> ...| | |R|K| |
		*/
		Move{Square{2, 0}},
	}

	//

	// Initial Board layout
	/*
		|R|N|B|Q|K|B|N|R|
		|p|p|p|p|p|p|p|p|
		| | | | | | | | |
		| | | | | | | | |
		| | | | | | | | |
		| | | | | | | | |
		|p|p|p|p|p|p|p|p|
		|R|N|B|Q|K|B|N|R|
	*/

	set := func(name string, col Color, c *ChessPiece) {
		sq := getSquare(name)
		m := new(ChessPiece)
		*m = *c
		m.state = Alive
		m.color = col
		b.squares[sq] = m
	}

	set("a1", White, &rook)
	set("a2", White, &knight)
	set("a3", White, &bishop)
	set("a4", White, &queen)
	set("a5", White, &king)
	set("a6", White, &bishop)
	set("a7", White, &knight)
	set("a8", White, &rook)
	set("b1", White, &pawn)
	set("b2", White, &pawn)
	set("b3", White, &pawn)
	set("b4", White, &pawn)
	set("b5", White, &pawn)
	set("b6", White, &pawn)
	set("b7", White, &pawn)
	set("b8", White, &pawn)
	set("h1", Black, &rook)
	set("h2", Black, &knight)
	set("h3", Black, &bishop)
	set("h4", Black, &queen)
	set("h5", Black, &king)
	set("h6", Black, &bishop)
	set("h7", Black, &knight)
	set("h8", Black, &rook)
	set("g1", Black, &pawn)
	set("g2", Black, &pawn)
	set("g3", Black, &pawn)
	set("g4", Black, &pawn)
	set("g5", Black, &pawn)
	set("g6", Black, &pawn)
	set("g7", Black, &pawn)
	set("g8", Black, &pawn)

	empty := ChessPiece{}
	empty.value = 0
	empty.symbol = "_"
	empty.legalMoves = []Move{}

	// arbitrarily we choose white, this has no importance
	set("c1", White, &empty)
	set("c2", White, &empty)
	set("c3", White, &empty)
	set("c4", White, &empty)
	set("c5", White, &empty)
	set("c6", White, &empty)
	set("c7", White, &empty)
	set("c8", White, &empty)
	set("d1", White, &empty)
	set("d2", White, &empty)
	set("d3", White, &empty)
	set("d4", White, &empty)
	set("d5", White, &empty)
	set("d6", White, &empty)
	set("d7", White, &empty)
	set("d8", White, &empty)
	set("e1", White, &empty)
	set("e2", White, &empty)
	set("e3", White, &empty)
	set("e4", White, &empty)
	set("e5", White, &empty)
	set("e6", White, &empty)
	set("e7", White, &empty)
	set("e8", White, &empty)
	set("f1", White, &empty)
	set("f2", White, &empty)
	set("f3", White, &empty)
	set("f4", White, &empty)
	set("f5", White, &empty)
	set("f6", White, &empty)
	set("f7", White, &empty)
	set("f8", White, &empty)
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
