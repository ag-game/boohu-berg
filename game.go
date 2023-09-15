package main

import (
	"container/heap"
	"fmt"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
)

var Version string = "v0.14-dev"

type game struct {
	Dungeon             *dungeon
	Player              *player
	Monsters            []*monster
	MonstersPosCache    []int // monster (dungeon index + 1) / no monster (0)
	Bands               []monsterBand
	BandData            []monsterBandData
	Events              *eventQueue
	Ev                  event
	EventIndex          int
	Depth               int
	ExploredLevels      int
	DepthPlayerTurn     int
	Turn                int
	Highlight           map[gruid.Point]bool // highlighted positions (e.g. targeted ray)
	Collectables        map[gruid.Point]collectable
	CollectableScore    int
	LastConsumables     []consumable
	Equipables          map[gruid.Point]equipable
	Rods                map[gruid.Point]rod
	Stairs              map[gruid.Point]stair
	Clouds              map[gruid.Point]cloud
	Fungus              map[gruid.Point]vegetation
	Doors               map[gruid.Point]bool
	TemporalWalls       map[gruid.Point]bool
	MagicalStones       map[gruid.Point]stone
	GeneratedUniques    map[monsterBand]int
	GeneratedEquipables map[equipable]bool
	GeneratedRods       map[rod]bool
	GenPlan             [MaxDepth + 1]genFlavour
	FoundEquipables     map[equipable]bool
	Simellas            map[gruid.Point]int
	WrongWall           map[gruid.Point]bool
	WrongFoliage        map[gruid.Point]bool
	WrongDoor           map[gruid.Point]bool
	ExclusionsMap       map[gruid.Point]bool
	Noise               map[gruid.Point]bool
	DreamingMonster     map[gruid.Point]bool
	Resting             bool
	RestingTurns        int
	Autoexploring       bool
	DijkstraMapRebuild  bool
	Targeting           gruid.Point
	PR                  *paths.PathRange
	PRauto              *paths.PathRange
	AutoTarget          gruid.Point
	AutoDir             direction
	AutoHalt            bool
	AutoNext            bool
	DrawBuffer          []UICell
	drawBackBuffer      []UICell
	DrawLog             []drawFrame
	Log                 []logEntry
	LogIndex            int
	LogNextTick         int
	InfoEntry           string
	Stats               stats
	Boredom             int
	Quit                bool
	Wizard              bool
	WizardMap           bool
	Version             string
	Opts                startOpts
	ui                  *gameui
}

type startOpts struct {
	Alternate     monsterKind
	StoneLevel    int
	SpecialBands  map[int][]monsterBandData
	UnstableLevel int
}

func (g *game) FreeCell() gruid.Point {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if c.T != FreeCell {
			continue
		}
		if g.Player != nil && g.Player.P == p {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		return p
	}
}

func (g *game) FreeCellForPlayer() gruid.Point {
	center := gruid.Point{DungeonWidth / 2, DungeonHeight / 2}
	bestpos := g.FreeCell()
	for i := 0; i < 2; i++ {
		p := g.FreeCell()
		if Distance(p, center) > Distance(bestpos, center) {
			bestpos = p
		}
	}
	return bestpos
}

func (g *game) FreeCellForStair(dist int) gruid.Point {
	iters := 0
	bestpos := g.Player.P
	for {
		p := g.FreeCellForStatic()
		adjust := 0
		for i := 0; i < 4; i++ {
			adjust += RandInt(dist)
		}
		adjust /= 4
		if Distance(p, g.Player.P) <= 6+adjust {
			continue
		}
		iters++
		if Distance(p, g.Player.P) > Distance(bestpos, g.Player.P) {
			bestpos = p
		}
		if iters == 2 {
			return bestpos
		}
	}
}

