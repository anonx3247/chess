package main

import (
	"errors"
	"fmt"
)

type MoveInstruction struct {
	move string
	init Square
	dest Square
}

func calcMove(m Move, start Square) (dest Square) {
	dest.x = start.x + m.x
	dest.y = start.y + m.y
	return
}

func (b Board) CheckMove(move string)

func (b Board) MoveToInstruction(move string) (i MoveInstruction, e error) {
	// first handle the destination
	destStr := move[len(move)-2:]
	i.dest = getSquare(destStr)

	// then check if the move is taking (i.e. is there an x in the string)
	taking := false
	// check if string has a 'x' in it
	if move[len(move)-3] == 'x' {
		taking = true
	}

	// determine the symbol of the piece we are moving
	symbol := ""
	if taking {
		// ex: Qxh8 -> Q
		symbol = move[:len(move)-3]
	} else {
		// ex: Qh8 -> Q
		symbol = move[0 : len(move)-2]
	}

	// find all the pieces with that symbol
	findPiece := func(sym string) []Square {
		pieces := make([]Square, 0)

		// i.e. if we have 'Ne5' for example
		if len(sym) == 1 {
			for square, piece := range b.squares {
				if piece.symbol == sym {
					pieces = append(pieces, square)
				}
			}
			// i.e. when we have Nge5 for example, use this information (the 'g') to find the correct piece.
		} else {
			for square, piece := range b.squares {
				if piece.symbol == sym[0:1] && square.Name()[0] == sym[1] {
					pieces = append(pieces, square)
				}
			}
		}
		return pieces
	}

	// get all the possible pieces that have the symbol 'symbol'
	possiblePieces := findPiece(symbol)
	if len(possiblePieces) == 0 {
		e = errors.New(fmt.Sprint("No matching pieces found for:", symbol))
	}

	// check each piece to see if they can reach the destination

	numberOfInitsFound := 0
	noPiecesFound := true
	for _, pieceSquare := range possiblePieces {
		//p := b.squares[pieceSquare]

		reachableSquares, err := b.GetAvailableSquaresFor(pieceSquare)
		e = err

		//reachableSquares := p.GetAvailableSquares(pieceSquare)
		for _, sq := range reachableSquares {
			if sq == i.dest {
				i.init = pieceSquare
				noPiecesFound = false
				numberOfInitsFound++
			}
		}
	}

	if noPiecesFound {
		e = errors.New(fmt.Sprint("No pieces could go to:", i.dest))
	}

	if numberOfInitsFound > 1 {
		e = errors.New(fmt.Sprint("Found", numberOfInitsFound, " init points"))
	}
	return
}

func (c ChessPiece) GetAvailableSquares(loc Square) []Square {
	moves := c.GetLegalMoves(loc)
	squares := make([]Square, 0)

	for _, move := range moves {
		squares = append(squares, calcMove(move, loc))
	}
	return squares
}

func (b Board) GetAvailableSquaresFor(sq Square) (squares []Square, e error) {
	piece := b.squares[sq]
	squares = make([]Square, 0)

	switch piece.symbol {
	case "p":
		squares = b.GetAvailableSquaresForPawn(sq)
		return
	case "N":
		squares = b.GetAvailableSquaresForKnight(sq)
		return
	case "B":
		squares = b.GetAvailableSquaresForBishop(sq)
		return
	case "R":
		squares = b.GetAvailableSquaresForRook(sq)
		return
	case "Q":
		squares = b.GetAvailableSquaresForQueen(sq)
		return
	case "K":
		squares = b.GetAvailableSquaresForKing(sq)
		return
	default:
		e = errors.New(fmt.Sprint("Unknown symbol:", piece.symbol))
		return

	}
}

func (b Board) GetAvailableSquaresForPawn(sq Square) (squares []Square) {
	squares = make([]Square, 0)
	piece := b.squares[sq]
	moves := piece.GetLegalMoves(sq)
	for _, move := range moves {
		kind, _, ud := moveType(move)
		dest := calcMove(move, sq)
		if (piece.color == White && ud == Up) || (piece.color == Black && ud == Down) {
			if kind == "diagonal" {
				// check if it's diagonals have pieces to kill
				if b.squares[dest].symbol != "_" {
					squares = append(squares, dest)
				}
			} else if kind == "vertical" {
				// check if there is a piece in front for vertical move
				if b.squares[dest].symbol == "_" {
					squares = append(squares, dest)
				}
			}

		}
	}
	return
}

func (b Board) GetAvailableSquaresForKnight(sq Square) []Square {
	p := b.squares[sq]
	// knights can go anywhere so theres nothing else to do
	return p.GetAvailableSquares(sq)
}

