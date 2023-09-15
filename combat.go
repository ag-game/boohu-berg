// combat utility functions

package main

import "codeberg.org/anaseto/gruid"

func (g *game) Absorb(armor int) int {
	absorb := 0
	for i := 0; i <= 2; i++ {
		absorb += RandInt(armor + 1)
	}
	q := absorb / 3
	r := absorb % 3
	if r == 2 {
		q++
	}
	return q
}

func (g *game) HitDamage(dt dmgType, base int, armor int) (attack int, clang bool) {
	min := base / 2
	attack = min + RandInt(base-min+1)
	absorb := g.Absorb(armor)
	if dt == DmgMagical {
		absorb = 2 * absorb / 3
	}
	attack -= absorb
	if absorb > 0 && absorb >= 2*armor/3 && RandInt(2) == 0 {
		clang = true
	}
	if attack < 0 {
		attack = 0
	}
	return attack, clang
}

func (m *monster) InflictDamage(g *game, damage, max int) {
	g.Stats.ReceivedHits++
	g.Stats.Damage += damage
	oldHP := g.Player.HP
	g.Player.HP -= damage
	g.ui.WoundedAnimation()
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical HP: %d (hit by %s)", g.Player.HP, m.Kind.Indefinite(false))
		g.ui.CriticalHPWarning()
	}
	if g.Player.HP <= 0 {
		return
	}
	stn, ok := g.MagicalStones[g.Player.P]
	if !ok {
		return
	}
	switch stn {
	case TeleStone:
		g.UseStone(g.Player.P)
		g.Teleportation(g.Ev)
	case FogStone:
		g.Fog(g.Player.P, 3, g.Ev)
		g.UseStone(g.Player.P)
	case QueenStone:
		g.MakeNoise(QueenStoneNoise, g.Player.P)
		dij := &normalPath{game: g}
		const radius = 2
		g.PR.BreadthFirstMap(dij, []gruid.Point{g.Player.P}, radius)
		for _, m := range g.Monsters {
			if !m.Exists() {
				continue
			}
			if m.State == Resting {
				continue
			}
			c := g.PR.BreadthFirstMapAt(m.P)
			if c > radius {
				continue
			}
			m.EnterConfusion(g, g.Ev)
		}
		//g.Confusion(g.Ev)
		g.UseStone(g.Player.P)
	case TreeStone:
		if !g.Player.HasStatus(StatusLignification) {
			g.UseStone(g.Player.P)
			g.EnterLignification(g.Ev)
			g.Print("You feel rooted to the ground.")
		}
	case ObstructionStone:
		neighbors := g.Dungeon.FreeNeighbors(g.Player.P)
		for _, p := range neighbors {
			mons := g.MonsterAt(p)
			if mons.Exists() {
				continue
			}
			g.CreateTemporalWallAt(p, g.Ev)
		}
		g.Printf("You see walls appear out of thin air around the stone.")
		g.UseStone(g.Player.P)
		g.ComputeLOS()
	}
}

func (g *game) MakeMonstersAware() {
	for _, m := range g.Monsters {
		if m.HP <= 0 {
			continue
		}
		if g.Player.LOS[m.P] {
			m.MakeAware(g)
			if m.State != Resting {
				m.GatherBand(g)
			}
		}
	}
}

func (g *game) MakeNoise(noise int, at gruid.Point) {
	dij := &normalPath{game: g}
	g.PR.BreadthFirstMap(dij, []gruid.Point{at}, noise)
	for _, m := range g.Monsters {
		if !m.Exists() {
			continue
		}
		if m.State == Hunting {
			continue
		}
		c := g.PR.BreadthFirstMapAt(m.P)
		if c > noise {
			continue
		}
		v := noise - c
		if v <= 0 {
			continue
		}
		if v > 25 {
			v = 25
		}
		r := RandInt(30)
		if m.State == Resting {
			v /= 2
		}
		if m.Status(MonsExhausted) {
			v = 2 * v / 3
		}
		if v > r {
			if g.Player.LOS[m.P] {
				m.MakeHunt(g)
			} else {
				m.Target = at
				m.State = Wandering
			}
			m.GatherBand(g)
		}
	}
}

