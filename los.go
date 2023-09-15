package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"codeberg.org/anaseto/gruid/rl"
)

func visionRange(p gruid.Point, radius int) gruid.Range {
	drg := gruid.NewRange(0, 0, DungeonWidth, DungeonHeight)
	delta := gruid.Point{radius, radius}
	return drg.Intersect(gruid.Range{Min: p.Sub(delta), Max: p.Add(delta).Shift(1, 1)})
}

type lighter struct {
	g       *game
	maxCost int
}

func (lt *lighter) Cost(src, from, to gruid.Point) int {
	g := lt.g
	wallcost := lt.maxCost
	// no terrain cost on origin
	if src == from {
		return paths.DistanceChebyshev(to, from)
	}
	// from terrain specific costs
	pfrom := point2Pos(from)
	c := g.Dungeon.Cell(pfrom)
	if c.T == WallCell {
		return wallcost
	}
	if _, ok := g.Clouds[pfrom]; ok {
		return wallcost
	}
	if _, ok := g.Doors[pfrom]; ok {
		if from != src {
			mons := g.MonsterAt(pfrom)
			if !mons.Exists() && pfrom != g.Player.Pos {
				return wallcost
			}
		}
	}
	if _, ok := g.Fungus[pfrom]; ok {
		return wallcost + paths.DistanceChebyshev(to, from) - 3
	}
	return paths.DistanceChebyshev(to, from)
}

func (lt *lighter) MaxCost(src gruid.Point) int {
	return lt.maxCost
}

func (g *game) LosRange() int {
	losRange := 6
	if g.Player.Armour == ShinyPlates {
		losRange++
	}
	if g.Player.Aptitudes[AptStealthyLOS] {
		losRange -= 2
	}
	if g.Player.Armour == HarmonistRobe {
		losRange -= 1
	}
	if g.Player.Weapon == Frundis {
		losRange -= 1
	}
	if g.Player.HasStatus(StatusShadows) {
		losRange = 1
	}
	if losRange < 1 {
		losRange = 1
	}
	return losRange
}

func (g *game) StopAuto() {
	if g.Autoexploring && !g.AutoHalt {
		g.Print("You stop exploring.")
	} else if g.AutoDir != NoDir {
		g.Print("You stop.")
	} else if g.AutoTarget != InvalidPos {
		g.Print("You stop.")
	}
	g.AutoHalt = true
	g.AutoDir = NoDir
	g.AutoTarget = InvalidPos
	if g.Resting {
		g.Stats.RestInterrupt++
		g.Resting = false
		g.Print("You could not sleep.")
	}
}

func (g *game) blocksSSCLOS(p gruid.Point) bool {
	return g.Dungeon.Cell(point2Pos(p)).T != WallCell
}

func (g *game) ComputeLOS() {
	if g.Player.LOS == nil {
		g.Player.LOS = map[position]bool{}
	}
	for k := range g.Player.LOS {
		delete(g.Player.LOS, k)
	}
	losRange := g.LosRange()
	p := pos2Point(g.Player.Pos)
	if g.Player.FOV == nil {
		g.Player.FOV = rl.NewFOV(visionRange(p, losRange))
	} else {
		g.Player.FOV.SetRange(visionRange(p, losRange))
	}
	lt := &lighter{g: g, maxCost: losRange}
	g.Player.FOV.SetRange(visionRange(p, losRange))
	lnodes := g.Player.FOV.VisionMap(lt, p)
	g.Player.FOV.SSCVisionMap(
		p, losRange,
		g.blocksSSCLOS,
		true,
	)
	for _, n := range lnodes {
		if !g.Player.FOV.Visible(n.P) {
			continue
		}
		if n.Cost <= losRange {
			pp := point2Pos(n.P)
			g.Player.LOS[pp] = true
			g.SeePosition(pp)
		}
	}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.Pos] {
			if mons.Seen {
				g.StopAuto()
				continue
			}
			mons.Seen = true
			g.Printf("You see %s (%v).", mons.Kind.Indefinite(false), mons.State)
			if mons.Kind.Dangerousness() > 10 {
				g.StoryPrint(mons.Kind.SeenStoryText())
			}
			g.StopAuto()
		}
	}
}

