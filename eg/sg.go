// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"
	"sort"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// sg tests the scene graph by building up a movable character that has multiple
// layers of sub-parts. The scene graph is demonstrated by changing the top level
// location, orientation and having it affect the sub-parts. Sg also tests adding
// and removing parts from a scene graph. Note that transparency sorting is
// handled automatically by the engine.
// This example has a bit more code due to playing around with what can best
// be described as merging and splitting voxels.
//
// CONTROLS:
//   WASD  : move model             : forward left back right
//   =-    : adjust boxes           : attach detach
//   0-5   : change cube size
//   P     : print number of cubes
func sg() {
	defer catchErrors()
	if err := vu.Run(&sgtag{}); err != nil {
		log.Printf("sg: error starting engine %s", err)
	}
}

// Globally unique "tag" that encapsulates example specific data.
type sgtag struct {
	scene  *vu.Ent
	meld   *meld
	run    float64
	spin   float64
	dt     float64
	reacts map[int]inputHandler
}

// inputHandler helps link keypresses to handling functions.
type inputHandler func(in *vu.Input, down int)

// Create is the engine callback for initial asset creation.
func (sg *sgtag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Scene Graph"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.1, 0.1, 0.1, 1.0))
	sg.scene = eng.AddScene()
	sg.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(0, 0, 6)
	sg.run = 10   // move so many cubes worth in one second.
	sg.spin = 270 // spin so many degrees in one second.
	sg.meld = newMeld(sg.scene.AddPart(), 1)

	// initialize the reactions
	sg.reacts = map[int]inputHandler{
		vu.KW:     sg.forward,
		vu.KA:     sg.left,
		vu.KS:     sg.back,
		vu.KD:     sg.right,
		vu.KEqual: sg.attach,
		vu.KMinus: sg.detach,
		vu.K0:     func(i *vu.Input, down int) { sg.setTr(down, 0) },
		vu.K1:     func(i *vu.Input, down int) { sg.setTr(down, 1) },
		vu.K2:     func(i *vu.Input, down int) { sg.setTr(down, 2) },
		vu.K3:     func(i *vu.Input, down int) { sg.setTr(down, 3) },
		vu.K4:     func(i *vu.Input, down int) { sg.setTr(down, 4) },
		vu.K5:     func(i *vu.Input, down int) { sg.setTr(down, 5) },
		vu.KP:     sg.stats,
	}
}
func (sg *sgtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	sg.dt = in.Dt
	for press, downLength := range in.Down {
		if react, ok := sg.reacts[press]; ok {
			react(in, downLength)
		}
	}
}

// User actions.
func (sg *sgtag) stats(i *vu.Input, down int) {
	if down == 1 {
		log.Printf("Cubes %d", sg.meld.health())
	}
}
func (sg *sgtag) left(i *vu.Input, down int)  { sg.meld.top.Spin(0, sg.dt*sg.spin, 0) }
func (sg *sgtag) right(i *vu.Input, down int) { sg.meld.top.Spin(0, sg.dt*-sg.spin, 0) }
func (sg *sgtag) back(i *vu.Input, down int) {
	sg.meld.top.Move(0, 0, sg.dt*sg.run, sg.scene.Cam().Look)
}
func (sg *sgtag) forward(i *vu.Input, down int) {
	sg.meld.top.Move(0, 0, sg.dt*-sg.run, sg.scene.Cam().Look)
}
func (sg *sgtag) attach(i *vu.Input, down int) { sg.meld.attach() }
func (sg *sgtag) detach(i *vu.Input, down int) { sg.meld.detach() }
func (sg *sgtag) setTr(down, lvl int) {
	if down == 1 {
		sg.meld.trash()
		sg.meld = newMeld(sg.scene.AddPart(), lvl)
	}
}

// =============================================================================

// meld plays around with a composite object made up of small pieces that
// are combined or separated as bits are added or removed.
type meld struct {
	top    *vu.Ent
	neo    *vu.Ent // un-injured meld
	center *vu.Ent // center always represented as one piece
	bits   []box   // panels and edge cubes.
	lvl    int
	mid    int // level entry number of cells.
}

