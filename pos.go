package main

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
)

func Distance(from, to gruid.Point) int {
	return paths.DistanceChebyshev(from, to)
}

func DistanceX(from, to gruid.Point) int {
	deltaX := Abs(to.X - from.X)
	return deltaX
}

func DistanceY(from, to gruid.Point) int {
	deltaY := Abs(to.Y - from.Y)
	return deltaY
}

type direction int

const (
	NoDir direction = iota
	E
	ENE
	NE
	NNE
	N
	NNW
	NW
	WNW
	W
	WSW
	SW
	SSW
	S
	SSE
	SE
	ESE
)

func KeyToDir(k keyAction) (dir direction) {
	switch k {
	case KeyW, KeyRunW:
		dir = W
	case KeyE, KeyRunE:
		dir = E
	case KeyS, KeyRunS:
		dir = S
	case KeyN, KeyRunN:
		dir = N
	case KeyNW, KeyRunNW:
		dir = NW
	case KeySW, KeyRunSW:
		dir = SW
	case KeyNE, KeyRunNE:
		dir = NE
	case KeySE, KeyRunSE:
		dir = SE
	}
	return dir
}

func To(p gruid.Point, dir direction) gruid.Point {
	to := p
	switch dir {
	case E, ENE, ESE:
		to = p.Add(gruid.Point{1, 0})
	case NE:
		to = p.Add(gruid.Point{1, -1})
	case NNE, N, NNW:
		to = p.Add(gruid.Point{0, -1})
	case NW:
		to = p.Add(gruid.Point{-1, -1})
	case WNW, W, WSW:
		to = p.Add(gruid.Point{-1, 0})
	case SW:
		to = p.Add(gruid.Point{-1, 1})
	case SSW, S, SSE:
		to = p.Add(gruid.Point{0, 1})
	case SE:
		to = p.Add(gruid.Point{1, 1})
	}
	return to
}

func Dir(p, from gruid.Point) direction {
	deltaX := Abs(p.X - from.X)
	deltaY := Abs(p.Y - from.Y)
	switch {
	case p.X > from.X && p.Y == from.Y:
		return E
	case p.X > from.X && p.Y < from.Y:
		switch {
		case deltaX > deltaY:
			return ENE
		case deltaX == deltaY:
			return NE
		default:
			return NNE
		}
	case p.X == from.X && p.Y < from.Y:
		return N
	case p.X < from.X && p.Y < from.Y:
		switch {
		case deltaY > deltaX:
			return NNW
		case deltaX == deltaY:
			return NW
		default:
			return WNW
		}
	case p.X < from.X && p.Y == from.Y:
		return W
	case p.X < from.X && p.Y > from.Y:
		switch {
		case deltaX > deltaY:
			return WSW
		case deltaX == deltaY:
			return SW
		default:
			return SSW
		}
	case p.X == from.X && p.Y > from.Y:
		return S
	case p.X > from.X && p.Y > from.Y:
		switch {
		case deltaY > deltaX:
			return SSE
		case deltaX == deltaY:
			return SE
		default:
			return ESE
		}
	default:
		panic(fmt.Sprintf("internal error: invalid position:%+v-%+v", p, from))
	}
}

func RandomNeighbor(p gruid.Point, diag bool) gruid.Point {
	if diag {
		return RandomNeighborDiagonals(p)
	}
	return RandomNeighborCardinal(p)
}

func RandomNeighborDiagonals(p gruid.Point) gruid.Point {
	neighbors := [8]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1}), p.Add(gruid.Point{1, -1}), p.Add(gruid.Point{-1, -1}), p.Add(gruid.Point{1, 1}), p.Add(gruid.Point{-1, 1})}
	var r int
	switch RandInt(8) {
	case 0:
		r = RandInt(len(neighbors[0:4]))
	case 1:
		r = RandInt(len(neighbors[0:2]))
	default:
		r = RandInt(len(neighbors[4:]))
	}
	return neighbors[r]
}

func RandomNeighborCardinal(p gruid.Point) gruid.Point {
	neighbors := [8]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1}), p.Add(gruid.Point{1, -1}), p.Add(gruid.Point{-1, -1}), p.Add(gruid.Point{1, 1}), p.Add(gruid.Point{-1, 1})}
	var r int
	switch RandInt(6) {
	case 0:
		r = RandInt(len(neighbors[0:4]))
	case 1:
		r = RandInt(len(neighbors))
	default:
		r = RandInt(len(neighbors[0:2]))
	}
	return neighbors[r]
}

func idx2Point(i int) gruid.Point {
	return gruid.Point{i % DungeonWidth, i / DungeonWidth}
}

func idx(p gruid.Point) int {
	return p.Y*DungeonWidth + p.X
}

func Laterals(p gruid.Point, dir direction) []gruid.Point {
	switch dir {
	case E, ENE, ESE:
		return []gruid.Point{p.Add(gruid.Point{1, -1}), p.Add(gruid.Point{1, 1})}
	case NE:
		return []gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{0, -1})}
	case N, NNE, NNW:
		return []gruid.Point{p.Add(gruid.Point{-1, -1}), p.Add(gruid.Point{1, -1})}
	case NW:
		return []gruid.Point{p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1})}
	case W, WNW, WSW:
		return []gruid.Point{p.Add(gruid.Point{-1, 1}), p.Add(gruid.Point{-1, -1})}
	case SW:
		return []gruid.Point{p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, 1})}
	case S, SSW, SSE:
		return []gruid.Point{p.Add(gruid.Point{-1, 1}), p.Add(gruid.Point{1, 1})}
	case SE:
		return []gruid.Point{p.Add(gruid.Point{0, 1}), p.Add(gruid.Point{1, 0})}
	default:
		// should not happen
		return []gruid.Point{}
	}
}
