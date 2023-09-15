package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"sort"
)

const unreachable = 9999

func valid(p gruid.Point) bool {
	return p.Y >= 0 && p.Y < DungeonHeight && p.X >= 0 && p.X < DungeonWidth
}

type dungeonPath struct {
	dungeon *dungeon
	wcost   int
	nbs     paths.Neighbors
}

func (dp *dungeonPath) Neighbors(p gruid.Point) []gruid.Point {
	return dp.nbs.All(p, valid)
}

func (dp *dungeonPath) Cost(from, to gruid.Point) int {
	if dp.dungeon.Cell(point2Pos(to)).T == WallCell {
		if dp.wcost > 0 {
			return dp.wcost
		}
		return 4
	}
	return 1
}

func (dp *dungeonPath) Estimation(from, to gruid.Point) int {
	return paths.DistanceChebyshev(from, to)
}

type playerPath struct {
	game *game
	nbs  paths.Neighbors
}

func (pp *playerPath) Neighbors(p gruid.Point) []gruid.Point {
	d := pp.game.Dungeon
	keep := func(np gruid.Point) bool {
		npos := point2Pos(np)
		if cld, ok := pp.game.Clouds[npos]; ok && cld == CloudFire && !(pp.game.WrongDoor[npos] || pp.game.WrongFoliage[npos]) {
			return false
		}
		return npos.valid() && ((d.Cell(npos).T == FreeCell && !pp.game.WrongWall[npos] || d.Cell(npos).T == WallCell && pp.game.WrongWall[npos]) || pp.game.Player.HasStatus(StatusDig)) &&
			d.Cell(npos).Explored
	}
	if pp.game.Player.HasStatus(StatusConfusion) {
		return pp.nbs.Cardinal(p, keep)
	}
	return pp.nbs.All(p, keep)
}

func (pp *playerPath) Cost(from, to gruid.Point) int {
	if !pp.game.ExclusionsMap[point2Pos(from)] && pp.game.ExclusionsMap[point2Pos(to)] {
		return unreachable
	}
	return 1
}

func (pp *playerPath) Estimation(from, to gruid.Point) int {
	return paths.DistanceChebyshev(from, to)
}

type noisePath struct {
	game *game
	nbs  paths.Neighbors
}

func (fp *noisePath) Neighbors(p gruid.Point) []gruid.Point {
	d := fp.game.Dungeon
	keep := func(np gruid.Point) bool {
		npos := point2Pos(np)
		return npos.valid() && d.Cell(npos).T != WallCell
	}
	return fp.nbs.All(p, keep)
}

func (fp *noisePath) Cost(from, to gruid.Point) int {
	return 1
}

type normalPath struct {
	game *game
	nbs  paths.Neighbors
}

func (np *normalPath) Neighbors(p gruid.Point) []gruid.Point {
	d := np.game.Dungeon
	keep := func(np gruid.Point) bool {
		npos := point2Pos(np)
		return npos.valid() && d.Cell(npos).T != WallCell
	}
	if np.game.Player.HasStatus(StatusConfusion) {
		return np.nbs.Cardinal(p, keep)
	}
	return np.nbs.All(p, keep)
}

func (np *normalPath) Cost(from, to gruid.Point) int {
	return 1
}

type autoexplorePath struct {
	game *game
	nbs  paths.Neighbors
}

func (ap *autoexplorePath) Neighbors(p gruid.Point) []gruid.Point {
	pos := point2Pos(p)
	if ap.game.ExclusionsMap[pos] {
		return nil
	}
	d := ap.game.Dungeon
	keep := func(np gruid.Point) bool {
		npos := point2Pos(np)
		if cld, ok := ap.game.Clouds[npos]; ok && cld == CloudFire && !(ap.game.WrongDoor[npos] || ap.game.WrongFoliage[npos]) {
			// XXX little info leak
			return false
		}
		return npos.valid() && (d.Cell(npos).T == FreeCell && !ap.game.WrongWall[npos] || d.Cell(npos).T == WallCell && ap.game.WrongWall[npos]) &&
			!ap.game.ExclusionsMap[npos]
	}
	if ap.game.Player.HasStatus(StatusConfusion) {
		return ap.nbs.Cardinal(p, keep)
	}
	return ap.nbs.All(p, keep)
}

func (ap *autoexplorePath) Cost(from, to gruid.Point) int {
	return 1
}

type monPath struct {
	game    *game
	monster *monster
	wall    bool
	nbs     paths.Neighbors
}

func (mp *monPath) Neighbors(p gruid.Point) []gruid.Point {
	d := mp.game.Dungeon
	keep := func(np gruid.Point) bool {
		npos := point2Pos(np)
		return npos.valid() && (d.Cell(npos).T != WallCell || mp.wall)
	}
	if mp.monster.Status(MonsConfused) {
		return mp.nbs.Cardinal(p, keep)
	}
	return mp.nbs.All(p, keep)
}

func (mp *monPath) Cost(from, to gruid.Point) int {
	g := mp.game
	mons := g.MonsterAt(point2Pos(to))
	if !mons.Exists() {
		if mp.wall && g.Dungeon.Cell(point2Pos(to)).T == WallCell && mp.monster.State != Hunting {
			return 6
		}
		return 1
	}
	if mons.Status(MonsLignified) {
		return 8
	}
	return 4
}

func (mp *monPath) Estimation(from, to gruid.Point) int {
	return paths.DistanceChebyshev(from, to)
}

func (m *monster) APath(g *game, from, to position) []position {
	mp := &monPath{game: g, monster: m}
	if m.Kind == MonsEarthDragon {
		mp.wall = true
	}
	path := points2Pos(g.PR.AstarPath(mp, pos2Point(from), pos2Point(to)))
	if len(path) == 0 {
		return nil
	}
	return path
}

func (g *game) PlayerPath(from, to position) []position {
	pp := &playerPath{game: g}
	path := points2Pos(g.PR.AstarPath(pp, pos2Point(from), pos2Point(to)))
	if len(path) == 0 {
		return nil
	}
	return path
}

func (g *game) SortedNearestTo(cells []position, to position) []position {
	ps := posSlice{}
	for _, p := range cells {
		pp := &dungeonPath{dungeon: g.Dungeon, wcost: unreachable}
		path := points2Pos(g.PR.AstarPath(pp, pos2Point(p), pos2Point(to)))
		if len(path) > 0 {
			ps = append(ps, posCost{p, len(path)})
		}
	}
	sort.Sort(ps)
	sorted := []position{}
	for _, pc := range ps {
		sorted = append(sorted, pc.pos)
	}
	return sorted
}

type posCost struct {
	pos  position
	cost int
}

type posSlice []posCost

func (ps posSlice) Len() int           { return len(ps) }
func (ps posSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps posSlice) Less(i, j int) bool { return ps[i].cost < ps[j].cost }
