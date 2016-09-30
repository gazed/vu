// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/gazed/vu"
	"github.com/gazed/vu/grid"
	"github.com/gazed/vu/math/lin"
)

// hx demonstrates hexagonal grids. Its main purpose is to determine
// which code can be provided generically by the engine (vu/grid.Hex)
// and which code is application specific (hexGrid, hexTile).
//
// CONTROLS:
//   Tab   : flip orientation
//   Lm    : show mouse hit
//   WASD  : move hex grid          : up left down right
//   ZX    : scale hex grid         : bigger smaller
func hx() {
	hx := &hxtag{}
	if err := vu.New(hx, "Hex Grid", 400, 100, 800, 600); err != nil {
		log.Printf("hx: error starting engine %s", err)
	}
	defer catchErrors()
}

// Encapsulate example specific data with a unique "tag".
type hxtag struct {
	cam  *vu.Camera // Fixed camera.
	hg   *hexGrid   // scales and positions the hexes.
	flat bool       // flat or pointy grid orientation. Start off pointy.
}

// Create is the engine callback for initial asset creation.
func (hx *hxtag) Create(eng vu.Eng, s *vu.State) {
	hx.cam = eng.Root().NewCam().SetUI()
	hx.cam.SetOrthographic(0, float64(s.W), 0, float64(s.H), 0, 50)

	// center hex tile.
	hx.hg = newHexGrid(eng.Root().NewPov())
	hx.hg.newTile(0, 0) //   Q  R  S  = 0  0  0

	// first ring around center: 6 hexes.
	hx.hg.newTile(-1, 0) // -1  0  1
	hx.hg.newTile(-1, 1) // -1  1  0
	hx.hg.newTile(1, -1) //  1 -1  0
	hx.hg.newTile(1, 0)  //  1  0 -1
	hx.hg.newTile(0, -1) //  0 -1  1
	hx.hg.newTile(0, 1)  //  0  1 -1

	// second ring around center: 12 hexes.
	hx.hg.newTile(2, 0)   //  2  0 -2
	hx.hg.newTile(1, 1)   //  1  1 -2
	hx.hg.newTile(0, 2)   //  0  2 -2
	hx.hg.newTile(-1, 2)  // -1  2 -1
	hx.hg.newTile(-2, 2)  // -2  2  0
	hx.hg.newTile(-2, 1)  // -2  1  1
	hx.hg.newTile(-2, 0)  // -2  0  2
	hx.hg.newTile(-1, -1) // -1 -1  2
	hx.hg.newTile(0, -2)  //  0 -2  2
	hx.hg.newTile(1, -2)  //  1 -2  1
	hx.hg.newTile(2, -2)  //  2 -2  0
	hx.hg.newTile(2, -1)  //  2 -1 -1

	// create the hilite marker last so it appears on top.
	hx.hg.hilite = hx.hg.models.NewPov()
	hx.hg.hilite.NewModel("uv", "msh:icon", "tex:halo")
	hx.hg.hilite.Cull = true
}

// Update is the regular engine callback.
func (hx *hxtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		hx.cam.SetOrthographic(0, float64(s.W), 0, float64(s.H), 0, 50)
	}
	hx.hg.updateLabels(in.Ut) // handle grid changes from last update.

	// scale and move the board.
	for press, down := range in.Down {
		switch {
		case press == vu.KZ:
			scale := hx.hg.scale()
			scale += 1.0
			hx.hg.setScale(scale)
		case press == vu.KX:
			scale := hx.hg.scale()
			scale -= 1.0
			hx.hg.setScale(scale)
		case press == vu.KW:
			hx.hg.move(0, 1)
		case press == vu.KS:
			hx.hg.move(0, -1)
		case press == vu.KA:
			hx.hg.move(-1, 0)
		case press == vu.KD:
			hx.hg.move(1, 0)
		case press == vu.KTab && down == 1:
			hx.hg.flipOrientation()
		case press == vu.KLm && down == 1:
			t := hx.hg.hit(in.Mx, in.My)
			hx.hg.mark(t)
		}
	}
	hx.hg.spinMark()
}

// hx example
// =============================================================================
// hexGrid playing surface.

// hexGrid represents a playing surface composed of hex tiles.
type hexGrid struct {
	labels *vu.Pov // root for the hex labels.
	models *vu.Pov // root for the hex models.
	hilite *vu.Pov // hilite the selected hex tile.
	flat   bool    // flat or pointy grid orientation. Start off pointy.

	// kludge to handle the fact that labels lag behind other model
	// changes. The hex model world positions are updated after the
	// call to Update()
	labelsDirty bool // set to true to update label positions.

	// tiles is indexed by the unique hexTile id.
	tiles map[uint64]*hexTile // all hex tiles in the grid.
}

// newHexGrid is expected to be called once on startup to create
// a single hexGrid instance.
func newHexGrid(root *vu.Pov) *hexGrid {
	hg := &hexGrid{labels: root}
	hg.models = root.NewPov().SetAt(400, 300, 0).SetScale(128, 128, 0)
	hg.tiles = map[uint64]*hexTile{}
	hg.labelsDirty = true
	return hg
}

// newTile creates a single hex tile instance.
func (hg *hexGrid) newTile(q, r int) *hexTile {
	hex := newHexTile(hg.labels, hg.models, int32(q), int32(r))
	hg.tiles[hex.id()] = hex
	return hex
}

// scale returns the current size of the hex tiles.
func (hg *hexGrid) scale() float64 {
	s, _, _ := hg.models.Scale()
	return s
}

// setScale changes the size of the hex tiles.
// The scale is set to the given value s.
func (hg *hexGrid) setScale(s float64) {
	hg.models.SetScale(s, s, 0)
	hg.labelsDirty = true
}

