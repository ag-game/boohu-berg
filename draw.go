package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"codeberg.org/anaseto/gruid"
)

var (
	UIWidth                = 100
	UIHeight               = 26
	DisableAnimations bool = false
)

type uicolor int

const (
	Color256Base03  uicolor = 234
	Color256Base02  uicolor = 235
	Color256Base01  uicolor = 240
	Color256Base00  uicolor = 241 // for dark on light background
	Color256Base0   uicolor = 244
	Color256Base1   uicolor = 245
	Color256Base2   uicolor = 254
	Color256Base3   uicolor = 230
	Color256Yellow  uicolor = 136
	Color256Orange  uicolor = 166
	Color256Red     uicolor = 160
	Color256Magenta uicolor = 125
	Color256Violet  uicolor = 61
	Color256Blue    uicolor = 33
	Color256Cyan    uicolor = 37
	Color256Green   uicolor = 64

	Color16Base03  uicolor = 8
	Color16Base02  uicolor = 0
	Color16Base01  uicolor = 10
	Color16Base00  uicolor = 11
	Color16Base0   uicolor = 12
	Color16Base1   uicolor = 14
	Color16Base2   uicolor = 7
	Color16Base3   uicolor = 15
	Color16Yellow  uicolor = 3
	Color16Orange  uicolor = 9
	Color16Red     uicolor = 1
	Color16Magenta uicolor = 5
	Color16Violet  uicolor = 13
	Color16Blue    uicolor = 4
	Color16Cyan    uicolor = 6
	Color16Green   uicolor = 2
)

// uicolors: http://ethanschoonover.com/solarized
var (
	ColorBase03  uicolor = Color256Base03
	ColorBase02  uicolor = Color256Base02
	ColorBase01  uicolor = Color256Base01
	ColorBase00  uicolor = Color256Base00 // for dark on light background
	ColorBase0   uicolor = Color256Base0
	ColorBase1   uicolor = Color256Base1
	ColorBase2   uicolor = Color256Base2
	ColorBase3   uicolor = Color256Base3
	ColorYellow  uicolor = Color256Yellow
	ColorOrange  uicolor = Color256Orange
	ColorRed     uicolor = Color256Red
	ColorMagenta uicolor = Color256Magenta
	ColorViolet  uicolor = Color256Violet
	ColorBlue    uicolor = Color256Blue
	ColorCyan    uicolor = Color256Cyan
	ColorGreen   uicolor = Color256Green
)

func (ui *gameui) Map256ColorTo16(c uicolor) uicolor {
	switch c {
	case Color256Base03:
		return Color16Base03
	case Color256Base02:
		return Color16Base02
	case Color256Base01:
		return Color16Base01
	case Color256Base00:
		return Color16Base00
	case Color256Base0:
		return Color16Base0
	case Color256Base1:
		return Color16Base1
	case Color256Base2:
		return Color16Base2
	case Color256Base3:
		return Color16Base3
	case Color256Yellow:
		return Color16Yellow
	case Color256Orange:
		return Color16Orange
	case Color256Red:
		return Color16Red
	case Color256Magenta:
		return Color16Magenta
	case Color256Violet:
		return Color16Violet
	case Color256Blue:
		return Color16Blue
	case Color256Cyan:
		return Color16Cyan
	case Color256Green:
		return Color16Green
	default:
		return c
	}
}

func (ui *gameui) Map16ColorTo256(c uicolor) uicolor {
	switch c {
	case Color16Base03:
		return Color256Base03
	case Color16Base02:
		return Color256Base02
	case Color16Base01:
		return Color256Base01
	case Color16Base00:
		return Color256Base00
	case Color16Base0:
		return Color256Base0
	case Color16Base1:
		return Color256Base1
	case Color16Base2:
		return Color256Base2
	case Color16Base3:
		return Color256Base3
	case Color16Yellow:
		return Color256Yellow
	case Color16Orange:
		return Color256Orange
	case Color16Red:
		return Color256Red
	case Color16Magenta:
		return Color256Magenta
	case Color16Violet:
		return Color256Violet
	case Color16Blue:
		return Color256Blue
	case Color16Cyan:
		return Color256Cyan
	case Color16Green:
		return Color256Green
	default:
		return c
	}
}

var (
	ColorBg,
	ColorBgBorder,
	ColorBgDark,
	ColorBgLOS,
	ColorFg,
	ColorFgAnimationHit,
	ColorFgCollectable,
	ColorFgConfusedMonster,
	ColorFgLignifiedMonster,
	ColorFgSlowedMonster,
	ColorFgDark,
	ColorFgExcluded,
	ColorFgExplosionEnd,
	ColorFgExplosionStart,
	ColorFgExplosionWallEnd,
	ColorFgExplosionWallStart,
	ColorFgHPcritical,
	ColorFgHPok,
	ColorFgHPwounded,
	ColorFgLOS,
	ColorFgMPcritical,
	ColorFgMPok,
	ColorFgMPpartial,
	ColorFgMagicPlace,
	ColorFgMonster,
	ColorFgPlace,
	ColorFgPlayer,
	ColorFgProjectile,
	ColorFgSimellas,
	ColorFgSleepingMonster,
	ColorFgStatusBad,
	ColorFgStatusGood,
	ColorFgStatusExpire,
	ColorFgStatusOther,
	ColorFgTargetMode,
	ColorFgWanderingMonster uicolor
)

func LinkColors() {
	ColorBg = ColorBase03
	ColorBgBorder = ColorBase02
	ColorBgDark = ColorBase03
	ColorBgLOS = ColorBase3
	ColorFg = ColorBase0
	ColorFgDark = ColorBase01
	ColorFgLOS = ColorBase0
	ColorFgAnimationHit = ColorMagenta
	ColorFgCollectable = ColorYellow
	ColorFgConfusedMonster = ColorGreen
	ColorFgLignifiedMonster = ColorYellow
	ColorFgSlowedMonster = ColorCyan
	ColorFgExcluded = ColorRed
	ColorFgExplosionEnd = ColorOrange
	ColorFgExplosionStart = ColorYellow
	ColorFgExplosionWallEnd = ColorMagenta
	ColorFgExplosionWallStart = ColorViolet
	ColorFgHPcritical = ColorRed
	ColorFgHPok = ColorGreen
	ColorFgHPwounded = ColorYellow
	ColorFgMPcritical = ColorMagenta
	ColorFgMPok = ColorBlue
	ColorFgMPpartial = ColorViolet
	ColorFgMagicPlace = ColorCyan
	ColorFgMonster = ColorRed
	ColorFgPlace = ColorMagenta
	ColorFgPlayer = ColorBlue
	ColorFgProjectile = ColorBlue
	ColorFgSimellas = ColorYellow
	ColorFgSleepingMonster = ColorViolet
	ColorFgStatusBad = ColorRed
	ColorFgStatusGood = ColorBlue
	ColorFgStatusExpire = ColorViolet
	ColorFgStatusOther = ColorYellow
	ColorFgTargetMode = ColorCyan
	ColorFgWanderingMonster = ColorOrange
}

func ApplyDarkLOS() {
	ColorBg = ColorBase03
	ColorBgBorder = ColorBase02
	ColorBgDark = ColorBase03
	ColorBgLOS = ColorBase02
	ColorFgDark = ColorBase01
	ColorFg = ColorBase0
	if Only8Colors {
		ColorFgLOS = ColorGreen
	} else {
		ColorFgLOS = ColorBase0
	}
}

func ApplyLightLOS() {
	if Only8Colors {
		ApplyDarkLOS()
		ColorBgLOS = ColorBase2
		ColorFgLOS = ColorBase00
	} else {
		ColorBg = ColorBase3
		ColorBgBorder = ColorBase2
		ColorBgDark = ColorBase3
		ColorBgLOS = ColorBase2
		ColorFgDark = ColorBase1
		ColorFgLOS = ColorBase00
		ColorFg = ColorBase00
	}
}

