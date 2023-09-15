package main

import "codeberg.org/anaseto/gruid"

func Neighbors(p gruid.Point, nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [8]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1}), p.Add(gruid.Point{1, -1}), p.Add(gruid.Point{-1, -1}), p.Add(gruid.Point{1, 1}), p.Add(gruid.Point{-1, 1})}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func CardinalNeighbors(p gruid.Point, nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [4]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1})}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func OutsideNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(p, nb, func(npos gruid.Point) bool {
		return !valid(npos)
	})

	return nb
}

func ValidNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(p, nb, valid)
	return nb
}

func (d *dungeon) IsFreeCell(p gruid.Point) bool {
	return valid(p) && d.Cell(p).T != WallCell
}

func (d *dungeon) FreeNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(p, nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalFreeNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(p, nb, d.IsFreeCell)
	return nb
}