func (g *game) FreeCellForStatic() gruid.Point {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForStatic")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if c.T != FreeCell {
			continue
		}
		if g.Player != nil && g.Player.P == p {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		if g.Doors[p] {
			continue
		}
		if g.Simellas[p] > 0 {
			continue
		}
		if _, ok := g.Collectables[p]; ok {
			continue
		}
		if _, ok := g.Stairs[p]; ok {
			continue
		}
		if _, ok := g.Rods[p]; ok {
			continue
		}
		if _, ok := g.Equipables[p]; ok {
			continue
		}
		if _, ok := g.MagicalStones[p]; ok {
			continue
		}
		return p
	}
}

func (g *game) FreeCellForMonster() gruid.Point {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForMonster")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if c.T != FreeCell {
			continue
		}
		if g.Player != nil && Distance(g.Player.P, p) < 8 {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		return p
	}
}

func (g *game) FreeCellForBandMonster(p gruid.Point) gruid.Point {
	count := 0
	for {
		count++
		if count > 1000 {
			return g.FreeCellForMonster()
		}
		neighbors := g.Dungeon.FreeNeighbors(p)
		r := RandInt(len(neighbors))
		p = neighbors[r]
		if g.Player != nil && Distance(g.Player.P, p) < 8 {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		return p
	}
}

func (g *game) FreeForStairs() gruid.Point {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeForStairs")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if c.T != FreeCell {
			continue
		}
		_, ok := g.Collectables[p]
		if ok {
			continue
		}
		return p
	}
}

const MaxDepth = 11
const WinDepth = 8

const (
	DungeonHeight = 21
	DungeonWidth  = 79
	DungeonNCells = DungeonWidth * DungeonHeight
)

func (g *game) GenDungeon() {
	g.Fungus = make(map[gruid.Point]vegetation)
	for {
		dg := GenRuinsMap
		switch RandInt(7) {
		//switch 4 {
		case 0:
			dg = GenCaveMap
		case 1:
			dg = GenRoomMap
		case 2:
			dg = GenCellularAutomataCaveMap
		case 3:
			dg = GenCaveMapTree
		case 4:
			dg = GenBSPMap
		}
		if g.Depth > 1 && dg.String() == g.Stats.DLayout[g.Depth-1] && RandInt(4) > 0 {
			// avoid too often the same layout in a row
			continue
		}
		dg.Use(g)
		break
	}
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:        42,
		MP:        3,
		Simellas:  0,
		Aptitudes: map[aptitude]bool{},
	}
	g.Player.Consumables = map[consumable]int{
		HealWoundsPotion: 1,
	}
	switch RandInt(7) {
	case 0:
		g.Player.Consumables[ExplosiveMagara] = 1
	case 1:
		g.Player.Consumables[NightMagara] = 1
	case 2:
		g.Player.Consumables[TeleportMagara] = 1
	case 3:
		g.Player.Consumables[SlowingMagara] = 1
	case 4:
		g.Player.Consumables[ConfuseMagara] = 1
	default:
		g.Player.Consumables[ConfusingDart] = 2
	}
	switch RandInt(12) {
	case 0, 1:
		g.Player.Consumables[TeleportationPotion] = 1
	case 2, 3:
		g.Player.Consumables[BerserkPotion] = 1
	case 4:
		g.Player.Consumables[SwiftnessPotion] = 1
	case 5:
		g.Player.Consumables[LignificationPotion] = 1
	case 6:
		g.Player.Consumables[WallPotion] = 1
	case 7:
		g.Player.Consumables[CBlinkPotion] = 1
	case 8:
		g.Player.Consumables[DigPotion] = 1
	case 9:
		g.Player.Consumables[SwapPotion] = 1
	case 10:
		g.Player.Consumables[ShadowsPotion] = 1
	case 11:
		g.Player.Consumables[AccuracyPotion] = 1
	}
	r := g.RandomRod()
	items := r.String()
	for c, n := range g.Player.Consumables {
		if n == 1 {
			items += ", " + c.String()
		} else {
			items += fmt.Sprintf(", %d %s", n, c.Plural())
		}
	}
	g.StoryPrintf("Started with %s", items)
	g.Player.Rods = map[rod]rodProps{r: {r.MaxCharge() - 1}}
	g.Player.Statuses = map[status]int{}
	g.Player.Expire = map[status]int{}

	// Testing
	//g.Player.Aptitudes[AptStealthyLOS] = true
	//g.Player.Aptitudes[AptStealthyMovement] = true
	//g.Player.Rods[RodSwapping] = rodProps{Charge: 3}
	//g.Player.Rods[RodFireball] = rodProps{Charge: 3}
	//g.Player.Rods[RodLightning] = rodProps{Charge: 3}
	//g.Player.Rods[RodLightningBolt] = rodProps{Charge: 3}
	//g.Player.Rods[RodShatter] = rodProps{Charge: 3}
	//g.Player.Rods[RodFog] = rodProps{Charge: 3}
	//g.Player.Rods[RodSleeping] = rodProps{Charge: 3}
	//g.Player.Consumables[BerserkPotion] = 5
	//g.Player.Consumables[MagicMappingPotion] = 1
	//g.Player.Consumables[ExplosiveMagara] = 5
	//g.Player.Consumables[NightMagara] = 5
	//g.Player.Consumables[SlowingMagara] = 5
	//g.Player.Consumables[ConfuseMagara] = 5
	//g.Player.Consumables[DigPotion] = 5
	//g.Player.Consumables[SwapPotion] = 5
	//g.Player.Consumables[DreamPotion] = 5
	//g.Player.Consumables[ShadowsPotion] = 5
	//g.Player.Consumables[TormentPotion] = 5
	//g.Player.Consumables[AccuracyPotion] = 5
	//g.Player.Weapon = ElecWhip
	//g.Player.Weapon = DancingRapier
	//g.Player.Weapon = Sabre
	//g.Player.Weapon = HarKarGauntlets
	//g.Player.Weapon = DefenderFlail
	//g.Player.Weapon = HopeSword
	//g.Player.Weapon = DragonSabre
	//g.Player.Weapon = FinalBlade
	//g.Player.Weapon = VampDagger
	//g.Player.Shield = EarthShield
	//g.Player.Shield = FireShield
	//g.Player.Shield = BashingShield
	//g.Player.Armour = TurtlePlates
	//g.Player.Armour = HarmonistRobe
	//g.Player.Armour = CelmistRobe
	//g.Player.Armour = ShinyPlates
	//g.Player.Armour = SmokingScales
}