// newMeld creates a meld at the starting size for the given level.
//    level 0: 1x1x1 :  1 cube
//    level 1: 2x2x2 :  8 edge cubes + 6 panels of 0x0 cubes + 0x0x0 center.
//    level 2: 3x3x3 : 20 edge cubes + 6 panels of 1x1 cubes + 1x1x1 center.
//    level 3: 4x4x4 : 32 edge cubes + 6 panels of 2x2 cubes + 2x2x2 center.
//    ...
func newMeld(pov *vu.Ent, level int) *meld {
	m := &meld{top: pov}
	m.lvl = level
	m.bits = []box{}
	m.mid = m.lvl*m.lvl*m.lvl*8 - (m.lvl-1)*(m.lvl-1)*(m.lvl-1)*8

	//
	if m.lvl == 0 {
		cube := newCube(m.top, 0, 0, 0, 1)
		cube.edgeSort(1)
		m.bits = append(m.bits, cube)
		return m
	}

	// create the panels. These are used in each level but the first.
	cubeSize := 1.0 / float64(m.lvl+1)
	centerOffset := cubeSize * 0.5
	panelCenter := float64(m.lvl) * centerOffset
	m.bits = append(m.bits, newBlock(m.top, panelCenter, 0.0, 0.0, m.lvl))
	m.bits = append(m.bits, newBlock(m.top, -panelCenter, 0.0, 0.0, m.lvl))
	m.bits = append(m.bits, newBlock(m.top, 0.0, panelCenter, 0.0, m.lvl))
	m.bits = append(m.bits, newBlock(m.top, 0.0, -panelCenter, 0.0, m.lvl))
	m.bits = append(m.bits, newBlock(m.top, 0.0, 0.0, panelCenter, m.lvl))
	m.bits = append(m.bits, newBlock(m.top, 0.0, 0.0, -panelCenter, m.lvl))

	// meld are made out of cubes and panels.
	mx := float64(-m.lvl)
	for cx := 0; cx <= m.lvl; cx++ {
		my := float64(-m.lvl)
		for cy := 0; cy <= m.lvl; cy++ {
			mz := float64(-m.lvl)
			for cz := 0; cz <= m.lvl; cz++ {

				// create the outer edges.
				newCells := 0
				if (cx == 0 || cx == m.lvl) && (cy == 0 || cy == m.lvl) && (cz == 0 || cz == m.lvl) {

					// corner cube
					newCells = 1
				} else if (cx == 0 || cx == m.lvl) && (cy == 0 || cy == m.lvl) ||
					(cx == 0 || cx == m.lvl) && (cz == 0 || cz == m.lvl) ||
					(cy == 0 || cy == m.lvl) && (cz == 0 || cz == m.lvl) {

					// edge cube
					newCells = 2
				} else if cx == 0 || cx == m.lvl || cy == 0 || cy == m.lvl || cz == 0 || cz == m.lvl {

					// side cubes are added to (controlled by) a panel.
					x, y, z := mx*centerOffset, my*centerOffset, mz*centerOffset
					if cx == m.lvl && x > y && x > z {
						m.bits[0].(*block).addCube(x, y, z, float64(cubeSize))
					} else if cx == 0 && x < y && x < z {
						m.bits[1].(*block).addCube(x, y, z, float64(cubeSize))
					} else if cy == m.lvl && y > x && y > z {
						m.bits[2].(*block).addCube(x, y, z, float64(cubeSize))
					} else if cy == 0 && y < x && y < z {
						m.bits[3].(*block).addCube(x, y, z, float64(cubeSize))
					} else if cz == m.lvl && z > x && z > y {
						m.bits[4].(*block).addCube(x, y, z, float64(cubeSize))
					} else if cz == 0 && z < x && z < y {
						m.bits[5].(*block).addCube(x, y, z, float64(cubeSize))
					}
				}
				if newCells > 0 {
					x, y, z := mx*centerOffset, my*centerOffset, mz*centerOffset
					cube := newCube(m.top, x, y, z, float64(cubeSize))
					cube.edgeSort(newCells)
					m.bits = append(m.bits, cube)
				}
				mz += 2
			}
			my += 2
		}
		mx += 2
	}
	m.addCenter()
	return m
}

// The interior center of the meld is a single cube the size of the previous level.
// This will be nothing on the first level.
func (m *meld) addCenter() {
	if m.lvl > 0 {
		cubeSize := 1.0 / float64(m.lvl+1)
		scale := float64(m.lvl-1) * cubeSize * 0.9 // leave a gap.
		m.center = m.top.AddPart().SetScale(scale, scale, scale)
		m.center.MakeModel("colored", "msh:box", "mat:red").SetAlpha(0.5)
	}
}

// health returns the number of cells in the meld outer layer.
func (m *meld) health() int {
	ccnt := 0
	for _, b := range m.bits {
		ccnt += b.box().ccnt
	}
	return ccnt
}