func (b Board) GetAvailableSquaresForBishop(sq Square) (squares []Square) {
	squares = make([]Square, 0)
	piece := b.squares[sq]
	moves := piece.GetLegalMoves(sq)
	// these will keep track of the minimum number of
	// squares along a diagonal before encountering another piece
	minRU := boardSize
	minRD := boardSize
	minLU := boardSize
	minLD := boardSize

	// find the minimum distance to a piece along each diagonal
	for _, move := range moves {
		kind, rl, ud := moveType(move)
		dest := calcMove(move, sq)
		// this may seem redundant or unneccessary but
		// since this method is used by the king and queen
		// we still need to filter out non-diagonal moves
		if kind == "diagonal" {
			if b.squares[dest].symbol != "_" {
				if rl && ud {
					newMinRU := move.x // each time we pick a positive coordinate
					if newMinRU < minRU {
						minRU = newMinRU
					}
				} else if rl && !ud {
					newMinRD := move.x
					if newMinRD < minRD {
						minRD = newMinRD
					}
				} else if !rl && ud {
					newMinLU := move.y
					if newMinLU < minLU {
						minLU = newMinLU
					}
				} else if !rl && !ud {
					newMinLD := -move.y
					if newMinLD < minLD {
						minLD = newMinLD
					}
				}
			}
		}
	}

	for _, move := range moves {
		_, rl, ud := moveType(move)
		dest := calcMove(move, sq)
		// we can take different colored pieces, we cannot take same colored pieces
		if b.squares[dest].symbol != "_" && b.squares[dest].color != piece.color {
			if rl && ud && move.x < minRU {
				squares = append(squares, dest)
			} else if rl && !ud && move.x < minRD {
				squares = append(squares, dest)
			} else if !rl && ud && move.y < minLU {
				squares = append(squares, dest)
			} else if !rl && !ud && -move.x < minLD {
				squares = append(squares, dest)
			}
		}
	}
	return
}

func (b Board) GetAvailableSquaresForRook(sq Square) (squares []Square) {
	squares = make([]Square, 0)
	piece := b.squares[sq]
	moves := piece.GetLegalMoves(sq)
	// these will keep track of the minimum number of
	// squares along an orthogonal before encountering another piece
	minR := boardSize
	minL := boardSize
	minU := boardSize
	minD := boardSize

	// find the minimum distance to a piece along each orthogonal
	for _, move := range moves {
		kind, rl, ud := moveType(move)
		dest := calcMove(move, sq)
		if b.squares[dest].symbol != "_" {
			if kind == "vertical" {
				if ud {
					newMinU := move.y
					if newMinU < minU {
						minU = newMinU
					}
				} else {
					newMinD := -move.y
					if newMinD < minD {
						minD = newMinD
					}
				}
			} else if kind == "horizontal" {
				if rl {
					newMinR := move.x
					if newMinR < minR {
						minR = newMinR
					}
				} else {
					newMinL := -move.x
					if newMinL < minL {
						minL = newMinL
					}
				}

			}
		}
	}

	for _, move := range moves {
		kind, rl, ud := moveType(move)
		dest := calcMove(move, sq)
		// we can take different colored pieces, we cannot take same colored pieces
		if b.squares[dest].symbol != "_" && b.squares[dest].color != piece.color {
			if kind == "horizontal" && rl && move.x < minR {
				squares = append(squares, dest)
			} else if kind == "horizontal" && !rl && -move.x < minL {
				squares = append(squares, dest)
			} else if kind == "vertical" && ud && move.y < minU {
				squares = append(squares, dest)
			} else if kind == "vertical" && !ud && -move.y < minD {
				squares = append(squares, dest)
			}
		}
	}
	return
}

func (b Board) GetAvailableSquaresForQueen(sq Square) (squares []Square) {
	bishopSquares := b.GetAvailableSquaresForBishop(sq)
	rookSquares := b.GetAvailableSquaresForRook(sq)
	squares = append(bishopSquares, rookSquares...)
	return
}

func (b Board) GetAvailableSquaresForKing(sq Square) (squares []Square) {
	// the same rules apply to the king and queen, but since `sq` refers
	// to the king, its set of moves will already be limited when calling `GetAvailableSquaresForQueen`
	squares = b.GetAvailableSquaresForQueen(sq)
	return
}

func moveType(m Move) (kind string, rl Direction, ud Direction) {
	if m.x == m.y && m.x > 0 {
		kind, rl, ud = "diagonal", Right, Up
	} else if m.x == -m.y && m.x > 0 {
		kind, rl, ud = "diagonal", Right, Down
	} else if m.x == m.y && m.x < 0 {
		kind, rl, ud = "diagonal", Left, Down
	} else if m.x == -m.y && m.x < 0 {
		kind, rl, ud = "diagonal", Left, Up
	} else if m.x == 0 && m.y > 0 {
		// here the 'Right' or 'Left' doesnt matter
		kind, rl, ud = "vertical", Right, Up
	} else if m.x == 0 && m.y < 0 {
		// here the 'Right' or 'Left' doesnt matter
		kind, rl, ud = "vertical", Right, Down
	} else if m.x > 0 && m.y == 0 {
		// here the 'Up' or 'Down' doesnt matter
		kind, rl, ud = "horizontal", Right, Up
	} else if m.x < 0 && m.y == 0 {
		// here the 'Up' or 'Down' doesnt matter
		kind, rl, ud = "horizontal", Left, Up
	} else {
		kind, rl, ud = "unknown", Right, Up
	}
	return
}