func (g *game) InOpenMons(mons *monster) bool {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.P)
	for _, p := range neighbors {
		if Distance(p, mons.P) > 1 {
			continue
		}
		if g.Dungeon.Cell(p).T == WallCell {
			return false
		}
	}
	return true
}

func (g *game) AttackMonster(mons *monster, ev event) {
	switch {
	case g.Player.HasStatus(StatusSwap) && !g.Player.HasStatus(StatusLignification) && !mons.Status(MonsLignified):
		g.SwapWithMonster(mons)
	case g.Player.Weapon == Frundis:
		if !g.HitMonster(DmgPhysical, g.Player.Attack(), mons, ev) {
			break
		}
		if RandInt(2) == 0 {
			mons.EnterConfusion(g, ev)
			g.PrintfStyled("Frundis glows… %s appears confused.", logPlayerHit, mons.Kind.Definite(false))
		}
	case g.Player.Weapon.Cleave():
		var neighbors []gruid.Point
		if g.Player.HasStatus(StatusConfusion) {
			neighbors = g.Dungeon.CardinalFreeNeighbors(g.Player.P)
		} else {
			neighbors = g.Dungeon.FreeNeighbors(g.Player.P)
		}
		for _, p := range neighbors {
			m := g.MonsterAt(p)
			if m.Exists() {
				g.HitMonster(DmgPhysical, g.Player.Attack(), m, ev)
			}
		}
	case g.Player.Weapon.Pierce():
		g.HitMonster(DmgPhysical, g.Player.Attack(), mons, ev)
		dir := Dir(mons.P, g.Player.P)
		behind := To(To(g.Player.P, dir), dir)
		if valid(behind) {
			m := g.MonsterAt(behind)
			if m.Exists() {
				g.HitMonster(DmgPhysical, g.Player.Attack(), m, ev)
			}
		}
	case g.Player.Weapon == ElecWhip:
		g.HitConnected(mons.P, DmgMagical, ev)
	case g.Player.Weapon == DancingRapier:
		ompos := mons.P
		g.HitMonster(DmgPhysical, g.Player.Attack(), mons, ev)
		if g.Player.HasStatus(StatusLignification) || mons.Status(MonsLignified) || mons.Kind == MonsTinyHarpy {
			break
		}
		dir := Dir(ompos, g.Player.P)
		behind := To(To(g.Player.P, dir), dir)
		if valid(behind) {
			m := g.MonsterAt(behind)
			if m.Exists() {
				g.HitMonster(DmgPhysical, g.Player.Attack()+3, m, ev)
			}
		}
		if mons.Exists() {
			mons.MoveTo(g, g.Player.P)
		}
		g.PlacePlayerAt(ompos)
	case g.Player.Weapon == HarKarGauntlets:
		g.HarKarAttack(mons, ev)
	case g.Player.Weapon == HopeSword:
		attack := g.Player.Attack()
		fact := -60 + 100*DefaultHealth/g.Player.HP
		if fact < 100 {
			fact = 100
		}
		if fact > 250 {
			fact = 250
		}
		attack *= fact
		attack /= 100
		g.HitMonster(DmgPhysical, attack, mons, ev)
	case g.Player.Weapon == DragonSabre:
		mfact := 100 * (mons.HPmax * mons.HPmax) / (45 * 45)
		bonus := -1 + 13*mfact/100
		g.HitMonster(DmgPhysical, g.Player.Attack()+bonus, mons, ev)
	case g.Player.Weapon == DefenderFlail:
		bonus := g.Player.Statuses[StatusSlay]
		g.HitMonster(DmgPhysical, g.Player.Attack()+bonus, mons, ev)
		g.Player.Statuses[StatusSlay]++
		g.PushEvent(&simpleEvent{ERank: ev.Rank() + 60, EAction: SlayEnd})
	default:
		g.HitMonster(DmgPhysical, g.Player.Attack(), mons, ev)
	}
}

