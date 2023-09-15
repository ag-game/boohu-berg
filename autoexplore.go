package main

import (
	"errors"

	"codeberg.org/anaseto/gruid"
)

var DijkstraMapCache [DungeonNCells]int

func (g *game) Autoexplore(ev event) error {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return errors.New("You cannot auto-explore while there are monsters in view.")
	}
	if g.ExclusionsMap[g.Player.P] {
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
		p := idx2Point(i)
		if c.T == WallCell {
			if len(np.Neighbors(p)) == 0 {
				continue
			}
		}
		_, okc := g.Collectables[p]
		if !c.Explored || g.Simellas[p] > 0 || okc {
			return false
		} else if _, ok := g.Rods[p]; ok {
			return false
		}
	}
	return true
}

func (g *game) AutoexploreSources() []gruid.Point {
	sources := []gruid.Point{}
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		p := idx2Point(i)
		if c.T == WallCell {
			if len(np.Neighbors(p)) == 0 {
				continue
			}
		}
		if g.ExclusionsMap[p] {
			continue
		}
		_, okc := g.Collectables[p]
		if !c.Explored || g.Simellas[p] > 0 || okc {
			sources = append(sources, idx2Point(i))
		} else if _, ok := g.Rods[p]; ok {
			sources = append(sources, idx2Point(i))
		}

	}
	return sources
}

func (g *game) BuildAutoexploreMap(sources []gruid.Point) {
	ap := &autoexplorePath{game: g}
	g.PRauto.BreadthFirstMap(ap, sources, unreachable)
	g.DijkstraMapRebuild = false
}

func (g *game) NextAuto() (next *gruid.Point, finished bool) {
	ap := &autoexplorePath{game: g}
	if g.PRauto.BreadthFirstMapAt(g.Player.P) > unreachable {
		return nil, false
	}
	neighbors := ap.Neighbors(g.Player.P)
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
	if ncost >= g.PRauto.BreadthFirstMapAt(g.Player.P) {
		finished = true
	}
	next = &n
	return next, finished
}
