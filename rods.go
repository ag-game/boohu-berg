package main

import (
	"codeberg.org/anaseto/gruid"
	"errors"
	"fmt"
)

type rod int

const (
	RodDigging rod = iota
	RodBlink
	RodTeleportOther
	RodFireBolt
	RodFireBall
	RodLightning
	RodFog
	RodObstruction
	RodShatter
	RodSleeping
	RodLignification
	RodHope
	RodSwapping
)

const NumRods = int(RodSwapping) + 1

func (r rod) Letter() rune {
	return '/'
}

func (r rod) Name() string {
	var text string
	switch r {
	case RodDigging:
		text = "digging"
	case RodBlink:
		text = "blinking"
	case RodTeleportOther:
		text = "teleport other"
	case RodFog:
		text = "fog"
	case RodFireBall:
		text = "fireball"
	case RodFireBolt:
		text = "fire bolt"
	case RodLightning:
		text = "lightning"
	case RodObstruction:
		text = "obstruction"
	case RodShatter:
		text = "shatter"
	case RodSleeping:
		text = "sleeping"
	case RodLignification:
		text = "lignification"
	case RodHope:
		text = "last hope"
	case RodSwapping:
		text = "swapping"
	}
	return text
}

func (r rod) String() string {
	return "rod of " + r.Name()
}

func (r rod) Desc() string {
	var text string
	switch r {
	case RodDigging:
		text = "digs through up to 3 walls in a given direction."
	case RodBlink:
		text = "makes you blink away within your line of sight. The rod is more susceptible to send you to the cells thar are most far from you."
	case RodTeleportOther:
		text = "teleports away one of your foes. Note that the monster remembers where it saw you last time."
	case RodFog:
		text = "creates a dense fog that reduces your line of sight. Monsters at more than 1 cell away from you will not be able to see you."
	case RodFireBall:
		text = "throws a 1-radius fireball at your foes. You cannot use it against yourself. It can burn foliage and doors."
	case RodFireBolt:
		text = "throws a fire bolt through one or more enemies. It can burn foliage and doors."
	case RodLightning:
		text = "deals electrical damage to foes connected to you. It can burn foliage and doors."
	case RodObstruction:
		text = "creates a temporary wall at targeted location."
	case RodShatter:
		text = "induces an explosion around a wall, hurting adjacent monsters. The wall can disintegrate. You cannot use against yourself."
	case RodSleeping:
		text = "induces deep sleeping and exhaustion for monsters in the targeted area. You cannot use it against yourself."
	case RodLignification:
		text = "lignifies a monster, so that it cannot move, but can still fight with improved resistance."
	case RodHope:
		text = "creates an energy channel against a targeted monster. The damage done is inversely proportional to your health. It can burn foliage and doors."
	case RodSwapping:
		text = "makes you swap positions with a targeted monster."
	}
	return fmt.Sprintf("The %s %s Rods sometimes regain charges as you go deeper. This rod can have up to %d charges.", r, text, r.MaxCharge())
}

type rodProps struct {
	Charge int
}

func (r rod) MaxCharge() (charges int) {
	switch r {
	case RodBlink:
		charges = 5
	case RodDigging, RodShatter:
		charges = 3
	default:
		charges = 4
	}
	return charges
}

func (r rod) Rate() int {
	rate := r.MaxCharge() - 2
	if rate < 1 {
		rate = 1
	}
	return rate
}

func (r rod) MPCost() (mp int) {
	return 1
	//switch r {
	//case RodBlink:
	//mp = 3
	//case RodTeleportOther, RodDigging, RodShatter:
	//mp = 5
	//default:
	//mp = 4
	//}
	//return mp
}

