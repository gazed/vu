// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"sort"
	"vu"
	"vu/math/lin"
)

// sg tests the scene graph by building up a movable character that has multiple
// layers of sub-parts.  The scene graph works when changing the top level location,
// orientation, and scale also affects the subparts. Sg also tests adding and removing
// parts from a scene graph. The main classes being tested are vu.Scene and vu.Part, eg:
//	    scene.AddPart()
//	    scene.RemPart(sg.tr.part)
//
// Note that the movement keys move the character and not the camera in this example.
func sg() {
	sg := &sgtag{}
	var err error
	if sg.eng, err = vu.New("Player character", 400, 100, 800, 600); err != nil {
		log.Printf("sg: error intitializing engine %s", err)
		return
	}
	sg.run = 10            // move so many cubes worth in one second.
	sg.spin = 270          // spin so many degrees in one second.
	sg.eng.SetDirector(sg) // override user input handling.
	sg.stagePlay()
	defer sg.eng.Shutdown()
	sg.eng.Action()
}

// Globally unique "tag" for this example.
type sgtag struct {
	eng    *vu.Eng
	scene  vu.Scene
	tr     *trooper
	reacts map[string]vu.Reaction
	run    float32
	spin   float32
}

func (sg *sgtag) stagePlay() {
	sg.scene = sg.eng.AddScene(vu.VP)
	sg.scene.SetPerspective(60, float32(800)/float32(600), 0.1, 50)
	sg.scene.SetLightLocation(0, 10, 0)
	sg.scene.SetLightColour(0.4, 0.7, 0.9)
	sg.scene.SetViewLocation(0, 0, 6)

	// load the floor model.
	floor := sg.scene.AddPart()
	floor.SetLocation(0, 0, 0)
	floor.SetFacade("floor", "gouraud", "floor")

	// load the trooper
	sg.tr = newTrooper(sg.eng, sg.scene, 1)

	// initialize the reactions
	sg.reacts = map[string]vu.Reaction{
		"W":   vu.NewReaction("Forward", func() { sg.forward() }),
		"A":   vu.NewReaction("Left", func() { sg.left() }),
		"S":   vu.NewReaction("Back", func() { sg.back() }),
		"D":   vu.NewReaction("Right", func() { sg.right() }),
		"KP+": vu.NewReaction("Attach", func() { sg.tr.attach() }),
		"KP-": vu.NewReaction("Detach", func() { sg.tr.detach() }),
		"0":   vu.NewReactOnce("SL0", func() { sg.setTr(0) }),
		"1":   vu.NewReactOnce("SL1", func() { sg.setTr(1) }),
		"2":   vu.NewReactOnce("SL2", func() { sg.setTr(2) }),
		"3":   vu.NewReactOnce("SL3", func() { sg.setTr(3) }),
		"4":   vu.NewReactOnce("SL4", func() { sg.setTr(4) }),
		"5":   vu.NewReactOnce("SL5", func() { sg.setTr(5) }),
		"P":   vu.NewReactOnce("Stats", func() { sg.stats() }),
	}

	// set some constant state.
	sg.eng.Enable(vu.BLEND, true)
	sg.eng.Enable(vu.CULL, true)
	sg.eng.Enable(vu.DEPTH, true)
	sg.eng.Color(0.1, 0.1, 0.1, 1.0)
	return
}

// Handle engine callbacks.
func (sg *sgtag) Focus(focus bool) {}
func (sg *sgtag) Resize(x, y, width, height int) {
	sg.eng.ResizeViewport(x, y, width, height)
	ratio := float32(width) / float32(height)
	sg.scene.SetPerspective(60, ratio, 0.1, 50)
}
func (sg *sgtag) Update(pressed []string, gt, dt float32) {
	for _, p := range pressed {
		if reaction, ok := sg.reacts[p]; ok {
			reaction.Do()
		}
	}
}

// User actions.
func (sg *sgtag) stats() { log.Printf("Health %d", sg.tr.health()) }
func (sg *sgtag) setTr(lvl int) {
	sg.scene.RemPart(sg.tr.part)
	sg.tr.trash()
	sg.tr = newTrooper(sg.eng, sg.scene, lvl)
}
func (sg *sgtag) left() {
	sg.tr.part.RotateY(sg.eng.Dt * sg.spin)
}
func (sg *sgtag) right() {
	sg.tr.part.RotateY(sg.eng.Dt * -sg.spin)
}
func (sg *sgtag) back() {
	sg.tr.part.Move(0, 0, sg.eng.Dt*sg.run)
}
func (sg *sgtag) forward() {
	sg.tr.part.Move(0, 0, sg.eng.Dt*-sg.run)
}