func SolarizedPalette() {
	ColorBase03 = Color16Base03
	ColorBase02 = Color16Base02
	ColorBase01 = Color16Base01
	ColorBase00 = Color16Base00
	ColorBase0 = Color16Base0
	ColorBase1 = Color16Base1
	ColorBase2 = Color16Base2
	ColorBase3 = Color16Base3
	ColorYellow = Color16Yellow
	ColorOrange = Color16Orange
	ColorRed = Color16Red
	ColorMagenta = Color16Magenta
	ColorViolet = Color16Violet
	ColorBlue = Color16Blue
	ColorCyan = Color16Cyan
	ColorGreen = Color16Green
}

const (
	Black uicolor = iota
	Maroon
	Green
	Olive
	Navy
	Purple
	Teal
	Silver
)

func Map16ColorTo8Color(c uicolor) uicolor {
	switch c {
	case Color16Base03:
		return Black
	case Color16Base02:
		return Black
	case Color16Base01:
		return Silver
	case Color16Base00:
		return Black
	case Color16Base0:
		return Silver
	case Color16Base1:
		return Silver
	case Color16Base2:
		return Silver
	case Color16Base3:
		return Silver
	case Color16Yellow:
		return Olive
	case Color16Orange:
		return Purple
	case Color16Red:
		return Maroon
	case Color16Magenta:
		return Purple
	case Color16Violet:
		return Teal
	case Color16Blue:
		return Navy
	case Color16Cyan:
		return Teal
	case Color16Green:
		return Green
	default:
		return c
	}
}

var Only8Colors bool

func Simple8ColorPalette() {
	Only8Colors = true
}

type drawFrame struct {
	Draws []cellDraw
	Time  time.Time
}

type cellDraw struct {
	Cell UICell
	X    int
	Y    int
}

func (ui *gameui) SetCell(x, y int, r rune, fg, bg uicolor) {
	ui.SetGenCell(x, y, r, fg, bg, false)
}

func (ui *gameui) SetGenCell(x, y int, r rune, fg, bg uicolor, inmap bool) {
	i := ui.GetIndex(x, y)
	if i >= UIHeight*UIWidth {
		return
	}
	c := UICell{R: r, Fg: fg, Bg: bg, InMap: inmap}
	ui.g.DrawBuffer[i] = c
}

func (ui *gameui) SetMapCell(x, y int, r rune, fg, bg uicolor) {
	ui.SetGenCell(x, y, r, fg, bg, true)
}

func (ui *gameui) DrawLogFrame() {
	if len(ui.g.drawBackBuffer) != len(ui.g.DrawBuffer) {
		ui.g.drawBackBuffer = make([]UICell, len(ui.g.DrawBuffer))
	}
	ui.g.DrawLog = append(ui.g.DrawLog, drawFrame{Time: time.Now()})
	for i := 0; i < len(ui.g.DrawBuffer); i++ {
		if ui.g.DrawBuffer[i] == ui.g.drawBackBuffer[i] {
			continue
		}
		c := ui.g.DrawBuffer[i]
		x, y := ui.GetPos(i)
		cdraw := cellDraw{Cell: c, X: x, Y: y}
		last := len(ui.g.DrawLog) - 1
		ui.g.DrawLog[last].Draws = append(ui.g.DrawLog[last].Draws, cdraw)
		ui.g.drawBackBuffer[i] = c
	}
}