// attach tries to fill in panels first.
func (m *meld) attach() {
	for _, b := range m.bits {
		if b.attach() {
			return
		}
	}
	m.evolve()
}

// detach tries to remove from edges first.
func (m *meld) detach() {
	if m.neo != nil {
		m.demerge()
		return
	}
	for _, b := range m.bits {
		if b.detach() {
			return
		}
	}
	m.devolve()
}

func (m *meld) merge() {
	m.trash()
	m.neo = m.top.AddPart().MakeModel("colored", "msh:box", "mat:blue")
	m.addCenter()
}

func (m *meld) demerge() {
	m.trash()
	m.addCenter()
	for _, b := range m.bits {
		b.reset(b.box().cmax)
	}
	m.bits[0].detach()
}

func (m *meld) trash() {
	for _, b := range m.bits {
		b.trash()
	}
	if m.center != nil {
		m.center.Dispose()
		m.center = nil
	}
	if m.neo != nil {
		m.neo.Dispose()
	}
	m.neo = nil
}

func (m *meld) evolve() {
	if m.neo == nil {
		// meld evolved - should be replaced by meld at next level
		m.merge()
	}
}

// meld devolved - should be replaced by meld at previous level
func (m *meld) devolve() {}

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
	cx, cy, cz     float64 // center of the box.
	csize          float64 // cell size where each side is the same dimension.
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

// blocks group 0 or more cubes into the center of one of the melds
// six sides.
type block struct {
	part  *vu.Ent // each panel is its own part.
	lvl   int     // used to scale slab.
	slab  *vu.Ent // un-injured panel is a single piece.
	cubes []*cube // injured panels are made of cubes.
	cbox
}

// newBlock creates a panel with no cubes.  The cubes are added later using
// panel.addCube().
func newBlock(part *vu.Ent, x, y, z float64, level int) *block {
	b := &block{}
	b.part = part.AddPart()
	b.lvl = level
	b.cubes = []*cube{}
	b.cx, b.cy, b.cz = x, y, z
	b.ccnt, b.cmax = 0, (level-1)*(level-1)*8
	b.mergec = func() { b.merge() }
	b.trashc = func() { b.trash() }
	b.addc = func() { b.addCell() }
	b.remc = func() { b.removeCell() }
	return b
}

// addCube is only used at the beginning to add cubes that are owned by this
// panel.
func (b *block) addCube(x, y, z, cubeSize float64) {
	b.csize = cubeSize
	c := newCube(b.part, x, y, z, b.csize)
	if (b.cx > b.cy && b.cx > b.cz) || (b.cx < b.cy && b.cx < b.cz) {
		c.panelSort(1, 0, 0, 4)
	} else if (b.cy > b.cx && b.cy > b.cz) || (b.cy < b.cx && b.cy < b.cz) {
		c.panelSort(0, 1, 0, 4)
	} else if (b.cz > b.cx && b.cz > b.cy) || (b.cz < b.cx && b.cz < b.cy) {
		c.panelSort(0, 0, 1, 4)
	}
	if c != nil {
		b.ccnt += 4
		b.cubes = append(b.cubes, c)
	}
}

func (b *block) addCell() {
	for addeven := 0; addeven < b.cubes[0].cmax; addeven++ {
		for _, c := range b.cubes {
			if c.ccnt <= addeven {
				c.attach()
				return
			}
		}
	}
	log.Printf("sg:panel addCell should never reach here. %d %d", b.ccnt, b.cmax)
}

func (b *block) removeCell() {
	for _, c := range b.cubes {
		if c.detach() {
			return
		}
	}
	log.Printf("sg:panel removeCell should never reach here.")
}

// merge turns all the cubes into a single slab.
func (b *block) merge() {
	b.trash()
	b.slab = b.part.AddPart().SetAt(b.cx, b.cy, b.cz)
	b.slab.MakeModel("colored", "msh:box", "mat:blue")
	scale := float64(b.lvl-1) * b.csize
	if (b.cx > b.cy && b.cx > b.cz) || (b.cx < b.cy && b.cx < b.cz) {
		b.slab.SetScale(b.csize, scale, scale)
	} else if (b.cy > b.cx && b.cy > b.cz) || (b.cy < b.cx && b.cy < b.cz) {
		b.slab.SetScale(scale, b.csize, scale)
	} else if (b.cz > b.cx && b.cz > b.cy) || (b.cz < b.cx && b.cz < b.cy) {
		b.slab.SetScale(scale, scale, b.csize)
	}
}

