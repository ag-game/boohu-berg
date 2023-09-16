// many ideas here from articles found at http://www.roguebasin.com/

package main

import (
	"sort"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
)

type dungeon struct {
	Gen   dungen
	Cells []cell
	PR    *paths.PathRange
}

type cell struct {
	T        terrain
	Explored bool
}

type terrain int

const (
	WallCell terrain = iota
	FreeCell
)

type dungen int

const (
	GenCaveMap dungen = iota
	GenRoomMap
	GenCellularAutomataCaveMap
	GenCaveMapTree
	GenRuinsMap
	GenBSPMap
)

func (dg dungen) Use(g *game) {
	switch dg {
	case GenCaveMap:
		g.GenCaveMap(DungeonHeight, DungeonWidth)
	case GenRoomMap:
		g.GenRoomMap(DungeonHeight, DungeonWidth)
	case GenCellularAutomataCaveMap:
		g.GenCellularAutomataCaveMap(DungeonHeight, DungeonWidth)
	case GenCaveMapTree:
		g.GenCaveMapTree(DungeonHeight, DungeonWidth)
	case GenRuinsMap:
		g.GenRuinsMap(DungeonHeight, DungeonWidth)
	case GenBSPMap:
		g.GenBSPMap(DungeonHeight, DungeonWidth)
	}
	g.Dungeon.Gen = dg
	g.Stats.DLayout[g.Depth] = dg.String()
}

func (dg dungen) String() (text string) {
	switch dg {
	case GenCaveMap:
		text = "OC"
	case GenRoomMap:
		text = "BR"
	case GenCellularAutomataCaveMap:
		text = "EC"
	case GenCaveMapTree:
		text = "TC"
	case GenRuinsMap:
		text = "RR"
	case GenBSPMap:
		text = "DT"
	}
	return text
}

func (dg dungen) Description() (text string) {
	switch dg {
	case GenCaveMap:
		text = "open cave"
	case GenRoomMap:
		text = "big rooms"
	case GenCellularAutomataCaveMap:
		text = "eight cave"
	case GenCaveMapTree:
		text = "tree-like cave"
	case GenRuinsMap:
		text = "ruined rooms"
	case GenBSPMap:
		text = "deserted town"
	}
	return text
}

type room struct {
	p gruid.Point
	w int
	h int
}

func (d *dungeon) Cell(p gruid.Point) cell {
	return d.Cells[idx(p)]
}

func (d *dungeon) Border(p gruid.Point) bool {
	return p.X == DungeonWidth-1 || p.Y == DungeonHeight-1 || p.X == 0 || p.Y == 0
}

func (d *dungeon) SetCell(p gruid.Point, t terrain) {
	d.Cells[idx(p)].T = t
}

func (d *dungeon) SetExplored(p gruid.Point) {
	d.Cells[idx(p)].Explored = true
}

func roomDistance(r1, r2 room) int {
	return Abs(r1.p.X-r2.p.X) + Abs(r1.p.Y-r2.p.Y)
}

func nearRoom(rooms []room, r room) room {
	closest := rooms[0]
	d := roomDistance(r, closest)
	for _, nextRoom := range rooms {
		nd := roomDistance(r, nextRoom)
		if nd < d {
			n := RandInt(10)
			if n > 3 {
				d = nd
				closest = nextRoom
			}
		}
	}
	return closest
}

func nearestRoom(rooms []room, r room) room {
	closest := rooms[0]
	d := roomDistance(r, closest)
	for _, nextRoom := range rooms {
		nd := roomDistance(r, nextRoom)
		if nd < d {
			n := RandInt(10)
			if n > 0 {
				d = nd
				closest = nextRoom
			}
		}
	}
	return closest
}

func intersectsRoom(rooms []room, r room) bool {
	for _, rr := range rooms {
		if (r.p.X+r.w-1 >= rr.p.X && rr.p.X+rr.w-1 >= r.p.X) &&
			(r.p.Y+r.h-1 >= rr.p.Y && rr.p.Y+rr.h-1 >= r.p.Y) {
			return true
		}
	}
	return false
}