func (ui *gameui) DrawWelcomeCommon() int {
	ui.DrawBufferInit()
	ui.Clear()
	col := 10
	line := 5
	rcol := col + 20
	ColorText := ColorFgHPok
	ui.DrawDark(fmt.Sprintf("       Boohu %s", Version), col, line-2, ColorText, false)
	ui.DrawDark("────│\\/\\/\\/\\/\\/\\/\\/│────", col, line, ColorText, false)
	line++
	ui.DrawDark("##", col, line, ColorFgDark, true)
	ui.DrawLOS("#", col+2, line, ColorFgLOS, true)
	ui.DrawLOS("#", col+3, line, ColorFgLOS, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark("####", rcol, line, ColorFgDark, true)
	line++
	ui.DrawDark("#.", col, line, ColorFgDark, true)
	ui.DrawLOS(".", col+2, line, ColorFgLOS, true)
	ui.DrawLOS(".", col+3, line, ColorFgLOS, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark(".", rcol, line, ColorFgDark, true)
	ui.DrawDark("♣", rcol+1, line, ColorFgSimellas, true)
	ui.DrawDark(".#", rcol+2, line, ColorFgDark, true)
	line++
	ui.DrawDark("##", col, line, ColorFgDark, true)
	ui.DrawLOS("!", col+2, line, ColorFgCollectable, true)
	ui.DrawLOS(".", col+3, line, ColorFgLOS, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark("│  BREAK       │", col+4, line, ColorText, false)
	ui.DrawDark(".###", rcol, line, ColorFgDark, true)
	line++
	ui.DrawDark(" #", col, line, ColorFgDark, true)
	ui.DrawLOS("g", col+2, line, ColorFgMonster, true)
	ui.DrawLOS("G", col+3, line, ColorFgMonster, true)
	ui.DrawDark("│  OUT OF      │", col+4, line, ColorText, false)
	ui.DrawDark("##  ", rcol, line, ColorFgDark, true)
	line++
	ui.DrawLOS("#", col, line, ColorFgLOS, true)
	ui.DrawLOS("#", col+1, line, ColorFgLOS, true)
	ui.DrawLOS("D", col+2, line, ColorFgMonster, true)
	ui.DrawLOS("g", col+3, line, ColorFgMonster, true)
	ui.DrawDark("│  HAREKA'S    │", col+4, line, ColorText, false)
	ui.DrawDark(".## ", rcol, line, ColorFgDark, true)
	line++
	ui.DrawLOS("#", col, line, ColorFgLOS, true)
	ui.DrawLOS("@", col+1, line, ColorFgPlayer, true)
	ui.DrawLOS("#", col+2, line, ColorFgLOS, true)
	ui.DrawDark("#", col+3, line, ColorFgDark, true)
	ui.DrawDark("│  UNDERGROUND │", col+4, line, ColorText, false)
	ui.DrawDark("\".##", rcol, line, ColorFgDark, true)
	line++
	ui.DrawLOS("#", col, line, ColorFgLOS, true)
	ui.DrawLOS(".", col+1, line, ColorFgLOS, true)
	ui.DrawLOS("#", col+2, line, ColorFgLOS, true)
	ui.DrawDark("#", col+3, line, ColorFgDark, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark("#.", rcol, line, ColorFgDark, true)
	ui.DrawDark(">", rcol+2, line, ColorFgPlace, true)
	ui.DrawDark("#", rcol+3, line, ColorFgDark, true)
	line++
	ui.DrawLOS("#", col, line, ColorFgLOS, true)
	ui.DrawLOS("[", col+1, line, ColorFgCollectable, true)
	ui.DrawLOS(".", col+2, line, ColorFgLOS, true)
	ui.DrawDark("##", col+3, line, ColorFgDark, true)
	ui.DrawDark("│              │", col+4, line, ColorFgHPok, false)
	ui.DrawDark("\"\"##", rcol, line, ColorFgDark, true)
	line++
	ui.DrawDark("────│/\\/\\/\\/\\/\\/\\/\\│────", col, line, ColorText, false)
	line++
	line++
	if runtime.GOARCH == "wasm" {
		ui.DrawDark("- (P)lay", col-3, line, ColorFg, false)
		ui.DrawDark("- (W)atch replay", col-3, line+1, ColorFg, false)
	} else {
		ui.DrawDark("───Press any key to continue───", col-3, line, ColorFg, false)
	}
	ui.Flush()
	return line
}

func (ui *gameui) DrawWelcome() {
	ui.DrawWelcomeCommon()
	ui.PressAnyKey()
}

func (ui *gameui) RestartDrawBuffers() {
	g := ui.g
	g.DrawBuffer = nil
	g.drawBackBuffer = nil
	ui.DrawBufferInit()
}

func (ui *gameui) DrawColored(text string, x, y int, fg, bg uicolor) {
	col := 0
	for _, r := range text {
		ui.SetCell(x+col, y, r, fg, bg)
		col++
	}
}

func (ui *gameui) DrawDark(text string, x, y int, fg uicolor, inmap bool) {
	col := 0
	for _, r := range text {
		if inmap {
			ui.SetMapCell(x+col, y, r, fg, ColorBgDark)
		} else {
			ui.SetCell(x+col, y, r, fg, ColorBgDark)
		}
		col++
	}
}

func (ui *gameui) DrawLOS(text string, x, y int, fg uicolor, inmap bool) {
	col := 0
	for _, r := range text {
		if inmap {
			ui.SetMapCell(x+col, y, r, fg, ColorBgLOS)
		} else {
			ui.SetCell(x+col, y, r, fg, ColorBgLOS)
		}
		col++
	}
}

func (ui *gameui) DrawKeysDescription(title string, actions []string) {
	ui.DrawDungeonView(NoFlushMode)

	if CustomKeys {
		ui.DrawStyledTextLine(fmt.Sprintf(" Default %s ", title), 0, HeaderLine)
	} else {
		ui.DrawStyledTextLine(fmt.Sprintf(" %s ", title), 0, HeaderLine)
	}
	for i := 0; i < len(actions)-1; i += 2 {
		bg := ui.ListItemBG(i / 2)
		ui.ClearLineWithColor(i/2+1, bg)
		ui.DrawColoredTextOnBG(fmt.Sprintf(" %-36s %s", actions[i], actions[i+1]), 0, i/2+1, ColorFg, bg)
	}
	lines := 1 + len(actions)/2
	ui.DrawTextLine(" press (x) to continue ", lines)
	ui.Flush()

	ui.WaitForContinue(lines)
}

func (ui *gameui) KeysHelp() {
	ui.DrawKeysDescription("Commands", []string{
		"Movement", "h/j/k/l/y/u/b/n or numpad or mouse left",
		"Wait a turn", "“.” or 5 or mouse left on @",
		"Rest (until status free or regen)", "r",
		"Descend stairs", "> or D",
		"Go to nearest stairs", "G",
		"Autoexplore", "o",
		"Examine", "x or mouse left",
		"Equip/Get weapon/armour/...", "e or g",
		"Quaff/Drink potion", "q or d",
		"Throw/Fire item", "t or f",
		"Evoke/Zap rod", "v or z",
		"Inventory summary", `i`,
		"View Character and Quest Information", `% or C`,
		"View previous messages", "m",
		"Write game statistics to file", "#",
		"Save and Quit", "S",
		"Quit without saving", "Q",
		"Change settings and key bindings", "=",
	})
}

func (ui *gameui) ExamineHelp() {
	ui.DrawKeysDescription("Examine/Travel/Targeting Commands", []string{
		"Move cursor", "h/j/k/l/y/u/b/n or numpad or mouse left",
		"Cycle through monsters", "+",
		"Cycle through stairs", ">",
		"Cycle through objects", "o",
		"Go to/select target", "“.” or enter or mouse left",
		"View target description", "v or d or mouse right",
		"Toggle exclude area from auto-travel", "e or mouse middle",
	})
}

const TextWidth = DungeonWidth - 2

func (ui *gameui) CharacterInfo() {
	g := ui.g
	ui.DrawDungeonView(NoFlushMode)

	b := bytes.Buffer{}
	b.WriteString(formatText("Every year, the elders send someone to collect medicinal simella plants in the Underground.  This year, the honor fell upon you, and so here you are.  According to the elders, deep in the Underground, a magical monolith will lead you back to your village.", TextWidth))
	b.WriteString("\n\n")
	b.WriteString(formatText(
		fmt.Sprintf("You are wielding %s. %s", Indefinite(g.Player.Weapon.String(), false), g.Player.Weapon.Desc()), TextWidth))
	b.WriteString("\n\n")
	b.WriteString(formatText(fmt.Sprintf("You are wearing %s. %s", g.Player.Armour.StringIndefinite(), g.Player.Armour.Desc()), TextWidth))
	b.WriteString("\n\n")
	if g.Player.Shield != NoShield {
		b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Shield, g.Player.Shield.Desc()), TextWidth))
		b.WriteString("\n\n")
	}
	b.WriteString(ui.AptitudesText())

	desc := b.String()
	lines := strings.Count(desc, "\n")
	for i := 0; i <= lines+2; i++ {
		if i >= DungeonWidth {
			ui.SetCell(DungeonWidth, i, '│', ColorFg, ColorBg)
		}
		ui.ClearLine(i)
	}
	ui.DrawText(desc, 0, 0)
	escspace := " press (x) to continue "
	if lines+2 >= DungeonHeight {
		ui.DrawTextLine(escspace, lines+2)
		ui.SetCell(DungeonWidth, lines+2, '┘', ColorFg, ColorBg)
	} else {
		ui.DrawTextLine(escspace, lines+2)
	}

	ui.Flush()
	ui.WaitForContinue(lines + 2)
}

func (ui *gameui) WizardInfo() {
	g := ui.g
	ui.Clear()
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "Monsters: %d (%d)\n", len(g.Monsters), g.MaxMonsters())
	fmt.Fprintf(b, "Danger: %d (%d)\n", g.Danger(), g.MaxDanger())
	ui.DrawText(b.String(), 0, 0)
	ui.Flush()
	ui.WaitForContinue(-1)
}

func (ui *gameui) AptitudesText() string {
	g := ui.g
	apts := []string{}
	for apt, b := range g.Player.Aptitudes {
		if b {
			apts = append(apts, apt.String())
		}
	}
	sort.Strings(apts)
	var text string
	if len(apts) > 0 {
		text = "Aptitudes:\n" + strings.Join(apts, "\n")
	} else {
		text = "You do not have any special aptitudes."
	}
	return text
}

func (ui *gameui) AddComma(see, s string) string {
	if len(s) > 0 {
		return s + ", "
	}
	return fmt.Sprintf("You %s %s", see, s)
}
func (ui *gameui) DescribePosition(p gruid.Point, targ Targeter) {
	g := ui.g
	var desc string
	switch {
	case !g.Dungeon.Cell(p).Explored:
		desc = "You do not know what is in there."
		g.InfoEntry = desc
		return
	case !targ.Reachable(g, p):
		desc = "This is out of reach."
		g.InfoEntry = desc
		return
	}
	mons := g.MonsterAt(p)
	c, okCollectable := g.Collectables[p]
	eq, okEq := g.Equipables[p]
	rod, okRod := g.Rods[p]
	if p == g.Player.P {
		desc = "This is you"
	}
	see := "see"
	if !g.Player.LOS[p] {
		see = "saw"
	}
	if g.Dungeon.Cell(p).T == WallCell && !g.WrongWall[p] || g.Dungeon.Cell(p).T == FreeCell && g.WrongWall[p] {
		desc = ui.AddComma(see, "")
		desc += "a wall"
		g.InfoEntry = desc + "."
		return
	}
	if mons.Exists() && g.Player.LOS[p] {
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("%s (%s)", mons.Kind.Indefinite(false), ui.MonsterInfo(mons))
	}
	strt, okStair := g.Stairs[p]
	stn, okStone := g.MagicalStones[p]
	switch {
	case g.Simellas[p] > 0:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("some simellas (%d)", g.Simellas[p])
	case okCollectable:
		if c.Quantity > 1 {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("%d %s", c.Quantity, c.Consumable)
		} else {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("%s", Indefinite(c.Consumable.String(), false))
		}
	case okEq:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("%s", Indefinite(eq.String(), false))
	case okRod:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("a %v", rod)
	case okStair:
		if strt == WinStair {
			desc = ui.AddComma(see, desc)
			desc += "glowing monolith"
		} else {
			desc = ui.AddComma(see, desc)
			desc += "stairs downwards"
		}
	case okStone:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprint(Indefinite(stn.String(), false))
	case g.Doors[p] || g.WrongDoor[p]:
		desc = ui.AddComma(see, desc)
		desc += "a door"
	}
	if cld, ok := g.Clouds[p]; ok && g.Player.LOS[p] {
		if cld == CloudFire {
			desc = ui.AddComma(see, desc)
			desc += "burning flames"
		} else if cld == CloudNight {
			desc = ui.AddComma(see, desc)
			desc += "night clouds"
		} else {
			desc = ui.AddComma(see, desc)
			desc += "a dense fog"
		}
	} else if _, ok := g.Fungus[p]; ok && !g.WrongFoliage[p] || !ok && g.WrongFoliage[p] {
		desc = ui.AddComma(see, desc)
		desc += "foliage"
	} else if desc == "" {
		desc = ui.AddComma(see, desc)
		desc += "the ground"
	}
	g.InfoEntry = desc + "."
}