func (g *game) AttractMonster(p gruid.Point) *monster {
	dir := Dir(p, g.Player.P)
	for cpos := To(p, dir); g.Player.LOS[cpos]; cpos = To(cpos, dir) {
		mons := g.MonsterAt(cpos)
		if mons.Exists() {
			mons.MoveTo(g, p)
			g.ui.TeleportAnimation(cpos, p, false)
			return mons
		}
	}
	return nil
}

func (g *game) HarKarAttack(mons *monster, ev event) {
	dir := Dir(mons.P, g.Player.P)
	p := g.Player.P
	for {
		p = To(p, dir)
		if !valid(p) || g.Dungeon.Cell(p).T != FreeCell {
			break
		}
		m := g.MonsterAt(p)
		if !m.Exists() {
			break
		}
	}
	if valid(p) && g.Dungeon.Cell(p).T == FreeCell && !g.Player.HasStatus(StatusLignification) {
		p = g.Player.P
		for {
			p = To(p, dir)
			if !valid(p) || g.Dungeon.Cell(p).T != FreeCell {
				break
			}
			m := g.MonsterAt(p)
			if !m.Exists() {
				break
			}
			g.HitMonster(DmgPhysical, g.Player.Attack(), m, ev)
		}
		if !valid(p) || g.Dungeon.Cell(p).T != FreeCell {
			return
		}
		g.PlacePlayerAt(p)
		behind := To(p, dir)
		m := g.MonsterAt(behind)
		if m.Exists() {
			g.HitMonster(DmgPhysical, g.Player.Attack(), m, ev)
		}
	} else {
		g.HitMonster(DmgPhysical, g.Player.Attack(), mons, ev)
	}
}