// move repositions the grid of hexes. The parameters are
// relative to the current location.
func (hg *hexGrid) move(x, y int) {
	bx, by, _ := hg.models.At()
	hg.models.SetAt(bx+float64(x), by+float64(y), 0)
	hg.labelsDirty = true
}

// flipOrientation changes from flat to pointy and back.
func (hg *hexGrid) flipOrientation() {
	hg.flat = !hg.flat
	size := 0.52 // greater than 0.5 for a gap between hexes.
	if hg.flat {
		for _, t := range hg.tiles {
			hx, hy := t.hex.ToFlat(size)
			t.model.SetAt(hx, hy, 0)
			t.model.Spin(0, 0, 30)
		}
	} else {
		for _, t := range hg.tiles {
			hx, hy := t.hex.ToPointy(size)
			t.model.SetAt(hx, hy, 0)
			t.model.Spin(0, 0, -30)
		}
	}
	hg.labelsDirty = true
	hg.hilite.Cull = true
}

// updateLabels is a kludge to update labels in the Update
// call after the Update where hex models change location.
func (hg *hexGrid) updateLabels(ut uint64) {
	if hg.labelsDirty && ut > 1 {
		for _, tile := range hg.tiles {
			tile.updateLabel()
			hg.labelsDirty = false
		}
	}
}

// mark updates the hilite marker to be on tile t.
// The highlite mark is turned off if t is nil.
func (hg *hexGrid) mark(t *hexTile) {
	if t != nil {
		x, y, z := t.model.At()
		hg.hilite.SetAt(x, y, z+0.1)
		hg.hilite.Cull = false
	} else {
		hg.hilite.Cull = true
	}
}

// spinMark adds some animation to the currently selected tile.
func (hg *hexGrid) spinMark() {
	if hg.hilite != nil && !hg.hilite.Cull {
		hg.hilite.Spin(0, 0, 2)
	}
}

// hit returns a hex tile if the mouse click was in a hex grid.
// This cheats by approimating the hex grid with a circle by using
// the innerRadius which is the distance from the center to an edge
// center. Clicking near a hex point will not cause a hit.
func (hg *hexGrid) hit(mx, my int) (t *hexTile) {
	size := hg.scale() * 0.5 // half the size of hex icon.
	innerRadius := size * math.Sin(lin.Rad(60)) / math.Sin(lin.Rad(90))
	radiusSquared := innerRadius * innerRadius
	for _, tile := range hg.tiles {
		tx, ty, _ := tile.model.World()
		dx, dy := float64(mx)-tx, float64(my)-ty
		dist2 := dx*dx + dy*dy
		if dist2 < radiusSquared {
			return tile
		}
	}
	return nil
}

// FUTURE: replace above hit test with exact hit. See:
//   http://www.playchilla.com/how-to-check-if-a-point-is-inside-a-hexagon
//   public function isInside(pos:Vec2Const):Boolean {
//    const q2x:Number = Math.abs(pos.x - _center.x);  // transform the test point locally and to quadrant 2
//    const q2y:Number = Math.abs(pos.y - _center.y);  // transform the test point locally and to quadrant 2
//    if (q2x > _hori || q2y > _vert*2) return false;  // bounding test (since q2 is in quadrant 2 only 2 tests are needed)
//
//    // finally the dot product can be reduced to this due to the hexagon symmetry
//    return 2 * _vert * _hori - _vert * q2x - _hori * q2y >= 0;
//   }
//
// Also consider/see:
//   http://www.redblobgames.com/grids/hexagons/implementation.html
//   FractionalHex pixel_to_hex(Layout layout, Point p) {
//     const Orientation& M = layout.orientation; // flat or pointy layout.
//     Point pt = Point((p.x - layout.origin.x) / layout.size.x,
//                      (p.y - layout.origin.y) / layout.size.y);
//     double q = M.b0 * pt.x + M.b1 * pt.y;
//     double r = M.b2 * pt.x + M.b3 * pt.y;
//     return FractionalHex(q, r, -q - r);
//   }

// hexGrid playing surface.
// =============================================================================
// hexTile for a hexGrid.

// hexTile represents a hex on the screen.
type hexTile struct {
	hex   *grid.Hex // location in hex cubic grid coordinates.
	label *vu.Pov   // optional tile label centered on tile.
	model *vu.Pov   // tile model with texture.
}

// newHexTile creates a single hex for the given grid location.
func newHexTile(root, board *vu.Pov, q, r int32) *hexTile {
	t := &hexTile{}
	t.hex = grid.NewHex(q, r)

	// A hex image is in a square that overlaps adjacent squares.
	hx, hy := t.hex.ToPointy(0.52) // greater than 0.5 for a gap between hexes.
	t.model = board.NewPov().SetAt(hx, hy, 0)
	t.model.NewModel("uv", "msh:icon", "tex:hextile")

	// labels use a different parent pov so they are not scaled with the board
	font := "lucidiaSu22"
	t.label = root.NewPov()
	label := t.label.NewLabel("uv", font, font+"White")
	label.SetStr(fmt.Sprintf("%d %d %d", t.hex.Q, t.hex.R, t.hex.S))
	return t
}

// id returns a unique identifier for this hex tile.
func (t *hexTile) id() uint64 { return t.hex.ID() }

// updateLabel ensures that the labels use the same world position
// as their associated hex model.
func (t *hexTile) updateLabel() {
	wx, wy, _ := t.model.World()
	textWidth := float64(t.label.Model().StrWidth())
	t.label.SetAt(wx-textWidth*0.5, wy-11, 0)
}