func (g *game) InitSpecialBands() {
	g.Opts.SpecialBands = map[int][]monsterBandData{}
	sb := MonsSpecialBands[RandInt(len(MonsSpecialBands))]
	depth := sb.minDepth + RandInt(sb.maxDepth-sb.minDepth+1)
	g.Opts.SpecialBands[depth] = sb.bands
	seb := MonsSpecialEndBands[RandInt(len(MonsSpecialEndBands))]
	if RandInt(4) == 0 {
		if RandInt(5) > 1 || depth == WinDepth {
			g.Opts.SpecialBands[WinDepth+1] = seb.bands
		} else {
			g.Opts.SpecialBands[WinDepth] = seb.bands
		}
	} else if RandInt(5) > 0 {
		if RandInt(3) > 0 {
			g.Opts.SpecialBands[MaxDepth] = seb.bands
		} else {
			g.Opts.SpecialBands[MaxDepth-1] = seb.bands
		}
	}
}

type genFlavour int

const (
	GenRod genFlavour = iota
	GenWeapon
	GenArmour
	GenWpArm
	GenExtraCollectables
)

func (g *game) InitFirstLevel() {
	g.Depth++ // start at 1
	g.InitPlayer()
	g.AutoTarget = InvalidPos
	g.Targeting = InvalidPos
	g.GeneratedRods = map[rod]bool{}
	g.GeneratedEquipables = map[equipable]bool{}
	g.FoundEquipables = map[equipable]bool{Robe: true, Dagger: true}
	g.GeneratedUniques = map[monsterBand]int{}
	g.Stats.KilledMons = map[monsterKind]int{}
	g.InitSpecialBands()
	if RandInt(4) > 0 {
		g.Opts.UnstableLevel = 1 + RandInt(MaxDepth)
	}
	if g.Opts.UnstableLevel >= 1 && g.Opts.UnstableLevel <= 3 {
		// it should happen less often in the first levels
		g.Opts.UnstableLevel += RandInt(MaxDepth - 2)
	}
	if RandInt(3) > 0 || RandInt(2) == 0 && g.Opts.UnstableLevel == 0 {
		g.Opts.StoneLevel = 1 + RandInt(MaxDepth)
	}
	if g.Opts.StoneLevel >= 1 && g.Opts.StoneLevel <= 3 {
		g.Opts.StoneLevel += RandInt(MaxDepth - 2)
	}
	if RandInt(3) == 0 {
		g.Opts.Alternate = MonsTinyHarpy
		if RandInt(10) == 0 {
			g.Opts.Alternate = MonsWorm
		}
	}
	g.Version = Version
	g.GenPlan = [MaxDepth + 1]genFlavour{
		1:  GenRod,
		2:  GenWeapon,
		3:  GenArmour,
		4:  GenRod,
		5:  GenExtraCollectables,
		6:  GenWpArm,
		7:  GenRod,
		8:  GenExtraCollectables,
		9:  GenWeapon,
		10: GenExtraCollectables,
		11: GenExtraCollectables,
	}
	permi := RandInt(7)
	switch permi {
	case 0, 1, 2, 3:
		g.GenPlan[permi+1], g.GenPlan[permi+2] = g.GenPlan[permi+2], g.GenPlan[permi+1]
	}
	if RandInt(4) == 0 {
		g.GenPlan[6], g.GenPlan[7] = g.GenPlan[7], g.GenPlan[6]
	}
	g.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	g.PRauto = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
}

