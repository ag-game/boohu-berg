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
	c := g.Dungeon.Cell(from)
	if c.T == WallCell {
		return wallcost
	}
	if _, ok := g.Clouds[from]; ok {
		return wallcost
	}
	if _, ok := g.Doors[from]; ok {
		if from != src {
			mons := g.MonsterAt(from)
			if !mons.Exists() && from != g.Player.P {
				return wallcost
			}
		}
	}
	if _, ok := g.Fungus[from]; ok {
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
	return g.Dungeon.Cell(p).T != WallCell
}

func (g *game) ComputeLOS() {
	if g.Player.LOS == nil {
		g.Player.LOS = map[gruid.Point]bool{}
	}
	for k := range g.Player.LOS {
		delete(g.Player.LOS, k)
	}
	losRange := g.LosRange()
	p := g.Player.P
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
			p := n.P
			g.Player.LOS[p] = true
			g.SeePosition(p)
		}
	}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.P] {
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

func (g *game) SeePosition(p gruid.Point) {
	if !g.Dungeon.Cell(p).Explored {
		see := "see"
		if c, ok := g.Collectables[p]; ok {
			if c.Quantity > 1 {
				g.Printf("You %s %d %s.", see, c.Quantity, c.Consumable.Plural())
			} else {
				g.Printf("You %s %s.", see, Indefinite(c.Consumable.String(), false))
			}
			g.StopAuto()
		} else if _, ok := g.Stairs[p]; ok {
			g.Printf("You %s stairs.", see)
			g.StopAuto()
		} else if eq, ok := g.Equipables[p]; ok {
			g.Printf("You %s %s.", see, Indefinite(eq.String(), false))
			g.StopAuto()
		} else if rd, ok := g.Rods[p]; ok {
			g.Printf("You %s %s.", see, Indefinite(rd.String(), false))
			g.StopAuto()
		} else if stn, ok := g.MagicalStones[p]; ok {
			g.Printf("You %s %s.", see, Indefinite(stn.String(), false))
			g.StopAuto()
		}
		g.FunAction()
		g.Dungeon.SetExplored(p)
		g.DijkstraMapRebuild = true
	} else {
		if g.WrongWall[p] {
			g.Printf("There is no longer a wall there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
		if cld, ok := g.Clouds[p]; ok && cld == CloudFire && (g.WrongDoor[p] || g.WrongFoliage[p]) {
			g.Printf("There are flames there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
	}
	if g.WrongWall[p] {
		delete(g.WrongWall, p)
		if g.Dungeon.Cell(p).T == FreeCell {
			delete(g.TemporalWalls, p)
		}
	}
	delete(g.WrongDoor, p)
	delete(g.WrongFoliage, p)
	delete(g.DreamingMonster, p)
}

func (g *game) ComputeExclusion(p gruid.Point, toggle bool) {
	exclusionRange := g.LosRange()
	g.ExclusionsMap[p] = toggle
	for d := 1; d <= exclusionRange; d++ {
		for x := -d + p.X; x <= d+p.X; x++ {
			for _, q := range []gruid.Point{{x, p.Y + d}, {x, p.Y - d}} {
				if !valid(q) {
					continue
				}
				g.ExclusionsMap[q] = toggle
			}
		}
		for y := -d + 1 + p.Y; y <= d-1+p.Y; y++ {
			for _, q := range []gruid.Point{{p.X + d, y}, {p.X - d, y}} {
				if !valid(q) {
					continue
				}
				g.ExclusionsMap[q] = toggle
			}
		}
	}
}

func (g *game) Ray(p gruid.Point) []gruid.Point {
	lt := &lighter{maxCost: g.LosRange(), g: g}
	lnodes := g.Player.FOV.Ray(lt, p)
	ps := []gruid.Point{}
	for i := len(lnodes) - 1; i > 0; i-- {
		ps = append(ps, lnodes[i].P)
	}
	return ps
}

func (g *game) ComputeRayHighlight(p gruid.Point) {
	g.Highlight = map[gruid.Point]bool{}
	ray := g.Ray(p)
	for _, q := range ray {
		g.Highlight[q] = true
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
	nodes := g.PR.BreadthFirstMap(dij, []gruid.Point{g.Player.P}, rg)
	count := 0
	noise := map[gruid.Point]bool{}
	rmax := 3
	if g.Player.Aptitudes[AptHear] {
		rmax--
	}
	for _, n := range nodes {
		p := n.P
		if g.Player.LOS[p] {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() && mons.State != Resting && RandInt(rmax) == 0 {
			switch mons.Kind {
			case MonsMirrorSpecter, MonsSatowalgaPlant:
				// no footsteps
			case MonsTinyHarpy, MonsWingedMilfid, MonsGiantBee:
				noise[p] = true
				g.Print("You hear the flapping of wings.")
				count++
			case MonsOgre, MonsCyclop, MonsBrizzia, MonsHydra, MonsEarthDragon, MonsTreeMushroom:
				noise[p] = true
				g.Print("You hear heavy footsteps.")
				count++
			case MonsWorm, MonsAcidMound:
				noise[p] = true
				g.Print("You hear a creep noise.")
				count++
			default:
				noise[p] = true
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