func (ui *gameui) ViewPositionDescription(p gruid.Point) {
	g := ui.g
	if !g.Dungeon.Cell(p).Explored {
		ui.DrawDescription("This place is unknown to you.")
		return
	}
	mons := g.MonsterAt(p)
	if mons.Exists() && g.Player.LOS[mons.P] {
		ui.HideCursor()
		ui.DrawMonsterDescription(mons)
		ui.SetCursor(p)
	} else if c, ok := g.Collectables[p]; ok {
		ui.DrawDescription(c.Consumable.Desc())
	} else if r, ok := g.Rods[p]; ok {
		ui.DrawDescription(r.Desc())
	} else if eq, ok := g.Equipables[p]; ok {
		ui.DrawDescription(eq.Desc())
	} else if strt, ok := g.Stairs[p]; ok {
		if strt == WinStair {
			desc := "This magical monolith will teleport you back to your village. It is said such monoliths were made some centuries ago by Marevor Helith. You can use it like stairs."
			if g.Depth < MaxDepth {
				desc += " Note that this is not the last floor, so you may want to find a stair and continue collecting simellas, if you're courageous enough."
			}
			ui.DrawDescription(desc)
		} else {
			desc := "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you."
			if g.Depth == WinDepth {
				desc += " If you're afraid, you could instead just win by taking the magical monolith somewhere in the same map."
			}
			ui.DrawDescription(desc)
		}
	} else if stn, ok := g.MagicalStones[p]; ok {
		ui.DrawDescription(stn.Description())
	} else if g.Doors[p] {
		ui.DrawDescription("A closed door blocks your line of sight. Doors open automatically when you or a monster stand on them. Doors are flammable.")
	} else if g.Simellas[p] > 0 {
		ui.DrawDescription("A simella is a plant with big white flowers which are used in the Underground for their medicinal properties. They can also make tasty infusions. You were actually sent here by your village to collect as many as possible of those plants.")
	} else if _, ok := g.Fungus[p]; ok && g.Dungeon.Cell(p).T == FreeCell {
		ui.DrawDescription("Blue dense foliage grows in the Underground. It is difficult to see through, and is flammable.")
	} else if g.Dungeon.Cell(p).T == WallCell {
		ui.DrawDescription("A wall is an impassable pile of rocks. It can be destructed by using some items.")
	} else {
		ui.DrawDescription("This is just plain ground.")
	}
}

func (ui *gameui) MonsterInfo(m *monster) string {
	infos := []string{}
	state := m.State.String()
	if m.Kind == MonsSatowalgaPlant && m.State == Wandering {
		state = "awaken"
	}
	infos = append(infos, state)
	for st, i := range m.Statuses {
		if i > 0 {
			infos = append(infos, monsterStatus(st).String())
		}
	}
	p := (m.HP * 100) / m.HPmax
	health := fmt.Sprintf("%d %% HP", p)
	infos = append(infos, health)
	return strings.Join(infos, ", ")
}

var CenteredCamera bool

func (ui *gameui) InView(p gruid.Point, targeting bool) bool {
	g := ui.g
	if targeting {
		return DistanceY(p, ui.cursor) <= 10 && DistanceX(p, ui.cursor) <= 39
	}
	return DistanceY(p, g.Player.P) <= 10 && DistanceX(p, g.Player.P) <= 39
}

func (ui *gameui) CameraOffset(p gruid.Point, targeting bool) (int, int) {
	g := ui.g
	if targeting {
		return p.X + 39 - ui.cursor.X, p.Y + 10 - ui.cursor.Y
	}
	return p.X + 39 - g.Player.P.X, p.Y + 10 - g.Player.P.Y
}

func (ui *gameui) InViewBorder(p gruid.Point, targeting bool) bool {
	g := ui.g
	if targeting {
		return DistanceY(p, ui.cursor) != 10 && DistanceX(p, ui.cursor) != 39
	}
	return DistanceY(p, g.Player.P) != 10 && DistanceX(p, g.Player.P) != 39
}

func (ui *gameui) DrawAtPosition(p gruid.Point, targeting bool, r rune, fg, bg uicolor) {
	g := ui.g
	if g.Highlight[p] || p == ui.cursor {
		bg, fg = fg, bg
	}
	if CenteredCamera {
		if !ui.InView(p, targeting) {
			return
		}
		x, y := ui.CameraOffset(p, targeting)
		ui.SetMapCell(x, y, r, fg, bg)
		if ui.InViewBorder(p, targeting) && g.Dungeon.Border(p) {
			for _, opos := range OutsideNeighbors(p) {
				xo, yo := ui.CameraOffset(opos, targeting)
				ui.SetMapCell(xo, yo, '#', ColorFg, ColorBgBorder)
			}
		}
		return
	}
	ui.SetMapCell(p.X, p.Y, r, fg, bg)
}

const BarCol = DungeonWidth + 2

func (ui *gameui) DrawDungeonView(m uiMode) {
	g := ui.g
	ui.Clear()
	d := g.Dungeon
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, DungeonHeight, '─', ColorFg, ColorBg)
	}
	for i := 0; i < DungeonHeight; i++ {
		ui.SetCell(DungeonWidth, i, '│', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, DungeonHeight, '┘', ColorFg, ColorBg)
	for i := range d.Cells {
		p := idx2Point(i)
		r, fgColor, bgColor := ui.PositionDrawing(p)
		ui.DrawAtPosition(p, m == TargetingMode, r, fgColor, bgColor)
	}
	line := 0
	if !ui.Small() {
		ui.SetMapCell(BarCol, line, '[', ColorFg, ColorBg)
		ui.DrawText(fmt.Sprintf(" %v", g.Player.Armour), BarCol+1, line)
		line++
		ui.SetMapCell(BarCol, line, ')', ColorFg, ColorBg)
		ui.DrawText(fmt.Sprintf(" %v", g.Player.Weapon), BarCol+1, line)
		line++
		if g.Player.Shield != NoShield {
			if g.Player.Weapon.TwoHanded() {
				ui.SetMapCell(BarCol, line, ']', ColorFg, ColorBg)
				ui.DrawText(" (unusable)", BarCol+1, line)
			} else {
				ui.SetMapCell(BarCol, line, ']', ColorFg, ColorBg)
				ui.DrawText(fmt.Sprintf(" %v", g.Player.Shield), BarCol+1, line)
			}
		}
		line++
		line++
	}
	if ui.Small() {
		ui.DrawStatusLine()
	} else {
		ui.DrawStatusBar(line)
		ui.DrawMenus()
	}
	if ui.Small() {
		ui.DrawLog(2)
	} else {
		ui.DrawLog(4)
	}
	if m != TargetingMode && m != NoFlushMode {
		ui.Flush()
	}
}