func (g *game) InitLevel() {
	// Starting data
	if g.Depth == 0 {
		g.InitFirstLevel()
	}

	// Dungeon terrain
	g.GenDungeon()

	g.MonstersPosCache = make([]int, DungeonNCells)
	g.Player.P = g.FreeCellForPlayer()

	g.WrongWall = map[gruid.Point]bool{}
	g.WrongFoliage = map[gruid.Point]bool{}
	g.WrongDoor = map[gruid.Point]bool{}
	g.ExclusionsMap = map[gruid.Point]bool{}
	g.TemporalWalls = map[gruid.Point]bool{}
	g.DreamingMonster = map[gruid.Point]bool{}

	// Monsters
	g.BandData = MonsBands
	if bd, ok := g.Opts.SpecialBands[g.Depth]; ok {
		g.BandData = bd
	}
	g.GenMonsters()

	// Collectables
	g.Collectables = make(map[gruid.Point]collectable)
	g.GenCollectables()

	// Equipment
	g.Equipables = make(map[gruid.Point]equipable)
	g.Rods = map[gruid.Point]rod{}
	switch g.GenPlan[g.Depth] {
	case GenWeapon:
		g.GenWeapon()
	case GenArmour:
		g.GenArmour()
	case GenWpArm:
		g.GenWeapon()
		g.GenArmour()
	case GenRod:
		g.GenerateRod()
	case GenExtraCollectables:
		for i := 0; i < 2; i++ {
			g.GenCollectable()
			g.CollectableScore-- // these are extra
		}
	}
	if g.Depth == 1 {
		// extra collectable
		g.GenCollectable()
		g.CollectableScore--
	}

	// Aptitudes/Mutations
	if g.Depth == 2 || g.Depth == 5 {
		apt, ok := g.RandomApt()
		if ok {
			g.ApplyAptitude(apt)
		}
	}

	// Stairs
	g.Stairs = make(map[gruid.Point]stair)
	nstairs := 2
	if RandInt(3) == 0 {
		if RandInt(2) == 0 {
			nstairs++
		} else {
			nstairs--
		}
	}
	if g.Depth >= WinDepth {
		nstairs = 1
	} else if g.Depth == WinDepth-1 && nstairs > 2 {
		nstairs = 2
	}
	for i := 0; i < nstairs; i++ {
		var p gruid.Point
		if g.Depth >= WinDepth && g.Depth != MaxDepth-1 {
			p = g.FreeCellForStair(60)
			g.Stairs[p] = WinStair
		}
		if g.Depth < MaxDepth {
			if g.Depth > 5 {
				p = g.FreeCellForStair(50)
			} else {
				p = g.FreeCellForStair(0)
			}
			g.Stairs[p] = NormalStair
		}
	}

	// Magical Stones
	g.MagicalStones = map[gruid.Point]stone{}
	nstones := 1
	switch RandInt(8) {
	case 0:
		nstones = 0
	case 1, 2, 3:
		nstones = 2
	case 4, 5, 6:
		nstones = 3
	}
	ustone := stone(0)
	if g.Depth == g.Opts.StoneLevel {
		ustone = stone(1 + RandInt(NumStones-1))
		nstones = 10 + RandInt(3)
		if RandInt(4) == 0 {
			g.Opts.StoneLevel = g.Opts.StoneLevel + RandInt(MaxDepth-g.Opts.StoneLevel) + 1
		}
	}
	for i := 0; i < nstones; i++ {
		p := g.FreeCellForStatic()
		var st stone
		if ustone != stone(0) {
			st = ustone
		} else {
			st = stone(1 + RandInt(NumStones-1))
		}
		g.MagicalStones[p] = st
	}

	// Simellas
	g.Simellas = make(map[gruid.Point]int)
	for i := 0; i < 5; i++ {
		p := g.FreeCellForStatic()
		const rounds = 5
		for j := 0; j < rounds; j++ {
			g.Simellas[p] += 1 + RandInt(g.Depth+g.Depth*g.Depth/6)
		}
		g.Simellas[p] /= rounds
		if g.Simellas[p] == 0 {
			g.Simellas[p] = 1
		}
	}

	// initialize LOS
	if g.Depth == 1 {
		g.Print("You're in Hareka's Underground searching for medicinal simellas. Good luck!")
		g.PrintStyled("► Type ? for help on keys or use the mouse and [buttons].", logSpecial)
	}
	if g.Depth == WinDepth {
		g.PrintStyled("You feel magic in the air. A first way out is close!", logSpecial)
	} else if g.Depth == MaxDepth {
		g.PrintStyled("If rumors are true, you have reached the bottom!", logSpecial)
	}
	g.ComputeLOS()
	g.MakeMonstersAware()

	// Frundis is somewhere in the level
	if g.FrundisInLevel() {
		g.PrintStyled("You hear some faint music… ♫ larilon, larila ♫ ♪", logSpecial)
	}

	// recharge rods
	if g.Depth > 1 {
		g.RechargeRods()
	}

	// clouds
	g.Clouds = map[gruid.Point]cloud{}

	// Events
	if g.Depth == 1 {
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		g.PushEvent(&simpleEvent{ERank: 0, EAction: PlayerTurn})
	} else {
		g.CleanEvents()
	}
	for i := range g.Monsters {
		g.PushEvent(&monsterEvent{ERank: g.Turn + RandInt(10), EAction: MonsterTurn, NMons: i})
	}
	if g.Depth == g.Opts.UnstableLevel {
		g.PrintStyled("You sense magic instability on this level.", logSpecial)
		for i := 0; i < 15; i++ {
			g.PushEvent(&cloudEvent{ERank: g.Turn + 100 + RandInt(900), EAction: ObstructionProgression})
		}
		if RandInt(4) == 0 {
			g.Opts.UnstableLevel = g.Opts.UnstableLevel + RandInt(MaxDepth-g.Opts.UnstableLevel) + 1
		}
	}
}