// trooper is an attempt to keep polygon growth linear while the player
// statistics grows exponentially. A trooper is rendered using a single mesh
// that is replicated 1 or more times depending on the health of the trooper.
type trooper struct {
	part   vu.Part
	lvl    int
	eng    *vu.Eng
	neo    vu.Part // un-injured trooper
	bits   []box   // injured troopers have panels and edge cubes.
	center vu.Part // center always represented as one piece
	mid    int     // level entry number of cells.
}

// newTrooper creates a trooper at the starting size for the given level.
//    level 0: 1x1x1 :  1 cube
//    level 1: 2x2x2 :  8 edge cubes + 6 panels of 0x0 cubes + 0x0x0 center.
//    level 2: 3x3x3 : 20 edge cubes + 6 panels of 1x1 cubes + 1x1x1 center.
//    level 3: 4x4x4 : 32 edge cubes + 6 panels of 2x2 cubes + 2x2x2 center.
//    ...
func newTrooper(eng *vu.Eng, scene vu.Scene, level int) *trooper {
	tr := &trooper{}
	tr.lvl = level
	tr.eng = eng
	tr.bits = []box{}
	tr.mid = tr.lvl*tr.lvl*tr.lvl*8 - (tr.lvl-1)*(tr.lvl-1)*(tr.lvl-1)*8
	tr.part = scene.AddPart()

	//
	if tr.lvl == 0 {
		cube := newCube(eng, tr.part, 0, 0, 0, 1)
		cube.edgeSort(1)
		tr.bits = append(tr.bits, cube)
		return tr
	}

	// create the panels. These are used in each level but the first.
	cubeSize := 1.0 / float32(tr.lvl+1)
	centerOffset := cubeSize * 0.5
	panelCenter := float32(tr.lvl) * centerOffset
	tr.bits = append(tr.bits, newPanel(eng, tr.part, panelCenter, 0.0, 0.0, tr.lvl))
	tr.bits = append(tr.bits, newPanel(eng, tr.part, -panelCenter, 0.0, 0.0, tr.lvl))
	tr.bits = append(tr.bits, newPanel(eng, tr.part, 0.0, panelCenter, 0.0, tr.lvl))
	tr.bits = append(tr.bits, newPanel(eng, tr.part, 0.0, -panelCenter, 0.0, tr.lvl))
	tr.bits = append(tr.bits, newPanel(eng, tr.part, 0.0, 0.0, panelCenter, tr.lvl))
	tr.bits = append(tr.bits, newPanel(eng, tr.part, 0.0, 0.0, -panelCenter, tr.lvl))

	// troopers are made out of cubes and panels.
	mx := float32(-tr.lvl)
	for cx := 0; cx <= tr.lvl; cx++ {
		my := float32(-tr.lvl)
		for cy := 0; cy <= tr.lvl; cy++ {
			mz := float32(-tr.lvl)
			for cz := 0; cz <= tr.lvl; cz++ {

				// create the outer edges.
				newCells := 0
				if (cx == 0 || cx == tr.lvl) && (cy == 0 || cy == tr.lvl) && (cz == 0 || cz == tr.lvl) {

					// corner cube
					newCells = 1
				} else if (cx == 0 || cx == tr.lvl) && (cy == 0 || cy == tr.lvl) ||
					(cx == 0 || cx == tr.lvl) && (cz == 0 || cz == tr.lvl) ||
					(cy == 0 || cy == tr.lvl) && (cz == 0 || cz == tr.lvl) {

					// edge cube
					newCells = 2
				} else if cx == 0 || cx == tr.lvl || cy == 0 || cy == tr.lvl || cz == 0 || cz == tr.lvl {

					// side cubes are added to (controlled by) a panel.
					x, y, z := mx*centerOffset, my*centerOffset, mz*centerOffset
					if cx == tr.lvl && x > y && x > z {
						tr.bits[0].(*panel).addCube(x, y, z, float32(cubeSize))
					} else if cx == 0 && x < y && x < z {
						tr.bits[1].(*panel).addCube(x, y, z, float32(cubeSize))
					} else if cy == tr.lvl && y > x && y > z {
						tr.bits[2].(*panel).addCube(x, y, z, float32(cubeSize))
					} else if cy == 0 && y < x && y < z {
						tr.bits[3].(*panel).addCube(x, y, z, float32(cubeSize))
					} else if cz == tr.lvl && z > x && z > y {
						tr.bits[4].(*panel).addCube(x, y, z, float32(cubeSize))
					} else if cz == 0 && z < x && z < y {
						tr.bits[5].(*panel).addCube(x, y, z, float32(cubeSize))
					}
				}
				if newCells > 0 {
					x, y, z := mx*centerOffset, my*centerOffset, mz*centerOffset
					cube := newCube(eng, tr.part, x, y, z, float32(cubeSize))
					cube.edgeSort(newCells)
					tr.bits = append(tr.bits, cube)
				}
				mz += 2
			}
			my += 2
		}
		mx += 2
	}
	tr.addCenter()
	return tr
}

