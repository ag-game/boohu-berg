package main

import (
	"codeberg.org/anaseto/gruid"
	"errors"
)

var DijkstraMapCache [DungeonNCells]int

func (g *game) Autoexplore(ev event) error {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return errors.New("You cannot auto-explore while there are monsters in view.")
	}
	if g.ExclusionsMap[g.Player.Pos] {
		return errors.New("You cannot auto-explore while in an excluded area.")
	}
	if g.AllExplored() {
		return errors.New("Nothing left to explore.")
	}
	sources := g.AutoexploreSources()
	if len(sources) == 0 {
		return errors.New("Some excluded places remain unexplored.")
	}
	g.BuildAutoexploreMap(sources)
	n, finished := g.NextAuto()
	if finished || n == nil {
		return errors.New("You cannot reach some places safely.")
	}
	g.Autoexploring = true
	g.AutoHalt = false
	return g.MovePlayer(*n, ev)
}

func (g *game) AllExplored() bool {
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if c.T == WallCell {
			if len(np.Neighbors(pos2Point(pos))) == 0 {
				continue
			}
		}
		_, okc := g.Collectables[pos]
		if !c.Explored || g.Simellas[pos] > 0 || okc {
			return false
		} else if _, ok := g.Rods[pos]; ok {
			return false
		}
	}
	return true
}

func (g *game) AutoexploreSources() []gruid.Point {
	sources := []gruid.Point{}
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if c.T == WallCell {
			if len(np.Neighbors(pos2Point(pos))) == 0 {
				continue
			}
		}
		if g.ExclusionsMap[pos] {
			continue
		}
		_, okc := g.Collectables[pos]
		if !c.Explored || g.Simellas[pos] > 0 || okc {
			sources = append(sources, idxtopoint(i))
		} else if _, ok := g.Rods[pos]; ok {
			sources = append(sources, idxtopoint(i))
		}

	}
	return sources
}

func (g *game) BuildAutoexploreMap(sources []gruid.Point) {
	ap := &autoexplorePath{game: g}
	g.PRauto.BreadthFirstMap(ap, sources, unreachable)
	g.DijkstraMapRebuild = false
}

func (g *game) NextAuto() (next *position, finished bool) {
	ap := &autoexplorePath{game: g}
	if g.PRauto.BreadthFirstMapAt(pos2Point(g.Player.Pos)) > unreachable {
		return nil, false
	}
	neighbors := ap.Neighbors(pos2Point(g.Player.Pos))
	if len(neighbors) == 0 {
		return nil, false
	}
	n := neighbors[0]
	ncost := g.PRauto.BreadthFirstMapAt(n)
	for _, p := range neighbors[1:] {
		cost := g.PRauto.BreadthFirstMapAt(p)
		if cost < ncost {
			n = p
			ncost = cost
		}
	}
	if ncost >= g.PRauto.BreadthFirstMapAt(pos2Point(g.Player.Pos)) {
		finished = true
	}
	np := point2Pos(n)
	next = &np
	return next, finished
}