// trash clears any visible parts from the panel. It is up to calling methods
// to ensure the cell count is correct.
func (b *block) trash() {
	if b.slab != nil {
		b.slab.Dispose()
		b.slab = nil
	}
	for _, cube := range b.cubes {
		cube.reset(0)
	}
}

// ===========================================================================

// cube is the building blocks for meld and panels.  Cube takes a size
// and location and creates an 8 part cube out of it.  Cubes can be queried
// as to their current number of cells which is between 0 (nothing visible),
// 1-7 (partial) and 8 (merged).
type cube struct {
	part    *vu.Ent   // each cube is its own set.
	cells   []*vu.Ent // max 8 cells per cube.
	centers csort     // precalculated center location of each cell.
	cbox
}

// newCube's are often started with 1 corner, 2 edges, or 4 bottom side pieces.
func newCube(meld *vu.Ent, x, y, z, cubeSize float64) *cube {
	c := &cube{}
	c.part = meld.AddPart()
	c.cells = []*vu.Ent{}
	c.cx, c.cy, c.cz, c.csize = x, y, z, cubeSize
	c.ccnt, c.cmax = 0, 8
	c.mergec = func() { c.merge() }
	c.trashc = func() { c.trash() }
	c.addc = func() { c.addCell() }
	c.remc = func() { c.removeCell() }

	// calculate the cell center locations (unsorted)
	qs := c.csize * 0.25
	c.centers = csort{
		&lin.V3{X: x - qs, Y: y - qs, Z: z - qs},
		&lin.V3{X: x - qs, Y: y - qs, Z: z + qs},
		&lin.V3{X: x - qs, Y: y + qs, Z: z - qs},
		&lin.V3{X: x - qs, Y: y + qs, Z: z + qs},
		&lin.V3{X: x + qs, Y: y - qs, Z: z - qs},
		&lin.V3{X: x + qs, Y: y - qs, Z: z + qs},
		&lin.V3{X: x + qs, Y: y + qs, Z: z - qs},
		&lin.V3{X: x + qs, Y: y + qs, Z: z + qs},
	}
	return c
}

func (c *cube) edgeSort(startCount int) {
	sort.Sort(c.centers)
	c.reset(startCount)
}

func (c *cube) panelSort(rx, ry, rz float64, startCount int) {
	sorter := &ssort{c.centers, rx, ry, rz}
	sort.Sort(sorter)
	c.reset(startCount)
}

// addCell creates and adds a new cell to the cube.
func (c *cube) addCell() {
	center := c.centers[c.ccnt-1]
	cell := c.part.AddPart().SetAt(center.X, center.Y, center.Z)
	cell.MakeModel("colored", "msh:box", "mat:green")
	scale := c.csize * 0.40 // leave a gap (0.5 for no gap).
	cell.SetScale(scale, scale, scale)
	c.cells = append(c.cells, cell)
}

// removeCell removes the last cell from the list of cube cells.
func (c *cube) removeCell() {
	last := len(c.cells)
	c.cells[last-1].Dispose()
	c.cells = c.cells[:last-1]
}

// merge removes all cells and replaces them with a single cube. Expected
// to only be called by attach.  The c.ccnt should be c.cmax before and after
// merge is called.
func (c *cube) merge() {
	c.trash()
	cell := c.part.AddPart().SetAt(c.cx, c.cy, c.cz)
	cell.MakeModel("colored", "msh:box", "mat:green")
	scale := c.csize - (c.csize * 0.15) // leave a gap (just c.csize for no gap)
	cell.SetScale(scale, scale, scale)
	c.cells = append(c.cells, cell)
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
func (c csort) Dtoc(v *lin.V3) float64 { return v.X*v.X + v.Y*v.Y + v.Z*v.Z }

// ssort is used to sort the panel cube quadrants so that the quadrants
// to the inside origin plane are first in the list. A reference normal is
// necessary since the panels get large enough that the points on the
// "outside" get picked up due to the angle.
type ssort struct {
	c       []*lin.V3 // list of quadrant centers.
	x, y, z float64   // reference plane.
}

func (s ssort) Len() int           { return len(s.c) }
func (s ssort) Swap(i, j int)      { s.c[i], s.c[j] = s.c[j], s.c[i] }
func (s ssort) Less(i, j int) bool { return s.Dtoc(s.c[i]) < s.Dtoc(s.c[j]) }
func (s ssort) Dtoc(v *lin.V3) float64 {
	normal := &lin.V3{X: s.x, Y: s.y, Z: s.z}
	dot := v.Dot(normal)
	dx := normal.X * dot
	dy := normal.Y * dot
	dz := normal.Z * dot
	return dx*dx + dy*dy + dz*dz
}