func (g *game) HitConnected(p gruid.Point, dt dmgType, ev event) {
	d := g.Dungeon
	conn := map[gruid.Point]bool{}
	stack := []gruid.Point{p}
	conn[p] = true
	nb := make([]gruid.Point, 0, 8)
	for len(stack) > 0 {
		p = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		mons := g.MonsterAt(p)
		if !mons.Exists() {
			continue
		}
		g.HitMonster(dt, g.Player.Attack(), mons, ev)
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
}

func (g *game) HitNoise(clang bool) int {
	noise := BaseHitNoise
	if g.Player.Weapon == Frundis {
		noise -= 5
	}
	if g.Player.Armour == HarmonistRobe {
		noise -= 3
	}
	if g.Player.Armour == Robe {
		noise -= 1
	}
	if clang {
		arnoise := g.Player.Armor()
		if arnoise > 7 {
			arnoise = 7
		}
		noise += arnoise
	}
	return noise
}

type dmgType int

const (
	DmgPhysical dmgType = iota
	DmgMagical
)

func (g *game) HitMonster(dt dmgType, dmg int, mons *monster, ev event) (hit bool) {
	maxacc := g.Player.Accuracy()
	if g.Player.Weapon == AssassinSabre && mons.HP > 0 {
		adjust := 6 * (-100 + 100*mons.HPmax/mons.HP) / 100
		if adjust > 25 {
			adjust = 25
		}
		maxacc += adjust
	} else if g.Player.Weapon == FinalBlade {
		maxacc += 10
	}
	acc := RandInt(maxacc)
	if g.Player.AccScore == 1 && acc >= maxacc/2 {
		acc -= RandInt(1 + maxacc/2)
	} else if g.Player.AccScore == -1 && acc < maxacc/2 {
		acc += RandInt(1 + maxacc/2)
	}
	if acc >= maxacc/2 {
		g.Player.AccScore = 1
	} else {
		g.Player.AccScore = -1
	}
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion || g.Player.HasStatus(StatusAccurate) {
		hit = true
		noise := BaseHitNoise
		if g.Player.Weapon == Dagger || g.Player.Weapon == VampDagger {
			noise -= 2
		}
		if g.Player.Armour == HarmonistRobe {
			noise -= 3
		}
		if g.Player.Weapon == Frundis {
			noise -= 5
		}
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += 2 + RandInt(4)
		}
		pa := dmg + bonus
		if g.Player.Weapon.Cleave() && g.InOpenMons(mons) {
			if g.Player.Attack() >= 15 {
				pa += 1 + RandInt(3)
			} else {
				pa += 1 + RandInt(2)
			}
		}
		marmor := mons.Armor
		marmor = 6 + marmor/2
		attack, clang := g.HitDamage(dt, pa, marmor)
		if clang {
			noise += marmor
		}
		g.MakeNoise(noise, mons.P)
		if mons.State == Resting {
			if g.Player.Weapon == Dagger || g.Player.Weapon == VampDagger {
				attack *= 4
			} else {
				attack *= 2
			}
		}
		var sclang string
		if clang {
			if marmor > 3 {
				sclang = " ♫ Clang!"
			} else {
				sclang = " ♪ Clang!"
			}
		}
		oldHP := mons.HP
		if g.Player.Weapon == FinalBlade {
			if mons.HP <= mons.HPmax/2 {
				attack = mons.HP
			}
		}
		mons.HP -= attack
		if g.Player.Weapon == VampDagger && mons.Kind.Living() {
			healing := attack
			if healing > 2*pa/3 {
				healing = 2 * pa / 3
			}
			if g.Player.HP+healing > g.Player.HPMax() {
				g.Player.HP = g.Player.HPMax()
			} else {
				g.Player.HP += healing
			}
		}
		g.ui.HitAnimation(mons.P, false)
		if mons.HP > 0 {
			g.PrintfStyled("You hit %s (%d dmg).%s", logPlayerHit, mons.Kind.Definite(false), attack, sclang)
		} else if oldHP > 0 {
			// test oldHP > 0 because of sword special attack
			g.PrintfStyled("You kill %s (%d dmg).%s", logPlayerHit, mons.Kind.Definite(false), attack, sclang)
			g.HandleKill(mons, ev)
		}
		if mons.Kind == MonsBrizzia && RandInt(4) == 0 && !g.Player.HasStatus(StatusNausea) &&
			Distance(mons.P, g.Player.P) == 1 {
			g.Player.Statuses[StatusNausea]++
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + 30 + RandInt(20), EAction: NauseaEnd})
			g.Print("The brizzia's corpse releases some nauseating gas. You feel sick.")
		}
		if mons.Kind == MonsTinyHarpy && mons.HP > 0 {
			mons.Blink(g)
		}
		g.HandleStone(mons)
		g.Stats.Hits++
	} else {
		g.Printf("You miss %s.", mons.Kind.Definite(false))
		g.Stats.Misses++
	}
	mons.MakeHuntIfHurt(g)
	return hit
}

func (g *game) HandleStone(mons *monster) {
	stn, ok := g.MagicalStones[mons.P]
	if !ok {
		return
	}
	switch stn {
	case TeleStone:
		if mons.Exists() {
			g.UseStone(mons.P)
			mons.TeleportAway(g)
		}
	case FogStone:
		g.Fog(mons.P, 3, g.Ev)
		g.UseStone(mons.P)
	case QueenStone:
		g.MakeNoise(QueenStoneNoise, mons.P)
		dij := &normalPath{game: g}
		const radius = 2
		g.PR.BreadthFirstMap(dij, []gruid.Point{mons.P}, radius)
		for _, m := range g.Monsters {
			if !m.Exists() {
				continue
			}
			if m.State == Resting {
				continue
			}
			c := g.PR.BreadthFirstMapAt(m.P)
			if c > radius {
				continue
			}
			m.EnterConfusion(g, g.Ev)
		}
		// _, ok := nm[g.Player.P]
		// if ok {
		// 	g.Confusion(g.Ev)
		// }
		g.UseStone(mons.P)
	case TreeStone:
		if mons.Exists() {
			g.UseStone(mons.P)
			mons.EnterLignification(g, g.Ev)
		}
	case ObstructionStone:
		if !mons.Exists() {
			g.CreateTemporalWallAt(mons.P, g.Ev)
		}
		neighbors := g.Dungeon.FreeNeighbors(mons.P)
		for _, p := range neighbors {
			if p == g.Player.P {
				continue
			}
			m := g.MonsterAt(p)
			if m.Exists() {
				continue
			}
			g.CreateTemporalWallAt(p, g.Ev)
		}
		g.Printf("You see walls appear out of thin air around the stone.")
		g.UseStone(mons.P)
		g.ComputeLOS()
	}
}

