package main

import (
	"codeberg.org/anaseto/gruid"
	"errors"
)

type Targeter interface {
	ComputeHighlight(*game, gruid.Point)
	Action(*game, gruid.Point) error
	Reachable(*game, gruid.Point) bool
	Done() bool
}

type examiner struct {
	done   bool
	stairs bool
}

func (ex *examiner) ComputeHighlight(g *game, p gruid.Point) {
	g.ComputePathHighlight(p)
}

func (g *game) ComputePathHighlight(p gruid.Point) {
	path := g.PlayerPath(g.Player.P, p)
	g.Highlight = map[gruid.Point]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (ex *examiner) Action(g *game, p gruid.Point) error {
	if !g.Dungeon.Cell(p).Explored {
		return errors.New("You do not know this place.")
	}
	if g.Dungeon.Cell(p).T == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot travel into a wall.")
	}
	path := g.PlayerPath(g.Player.P, p)
	if len(path) == 0 {
		if ex.stairs {
			return errors.New("There is no safe path to the nearest stairs.")
		}
		return errors.New("There is no safe path to this place.")
	}
	if c := g.Dungeon.Cell(p); c.Explored && c.T == FreeCell {
		g.AutoTarget = p
		g.Targeting = p
		ex.done = true
		return nil
	}
	return errors.New("Invalid destination.")
}

func (ex *examiner) Reachable(g *game, p gruid.Point) bool {
	return true
}

func (ex *examiner) Done() bool {
	return ex.done
}

type chooser struct {
	done         bool
	area         bool
	minDist      bool
	needsFreeWay bool
	free         bool
	flammable    bool
	wall         bool
}

func (ch *chooser) ComputeHighlight(g *game, p gruid.Point) {
	g.ComputeRayHighlight(p)
	if !ch.area {
		return
	}
	neighbors := g.Dungeon.FreeNeighbors(p)
	for _, p := range neighbors {
		g.Highlight[p] = true
	}
}

func (ch *chooser) Reachable(g *game, p gruid.Point) bool {
	return g.Player.LOS[p]
}

func (ch *chooser) Action(g *game, p gruid.Point) error {
	if !ch.Reachable(g, p) {
		return errors.New("You cannot target that place.")
	}
	if ch.minDist && Distance(p, g.Player.P) <= 1 {
		return errors.New("Invalid target: too close.")
	}
	c := g.Dungeon.Cell(p)
	if c.T == WallCell {
		return errors.New("You cannot target a wall.")
	}
	if (ch.area || ch.needsFreeWay) && !ch.freeWay(g, p) {
		return errors.New("Invalid target: there are monsters in the way.")
	}
	mons := g.MonsterAt(p)
	if ch.free {
		if mons.Exists() {
			return errors.New("Invalid target: there is a monster there.")
		}
		if g.Player.P == p {
			return errors.New("Invalid target: you are here.")
		}
		g.Player.Target = p
		ch.done = true
		return nil
	}
	if mons.Exists() || ch.flammable && ch.flammableInWay(g, p) {
		g.Player.Target = p
		ch.done = true
		return nil
	}
	if ch.flammable && ch.flammableInWay(g, p) {
		g.Player.Target = p
		ch.done = true
		return nil
	}
	if !ch.area {
		return errors.New("You must target a monster.")
	}
	neighbors := ValidNeighbors(p)
	for _, npos := range neighbors {
		nc := g.Dungeon.Cell(npos)
		if !ch.wall && nc.T == WallCell {
			continue
		}
		mons := g.MonsterAt(npos)
		_, okFungus := g.Fungus[p]
		_, okDoors := g.Doors[p]
		if ch.flammable && (okFungus || okDoors) || mons.Exists() || nc.T == WallCell {
			g.Player.Target = p
			ch.done = true
			return nil
		}
	}
	if ch.flammable && ch.wall {
		return errors.New("Invalid target: no monsters, walls nor flammable terrain in the area.")
	}
	if ch.flammable {
		return errors.New("Invalid target: no monsters nor flammable terrain in the area.")
	}
	if ch.wall {
		return errors.New("Invalid target: no monsters nor walls in the area.")
	}
	return errors.New("Invalid target: no monsters in the area.")
}

func (ch *chooser) Done() bool {
	return ch.done
}

func (ch *chooser) freeWay(g *game, p gruid.Point) bool {
	ray := g.Ray(p)
	tpos := p
	for _, rpos := range ray {
		mons := g.MonsterAt(rpos)
		if !mons.Exists() {
			continue
		}
		tpos = mons.P
	}
	return tpos == p
}

func (ch *chooser) flammableInWay(g *game, p gruid.Point) bool {
	ray := g.Ray(p)
	for _, rpos := range ray {
		if rpos == g.Player.P {
			continue
		}
		if _, ok := g.Fungus[rpos]; ok {
			return true
		}
		if _, ok := g.Doors[rpos]; ok {
			return true
		}
	}
	return false
}

type wallChooser struct {
	done    bool
	minDist bool
}

func (ch *wallChooser) ComputeHighlight(g *game, p gruid.Point) {
	g.ComputeRayHighlight(p)
}

func (ch *wallChooser) Reachable(g *game, p gruid.Point) bool {
	return g.Player.LOS[p]
}

func (ch *wallChooser) Action(g *game, p gruid.Point) error {
	if !ch.Reachable(g, p) {
		return errors.New("You cannot target that place.")
	}
	ray := g.Ray(p)
	if len(ray) == 0 {
		return errors.New("You are not a wall.")
	}
	if g.Dungeon.Cell(ray[0]).T != WallCell {
		return errors.New("You must target a wall.")
	}
	if ch.minDist && Distance(g.Player.P, p) <= 1 {
		return errors.New("You cannot target an adjacent wall.")
	}
	for _, p := range ray[1:] {
		mons := g.MonsterAt(p)
		if mons.Exists() {
			return errors.New("There are monsters in the way.")
		}
	}
	g.Player.Target = p
	ch.done = true
	return nil
}

func (ch *wallChooser) Done() bool {
	return ch.done
}