func (d *dungeon) connectRooms(r1, r2 room) {
	x := r1.p.X
	if x < r2.p.X {
		x += r1.w - 1
	}
	y := r1.p.Y
	if y < r2.p.Y {
		y += r1.h - 1
	}
	d.SetCell(gruid.Point{x, y}, FreeCell)
	count := 0
	for {
		count++
		if count > 1000 {
			panic("ConnectRooms")
		}
		if x < r2.p.X {
			x++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if x > r2.p.X {
			x--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if y < r2.p.Y {
			y++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if y > r2.p.Y {
			y--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		break
	}
	d.SetCell(r2.p, FreeCell)
}

func (d *dungeon) connectRoomsDiagonally(r1, r2 room) {
	x := r1.p.X
	if x < r2.p.X {
		x += r1.w - 1
	}
	y := r1.p.Y
	if y < r2.p.Y {
		y += r1.h - 1
	}
	d.SetCell(gruid.Point{x, y}, FreeCell)
	count := 0
	for {
		count++
		if count > 1000 {
			panic("ConnectRooms")
		}
		if x < r2.p.X && y < r2.p.Y {
			x++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			y++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if x > r2.p.X && y < r2.p.Y {
			x--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			y++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if x > r2.p.X && y > r2.p.Y {
			x--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			y--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if x < r2.p.X && y > r2.p.Y {
			x++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			y--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if x < r2.p.X {
			x++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if x > r2.p.X {
			x--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if y < r2.p.Y {
			y++
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		if y > r2.p.Y {
			y--
			d.SetCell(gruid.Point{x, y}, FreeCell)
			continue
		}
		break
	}
	d.SetCell(r2.p, FreeCell)
}

func (d *dungeon) Area(area []gruid.Point, p gruid.Point, radius int) []gruid.Point {
	area = area[:0]
	for x := p.X - radius; x <= p.X+radius; x++ {
		for y := p.Y - radius; y <= p.Y+radius; y++ {
			q := gruid.Point{x, y}
			if valid(q) {
				area = append(area, q)
			}
		}
	}
	return area
}

func (d *dungeon) ConnectRoomsShortestPath(r1, r2 room) {
	var r1pos, r2pos gruid.Point
	r1pos.X = r1.p.X + RandInt(r1.w)
	if r1pos.X < r2.p.X {
		r1pos.X = r1.p.X + r1.w - 1
	}
	r1pos.Y = r1.p.Y + RandInt(r1.h)
	if r1pos.Y < r2.p.Y {
		r1pos.Y = r1.p.Y + r1.h - 1
	}
	r2pos.X = r2.p.X + RandInt(r2.w)
	if r2pos.X < r1.p.X {
		r2pos.X = r2.p.X + r2.w - 1
	}
	r2pos.Y = r2.p.Y + RandInt(r2.h)
	if r2pos.Y < r1.p.Y {
		r2pos.Y = r2.p.Y + r2.h - 1
	}
	mp := &dungeonPath{dungeon: d}
	path := d.PR.AstarPath(mp, r1pos, r2pos)
	for _, p := range path {
		d.SetCell(p, FreeCell)
	}
}

func (d *dungeon) ConnectIsolatedRoom(doorpos gruid.Point) {
	for i := 0; i < 200; i++ {
		p := d.FreeCell()
		dp := &dungeonPath{dungeon: d, wcost: unreachable}
		path := d.PR.AstarPath(dp, p, doorpos)
		wall := false
		for _, p := range path {
			if d.Cell(p).T == WallCell {
				wall = true
				break
			}
		}
		if !wall {
			continue
		}
		for _, p := range path {
			d.SetCell(p, FreeCell)
		}
		break
	}
}

func (d *dungeon) DigRoom(r room) {
	for i := r.p.X; i < r.p.X+r.w; i++ {
		for j := r.p.Y; j < r.p.Y+r.h; j++ {
			rpos := gruid.Point{i, j}
			if valid(rpos) {
				d.SetCell(rpos, FreeCell)
			}
		}
	}
}

func (d *dungeon) PutCols(r room) {
	for i := r.p.X + 1; i < r.p.X+r.w-1; i += 2 {
		for j := r.p.Y + 1; j < r.p.Y+r.h-1; j += 2 {
			rpos := gruid.Point{i, j}
			if valid(rpos) {
				d.SetCell(rpos, WallCell)
			}
		}
	}
}

func (d *dungeon) PutDiagCols(r room) {
	n := RandInt(2)
	for i := r.p.X + 1; i < r.p.X+r.w-1; i++ {
		m := n
		for j := r.p.Y + 1; j < r.p.Y+r.h-1; j++ {
			rpos := gruid.Point{i, j}
			if valid(rpos) && m%2 == 0 {
				d.SetCell(rpos, WallCell)
			}
			m++
		}
		n++
	}
}

func (d *dungeon) IsAreaFree(p gruid.Point, h, w int) bool {
	for i := p.X; i < p.X+w; i++ {
		for j := p.Y; j < p.Y+h; j++ {
			rpos := gruid.Point{i, j}
			if !valid(rpos) || d.Cell(rpos).T != FreeCell {
				return false
			}
		}
	}
	return true
}

func (d *dungeon) RoomDigCanditate(p gruid.Point, h, w int) (ret bool) {
	for i := p.X; i < p.X+w; i++ {
		for j := p.Y; j < p.Y+h; j++ {
			rpos := gruid.Point{i, j}
			if !valid(rpos) {
				return false
			}
			if d.Cell(rpos).T == FreeCell {
				ret = true
			}
		}
	}
	return ret
}

func (d *dungeon) IsolatedRoomDigCanditate(p gruid.Point, h, w int) (ret bool) {
	for i := p.X; i < p.X+w; i++ {
		for j := p.Y; j < p.Y+h; j++ {
			rpos := gruid.Point{i, j}
			if !valid(rpos) {
				return false
			}
			if d.Cell(rpos).T == FreeCell {
				return false
			}
		}
	}
	return true
}

func (d *dungeon) DigArea(p gruid.Point, h, w int) {
	for i := p.X; i < p.X+w; i++ {
		for j := p.Y; j < p.Y+h; j++ {
			rpos := gruid.Point{i, j}
			if !valid(rpos) {
				continue
			}
			d.SetCell(rpos, FreeCell)
		}
	}
}

func (d *dungeon) BlockArea(p gruid.Point, h, w int) {
	// not used now
	for i := p.X; i < p.X+w; i++ {
		for j := p.Y; j < p.Y+h; j++ {
			rpos := gruid.Point{i, j}
			if !valid(rpos) {
				continue
			}
			d.SetCell(rpos, WallCell)
		}
	}
}

func (d *dungeon) BuildRoom(p gruid.Point, w, h int, outside bool) map[gruid.Point]bool {
	spos := gruid.Point{p.X - 1, p.Y - 1}
	if outside && !d.IsAreaFree(spos, h+2, w+2) {
		return nil
	}
	for i := p.X; i < p.X+w; i++ {
		d.SetCell(gruid.Point{i, p.Y}, WallCell)
		d.SetCell(gruid.Point{i, p.Y + h - 1}, WallCell)
	}
	for i := p.Y; i < p.Y+h; i++ {
		d.SetCell(gruid.Point{p.X, i}, WallCell)
		d.SetCell(gruid.Point{p.X + w - 1, i}, WallCell)
	}
	if RandInt(2) == 0 || !outside {
		n := RandInt(2)
		for x := p.X + 1; x < p.X+w-1; x++ {
			m := n
			for y := p.Y + 1; y < p.Y+h-1; y++ {
				if m%2 == 0 {
					d.SetCell(gruid.Point{x, y}, WallCell)
				}
				m++
			}
			n++
		}
	} else {
		n := RandInt(2)
		m := RandInt(2)
		//if n == 0 && m == 0 {
		//// round room
		//d.SetCell(p, FreeCell)
		//d.SetCell(position{p.X, p.Y + h - 1}, FreeCell)
		//d.SetCell(position{p.X + w - 1, p.Y}, FreeCell)
		//d.SetCell(position{p.X + w - 1, p.Y + h - 1}, FreeCell)
		//}
		for x := p.X + 1 + m; x < p.X+w-1; x += 2 {
			for y := p.Y + 1 + n; y < p.Y+h-1; y += 2 {
				d.SetCell(gruid.Point{x, y}, WallCell)
			}
		}

	}
	area := make([]gruid.Point, 9)
	if outside {
		for _, p := range [4]gruid.Point{p, {p.X, p.Y + h - 1}, {p.X + w - 1, p.Y}, {p.X + w - 1, p.Y + h - 1}} {
			if d.WallAreaCount(area, p, 1) == 4 {
				d.SetCell(p, FreeCell)
			}
		}
	}
	doorsc := [4]gruid.Point{
		{p.X + w/2, p.Y},
		{p.X + w/2, p.Y + h - 1},
		{p.X, p.Y + h/2},
		{p.X + w - 1, p.Y + h/2},
	}
	doors := make(map[gruid.Point]bool)
	for i := 0; i < 3+RandInt(2); i++ {
		dpos := doorsc[RandInt(4)]
		doors[dpos] = true
		d.SetCell(dpos, FreeCell)
	}
	return doors
}

func (d *dungeon) BuildSomeRoom(w, h int) map[gruid.Point]bool {
	for i := 0; i < 200; i++ {
		p := d.FreeCell()
		doors := d.BuildRoom(p, w, h, true)
		if doors != nil {
			return doors
		}
	}
	return nil
}

func (d *dungeon) DigSomeRoom(w, h int) map[gruid.Point]bool {
	for i := 0; i < 200; i++ {
		p := d.FreeCell()
		dpos := gruid.Point{p.X - 1, p.Y - 1}
		if !d.RoomDigCanditate(dpos, h+2, w+2) {
			continue
		}
		d.DigArea(dpos, h+2, w+2)
		doors := d.BuildRoom(p, w, h, true)
		if doors != nil {
			return doors
		}
	}
	return nil
}

func (d *dungeon) DigIsolatedRoom(w, h int) map[gruid.Point]bool {
	i := RandInt(DungeonNCells)
	for j := 0; j < DungeonNCells; j++ {
		i = (i + 1) % DungeonNCells
		p := idx2Point(i)
		if d.Cells[i].T == FreeCell {
			continue
		}
		dpos := gruid.Point{p.X - 1, p.Y - 1}
		if !d.IsolatedRoomDigCanditate(dpos, h+2, w+2) {
			continue
		}
		d.DigArea(p, h, w)
		doors := d.BuildRoom(p, w, h, false)
		if doors != nil {
			return doors
		}
	}
	return nil
}

func (d *dungeon) ResizeRoom(r room) room {
	if DungeonWidth-r.p.X < r.w {
		r.w = DungeonWidth - r.p.X
	}
	if DungeonHeight-r.p.Y < r.h {
		r.h = DungeonHeight - r.p.Y
	}
	return r
}

func (g *game) GenRuinsMap(h, w int) {
	d := &dungeon{}
	d.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	d.Cells = make([]cell, h*w)
	rooms := []room{}
	for i := 0; i < 43; i++ {
		var ro room
		count := 100
		for count > 0 {
			count--
			ro = room{
				p: gruid.Point{RandInt(w - 1), RandInt(h - 1)},
				w: 3 + RandInt(5),
				h: 2 + RandInt(3)}
			ro = d.ResizeRoom(ro)
			if !intersectsRoom(rooms, ro) {
				break
			}
		}

		d.DigRoom(ro)
		if RandInt(60) == 0 {
			if RandInt(2) == 0 {
				d.PutCols(ro)
			} else {
				d.PutDiagCols(ro)
			}
		}
		if len(rooms) > 0 {
			r := RandInt(100)
			if r > 75 {
				d.connectRooms(nearRoom(rooms, ro), ro)
			} else if r > 25 {
				d.ConnectRoomsShortestPath(nearRoom(rooms, ro), ro)
			} else {
				d.connectRoomsDiagonally(nearRoom(rooms, ro), ro)
			}
		}
		rooms = append(rooms, ro)
	}
	doors := d.DigSomeRooms(5)
	g.Dungeon = d
	g.Fungus = make(map[gruid.Point]vegetation)
	g.DigFungus(1 + RandInt(2))
	g.PutDoors(30)
	g.PutDoorsList(doors, 20)
}

func (g *game) DigFungus(n int) {
	d := g.Dungeon
	count := 0
	fungus := g.Foliage(DungeonHeight, DungeonWidth)
loop:
	for i := 0; i < 100; i++ {
		if count > 100 {
			break loop
		}
		if n <= 0 {
			break
		}
		p := d.FreeCell()
		if _, ok := fungus[p]; ok {
			continue
		}
		conn, count := d.Connected(p, func(npos gruid.Point) bool {
			_, ok := fungus[npos]
			//return ok && d.IsFreeCell(npos)
			return ok
		})
		if count < 3 {
			continue
		}
		if len(conn) > 150 {
			continue
		}
		for cpos := range conn {
			d.SetCell(cpos, FreeCell)
			g.Fungus[cpos] = foliage
			count++
		}
		n--
	}
}

type roomSlice []room

func (rs roomSlice) Len() int      { return len(rs) }
func (rs roomSlice) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs roomSlice) Less(i, j int) bool {
	return rs[i].p.Y < rs[j].p.Y || rs[i].p.Y == rs[j].p.Y && rs[i].p.X < rs[j].p.X
}

func (g *game) GenRoomMap(h, w int) {
	d := &dungeon{}
	d.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	d.Cells = make([]cell, h*w)
	rooms := []room{}
	cols := 0
	for i := 0; i < 35; i++ {
		var ro room
		count := 100
		for count > 0 {
			count--
			ro = room{
				p: gruid.Point{RandInt(w - 1), RandInt(h - 1)},
				w: 5 + RandInt(4),
				h: 3 + RandInt(3)}
			ro = d.ResizeRoom(ro)
			if !intersectsRoom(rooms, ro) {
				break
			}
		}

		d.DigRoom(ro)
		if RandInt(10+15*cols) == 0 {
			if RandInt(2) == 0 {
				d.PutCols(ro)
			} else {
				d.PutDiagCols(ro)
			}
			cols++
		}
		rooms = append(rooms, ro)
	}
	sort.Sort(roomSlice(rooms))
	for i, ro := range rooms {
		if i == 0 {
			continue
		}
		r := RandInt(100)
		if r > 50 {
			d.connectRooms(nearestRoom(rooms[:i], ro), ro)
		} else if r > 25 {
			d.ConnectRoomsShortestPath(nearRoom(rooms[:i], ro), ro)
		} else {
			d.connectRoomsDiagonally(nearestRoom(rooms[:i], ro), ro)
		}
	}
	g.Dungeon = d
	doors := d.DigSomeRooms(5)
	g.PutDoors(90)
	g.PutDoorsList(doors, 10)
}

func (g *game) PutDoorsList(doors map[gruid.Point]bool, threshold int) {
	for p := range doors {
		if g.DoorCandidate(p) && RandInt(100) > threshold {
			g.Doors[p] = true
			delete(g.Fungus, p)
		}
	}
}

func (d *dungeon) FreeCell() gruid.Point {
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
		if c.T == FreeCell {
			return p
		}
	}
}

func (d *dungeon) WallCell() gruid.Point {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("WallCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if c.T == WallCell {
			return p
		}
	}
}

func (g *game) GenCaveMap(h, w int) {
	d := &dungeon{}
	d.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	d.Cells = make([]cell, h*w)
	p := gruid.Point{40, 10}
	max := 21 * 42
	d.SetCell(p, FreeCell)
	cells := 1
	notValid := 0
	lastValid := p
	diag := RandInt(4) == 0
	for cells < max {
		npos := RandomNeighbor(p, diag)
		if !valid(p) && valid(npos) && d.Cell(npos).T == WallCell {
			p = lastValid
			continue
		}
		p = npos
		if valid(p) {
			if d.Cell(p).T != FreeCell {
				d.SetCell(p, FreeCell)
				cells++
			}
			lastValid = p
		} else {
			notValid++
		}
		if notValid > 200 {
			notValid = 0
			p = lastValid
		}
	}
	cells = 1
	max = DungeonHeight * 1
	digs := 0
	i := 0
	block := make([]gruid.Point, 0, 64)
loop:
	for cells < max {
		i++
		if digs > 3 {
			break
		}
		if i > 1000 {
			break
		}
		diag = RandInt(2) == 0
		block = d.DigBlock(block, diag)
		if len(block) == 0 {
			continue loop
		}
		if len(block) < 4 || len(block) > 10 {
			continue loop
		}
		for _, p := range block {
			d.SetCell(p, FreeCell)
			cells++
		}
		digs++
	}
	doors := make(map[gruid.Point]bool)
	rooms := 0
	if RandInt(4) > 0 {
		w, h := GenCaveRoomSize()
		rooms++
		for p := range d.BuildSomeRoom(w, h) {
			doors[p] = true
		}
		if RandInt(7) == 0 {
			rooms++
			w, h := GenCaveRoomSize()
			for p := range d.BuildSomeRoom(w, h) {
				doors[p] = true
			}

		}
	}
	if RandInt(1+rooms) == 0 {
		w, h := GenLittleRoomSize()
		i := 0
		for p := range d.DigIsolatedRoom(w, h) {
			doors[p] = true
			if i == 0 {
				d.ConnectIsolatedRoom(p)
			}
			i++
		}

	}
	g.Dungeon = d
	g.Fungus = g.Foliage(DungeonHeight, DungeonWidth)
	g.PutDoors(5)
	for p := range doors {
		if g.DoorCandidate(p) && RandInt(100) > 20 {
			g.Doors[p] = true
			delete(g.Fungus, p)
		}
	}
}

func GenCaveRoomSize() (int, int) {
	return 7 + 2*RandInt(2), 5 + 2*RandInt(2)
}

func GenLittleRoomSize() (int, int) {
	return 7, 5
}

func (d *dungeon) HasFreeNeighbor(p gruid.Point) bool {
	neighbors := ValidNeighbors(p)
	for _, p := range neighbors {
		if d.Cell(p).T == FreeCell {
			return true
		}
	}
	return false
}

func (g *game) HasFreeExploredNeighbor(p gruid.Point) bool {
	d := g.Dungeon
	neighbors := ValidNeighbors(p)
	for _, p := range neighbors {
		c := d.Cell(p)
		if c.T == FreeCell && c.Explored && !g.WrongWall[p] {
			return true
		}
	}
	return false
}

func (d *dungeon) DigBlock(block []gruid.Point, diag bool) []gruid.Point {
	p := d.WallCell()
	block = block[:0]
	for {
		block = append(block, p)
		if d.HasFreeNeighbor(p) {
			break
		}
		p = RandomNeighbor(p, diag)
		if !valid(p) {
			block = block[:0]
			p = d.WallCell()
			continue
		}
		if !valid(p) {
			return nil
		}
	}
	return block
}

func (g *game) GenCaveMapTree(h, w int) {
	d := &dungeon{}
	d.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	d.Cells = make([]cell, h*w)
	center := gruid.Point{40, 10}
	d.SetCell(center, FreeCell)
	d.SetCell(center.Shift(1, 0), FreeCell)
	d.SetCell(center.Shift(1, -1), FreeCell)
	d.SetCell(center.Shift(0, 1), FreeCell)
	d.SetCell(center.Shift(1, 1), FreeCell)
	d.SetCell(center.Shift(0, -1), FreeCell)
	d.SetCell(center.Shift(-1, -1), FreeCell)
	d.SetCell(center.Shift(-1, 0), FreeCell)
	d.SetCell(center.Shift(-1, 1), FreeCell)
	max := 21 * 23
	cells := 1
	diag := RandInt(2) == 0
	block := make([]gruid.Point, 0, 64)
loop:
	for cells < max {
		block = d.DigBlock(block, diag)
		if len(block) == 0 {
			continue loop
		}
		for _, p := range block {
			if d.Cell(p).T != FreeCell {
				d.SetCell(p, FreeCell)
				cells++
			}
		}
	}

	doors := d.DigSomeRooms(5)
	g.Dungeon = d
	g.Fungus = make(map[gruid.Point]vegetation)
	g.DigFungus(1 + RandInt(2))
	g.PutDoors(5)
	g.PutDoorsList(doors, 20)
}

func (d *dungeon) DigSomeRooms(chances int) map[gruid.Point]bool {
	doors := make(map[gruid.Point]bool)
	if RandInt(chances) > 0 {
		w, h := GenCaveRoomSize()
		for p := range d.DigSomeRoom(w, h) {
			doors[p] = true
		}
		if RandInt(3) == 0 {
			w, h := GenCaveRoomSize()
			for p := range d.DigSomeRoom(w, h) {
				doors[p] = true
			}
		}
	}
	return doors
}

func (d *dungeon) WallAreaCount(area []gruid.Point, p gruid.Point, radius int) int {
	area = d.Area(area, p, radius)
	count := 0
	for _, npos := range area {
		if d.Cell(npos).T == WallCell {
			count++
		}
	}
	switch radius {
	case 1:
		count += 9 - len(area)
	case 2:
		count += 25 - len(area)
	}
	return count
}

func (d *dungeon) Connected(p gruid.Point, nf func(gruid.Point) bool) (map[gruid.Point]bool, int) {
	conn := map[gruid.Point]bool{}
	stack := []gruid.Point{p}
	count := 0
	conn[p] = true
	nb := make([]gruid.Point, 0, 8)
	for len(stack) > 0 {
		p = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		nb = Neighbors(p, nb, nf)
		for _, npos := range nb {
			if !conn[npos] {
				conn[npos] = true
				stack = append(stack, npos)
			}
		}
	}
	return conn, count
}

func (d *dungeon) connex() bool {
	p := d.FreeCell()
	conn, _ := d.Connected(p, d.IsFreeCell)
	for i, c := range d.Cells {
		if c.T == FreeCell && !conn[idx2Point(i)] {
			return false
		}
	}
	return true
}

func (g *game) RunCellularAutomataCave(h, w int) bool {
	d := &dungeon{}
	d.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	d.Cells = make([]cell, h*w)
	for i := range d.Cells {
		r := RandInt(100)
		p := idx2Point(i)
		if r >= 45 {
			d.SetCell(p, FreeCell)
		} else {
			d.SetCell(p, WallCell)
		}
	}
	bufm := &dungeon{}
	bufm.Cells = make([]cell, h*w)
	area := make([]gruid.Point, 0, 25)
	for i := 0; i < 5; i++ {
		for j := range bufm.Cells {
			p := idx2Point(j)
			c1 := d.WallAreaCount(area, p, 1)
			if c1 >= 5 {
				bufm.SetCell(p, WallCell)
			} else {
				bufm.SetCell(p, FreeCell)
			}
			if i == 3 {
				c2 := d.WallAreaCount(area, p, 2)
				if c2 <= 2 {
					bufm.SetCell(p, WallCell)
				}
			}
		}
		copy(d.Cells, bufm.Cells)
	}
	var conn map[gruid.Point]bool
	var count int
	var winner gruid.Point
	for i := 0; i < 15; i++ {
		p := d.FreeCell()
		if conn[p] {
			continue
		}
		var ncount int
		conn, ncount = d.Connected(p, d.IsFreeCell)
		if ncount > count {
			count = ncount
			winner = p
		}
		if count >= 37*DungeonHeight*DungeonWidth/100 {
			break
		}
	}
	conn, count = d.Connected(winner, d.IsFreeCell)
	if count <= 37*DungeonHeight*DungeonWidth/100 {
		return false
	}
	for i, c := range d.Cells {
		p := idx2Point(i)
		if c.T == FreeCell && !conn[p] {
			d.SetCell(p, WallCell)
		}
	}
	digs := 0
	max := DungeonHeight / 2
	cells := 1
	i := 0
	block := make([]gruid.Point, 0, 64)
loop:
	for cells < max {
		i++
		if digs > 3 {
			break
		}
		if i > 1000 {
			break
		}
		diag := RandInt(2) == 0
		block = d.DigBlock(block, diag)
		if len(block) == 0 {
			continue loop
		}
		if len(block) < 4 || len(block) > 10 {
			continue loop
		}
		for _, p := range block {
			d.SetCell(p, FreeCell)
			cells++
		}
		digs++
	}
	doors := make(map[gruid.Point]bool)
	if RandInt(5) > 0 {
		w, h := GenLittleRoomSize()
		i := 0
		for p := range d.DigIsolatedRoom(w, h) {
			doors[p] = true
			if i == 0 {
				d.ConnectIsolatedRoom(p)
			}
			i++
		}
		if RandInt(4) == 0 {
			w, h := GenCaveRoomSize()
			i := 0
			for p := range d.DigIsolatedRoom(w, h) {
				doors[p] = true
				if i == 0 {
					d.ConnectIsolatedRoom(p)
				}
				i++
			}
		}
	}
	g.Dungeon = d
	g.PutDoors(10)
	for p := range doors {
		if g.DoorCandidate(p) && RandInt(100) > 20 {
			g.Doors[p] = true
			delete(g.Fungus, p)
		}
	}
	return true
}

func (g *game) GenCellularAutomataCaveMap(h, w int) {
	count := 0
	for {
		count++
		if count > 100 {
			panic("genCellularAutomataCaveMap")
		}
		if g.RunCellularAutomataCave(h, w) {
			break
		}
	}
	g.Fungus = g.Foliage(DungeonHeight, DungeonWidth)
}

func (d *dungeon) SimpleRoom(r room) map[gruid.Point]bool {
	for i := r.p.X; i < r.p.X+r.w; i++ {
		d.SetCell(gruid.Point{i, r.p.Y}, WallCell)
		d.SetCell(gruid.Point{i, r.p.Y + r.h - 1}, WallCell)
	}
	for i := r.p.Y; i < r.p.Y+r.h; i++ {
		d.SetCell(gruid.Point{r.p.X, i}, WallCell)
		d.SetCell(gruid.Point{r.p.X + r.w - 1, i}, WallCell)
	}
	doorsc := [4]gruid.Point{
		{r.p.X + r.w/2, r.p.Y},
		{r.p.X + r.w/2, r.p.Y + r.h - 1},
		{r.p.X, r.p.Y + r.h/2},
		{r.p.X + r.w - 1, r.p.Y + r.h/2},
	}
	doors := make(map[gruid.Point]bool)
	for i := 0; i < 3+RandInt(2); i++ {
		dpos := doorsc[RandInt(4)]
		doors[dpos] = true
		d.SetCell(dpos, FreeCell)
	}
	return doors
}

func (g *game) ExtendEdgeRoom(r room, doors map[gruid.Point]bool) room {
	if g.Dungeon.Cell(r.p).T != WallCell {
		return r
	}
	extend := false
	if r.p.X+r.w+1 == DungeonWidth {
		for i := r.p.Y + 1; i < r.p.Y+r.h-1; i++ {
			g.Dungeon.SetCell(gruid.Point{DungeonWidth - 2, i}, FreeCell)
			g.Dungeon.SetCell(gruid.Point{DungeonWidth - 1, i}, WallCell)
		}
		g.Dungeon.SetCell(gruid.Point{DungeonWidth - 1, r.p.Y}, WallCell)
		g.Dungeon.SetCell(gruid.Point{DungeonWidth - 1, r.p.Y + r.h - 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{DungeonWidth - 2, r.p.Y + 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{DungeonWidth - 2, r.p.Y + r.h - 2}, WallCell)
		r.w++
		extend = true
	}
	if r.p.X == 1 {
		for i := r.p.Y + 1; i < r.p.Y+r.h-1; i++ {
			g.Dungeon.SetCell(gruid.Point{1, i}, FreeCell)
			g.Dungeon.SetCell(gruid.Point{0, i}, WallCell)
		}
		g.Dungeon.SetCell(gruid.Point{0, r.p.Y}, WallCell)
		g.Dungeon.SetCell(gruid.Point{0, r.p.Y + r.h - 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{1, r.p.Y + 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{1, r.p.Y + r.h - 2}, WallCell)
		r.w++
		r.p.X--
		extend = true
	}
	if r.p.Y+r.h+1 == DungeonHeight {
		for i := r.p.X + 1; i < r.p.X+r.w-1; i++ {
			g.Dungeon.SetCell(gruid.Point{i, DungeonHeight - 2}, FreeCell)
			g.Dungeon.SetCell(gruid.Point{i, DungeonHeight - 1}, WallCell)
		}
		g.Dungeon.SetCell(gruid.Point{r.p.X, DungeonHeight - 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{r.p.X + r.w - 1, DungeonHeight - 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{r.p.X + 1, DungeonHeight - 2}, WallCell)
		g.Dungeon.SetCell(gruid.Point{r.p.X + r.w - 2, DungeonHeight - 2}, WallCell)
		r.h++
		extend = true
	}
	if r.p.Y == 1 {
		for i := r.p.X + 1; i < r.p.X+r.w-1; i++ {
			g.Dungeon.SetCell(gruid.Point{i, 1}, FreeCell)
			g.Dungeon.SetCell(gruid.Point{i, 0}, WallCell)
		}
		g.Dungeon.SetCell(gruid.Point{r.p.X, 0}, WallCell)
		g.Dungeon.SetCell(gruid.Point{r.p.X + r.w - 1, 0}, WallCell)
		g.Dungeon.SetCell(gruid.Point{r.p.X + 1, 1}, WallCell)
		g.Dungeon.SetCell(gruid.Point{r.p.X + r.w - 2, 1}, WallCell)
		r.h++
		r.p.Y--
		extend = true
	}
	if !extend {
		return r
	}
	for p := range doors {
		if p.X == 1 || p.X == DungeonWidth-2 || p.Y == 1 || p.Y == DungeonHeight-2 {
			delete(g.Doors, p)
			continue
		}
	}
	doorsc := [4]gruid.Point{
		{r.p.X + r.w/2, r.p.Y},
		{r.p.X + r.w/2, r.p.Y + r.h - 1},
		{r.p.X, r.p.Y + r.h/2},
		{r.p.X + r.w - 1, r.p.Y + r.h/2},
	}
	ndoorsc := []gruid.Point{}
	ndoors := 0
	for _, p := range doorsc {
		if p.X == 0 || p.X == DungeonWidth-1 || p.Y == 0 || p.Y == DungeonHeight-1 {
			continue
		}
		if g.Doors[p] {
			ndoors++
		}
		ndoorsc = append(ndoorsc, p)
	}
	for i := 0; i < 1+RandInt(2-ndoors); i++ {
		dpos := ndoorsc[RandInt(len(ndoorsc))]
		g.Doors[dpos] = true
		g.Dungeon.SetCell(dpos, FreeCell)
	}
	return r
}

func (g *game) DivideRoomVertically(r room) {
	if g.Dungeon.Cell(r.p).T != WallCell {
		return
	}
	if r.w <= 6 {
		return
	}
	if r.h < 5 {
		return
	}
	dx := 2 + RandInt(r.w/2-2)
	if RandInt(2) == 0 {
		dx = r.w - 3 - RandInt(r.w/2-3)
	}
	if dx == 2 && r.p.X == 0 {
		return
	}
	if dx == r.w-3 && r.p.X+r.w == DungeonWidth {
		return
	}
	free := true
loop:
	for i := r.p.Y + 1; i < r.p.Y+r.h-1; i++ {
		for j := dx - 1; j <= dx+1; j++ {
			if g.Dungeon.Cell(gruid.Point{r.p.X + j, i}).T == WallCell {
				free = false
				break loop
			}
		}
	}
	if !free {
		return
	}
	for i := r.p.Y + 1; i < r.p.Y+r.h-1; i++ {
		g.Dungeon.SetCell(gruid.Point{r.p.X + dx, i}, WallCell)
	}
	doorpos := gruid.Point{r.p.X + dx, r.p.Y + r.h/2}
	g.Doors[doorpos] = true
	g.Dungeon.SetCell(doorpos, FreeCell)
}

func (g *game) DivideRoomHorizontally(r room) {
	if g.Dungeon.Cell(r.p).T != WallCell {
		return
	}
	if r.h <= 6 {
		return
	}
	if r.w < 5 {
		return
	}
	dy := 2 + RandInt(r.h/2-2)
	if RandInt(2) == 0 {
		dy = r.h - 3 - RandInt(r.h/2-3)
	}
	if dy == 2 && r.p.Y == 0 {
		return
	}
	if dy == r.h-3 && r.p.Y+r.h == DungeonHeight {
		return
	}
	free := true
loop:
	for i := r.p.X + 1; i < r.p.X+r.w-1; i++ {
		for j := dy - 1; j <= dy+1; j++ {
			if g.Dungeon.Cell(gruid.Point{i, r.p.Y + j}).T == WallCell {
				free = false
				break loop
			}
		}
	}
	if !free {
		return
	}
	for i := r.p.X + 1; i < r.p.X+r.w-1; i++ {
		g.Dungeon.SetCell(gruid.Point{i, r.p.Y + dy}, WallCell)
	}
	doorpos := gruid.Point{r.p.X + r.w/2, r.p.Y + dy}
	g.Doors[doorpos] = true
	g.Dungeon.SetCell(doorpos, FreeCell)
}

func (g *game) GenBSPMap(height, width int) {
	rooms := []room{}
	crooms := []room{{p: gruid.Point{1, 1}, w: DungeonWidth - 2, h: DungeonHeight - 2}}
	big := 0
	for len(crooms) > 0 {
		r := crooms[0]
		crooms = crooms[1:]
		if r.h <= 8 && r.w <= 12 {
			switch RandInt(6) {
			case 0:
				if r.h >= 6 {
					r.h--
					if RandInt(2) == 0 {
						r.p.Y++
					}
				}
			case 1:
				if r.w >= 8 {
					r.w--
					if RandInt(2) == 0 {
						r.p.X++
					}
				}
			}
			if r.h > 2 && r.w > 2 {
				rooms = append(rooms, r)
			}
			continue
		}
		if RandInt(2+big) == 0 && (r.h <= 12 && r.w <= 20) {
			big++
			switch RandInt(4) {
			case 0:
				r.h--
				if RandInt(2) == 0 {
					r.p.Y++
				}
			case 1:
				r.w--
				if RandInt(2) == 0 {
					r.p.X++
				}
			}
			if r.h > 2 && r.w > 2 {
				rooms = append(rooms, r)
			}
			continue
		}
		horizontal := false
		if r.h > 8 && r.w > 10 && r.w < 40 && RandInt(4) == 0 {
			horizontal = true
		} else if r.h > 8 && r.w <= 10+RandInt(5) {
			horizontal = true
		}
		if horizontal {
			h := r.h/2 - r.h/4 + RandInt(1+r.h/2)
			if h <= 3 {
				h++
			}
			if r.h-h-1 <= 3 {
				h--
			}
			crooms = append(crooms, room{r.p, r.w, h}, room{gruid.Point{r.p.X, r.p.Y + 1 + h}, r.w, r.h - h - 1})
		} else {
			w := r.w/2 - r.w/4 + RandInt(1+r.w/2)
			if w <= 3 {
				w++
			}
			if r.w-w-1 <= 3 {
				w--
			}
			crooms = append(crooms, room{r.p, w, r.h}, room{gruid.Point{r.p.X + 1 + w, r.p.Y}, r.w - w - 1, r.h})
		}
	}

	d := &dungeon{}
	d.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	d.Cells = make([]cell, height*width)
	for i := 0; i < DungeonNCells; i++ {
		d.SetCell(idx2Point(i), FreeCell)
	}
	g.Dungeon = d
	g.Doors = map[gruid.Point]bool{}
	special := 0
	empty := 0
	for i, r := range rooms {
		var doors map[gruid.Point]bool
		if RandInt(2+special/3) == 0 && r.w%2 == 1 && r.h%2 == 1 && r.w >= 5 && r.h >= 5 {
			doors = d.BuildRoom(r.p, r.w, r.h, true)
			special++
		} else if empty > 0 || RandInt(20) > 0 {
			doors = d.SimpleRoom(r)
			if RandInt(2) == 0 && r.w >= 7 && r.h >= 7 {
				rn := r
				rn.p.X++
				rn.p.Y++
				rn.h--
				rn.h--
				rn.w--
				rn.w--
				if RandInt(2) == 0 {
					d.PutCols(rn)
				} else {
					d.PutDiagCols(rn)
				}
			} else if RandInt(1+special/2) == 0 && r.w >= 11 && r.h >= 9 {
				sx := (r.w - 11) / 2
				sy := (r.h - 9) / 2
				doors = d.BuildRoom(gruid.Point{r.p.X + 2 + sx, r.p.Y + 2 + sy}, 7, 5, true)
				special++
			}
		} else {
			empty++
		}
		for p := range doors {
			if g.DoorCandidate(p) && RandInt(100) > 10 {
				g.Doors[p] = true
			}
		}
		if RandInt(2) == 0 {
			r = g.ExtendEdgeRoom(r, doors)
			rooms[i] = r
		}
		if RandInt(5) > 0 {
			if RandInt(2) == 0 {
				g.DivideRoomVertically(r)
			} else {
				g.DivideRoomHorizontally(r)
			}
		}
	}
	g.Fungus = make(map[gruid.Point]vegetation)
	g.DigFungus(RandInt(3))
	for i := 0; i <= RandInt(2); i++ {
		r := rooms[RandInt(len(rooms))]
		for x := r.p.X + 1; x < r.p.X+r.w-1; x++ {
			for y := r.p.Y + 1; y < r.p.Y+r.h-1; y++ {
				g.Fungus[gruid.Point{x, y}] = foliage
			}
		}
	}
}

type vegetation int

const (
	foliage vegetation = iota
)

func (g *game) Foliage(h, w int) map[gruid.Point]vegetation {
	// use same structure as for the dungeon
	// walls will become foliage
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	for i := range d.Cells {
		r := RandInt(100)
		p := idx2Point(i)
		if r >= 43 {
			d.SetCell(p, WallCell)
		} else {
			d.SetCell(p, FreeCell)
		}
	}
	area := make([]gruid.Point, 0, 25)
	for i := 0; i < 6; i++ {
		bufm := &dungeon{}
		bufm.Cells = make([]cell, h*w)
		copy(bufm.Cells, d.Cells)
		for j := range bufm.Cells {
			p := idx2Point(j)
			c1 := d.WallAreaCount(area, p, 1)
			if i < 4 {
				if c1 <= 4 {
					bufm.SetCell(p, FreeCell)
				} else {
					bufm.SetCell(p, WallCell)
				}
			}
			if i == 4 {
				if c1 > 6 {
					bufm.SetCell(p, WallCell)
				}
			}
			if i == 5 {
				c2 := d.WallAreaCount(area, p, 2)
				if c2 < 5 && c1 <= 2 {
					bufm.SetCell(p, FreeCell)
				}
			}
		}
		d.Cells = bufm.Cells
	}
	fungus := make(map[gruid.Point]vegetation)
	for i, c := range d.Cells {
		if _, ok := g.Doors[idx2Point(i)]; !ok && c.T == FreeCell {
			fungus[idx2Point(i)] = foliage
		}
	}
	return fungus
}

func (g *game) DoorCandidate(p gruid.Point) bool {
	d := g.Dungeon
	if !valid(p) || d.Cell(p).T != FreeCell {
		return false
	}
	return valid(p.Shift(-1, 0)) && valid(p.Shift(1, 0)) &&
		d.Cell(p.Shift(-1, 0)).T == FreeCell && d.Cell(p.Shift(1, 0)).T == FreeCell &&
		!g.Doors[p.Shift(-1, 0)] && !g.Doors[p.Shift(1, 0)] &&
		(!valid(p.Shift(0, -1)) || d.Cell(p.Shift(0, -1)).T == WallCell) &&
		(!valid(p.Shift(0, 1)) || d.Cell(p.Shift(0, 1)).T == WallCell) &&
		((valid(p.Shift(-1, -1)) && d.Cell(p.Shift(-1, -1)).T == FreeCell) ||
			(valid(p.Shift(-1, 1)) && d.Cell(p.Shift(-1, 1)).T == FreeCell) ||
			(valid(p.Shift(1, -1)) && d.Cell(p.Shift(1, -1)).T == FreeCell) ||
			(valid(p.Shift(1, 1)) && d.Cell(p.Shift(1, 1)).T == FreeCell)) ||
		valid(p.Shift(0, -1)) && valid(p.Shift(0, 1)) &&
			d.Cell(p.Shift(0, -1)).T == FreeCell && d.Cell(p.Shift(0, 1)).T == FreeCell &&
			!g.Doors[p.Shift(0, -1)] && !g.Doors[p.Shift(0, 1)] &&
			(!valid(p.Shift(1, 0)) || d.Cell(p.Shift(1, 0)).T == WallCell) &&
			(!valid(p.Shift(-1, 0)) || d.Cell(p.Shift(-1, 0)).T == WallCell) &&
			((valid(p.Shift(-1, -1)) && d.Cell(p.Shift(-1, -1)).T == FreeCell) ||
				(valid(p.Shift(-1, 1)) && d.Cell(p.Shift(-1, 1)).T == FreeCell) ||
				(valid(p.Shift(1, -1)) && d.Cell(p.Shift(1, -1)).T == FreeCell) ||
				(valid(p.Shift(1, 1)) && d.Cell(p.Shift(1, 1)).T == FreeCell))
}

func (g *game) PutDoors(percentage int) {
	g.Doors = map[gruid.Point]bool{}
	for i := range g.Dungeon.Cells {
		p := idx2Point(i)
		if g.DoorCandidate(p) && RandInt(100) < percentage {
			g.Doors[p] = true
			delete(g.Fungus, p)
		}
	}
}