func (r rod) Use(g *game, ev event) error {
	rods := g.Player.Rods
	if rods[r].Charge <= 0 {
		return errors.New("No charges remaining on this rod.")
	}
	if r.MPCost() > g.Player.MP {
		return errors.New("Not enough magic points for using this rod.")
	}
	if g.Player.HasStatus(StatusBerserk) {
		return errors.New("You cannot use rods while berserk.")
	}
	var err error
	switch r {
	case RodBlink:
		err = g.EvokeRodBlink(ev)
	case RodTeleportOther:
		err = g.EvokeRodTeleportOther(ev)
	case RodFireBolt:
		err = g.EvokeRodFireBolt(ev)
	case RodFireBall:
		err = g.EvokeRodFireball(ev)
	case RodLightning:
		err = g.EvokeRodLightning(ev)
	case RodFog:
		err = g.EvokeRodFog(ev)
	case RodDigging:
		err = g.EvokeRodDigging(ev)
	case RodObstruction:
		err = g.EvokeRodObstruction(ev)
	case RodShatter:
		err = g.EvokeRodShatter(ev)
	case RodSleeping:
		err = g.EvokeRodSleeping(ev)
	case RodLignification:
		err = g.EvokeRodLignification(ev)
	case RodHope:
		err = g.EvokeRodHope(ev)
	case RodSwapping:
		err = g.EvokeRodSwapping(ev)
	}

	if err != nil {
		return err
	}
	rp := rods[r]
	rp.Charge--
	rods[r] = rp
	g.Player.MP -= r.MPCost()
	g.StoryPrintf("Evoked your %s.", r)
	g.Stats.UsedRod[r]++
	g.Stats.Evocations++
	g.FunAction()
	ev.Renew(g, 7)
	return nil
}

func (g *game) EvokeRodBlink(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	g.Blink(ev)
	return nil
}

func (g *game) BlinkPos() gruid.Point {
	losPos := []gruid.Point{}
	for p, b := range g.Player.LOS {
		if !b {
			continue
		}
		if g.Dungeon.Cell(p).T != FreeCell {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		losPos = append(losPos, p)
	}
	if len(losPos) == 0 {
		return InvalidPos
	}
	npos := losPos[RandInt(len(losPos))]
	for i := 0; i < 4; i++ {
		p := losPos[RandInt(len(losPos))]
		if Distance(npos, g.Player.P) < Distance(p, g.Player.P) {
			npos = p
		}
	}
	return npos
}

func (g *game) Blink(ev event) {
	if g.Player.HasStatus(StatusLignification) {
		return
	}
	npos := g.BlinkPos()
	if !valid(npos) {
		// should not happen
		g.Print("You could not blink.")
		return
	}
	opos := g.Player.P
	g.Print("You blink away.")
	g.ui.TeleportAnimation(opos, npos, true)
	g.PlacePlayerAt(npos)
}

func (g *game) EvokeRodTeleportOther(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	mons.TeleportAway(g)
	return nil
}

func (g *game) EvokeRodSleeping(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Print("A sleeping ball emerges straight out of the rod.")
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgSleepingMonster)
	for _, p := range append(neighbors, g.Player.Target) {
		mons := g.MonsterAt(p)
		if !mons.Exists() {
			continue
		}
		if mons.State != Resting {
			g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		}
		mons.State = Resting
		mons.ExhaustTime(g, 40+RandInt(10))
	}
	return nil
}