func (ui *gameui) PositionDrawing(p gruid.Point) (r rune, fgColor, bgColor uicolor) {
	g := ui.g
	m := g.Dungeon
	c := m.Cell(p)
	fgColor = ColorFg
	bgColor = ColorBg
	if !c.Explored && !g.Wizard {
		r = ' '
		bgColor = ColorBgDark
		if g.HasFreeExploredNeighbor(p) {
			r = '¤'
			fgColor = ColorFgDark
		}
		if g.DreamingMonster[p] {
			r = '☻'
			fgColor = ColorFgSleepingMonster
		}
		if g.Noise[p] {
			r = '♫'
			fgColor = ColorFgWanderingMonster
		}
		return
	}
	if g.Wizard {
		if !c.Explored && g.HasFreeExploredNeighbor(p) && !g.WizardMap {
			r = '¤'
			fgColor = ColorFgDark
			bgColor = ColorBgDark
			return
		}
		if c.T == WallCell {
			if len(g.Dungeon.FreeNeighbors(p)) == 0 {
				r = ' '
				return
			}
		}
	}
	if g.Player.LOS[p] && !g.WizardMap {
		fgColor = ColorFgLOS
		bgColor = ColorBgLOS
	} else {
		fgColor = ColorFgDark
		bgColor = ColorBgDark
	}
	if g.ExclusionsMap[p] && c.T != WallCell {
		fgColor = ColorFgExcluded
	}
	switch {
	case c.T == WallCell && (!g.WrongWall[p] || g.Wizard) || c.T == FreeCell && g.WrongWall[p] && !g.Wizard:
		r = '#'
		if g.TemporalWalls[p] {
			fgColor = ColorFgMagicPlace
		}
	case p == g.Player.P && !g.WizardMap:
		r = '@'
		fgColor = ColorFgPlayer
	default:
		r = '.'
		if _, ok := g.Fungus[p]; ok && !g.WrongFoliage[p] || !ok && g.WrongFoliage[p] {
			r = '"'
		}
		if cld, ok := g.Clouds[p]; ok && g.Player.LOS[p] {
			r = '§'
			if cld == CloudFire {
				fgColor = ColorFgWanderingMonster
			} else if cld == CloudNight {
				fgColor = ColorFgSleepingMonster
			}
		}
		if c, ok := g.Collectables[p]; ok {
			r = c.Consumable.Letter()
			fgColor = ColorFgCollectable
		} else if eq, ok := g.Equipables[p]; ok {
			r = eq.Letter()
			fgColor = ColorFgCollectable
		} else if rod, ok := g.Rods[p]; ok {
			r = rod.Letter()
			fgColor = ColorFgCollectable
		} else if strt, ok := g.Stairs[p]; ok {
			r = '>'
			if strt == WinStair {
				fgColor = ColorFgMagicPlace
				r = 'Δ'
			} else {
				fgColor = ColorFgPlace
			}
		} else if stn, ok := g.MagicalStones[p]; ok {
			r = '_'
			if stn == InertStone {
				fgColor = ColorFgPlace
			} else {
				fgColor = ColorFgMagicPlace
			}
		} else if _, ok := g.Simellas[p]; ok {
			r = '♣'
			fgColor = ColorFgSimellas
		} else if _, ok := g.Doors[p]; ok {
			r = '+'
			fgColor = ColorFgPlace
		}
		if (g.Player.LOS[p] || g.Wizard) && !g.WizardMap {
			m := g.MonsterAt(p)
			if m.Exists() {
				r = m.Kind.Letter()
				if m.Status(MonsLignified) {
					fgColor = ColorFgLignifiedMonster
				} else if m.Status(MonsConfused) {
					fgColor = ColorFgConfusedMonster
				} else if m.Status(MonsSlow) {
					fgColor = ColorFgSlowedMonster
				} else if m.State == Resting {
					fgColor = ColorFgSleepingMonster
				} else if m.State == Wandering {
					fgColor = ColorFgWanderingMonster
				} else {
					fgColor = ColorFgMonster
				}
			}
		} else if !g.Wizard && g.Noise[p] {
			r = '♫'
			fgColor = ColorFgWanderingMonster
		} else if !g.Wizard && g.DreamingMonster[p] {
			r = '☻'
			fgColor = ColorFgSleepingMonster
		}
	}
	return
}

func (ui *gameui) DrawStatusBar(line int) {
	g := ui.g
	sts := statusSlice{}
	if cld, ok := g.Clouds[g.Player.P]; ok && cld == CloudFire {
		g.Player.Statuses[StatusFlames] = 1
		defer func() {
			g.Player.Statuses[StatusFlames] = 0
		}()
	}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st)
		}
	}
	sort.Sort(sts)
	hpColor := ColorFgHPok
	switch {
	case g.Player.HP*100/g.Player.HPMax() < 30:
		hpColor = ColorFgHPcritical
	case g.Player.HP*100/g.Player.HPMax() < 70:
		hpColor = ColorFgHPwounded
	}
	mpColor := ColorFgMPok
	switch {
	case g.Player.MP*100/g.Player.MPMax() < 30:
		mpColor = ColorFgMPcritical
	case g.Player.MP*100/g.Player.MPMax() < 70:
		mpColor = ColorFgMPpartial
	}
	ui.DrawColoredText(fmt.Sprintf("HP: %d", g.Player.HP), BarCol, line, hpColor)
	line++
	ui.DrawColoredText(fmt.Sprintf("MP: %d", g.Player.MP), BarCol, line, mpColor)
	line++
	line++
	ui.DrawText(fmt.Sprintf("Simellas: %d", g.Player.Simellas), BarCol, line)
	line++
	if g.Depth == -1 {
		ui.DrawText("Depth: Out!", BarCol, line)
	} else {
		ui.DrawText(fmt.Sprintf("Depth: %d", g.Depth), BarCol, line)
	}
	line++
	ui.DrawText(fmt.Sprintf("Turns: %.1f", float64(g.Turn)/10), BarCol, line)
	line++
	for _, st := range sts {
		fg := ColorFgStatusOther
		if st.Good() {
			fg = ColorFgStatusGood
			t := 13
			if g.Player.Statuses[StatusBerserk] > 0 {
				t -= 3
			}
			if g.Player.Statuses[StatusSlow] > 0 {
				t += 3
			}
			if g.Player.Expire[st] >= g.Ev.Rank() && g.Player.Expire[st]-g.Ev.Rank() <= t {
				fg = ColorFgStatusExpire
			}
		} else if st.Bad() {
			fg = ColorFgStatusBad
		}
		if g.Player.Statuses[st] > 1 {
			ui.DrawColoredText(fmt.Sprintf("%s(%d)", st, g.Player.Statuses[st]), BarCol, line, fg)
		} else {
			ui.DrawColoredText(st.String(), BarCol, line, fg)
		}
		line++
	}
}

func (ui *gameui) DrawStatusLine() {
	g := ui.g
	sts := statusSlice{}
	if cld, ok := g.Clouds[g.Player.P]; ok && cld == CloudFire {
		g.Player.Statuses[StatusFlames] = 1
		defer func() {
			g.Player.Statuses[StatusFlames] = 0
		}()
	}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st)
		}
	}
	sort.Sort(sts)
	hpColor := ColorFgHPok
	switch {
	case g.Player.HP*100/g.Player.HPMax() < 30:
		hpColor = ColorFgHPcritical
	case g.Player.HP*100/g.Player.HPMax() < 70:
		hpColor = ColorFgHPwounded
	}
	mpColor := ColorFgMPok
	switch {
	case g.Player.MP*100/g.Player.MPMax() < 30:
		mpColor = ColorFgMPcritical
	case g.Player.MP*100/g.Player.MPMax() < 70:
		mpColor = ColorFgMPpartial
	}
	line := DungeonHeight
	col := 2
	ui.DrawText(" ", col, line)
	col++
	ui.SetMapCell(col, line, ')', ColorFg, ColorBg)
	col++
	weapon := fmt.Sprintf("%s ", g.Player.Weapon.Short())
	ui.DrawText(weapon, col, line)
	col += utf8.RuneCountInString(weapon)
	ui.SetMapCell(col, line, '[', ColorFg, ColorBg)
	col++
	armour := fmt.Sprintf("%s ", g.Player.Armour.Short())
	ui.DrawText(armour, col, line)
	col += utf8.RuneCountInString(armour)
	if g.Player.Shield != NoShield {
		ui.SetMapCell(col, line, ']', ColorFg, ColorBg)
		col++
		shield := fmt.Sprintf("%s ", g.Player.Shield.Short())
		ui.DrawText(shield, col, line)
		col += utf8.RuneCountInString(shield)
	}
	ui.SetMapCell(col, line, '♣', ColorFg, ColorBg)
	col++
	simellas := fmt.Sprintf(":%d ", g.Player.Simellas)
	ui.DrawText(simellas, col, line)
	col += utf8.RuneCountInString(simellas)
	var depth string
	if g.Depth == -1 {
		depth = "D: Out! "
	} else {
		depth = fmt.Sprintf("D:%d ", g.Depth)
	}
	ui.DrawText(depth, col, line)
	col += utf8.RuneCountInString(depth)
	turns := fmt.Sprintf("T:%.1f ", float64(g.Turn)/10)
	ui.DrawText(turns, col, line)
	col += utf8.RuneCountInString(turns)
	hp := fmt.Sprintf("HP:%2d ", g.Player.HP)
	ui.DrawColoredText(hp, col, line, hpColor)
	col += utf8.RuneCountInString(hp)
	mp := fmt.Sprintf("MP:%d ", g.Player.MP)
	ui.DrawColoredText(mp, col, line, mpColor)
	col += utf8.RuneCountInString(mp)
	if len(sts) > 0 {
		ui.DrawText("| ", col, line)
		col += 2
	}
	for _, st := range sts {
		fg := ColorFgStatusOther
		if st.Good() {
			fg = ColorFgStatusGood
			t := 13
			if g.Player.Statuses[StatusBerserk] > 0 {
				t -= 3
			}
			if g.Player.Statuses[StatusSlow] > 0 {
				t += 3
			}
			if g.Player.Expire[st] >= g.Ev.Rank() && g.Player.Expire[st]-g.Ev.Rank() <= t {
				fg = ColorFgStatusExpire
			}
		} else if st.Bad() {
			fg = ColorFgStatusBad
		}
		var sttext string
		if g.Player.Statuses[st] > 1 {
			sttext = fmt.Sprintf("%s(%d) ", st.Short(), g.Player.Statuses[st])
		} else {
			sttext = fmt.Sprintf("%s ", st.Short())
		}
		ui.DrawColoredText(sttext, col, line, fg)
		col += utf8.RuneCountInString(sttext)
	}
}