// The interior center of the trooper is a single cube the size of the previous level.
// This will be nothing on the first level.
func (tr *trooper) addCenter() {
	if tr.lvl > 0 {
		tr.center = tr.part.AddPart()
		tr.center.SetFacade("cube", "flat", "red")
		cubeSize := 1.0 / float32(tr.lvl+1)
		scale := float32(tr.lvl-1) * cubeSize * 0.9 // leave a gap.
		tr.center.SetScale(scale, scale, scale)
	}
}

// health returns the number of cells in the troopers outer layer.
func (tr *trooper) health() int {
	ccnt := 0
	for _, b := range tr.bits {
		ccnt += b.box().ccnt
	}
	return ccnt
}

// attach currently tries to fill in panels first.
func (tr *trooper) attach() {
	for _, b := range tr.bits {
		if b.attach() {
			return
		}
	}
	tr.evolve()
}

// detach currently tries to remove from edges first.
func (tr *trooper) detach() {
	if tr.neo != nil {
		tr.demerge()
		return
	}
	for _, b := range tr.bits {
		if b.detach() {
			return
		}
	}
	tr.devolve()
}

func (tr *trooper) merge() {
	tr.trash()
	tr.neo = tr.part.AddPart()
	tr.neo.SetFacade("cube", "flat", "blue")
	tr.addCenter()
}

func (tr *trooper) demerge() {
	tr.trash()
	tr.addCenter()
	for _, b := range tr.bits {
		b.reset(b.box().cmax)
	}
	tr.bits[0].detach()
}

func (tr *trooper) trash() {
	for _, b := range tr.bits {
		b.trash()
	}
	if tr.center != nil {
		tr.part.RemPart(tr.center)
		tr.center = nil
	}
	tr.part.RemPart(tr.neo)
	tr.neo = nil
}

func (tr *trooper) evolve() {
	if tr.neo == nil {
		// trooper evolved - should be replaced by trooper at next level
		tr.merge()
	}
}

func (tr *trooper) devolve() {
	// trooper devolved - should be replaced by trooper at previous level
}

// ===========================================================================

type box interface {
	attach() bool
	detach() bool
	trash()
	merge()
	reset(count int)
	box() *cbox
}

// cbox is a base class for panels and cubes. It just pulls some common code
// into one spot to remove duplication.
type cbox struct {
	ccnt, cmax     int     // number of cells.
	cx, cy, cz     float32 // center of the box.
	csize          float32 // cell size where each side is the same dimension.
	trashc, mergec func()  // set by super class.
	addc, remc     func()  // set by super class.
}

// attach adds a cell to the cube, merging the cube when the cube is full.
// Attach returns true if a cell was added. A return of false indicates a
// full cube.
func (c *cbox) attach() bool {
	if c.ccnt >= 0 && c.ccnt < c.cmax {
		c.ccnt++ // only spot where this is incremented.
		if c.ccnt == c.cmax {
			c.mergec() // c.merge()
		} else {
			c.addc() // c.addCell()
		}
		return true
	}
	return false
}

// detach removes a cell from the cube, demerging a full cubes if necessary.
// Detach returns true if a cell was detached.  A return of false indicates
// an empty cube.
func (c *cbox) detach() bool {
	if c.ccnt > 0 && c.ccnt <= c.cmax {
		if c.ccnt == c.cmax {
			c.reset(c.cmax - 1)
		} else {
			c.remc() // c.removeCell()
			c.ccnt-- // only spot where this is decremented.
		}
		return true
	}
	return false
}

