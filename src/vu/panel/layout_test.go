// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package panel

import (
	"testing"
)

func TestFixedLayout(t *testing.T) {
	s := newSection(200, 100)
	c1 := s.AddControl().(*control)
	c2 := s.AddControl().(*control)
	s1 := s.AddSection().(*section)
	s2 := s.AddSection().(*section)
	s.SetLayout(&FixedLayout{map[uint]*FixedInfo{
		c1.Id(): &FixedInfo{10, 10, 50, 25, BL},
		c2.Id(): &FixedInfo{10, 10, 50, 25, BR},
		s1.Id(): &FixedInfo{10, 10, 50, 25, TL},
		s2.Id(): &FixedInfo{10, 10, 50, 25, TR},
	}})
	s.align()
	c, x, y, w, h := c1, 10, 10, 50, 25
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}
	c, x, y, w, h = c2, 140, 10, 50, 25
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}
	p, x, y, w, h := s1, 10, 65, 50, 25
	if x != p.sx || y != p.sy || w != p.sw || h != p.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, p.sx, p.sy, p.sw, p.sh)
	}
	p, x, y, w, h = s2, 140, 65, 50, 25
	if x != p.sx || y != p.sy || w != p.sw || h != p.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, p.sx, p.sy, p.sw, p.sh)
	}

	// test center locations
	s.layout.(*FixedLayout).Info[c1.Id()].Mode = CL
	s.layout.(*FixedLayout).Info[c2.Id()].Mode = CR
	s.layout.(*FixedLayout).Info[s1.Id()].Mode = CT
	s.layout.(*FixedLayout).Info[s2.Id()].Mode = CB
	s.align()
	c, x, y, w, h = c1, 10, 48, 50, 25
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}
	c, x, y, w, h = c2, 140, 48, 50, 25
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}
	p, x, y, w, h = s1, 85, 65, 50, 25
	if x != p.sx || y != p.sy || w != p.sw || h != p.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, p.sx, p.sy, p.sw, p.sh)
	}
	p, x, y, w, h = s2, 85, 10, 50, 25
	if x != p.sx || y != p.sy || w != p.sw || h != p.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, p.sx, p.sy, p.sw, p.sh)
	}
}

func TestGridLayout(t *testing.T) {
	s := newSection(200, 100)
	s.sx, s.sy, s.sw, s.sh = 0, 0, 200, 100 // fake a parent layout resize.
	s.AddControl()
	s.AddControl()
	s.AddControl()
	c := s.AddControl().(*control)

	// Test with 4 rows of 1 widget.
	s.align()
	x, y, w, h := 0, 75, 200, 25
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}

	// Retest with 1 row of 4 widgets.
	s.layout.(*GridLayout).Columns = 4
	s.align()
	x, y, w, h = 150, 0, 50, 100
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}

	// Retest with 2 rows of 2 widgets.
	s.layout.(*GridLayout).Columns = 2
	s.align()
	x, y, w, h = 100, 50, 100, 50
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}

	// Retest with margin.
	s.layout.(*GridLayout).Margin = 5
	s.align()
	x, y, w, h = 105, 55, 95, 45
	if x != c.sx || y != c.sy || w != c.sw || h != c.sh {
		t.Errorf("Wanted %d %d %d %d got %d %d %d %d", x, y, w, h, c.sx, c.sy, c.sw, c.sh)
	}
}