func (ui *gameui) LogColor(e logEntry) uicolor {
	fg := ColorFg
	switch e.Style {
	case logCritic:
		fg = ColorRed
	case logPlayerHit:
		fg = ColorGreen
	case logMonsterHit:
		fg = ColorOrange
	case logSpecial:
		fg = ColorMagenta
	case logStatusEnd:
		fg = ColorViolet
	case logError:
		fg = ColorRed
	}
	return fg
}

func (ui *gameui) DrawLog(lines int) {
	g := ui.g
	min := len(g.Log) - lines
	if min < 0 {
		min = 0
	}
	l := len(g.Log) - 1
	if l < lines {
		lines = l + 1
	}
	for i := lines; i > 0 && l >= 0; i-- {
		cols := 0
		first := true
		to := l
		for l >= 0 {
			e := g.Log[l]
			el := utf8.RuneCountInString(e.String())
			if e.Tick {
				el += 2
			}
			cols += el + 1
			if !first && cols > DungeonWidth {
				l++
				break
			}
			if e.Tick || l <= i {
				break
			}
			first = false
			l--
		}
		if l < 0 {
			l = 0
		}
		col := 0
		for ln := l; ln <= to; ln++ {
			e := g.Log[ln]
			fguicolor := ui.LogColor(e)
			if e.Tick {
				ui.DrawColoredText("•", 0, DungeonHeight+i, ColorYellow)
				col += 2
			}
			ui.DrawColoredText(e.String(), col, DungeonHeight+i, fguicolor)
			col += utf8.RuneCountInString(e.String()) + 1
		}
		l--
	}
}

func InRuneSlice(r rune, s []rune) bool {
	for _, rr := range s {
		if r == rr {
			return true
		}
	}
	return false
}

func (ui *gameui) RunesForKeyAction(k keyAction) string {
	runes := []rune{}
	for r, ka := range GameConfig.RuneNormalModeKeys {
		if k == ka && !InRuneSlice(r, runes) {
			runes = append(runes, r)
		}
	}
	for r, ka := range GameConfig.RuneTargetModeKeys {
		if k == ka && !InRuneSlice(r, runes) {
			runes = append(runes, r)
		}
	}
	chars := strings.Split(string(runes), "")
	sort.Strings(chars)
	text := strings.Join(chars, " or ")
	return text
}

type keyConfigAction int

const (
	NavigateKeys keyConfigAction = iota
	ChangeKeys
	ResetKeys
	QuitKeyConfig
)

func (ui *gameui) ChangeKeys() {
	g := ui.g
	lines := DungeonHeight
	nmax := len(configurableKeyActions) - lines
	n := 0
	s := 0
loop:
	for {
		ui.DrawDungeonView(NoFlushMode)
		if n >= nmax {
			n = nmax
		}
		if n < 0 {
			n = 0
		}
		to := n + lines
		if to >= len(configurableKeyActions) {
			to = len(configurableKeyActions)
		}
		for i := n; i < to; i++ {
			ka := configurableKeyActions[i]
			desc := ka.NormalModeDescription()
			if !ka.NormalModeKey() {
				desc = ka.TargetingModeDescription()
			}
			bg := ui.ListItemBG(i)
			ui.ClearLineWithColor(i-n, bg)
			desc = fmt.Sprintf(" %-36s %s", desc, ui.RunesForKeyAction(ka))
			if i == s {
				ui.DrawColoredTextOnBG(desc, 0, i-n, ColorYellow, bg)
			} else {
				ui.DrawColoredTextOnBG(desc, 0, i-n, ColorFg, bg)
			}
		}
		ui.ClearLine(lines)
		ui.DrawStyledTextLine(" add key (a) up/down (arrows/u/d) reset (R) quit (x) ", lines, FooterLine)
		ui.Flush()

		var action keyConfigAction
		s, action = ui.KeyMenuAction(s)
		if s >= len(configurableKeyActions) {
			s = len(configurableKeyActions) - 1
		}
		if s < 0 {
			s = 0
		}
		if s < n+1 {
			n -= 12
		}
		if s > n+lines-2 {
			n += 12
		}
		switch action {
		case ChangeKeys:
			ui.DrawStyledTextLine(" insert new key ", lines, FooterLine)
			ui.Flush()
			r := ui.ReadRuneKey()
			if r == 0 {
				continue loop
			}
			if FixedRuneKey(r) {
				g.Printf("You cannot rebind “%c”.", r)
				continue loop
			}
			CustomKeys = true
			ka := configurableKeyActions[s]
			if ka.NormalModeKey() {
				GameConfig.RuneNormalModeKeys[r] = ka
			} else {
				delete(GameConfig.RuneNormalModeKeys, r)
			}
			if ka.TargetingModeKey() {
				GameConfig.RuneTargetModeKeys[r] = ka
			} else {
				delete(GameConfig.RuneTargetModeKeys, r)
			}
			err := g.SaveConfig()
			if err != nil {
				g.Print(err.Error())
			}
		case QuitKeyConfig:
			break loop
		case ResetKeys:
			ApplyDefaultKeyBindings()
			err := g.SaveConfig()
			//err := g.RemoveDataFile("config.gob")
			if err != nil {
				g.Print(err.Error())
			}
		}
	}
}

func (ui *gameui) DrawPreviousLogs() {
	g := ui.g
	bottom := 4
	if ui.Small() {
		bottom = 2
	}
	lines := DungeonHeight + bottom
	nmax := len(g.Log) - lines
	n := nmax
loop:
	for {
		ui.DrawDungeonView(NoFlushMode)
		if n >= nmax {
			n = nmax
		}
		if n < 0 {
			n = 0
		}
		to := n + lines
		if to >= len(g.Log) {
			to = len(g.Log)
		}
		for i := 0; i < bottom; i++ {
			ui.SetCell(DungeonWidth, DungeonHeight+i, '│', ColorFg, ColorBg)
		}
		for i := n; i < to; i++ {
			e := g.Log[i]
			fguicolor := ui.LogColor(e)
			ui.ClearLine(i - n)
			rc := utf8.RuneCountInString(e.String())
			if e.Tick {
				rc += 2
			}
			if rc >= DungeonWidth {
				for j := DungeonWidth; j < 103; j++ {
					ui.SetCell(j, i-n, ' ', ColorFg, ColorBg)
				}
			}
			if e.Tick {
				ui.DrawColoredText("•", 0, i-n, ColorYellow)
				ui.DrawColoredText(e.String(), 2, i-n, fguicolor)
			} else {
				ui.DrawColoredText(e.String(), 0, i-n, fguicolor)
			}
		}
		for i := len(g.Log); i < DungeonHeight+bottom; i++ {
			ui.ClearLine(i - n)
		}
		ui.ClearLine(lines)
		s := fmt.Sprintf(" half-page up/down (u/d) quit (x) — (%d/%d) \n", len(g.Log)-to, len(g.Log))
		ui.DrawStyledTextLine(s, lines, FooterLine)
		ui.Flush()
		var quit bool
		n, quit = ui.Scroll(n)
		if quit {
			break loop
		}
	}
}

