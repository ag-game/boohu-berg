package main

import (
	"math/rand"
	"sort"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
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
	if dp.dungeon.Cell(to).T == WallCell {
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
	goal gruid.Point
}

func (pp *playerPath) Neighbors(p gruid.Point) []gruid.Point {
	d := pp.game.Dungeon
	keep := func(np gruid.Point) bool {
		if cld, ok := pp.game.Clouds[np]; ok && cld == CloudFire && !(pp.game.WrongDoor[np] || pp.game.WrongFoliage[np]) {
			return false
		}
		return valid(np) && ((d.Cell(np).T == FreeCell && !pp.game.WrongWall[np] || d.Cell(np).T == WallCell && pp.game.WrongWall[np]) || pp.game.Player.HasStatus(StatusDig)) &&
			d.Cell(np).Explored
	}
	var nb []gruid.Point
	if pp.game.Player.HasStatus(StatusConfusion) {
		nb = pp.nbs.Cardinal(p, keep)
	} else {
		nb = pp.nbs.All(p, keep)
	}
	sort.Slice(nb, func(i, j int) bool {
		return paths.DistanceManhattan(nb[i], pp.goal) <= paths.DistanceManhattan(nb[j], pp.goal)
	})
	return nb
}

func (pp *playerPath) Cost(from, to gruid.Point) int {
	if !pp.game.ExclusionsMap[from] && pp.game.ExclusionsMap[to] {
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
		return valid(np) && d.Cell(np).T != WallCell
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
		return valid(np) && d.Cell(np).T != WallCell
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
	if ap.game.ExclusionsMap[p] {
		return nil
	}
	d := ap.game.Dungeon
	keep := func(np gruid.Point) bool {
		if cld, ok := ap.game.Clouds[np]; ok && cld == CloudFire && !(ap.game.WrongDoor[np] || ap.game.WrongFoliage[np]) {
			// XXX little info leak
			return false
		}
		return valid(np) && (d.Cell(np).T == FreeCell && !ap.game.WrongWall[np] || d.Cell(np).T == WallCell && ap.game.WrongWall[np]) &&
			!ap.game.ExclusionsMap[np]
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
		return valid(np) && (d.Cell(np).T != WallCell || mp.wall)
	}
	var nb []gruid.Point
	if mp.monster.Status(MonsConfused) {
		nb = mp.nbs.Cardinal(p, keep)
	} else {
		nb = mp.nbs.All(p, keep)
	}
	rand.Shuffle(len(nb), func(i, j int) {
		nb[i], nb[j] = nb[j], nb[i]
	})
	return nb
}

func (mp *monPath) Cost(from, to gruid.Point) int {
	g := mp.game
	mons := g.MonsterAt(to)
	if !mons.Exists() {
		if mp.wall && g.Dungeon.Cell(to).T == WallCell && mp.monster.State != Hunting {
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

func (m *monster) APath(g *game, from, to gruid.Point) []gruid.Point {
	mp := &monPath{game: g, monster: m}
	if m.Kind == MonsEarthDragon {
		mp.wall = true
	}
	path := g.PR.AstarPath(mp, from, to)
	if len(path) == 0 {
		return nil
	}
	return path
}

func (g *game) PlayerPath(from, to gruid.Point) []gruid.Point {
	pp := &playerPath{game: g, goal: to}
	path := g.PR.AstarPath(pp, from, to)
	if len(path) == 0 {
		return nil
	}
	return path
}

func (g *game) SortedNearestTo(cells []gruid.Point, to gruid.Point) []gruid.Point {
	ps := posSlice{}
	for _, p := range cells {
		pp := &dungeonPath{dungeon: g.Dungeon, wcost: unreachable}
		path := g.PR.AstarPath(pp, p, to)
		if len(path) > 0 {
			ps = append(ps, posCost{p, len(path)})
		}
	}
	sort.Sort(ps)
	sorted := []gruid.Point{}
	for _, pc := range ps {
		sorted = append(sorted, pc.p)
	}
	return sorted
}

type posCost struct {
	p    gruid.Point
	cost int
}

type posSlice []posCost

func (ps posSlice) Len() int           { return len(ps) }
func (ps posSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps posSlice) Less(i, j int) bool { return ps[i].cost < ps[j].cost }