// reset clears the cbox and ensures the cell count is the given value.
func (c *cbox) reset(cellCount int) {
	c.trashc()
	c.ccnt = 0 // only spot where this is reset to 0
	if cellCount > c.cmax {
		cellCount = c.cmax
	}
	for cnt := 0; cnt < cellCount; cnt++ {
		c.attach()
	}
}

// box allows direct access to the cbox from a super class.
func (c *cbox) box() *cbox { return c }

// ===========================================================================

// panels group 0 or more cubes into the center of one of the troopers
// six sides.
type panel struct {
	eng   *vu.Eng // needed to create new cells.
	part  vu.Part // each panel is its own part.
	lvl   int     // used to scale slab.
	slab  vu.Part // un-injured panel is a single piece.
	cubes []*cube // injured panels are made of cubes.
	cbox
}

// newPanel creates a panel with no cubes.  The cubes are added later using
// panel.addCube().
func newPanel(eng *vu.Eng, part vu.Part, x, y, z float32, level int) *panel {
	p := &panel{}
	p.eng = eng
	p.part = part.AddPart()
	p.lvl = level
	p.cubes = []*cube{}
	p.cx, p.cy, p.cz = x, y, z
	p.ccnt, p.cmax = 0, (level-1)*(level-1)*8
	p.mergec = func() { p.merge() }
	p.trashc = func() { p.trash() }
	p.addc = func() { p.addCell() }
	p.remc = func() { p.removeCell() }
	return p
}

// addCube is only used at the begining to add cubes that are owned by this
// panel.
func (p *panel) addCube(x, y, z, cubeSize float32) {
	p.csize = cubeSize
	c := newCube(p.eng, p.part, x, y, z, p.csize)
	if (p.cx > p.cy && p.cx > p.cz) || (p.cx < p.cy && p.cx < p.cz) {
		c.panelSort(1, 0, 0, 4)
	} else if (p.cy > p.cx && p.cy > p.cz) || (p.cy < p.cx && p.cy < p.cz) {
		c.panelSort(0, 1, 0, 4)
	} else if (p.cz > p.cx && p.cz > p.cy) || (p.cz < p.cx && p.cz < p.cy) {
		c.panelSort(0, 0, 1, 4)
	}
	if c != nil {
		p.ccnt += 4
		p.cubes = append(p.cubes, c)
	}
}

func (p *panel) addCell() {
	for addeven := 0; addeven < p.cubes[0].cmax; addeven++ {
		for _, c := range p.cubes {
			if c.ccnt <= addeven {
				c.attach()
				return
			}
		}
	}
	log.Printf("sg:panel addCell should never reach here. %d %d", p.ccnt, p.cmax)
}

func (p *panel) removeCell() {
	for _, c := range p.cubes {
		if c.detach() {
			return
		}
	}
	log.Printf("sg:panel removeCell should never reach here.")
}

// merge turns all the cubes into a single slab.
func (p *panel) merge() {
	p.trash()
	p.slab = p.part.AddPart()
	p.slab.SetFacade("cube", "flat", "blue")
	p.slab.SetLocation(p.cx, p.cy, p.cz)
	scale := float32(p.lvl-1) * p.csize
	if (p.cx > p.cy && p.cx > p.cz) || (p.cx < p.cy && p.cx < p.cz) {
		p.slab.SetScale(p.csize, scale, scale)
	} else if (p.cy > p.cx && p.cy > p.cz) || (p.cy < p.cx && p.cy < p.cz) {
		p.slab.SetScale(scale, p.csize, scale)
	} else if (p.cz > p.cx && p.cz > p.cy) || (p.cz < p.cx && p.cz < p.cy) {
		p.slab.SetScale(scale, scale, p.csize)
	}
}

// trash clears any visible parts from the panel. It is up to calling methods
// to ensure the cell count is correct.
func (p *panel) trash() {
	if p.slab != nil {
		p.part.RemPart(p.slab)
		p.slab = nil
	}
	for _, cube := range p.cubes {
		cube.reset(0)
	}
}

// ===========================================================================

// cube is the building blocks for troopers and panels.  Cube takes a size
// and location and creates an 8 part cube out of it.  Cubes can be queried
// as to their current number of cells which is between 0 (nothing visible),
// 1-7 (partial) and 8 (merged).
type cube struct {
	eng     *vu.Eng   // needed to create new cells.
	part    vu.Part   // each cube is its own set.
	cells   []vu.Part // max 8 cells per cube.
	centers csort     // precalculated center location of each cell.
	cbox
}