func (ui *gameui) DrawMonsterDescription(mons *monster) {
	s := mons.Kind.Desc()
	s += " " + fmt.Sprintf("They can hit for up to %d damage.", mons.Kind.BaseAttack())
	s += " " + fmt.Sprintf("They have around %d HP.", mons.Kind.MaxHP())
	ui.DrawDescription(s)
}

func (ui *gameui) DrawConsumableDescription(c consumable) {
	ui.DrawDescription(c.Desc())
}

func (ui *gameui) DrawDescription(desc string) {
	ui.DrawDungeonView(NoFlushMode)
	desc = formatText(desc, TextWidth)
	lines := strings.Count(desc, "\n")
	for i := 0; i <= lines+2; i++ {
		ui.ClearLine(i)
	}
	ui.DrawText(desc, 0, 0)
	ui.DrawTextLine(" press (x) to continue ", lines+2)
	ui.Flush()
	ui.WaitForContinue(lines + 2)
	ui.DrawDungeonView(NoFlushMode)
}

func (ui *gameui) DrawText(text string, x, y int) {
	ui.DrawColoredText(text, x, y, ColorFg)
}

func (ui *gameui) DrawColoredText(text string, x, y int, fg uicolor) {
	ui.DrawColoredTextOnBG(text, x, y, fg, ColorBg)
}

func (ui *gameui) DrawColoredTextOnBG(text string, x, y int, fg, bg uicolor) {
	col := 0
	for _, r := range text {
		if r == '\n' {
			y++
			col = 0
			continue
		}
		if x+col >= UIWidth {
			break
		}
		ui.SetCell(x+col, y, r, fg, bg)
		col++
	}
}

func (ui *gameui) DrawLine(lnum int) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, lnum, '┤', ColorFg, ColorBg)
}

func (ui *gameui) DrawTextLine(text string, lnum int) {
	ui.DrawStyledTextLine(text, lnum, NormalLine)
}

type linestyle int

const (
	NormalLine linestyle = iota
	HeaderLine
	FooterLine
)

func (ui *gameui) DrawInfoLine(text string) {
	ui.ClearLineWithColor(DungeonHeight+1, ColorBgBorder)
	ui.DrawColoredTextOnBG(text, 0, DungeonHeight+1, ColorBlue, ColorBgBorder)
}

func (ui *gameui) DrawStyledTextLine(text string, lnum int, st linestyle) {
	nchars := utf8.RuneCountInString(text)
	dist := (DungeonWidth - nchars) / 2
	for i := 0; i < dist; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	switch st {
	case HeaderLine:
		ui.DrawColoredText(text, dist, lnum, ColorYellow)
	case FooterLine:
		ui.DrawColoredText(text, dist, lnum, ColorCyan)
	default:
		ui.DrawColoredText(text, dist, lnum, ColorFg)
	}
	for i := dist + nchars; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	switch st {
	case HeaderLine:
		ui.SetCell(DungeonWidth, lnum, '┐', ColorFg, ColorBg)
	case FooterLine:
		ui.SetCell(DungeonWidth, lnum, '┘', ColorFg, ColorBg)
	default:
		ui.SetCell(DungeonWidth, lnum, '┤', ColorFg, ColorBg)
	}
}

func (ui *gameui) ClearLine(lnum int) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, ' ', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, lnum, '│', ColorFg, ColorBg)
}

func (ui *gameui) ClearLineWithColor(lnum int, bg uicolor) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, ' ', ColorFg, bg)
	}
	ui.SetCell(DungeonWidth, lnum, '│', ColorFg, ColorBg)
}

func (ui *gameui) ListItemBG(i int) uicolor {
	bg := ColorBg
	if i%2 == 1 {
		bg = ColorBgBorder
	}
	return bg
}

func (ui *gameui) ConsumableItem(i, lnum int, c consumable, fg uicolor) {
	g := ui.g
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, lnum, fg, bg)
}

func (ui *gameui) SelectProjectile(ev event) error {
	g := ui.g
	desc := false
	for {
		cs := g.SortedProjectiles()
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuThrow.String(), MenuCols[MenuThrow][0], DungeonHeight, ColorCyan)
		}
		if desc {
			ui.DrawColoredText("Describe", 0, 0, ColorBlue)
			col := utf8.RuneCountInString("Describe")
			ui.DrawText(" which projectile? (press ? or click here for throwing menu)", col, 0)
		} else {
			ui.DrawColoredText("Throw", 0, 0, ColorOrange)
			col := utf8.RuneCountInString("Throw")
			ui.DrawText(" which projectile? (press ? or click here for describe menu)", col, 0)
		}
		for i, c := range cs {
			ui.ConsumableItem(i, i+1, c, ColorFg)
		}
		ui.DrawTextLine(" press (x) to cancel ", len(cs)+1)
		ui.Flush()
		index, alt, err := ui.Select(len(cs))
		if alt {
			desc = !desc
			continue
		}
		if err == nil {
			ui.ConsumableItem(index, index+1, cs[index], ColorYellow)
			ui.Flush()
			time.Sleep(75 * time.Millisecond)
			if desc {
				ui.DrawDescription(cs[index].Desc())
				continue
			}
			err = cs[index].Use(g, ev)
		}
		return err
	}
}

func (ui *gameui) SelectPotion(ev event) error {
	g := ui.g
	desc := false
	for {
		cs := g.SortedPotions()
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuDrink.String(), MenuCols[MenuDrink][0], DungeonHeight, ColorCyan)
		}
		if desc {
			ui.DrawColoredText("Describe", 0, 0, ColorBlue)
			col := utf8.RuneCountInString("Describe")
			ui.DrawText(" which potion? (press ? or click here for quaff menu)", col, 0)
		} else {
			ui.DrawColoredText("Drink", 0, 0, ColorGreen)
			col := utf8.RuneCountInString("Drink")
			ui.DrawText(" which potion? (press ? or click here for description menu)", col, 0)
		}
		for i, c := range cs {
			ui.ConsumableItem(i, i+1, c, ColorFg)
		}
		ui.DrawTextLine(" press (x) to cancel ", len(cs)+1)
		ui.Flush()
		index, alt, err := ui.Select(len(cs))
		if alt {
			desc = !desc
			continue
		}
		if err == nil {
			ui.ConsumableItem(index, index+1, cs[index], ColorYellow)
			ui.Flush()
			time.Sleep(75 * time.Millisecond)
			if desc {
				ui.DrawDescription(cs[index].Desc())
				continue
			}
			err = cs[index].Use(g, ev)
		}
		return err
	}
}

func (ui *gameui) RodItem(i, lnum int, r rod, fg uicolor) {
	g := ui.g
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	mc := r.MaxCharge()
	if g.Player.Armour == CelmistRobe {
		mc += 2
	}
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d/%d charges, %d mana cost)",
		rune(i+97), r, g.Player.Rods[r].Charge, mc, r.MPCost()), 0, lnum, fg, bg)
}