func (g *game) CleanEvents() {
	evq := &eventQueue{}
	for g.Events.Len() > 0 {
		iev := g.PopIEvent()
		switch iev.Event.(type) {
		case *monsterEvent:
		case *cloudEvent:
		default:
			heap.Push(evq, iev)
		}
	}
	g.Events = evq
}

func (g *game) StairsSlice() []gruid.Point {
	stairs := []gruid.Point{}
	for stairPos := range g.Stairs {
		if g.Dungeon.Cell(stairPos).Explored {
			stairs = append(stairs, stairPos)
		}
	}
	return stairs
}

func (g *game) GenCollectable() {
	rounds := 100
	if len(g.LastConsumables) > 3 {
		g.LastConsumables = g.LastConsumables[1:]
	}
	for {
	loopcons:
		for c, data := range ConsumablesCollectData {
			r := RandInt(data.rarity * rounds)
			if r != 0 {
				continue
			}

			// avoid too many of the same
			for _, co := range g.LastConsumables {
				if co == c && RandInt(4) > 0 {
					continue loopcons
				}
			}
			g.LastConsumables = append(g.LastConsumables, c)
			g.CollectableScore++
			p := g.FreeCellForStatic()
			g.Collectables[p] = collectable{Consumable: c, Quantity: data.quantity}
			return
		}
	}

}