// newCube's are often started with 1 corner, 2 edges, or 4 bottom side pieces.
func newCube(eng *vu.Eng, part vu.Part, x, y, z, cubeSize float32) *cube {
	c := &cube{}
	c.eng = eng
	c.part = part.AddPart()
	c.cells = []vu.Part{}
	c.cx, c.cy, c.cz, c.csize = x, y, z, cubeSize
	c.ccnt, c.cmax = 0, 8
	c.mergec = func() { c.merge() }
	c.trashc = func() { c.trash() }
	c.addc = func() { c.addCell() }
	c.remc = func() { c.removeCell() }

	// calculate the cell center locations (unsorted)
	qs := c.csize * 0.25
	c.centers = csort{
		&lin.V3{x - qs, y - qs, z - qs},
		&lin.V3{x - qs, y - qs, z + qs},
		&lin.V3{x - qs, y + qs, z - qs},
		&lin.V3{x - qs, y + qs, z + qs},
		&lin.V3{x + qs, y - qs, z - qs},
		&lin.V3{x + qs, y - qs, z + qs},
		&lin.V3{x + qs, y + qs, z - qs},
		&lin.V3{x + qs, y + qs, z + qs},
	}
	return c
}

func (c *cube) edgeSort(startCount int) {
	sort.Sort(c.centers)
	c.reset(startCount)
}

func (c *cube) panelSort(rx, ry, rz float32, startCount int) {
	sorter := &ssort{c.centers, rx, ry, rz}
	sort.Sort(sorter)
	c.reset(startCount)
}

// addCell creates and adds a new cell to the cube.
func (c *cube) addCell() {
	cell := c.part.AddPart()
	cell.SetFacade("cube", "flat", "green")
	center := c.centers[c.ccnt-1]
	cell.SetLocation(center.X, center.Y, center.Z)
	scale := c.csize * 0.40 // leave a gap (0.5 for no gap).
	cell.SetScale(scale, scale, scale)
	c.cells = append(c.cells, cell)
}

// removeCell removes the last cell from the list of cube cells.
func (c *cube) removeCell() {
	last := len(c.cells)
	c.part.RemPart(c.cells[last-1])
	c.cells = c.cells[:last-1]
}

// merge removes all cells and replaces them with a single cube. Expected
// to only be called by attach.  The c.ccnt should be c.cmax before and after
// merge is called.
func (c *cube) merge() {
	c.trash()
	cell := c.part.AddPart()
	cell.SetFacade("cube", "flat", "green")
	cell.SetLocation(c.cx, c.cy, c.cz)
	scale := c.csize - (c.csize * 0.15) // leave a gap (just c.csize for no gap)
	cell.SetScale(scale, scale, scale)
	c.cells = append(c.cells, cell)
	// TODO show merge animation.
}

// removes all visible cube parts.
func (c *cube) trash() {
	for len(c.cells) > 0 {
		c.removeCell()
	}
}

// csort is used to sort the cube quadrants so that the quadrants closest
// to the origin are first in the list.  This way the cells added first and
// removed last are those closest to the center.
//
// A reference point is necessary since the origin gets to far away for
// a flat panel to orient the quads properly.
type csort []*lin.V3 // list of quadrant centers.

func (c csort) Len() int               { return len(c) }
func (c csort) Swap(i, j int)          { c[i], c[j] = c[j], c[i] }
func (c csort) Less(i, j int) bool     { return c.Dtoc(c[i]) < c.Dtoc(c[j]) }
func (c csort) Dtoc(v *lin.V3) float32 { return v.X*v.X + v.Y*v.Y + v.Z*v.Z }

// ssort is used to sort the panel cube quadrants so that the quadrants
// to the inside origin plane are first in the list. A reference normal is
// necessary since the panels get large enough that the points on the
// "outside" get picked up due to the angle.
type ssort struct {
	c       []*lin.V3 // list of quadrant centers.
	x, y, z float32   // reference plane.
}

func (s ssort) Len() int           { return len(s.c) }
func (s ssort) Swap(i, j int)      { s.c[i], s.c[j] = s.c[j], s.c[i] }
func (s ssort) Less(i, j int) bool { return s.Dtoc(s.c[i]) < s.Dtoc(s.c[j]) }
func (s ssort) Dtoc(v *lin.V3) float32 {
	normal := &lin.V3{s.x, s.y, s.z}
	dot := v.Dot(normal)
	dx := normal.X * dot
	dy := normal.Y * dot
	dz := normal.Z * dot
	return dx*dx + dy*dy + dz*dz
}
