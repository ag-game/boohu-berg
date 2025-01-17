//go:build tk

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"unicode/utf8"

	"codeberg.org/anaseto/gothic"
	"codeberg.org/anaseto/gruid"
)

type gameui struct {
	g         *game
	ir        *gothic.Interpreter
	cursor    gruid.Point
	stty      string
	cache     map[UICell]*image.RGBA
	width     int
	height    int
	mousepos  gruid.Point
	menuHover menu
	itemHover int
	canvas    *image.RGBA
}

func (ui *gameui) Init() error {
	ui.canvas = image.NewRGBA(image.Rect(0, 0, UIWidth*16, UIHeight*24))
	ui.ir = gothic.NewInterpreter(`
wm title . "Boohu Tk"
wm resizable . 0 0
set width [expr {16 * 100}]
set height [expr {24 * 26}]
wm geometry . =${width}x$height
set can [canvas .c -width $width -height $height -background #002b36]
grid $can -row 0 -column 0
focus $can
image create photo gamescreen -width $width -height $height -palette 256/256/256
image create photo bufscreen -width $width -height $height -palette 256/256/256
$can create image 0 0 -anchor nw -image gamescreen
`)
	ui.InitElements()
	ui.ir.RegisterCommand("GetKey", func(c, keysym string) {
		var s string
		if c != "" {
			s = c
		} else {
			s = keysym
		}
		if len(ch) < cap(ch) {
			ch <- uiInput{key: s}
		}
	})
	ui.ir.RegisterCommand("MouseDown", func(x, y, b int) {
		if len(ch) < cap(ch) {
			ch <- uiInput{mouse: true, mouseX: (x - 1) / ui.width, mouseY: (y - 1) / ui.height, button: b - 1}
		}
	})
	ui.ir.RegisterCommand("MouseMotion", func(x, y int) {
		nx := (x - 1) / ui.width
		ny := (y - 1) / ui.height
		if nx != ui.mousepos.X || ny != ui.mousepos.Y {
			ui.mousepos.X = nx
			ui.mousepos.Y = ny
			if len(ch) < cap(ch) {
				ch <- uiInput{mouse: true, mouseX: nx, mouseY: ny, button: -1}
			}
		}
	})
	ui.ir.Eval(`
bind .c <Key> {
	GetKey %A %K
}
bind .c <Motion> {
	MouseMotion %x %y
}
bind .c <ButtonPress> {
	MouseDown %x %y %b
}
`)
	ui.menuHover = -1

	SolarizedPalette()
	ui.HideCursor()
	settingsActions = append(settingsActions, toggleTiles)
	GameConfig.Tiles = true
	return nil
}

func (ui *gameui) InitElements() error {
	ui.width = 16
	ui.height = 24
	ui.cache = make(map[UICell]*image.RGBA)
	return nil
}

var ch chan uiInput
var interrupt chan bool

func init() {
	ch = make(chan uiInput, 5)
	interrupt = make(chan bool)
}

func (ui *gameui) Close() {
}

func (ui *gameui) Flush() {
	ui.DrawLogFrame()
	// very ugly optimisation
	xdgnmin := UIWidth - 1
	xdgnmax := 0
	ydgnmin := UIHeight - 1
	ydgnmax := 0
	xlogmin := UIWidth - 1
	xlogmax := 0
	ylogmin := UIHeight - 1
	ylogmax := 0
	xbarmin := UIWidth - 1
	xbarmax := 0
	ybarmin := UIHeight - 1
	ybarmax := 0
	for _, cdraw := range ui.g.DrawLog[len(ui.g.DrawLog)-1].Draws {
		cell := cdraw.Cell
		x, y := cdraw.X, cdraw.Y
		ui.Draw(cell, x, y)
		switch {
		case x < DungeonWidth && y < DungeonHeight:
			if x < xdgnmin {
				xdgnmin = x
			}
			if x > xdgnmax {
				xdgnmax = x
			}
			if y < ydgnmin {
				ydgnmin = y
			}
			if y > ydgnmax {
				ydgnmax = y
			}
		case x > DungeonWidth:
			if x < xbarmin {
				xbarmin = x
			}
			if x > xbarmax {
				xbarmax = x
			}
			if y < ybarmin {
				ybarmin = y
			}
			if y > ybarmax {
				ybarmax = y
			}
		default:
			if x < xlogmin {
				xlogmin = x
			}
			if x > xlogmax {
				xlogmax = x
			}
			if y < ylogmin {
				ylogmin = y
			}
			if y > ylogmax {
				ylogmax = y
			}
		}
	}
	ui.UpdateRectangle(xdgnmin, ydgnmin, xdgnmax, ydgnmax)
	ui.UpdateRectangle(xbarmin, ybarmin, xbarmax, ybarmax)
	ui.UpdateRectangle(xlogmin, ylogmin, xlogmax, ylogmax)
}

func (ui *gameui) UpdateRectangle(xmin, ymin, xmax, ymax int) {
	if xmin > xmax || ymin > ymax {
		return
	}
	pngbuf := &bytes.Buffer{}
	subimg := ui.canvas.SubImage(image.Rect(xmin*16, ymin*24, (xmax+1)*16, (ymax+1)*24))
	png.Encode(pngbuf, subimg)
	png := base64.StdEncoding.EncodeToString(pngbuf.Bytes())
	ui.ir.Eval("gamescreen put %{0%s} -format png -to %{1%d} %{2%d} %{3%d} %{4%d}", png,
		xmin*16, ymin*24, (xmax+1)*16, (ymax+1)*24) // TODO: optimize this more
}

func (ui *gameui) ApplyToggleLayout() {
	ui.ApplyToggleLayoutWithClear(true)
}

func (ui *gameui) ApplyToggleLayoutWithClear(clear bool) {
	GameConfig.Small = !GameConfig.Small
	if GameConfig.Small {
		ui.ir.Eval("wm geometry . =1280x576")
		if clear {
			ui.Clear()
			ui.Flush()
		}
		UIHeight = 24
		UIWidth = 80
	} else {
		ui.ir.Eval("wm geometry . =${width}x$height")
		UIHeight = 26
		UIWidth = 100
	}
	ui.cache = make(map[UICell]*image.RGBA)
	ui.g.DrawBuffer = make([]UICell, UIWidth*UIHeight)
	if clear {
		ui.Clear()
	}
}

func (ui *gameui) Draw(cell UICell, x, y int) {
	var img *image.RGBA
	if im, ok := ui.cache[cell]; ok {
		img = im
	} else {
		img = getImage(cell)
		ui.cache[cell] = img
	}
	draw.Draw(ui.canvas, image.Rect(x*ui.width, ui.height*y, (x+1)*ui.width, (y+1)*ui.height), img, image.Point{0, 0}, draw.Over)
}

func (ui *gameui) PollEvent() (in uiInput) {
	select {
	case in = <-ch:
	case in.interrupt = <-interrupt:
	}
	switch in.key {
	case "KP_Enter", "Return", "\r", "\n":
		in.key = "."
	case "Left", "KP_Left":
		in.key = "4"
	case "Right", "KP_Right":
		in.key = "6"
	case "Up", "KP_Up", "BackSpace":
		in.key = "8"
	case "Down", "KP_Down":
		in.key = "2"
	case "KP_Home":
		in.key = "7"
	case "KP_End":
		in.key = "1"
	case "KP_Prior", "Prior":
		in.key = "9"
	case "KP_Next", "Next":
		in.key = "3"
	case "KP_Begin", "KP_Delete":
		in.key = "5"
	default:
		if utf8.RuneCountInString(in.key) != 1 {
			in.key = ""
		}
	}
	return in
}