func (g *game) GenCollectables() {
	score := g.CollectableScore - 2*(g.Depth-1)
	n := 2
	if score >= 0 && RandInt(4) == 0 {
		n--
	}
	if score <= 0 && RandInt(4) == 0 {
		n++
	}
	if score > 0 && n >= 2 {
		n--
	}
	if score < 0 && n <= -2 {
		n++
	}
	for i := 0; i < n; i++ {
		g.GenCollectable()
	}
}

func (g *game) GenShield() {
	ars := [4]shield{ConfusingShield, BashingShield, EarthShield, FireShield}
	for {
		i := RandInt(len(ars))
		if g.GeneratedEquipables[ars[i]] {
			// do not generate duplicates
			continue
		}
		p := g.FreeCellForStatic()
		g.Equipables[p] = ars[i]
		g.GeneratedEquipables[ars[i]] = true
		break
	}
}

func (g *game) GenArmour() {
	ars := [6]armour{SmokingScales, ShinyPlates, TurtlePlates, SpeedRobe, CelmistRobe, HarmonistRobe}
	for {
		i := RandInt(len(ars))
		if g.GeneratedEquipables[ars[i]] {
			// do not generate duplicates
			continue
		}
		p := g.FreeCellForStatic()
		g.Equipables[p] = ars[i]
		g.GeneratedEquipables[ars[i]] = true
		break
	}
}

func (g *game) GenWeapon() {
	wps := [WeaponNum - 1]weapon{Axe, BattleAxe, Spear, Halberd, AssassinSabre, DancingRapier, HopeSword, Frundis, ElecWhip, HarKarGauntlets, VampDagger, DragonSabre, FinalBlade, DefenderFlail}
	onehanded := false
	for {
		i := RandInt(len(wps))
		if g.GeneratedEquipables[wps[i]] {
			// do not generate duplicates
			continue
		}
		p := g.FreeCellForStatic()
		g.Equipables[p] = wps[i]
		if !wps[i].TwoHanded() {
			onehanded = true
		}
		g.GeneratedEquipables[wps[i]] = true
		break
	}
	if onehanded {
		g.GenShield()
	}
}

func (g *game) FrundisInLevel() bool {
	for _, eq := range g.Equipables {
		if wp, ok := eq.(weapon); ok && wp == Frundis {
			return true
		}
	}
	return false
}

func (g *game) Descend() bool {
	g.LevelStats()
	if strt, ok := g.Stairs[g.Player.P]; ok && strt == WinStair {
		g.StoryPrint("Escaped!")
		g.ExploredLevels = g.Depth
		g.Depth = -1
		return true
	}
	g.Print("You descend deeper in the dungeon.")
	g.StoryPrint("Descended deeper in the dungeon.")
	g.Depth++
	g.DepthPlayerTurn = 0
	g.Boredom = 0
	g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: PlayerTurn})
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) WizardMode() {
	g.Wizard = true
	g.Player.Consumables[DescentPotion] = 15
	g.PrintStyled("You are now in wizard mode and cannot obtain winner status.", logSpecial)
}