func (g *game) EvokeRodFireBolt(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{flammable: true}); err != nil {
		return err
	}
	ray := g.Ray(g.Player.Target)
	g.MakeNoise(MagicCastNoise, g.Player.P)
	g.Print("Whoosh! A fire bolt emerges straight out of the rod.")
	g.ui.FireBoltAnimation(ray)
	for _, p := range ray {
		g.Burn(p, ev)
		mons := g.MonsterAt(p)
		if !mons.Exists() {
			continue
		}
		dmg := 0
		for i := 0; i < 2; i++ {
			dmg += RandInt(21)
		}
		dmg /= 2
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by the bolt.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(MagicHitNoise, mons.P)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodFireball(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true, flammable: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.MakeNoise(MagicExplosionNoise, g.Player.Target)
	g.Printf("A fireball emerges straight out of the rod... %s", g.ExplosionSound())
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionStart)
	g.ui.ExplosionAnimation(FireExplosion, g.Player.Target)
	for _, p := range append(neighbors, g.Player.Target) {
		g.Burn(p, ev)
		mons := g.MonsterAt(p)
		if mons == nil {
			continue
		}
		dmg := 0
		for i := 0; i < 2; i++ {
			dmg += RandInt(24)
		}
		dmg /= 2
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by the fireball.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(MagicHitNoise, mons.P)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodLightning(ev event) error {
	d := g.Dungeon
	conn := map[gruid.Point]bool{}
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(g.Player.P, nb, func(npos gruid.Point) bool {
		return valid(npos) && d.Cell(npos).T != WallCell
	})

	stack := []gruid.Point{}
	for _, p := range nb {
		mons := g.MonsterAt(p)
		if !mons.Exists() {
			continue
		}
		stack = append(stack, p)
		conn[p] = true
	}
	if len(stack) == 0 {
		return errors.New("There are no adjacent monsters.")
	}
	g.MakeNoise(MagicCastNoise, g.Player.P)
	g.Print("Whoosh! Lightning emerges straight out of the rod.")
	var p gruid.Point
	targets := []gruid.Point{}
	for len(stack) > 0 {
		p = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		g.Burn(p, ev)
		mons := g.MonsterAt(p)
		if !mons.Exists() {
			continue
		}
		targets = append(targets, p)
		dmg := 0
		for i := 0; i < 2; i++ {
			dmg += RandInt(17)
		}
		dmg /= 2
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by lightning.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(MagicHitNoise, mons.P)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
		nb = Neighbors(p, nb, func(npos gruid.Point) bool {
			return valid(npos) && d.Cell(npos).T != WallCell
		})

		for _, npos := range nb {
			if !conn[npos] {
				conn[npos] = true
				stack = append(stack, npos)
			}
		}
	}
	g.ui.LightningHitAnimation(targets)

	return nil
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (g *game) EvokeRodFog(ev event) error {
	g.Fog(g.Player.P, 3, ev)
	g.Print("You are surrounded by a dense fog.")
	return nil
}

func (g *game) Fog(at gruid.Point, radius int, ev event) {
	dij := &normalPath{game: g}
	nodes := g.PR.BreadthFirstMap(dij, []gruid.Point{at}, radius)
	for _, n := range nodes {
		p := n.P
		_, ok := g.Clouds[p]
		if !ok {
			g.Clouds[p] = CloudFog
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: CloudEnd, P: p})
		}
	}
	g.ComputeLOS()
}

func (g *game) EvokeRodDigging(ev event) error {
	if err := g.ui.ChooseTarget(&wallChooser{}); err != nil {
		return err
	}
	p := g.Player.Target
	for i := 0; i < 3; i++ {
		g.Dungeon.SetCell(p, FreeCell)
		g.Stats.Digs++
		g.MakeNoise(WallNoise, p)
		g.Fog(p, 1, ev)
		p = To(p, Dir(p, g.Player.P))
		if !g.Player.LOS[p] {
			g.WrongWall[p] = true
		}
		if !valid(p) || g.Dungeon.Cell(p).T != WallCell {
			break
		}
	}
	g.Print("You see the wall disintegrate with a crash.")
	g.ComputeLOS()
	g.MakeMonstersAware()
	return nil
}

func (g *game) EvokeRodShatter(ev event) error {
	if err := g.ui.ChooseTarget(&wallChooser{minDist: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Dungeon.SetCell(g.Player.Target, FreeCell)
	g.Stats.Digs++
	g.ComputeLOS()
	g.MakeMonstersAware()
	g.MakeNoise(WallNoise, g.Player.Target)
	g.Printf("%s The wall disappeared.", g.CrackSound())
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionWallStart)
	g.ui.ExplosionAnimation(WallExplosion, g.Player.Target)
	g.Fog(g.Player.Target, 2, ev)
	for _, p := range neighbors {
		mons := g.MonsterAt(p)
		if !mons.Exists() {
			continue
		}
		dmg := 0
		for i := 0; i < 3; i++ {
			dmg += RandInt(30)
		}
		dmg /= 3
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by the explosion.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(ExplosionHitNoise, mons.P)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodObstruction(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{free: true}); err != nil {
		return err
	}
	g.TemporalWallAt(g.Player.Target, ev)
	g.Printf("You see a wall appear out of thin air.")
	return nil
}

func (g *game) EvokeRodLignification(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in targeter)
	if mons.Status(MonsLignified) {
		return errors.New("You cannot target a lignified monster.")
	}
	mons.EnterLignification(g, ev)
	return nil
}

func (g *game) TemporalWallAt(p gruid.Point, ev event) {
	if g.Dungeon.Cell(p).T == WallCell {
		return
	}
	if !g.Player.LOS[p] {
		g.WrongWall[p] = true
	}
	g.CreateTemporalWallAt(p, ev)
	g.ComputeLOS()
}

func (g *game) CreateTemporalWallAt(p gruid.Point, ev event) {
	g.Dungeon.SetCell(p, WallCell)
	delete(g.Clouds, p)
	g.TemporalWalls[p] = true
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + 200 + RandInt(50), P: p, EAction: ObstructionEnd})
}