func (g *game) HandleKill(mons *monster, ev event) {
	g.Stats.Killed++
	g.Stats.KilledMons[mons.Kind]++
	if mons.Kind == MonsExplosiveNadre {
		mons.Explode(g, ev)
	}
	if g.Doors[mons.P] {
		g.ComputeLOS()
	}
	if mons.Kind.Dangerousness() > 10 {
		g.StoryPrintf("Killed %s.", mons.Kind.Indefinite(false))
	}
}

const (
	WallNoise           = 18
	TemporalWallNoise   = 13
	ExplosionHitNoise   = 13
	ExplosionNoise      = 18
	MagicHitNoise       = 15
	BarkNoise           = 13
	MagicExplosionNoise = 16
	MagicCastNoise      = 16
	BaseHitNoise        = 11
	ShieldBlockNoise    = 17
	QueenStoneNoise     = 19
)

func (g *game) ArmourClang() (sclang string) {
	if g.Player.Armor() > 3 {
		sclang = " Clang!"
	} else {
		sclang = " Smash!"
	}
	return sclang
}

func (g *game) BlockEffects(m *monster) {
	g.Stats.Blocks++
	// only one shield block per turn
	g.Player.Blocked = true
	g.PushEvent(&simpleEvent{ERank: g.Ev.Rank() + 10, EAction: BlockEnd})
	switch g.Player.Shield {
	case EarthShield:
		dir := Dir(m.P, g.Player.P)
		lat := Laterals(g.Player.P, dir)
		for _, p := range lat {
			if !valid(p) {
				continue
			}
			if RandInt(3) == 0 && g.Dungeon.Cell(p).T == WallCell {
				g.Dungeon.SetCell(p, FreeCell)
				g.Stats.Digs++
				g.MakeNoise(WallNoise+3, p)
				g.Fog(p, 1, g.Ev)
				g.Printf("%s The sound of blocking breaks a wall.", g.CrackSound())
			}
		}
	case BashingShield:
		if m.Kind == MonsSatowalgaPlant || Distance(m.P, g.Player.P) > 1 {
			break
		}
		if RandInt(5) == 0 {
			break
		}
		dir := Dir(m.P, g.Player.P)
		p := m.P
		npos := p
		i := 0
		for {
			i++
			npos = To(npos, dir)
			if !valid(npos) || g.Dungeon.Cell(npos).T == WallCell {
				break
			}
			mons := g.MonsterAt(npos)
			if mons.Exists() {
				continue
			}
			p = npos
			if i >= 3 {
				break
			}
		}
		m.Exhaust(g)
		if p != m.P {
			m.MoveTo(g, p)
			g.Printf("%s is repelled.", m.Kind.Definite(true))
		}
	case ConfusingShield:
		if Distance(m.P, g.Player.P) > 1 {
			break
		}
		if RandInt(4) == 0 {
			m.EnterConfusion(g, g.Ev)
			g.Printf("%s appears confused.", m.Kind.Definite(true))
		}
	case FireShield:
		dir := Dir(m.P, g.Player.P)
		burnpos := To(g.Player.P, dir)
		if RandInt(4) == 0 {
			g.Print("Sparks emerge out of the shield.")
			g.Burn(burnpos, g.Ev)
		}
	}
}
