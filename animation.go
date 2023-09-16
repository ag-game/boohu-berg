package main

import (
	"time"

	"codeberg.org/anaseto/gruid"
)

func (ui *gameui) SwappingAnimation(mpos, ppos gruid.Point) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	_, fgm, bgColorm := ui.PositionDrawing(mpos)
	_, _, bgColorp := ui.PositionDrawing(ppos)
	ui.DrawAtPosition(mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	ui.DrawAtPosition(mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
}

func (ui *gameui) TeleportAnimation(from, to gruid.Point, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(from)
	_, _, bgColort := ui.PositionDrawing(to)
	ui.DrawAtPosition(from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	if showto {
		ui.DrawAtPosition(from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(to, true, 'Φ', ColorCyan, bgColort)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (ui *gameui) ProjectileTrajectoryAnimation(ray []gruid.Point, fg uicolor) {
	if DisableAnimations {
		return
	}
	for i := len(ray) - 1; i >= 0; i-- {
		p := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(p)
		ui.DrawAtPosition(p, true, '•', fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(p, true, r, fgColor, bgColor)
	}
}

func (ui *gameui) MonsterProjectileAnimation(ray []gruid.Point, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		p := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(p)
		ui.DrawAtPosition(p, true, r, fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(p, true, or, fgColor, bgColor)
	}
}

func (ui *gameui) ExplosionAnimationAt(p gruid.Point, fg uicolor) {
	g := ui.g
	_, _, bgColor := ui.PositionDrawing(p)
	mons := g.MonsterAt(p)
	r := ';'
	switch RandInt(9) {
	case 0, 6:
		r = ','
	case 1:
		r = '}'
	case 2:
		r = '%'
	case 3, 7:
		r = ':'
	case 4:
		r = '\\'
	case 5:
		r = '~'
	}
	if mons.Exists() || g.Player.P == p {
		r = '√'
	}
	//ui.DrawAtPosition(p, true, r, fg, bgColor)
	ui.DrawAtPosition(p, true, r, bgColor, fg)
}

func (ui *gameui) ExplosionAnimation(es explosionStyle, p gruid.Point) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(20 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	if es == WallExplosion || es == AroundWallExplosion {
		colors[0] = ColorFgExplosionWallStart
		colors[1] = ColorFgExplosionWallEnd
	}
	for i := 0; i < 3; i++ {
		nb := g.Dungeon.FreeNeighbors(p)
		if es != AroundWallExplosion {
			nb = append(nb, p)
		}
		for _, npos := range nb {
			fg := colors[RandInt(2)]
			if !g.Player.LOS[npos] {
				continue
			}
			ui.ExplosionAnimationAt(npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *gameui) TormentExplosionAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(20 * time.Millisecond)
	colors := [3]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd, ColorFgMagicPlace}
	for i := 0; i < 3; i++ {
		for npos, b := range g.Player.LOS {
			if !b {
				continue
			}
			fg := colors[RandInt(3)]
			ui.ExplosionAnimationAt(npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *gameui) WallExplosionAnimation(p gruid.Point) {
	if DisableAnimations {
		return
	}
	colors := [2]uicolor{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := ui.PositionDrawing(p)
		//ui.DrawAtPosition(p, true, '☼', fg, bgColor)
		ui.DrawAtPosition(p, true, '☼', bgColor, fg)
		ui.Flush()
		time.Sleep(25 * time.Millisecond)
	}
}

func (ui *gameui) FireBoltAnimation(ray []gruid.Point) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			p := ray[i]
			_, _, bgColor := ui.PositionDrawing(p)
			mons := g.MonsterAt(p)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			if mons.Exists() {
				r = '√'
			}
			//ui.DrawAtPosition(p, true, r, fg, bgColor)
			ui.DrawAtPosition(p, true, r, bgColor, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *gameui) SlowingMagaraAnimation(ray []gruid.Point) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgConfusedMonster, ColorFgMagicPlace}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			p := ray[i]
			_, _, bgColor := ui.PositionDrawing(p)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			ui.DrawAtPosition(p, true, r, bgColor, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *gameui) ProjectileSymbol(dir direction) (r rune) {
	switch dir {
	case E, ENE, ESE, WNW, W, WSW:
		r = '—'
	case NE, SW:
		r = '/'
	case NNE, N, NNW, SSW, S, SSE:
		r = '|'
	case NW, SE:
		r = '\\'
	}
	return r
}

func (ui *gameui) ThrowAnimation(ray []gruid.Point, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := len(ray) - 1; i >= 0; i-- {
		p := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(p)
		ui.DrawAtPosition(p, true, ui.ProjectileSymbol(Dir(p, g.Player.P)), ColorFgProjectile, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(p, true, r, fgColor, bgColor)
	}
	if hit {
		p := ray[0]
		ui.HitAnimation(p, true)
	}
}

func (ui *gameui) MonsterJavelinAnimation(ray []gruid.Point, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		p := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(p)
		ui.DrawAtPosition(p, true, ui.ProjectileSymbol(Dir(p, g.Player.P)), ColorFgMonster, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(p, true, r, fgColor, bgColor)
	}
}

func (ui *gameui) HitAnimation(p gruid.Point, targeting bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	if !g.Player.LOS[p] {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	_, _, bgColor := ui.PositionDrawing(p)
	mons := g.MonsterAt(p)
	if mons.Exists() || p == g.Player.P {
		ui.DrawAtPosition(p, targeting, '√', ColorFgAnimationHit, bgColor)
	} else {
		ui.DrawAtPosition(p, targeting, '∞', ColorFgAnimationHit, bgColor)
	}
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
}

func (ui *gameui) LightningHitAnimation(targets []gruid.Point) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 2; j++ {
		for _, p := range targets {
			_, _, bgColor := ui.PositionDrawing(p)
			mons := g.MonsterAt(p)
			if mons.Exists() || p == g.Player.P {
				ui.DrawAtPosition(p, false, '√', bgColor, colors[RandInt(2)])
			} else {
				ui.DrawAtPosition(p, false, '∞', bgColor, colors[RandInt(2)])
			}
		}
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
	}
}

func (ui *gameui) WoundedAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	r, _, bg := ui.PositionDrawing(g.Player.P)
	ui.DrawAtPosition(g.Player.P, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g.Player.P, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		time.Sleep(50 * time.Millisecond)
	}
}

func (ui *gameui) DrinkingPotionAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	time.Sleep(25 * time.Millisecond)
	r, fg, bg := ui.PositionDrawing(g.Player.P)
	ui.DrawAtPosition(g.Player.P, false, r, ColorGreen, bg)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	ui.DrawAtPosition(g.Player.P, false, r, ColorYellow, bg)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	ui.DrawAtPosition(g.Player.P, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) StatusEndAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	r, _, bg := ui.PositionDrawing(g.Player.P)
	ui.DrawAtPosition(g.Player.P, false, r, ColorViolet, bg)
	ui.Flush()
	time.Sleep(100 * time.Millisecond)
}

func (ui *gameui) MenuSelectedAnimation(m menu, ok bool) {
	if DisableAnimations {
		return
	}
	if !ui.Small() {
		var message string
		if m == MenuInteract {
			message = ui.UpdateInteractButton()
		} else {
			message = m.String()
		}
		if message == "" {
			return
		}
		if ok {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorCyan)
		} else {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorMagenta)
		}
		ui.Flush()
		time.Sleep(25 * time.Millisecond)
		ui.DrawColoredText(m.String(), MenuCols[m][0], DungeonHeight, ColorViolet)
	}
}

func (ui *gameui) MagicMappingAnimation(border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		p := idx2Point(i)
		r, fg, bg := ui.PositionDrawing(p)
		ui.DrawAtPosition(p, false, r, fg, bg)
	}
	ui.Flush()
	time.Sleep(12 * time.Millisecond)
}