func (g *game) EvokeRodHope(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{needsFreeWay: true}); err != nil {
		return err
	}
	g.MakeNoise(MagicCastNoise, g.Player.P)
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionStart)
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	attack := -20 + 30*DefaultHealth/g.Player.HP
	if attack > 130 {
		attack = 130
	}
	dmg := 0
	for i := 0; i < 5; i++ {
		dmg += RandInt(attack)
	}
	dmg /= 5
	if dmg < 0 {
		// should not happen
		dmg = 0
	}
	mons.HP -= dmg
	g.Burn(g.Player.Target, ev)
	g.ui.HitAnimation(g.Player.Target, true)
	g.Printf("An energy channel hits %s (%d dmg).", mons.Kind.Definite(false), dmg)
	if mons.HP <= 0 {
		g.Printf("%s dies.", mons.Kind.Indefinite(true))
		g.HandleKill(mons, ev)
	}
	return nil
}

func (g *game) EvokeRodSwapping(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this rod while lignified.")
	}
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	if mons.Status(MonsLignified) {
		return errors.New("You cannot target a lignified monster.")
	}
	g.SwapWithMonster(mons)
	return nil
}

func (g *game) SwapWithMonster(mons *monster) {
	ompos := mons.P
	g.Printf("You swap positions with the %s.", mons.Kind)
	g.ui.SwappingAnimation(mons.P, g.Player.P)
	mons.MoveTo(g, g.Player.P)
	g.PlacePlayerAt(ompos)
	mons.MakeAware(g)
}

func (g *game) GeneratedRodsCount() int {
	count := 0
	for _, b := range g.GeneratedRods {
		if b {
			count++
		}
	}
	return count
}

func (g *game) RandomRod() rod {
	r := rod(RandInt(NumRods))
	return r
}

func (g *game) GenerateRod() {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("GenerateRod")
		}
		p := g.FreeCellForStatic()
		r := g.RandomRod()
		if _, ok := g.Player.Rods[r]; !ok && !g.GeneratedRods[r] {
			g.GeneratedRods[r] = true
			g.Rods[p] = r
			return
		}
	}
}

func (g *game) RechargeRods() {
	for r, props := range g.Player.Rods {
		max := r.MaxCharge()
		if g.Player.Armour == CelmistRobe {
			max += 2
		}
		if props.Charge < max {
			rchg := RandInt(1 + r.Rate())
			if rchg == 0 && RandInt(2) == 0 {
				rchg++
			}
			if g.Player.Armour == CelmistRobe {
				if RandInt(10) > 0 {
					rchg++
				}
				if RandInt(3) == 0 {
					rchg++
				}
			}
			props.Charge += rchg
			g.Player.Rods[r] = props
		}
		if props.Charge > max {
			props.Charge = max
			g.Player.Rods[r] = props
		}
	}
}