func (g *game) SeePosition(pos position) {
	if !g.Dungeon.Cell(pos).Explored {
		see := "see"
		if c, ok := g.Collectables[pos]; ok {
			if c.Quantity > 1 {
				g.Printf("You %s %d %s.", see, c.Quantity, c.Consumable.Plural())
			} else {
				g.Printf("You %s %s.", see, Indefinite(c.Consumable.String(), false))
			}
			g.StopAuto()
		} else if _, ok := g.Stairs[pos]; ok {
			g.Printf("You %s stairs.", see)
			g.StopAuto()
		} else if eq, ok := g.Equipables[pos]; ok {
			g.Printf("You %s %s.", see, Indefinite(eq.String(), false))
			g.StopAuto()
		} else if rd, ok := g.Rods[pos]; ok {
			g.Printf("You %s %s.", see, Indefinite(rd.String(), false))
			g.StopAuto()
		} else if stn, ok := g.MagicalStones[pos]; ok {
			g.Printf("You %s %s.", see, Indefinite(stn.String(), false))
			g.StopAuto()
		}
		g.FunAction()
		g.Dungeon.SetExplored(pos)
		g.DijkstraMapRebuild = true
	} else {
		if g.WrongWall[pos] {
			g.Printf("There is no longer a wall there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
		if cld, ok := g.Clouds[pos]; ok && cld == CloudFire && (g.WrongDoor[pos] || g.WrongFoliage[pos]) {
			g.Printf("There are flames there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
	}
	if g.WrongWall[pos] {
		delete(g.WrongWall, pos)
		if g.Dungeon.Cell(pos).T == FreeCell {
			delete(g.TemporalWalls, pos)
		}
	}
	if _, ok := g.WrongDoor[pos]; ok {
		delete(g.WrongDoor, pos)
	}
	if _, ok := g.WrongFoliage[pos]; ok {
		delete(g.WrongFoliage, pos)
	}
	if _, ok := g.DreamingMonster[pos]; ok {
		delete(g.DreamingMonster, pos)
	}
}

func (g *game) ComputeExclusion(pos position, toggle bool) {
	exclusionRange := g.LosRange()
	g.ExclusionsMap[pos] = toggle
	for d := 1; d <= exclusionRange; d++ {
		for x := -d + pos.X; x <= d+pos.X; x++ {
			for _, pos := range []position{{x, pos.Y + d}, {x, pos.Y - d}} {
				if !pos.valid() {
					continue
				}
				g.ExclusionsMap[pos] = toggle
			}
		}
		for y := -d + 1 + pos.Y; y <= d-1+pos.Y; y++ {
			for _, pos := range []position{{pos.X + d, y}, {pos.X - d, y}} {
				if !pos.valid() {
					continue
				}
				g.ExclusionsMap[pos] = toggle
			}
		}
	}
}

func (g *game) Ray(pos position) []position {
	lt := &lighter{maxCost: g.LosRange(), g: g}
	lnodes := g.Player.FOV.Ray(lt, pos2Point(pos))
	ps := []position{}
	for i := len(lnodes) - 1; i > 0; i-- {
		ps = append(ps, point2Pos(lnodes[i].P))
	}
	return ps
}

func (g *game) ComputeRayHighlight(pos position) {
	g.Highlight = map[position]bool{}
	ray := g.Ray(pos)
	for _, p := range ray {
		g.Highlight[p] = true
	}
}

func (g *game) ComputeNoise() {
	dij := &noisePath{game: g}
	rg := g.LosRange() + 2
	if rg <= 5 {
		rg++
	}
	if g.Player.Aptitudes[AptHear] {
		rg++
	}
	nodes := g.PR.BreadthFirstMap(dij, []gruid.Point{pos2Point(g.Player.Pos)}, rg)
	count := 0
	noise := map[position]bool{}
	rmax := 3
	if g.Player.Aptitudes[AptHear] {
		rmax--
	}
	for _, n := range nodes {
		pos := point2Pos(n.P)
		if g.Player.LOS[pos] {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() && mons.State != Resting && RandInt(rmax) == 0 {
			switch mons.Kind {
			case MonsMirrorSpecter, MonsSatowalgaPlant:
				// no footsteps
			case MonsTinyHarpy, MonsWingedMilfid, MonsGiantBee:
				noise[pos] = true
				g.Print("You hear the flapping of wings.")
				count++
			case MonsOgre, MonsCyclop, MonsBrizzia, MonsHydra, MonsEarthDragon, MonsTreeMushroom:
				noise[pos] = true
				g.Print("You hear heavy footsteps.")
				count++
			case MonsWorm, MonsAcidMound:
				noise[pos] = true
				g.Print("You hear a creep noise.")
				count++
			default:
				noise[pos] = true
				g.Print("You hear footsteps.")
				count++
			}
		}
	}
	if count > 0 {
		g.StopAuto()
	}
	g.Noise = noise
}