func (g *game) ApplyRest() {
	g.Player.HP = g.Player.HPMax()
	g.Player.MP = g.Player.MPMax()
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			continue
		}
		mons.HP = mons.HPmax
	}
	adjust := 0
	if g.Player.Armour == HarmonistRobe {
		// the harmonist robe mitigates the sound of your snorts
		adjust = 100
	}
	if g.DepthPlayerTurn < 100+adjust && RandInt(5) > 2 || g.DepthPlayerTurn >= 100+adjust && g.DepthPlayerTurn < 250+adjust && RandInt(2) == 0 ||
		g.DepthPlayerTurn >= 250+adjust && RandInt(3) > 0 {
		rmons := []int{}
		for i, mons := range g.Monsters {
			if mons.Exists() && mons.State == Resting {
				rmons = append(rmons, i)
			}
		}
		if len(rmons) > 0 {
			g.Monsters[rmons[RandInt(len(rmons))]].NaturalAwake(g)
		}
	}
	g.Stats.Rest++
	g.PrintStyled("You feel fresh again. Some monsters might have awoken.", logStatusEnd)
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		const enoughRestTurns = 15
		mons := g.MonsterInLOS()
		sr := g.StatusRest()
		if mons == nil && (sr || g.NeedsRegenRest() && g.RestingTurns >= 0) && g.RestingTurns < enoughRestTurns {
			g.WaitTurn(ev)
			if !sr && g.RestingTurns >= 0 {
				g.RestingTurns++
			}
			return true
		}
		if g.RestingTurns >= enoughRestTurns {
			g.ApplyRest()
		} else if mons != nil {
			g.Stats.RestInterrupt++
			g.Print("You could not sleep.")
		}
		g.Resting = false
	} else if g.Autoexploring {
		if g.ui.ExploreStep() {
			g.AutoHalt = true
			g.Print("Stopping, then.")
		}
		switch {
		case g.AutoHalt:
			// stop exploring
		default:
			var n *gruid.Point
			var finished bool
			if g.DijkstraMapRebuild {
				if g.AllExplored() {
					g.Print("You finished exploring.")
					break
				}
				sources := g.AutoexploreSources()
				g.BuildAutoexploreMap(sources)
			}
			n, finished = g.NextAuto()
			if finished {
				n = nil
			}
			if finished && g.AllExplored() {
				g.Print("You finished exploring.")
			} else if n == nil {
				g.Print("You could not safely reach some places.")
			}
			if n != nil {
				err := g.MovePlayer(*n, ev)
				if err != nil {
					g.Print(err.Error())
					break
				}
				return true
			}
		}
		g.Autoexploring = false
	} else if valid(g.AutoTarget) {
		if !g.ui.ExploreStep() && g.MoveToTarget(ev) {
			return true
		} else {
			g.AutoTarget = InvalidPos
		}
	} else if g.AutoDir != NoDir {
		if !g.ui.ExploreStep() && g.AutoToDir(ev) {
			return true
		} else {
			g.AutoDir = NoDir
		}
	}
	return false
}

func (g *game) EventLoop() {
loop:
	for {
		if g.Player.HP <= 0 {
			if g.Wizard {
				g.Player.HP = g.Player.HPMax()
			} else {
				g.LevelStats()
				err := g.RemoveSaveFile()
				if err != nil {
					g.PrintfStyled("Error removing save file: %v", logError, err.Error())
				}
				g.ui.Death()
				break loop
			}
		}
		if g.Events.Len() == 0 {
			break loop
		}
		ev := g.PopIEvent().Event
		g.Turn = ev.Rank()
		g.Ev = ev
		ev.Action(g)
		if g.AutoNext {
			continue loop
		}
		if g.Quit {
			break loop
		}
	}
}