func (ui *gameui) SelectRod(ev event) error {
	g := ui.g
	desc := false
	for {
		rs := g.SortedRods()
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuEvoke.String(), MenuCols[MenuEvoke][0], DungeonHeight, ColorCyan)
		}
		if desc {
			ui.DrawColoredText("Describe", 0, 0, ColorBlue)
			col := utf8.RuneCountInString("Describe")
			ui.DrawText(" which rod? (press ? or click here for evocation menu)", col, 0)
		} else {
			ui.DrawColoredText("Evoke", 0, 0, ColorCyan)
			col := utf8.RuneCountInString("Evoke")
			ui.DrawText(" which rod? (press ? or click here for description menu)", col, 0)
		}
		for i, r := range rs {
			ui.RodItem(i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press (x) to cancel ", len(rs)+1)
		ui.Flush()
		index, alt, err := ui.Select(len(rs))
		if alt {
			desc = !desc
			continue
		}
		if err == nil {
			ui.RodItem(index, index+1, rs[index], ColorYellow)
			ui.Flush()
			time.Sleep(75 * time.Millisecond)
			if desc {
				ui.DrawDescription(rs[index].Desc())
				continue
			}
			err = rs[index].Use(g, ev)
		}
		return err
	}
}

func (ui *gameui) AbbreviatedItem(lnum, col int, name string, q int) {
	bg := ui.ListItemBG(lnum)
	if col == 0 {
		ui.ClearLineWithColor(lnum, bg)
	}
	ui.DrawColoredTextOnBG(fmt.Sprintf("- %s (%d)", name, q), col, lnum, ColorFg, bg)
}

func (ui *gameui) ViewAllTitle(lnum, col int, name string) {
	bg := ui.ListItemBG(lnum)
	if col == 0 {
		ui.ClearLineWithColor(lnum, bg)
	}
	ui.DrawColoredTextOnBG(name, col, lnum, ColorYellow, bg)
}

func (ui *gameui) NextCell(lnum, col int) (int, int) {
	if lnum >= DungeonHeight-1 {
		lnum = 0
		col = DungeonWidth / 2
	} else {
		lnum++
	}
	return lnum, col
}

func (ui *gameui) ViewAll() {
	g := ui.g
	sr := g.SortedRods()
	lnum := 0
	col := 0
	ui.ViewAllTitle(lnum, col, "Rods")
	lnum, col = ui.NextCell(lnum, col)
	for _, c := range sr {
		ui.AbbreviatedItem(lnum, col, c.Name(), ui.g.Player.Rods[c].Charge)
		lnum, col = ui.NextCell(lnum, col)
	}
	ui.ViewAllTitle(lnum, col, "Projectiles")
	lnum, col = ui.NextCell(lnum, col)
	sp := g.SortedProjectiles()
	for _, c := range sp {
		ui.AbbreviatedItem(lnum, col, c.String(), g.Player.Consumables[c])
		lnum, col = ui.NextCell(lnum, col)
	}
	ui.ViewAllTitle(lnum, col, "Potions")
	lnum, col = ui.NextCell(lnum, col)
	sp = g.SortedPotions()
	for _, c := range sp {
		p := c.(potion)
		ui.AbbreviatedItem(lnum, col, p.Name(), g.Player.Consumables[c])
		lnum, col = ui.NextCell(lnum, col)
	}
	lmax := lnum
	if col != 0 {
		lmax = DungeonHeight
	}
	ui.DrawTextLine(" press (x) to cancel ", lmax)
	ui.Flush()
	ui.WaitForContinue(lmax)
}

func (ui *gameui) ActionItem(i, lnum int, ka keyAction, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	desc := ka.NormalModeDescription()
	if !ka.NormalModeKey() {
		desc = ka.TargetingModeDescription()
	}
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s", rune(i+97), desc), 0, lnum, fg, bg)
}

var menuActions = []keyAction{
	KeyCharacterInfo,
	KeyLogs,
	KeyMenuCommandHelp,
	KeyMenuTargetingHelp,
	KeyConfigure,
	KeySave,
	KeyQuit,
}

func (ui *gameui) SelectAction(actions []keyAction, ev event) (keyAction, error) {
	for {
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuOther.String(), MenuCols[MenuOther][0], DungeonHeight, ColorCyan)
		}
		ui.DrawColoredText("Choose", 0, 0, ColorCyan)
		col := utf8.RuneCountInString("Choose")
		ui.DrawText(" which action?", col, 0)
		for i, r := range actions {
			ui.ActionItem(i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press (x) to cancel ", len(actions)+1)
		ui.Flush()
		index, alt, err := ui.Select(len(actions))
		if alt {
			continue
		}
		if err != nil {
			ui.DrawDungeonView(NoFlushMode)
			return KeyExamine, err
		}
		ui.ActionItem(index, index+1, actions[index], ColorYellow)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
		ui.DrawDungeonView(NoFlushMode)
		return actions[index], nil
	}
}

type setting int

const (
	setKeys setting = iota
	invertLOS
	toggleLayout
	toggleTiles
)

func (s setting) String() (text string) {
	switch s {
	case setKeys:
		text = "Change key bindings"
	case invertLOS:
		text = "Toggle dark/light LOS"
	case toggleLayout:
		text = "Toggle normal/compact layout"
	case toggleTiles:
		text = "Toggle Tiles/Ascii display"
	}
	return text
}

var settingsActions = []setting{
	setKeys,
	invertLOS,
	toggleLayout,
}

func (ui *gameui) ConfItem(i, lnum int, s setting, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s", rune(i+97), s), 0, lnum, fg, bg)
}

func (ui *gameui) SelectConfigure(actions []setting) (setting, error) {
	for {
		ui.ClearLine(0)
		ui.DrawColoredText("Perform", 0, 0, ColorCyan)
		col := utf8.RuneCountInString("Perform")
		ui.DrawText(" which change?", col, 0)
		for i, r := range actions {
			ui.ConfItem(i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press (x) to cancel ", len(actions)+1)
		ui.Flush()
		index, alt, err := ui.Select(len(actions))
		if alt {
			continue
		}
		if err != nil {
			ui.DrawDungeonView(NoFlushMode)
			return setKeys, err
		}
		ui.ConfItem(index, index+1, actions[index], ColorYellow)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
		ui.DrawDungeonView(NoFlushMode)
		return actions[index], nil
	}
}

func (ui *gameui) HandleSettingAction() error {
	g := ui.g
	s, err := ui.SelectConfigure(settingsActions)
	if err != nil {
		return err
	}
	switch s {
	case setKeys:
		ui.ChangeKeys()
	case invertLOS:
		GameConfig.DarkLOS = !GameConfig.DarkLOS
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
		if GameConfig.DarkLOS {
			ApplyDarkLOS()
		} else {
			ApplyLightLOS()
		}
	case toggleLayout:
		ui.ApplyToggleLayout()
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
	case toggleTiles:
		ui.ApplyToggleTiles()
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
	}
	return nil
}

func (ui *gameui) WizardItem(i, lnum int, s wizardAction, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s", rune(i+97), s), 0, lnum, fg, bg)
}

func (ui *gameui) SelectWizardMagic(actions []wizardAction) (wizardAction, error) {
	for {
		ui.ClearLine(0)
		ui.DrawColoredText("Evoke", 0, 0, ColorCyan)
		col := utf8.RuneCountInString("Evoke")
		ui.DrawText(" which magic?", col, 0)
		for i, r := range actions {
			ui.WizardItem(i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press (x) to cancel ", len(actions)+1)
		ui.Flush()
		index, alt, err := ui.Select(len(actions))
		if alt {
			continue
		}
		if err != nil {
			ui.DrawDungeonView(NoFlushMode)
			return WizardInfoAction, err
		}
		ui.WizardItem(index, index+1, actions[index], ColorYellow)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
		ui.DrawDungeonView(NoFlushMode)
		return actions[index], nil
	}
}

func (ui *gameui) DrawMenus() {
	line := DungeonHeight
	for i, cols := range MenuCols[0 : len(MenuCols)-1] {
		if cols[0] >= 0 {
			if menu(i) == ui.menuHover {
				ui.DrawColoredText(menu(i).String(), cols[0], line, ColorBlue)
			} else {
				ui.DrawColoredText(menu(i).String(), cols[0], line, ColorViolet)
			}
		}
	}
	interactMenu := ui.UpdateInteractButton()
	if interactMenu == "" {
		return
	}
	i := len(MenuCols) - 1
	cols := MenuCols[i]
	if menu(i) == ui.menuHover {
		ui.DrawColoredText(interactMenu, cols[0], line, ColorBlue)
	} else {
		ui.DrawColoredText(interactMenu, cols[0], line, ColorViolet)
	}
}
